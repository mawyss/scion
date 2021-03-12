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

package stateful

import (
	"context"
	"math"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	base "github.com/scionproto/scion/go/cs/reservation"
	"github.com/scionproto/scion/go/cs/reservation/segment"
	"github.com/scionproto/scion/go/cs/reservationstorage/backend"
	"github.com/scionproto/scion/go/cs/reservationstorage/backend/mock_backend"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/colibri/reservation"
	"github.com/scionproto/scion/go/lib/util"
	"github.com/scionproto/scion/go/lib/xtest"
)

func TestSumMaxBlockedBW(t *testing.T) {
	cases := map[string]struct {
		blockedBW uint64
		rsvsFcn   func() []*segment.Reservation
		excludeID string
	}{
		"empty": {
			blockedBW: 0,
			rsvsFcn: func() []*segment.Reservation {
				return nil
			},
			excludeID: "ff0000010001beefcafe",
		},
		"one reservation": {
			blockedBW: reservation.BWCls(5).ToKbps(),
			rsvsFcn: func() []*segment.Reservation {
				rsv := testNewRsv(t, "ff00:1:1", "01234567", 1, 2, 5, 5, 5)
				_, err := rsv.NewIndexAtSource(util.SecsToTime(3), 1, 1, 1, 1, reservation.CorePath)
				require.NoError(t, err)
				_, err = rsv.NewIndexAtSource(util.SecsToTime(3), 1, 1, 1, 1, reservation.CorePath)
				require.NoError(t, err)
				return []*segment.Reservation{rsv}
			},
			excludeID: "ff0000010001beefcafe",
		},
		"one reservation but excluded": {
			blockedBW: 0,
			rsvsFcn: func() []*segment.Reservation {
				rsv := testNewRsv(t, "ff00:1:1", "beefcafe", 1, 2, 5, 5, 5)
				_, err := rsv.NewIndexAtSource(util.SecsToTime(3), 1, 1, 1, 1, reservation.CorePath)
				require.NoError(t, err)
				_, err = rsv.NewIndexAtSource(util.SecsToTime(3), 1, 1, 1, 1, reservation.CorePath)
				require.NoError(t, err)
				return []*segment.Reservation{rsv}
			},
			excludeID: "ff0000010001beefcafe",
		},
		"many reservations": {
			blockedBW: 309, // 181 + 128
			rsvsFcn: func() []*segment.Reservation {
				rsv := testNewRsv(t, "ff00:1:1", "beefcafe", 1, 2, 5, 5, 5)
				_, err := rsv.NewIndexAtSource(util.SecsToTime(3), 1, 17, 7, 1,
					reservation.CorePath)
				require.NoError(t, err)
				rsvs := []*segment.Reservation{rsv}

				rsv = testNewRsv(t, "ff00:1:1", "01234567", 1, 2, 5, 5, 5)
				_, err = rsv.NewIndexAtSource(util.SecsToTime(3), 1, 8, 8, 1, reservation.CorePath)
				require.NoError(t, err)
				_, err = rsv.NewIndexAtSource(util.SecsToTime(3), 1, 7, 7, 1, reservation.CorePath)
				require.NoError(t, err)
				rsvs = append(rsvs, rsv)

				rsv = testNewRsv(t, "ff00:1:2", "01234567", 1, 2, 5, 5, 5)
				_, err = rsv.NewIndexAtSource(util.SecsToTime(2), 1, 7, 7, 1, reservation.CorePath)
				require.NoError(t, err)
				rsvs = append(rsvs, rsv)

				return rsvs
			},
			excludeID: "ff0000010001beefcafe",
		},
	}

	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			excludedID, err := reservation.SegmentIDFromRaw(xtest.MustParseHexString(tc.excludeID))
			require.NoError(t, err)
			sum := sumMaxBlockedBW(tc.rsvsFcn(), *excludedID)
			require.Equal(t, tc.blockedBW, sum)
		})
	}
}

