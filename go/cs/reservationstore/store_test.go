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

package reservationstore_test

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"runtime/trace"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	base "github.com/scionproto/scion/go/cs/reservation"
	"github.com/scionproto/scion/go/cs/reservation/e2e"
	"github.com/scionproto/scion/go/cs/reservation/segment"
	"github.com/scionproto/scion/go/cs/reservation/segment/admission"
	stateless "github.com/scionproto/scion/go/cs/reservation/segment/admission/impl"
	"github.com/scionproto/scion/go/cs/reservation/sqlite"
	"github.com/scionproto/scion/go/cs/reservation/test"
	"github.com/scionproto/scion/go/cs/reservationstorage"
	"github.com/scionproto/scion/go/cs/reservationstorage/backend"
	"github.com/scionproto/scion/go/cs/reservationstore"
	"github.com/scionproto/scion/go/lib/colibri/reservation"
	"github.com/scionproto/scion/go/lib/util"
	"github.com/scionproto/scion/go/lib/xtest"
)

const REPS = 10

func TestStore(t *testing.T) {
	var s reservationstorage.Store = &reservationstore.Store{}
	_ = s
}

func TestDebugAdmitSegmentReservation(t *testing.T) {
	timeAdmitSegmentReservation(t, 1)
}

func TestDebugAdmitE2EReservation(t *testing.T) {
	timeAdmitE2EReservation(t, 1)
}

type performanceTestCase struct {
	TestName string
	Xmin     int
	Xmax     int
	Xstride  int
	Xlabel   string
	YLabels  []string // leave empty to generate default "sample 1", "sample 2", ...

	Repetitions int
	Function    func(*testing.T, int) time.Duration
	Filter      func(*testing.T, []time.Duration) []time.Duration

	DebugPrintProgress bool
}

func TestPerformanceCOLIBRI(t *testing.T) {

	if os.Getenv("COLIBRI_PERF_TESTS") == "" {
		t.SkipNow()
	}

	testCases := []performanceTestCase{
		// {
		// 	TestName:    "segmentAdmitManyASes",
		// 	Xmin:        1,
		// 	Xmax:        100,
		// 	Xstride:     1,
		// 	Xlabel:      "# ASes",
		// 	Repetitions: REPS,
		// 	Function:    timeAdmitSegmentReservation,
		// 	Filter:      identity,
		// },
		// {
		// 	TestName:    "segmentAdmitManyASesAverages",
		// 	Xmin:        1,
		// 	Xmax:        100,
		// 	Xstride:     1,
		// 	Xlabel:      "# ASes",
		// 	YLabels:     []string{"ave. µsecs"},
		// 	Repetitions: REPS,
		// 	Function:    timeAdmitSegmentReservation,
		// 	Filter:      getAverage,
		// },
		// {
		// 	TestName:    "segmentAdmitManyASesQuartiles",
		// 	Xmin:        1,
		// 	Xmax:        100,
		// 	Xstride:     1,
		// 	Xlabel:      "# ASes",
		// 	YLabels:     []string{"q1", "median", "q3"},
		// 	Repetitions: REPS,
		// 	Function:    timeAdmitSegmentReservation,
		// 	Filter:      getQuartiles,
		// },
		{
			TestName:           "e2eAdmitManyEndhosts",
			Xmin:               1,
			Xmax:               10000,
			Xstride:            100,
			Xlabel:             "# endhosts",
			Repetitions:        REPS,
			Function:           timeAdmitE2EReservation,
			Filter:             identity,
			DebugPrintProgress: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.TestName, func(t *testing.T) {
			t.Parallel()
			doPerformanceTest(t, tc)
		})
	}
}

func doPerformanceTest(t *testing.T, tc performanceTestCase) {
	Xs := delta(tc.Xmin, tc.Xmax, tc.Xstride)
	values := mapWithFunction(t,
		repeatWithFilter(t, tc.Function, tc.Repetitions, tc.Filter), Xs, tc.DebugPrintProgress)
	var columnTitles []string
	if len(tc.YLabels) > 0 {
		columnTitles = append([]string{tc.Xlabel}, tc.YLabels...)
	} else {
		columnTitles = make([]string, tc.Repetitions+1)
		columnTitles[0] = tc.Xlabel
		for i := 0; i < len(values[0]); i++ {
			columnTitles[i+1] = fmt.Sprintf("sample %d", i)
		}
	}
	toCSV(t, fmt.Sprintf("%s.csv", tc.TestName), columnTitles, Xs, values)
}

