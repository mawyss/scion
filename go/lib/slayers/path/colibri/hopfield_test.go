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

package colibri_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/scionproto/scion/go/lib/slayers/path/colibri"
)

func TestColibriHopfieldSerializeDecode(t *testing.T) {
	buffer := make([]byte, colibri.LenHopField)
	hf := &colibri.HopField{
		IngressId: 35,
		EgressId:  24,
		Mac:       []byte{0xf2, 0x83, 0x54, 0xaa},
	}
	assert.NoError(t, hf.SerializeTo(buffer))
	hf2 := &colibri.HopField{}
	assert.NoError(t, hf2.DecodeFromBytes(buffer))
	assert.Equal(t, hf, hf2)
}
