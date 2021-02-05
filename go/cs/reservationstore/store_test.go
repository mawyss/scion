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

package reservationstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	base "github.com/scionproto/scion/go/cs/reservation"
	"github.com/scionproto/scion/go/cs/reservation/segment"
	"github.com/scionproto/scion/go/cs/reservation/segment/admission"
	stateless "github.com/scionproto/scion/go/cs/reservation/segment/admission/impl"
	"github.com/scionproto/scion/go/cs/reservation/sqlite"
	"github.com/scionproto/scion/go/cs/reservation/test"
	"github.com/scionproto/scion/go/cs/reservationstorage"
	"github.com/scionproto/scion/go/cs/reservationstorage/backend"
	"github.com/scionproto/scion/go/lib/colibri/reservation"
	"github.com/scionproto/scion/go/lib/util"
	"github.com/scionproto/scion/go/lib/xtest"
)

func TestStore(t *testing.T) {
	var s reservationstorage.Store = &Store{}
	_ = s
}

func TestAdmitSegmentReservation(t *testing.T) {
	db := newDB(t)
	cap := newCapacities()
	admitter := newAdmitter(cap)
	s := NewStore(db, admitter)

	ctx := context.Background()
	req := newTestRequest(t, 1, 2, 5, 7)

	res, err := s.AdmitSegmentReservation(ctx, req)
	_ = res
	require.NoError(t, err)
}

func newDB(t *testing.T) backend.DB {
	t.Helper()
	db, err := sqlite.New("file::memory:")
	require.NoError(t, err)
	// db.SetMaxOpenConns(10)
	return db
}

func newCapacities() base.Capacities {
	return &testCapacities{
		Cap:    1024, // 1MBps
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

// newTestRequest creates a request ID ff00:1:1 beefcafe
func newTestRequest(t *testing.T, ingress, egress uint16,
	minBW, maxBW reservation.BWCls) *segment.SetupReq {

	ID, err := reservation.SegmentIDFromRaw(xtest.MustParseHexString("ff0000010001beefcafe"))
	require.NoError(t, err)
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
