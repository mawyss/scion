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
	"golang.org/x/net/ipv4"

	"github.com/scionproto/scion/go/lib/serrors"
)

type TrafficClass int
type Scheduler int

// Identify the possible traffic classes. The class "ClsOthers" must be zero.
const (
	ClsOthers TrafficClass = iota
	ClsColibri
	ClsEpic
	ClsScmp
	ClsScion
)

const (
	SchedOthersOnly Scheduler = iota
	SchedRoundRobin
	SchedColibriPrio
)

// queueBufferSize denotes the buffer size of a queue
const queueBufferSize = 64

// maxSchedSize denotes the maximum amount of packets that are scheduled in one round.
const maxSchedSize = 4

type Queues struct {
	mapping  map[TrafficClass]chan ipv4.Message
	nonempty chan bool
}

func NewQueues() *Queues {
	qs := &Queues{}
	qs.nonempty = make(chan bool, 1)
	qs.mapping = make(map[TrafficClass]chan ipv4.Message)

	qs.mapping[ClsOthers] = make(chan ipv4.Message, queueBufferSize)
	qs.mapping[ClsColibri] = make(chan ipv4.Message, queueBufferSize)
	qs.mapping[ClsEpic] = make(chan ipv4.Message, queueBufferSize)
	qs.mapping[ClsScmp] = make(chan ipv4.Message, queueBufferSize)
	qs.mapping[ClsScion] = make(chan ipv4.Message, queueBufferSize)
	return qs
}

func (qs *Queues) Enqueue(tc TrafficClass, m ipv4.Message) error {
	q, ok := qs.mapping[tc]
	if !ok {
		return serrors.New("unknown traffic class")
	}

	select {
	case q <- m:
		qs.setToNonempty()
		return nil
	default:
		return serrors.New("queue full", "traffic class", tc)
	}
}

func (qs *Queues) setToNonempty() {
	select {
	case qs.nonempty <- true:

	default:
	}
}

func (qs *Queues) WaitUntilNonempty() {
	select {
	case <-qs.nonempty:
	}
}

func (qs *Queues) dequeue(tc TrafficClass, batchSize int, ms []ipv4.Message) (int, error) {
	q, ok := qs.mapping[tc]
	if !ok {
		return 0, serrors.New("unknown traffic class")
	}

	var counter int
L:
	for counter = 0; counter < batchSize; counter++ {
		select {
		case m := <-q:
			//log.Debug("Deque some packet", "message", m)
			ms[counter] = m
		default:
			//log.Debug("No packet to dequeue")
			break L
		}
	}
	//log.Debug("Dequeued n packets", "traffic class", tc, "n", counter)
	// reset using ms = ms[:0]
	return counter, nil
}

func (qs *Queues) Schedule(s Scheduler) ([]ipv4.Message, error) {
	switch s {
	case SchedRoundRobin:
		return qs.scheduleRoundRobin()
	case SchedColibriPrio:
		return qs.scheduleColibriPrio()
	case SchedOthersOnly:
		return qs.scheduleOthersOnly()
	default:
		return qs.scheduleRoundRobin()
	}
}

func (qs *Queues) scheduleOthersOnly() ([]ipv4.Message, error) {
	messageBuffer := make([]ipv4.Message, maxSchedSize)
	n, err := qs.dequeue(ClsOthers, maxSchedSize, messageBuffer)
	if err != nil {
		return nil, err
	}
	return messageBuffer[:n], nil
}

func (qs *Queues) scheduleRoundRobin() ([]ipv4.Message, error) {
	return nil, nil
}

func (qs *Queues) scheduleColibriPrio() ([]ipv4.Message, error) {
	// Strict priority: Colibri

	// Remaining classes are scheduled if Colibri queue is empty

	return nil, nil
}
