// Copyright 2020 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package colibri

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"math"
	"time"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/serrors"
	"github.com/scionproto/scion/go/lib/slayers"
	"github.com/scionproto/scion/go/lib/slayers/path/colibri"
	//"github.com/scionproto/scion/go/lib/slayers/path/scion"
)

const (
	// packetLifetimeMs denotes the maximal lifetime of a packet in milliseconds
	packetLifetimeMs uint16 = 2000
	// clockSkewMs denotes the maximal clock skew in milliseconds
	clockSkewMs uint16 = 1000
)

// CreateColibriTimestamp creates the COLIBRI packetTimestamp from tsRel, coreID, and coreCounter.
func CreateColibriTimestamp(tsRel uint32, coreID uint8, coreCounter uint32) (packetTimestamp uint64) {
	// 0                   1                   2                   3
	// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                             TsRel                             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |    CoreID     |                  CoreCounter                  |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b[4:8], coreCounter)
	binary.BigEndian.PutUint16(b[3:5], uint16(coreID))
	binary.BigEndian.PutUint32(b[:4], tsRel)
	packetTimestamp = binary.BigEndian.Uint64(b[:8])
	return
}

// CreateColibriTimestampCustom creates the COLIBRI packetTimestamp from tsRel and pckId.
func CreateColibriTimestampCustom(tsRel uint32, pckId uint32) (packetTimestamp uint64) {
	// 0                   1                   2                   3
	// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                             TsRel                             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                             PckId                             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b[:4], tsRel)
	binary.BigEndian.PutUint32(b[4:8], pckId)
	packetTimestamp = binary.BigEndian.Uint64(b[:8])
	return
}

// ParseColibriTimestamp reads tsRel, coreID, and coreCounter from the packetTimestamp.
func ParseColibriTimestamp(packetTimestamp uint64) (tsRel uint32, coreID uint8, coreCounter uint32) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b[:8], packetTimestamp)
	tsRel = binary.BigEndian.Uint32(b[:4])
	coreID = uint8(binary.BigEndian.Uint16(b[3:5]))
	coreCounter = binary.BigEndian.Uint32(b[4:8]) % (1 << 24)
	return tsRel, coreID, coreCounter
}

// CreateTsRel returns tsRel, which encodes the current time (the time when this function is called)
// relative to the expiration time minus 16 seconds. The input expiration tick must be specified in
// ticks of four seconds since Unix time.
// If the current time is not between the expiration time minus 16 seconds and the expiration time,
// an error is returned.
func CreateTsRel(expirationTick uint32) (uint32, error) {
	expirationNano := 4*uint64(expirationTick) * uint64(math.Pow10(9))
	timestampNano := (4*uint64(expirationTick) - 16) * uint64(math.Pow10(9))
	nowNano := uint64(time.Now().UnixNano())
	if nowNano > expirationNano {
		return 0, serrors.New("provided packet expiration time is in the past",
			"expiration", expirationNano, "now", nowNano)
	}
	if nowNano < timestampNano {
		return 0, serrors.New("provided packet expiration time is too far in the future",
		"timestampNano", timestampNano, "now", nowNano)
	}
	diff := nowNano - timestampNano
	tsRel := max(0, (diff/4)-1)
	return uint32(tsRel), nil
}

// VerifyTimestamp checks whether a COLIBRI packet is fresh. This means that the time the packet
// was sent from the source host, which is encoded by the expiration tick and the packetTimestamp,
// does not date back more than the maximal packet lifetime of two seconds. The function also takes
// a possible clock drift between the packet source and the verifier of up to one second into
// account.
func VerifyTimestamp(expirationTick uint32, packetTimestamp uint64) bool {
	tsRel, _, _ := ParseColibriTimestamp(packetTimestamp)
	timestampNano := (4*uint64(expirationTick) - 16) * uint64(math.Pow10(9))
	timestampSenderNano := timestampNano + (1+uint64(tsRel))*4
	
	nowMs := uint64(time.Now().UnixNano() / 1000000)
	tsSenderMs := timestampSenderNano / 1000000

	if (nowMs < tsSenderMs-uint64(clockSkewMs)) ||
		(nowMs > tsSenderMs+uint64(packetLifetimeMs)+uint64(clockSkewMs)) {
		return false
	} else {
		return true
	}
}

func CalculateColibriMac(sharedKey []byte, inf *colibri.InfoField, currHop *colibri.HopField) ([]byte, error) {

	// Initialize cryptographic MAC function
	f, err := initColibriMac(sharedKey)
	if err != nil {
		return nil, err
	}
	// Prepare the input for the MAC function
	input, err := prepareMacInput(s, timestamp)
	if err != nil {
		return nil, err
	}
	if len(input) < 16 || len(input)%16 != 0 {
		return nil, serrors.New("colibri mac input has invalid length", "expected", 16,
			"is", len(input))
	}
	// Calculate Colibri MAC = first 4 bytes of the last CBC block
	mac := make([]byte, len(input))
	f.CryptBlocks(mac, input)
	return mac[len(mac)-16 : len(mac)-12], nil
}

func VerifyMAC(auth []byte, timestamp uint32, inf *colibri.InfoField,
	currHop *colibri.HopField, s *slayers.SCION) (bool, error) {

	return bytes.Equal(nil, nil), nil
}

func initColibriMac(key []byte) (cipher.BlockMode, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, serrors.New("Unable to initialize AES cipher")
	}

	// Zero initialization vector
	zeroInitVector := make([]byte, 16)
	// CBC-MAC = CBC-Encryption with zero initialization vector
	mode := cipher.NewCBCEncrypter(block, zeroInitVector)
	return mode, nil
}

func inputToBytes(timestamp uint32, packetTimestamp uint64,
	srcIA addr.IA, srcAddr []byte, srcAddrLen uint8, payloadLen uint16) ([]byte, error) {

	// TODO

	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   | flags (1B) | timestamp (4B) | packetTimestamp (8B)  |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   | srcIA (8B) | srcAddr (4/8/12/16B) | payloadLen (2B) |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   | zero padding (0-15B)                                |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// The "flags" field only encodes the length of the source address.

	l := int((srcAddrLen + 1) * 4)
	if srcAddrLen > 3 || l != len(srcAddr) {
		return nil, serrors.New("srcAddrLen must be between 0 and 3, and encode the "+
			"srcAddr length", "srcAddrLen", srcAddrLen, "len(srcAddr)", len(srcAddr),
			"l", l)
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

func prepareMacInput(s *slayers.SCION, timestamp uint32) ([]byte, error) {
	if s == nil {
		return nil, serrors.New("SCION common+address header must not be nil")
	}
	payloadLen := s.PayloadLen
	srcIA := s.SrcIA
	srcAddrLen := uint8(s.SrcAddrLen)
	srcAddr := s.RawSrcAddr
	//return inputToBytes(timestamp, srcIA, srcAddr, srcAddrLen, payloadLen)
	return nil, nil
}

func max(x, y uint64) uint64 {
	if x < y {
		return y
	}
	return x
}
