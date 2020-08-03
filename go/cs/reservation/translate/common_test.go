// Copyright 2020 ETH Zurich, Anapaya Systems
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

package translate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/scionproto/scion/go/cs/reservation/segment"
	"github.com/scionproto/scion/go/lib/colibri/reservation"
	"github.com/scionproto/scion/go/lib/ctrl/colibri_mgmt"
	"github.com/scionproto/scion/go/lib/spath"
	"github.com/scionproto/scion/go/lib/xtest"
	"github.com/scionproto/scion/go/proto"
)

func newID() *colibri_mgmt.SegmentReservationID {
	return &colibri_mgmt.SegmentReservationID{
		ASID:   xtest.MustParseHexString("ff00cafe0001"),
		Suffix: xtest.MustParseHexString("deadbeef"),
	}
}

func newE2EID() *colibri_mgmt.E2EReservationID {
	return &colibri_mgmt.E2EReservationID{
		ASID:   xtest.MustParseHexString("ff00cafe0001"),
		Suffix: xtest.MustParseHexString("0123456789abcdef0123"),
	}
}

func newTestBase(idx uint8) *colibri_mgmt.SegmentBase {
	return &colibri_mgmt.SegmentBase{
		ID:    newID(),
		Index: idx,
	}
}

func newTestE2EBase(idx uint8) *colibri_mgmt.E2EBase {
	return &colibri_mgmt.E2EBase{
		ID:    newE2EID(),
		Index: idx,
	}
}

func newSetup() *colibri_mgmt.SegmentSetup {
	return &colibri_mgmt.SegmentSetup{
		Base:     newTestBase(1),
		MinBW:    1,
		MaxBW:    2,
		SplitCls: 3,
		StartProps: colibri_mgmt.PathEndProps{
			Local:    true,
			Transfer: false,
		},
		EndProps: colibri_mgmt.PathEndProps{
			Local:    false,
			Transfer: true,
		},
		InfoField: xtest.MustParseHexString("16ebdb4f0d042500"),
		AllocationTrail: []*colibri_mgmt.AllocationBead{
			{
				AllocBW: 5,
				MaxBW:   6,
			},
		},
	}
}

func newTelesSetup() *colibri_mgmt.SegmentTelesSetup {
	return &colibri_mgmt.SegmentTelesSetup{
		Setup:  newSetup(),
		BaseID: newID(),
	}
}

func newIndexConfirmation() *colibri_mgmt.SegmentIndexConfirmation {
	return &colibri_mgmt.SegmentIndexConfirmation{
		Base:  newTestBase(2),
		State: proto.ReservationIndexState_active,
	}
}

func newCleanup() *colibri_mgmt.SegmentCleanup {
	return &colibri_mgmt.SegmentCleanup{
		Base: newTestBase(1),
	}
}

func newE2ESetup() *colibri_mgmt.E2ESetup {
	return &colibri_mgmt.E2ESetup{
		Base:  newTestE2EBase(1),
		Token: xtest.MustParseHexString("16ebdb4f0d042500003f001002bad1ce003f001002facade"),
	}
}

func newE2ECleanup() *colibri_mgmt.E2ECleanup {
	return &colibri_mgmt.E2ECleanup{
		Base: newTestE2EBase(1),
	}
}

// new path with one segment consisting on 3 hopfields: (0,2)->(1,2)->(1,0)
func newPath() *spath.Path {
	path := &spath.Path{
		InfOff: 0,
		HopOff: spath.InfoFieldLength + spath.HopFieldLength, // second hop field
		Raw:    make([]byte, spath.InfoFieldLength+3*spath.HopFieldLength),
	}
	inf := spath.InfoField{ConsDir: true, ISD: 1, Hops: 3}
	inf.Write(path.Raw)

	hf := &spath.HopField{ConsEgress: 2}
	hf.Write(path.Raw[spath.InfoFieldLength:])
	hf = &spath.HopField{ConsIngress: 1, ConsEgress: 2}
	hf.Write(path.Raw[spath.InfoFieldLength+spath.HopFieldLength:])
	hf = &spath.HopField{ConsIngress: 1}
	hf.Write(path.Raw[spath.InfoFieldLength+spath.HopFieldLength*2:])

	return path
}

func checkRequest(t *testing.T, segSetup *colibri_mgmt.SegmentSetup, r *segment.SetupReq,
	ts time.Time) {

	t.Helper()
	require.Equal(t, (*segment.Reservation)(nil), r.Reservation)
	require.Equal(t, ts, r.Timestamp)
	checkIDs(t, segSetup.Base.ID, &r.ID)
	require.Equal(t, segSetup.Base.Index, uint8(r.Index))
	require.Equal(t, segSetup.MinBW, uint8(r.MinBW))
	require.Equal(t, segSetup.MaxBW, uint8(r.MaxBW))
	require.Equal(t, segSetup.SplitCls, uint8(r.SplitCls))
	require.Equal(t, reservation.NewPathEndProps(
		segSetup.StartProps.Local, segSetup.StartProps.Transfer,
		segSetup.EndProps.Local, segSetup.EndProps.Transfer), r.PathProps)
	require.Len(t, r.AllocTrail, len(segSetup.AllocationTrail))
	for i := range segSetup.AllocationTrail {
		require.Equal(t, segSetup.AllocationTrail[i].AllocBW, uint8(r.AllocTrail[i].AllocBW))
		require.Equal(t, segSetup.AllocationTrail[i].MaxBW, uint8(r.AllocTrail[i].MaxBW))
	}
}

func checkIDs(t *testing.T, ctrlID *colibri_mgmt.SegmentReservationID, id *reservation.SegmentID) {
	t.Helper()
	expectedID := append(ctrlID.ASID, ctrlID.Suffix...)
	require.Equal(t, expectedID, id.ToRaw())
}

func checkE2EIDs(t *testing.T, ctrlID *colibri_mgmt.E2EReservationID, id *reservation.E2EID) {
	t.Helper()
	expectedID := append(ctrlID.ASID, ctrlID.Suffix...)
	require.Equal(t, expectedID, id.ToRaw())
}