func TestAvailableBW(t *testing.T) {
	req := newTestRequest(t, 1, 2, 5, 7)

	cases := map[string]struct {
		availBW uint64
		delta   float64
		req     *segment.SetupReq
		setupDB func(db *mock_backend.MockDB)
	}{
		"empty DB": {
			availBW: 1024,
			delta:   1,
			req:     req,
			setupDB: func(db *mock_backend.MockDB) {
				db.EXPECT().GetInterfaceUsageIngress(gomock.Any(), gomock.Any()).Return(
					uint64(0), nil)
				db.EXPECT().GetInterfaceUsageEgress(gomock.Any(), gomock.Any()).Return(
					uint64(0), nil)
				db.EXPECT().GetSegmentRsvFromID(gomock.Any(), &req.ID).Return(
					nil, nil)
			},
		},
		"this reservation in DB": {
			// as the only reservation in DB has the same ID as the request, the availableBW
			// function should return the same value as with an empty DB.
			availBW: 1024,
			delta:   1,
			req:     req,
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{testNewRsv(t, "ff00:1:1", "beefcafe", 1, 2, 5, 5, 5)}
				db.EXPECT().GetInterfaceUsageIngress(gomock.Any(), gomock.Any()).Return(
					reservation.BWCls(5).ToKbps(), nil)
				db.EXPECT().GetInterfaceUsageEgress(gomock.Any(), gomock.Any()).Return(
					reservation.BWCls(5).ToKbps(), nil)
				db.EXPECT().GetSegmentRsvFromID(gomock.Any(), &req.ID).Return(
					rsvs[0], nil)
			},
		},
		"other reservation in DB": {
			availBW: 1024 - 64,
			delta:   1,
			req:     req,
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "beefcafe", 1, 2, 5, 5, 5),
					testNewRsv(t, "ff00:1:2", "beefcafe", 1, 2, 5, 5, 5),
				}
				db.EXPECT().GetInterfaceUsageIngress(gomock.Any(), gomock.Any()).Return(
					reservation.BWCls(5).ToKbps()*2, nil)
				db.EXPECT().GetInterfaceUsageEgress(gomock.Any(), gomock.Any()).Return(
					reservation.BWCls(5).ToKbps()*2, nil)
				db.EXPECT().GetSegmentRsvFromID(gomock.Any(), &req.ID).Return(
					rsvs[0], nil)
			},
		},
		"change delta": {
			availBW: (1024 - 64) / 2,
			delta:   .5,
			req:     req,
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "beefcafe", 1, 2, 5, 5, 5),
					testNewRsv(t, "ff00:1:2", "beefcafe", 1, 2, 5, 5, 5),
				}
				db.EXPECT().GetInterfaceUsageIngress(gomock.Any(), gomock.Any()).Return(
					reservation.BWCls(5).ToKbps()*2, nil)
				db.EXPECT().GetInterfaceUsageEgress(gomock.Any(), gomock.Any()).Return(
					reservation.BWCls(5).ToKbps()*2, nil)
				db.EXPECT().GetSegmentRsvFromID(gomock.Any(), &req.ID).Return(
					rsvs[0], nil)
			},
		},
	}
	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			db, finish := newTestDB(t)
			adm := newTestAdmitter(t)
			defer finish()

			adm.Delta = tc.delta
			ctx := context.Background()
			tc.setupDB(db.(*mock_backend.MockDB))
			avail, err := adm.availableBW(ctx, db, tc.req)
			require.NoError(t, err)
			require.Equal(t, tc.availBW, avail)
		})
	}
}

