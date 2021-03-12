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

package backend

import (
	"context"
	"database/sql"
	"io"
	"time"

	"github.com/scionproto/scion/go/cs/reservation/e2e"
	"github.com/scionproto/scion/go/cs/reservation/segment"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/colibri/reservation"
	"github.com/scionproto/scion/go/lib/infra/modules/db"
)

// ReserverOnly has the methods available to the AS that starts the reservation.
type ReserverOnly interface {
	// GetSegmentRsvsFromSrcDstIA returns all reservations that start at src AS and end in dst AS.
	GetSegmentRsvsFromSrcDstIA(ctx context.Context, srcIA, dstIA addr.IA) (
		[]*segment.Reservation, error)
	// GetSegmentRsvFromPath searches for a segment reservation with the specified path.
	GetSegmentRsvFromPath(ctx context.Context, path segment.ReservationTransparentPath) (
		*segment.Reservation, error)

	// NewSegmentRsv creates a new segment reservation in the DB, with an unused reservation ID.
	// The created ID is set in the reservation pointer argument. Used by setup req.
	NewSegmentRsv(ctx context.Context, rsv *segment.Reservation) error
}

// TransitOnly represents an AS in-path of a reservation, not the one originating it.
type TransitOnly interface {
	// GetAllSegmentRsvs returns all segment reservations. Used by setup req.
	GetAllSegmentRsvs(ctx context.Context) ([]*segment.Reservation, error)
	// GetSegmentRsvsFromIFPair returns all segment reservations that enter this AS at
	// the specified ingress and exit at that egress. Used by setup req.
	GetSegmentRsvsFromIFPair(ctx context.Context, ingress, egress *uint16) (
		[]*segment.Reservation, error)
}

// ReserverAndTransit contains the functionality for any AS that has a COLIBRI service.
type ReserverAndTransit interface {
	// GetSegmentRsvFromID will return the reservation with that ID.
	// Used by setup/renew req/resp. and any request.
	GetSegmentRsvFromID(ctx context.Context, ID *reservation.SegmentID) (
		*segment.Reservation, error)
	// PersistSegmentRsv ensures the DB contains the reservation as represented in rsv.
	PersistSegmentRsv(ctx context.Context, rsv *segment.Reservation) error
	// DeleteSegmentRsv removes the segment reservation. Used in teardown.
	DeleteSegmentRsv(ctx context.Context, ID *reservation.SegmentID) error

	// DeleteExpiredIndices will remove expired indices from the DB. If a reservation is left
	// without any index after removing the expired ones, it will also be removed. This applies to
	// both segment and e2e reservations.
	// Used on schedule.
	DeleteExpiredIndices(ctx context.Context, now time.Time) (int, error)

	// GetE2ERsvFromID finds the end to end resevation given its ID.
	GetE2ERsvFromID(ctx context.Context, ID *reservation.E2EID) (*e2e.Reservation, error)
	// GetE2ERsvsOnSegRsv returns the e2e reservations running on top of a given segment one.
	GetE2ERsvsOnSegRsv(ctx context.Context, ID *reservation.SegmentID) ([]*e2e.Reservation, error)
	// PersistE2ERsv makes the DB reflect the same contents as the rsv parameter.
	PersistE2ERsv(ctx context.Context, rsv *e2e.Reservation) error
}

// OptimizedStore is implemented by all DBs.
type OptimizedStore interface {
	// GetRsvsPerSource returns the reservations grouped by source that have this ingress OR egress.
	GetRsvsPerSource(ctx context.Context, ingress, egress uint16) (
		map[addr.AS][]*segment.Reservation, error)

	// GetMaxBlockedBWPerSource returns the blocked BW per source AS ID.
	GetMaxBlockedBWPerSource(ctx context.Context, skipRsv reservation.SegmentID) (
		map[addr.AS]uint64, error)

	// TODO(juagargi) the two ^^ previous calls are obsolete

	// GetInterfaceUsageIngress returns the bandwidth already blocked in ingress interface `ifid`.
	GetInterfaceUsageIngress(ctx context.Context, ifid uint16) (uint64, error)

	// GetInterfaceUsageEgress returns the bandwidth already blocked in egress interface `ifid`.
	GetInterfaceUsageEgress(ctx context.Context, ifid uint16) (uint64, error)

	// GetTransitDem returns the stored transit demand between ingress and egress.
	GetTransitDem(ctx context.Context, ingress, egress uint16) (uint64, error)

	// PersistTransit stores the transit state between ingress and egress.
	PersistTransit(ctx context.Context, ingress, egress uint16, transit uint64) error

	// GetSourceState returns the srcDem and srcAlloc for a source,ingress,egress tuple.
	GetSourceState(ctx context.Context, source addr.AS, ingress, egress uint16) (
		uint64, uint64, error)

	// PersistSourceState stores the source state.
	PersistSourceState(ctx context.Context, source addr.AS, ingress, egress uint16,
		srcDem, srcAlloc uint64) error

	GetInDemand(ctx context.Context, source addr.AS, ingress uint16) (uint64, error)
	GetEgDemand(ctx context.Context, source addr.AS, egress uint16) (uint64, error)

	// GetTransitAlloc returns the denominator of the linkRatio formula.
	GetTransitAlloc(ctx context.Context, ingress, egress uint16) (uint64, error)

	// TODO(juagargi) missing state persistance functions for InterfaceUsageIngress/Egress,
	// InDemand/EgDemand, and TransitAlloc

	// TODO(juagargi) missing calling those ^^ functions when modifying the reservations / indices
}

type ColibriStorage interface {
	ReserverOnly
	TransitOnly
	ReserverAndTransit
	OptimizedStore
}

type Transaction interface {
	ColibriStorage
	Commit() error
	Rollback() error
}

// DB is the interface for any reservation backend.
type DB interface {
	BeginTransaction(ctx context.Context, opts *sql.TxOptions) (Transaction, error)
	ColibriStorage
	db.LimitSetter
	io.Closer
}
