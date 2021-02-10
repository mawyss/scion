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

	"github.com/scionproto/scion/go/cs/reservation/e2e"
	"github.com/scionproto/scion/go/cs/reservation/segment"
	"github.com/scionproto/scion/go/cs/reservation/segmenttest"
	"github.com/scionproto/scion/go/cs/reservationstorage/backend"
	"github.com/scionproto/scion/go/lib/colibri/reservation"
	"github.com/scionproto/scion/go/lib/util"
	"github.com/scionproto/scion/go/lib/xtest"
)

// AddSegmentReservation adds `count` segment reservation ff00:0:1-00000001 to `db`.
func AddSegmentReservation(t testing.TB, db backend.DB, count int) {
	t.Helper()
	ctx := context.Background()

	r := newTestSegmentReservation(t) // ff00:0:1, suffix=1
	for i := 0; i < count; i++ {
		r.Path = segmenttest.NewPathFromComponents(0, "1-ff00:0:1", i, 1, "1-ff00:0:2", 0)
		r.Indices = segment.Indices{}
		err := db.NewSegmentRsv(ctx, r)
		require.NoError(t, err, "iteration i = %d", i)
	}
}

func AddE2EReservation(t testing.TB, db backend.DB, count int) {
	t.Helper()
	ctx := context.Background()

	for i := 0; i < count; i++ {
		r := newTestE2EReservation(t)

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

func newTestSegmentReservation(t testing.TB) *segment.Reservation {
	t.Helper()
	r := segment.NewReservation()
	r.Path = segment.ReservationTransparentPath{}
	r.ID.ASID = xtest.MustParseAS("ff00:0:1")
	r.Ingress = 0
	r.Egress = 1
	r.TrafficSplit = 3
	r.PathEndProps = reservation.EndLocal | reservation.StartLocal
	expTime := util.SecsToTime(1)
	_, err := r.NewIndexAtSource(expTime, 1, 3, 2, 5, reservation.CorePath)
	require.NoError(t, err)
	err = r.SetIndexConfirmed(0)
	require.NoError(t, err)
	return r
}

// newTestE2EReservation adds an E2E reservation, that uses three segment reservations on
// ff00:1:1-00000001, ff00:2:2-00000002 and ff00:3:3-00000003 .
// The E2E reservation is transit on the first leg.
func newTestE2EReservation(t testing.TB) *e2e.Reservation {
	t.Helper()

	rsv := &e2e.Reservation{
		ID: *e2eIDFromRaw(t, "ff00:1:1", "00000000000000000001"),
		SegmentReservations: []*segment.Reservation{
			newTestSegmentReservation(t),
			// newTestSegmentReservation(t), // TODO change 2nd and 3rd seg rsvs.
			// newTestSegmentReservation(t),
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

func segmentIDFromRaw(t testing.TB, rawID string) *reservation.SegmentID {
	t.Helper()
	ID, err := reservation.SegmentIDFromRaw(xtest.MustParseHexString(rawID))
	require.NoError(t, err)
	return ID
}

func e2eIDFromRaw(t testing.TB, ASID, suffix string) *reservation.E2EID {
	t.Helper()
	id, err := reservation.NewE2EID(xtest.MustParseAS(ASID), xtest.MustParseHexString(suffix))
	require.NoError(t, err)
	return id
}
