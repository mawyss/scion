package libepic

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"hash"
	"time"

	"github.com/dchest/cmac"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/log"
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
	clockSkewMs        uint16 = 1000
)

func initEpicMac(key []byte) (hash.Hash, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, common.NewBasicError(ErrCipherFailure, err)
	}

	// Note: CMAC is not ideal for EPIC due to its subkey generation overhead.
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
// The "flags" field encodes the length of the source address.
func inputToBytes(timestamp uint32, packetTimestamp uint64, 
	scrIA addr.IA, srcAddr []byte, srcAddrLen uint8, payloadLen uint16) []byte {

	// todo: Check validity of srcAddrLen

	// Create a multiple of 16 such that the input fits in
	nrBlocks := int((23 + 4*(int(srcAddrLen) + 1))/16)
	input := make([]byte, 16 * nrBlocks)


	//binary.BigEndian.PutUint16(input[2:4], segID)
	//binary.BigEndian.PutUint32(input[4:8], timestamp)
	//input[9] = expTime
	//binary.BigEndian.PutUint16(input[10:12], consIngress)
	//binary.BigEndian.PutUint16(input[12:14], consEgress)
	return input
}

func prepareMacInput(epicpath *epic.EpicPath) []byte {
	log.Debug("test: ")
	return nil
}

func CreateEpicTimestamp(tsRel uint32, coreID uint8, coreCounter uint32) (packetTimestamp uint64) {
	b := make([]byte, 9)
	binary.BigEndian.PutUint16(b[3:5], uint16(coreID))
	binary.BigEndian.PutUint32(b[:4], tsRel)
	binary.BigEndian.PutUint32(b[5:9], coreCounter)
	packetTimestamp = binary.BigEndian.Uint64(b[:8])
	return
}

func ParseEpicTimestamp(packetTimestamp uint64) (tsRel uint32, coreID uint8, coreCounter uint32) {
	b := make([]byte, 9)
	binary.BigEndian.PutUint64(b[:8], packetTimestamp)
	tsRel = binary.BigEndian.Uint32(b[:4])
	coreID = uint8(binary.BigEndian.Uint16(b[3:5]))
	coreCounter = binary.BigEndian.Uint32(b[5:9])
	return tsRel, coreID, coreCounter
}

// todo: rename camel case
func VerifyTimestamp(timestamp uint32, packetTimestamp uint64) bool {
	// Unix time in milliseconds when the packet was timestamped by the sender
	tsRel, _, _ := ParseEpicTimestamp(packetTimestamp)
	ts_info_ms := uint64(timestamp) * 10^3
	ts_sender_ms := ts_info_ms + ((uint64(tsRel)+1)*21)/1000

	// Current unix time in milliseconds
	now_ms := uint64(time.Now().Unix())*10^3

	// Verification
	if (now_ms < ts_sender_ms - uint64(clockSkewMs)) || (now_ms > ts_sender_ms + uint64(packetLifetimeMs) + uint64(clockSkewMs)) {
		return false
	} else {
		return true
	}
}

// VerifyHVF verifies the correctness of the PHVF (if "last" is false)
// or the LHVF (if "last" is true).
func VerifyHVF(auth []byte, epicpath *epic.EpicPath, last bool) bool {
	// Initialize cryptographic MAC function
	f, err := initEpicMac(auth)
	if err != nil {
		return false
	}

	// Prepare the input for the MAC function
	input := prepareMacInput(epicpath)

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

func IsPenultimateHop(scionRaw *scion.Raw) bool {
	numberHops := scionRaw.NumHops
	currentHop := int(scionRaw.PathMeta.CurrHF)

	return currentHop == numberHops - 2
}

func IsLastHop(scionRaw *scion.Raw) bool {
	numberHops := scionRaw.NumHops
	currentHop := int(scionRaw.PathMeta.CurrHF)

	return currentHop == numberHops - 1
}