func TestTubeRatio(t *testing.T) {
	cases := map[string]struct {
		tubeRatio      float64
		req            *segment.SetupReq
		setupDB        func(db *mock_backend.MockDB)
		globalCapacity uint64
		interfaces     []uint16
	}{
		"empty": {
			tubeRatio: 1,
			req:       newTestRequest(t, 1, 2, 5, 5),
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{}
				req := newTestRequest(t, 1, 2, 5, 5)
				prepareMockForTubeRatio(db, rsvs, req, 1024*1024)
			},
			globalCapacity: 1024 * 1024,
			interfaces:     []uint16{1, 2, 3},
		},
		"one source, one ingress": {
			tubeRatio: 1,
			req:       newTestRequest(t, 1, 2, 5, 5),
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "00000001", 1, 2, 5, 5, 5),
				}
				req := newTestRequest(t, 1, 2, 5, 5)
				prepareMockForTubeRatio(db, rsvs, req, 1024*1024)
			},
			globalCapacity: 1024 * 1024,
			interfaces:     []uint16{1, 2, 3},
		},
		"one source, two ingress": {
			tubeRatio: .5,
			req:       newTestRequest(t, 1, 2, 3, 3), // 64Kbps
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "00000001", 1, 2, 3, 3, 3), // 64Kbps
					testNewRsv(t, "ff00:1:1", "00000002", 3, 2, 5, 5, 5), // 128Kbps
				}
				req := newTestRequest(t, 1, 2, 3, 3)
				prepareMockForTubeRatio(db, rsvs, req, 1024*1024)
			},
			globalCapacity: 1024 * 1024,
			interfaces:     []uint16{1, 2, 3},
		},
		"two sources, request already present": {
			tubeRatio: .5,
			req:       newTestRequest(t, 1, 2, 5, 5),
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "beefcafe", 1, 2, 5, 9, 9), // will be ignored
					testNewRsv(t, "ff00:1:1", "00000002", 3, 2, 5, 5, 5),
				}
				req := newTestRequest(t, 1, 2, 5, 5)
				prepareMockForTubeRatio(db, rsvs, req, 1024*1024)

			},
			globalCapacity: 1024 * 1024,
			interfaces:     []uint16{1, 2, 3},
		},
		"multiple sources, multiple ingress": {
			tubeRatio: .75,
			req:       newTestRequest(t, 1, 2, 5, 5),
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "00000001", 1, 2, 5, 5, 5),
					testNewRsv(t, "ff00:1:2", "00000001", 1, 2, 5, 5, 5),
					testNewRsv(t, "ff00:1:1", "00000002", 3, 2, 5, 5, 5),
				}
				req := newTestRequest(t, 1, 2, 5, 5)
				prepareMockForTubeRatio(db, rsvs, req, 1024*1024)
			},
			globalCapacity: 1024 * 1024,
			interfaces:     []uint16{1, 2, 3},
		},
		"exceeding ingress capacity": {
			tubeRatio: 10. / 13., // 10 / (10 + 0 + 3)
			req:       newTestRequest(t, 1, 2, 5, 5),
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "00000001", 1, 2, 5, 5, 5),
					testNewRsv(t, "ff00:1:2", "00000001", 1, 2, 5, 5, 5),
					testNewRsv(t, "ff00:1:1", "00000002", 3, 2, 5, 5, 5),
				}
				req := newTestRequest(t, 1, 2, 5, 5)
				prepareMockForTubeRatio(db, rsvs, req, 10)
			},
			globalCapacity: 10,
			interfaces:     []uint16{1, 2, 3},
		},
		"with many other irrelevant reservations": {
			tubeRatio: .75,
			req:       newTestRequest(t, 1, 2, 5, 5),
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "00000001", 1, 2, 5, 5, 5),
					testNewRsv(t, "ff00:1:2", "00000001", 1, 2, 5, 5, 5), // 128 Kbps
					testNewRsv(t, "ff00:1:1", "00000002", 3, 2, 5, 5, 5),
					testNewRsv(t, "ff00:1:3", "00000001", 4, 5, 5, 9, 9),
					testNewRsv(t, "ff00:1:3", "00000002", 4, 5, 5, 9, 9),
					testNewRsv(t, "ff00:1:4", "00000001", 5, 4, 5, 9, 9),
					testNewRsv(t, "ff00:1:4", "00000002", 5, 4, 5, 9, 9),
				}
				req := newTestRequest(t, 1, 2, 5, 5)
				prepareMockForTubeRatio(db, rsvs, req, 1024*1024)
			},
			globalCapacity: 1024 * 1024,
			interfaces:     []uint16{1, 2, 3, 4, 5},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			db, finish := newTestDB(t)
			adm := newTestAdmitter(t)
			defer finish()

			adm.Capacities = &testCapacities{
				Cap:    tc.globalCapacity,
				Ifaces: tc.interfaces,
			}
			tc.setupDB(db.(*mock_backend.MockDB))

			ctx := context.Background()
			ratio, err := adm.tubeRatio(ctx, db, tc.req)
			require.NoError(t, err)
			require.Equal(t, tc.tubeRatio, ratio)
		})
	}
}

