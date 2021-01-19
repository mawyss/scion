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

func TestColibriSerializeDecode(t *testing.T) {
	for i := 2; i < 11; i++ {
		bufferLength := 8 + colibri.LenInfoField + i*colibri.LenHopField
		buffer := randBytes(uint16(bufferLength))
		// Remove the "reserved" flags from the info field
		buffer[8] = buffer[8] & uint8(0xE0)
		buffer[9] = 0
		buffer[30] = buffer[30] & uint8(0xF0)
		buffer[31] = 0
		// Set correct number of hop fields
		buffer[10] = uint8(i - 1)
		buffer[11] = uint8(i)

		// Test ColibriPath
		col := &colibri.ColibriPath{}
		assert.NoError(t, col.DecodeFromBytes(buffer))

		buffer2 := make([]byte, col.Len())
		assert.NoError(t, col.SerializeTo(buffer2))
		assert.Equal(t, buffer, buffer2)

		// Test ColibriPathMinimal
		colMin := &colibri.ColibriPathMinimal{}
		colMin2 := &colibri.ColibriPathMinimal{}
		assert.NoError(t, colMin.DecodeFromBytes(buffer))
		buffer2 = make([]byte, colMin.Len())
		assert.NoError(t, colMin.SerializeTo(buffer2))
		assert.NoError(t, colMin2.DecodeFromBytes(buffer2))
		assert.Equal(t, colMin, colMin2)
	}
}

func TestColibriReverse(t *testing.T) {
	for i := 2; i < 11; i++ {
		bufferLength := 8 + colibri.LenInfoField + i*colibri.LenHopField
		buffer := randBytes(uint16(bufferLength))
		// Set correct number of hop fields
		buffer[10] = uint8(i - 1)
		buffer[11] = uint8(i)

		old := &colibri.ColibriPath{}
		new := &colibri.ColibriPath{}
		assert.NoError(t, old.DecodeFromBytes(buffer))
		assert.NoError(t, new.DecodeFromBytes(buffer))

		rev, err := new.Reverse()
		new = rev.(*colibri.ColibriPath)
		assert.NoError(t, err)

		assert.Equal(t, old.InfoField.R, !new.InfoField.R)
		for j := 0; j < i/2+1; j++ {
			assert.Equal(t, old.HopFields[j], new.HopFields[i-1-j])
		}

		revrev, err := rev.Reverse()
		assert.Equal(t, revrev, old)
	}
}
