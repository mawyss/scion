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

func TestColibriInfofieldSerializeDecode(t *testing.T) {
	buffer := []byte("073f5f1c20df5381f2e896d0")
	// Remove the "reserved" flags
	buffer[0] = buffer[0] & uint8(0xE0)
	buffer[1] = buffer[1] & uint8(0x0F)

	inf := &colibri.InfoField{}
	assert.NoError(t, inf.DecodeFromBytes(buffer))

	buffer2 := make([]byte, colibri.LenInfoField)
	assert.NoError(t, inf.SerializeTo(buffer2))
	assert.Equal(t, buffer, buffer2)
}