/////////////////////////////////////////

func delta(first, last, stride int) []int {
	ret := make([]int, (last-first)/stride+1)
	for i, j := first, 0; i <= last; i += stride {
		ret[j] = i
		j++
	}
	return ret
}

func mapWithFunction(t *testing.T,
	fn func(*testing.T, int) []time.Duration,
	xValues []int,
	printProgress bool) [][]time.Duration {

	values := make([][]time.Duration, len(xValues))
	// TODO worker pool here
	for i, x := range xValues {
		row := fn(t, x)
		values[i] = row
		if printProgress {
			t.Logf("done [%v] %v\n", time.Now(), x)
		}
	}
	return values
}

// returns a function applicable to map
func repeatWithFilter(t *testing.T,
	sampler func(*testing.T, int) time.Duration,
	repeatCount int,
	filter func(*testing.T, []time.Duration) []time.Duration) func(*testing.T, int) []time.Duration {

	ret := func(t *testing.T, x int) []time.Duration {
		samples := make([]time.Duration, repeatCount)
		// TODO worker pool here
		// type results struct{
		// 	iteration int
		// 	result time.Duration
		// }
		// const numJobs = 10
		// results := make (chan results, numJobs)
		for i := 0; i < repeatCount; i++ {
			samples[i] = sampler(t, x)
		}
		return filter(t, samples)
	}
	return ret
}

func identity(t *testing.T, values []time.Duration) []time.Duration {
	return values
}

// values contains N repetitions of scalars. Returns just 1 value in the slice.
func getAverage(t *testing.T, values []time.Duration) []time.Duration {
	require.Greater(t, len(values), 0)

	average := time.Duration(0)
	for i := 0; i < len(values); i++ {
		average += values[i]
	}
	average = time.Duration(average.Nanoseconds() / int64(len(values)))
	return []time.Duration{average}
}

// returns Q1, median, Q3 per dimension (three values in the slice).
func getQuartiles(t *testing.T, values []time.Duration) []time.Duration {
	medianFun := func(values []time.Duration) time.Duration {
		require.Greater(t, len(values), 0)
		l := len(values) / 2
		if len(values)%2 == 1 {
			return values[l]
		} else {
			return (values[l-1] + values[l]) / 2
		}
	}
	// sort in place
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	median := medianFun(values)
	var a, b int // indices for Q1 and Q3 segments
	l := len(values) / 2
	if len(values)%2 == 1 {
		a = l
		b = l + 1
	} else {
		a = l - 1
		b = l + 1
	}
	q1 := medianFun(values[:a])
	q3 := medianFun(values[b:])
	return []time.Duration{q1, median, q3}
}

func toCSV(t *testing.T, filename string, columnTitles []string, xValues []int, values [][]time.Duration) {
	width := len(values[0]) + 1
	require.Equal(t, len(xValues), len(values))
	require.Equal(t, len(columnTitles), width)

	f, err := os.Create(filename)
	require.NoError(t, err)
	defer f.Close()

	w := csv.NewWriter(f)
	err = w.Write(columnTitles)
	require.NoError(t, err)
	for i := 0; i < len(values); i++ {
		require.Equal(t, len(values[i]), width-1, "failed at row %d", i)
		row := make([]string, width)
		row[0] = fmt.Sprintf("%d", xValues[i])
		for j := 1; j < width; j++ {
			row[j] = fmt.Sprintf("%d", values[i][j-1].Microseconds())
		}
		w.Write(row)
	}
	w.Flush()
}

// func BenchmarkAdmitSegmentReservation10(b *testing.B) {
// 	benchmarkAdmitSegmentReservation(b, 10)
// }

// func BenchmarkAdmitSegmentReservation100(b *testing.B) {
// 	benchmarkAdmitSegmentReservation(b, 100)
// }

// func BenchmarkAdmitSegmentReservation1000(b *testing.B) {
// 	benchmarkAdmitSegmentReservation(b, 1000)
// }

// func BenchmarkAdmitSegmentReservation5000(b *testing.B) {
// 	benchmarkAdmitSegmentReservation(b, 5000)
// }

// func BenchmarkThings(b *testing.B) {
// 	b.ResetTimer()
// 	time.Sleep(5 * time.Second)
// }

