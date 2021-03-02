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

// Package colibri contains methods for the creation and verification of the colibri packet
// timestamp and validation fields.
package colibri

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"math"
	"time"

	"github.com/scionproto/scion/go/lib/serrors"
	"github.com/scionproto/scion/go/lib/slayers"
	"github.com/scionproto/scion/go/lib/slayers/path/colibri"
)

const (
	// packetLifetimeMs denotes the maximal lifetime of a packet in milliseconds
	packetLifetimeMs uint16 = 2000
	// clockSkewMs denotes the maximal clock skew in milliseconds
	clockSkewMs uint16 = 1000
	// lengthInputData denotes the length of InputData in bytes
	lengthInputData = 30
	// lengthInputDataRound16 denotes the lengthInputData rounded to the next multiple of 16
	lengthInputDataRound16 = 32
)

// CreateColibriTimestamp creates the COLIBRI packetTimestamp from tsRel, coreID, and coreCounter.
func CreateColibriTimestamp(tsRel uint32, coreID uint8,
	coreCounter uint32) (packetTimestamp uint64) {
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
func ParseColibriTimestamp(packetTimestamp uint64) (tsRel uint32, coreID uint8,
	coreCounter uint32) {

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b[:8], packetTimestamp)
	tsRel = binary.BigEndian.Uint32(b[:4])
	coreID = uint8(binary.BigEndian.Uint16(b[3:5]))
	coreCounter = binary.BigEndian.Uint32(b[4:8]) % (1 << 24)
	return tsRel, coreID, coreCounter
}

// ParseColibriTimestampCustom reads tsRel and pckId from the packetTimestamp.
func ParseColibriTimestampCustom(packetTimestamp uint64) (tsRel uint32, pckId uint32) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b[:8], packetTimestamp)
	tsRel = binary.BigEndian.Uint32(b[:4])
	pckId = binary.BigEndian.Uint32(b[4:8])
	return
}

