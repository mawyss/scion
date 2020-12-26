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
	"encoding/binary"

	"github.com/scionproto/scion/go/lib/serrors"
)

const LenInfoField int = 24

type InfoField struct {
	// C denotes the control plane flag.
	C bool
	// R denotes the reverse flag.
	R bool
	// S denotes the segment flag.
	S bool
	// CurrHF denotes the current hop field.
	CurrHF uint8
	// HFCount denotes the total number of hop fields.
	HFCount uint8
	// ResIdSuffix (12 bytes) denotes the reservation ID suffix.
	ResIdSuffix []byte
	// ExpTick denotes the expiration tick, where one tick corresponds to 4 seconds.
	ExpTick uint32
	// BwCls denotes the bandwidth class of the reservation.
	BwCls uint8
	// Rlx denotes the request latency class of the reservation.
	Rlc uint8
	// Ver (4 bits) denotes the reservation version.
	Ver uint8
}

func (inf *InfoField) DecodeFromBytes(b []byte) error {
	if inf == nil {
		return serrors.New("colibri info field must not be nil")
	}
	if len(b) < LenInfoField {
		return serrors.New("raw colibri info field buffer too small")
	}
	flags := binary.BigEndian.Uint16(b[:2])
	inf.C = (flags & (uint16(1) << 15)) != 0
	inf.R = (flags & (uint16(1) << 14)) != 0
	inf.S = (flags & (uint16(1) << 13)) != 0
	inf.CurrHF = uint8(binary.BigEndian.Uint16(b[1:3]))
	inf.HFCount = uint8(binary.BigEndian.Uint16(b[2:4]))
	inf.ResIdSuffix = make([]byte, 12)
	copy(inf.ResIdSuffix, b[4:16])
	inf.ExpTick = binary.BigEndian.Uint32(b[16:20])
	inf.BwCls = uint8(binary.BigEndian.Uint16(b[19:21]))
	inf.Rlc = uint8(binary.BigEndian.Uint16(b[20:22]))
	inf.Ver = (uint8(binary.BigEndian.Uint16(b[21:23])) & uint8(0xF0)) >> 4
	return nil
}

func (inf *InfoField) SerializeTo(b []byte) error {
	if inf == nil {
		return serrors.New("colibri info field must not be nil")
	}
	if len(b) < LenInfoField {
		return serrors.New("raw colibri info field buffer too small")
	}
	if len(inf.ResIdSuffix) != 12 {
		return serrors.New("colibri ResIdSuffix must be 12 bytes long",
			"is", len(inf.ResIdSuffix))
	}
	var flags uint16
	if inf.C {
		flags += uint16(1) << 15
	}
	if inf.R {
		flags += uint16(1) << 14
	}
	if inf.S {
		flags += uint16(1) << 13
	}
	binary.BigEndian.PutUint16(b[2:4], uint16(inf.HFCount))
	binary.BigEndian.PutUint16(b[1:3], uint16(inf.CurrHF))
	binary.BigEndian.PutUint16(b[:2], flags)
	copy(b[4:16], inf.ResIdSuffix)
	binary.BigEndian.PutUint16(b[20:22], uint16(inf.Rlc))
	binary.BigEndian.PutUint16(b[19:21], uint16(inf.BwCls))
	binary.BigEndian.PutUint32(b[16:20], inf.ExpTick)
	var endFlags uint16
	endFlags += uint16(inf.Ver<<4) << 8
	binary.BigEndian.PutUint16(b[22:24], endFlags)
	return nil
}

// SerializeToMac serializes the InfoField for the use in the colibri static MAC calculation.
// It does the same as SerializeTo(), but replaces the fields that should not be part of the MAC
// input with zeroes.
func (inf *InfoField) SerializeToMac(b []byte) error {
	if inf == nil {
		return serrors.New("colibri info field must not be nil")
	}
	if len(b) < LenInfoField {
		return serrors.New("raw colibri info field buffer too small")
	}
	if len(inf.ResIdSuffix) != 12 {
		return serrors.New("colibri ResIdSuffix must be 12 bytes long",
			"is", len(inf.ResIdSuffix))
	}
	var flags uint16
	if inf.C {
		flags += uint16(1) << 15
	}
	if inf.R {
		flags += uint16(1) << 14
	}
	if inf.S {
		flags += uint16(1) << 13
	}
	binary.BigEndian.PutUint16(b[2:4], uint16(inf.HFCount))
	binary.BigEndian.PutUint16(b[1:3], uint16(inf.CurrHF))
	binary.BigEndian.PutUint16(b[:2], flags)
	copy(b[4:16], inf.ResIdSuffix)
	binary.BigEndian.PutUint16(b[20:22], uint16(inf.Rlc))
	binary.BigEndian.PutUint16(b[19:21], uint16(inf.BwCls))
	binary.BigEndian.PutUint32(b[16:20], inf.ExpTick)
	var endFlags uint16
	endFlags += uint16(inf.Ver<<4) << 8
	binary.BigEndian.PutUint16(b[22:24], endFlags)
	return nil
}
