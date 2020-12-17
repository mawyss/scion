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
const LenMinColibri int = 8 + LenInfoField + 2*LenHopField

func RegisterPath() {
	path.RegisterPath(path.Metadata{
		Type: PathType,
		Desc: "Colibri",
		New: func() path.Path {
			return &ColibriPath{
				InfoField: &InfoField{},
			}
		},
	})
}

type ColibriPath struct {
	// PacketTimestamp denotes the high-precision timestamp.
	PacketTimestamp uint64
	// InfoField denotes the COLIBRI info field.
	InfoField *InfoField
	// HopFields denote the COLIBRI hop fields.
	HopFields []*HopField
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
	if err := c.InfoField.DecodeFromBytes(b[8 : 8+LenInfoField]); err != nil {
		return err
	}
	nrHopFields := int(c.InfoField.HFCount)
	if 8+LenInfoField+(nrHopFields*LenHopField) > len(b) {
		return serrors.New("raw colibri path is smaller than what is " +
			"indicated by HFCount in the info field")
	}
	c.HopFields = make([]*HopField, nrHopFields)
	for i := 0; i < nrHopFields; i++ {
		start := 8 + LenInfoField + i*LenHopField
		end := start + LenHopField
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
	if err := c.InfoField.SerializeTo(b[8 : 8+LenInfoField]); err != nil {
		return err
	}
	for i, hf := range c.HopFields {
		start := 8 + LenInfoField + i*LenHopField
		end := start + LenHopField
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