// CreateTsRel returns tsRel, which encodes the current time (the time when this function is called)
// relative to the expiration time minus 16 seconds. The input expiration tick must be specified in
// ticks of four seconds since Unix time.
// If the current time is not between the expiration time minus 16 seconds and the expiration time,
// an error is returned.
func CreateTsRel(expirationTick uint32) (uint32, error) {
	expirationNano := 4 * uint64(expirationTick) * uint64(math.Pow10(9))
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

// VerifyExpirationTick returns whether the expiration time has not been reached yet.
func VerifyExpirationTick(expirationTick uint32) bool {
	expTime := 4 * int64(expirationTick)
	now := time.Now().Unix()
	return now <= expTime
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

// VerifyMAC verifies the authenticity of the MAC in the colibri hop field. If the MAC is correct,
// nil is returned, otherwise VerifyMAC returns an error.
func VerifyMAC(privateKey []byte, packetTimestamp uint64, inf *colibri.InfoField,
	currHop *colibri.HopField, s *slayers.SCION) error {

	var mac []byte
	var err error

	switch inf.C {
	case true:
		mac, err = CalculateColibriMacStatic(privateKey, inf, currHop, s)
		if err != nil {
			return err
		}
	case false:
		auth, err := CalculateColibriMacSigma(privateKey, inf, currHop, s)
		if err != nil {
			return err
		}
		mac, err = CalculateColibriMacPacket(auth, packetTimestamp, inf, s)
		if err != nil {
			return err
		}
	}

	if !bytes.Equal(mac[:4], currHop.Mac[:4]) {
		return serrors.New("colibri mac verification failed", "calculated", mac[:4],
			"packet", currHop.Mac[:4])
	}

	return nil
}

// CalculateColibriMacStatic calculates the static colibri MAC.
func CalculateColibriMacStatic(privateKey []byte, inf *colibri.InfoField,
	currHop *colibri.HopField, s *slayers.SCION) ([]byte, error) {

	// Initialize cryptographic MAC function
	f, err := initColibriMac(privateKey)
	if err != nil {
		return nil, err
	}
	// Prepare the input for the MAC function
	input, err := prepareMacInputStatic(s, inf, currHop)
	if err != nil {
		return nil, err
	}
	if len(input) < 16 || len(input)%16 != 0 {
		return nil, serrors.New("colibri static mac input has invalid length", "expected", 16,
			"is", len(input))
	}
	// Calculate CBC-MAC = first 4 bytes of the last CBC block
	mac := make([]byte, len(input))
	f.CryptBlocks(mac, input)
	return mac[len(mac)-16 : len(mac)-12], nil
}

// CalculateColibriMacSigma calculates the "sigma" authenticator.
func CalculateColibriMacSigma(privateKey []byte, inf *colibri.InfoField,
	currHop *colibri.HopField, s *slayers.SCION) ([]byte, error) {

	// Initialize cryptographic MAC function
	f, err := initColibriMac(privateKey)
	if err != nil {
		return nil, err
	}
	// Prepare the input for the MAC function
	input, err := prepareMacInputSigma(s, inf, currHop)
	if err != nil {
		return nil, err
	}

	// Calculate CBC-MAC = last CBC block
	mac := make([]byte, len(input))
	f.CryptBlocks(mac, input)
	return mac[len(mac)-16:], nil
}

// CalculateColibriMacPacket calculates the per-packet colibri MAC.
func CalculateColibriMacPacket(auth []byte, packetTimestamp uint64,
	inf *colibri.InfoField, s *slayers.SCION) ([]byte, error) {

	// Initialize cryptographic MAC function
	f, err := initColibriMac(auth)
	if err != nil {
		return nil, err
	}
	// Prepare the input for the MAC function
	input, err := prepareMacInputPacket(packetTimestamp, inf, s)
	if err != nil {
		return nil, err
	}

	// Calculate CBC-MAC = first 4 bytes of the last CBC block
	mac := make([]byte, len(input))
	f.CryptBlocks(mac, input)
	return mac[len(mac)-16 : len(mac)-12], nil
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

func prepareMacInputStatic(s *slayers.SCION, inf *colibri.InfoField,
	hop *colibri.HopField) ([]byte, error) {

	// Create buffer large enough to store InputData, with length aligned to 16 bytes
	buffer := make([]byte, lengthInputDataRound16)
	err := prepareInputData(s, inf, hop, buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func prepareMacInputSigma(s *slayers.SCION, inf *colibri.InfoField,
	hop *colibri.HopField) ([]byte, error) {

	// Check consistency of SL and DL with the actual address lengths
	srcLen := len(s.RawSrcAddr)
	dstLen := len(s.RawDstAddr)
	consistent := (4*(int(s.DstAddrLen)+1) == dstLen) &&
		(4*(int(s.SrcAddrLen)+1) == srcLen)
	if !consistent {
		return nil, serrors.New("SL/DL not consistent with actual address lengths",
			"DL", s.DstAddrLen, "SL", s.SrcAddrLen)
	}

	// Write SL/ST/DL/DT into one single byte
	flags := uint8(s.DstAddrType&0x3)<<6 | uint8(s.DstAddrLen&0x3)<<4 |
		uint8(s.SrcAddrType&0x3)<<2 | uint8(s.SrcAddrLen&0x3)

	// The MAC input consists of the InputData plus the host addresses and the flags, rounded
	// up to the next multiple of 16 bytes
	bufLen := lengthInputData + 1 + srcLen + dstLen
	nrBlocks := uint8(math.Ceil(float64(bufLen) / 16))
	buffer := make([]byte, 16*nrBlocks)

	err := prepareInputData(s, inf, hop, buffer)
	if err != nil {
		return nil, err
	}
	buffer[lengthInputData] = flags
	copy(buffer[lengthInputData+1:], s.RawSrcAddr)
	copy(buffer[lengthInputData+1+srcLen:], s.RawDstAddr)

	return buffer, nil
}

func prepareMacInputPacket(packetTimestamp uint64, inf *colibri.InfoField,
	s *slayers.SCION) ([]byte, error) {

	if inf == nil {
		return nil, serrors.New("invalid input")
	}

	input := make([]byte, 16)
	binary.BigEndian.PutUint64(input[0:8], packetTimestamp)

	baseHdrLen := uint64(slayers.CmnHdrLen + s.AddrHdrLen())
	hfcount := uint64(inf.HFCount)
	colHdrLen := 32 + (hfcount * 8)
	payloadLen := uint64(inf.OrigPayLen)
	total64 := baseHdrLen + colHdrLen + payloadLen
	if total64 > (1 << 16) {
		return nil, serrors.New("total packet length bigger than 2^16")
	}
	total16 := uint16(total64)

	binary.BigEndian.PutUint16(input[8:10], total16)

	return input, nil
}

// prepareInputData writes InputData to the given buffer.
func prepareInputData(s *slayers.SCION, inf *colibri.InfoField,
	hop *colibri.HopField, buffer []byte) error {

	if s == nil || inf == nil || hop == nil {
		return serrors.New("invalid input")
	}
	if len(buffer) < lengthInputData {
		return serrors.New("provided buffer is too small")
	}

	copy(buffer[0:12], inf.ResIdSuffix)
	binary.BigEndian.PutUint32(buffer[12:16], inf.ExpTick)
	buffer[16] = inf.BwCls
	buffer[17] = inf.Rlc
	buffer[18] = 0

	// Version | C | 0
	var flags uint8
	if inf.C {
		flags = uint8(1) << 3
	}
	flags += inf.Ver << 4
	buffer[19] = flags
	srcAs := uint64(s.SrcIA.A)
	binary.BigEndian.PutUint64(buffer[22:30], srcAs)
	binary.BigEndian.PutUint16(buffer[20:22], hop.IngressId)
	binary.BigEndian.PutUint16(buffer[22:24], hop.EgressId)

	return nil
}

func max(x, y uint64) uint64 {
	if x < y {
		return y
	}
	return x
}
