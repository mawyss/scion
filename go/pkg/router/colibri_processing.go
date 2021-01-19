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

	"github.com/scionproto/scion/go/lib/addr"
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

	// colibriPathMinimal is the optimized representation of the colibri path type.
	colibriPathMinimal *colibri.ColibriPathMinimal
}

func (c *colibriPacketProcessor) process() (processResult, error) {
	if c == nil {
		return processResult{}, serrors.New("colibri packet processor must not be nil")
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

func (c *colibriPacketProcessor) getPath() (processResult, error) {
	var ok bool
	c.colibriPathMinimal, ok = c.scionLayer.Path.(*colibri.ColibriPathMinimal)
	if !ok {
		return processResult{}, serrors.New("getting minimal colibri path information failed")
	}
	return processResult{}, nil

}

func (c *colibriPacketProcessor) basicValidation() (processResult, error) {
	R := c.colibriPathMinimal.InfoField.R
	S := c.colibriPathMinimal.InfoField.S
	C := c.colibriPathMinimal.InfoField.C

	// Consistency of flags: S implies C
	if S && !C {
		return processResult{}, serrors.New("invalid flags", "S", S, "R", R, "C", C)
	}
	// Correct ingress interface
	if (R && c.ingressID != c.colibriPathMinimal.CurrHopField.EgressId) ||
		(!R && c.ingressID != c.colibriPathMinimal.CurrHopField.IngressId) {

		return processResult{}, serrors.New("invalid ingress identifier")
	}
	// Valid packet length
	if (!R && c.scionLayer.PayloadLen != c.colibriPathMinimal.InfoField.OrigPayLen) ||
		(int(c.scionLayer.PayloadLen) != len(c.scionLayer.Payload)) {

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

	// Reservation not expired
	expTick := c.colibriPathMinimal.InfoField.ExpTick
	notExpired := libcolibri.VerifyExpirationTick(expTick)
	if !notExpired {
		return processResult{}, serrors.New("packet expired")
	}

	// Packet freshness
	if !C {
		timestamp := c.colibriPathMinimal.PacketTimestamp
		isFresh := libcolibri.VerifyTimestamp(expTick, timestamp)
		if !isFresh {
			return processResult{}, serrors.New("verification of packet timestamp failed")
		}
	}

	return processResult{}, nil
}

func (c *colibriPacketProcessor) cryptographicValidation() (processResult, error) {
	privateKey := c.d.ColibriKey
	colHeader := c.colibriPathMinimal
	err := libcolibri.VerifyMAC(privateKey, colHeader.PacketTimestamp, colHeader.InfoField,
		colHeader.CurrHopField, &c.scionLayer)
	return processResult{}, err
}

func (c *colibriPacketProcessor) forward() (processResult, error) {
	egressId, err := c.egressInterface()
	if err != nil {
		return processResult{}, err
	}

	if c.ingressID == 0 {
		// Received packet from within AS
		return c.forwardToLocalEgress(egressId)
	}

	// Received packet from outside of the AS
	if c.colibriPathMinimal.InfoField.C {
		// Control plane forwarding
		// Assumption: in case there are multiple COLIBRI services, they are always synchronized
		return c.forwardToColibriSvc()
	} else {
		// Data plane forwarding
		if c.destinedToLocalHost(egressId) {
			return c.forwardToLocalHost()
		} else {
			if r, err := c.forwardToLocalEgress(egressId); err == nil {
				return r, err
			}
			return c.forwardToRemoteEgress(egressId)
		}
	}
}

func (c *colibriPacketProcessor) egressInterface() (uint16, error) {
	if c == nil {
		return 0, serrors.New("colibri packet processor must not be nil")
	}
	if c.colibriPathMinimal == nil {
		return 0, serrors.New("colibri path must not be nil")
	}
	if c.colibriPathMinimal.CurrHopField == nil {
		return 0, serrors.New("colibri hop field must not be nil")
	}

	if c.colibriPathMinimal.InfoField.R {
		return c.colibriPathMinimal.CurrHopField.IngressId, nil
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
	if conn, ok := c.d.external[egressId]; ok {
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
		return processResult{}, serrors.New("no external interface with this id",
			"egressId", egressId)
	}
}

func (c *colibriPacketProcessor) forwardToRemoteEgress(egressId uint16) (processResult, error) {
	// AS transit: the packet will leave the AS from another border router.
	if a, ok := c.d.internalNextHops[egressId]; ok {
		return processResult{OutConn: c.d.internal, OutAddr: a, OutPkt: c.rawPkt}, nil
	} else {
		return processResult{}, serrors.New("no remote border router with this egress id",
			"egressId", egressId)
	}
}

func (c *colibriPacketProcessor) forwardToLocalHost() (processResult, error) {
	// Inbound: packet destined to a host in the local IA.
	a, err := c.d.resolveLocalDst(c.scionLayer)
	if err != nil {
		return processResult{}, err
	}
	return processResult{OutConn: c.d.internal, OutAddr: a, OutPkt: c.rawPkt}, nil
}

func (c *colibriPacketProcessor) forwardToColibriSvc() (processResult, error) {
	// Inbound: packet destined to the local colibri service.

	// Get address of colibri service (pick one at random if there are multiple COLIBRI services)
	// Assumption: in case there are multiple COLIBRI services, they are always synchronized
	a, ok := c.d.svc.Any(addr.SvcCOL.Base())
	if !ok {
		return processResult{}, serrors.New("no colibri service registered at border router")
	}
	return processResult{OutConn: c.d.internal, OutAddr: a, OutPkt: c.rawPkt}, nil
}

func (c *colibriPacketProcessor) destinedToLocalHost(egressId uint16) bool {
	isLast, _ := c.colibriPathMinimal.IsLastHop()
	return c.scionLayer.DstIA.Equal(c.d.localIA) && egressId == 0 && isLast
}