func TestLinkRatio(t *testing.T) {
	cases := map[string]struct {
		linkRatio float64
		req       *segment.SetupReq
		setupDB   func(db *mock_backend.MockDB)
	}{
		"empty": {
			linkRatio: 1.,
			req:       testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 5, 5),
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{}
				req := testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 5, 5)
				prepareMockForLinkRatio(db, rsvs, req, 1024*1024)
			},
		},
		"same request": {
			linkRatio: 1.,
			req:       testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 5, 5),
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "beefcafe", 1, 2, 5, 5, 5),
				}
				req := testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 5, 5)
				prepareMockForLinkRatio(db, rsvs, req, 1024*1024)
			},
		},
		"same source": {
			linkRatio: .5,
			req:       testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 5, 5),
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "beefcafe", 1, 2, 5, 5, 5),
					testNewRsv(t, "ff00:1:1", "00000001", 1, 2, 5, 5, 5),
				}
				req := testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 5, 5)
				prepareMockForLinkRatio(db, rsvs, req, 1024*1024)
			},
		},
		"different sources": {
			linkRatio: 1. / 3.,
			req:       testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 5, 5),
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:2", "00000001", 1, 2, 5, 5, 5),
					testNewRsv(t, "ff00:1:3", "00000001", 1, 2, 5, 5, 5),
				}
				req := testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 5, 5)
				prepareMockForLinkRatio(db, rsvs, req, 1024*1024)
			},
		},
		"different egress interface": {
			linkRatio: 1., // 64 / 64  => srcAlloc(ff00:1:1, 1, 2) = 0 + prevBW = 0 + 64 = 64
			req:       testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 5, 5),
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "00000001", 1, 3, 5, 5, 5),
				}
				req := testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 5, 5)
				prepareMockForLinkRatio(db, rsvs, req, 1024*1024)
			},
		},
		"smaller prevBW": {
			linkRatio: 1. / 3.,
			req:       testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 3, 3), // 64 Kbps
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "00000001", 1, 2, 5, 5, 5), // 128 Kbps
				}
				req := testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 3, 3)
				prepareMockForLinkRatio(db, rsvs, req, 1024*1024)
			},
		},
		"bigger prevBW": {
			linkRatio: 2. / 3.,
			req:       testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 7, 7), // 256 Kbps
			setupDB: func(db *mock_backend.MockDB) {
				rsvs := []*segment.Reservation{
					testNewRsv(t, "ff00:1:1", "00000001", 1, 2, 5, 5, 5), // 128 Kbps
				}
				req := testAddAllocTrail(newTestRequest(t, 1, 2, 5, 5), 7, 7)
				prepareMockForLinkRatio(db, rsvs, req, 1024*1024)
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			db, finish := newTestDB(t)
			adm := newTestAdmitter(t)
			defer finish()

			adm.Capacities = &testCapacities{
				Cap:    1024 * 1024,
				Ifaces: []uint16{1, 2, 3},
			}
			tc.setupDB(db.(*mock_backend.MockDB))

			ctx := context.Background()
			linkRatio, err := adm.linkRatio(ctx, db, tc.req)
			require.NoError(t, err)
			require.Equal(t, tc.linkRatio, linkRatio)
		})
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

func newTestDB(t *testing.T) (backend.DB, func()) {
	mctlr := gomock.NewController(t)
	db := mock_backend.NewMockDB(mctlr)

	return db, mctlr.Finish
}

func newTestAdmitter(t *testing.T) *StatefulAdmission {
	return &StatefulAdmission{
		Capacities: &testCapacities{
			Cap:    1024, // 1MBps
			Ifaces: []uint16{1, 2},
		},
		Delta: 1,
	}
}

