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
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/scionproto/scion/go/lib/slayers/path/colibri"
)

func TestColibriHopfieldSerializeDecode(t *testing.T) {
	buffer := make([]byte, colibri.LenHopField)
	for i := 1; i < 10; i++ {
		hf := &colibri.HopField{
			IngressId: randUint16(),
			EgressId:  randUint16(),
			Mac:       randBytes(4),
		}
		assert.NoError(t, hf.SerializeTo(buffer))
		hf2 := &colibri.HopField{}
		assert.NoError(t, hf2.DecodeFromBytes(buffer))
		assert.Equal(t, hf, hf2)
	}
}

func randUint64() uint64 {
	return rand.Uint64()
}

func randUint32() uint32 {
	return rand.Uint32()
}

func randUint16() uint16 {
	return uint16(randUint32())
}

func randUint8() uint8 {
	return uint8(randUint32())
}

func randBool() bool {
	if randUint16()%2 == 1 {
		return true
	}
	return false
}

func randBytes(l uint16) []byte {
	r := make([]byte, l)
	rand.Read(r)
	return r
}
