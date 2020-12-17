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
	"github.com/scionproto/scion/go/lib/slayers/path"
)

const PathType path.Type = 4
const LenInfoField int = 24
const LenHopField int = 8
const LenMinColibri int = 8 + LenInfoField + 2*LenHopField

func RegisterPath() {
	path.RegisterPath(path.Metadata{
		Type: PathType,
		Desc: "Colibri",
		New: func() path.Path {
			return &ColibriPath {
				InfoField: &InfoField{},
			}
		},
	})
}

type ColibriPath struct {
	// PacketTimestamp denotes the high-precision timestamp.
	PacketTimestamp uint64
	// InfoField denotes the COLIBRI info field.
	InfoField       *InfoField
	// HopFields denote the COLIBRI hop fields.
	HopFields       []*HopField
}

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
	return nil
}

func (inf *InfoField) SerializeTo(b []byte) error {
	return nil
}

type HopField struct {
	// IngressId denotes the ingress interface in the direction of the reservation (R=0).
	IngressId uint16
	// EgressId denotes the egress interface in the direction of the reservation (R=0).
	EgressId uint16
	// Mac (4 bytes) denotes the MAC (static or per-packet MAC, depending on the S flag).
	Mac []byte
}

func (hf *HopField) DecodeFromBytes(b []byte) error {
	return nil
}

func (hf *HopField) SerializeTo(b []byte) error {
	return nil
}

func (c *ColibriPath) DecodeFromBytes(b []byte) error {
	if c == nil {
		return serrors.New("colibri path must not be nil")
	}
	if len(b) < LenMinColibri {
		return serrors.New("raw colibri path too short", "is:", len(b),
			"needs:", LenMinColibri)
	}

	c.PacketTimestamp = binary.BigEndian.Uint64(b[:8])
	if c.InfoField == nil {
		c.InfoField = &InfoField{}
	}
	nrHopFields := int(c.InfoField.HFCount)
	if 8 + LenInfoField + nrHopFields*LenHopField > len(b) {
		return serrors.New("raw colibri path is smaller than what is " +
			"indicated by HFCount in the info field")
	}
	if err := c.InfoField.DecodeFromBytes(b[8:8+LenInfoField]); err != nil {
		return err
	}
	c.HopFields = make([]*HopField, nrHopFields)
	for i := 0; i < nrHopFields; i++ {
		start := 8 + LenInfoField + i*LenHopField
		end := start + (i+1)*LenHopField
		c.HopFields[i] = &HopField{}
		if err := c.HopFields[i].DecodeFromBytes(b[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (c *ColibriPath) SerializeTo(b []byte) error {
	if c == nil {
		return serrors.New("colibri path must not be nil")
	}
	if c.InfoField == nil {
		return serrors.New("the info field must not be nil")
	}
	if len(c.HopFields) < 2 {
		return serrors.New("a colibri path must have at least two hop fields")
	}
	if len(b) < c.Len() {
		return serrors.New("buffer for ColibriPath too short", "is:", len(b),
			"needs:", c.Len())
	}

	binary.BigEndian.PutUint64(b[0:8], c.PacketTimestamp)
	if err := c.InfoField.SerializeTo(b[8:8+LenInfoField]); err != nil {
		return err
	}
	for i, hf := range c.HopFields {
		start := 8 + LenInfoField + i*LenHopField
		end := start + (i+1)*LenHopField
		if err := hf.SerializeTo(b[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (c *ColibriPath) Reverse() (path.Path, error) {
	if c == nil {
		return nil, serrors.New("colibri path must not be nil")
	}
	if c.InfoField == nil {
		return nil, serrors.New("the info field must not be nil")
	}
	c.InfoField.R = !c.InfoField.R
	return c, nil
}

func (c *ColibriPath) Len() int {
	if c == nil {
		return 0
	}
	return 8 + LenInfoField + len(c.HopFields)*LenHopField
}

func (c *ColibriPath) Type() path.Type {
	return PathType
}