// func BenchmarkAdmitSegmentReservation10000(b *testing.B) {
// 	benchmarkAdmitSegmentReservation(b, 10000)
// }

// func BenchmarkAdmitSegmentReservation100000(b *testing.B) {
// 	benchmarkAdmitSegmentReservation(b, 100000)
// }

func benchmarkAdmitSegmentReservation(b *testing.B, count int) {
	db := newDB(b)
	defer db.Close()

	cap := newCapacities()
	admitter := newAdmitter(cap)
	s := reservationstore.NewStore(db, admitter)

	AddSegmentReservation(b, db, count)
	ctx := context.Background()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		req := newTestSegmentRequest(b, 1, 2, 5, 7)
		// req.AllocTrail = newAllocationBeads(1, 2)
		trace.WithRegion(ctx, "AdmitSegmentReservation", func() {
			_, err := s.AdmitSegmentReservation(ctx, req)
			require.NoError(b, err, "iteration n = %d", n)
		})
	}
}

func timeAdmitSegmentReservation(t *testing.T, count int) time.Duration {
	db := newDB(t)
	cap := newCapacities()
	admitter := newAdmitter(cap)
	s := reservationstore.NewStore(db, admitter)

	AddSegmentReservation(t, db, count)
	ctx := context.Background()
	req := newTestSegmentRequest(t, 1, 2, 5, 7)

	t0 := time.Now()
	_, err := s.AdmitSegmentReservation(ctx, req)
	t1 := time.Since(t0)
	require.NoError(t, err)
	return t1
}

func timeAdmitE2EReservation(t *testing.T, count int) time.Duration {
	db := newDB(t)
	defer db.Close()

	cap := newCapacities()
	admitter := newAdmitter(cap)
	s := reservationstore.NewStore(db, admitter)

	AddSegmentReservation(t, db, 2)
	ctx := context.Background()
	req := newTestE2ERequest(t)
	token, err := reservation.TokenFromRaw(
		xtest.MustParseHexString("16ebdb4f0d042500003f001002bad1ce003f001002facade"))
	require.NoError(t, err)
	successfulReq := &e2e.SetupReqSuccess{
		SetupReq: *req,
		Token:    *token,
	}
	// activate segment reservations related to the e2e
	seg, err := db.GetSegmentRsvFromID(ctx, segmentIDFromRaw(t, "ff000000000100000001"))
	require.NoError(t, err)
	require.NotNil(t, seg)
	idx, err := seg.NewIndexFromToken(token, 5, 5)
	require.NoError(t, err)
	err = seg.SetIndexConfirmed(idx)
	require.NoError(t, err)
	err = seg.SetIndexActive(idx)
	require.NoError(t, err)
	err = db.PersistSegmentRsv(ctx, seg)
	require.NoError(t, err)

	// add `count` E2E other segments in DB
	AddE2EReservation(t, db, count)
	backend := db.(*sqlite.Backend)
	e2eRsvsCount, err := backend.DebugCountE2ERsvs(ctx)
	require.NoError(t, err)
	require.Equal(t, e2eRsvsCount, count)

	// now perform the actual E2E admission
	t0 := time.Now()
	_, err = s.AdmitE2EReservation(ctx, successfulReq)
	t1 := time.Since(t0)
	require.NoError(t, err)
	return t1
}

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

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
func newTestSegmentRequest(t testing.TB, ingress, egress uint16,
	minBW, maxBW reservation.BWCls) *segment.SetupReq {

	t.Helper()

	ID := segmentIDFromRaw(t, "ff0000010001beefcafe")
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

func newTestE2ERequest(t testing.TB) *e2e.SetupReq {
	t.Helper()

	id := e2eIDFromRaw(t, "ff00:0:1", "beefcafebeefcafebeef")
	path := test.NewTestPath()
	baseReq, err := e2e.NewRequest(util.SecsToTime(1), id, 1, path)
	require.NoError(t, err)
	segmentRsvs := []reservation.SegmentID{
		*segmentIDFromRaw(t, "ff000000000100000001"),
		*segmentIDFromRaw(t, "ff0000020002beefcafe"),
		*segmentIDFromRaw(t, "ff0000030003beefcafe"),
	}
	ASCountPerSegment := []uint8{4, 4, 5}
	trail := []reservation.BWCls{5, 5}
	setup, err := e2e.NewSetupRequest(baseReq, segmentRsvs, ASCountPerSegment, 5, trail)
	require.NoError(t, err)
	return setup
}
