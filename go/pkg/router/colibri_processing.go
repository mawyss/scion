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

package router

import (
	"github.com/google/gopacket"
	
	libcolibri "github.com/scionproto/scion/go/lib/colibri"
	"github.com/scionproto/scion/go/lib/serrors"
	"github.com/scionproto/scion/go/lib/slayers"
	"github.com/scionproto/scion/go/lib/slayers/path/colibri"
)

type colibriPacketProcessor struct {
	// d is a reference to the dataplane instance that initiated this processor.
	d *DataPlane
	// ingressID is the interface ID this packet came in, determined from the
	// socket.
	ingressID uint16
	// rawPkt is the raw packet, it is updated during processing to contain the
	// message to send out.
	rawPkt []byte
	// scionLayer is the SCION gopacket layer (common/address header).
	scionLayer slayers.SCION
	// origPacket is the raw original packet, must not be modified.
	origPacket []byte
	// buffer is the buffer that can be used to serialize gopacket layers.
	buffer gopacket.SerializeBuffer

	// hopField is the current hopField field, is updated during processing.
	colibriPathMinimal *colibri.ColibriPathMinimal
}

func (c *colibriPacketProcessor) process() (processResult, error) {
	// Get path
	if pRes, err := c.getPath(); err != nil {
		return pRes, err
	}
	// Basic validation checks
	if pRes, err := c.basicValidation(); err != nil {
		return pRes, err
	}
	// Check hop field MAC
	if pRes, err := c.cryptographicValidation(); err != nil {
		return pRes, err
	}

	// TODO: Forwarding...

	return processResult{}, nil
}

func (c *colibriPacketProcessor) basicValidation() (processResult, error) {
	R := c.colibriPathMinimal.InfoField.R
	S := c.colibriPathMinimal.InfoField.R
	C := c.colibriPathMinimal.InfoField.R

	// Consistency of R, S, and C flags: (S or R) implies C
	if (S || R) && !C {
		return processResult{}, serrors.New("invalid flags", "S", S, "R", R, "C", C)
	}
	// Correct ingress interface
	if R && c.ingressID != c.colibriPathMinimal.CurrHopField.IngressId {
		return processResult{}, serrors.New("invalid ingress identifier")
	}
	if !R && c.ingressID != c.colibriPathMinimal.CurrHopField.EgressId {
		return processResult{}, serrors.New("invalid ingress identifier")
	}
	// Valid packet length
	if int(c.scionLayer.PayloadLen) != len(c.scionLayer.Payload) {
		return processResult{}, serrors.New("packet length validation failed")
	}
	// Colibri path has at least two hop fields
	hfCount := c.colibriPathMinimal.InfoField.HFCount
	currHF := c.colibriPathMinimal.InfoField.CurrHF
	if hfCount < 2 {
		return processResult{}, serrors.New("colibri path needs to have at least 2 hop fields",
			"has", hfCount)
	}
	// Valid current hop field index
	if currHF >= hfCount {
		return processResult{}, serrors.New("invalid current hop field index",
			"currHF", currHF, "hfCount", hfCount)
	}

	// TODO: Packet freshness

	return processResult{}, nil
}

func (c *colibriPacketProcessor) cryptographicValidation() (processResult, error) {

	// TODO: MAC computation
	return processResult{}, nil
}

func (c *colibriPacketProcessor) getPath() (processResult, error) {
	var ok bool
	c.colibriPathMinimal, ok = c.scionLayer.Path.(*colibri.ColibriPathMinimal)
	if !ok {
		return processResult{}, serrors.New("getting minimal colibri path information failed")
	}
	return processResult{}, nil

}