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
	"math"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	libcolibri "github.com/scionproto/scion/go/lib/colibri"
	"github.com/scionproto/scion/go/lib/slayers"
	"github.com/scionproto/scion/go/lib/slayers/path/colibri"
	"github.com/scionproto/scion/go/lib/slayers/path/scion"
	"github.com/scionproto/scion/go/lib/xtest"
)

func TestStaticMacInputGeneration(t *testing.T) {
	// TODO

	/*
	want := []byte(
		"\x00\x00")
	s := createScionCmnAddrHdr()
	c := createColibriPath()
	got, err := libcolibri.PrepareMacInputStatic(s, c.InfoField, c.HopFields[0])
	assert.NoError(t, err)
	assert.Equal(t, want, got)
	*/
}

func TestPacketMacInputGeneration(t *testing.T) {
	// TODO
}

func TestTimestamp(t *testing.T) {
	testCases := []uint64{0, 1, math.MaxInt32}
	for i := 1; i <= 10; i++ {
		testCases = append(testCases, randUint64())
	}
	for _, want := range testCases {
		tsRel, coreID, coreCounter := libcolibri.ParseColibriTimestamp(want)
		got := libcolibri.CreateColibriTimestamp(tsRel, coreID, coreCounter)
		assert.Equal(t, want, got)
	}
}

func TestTsRel(t *testing.T) {
	// expTick encodes the current time plus something between 4 and 8 seconds.
	expTick := uint32(time.Now().Unix() / 4) + 2

	// Incrementing tsRel by one corresponds to adding 4 seconds

	// TODO: check this test

	testCases := map[uint32]bool{
		0:           false,
		expTick - 2: false,
		expTick - 1: true,
		expTick:     true,
		expTick + 1: true,
		expTick + 2: true,
		expTick + 3: false,
	}
	for tsRel, want := range testCases {
		_, err := libcolibri.CreateTsRel(tsRel)
		if want == true {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestTimestampVerification(t *testing.T) {

}

func TestStaticHVFVerification(t *testing.T) {

}

func TestPacketHVFVerification(t *testing.T) {

}

func createScionCmnAddrHdr() *slayers.SCION {
	spkt := &slayers.SCION{
		SrcAddrLen: 0,
		SrcIA:      xtest.MustParseIA("2-ff00:0:222"),
		PayloadLen: 120,
	}
	ip4Addr := &net.IPAddr{IP: net.ParseIP("10.0.0.100")}
	spkt.SetSrcAddr(ip4Addr)
	return spkt
}

func createColibriPath() *colibri.ColibriPath {
	ts := libcolibri.CreateColibriTimestamp(1, 2, 3)
	colibripath := &colibri.ColibriPath{
		PacketTimestamp: ts,
		InfoField:       &colibri.InfoField{
			CurrHF:      0,
			HFCount:     10,
			ResIdSuffix: randBytes(12),
			ExpTick:     uint32(time.Now().Unix()/4),
			BwCls:       randBytes(1)[0],
			Rlc:         randBytes(1)[0],
			Ver:         randBytes(1)[0],
		},
	}
	hopfields := make([]*colibri.HopField, 10)
	for i, _ := range hopfields {
		hf := &colibri.HopField{
			IngressId: 1,
			EgressId:  2,
			Mac:       []byte{1, 2, 3, 4},
		}
		hopfields[i] = hf
	}
	colibripath.HopFields = hopfields
	return colibripath
}

func createScionPath(currHF uint8, numHops int) *scion.Raw {
	scionRaw := &scion.Raw{
		Base: scion.Base{
			PathMeta: scion.MetaHdr{
				CurrHF: currHF,
			},
			NumHops: numHops,
		},
	}
	return scionRaw
}

func randUint64() uint64 {
	return uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
}

func randBytes(l uint16) []byte {
	r := make([]byte, l)
	rand.Read(r)
	return r
}