// newTestRequest creates a request ID ff00:1:1 beefcafe
func newTestRequest(t *testing.T, ingress, egress uint16,
	minBW, maxBW reservation.BWCls) *segment.SetupReq {

	ID, err := reservation.SegmentIDFromRaw(xtest.MustParseHexString("ff0000010001beefcafe"))
	require.NoError(t, err)
	return &segment.SetupReq{
		Request: segment.Request{
			RequestMetadata: base.RequestMetadata{},
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

func testNewRsv(t *testing.T, srcAS string, suffix string, ingress, egress uint16,
	minBW, maxBW, allocBW reservation.BWCls) *segment.Reservation {

	ID, err := reservation.NewSegmentID(xtest.MustParseAS(srcAS),
		xtest.MustParseHexString(suffix))
	require.NoError(t, err)
	rsv := &segment.Reservation{
		ID: *ID,
		Indices: segment.Indices{
			segment.Index{
				Idx:        10,
				Expiration: util.SecsToTime(2),
				MinBW:      minBW,
				MaxBW:      maxBW,
				AllocBW:    allocBW,
			},
		},
		Ingress:      ingress,
		Egress:       egress,
		PathType:     reservation.UpPath,
		PathEndProps: reservation.StartLocal | reservation.EndLocal | reservation.EndTransfer,
		TrafficSplit: 2,
	}
	err = rsv.SetIndexConfirmed(10)
	require.NoError(t, err)
	err = rsv.SetIndexActive(10)
	require.NoError(t, err)
	return rsv
}

// testAddAllocTrail adds an allocation trail to a reservation. The beads parameter represents
// the trail like: alloc0,max0,alloc1,max1,...
func testAddAllocTrail(req *segment.SetupReq, beads ...reservation.BWCls) *segment.SetupReq {
	if len(beads)%2 != 0 {
		panic("the beads must be even")
	}
	for i := 0; i < len(beads); i += 2 {
		beads := reservation.AllocationBead{
			AllocBW: beads[i],
			MaxBW:   beads[i+1],
		}
		req.AllocTrail = append(req.AllocTrail, beads)
	}
	return req
}

func getMaxBWPerSource(t *testing.T, rsvs []*segment.Reservation, skipASID, skipSuffix string) (
	map[addr.AS]uint64, error) {

	skipRsv, err := reservation.NewSegmentID(xtest.MustParseAS(skipASID),
		xtest.MustParseHexString(skipSuffix))
	require.NoError(t, err)
	maxBWPerSrc := make(map[addr.AS]uint64)
	for _, r := range rsvs {
		if r.ID != *skipRsv {
			maxBWPerSrc[r.ID.ASID] += r.MaxBlockedBW()
		}
	}
	return maxBWPerSrc, nil
}

type sourceIngressEgress struct {
	Source  addr.AS
	Ingress uint16
	Egress  uint16
}
type sourceState struct {
	SrcDem   uint64
	SrcAlloc uint64
}

func prepareForMock(rsvs []*segment.Reservation, req *segment.SetupReq, globalCapacity uint64) (
	sameIDAsRequest *segment.Reservation,
	sourceStateMap map[sourceIngressEgress]sourceState,
	inMap map[addr.AS]map[uint16]uint64,
	egMap map[addr.AS]map[uint16]uint64,
	transitDem map[uint16]uint64,
	transitAlloc uint64) {

	// source,ingress,egress map
	sourceStateMap = make(map[sourceIngressEgress]sourceState)

	inMap = make(map[addr.AS]map[uint16]uint64) // inDem for each source
	egMap = make(map[addr.AS]map[uint16]uint64)

	// transitDem is map[ingress] to req.Egress
	transitDem = make(map[uint16]uint64)

	// transitAlloc goes from req.Ingress to req.Egress

	for _, r := range rsvs {
		if r.ID == req.ID {
			sameIDAsRequest = r
		}

		key := sourceIngressEgress{
			Source:  r.ID.ASID,
			Ingress: r.Ingress,
			Egress:  r.Egress,
		}
		state := sourceStateMap[key]
		state.SrcDem += minBW(r.MaxRequestedBW(), globalCapacity)
		state.SrcAlloc += r.MaxBlockedBW()
		sourceStateMap[key] = state

		if inMap[r.ID.ASID] == nil {
			inMap[r.ID.ASID] = make(map[uint16]uint64)
		}
		inMap[r.ID.ASID][r.Ingress] += minBW(r.MaxRequestedBW(), globalCapacity)
		if egMap[r.ID.ASID] == nil {
			egMap[r.ID.ASID] = make(map[uint16]uint64)
		}
		egMap[r.ID.ASID][r.Egress] += minBW(r.MaxRequestedBW(), globalCapacity)
	}
	// the transitDem and transitAlloc need the scale factors to be computed:
	for _, r := range rsvs {
		inScalFctr := float64(minBW(inMap[r.ID.ASID][r.Ingress], globalCapacity)) /
			float64(inMap[r.ID.ASID][r.Ingress])
		egScalFctr := float64(minBW(egMap[r.ID.ASID][r.Egress], globalCapacity)) /
			float64(egMap[r.ID.ASID][r.Egress])
		key := sourceIngressEgress{
			Source:  r.ID.ASID,
			Ingress: r.Ingress,
			Egress:  r.Egress,
		}
		state := sourceStateMap[key]

		if r.Egress == req.Egress {
			transitDem[r.Ingress] += uint64(math.Round(float64(state.SrcDem) *
				math.Min(inScalFctr, egScalFctr)))
			if r.Ingress == req.Ingress {
				transitAlloc += r.MaxBlockedBW()
			}
		}
	}
	return sameIDAsRequest, sourceStateMap, inMap, egMap, transitDem, transitAlloc
}

func prepareMockForTubeRatio(db *mock_backend.MockDB, rsvs []*segment.Reservation, req *segment.SetupReq,
	globalCapacity uint64) {

	sameIDAsRequest, sourceStateMap, inMap, egMap, transitDem, transitAlloc :=
		prepareForMock(rsvs, req, globalCapacity)
	_ = transitAlloc

	db.EXPECT().GetSegmentRsvFromID(gomock.Any(), &req.ID).AnyTimes().Return(sameIDAsRequest, nil)

	db.EXPECT().GetTransitDem(gomock.Any(), gomock.Any(), req.Egress).AnyTimes().
		DoAndReturn(
			func(_ context.Context, ingress, _ uint16) (uint64, error) {
				return transitDem[ingress], nil
			})

	db.EXPECT().GetSourceState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(
			func(_ context.Context, source addr.AS, ingress, egress uint16) (uint64, uint64, error) {
				key := sourceIngressEgress{
					Source:  source,
					Ingress: ingress,
					Egress:  egress,
				}
				state := sourceStateMap[key]
				return state.SrcDem, state.SrcAlloc, nil
			})

	db.EXPECT().GetInDemand(gomock.Any(), req.ID.ASID, gomock.Any()).AnyTimes().
		DoAndReturn(
			func(_ context.Context, _ addr.AS, ingress uint16) (uint64, error) {
				return inMap[req.ID.ASID][ingress], nil
			})

	db.EXPECT().GetEgDemand(gomock.Any(), req.ID.ASID, gomock.Any()).AnyTimes().
		DoAndReturn(
			func(_ context.Context, _ addr.AS, egress uint16) (uint64, error) {
				return egMap[req.ID.ASID][egress], nil
			})
}

func prepareMockForLinkRatio(db *mock_backend.MockDB, rsvs []*segment.Reservation, req *segment.SetupReq,
	globalCapacity uint64) {

	sameIDAsRequest, sourceStateMap, inMap, egMap, transitDem, transitAlloc :=
		prepareForMock(rsvs, req, globalCapacity)
	_ = sameIDAsRequest
	_ = transitDem

	db.EXPECT().GetSegmentRsvFromID(gomock.Any(), &req.ID).AnyTimes().Return(sameIDAsRequest, nil)

	db.EXPECT().GetTransitAlloc(gomock.Any(), req.Ingress, req.Egress).AnyTimes().
		Return(transitAlloc, nil)

	db.EXPECT().GetSourceState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(
			func(_ context.Context, source addr.AS, ingress, egress uint16) (uint64, uint64, error) {
				key := sourceIngressEgress{
					Source:  source,
					Ingress: ingress,
					Egress:  egress,
				}
				state := sourceStateMap[key]
				return state.SrcDem, state.SrcAlloc, nil
			})

	db.EXPECT().GetInDemand(gomock.Any(), req.ID.ASID, gomock.Any()).AnyTimes().
		DoAndReturn(
			func(_ context.Context, _ addr.AS, ingress uint16) (uint64, error) {
				return inMap[req.ID.ASID][ingress], nil
			})

	db.EXPECT().GetEgDemand(gomock.Any(), req.ID.ASID, gomock.Any()).AnyTimes().
		DoAndReturn(
			func(_ context.Context, _ addr.AS, egress uint16) (uint64, error) {
				return egMap[req.ID.ASID][egress], nil
			})
}
