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
	if c == nil {
		return 0, serrors.New("colibri packet processor must not be nil")
	}

	// Get path
	if r, err := c.getPath(); err != nil {
		return r, err
	}
	// Basic validation checks
	if r, err := c.basicValidation(); err != nil {
		return r, err
	}
	// Check hop field MAC
	if r, err := c.cryptographicValidation(); err != nil {
		return r, err
	}
	// Forward the packet to the correct entity
	return c.forward()
}

func (c *colibriPacketProcessor) forward() (processResult, error) {
	egressId := c.egressInterface()

	if c.ingressID == 0 {
		// Received packet from within AS
		return c.forwardToLocalEgress(egressId)
	}

	// Received packet from outside of the AS
	if c.colibriPathMinimal.InfoField.C {
		// Control plane forwarding
		c.forwardToColibriSvc() // TODO
	} else {
		// Data plane forwarding
		if c.destinedToLocalHost(egressId) {
			c.forwardToLocalHost() // TODO
		} else {
			if r, err := c.forwardToLocalEgress(egressId); err == nil {
				return r, err
			}
			c.forwardToRemoteEgress(egressId)
		}
	}


	// Inbound: packet destined to the local IA. (SCION: resolveInbound())
	if c.scionLayer.DstIA.Equal(c.d.localIA) && c.colibriPathMinimal.IsLastHop() {
		a, err := c.d.resolveLocalDst(c.scionLayer)
		if err != nil {
			return processResult{}, err
		}
		return processResult{OutConn: c.d.internal, OutAddr: a, OutPkt: c.rawPkt}, nil
	}

	// Outbound: pkts leaving the local IA.
	// and
	// BRTransit: pkts leaving from the same BR different interface.

	// Only if we want SCMP for COLIBRI:
	// Checks if the egress interface is registered (either at same BR or at other BR of the same AS)
	//if r, err := p.validateEgressID(); err != nil {
	//	return r, err
	//}

	// Only if we want SCMP for COLIBRI:
	// Check if BFD session for this egress interface is up.
	//if r, err := p.validateEgressUp(); err != nil {
	//	return r, err
	//}

	egressID := p.egressInterface()
	if c, ok := p.d.external[egressID]; ok {
		// Increase the hop pointer, serialize the whole address/common/pathtype header structs into rawPkt
		if err := p.processEgress(); err != nil {
			return processResult{}, err
		}
		return processResult{EgressID: egressID, OutConn: c, OutPkt: p.rawPkt}, nil
	}

	// ASTransit: pkts leaving from another AS BR.
	if a, ok := p.d.internalNextHops[egressID]; ok {
		return processResult{OutConn: p.d.internal, OutAddr: a, OutPkt: p.rawPkt}, nil
	}

	/*
	errCode := slayers.SCMPCodeUnknownHopFieldEgress
	if !p.infoField.ConsDir {
		errCode = slayers.SCMPCodeUnknownHopFieldIngress
	}
	return p.packSCMP(
		&slayers.SCMP{
			TypeCode: slayers.CreateSCMPTypeCode(slayers.SCMPTypeParameterProblem, errCode),
		},
		&slayers.SCMPParameterProblem{Pointer: p.currentHopPointer()},
		cannotRoute,
	)
	*/
}

func (c *colibriPacketProcessor) basicValidation() (processResult, error) {
	R := c.colibriPathMinimal.InfoField.R
	S := c.colibriPathMinimal.InfoField.S
	C := c.colibriPathMinimal.InfoField.C

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

func (c *colibriPacketProcessor) egressInterface() (uint16, error) {
	if c == nil {
		return 0, serrors.New("colibri packet processor must not be nil")
	}
	if c.colibriPathMinimal == nil {
		return 0, serrors.New("colibri path must not be nil")
	}
	if c.colibriPathMinimal.CurrHopField == nil {
		return 0, serrors.New("colibri info field must not be nil")
	}
	if c.colibriPathMinimal.CurrHopField == nil {
		return 0, serrors.New("colibri hop field must not be nil")
	}

	if c.colibriPathMinimal.InfoField.R {
		return c.colibriPathMinimal.CurrHopField.Ingress, nil
	} else {
		return c.colibriPathMinimal.CurrHopField.EgressId, nil
	}
}

func (c *colibriPacketProcessor) ingressInterface() (uint16, error) {
	if c == nil {
		return 0, serrors.New("colibri packet processor must not be nil")
	}
	if c.colibriPathMinimal == nil {
		return 0, serrors.New("colibri path must not be nil")
	}
	if c.colibriPathMinimal.CurrHopField == nil {
		return 0, serrors.New("colibri info field must not be nil")
	}
	if c.colibriPathMinimal.CurrHopField == nil {
		return 0, serrors.New("colibri hop field must not be nil")
	}

	if c.colibriPathMinimal.InfoField.R {
		return c.colibriPathMinimal.CurrHopField.EgressId, nil
	} else {
		return c.colibriPathMinimal.CurrHopField.IngressId, nil
	}
}

func (c *colibriPacketProcessor) forwardToLocalEgress(egressId uint16) (processResult, error) {
	// BR transit: the packet will leave the AS through the same border router, but through a
	// different interface.
	if conn, ok := c.d.external[egressID]; ok {
		// Increase/decrease (depending on "R" flag) the hop field index.
		if err := c.colibriPathMinimal.UpdateCurrHF(); err != nil {
			return processResult{}, err
		}
		// Serialize updated hop field index into rawPkt.
		if err := c.colibriPathMinimal.SerializeToInternal(); err != nil {
			return processResult{}, err
		}
		return processResult{EgressID: egressId, OutConn: conn, OutPkt: c.rawPkt}, nil
	} else {
		return processResult{}, serrors.New("no external interface with this id", "egressId", egressId)
	}
}

func (c *colibriPacketProcessor) forwardToRemoteEgress(egressId uint16) (processResult, error) {
	// AS transit: the packet will leave the AS from another border router.
	if a, ok := c.d.internalNextHops[egressID]; ok {
		return processResult{OutConn: c.d.internal, OutAddr: a, OutPkt: c.rawPkt}, nil
	} else {
		return processResult{}, serrors.New("no remote border router with this egress id",
			"egressId", egressId)
	}
}

func (c *colibriPacketProcessor) destinedToLocalHost(egressId) bool {
	return c.scionLayer.DstIA.Equal(c.d.localIA) && egressId == 0 &&
		c.colibriPathMinimal.IsLastHop()
}
