// Copyright 2016 ETH Zurich
// Copyright 2019 ETH Zurich, Anapaya Systems
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

// Package te adds traffic engineering capabilities to the border router.
package te

import (
	"github.com/scionproto/scion/go/lib/serrors"
)

type TrafficClass int

const (
	ClsOthers  TrafficClass = iota
	ClsColibri
	ClsEpic
	ClsScmp
	ClsScion
)

const queueBufferSize = 64

type Queues struct {
	mapping map[TrafficClass] chan []byte
}

func NewQueues() *Queues {
	qs := &Queues{}
	qs.mapping = make(map[TrafficClass] chan []byte)

	qs.mapping[ClsOthers] = make(chan []byte, queueBufferSize)
	qs.mapping[ClsColibri] = make(chan []byte, queueBufferSize)
	qs.mapping[ClsEpic] = make(chan []byte, queueBufferSize)
	qs.mapping[ClsScmp] = make(chan []byte, queueBufferSize)
	qs.mapping[ClsScion] = make(chan []byte, queueBufferSize)
	return qs
}

func (qs *Queues) Enqueue(tc TrafficClass, op []byte) error {
	q, ok := qs.mapping[tc]
	if !ok {
		return serrors.New("unknown traffic class")
	}

	select {
	case q <- op:
		return nil
	default:
		return serrors.New("queue full", "traffic class", tc)
	}
}

func (qs *Queues) Dequeue(tc TrafficClass, batchSize int, ops [][]byte) (int, error) {
	q, ok := qs.mapping[tc]
	if !ok {
		return 0, serrors.New("unknown traffic class")
	}

	var counter int
	for counter = 0; counter < batchSize; counter++ {
		select {
		case op := <- q:
			ops[counter] = op
		default:
			break
		}
	}
	// reset using ops = ops[:0]
	return counter, nil
}

func (qs *Queues) Schedule() {
	// Strict priority: Colibri

	// Remaining classes are scheduled if Colibri queue is empty
	
}