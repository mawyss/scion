// Copyright 2021 ETH Zurich, Anapaya Systems
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

package reservationstore_test

import (
	"context"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	base "github.com/scionproto/scion/go/cs/reservation"
	"github.com/scionproto/scion/go/cs/reservation/e2e"
	"github.com/scionproto/scion/go/cs/reservation/segment"
	"github.com/scionproto/scion/go/cs/reservation/segment/admission"
	stateless "github.com/scionproto/scion/go/cs/reservation/segment/admission/impl"
	"github.com/scionproto/scion/go/cs/reservation/segmenttest"
	"github.com/scionproto/scion/go/cs/reservation/sqlite"
	"github.com/scionproto/scion/go/cs/reservation/test"
	"github.com/scionproto/scion/go/cs/reservationstorage/backend"
	"github.com/scionproto/scion/go/lib/colibri/reservation"
	"github.com/scionproto/scion/go/lib/util"
	"github.com/scionproto/scion/go/lib/xtest"
)

// AddSegmentReservation adds `count` segment reservation ASID-newsuffix to `db`.
func AddSegmentReservation(t testing.TB, db backend.DB, ASID string, count int) {
	t.Helper()
	ctx := context.Background()

	r := newTestSegmentReservation(t, ASID) // the suffix will be overwritten
	for i := 0; i < count; i++ {
		r.Path = segmenttest.NewPathFromComponents(0, "1-"+ASID, i, 1, "1-ff00:0:2", 0)
		err := db.NewSegmentRsv(ctx, r)
		require.NoError(t, err, "iteration i = %d", i)
	}
}

// AddE2EReservation ads `count` E2E reservations to the DB.
func AddE2EReservation(t testing.TB, db backend.DB, ASID string, count int) {
	t.Helper()
	ctx := context.Background()

	for i := 0; i < count; i++ {
		r := newTestE2EReservation(t, ASID)

		auxBuff := make([]byte, 8)
		binary.BigEndian.PutUint64(auxBuff, uint64(i+1))
		copy(r.ID.Suffix[2:], auxBuff)
		for _, seg := range r.SegmentReservations {
			err := db.PersistSegmentRsv(ctx, seg)
			require.NoError(t, err)
		}
		err := db.PersistE2ERsv(ctx, r)
		require.NoError(t, err)
	}
}

type testCapacities struct {
	Cap    uint64
	Ifaces []uint16
}

var _ base.Capacities = (*testCapacities)(nil)

func (c *testCapacities) IngressInterfaces() []uint16           { return c.Ifaces }
func (c *testCapacities) EgressInterfaces() []uint16            { return c.Ifaces }
func (c *testCapacities) Capacity(from, to uint16) uint64       { return c.Cap }
func (c *testCapacities) CapacityIngress(ingress uint16) uint64 { return c.Cap }
func (c *testCapacities) CapacityEgress(egress uint16) uint64   { return c.Cap }

// newTestSegmentReservation creates a segment reservation
func newTestSegmentReservation(t testing.TB, ASID string) *segment.Reservation {
	t.Helper()
	r := segment.NewReservation()
	r.Path = segment.ReservationTransparentPath{}
	r.ID.ASID = xtest.MustParseAS(ASID)
	r.Ingress = 0
	r.Egress = 1
	r.TrafficSplit = 3
	r.PathEndProps = reservation.EndLocal | reservation.StartLocal
	expTime := util.SecsToTime(1)
	_, err := r.NewIndexAtSource(expTime, 1, 3, 2, 5, reservation.CorePath)
	require.NoError(t, err)
	err = r.SetIndexConfirmed(0)
	require.NoError(t, err)
	err = r.SetIndexActive(0)
	require.NoError(t, err)
	return r
}

// newTestE2EReservation adds an E2E reservation, that uses three segment reservations on
// ASID-00000001, ff00:2:2-00000002 and ff00:3:3-00000003 .
// The E2E reservation is transit on the first leg.
func newTestE2EReservation(t testing.TB, ASID string) *e2e.Reservation {
	t.Helper()

	rsv := &e2e.Reservation{
		ID: *e2eIDFromRaw(t, ASID, "00000000000000000001"),
		SegmentReservations: []*segment.Reservation{
			newTestSegmentReservation(t, ASID),
		},
	}
	_, err := rsv.NewIndex(util.SecsToTime(1))
	require.NoError(t, err)
	return rsv
}

