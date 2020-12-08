package libepic

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"hash"
	"math"
	"time"

	"github.com/dchest/cmac"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	//"github.com/scionproto/scion/go/lib/log"
	"github.com/scionproto/scion/go/lib/serrors"
	"github.com/scionproto/scion/go/lib/slayers"
	"github.com/scionproto/scion/go/lib/slayers/path/epic"
	"github.com/scionproto/scion/go/lib/slayers/path/scion"
)

const (
	// Error messages
	ErrCipherFailure common.ErrMsg = "Unable to initialize AES cipher"
	ErrMacFailure    common.ErrMsg = "Unable to initialize Mac"
	// Maximal lifetime of a packet in milliseconds
	packetLifetimeMs uint16 = 2000
	// Maximal clock skew in milliseconds
	clockSkewMs uint16 = 1000
)

func initEpicMac(key []byte) (hash.Hash, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, common.NewBasicError(ErrCipherFailure, err)
	}
	// CMAC is not ideal for EPIC due to its subkey generation overhead.
	// We might want to change this in the future.
	mac, err := cmac.New(block)
	if err != nil {
		return nil, common.NewBasicError(ErrMacFailure, err)
	}
	return mac, nil
}

//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   | flags (1B) | timestamp (4B) | packetTimestamp (8B)  |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   | srcIA (8B) | srcAddr (4/8/12/16B) | payloadLen (2B) |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   | zero padding (0-15B)                                |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// The "flags" field only encodes the length of the source address.
func inputToBytes(timestamp uint32, packetTimestamp uint64,
	srcIA addr.IA, srcAddr []byte, srcAddrLen uint8, payloadLen uint16) ([]byte, error) {

	l := int((srcAddrLen + 1) * 4)
	if srcAddrLen > 3 || l != len(srcAddr) {
		return nil, serrors.New("srcAddrLen must be between 0 and 3, and encode the "+
			"srcAddr length", "srcAddrLen", srcAddrLen, "len(srcAddr)", len(srcAddr))
	}

	// Create a multiple of 16 such that the input fits in
	nrBlocks := uint8(math.Ceil((23 + float64(l)) / 16))
	input := make([]byte, 16*nrBlocks)

	// Fill input
	input[0] = srcAddrLen
	binary.BigEndian.PutUint32(input[1:5], timestamp)
	binary.BigEndian.PutUint64(input[5:13], packetTimestamp)
	binary.BigEndian.PutUint64(input[13:21], uint64(srcIA.A))
	binary.BigEndian.PutUint16(input[13:15], uint16(srcIA.I))
	copy(input[21:(21+l)], srcAddr[:l])
	binary.BigEndian.PutUint16(input[(21+l):(23+l)], payloadLen)
	return input, nil
}

func prepareMacInput(epicpath *epic.EpicPath, s *slayers.SCION, timestamp uint32) ([]byte, error) {
	if epicpath == nil {
		return nil, serrors.New("epicpath must not be nil")
	}
	if s == nil {
		return nil, serrors.New("SCION common+address header must not be nil")
	}
	packetTimestamp := epicpath.PacketTimestamp
	payloadLen := s.PayloadLen
	srcIA := s.SrcIA
	srcAddrLen := uint8(s.SrcAddrLen)
	srcAddr := s.RawSrcAddr
	return inputToBytes(timestamp, packetTimestamp, srcIA, srcAddr, srcAddrLen, payloadLen)
}

// todo: add function to create epic timestamp from PckId

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                             TsRel                             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    CoreID     |                  CoreCounter                  |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
func CreateEpicTimestamp(tsRel uint32, coreID uint8, coreCounter uint32) (packetTimestamp uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b[4:8], coreCounter)
	binary.BigEndian.PutUint16(b[3:5], uint16(coreID))
	binary.BigEndian.PutUint32(b[:4], tsRel)
	packetTimestamp = binary.BigEndian.Uint64(b[:8])
	return
}

func ParseEpicTimestamp(packetTimestamp uint64) (tsRel uint32, coreID uint8, coreCounter uint32) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b[:8], packetTimestamp)
	tsRel = binary.BigEndian.Uint32(b[:4])
	coreID = uint8(binary.BigEndian.Uint16(b[3:5]))
	coreCounter = binary.BigEndian.Uint32(b[4:8]) % (1<<24)
	return tsRel, coreID, coreCounter
}

func VerifyTimestamp(timestamp uint32, packetTimestamp uint64) bool {
	// Unix time in milliseconds when the packet was timestamped by the sender
	tsRel, _, _ := ParseEpicTimestamp(packetTimestamp)
	tsInfoMs := uint64(timestamp) * 1000
	tsSenderMs := tsInfoMs + ((uint64(tsRel)+1)*21)/1000

	// Current unix time in milliseconds
	nowMs := uint64(time.Now().Unix()) * 1000

	// Verification
	if (nowMs < tsSenderMs-uint64(clockSkewMs)) ||
		(nowMs > tsSenderMs+uint64(packetLifetimeMs)+uint64(clockSkewMs)) {
		return false
	} else {
		return true
	}
}

// VerifyHVF verifies the correctness of the PHVF (if "last" is false)
// or the LHVF (if "last" is true).
func VerifyHVF(auth []byte, epicpath *epic.EpicPath, s *slayers.SCION,
	timestamp uint32, last bool) bool {

	// Initialize cryptographic MAC function
	f, err := initEpicMac(auth)
	if err != nil {
		return false
	}
	// Prepare the input for the MAC function
	input, err := prepareMacInput(epicpath, s, timestamp)
	if err != nil {
		return false
	}
	// Calculate MAC ("Write" must not return an error: https://godoc.org/hash#Hash)
	if _, err := f.Write(input); err != nil {
		panic(err)
	}
	mac := f.Sum(nil)[:16]
	// Check if the HVF is valid
	var hvf []byte
	if last {
		hvf = epicpath.LHVF
	} else {
		hvf = epicpath.PHVF
	}
	return bytes.Equal(hvf, mac)
}

func IsPenultimateHop(scionRaw *scion.Raw) (bool, error) {
	if scionRaw == nil {
		return true, serrors.New("scion path must not be nil")
	}
	numberHops := scionRaw.NumHops
	currentHop := int(scionRaw.PathMeta.CurrHF)
	return currentHop == numberHops-2, nil
}

func IsLastHop(scionRaw *scion.Raw) (bool, error) {
	if scionRaw == nil {
		return true, serrors.New("scion path must not be nil")
	}
	numberHops := scionRaw.NumHops
	currentHop := int(scionRaw.PathMeta.CurrHF)
	return currentHop == numberHops-1, nil
}
