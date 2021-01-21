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
			return &ColibriPathMinimal{
				InfoField:    &InfoField{},
				CurrHopField: &HopField{},
			}
		},
	})
}

// ColibriPathMinimal denotes a COLIBRI path representation optimized for the border router. Only
// the current hop field is parsed, the border router does not need the other ones.
type ColibriPathMinimal struct {
	// PacketTimestamp denotes the high-precision timestamp.
	PacketTimestamp uint64
	// InfoField denotes the COLIBRI info field.
	InfoField *InfoField
	// CurrHopField denotes the current COLIBRI hop field.
	CurrHopField *HopField
	// Raw contains the raw bytes of the COLIBRI path type header. It is set during the execution
	// of DecodeFromBytes.
	Raw []byte
}

func (c *ColibriPathMinimal) DecodeFromBytes(b []byte) error {
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
	currHF := int(c.InfoField.CurrHF)
	if 8+LenInfoField+(nrHopFields*LenHopField) > len(b) {
		return serrors.New("raw colibri path is smaller than what is " +
			"indicated by HFCount in the info field")
	}
	if currHF >= nrHopFields {
		return serrors.New("colibri currHF >= nrHopFields", "currHF", currHF,
			"nrHopFields", nrHopFields)
	}
	c.CurrHopField = &HopField{}
	start := 8 + LenInfoField + currHF*LenHopField
	end := start + LenHopField
	if err := c.CurrHopField.DecodeFromBytes(b[start:end]); err != nil {
		return err
	}
	c.Raw = b[:c.Len()]
	return nil
}

// SerializeToInternal serializes the COLIBRI timestamp and info field to the Raw buffer. No hop
// field is serialized.
func (c *ColibriPathMinimal) SerializeToInternal() error {
	if c == nil {
		return serrors.New("colibri path must not be nil")
	}
	if c.InfoField == nil {
		return serrors.New("the colibri info field must not be nil")
	}
	if c.CurrHopField == nil {
		return serrors.New("the colibri hop field must not be nil")
	}
	if c.Raw == nil {
		return serrors.New("internal Raw buffer must not be nil")
	}
	if c.InfoField.HFCount < 2 {
		return serrors.New("a colibri path must have at least two hop fields")
	}
	if len(c.Raw) < c.Len() {
		return serrors.New("internal Raw buffer for ColibriPath too short", "is:", len(c.Raw),
			"needs:", c.Len())
	}
	binary.BigEndian.PutUint64(c.Raw[0:8], c.PacketTimestamp)
	if err := c.InfoField.SerializeTo(c.Raw[8 : 8+LenInfoField]); err != nil {
		return err
	}
	return nil
}

// SerializeTo serializes the COLIBRI timestamp and info field to the Raw buffer. No hop field is
// serialized. Then the Raw buffer is copied to b.
func (c *ColibriPathMinimal) SerializeTo(b []byte) error {
	if len(b) < c.Len() {
		return serrors.New("buffer for ColibriPath too short", "is:", len(b),
			"needs:", c.Len())
	}
	if err := c.SerializeToInternal(); err != nil {
		return err
	}
	copy(b[:c.Len()], c.Raw[:c.Len()])
	return nil
}

// Reverse reverses the path: the R-flag is toggled and the order of the Hop Fields is inverted.
// The currHF field is updated to still point to the current hop field.
// This is reflected in the underlying Raw buffer, as well as the updated Info and Hop Field.
func (c *ColibriPathMinimal) Reverse() (path.Path, error) {
	// XXX(mawyss): The current implementation is not the most performant, as it parses the entire
	// path first. If this becomes a performance bottleneck, the implementation should be changed to
	// work directly on the ColibriPathMinimal.Raw buffer.

	if c == nil {
		return nil, serrors.New("colibri path must not be nil")
	}
	if c.InfoField == nil {
		return nil, serrors.New("the colibri info field must not be nil")
	}

	colibriPath, err := c.ToColibriPath()
	if err != nil {
		return nil, err
	}

	reversed, err := colibriPath.Reverse()
	if err != nil {
		return nil, err
	}

	if err := reversed.SerializeTo(c.Raw); err != nil {
		return nil, err
	}

	err = c.DecodeFromBytes(c.Raw)
	return c, err
}

func (c *ColibriPathMinimal) Len() int {
	if c == nil || c.InfoField == nil || c.CurrHopField == nil {
		return 0
	}
	nrHopFields := int(c.InfoField.HFCount)
	return 8 + LenInfoField + nrHopFields*LenHopField
}

func (c *ColibriPathMinimal) Type() path.Type {
	return PathType
}

// UpdateCurrHF increases the CurrHF index.
// The CurrHopField reference is not updated and will still point to the old hop field.
func (c *ColibriPathMinimal) UpdateCurrHF() error {
	if c == nil {
		return serrors.New("colibri path must not be nil")
	}
	if c.InfoField == nil {
		return serrors.New("the colibri info field must not be nil")
	}

	if c.InfoField.CurrHF+1 >= c.InfoField.HFCount {
		return serrors.New("colibri path already at end")
	}
	c.InfoField.CurrHF = c.InfoField.CurrHF + 1
	return nil
}

// IsLastHop returns whether the currHF index denotes the last hop.
func (c *ColibriPathMinimal) IsLastHop() (bool, error) {
	if c == nil {
		return false, serrors.New("colibri path must not be nil")
	}
	if c.InfoField == nil {
		return false, serrors.New("the colibri info field must not be nil")
	}

	if c.InfoField.CurrHF+1 == c.InfoField.HFCount {
		return true, nil
	}
	return false, nil
}

// ToColibriPath converts ColibriPathMinimal to a ColibriPath.
func (c *ColibriPathMinimal) ToColibriPath() (*ColibriPath, error) {
	if c == nil {
		return nil, serrors.New("colibri path must not be nil")
	}

	// Serialize ColibriPathMinimal to ensure potential changes are written to the buffer.
	if err := c.SerializeToInternal(); err != nil {
		return nil, err
	}

	colibriPath := &ColibriPath{}
	if err := colibriPath.DecodeFromBytes(c.Raw); err != nil {
		return nil, err
	}
	return colibriPath, nil
}