// newAllocationBeads (1,2,3,4) returns two beads {alloc: 1, max: 2}, {alloc:3, max:4}
func newAllocationBeads(beads ...reservation.BWCls) reservation.AllocationBeads {
	if len(beads)%2 != 0 {
		panic("must have an even number of parameters")
	}
	ret := make(reservation.AllocationBeads, len(beads)/2)
	for i := 0; i < len(beads); i += 2 {
		ret[i/2] = reservation.AllocationBead{AllocBW: beads[i], MaxBW: beads[i+1]}
	}
	return ret
}

func segmentIDFromRaw(t testing.TB, ASID, suffix string) *reservation.SegmentID {
	t.Helper()
	ID, err := reservation.NewSegmentID(xtest.MustParseAS(ASID), xtest.MustParseHexString(suffix))
	require.NoError(t, err)
	return ID
}

func e2eIDFromRaw(t testing.TB, ASID, suffix string) *reservation.E2EID {
	t.Helper()
	ID, err := reservation.NewE2EID(xtest.MustParseAS(ASID), xtest.MustParseHexString(suffix))
	require.NoError(t, err)
	return ID
}

func newDB(t testing.TB) backend.DB {
	t.Helper()
	db, err := sqlite.New("file::memory:")
	require.NoError(t, err)
	// db.SetMaxOpenConns(10)
	return db
}

func newCapacities() base.Capacities {
	return &testCapacities{
		Cap:    1024 * 1024, // 1GBps
		Ifaces: []uint16{1, 2},
	}
}

func newAdmitter(cap base.Capacities) admission.Admitter {
	admitter := &stateless.StatelessAdmission{
		Capacities: cap,
		Delta:      1,
	}
	return admitter
}

// newTestRequest creates a request ID ff00:1:1 beefcafe
func newTestSegmentRequest(t testing.TB, ASID string, ingress, egress uint16,
	minBW, maxBW reservation.BWCls) *segment.SetupReq {

	t.Helper()

	ID := segmentIDFromRaw(t, ASID, "beefcafe")
	path := test.NewTestPath()
	meta, err := base.NewRequestMetadata(path)
	require.NoError(t, err)
	return &segment.SetupReq{
		Request: segment.Request{
			RequestMetadata: *meta,
			ID:              *ID,
			Timestamp:       util.SecsToTime(1),
			Ingress:         ingress,
			Egress:          egress,
		},
		MinBW:     minBW,
		MaxBW:     maxBW,
		SplitCls:  2,
		PathProps: reservation.StartLocal | reservation.EndLocal,
	}
}

func newTestE2ESetupRequest(t testing.TB, ASID string) *e2e.SetupReq {
	t.Helper()

	ID := e2eIDFromRaw(t, ASID, "beefcafebeefcafebeef")
	path := test.NewTestPath()
	baseReq, err := e2e.NewRequest(util.SecsToTime(1), ID, 1, path)
	require.NoError(t, err)
	segmentRsvs := []reservation.SegmentID{
		*segmentIDFromRaw(t, ASID, "00000001"),
		*segmentIDFromRaw(t, "ff00:2:2", "beefcafe"),
		*segmentIDFromRaw(t, "ff00:3:3", "beefcafe"),
	}
	ASCountPerSegment := []uint8{4, 4, 5}
	trail := []reservation.BWCls{5, 5}
	setup, err := e2e.NewSetupRequest(baseReq, segmentRsvs, ASCountPerSegment, 5, trail)
	require.NoError(t, err)
	return setup
}

func newTestE2ESuccessReq(t testing.TB, ASID string) *e2e.SetupReqSuccess {
	token, err := reservation.TokenFromRaw(
		xtest.MustParseHexString("16ebdb4f0d042500003f001002bad1ce003f001002facade"))
	require.NoError(t, err)
	return &e2e.SetupReqSuccess{
		SetupReq: *newTestE2ESetupRequest(t, ASID),
		Token:    *token,
	}
}
