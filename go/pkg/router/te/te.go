// Copyright 2021 ETH Zurich
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
// Traffic engineering can be enabled/disabled in 'scion/go/posix-router/main.go'.
// The purpose of this package is to demonstrate the feasibility of integrating
// scheduling into the border router. Only basic scheduling algorithms are
// implemented, more elaborate ones might be necessary in the future.
package te

import (
	"golang.org/x/net/ipv4"

	"github.com/scionproto/scion/go/lib/serrors"
)

type TrafficClass int
type Scheduler int

// Identify the possible traffic classes. The class "ClsOthers" must be zero. Packets where the
// traffic class is not set will be assigned to this class by default. Every class must be
// registered inside the 'NewQueues()' function.
const (
	ClsOthers TrafficClass = iota
	ClsColibri
	ClsEpic
	ClsBfd
	ClsScmp
	ClsScion
)

// Identify the possible scheduling algorithms.
const (
	SchedOthersOnly Scheduler = iota
	SchedRoundRobin
	SchedColibriPrio
	SchedStrictPriority
)

// queueSize denotes the buffer size of a queue, i.e., the maximum number of packets that a queue
// can store.
const queueSize = 64

// Queues describes the queues (one for each traffic class) for a certain router interface.
// The 'mapping' translates traffic classes to their respective queue. The 'nonempty' channel is
// used to signal to the scheduler that packets are ready, which is necessary to omit
// computationally expensive spinning over the queues.
type Queues struct {
	mapping  map[TrafficClass]chan ipv4.Message
	nonempty chan bool
}

// NewQueues creates new queues.
func NewQueues() *Queues {
	qs := &Queues{}
	qs.nonempty = make(chan bool, 1)
	qs.mapping = make(map[TrafficClass]chan ipv4.Message)

	qs.mapping[ClsOthers] = make(chan ipv4.Message, queueSize)
	qs.mapping[ClsColibri] = make(chan ipv4.Message, queueSize)
	qs.mapping[ClsEpic] = make(chan ipv4.Message, queueSize)
	qs.mapping[ClsBfd] = make(chan ipv4.Message, queueSize)
	qs.mapping[ClsScmp] = make(chan ipv4.Message, queueSize)
	qs.mapping[ClsScion] = make(chan ipv4.Message, queueSize)
	return qs
}

// Enqueue adds the message 'm' to the queue corresponding to traffic class 'tc'.
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

// setToNonempty signals to the scheduler that new messages are ready to be scheduled.
func (qs *Queues) setToNonempty() {
	select {
	case qs.nonempty <- true:

	default:
	}
}

// WaitUntilNonempty blocks until new messages are ready to be scheduled.
func (qs *Queues) WaitUntilNonempty() {
	select {
	case <-qs.nonempty:
	}
}

// dequeue reads up to 'batchSize' number of messages from the queue corresponding to the traffic
// class 'tc' and stores them in the message buffer 'ms'.
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
			ms[counter] = m
		default:
			break L
		}
	}
	return counter, nil
}

// Schedule dequeues messages from the queues and prioritizes them according the the specified
// Scheduler. It returns the messages that will be forwarded.
func (qs *Queues) Schedule(s Scheduler) ([]ipv4.Message, error) {
	switch s {
	case SchedRoundRobin:
		return qs.scheduleRoundRobin()
	case SchedColibriPrio:
		return qs.scheduleColibriPrio()
	case SchedOthersOnly:
		return qs.scheduleOthersOnly()
	case SchedStrictPriority:
		return qs.scheduleStrictPriority()
	default:
		return qs.scheduleRoundRobin()
	}
}

// scheduleOthersOnly only forwards packets from the 'Others' queue, all other queues are ignored.
func (qs *Queues) scheduleOthersOnly() ([]ipv4.Message, error) {
	maxSchedSize := 8
	messageBuffer := make([]ipv4.Message, maxSchedSize)

	read, err := qs.dequeue(ClsOthers, maxSchedSize, messageBuffer)
	if err != nil {
		return nil, err
	}

	if read > 0 {
		qs.setToNonempty()
	}

	return messageBuffer[:read], nil
}

// scheduleStrictPriority schedules packets based on a strict hierarchy, where a message from a
// queue is only scheduled if all higher priority queues are empty.
// The priorities are: COLIBRI > EPIC > BFD > SCMP > SCION > Others.
func (qs *Queues) scheduleStrictPriority() ([]ipv4.Message, error) {
	maxSchedSize := 8
	messageBuffer := make([]ipv4.Message, maxSchedSize)

	read := 0
	n, err := qs.dequeue(ClsColibri, maxSchedSize-read, messageBuffer[read:])
	if err != nil {
		return nil, err
	}
	read = read + n

	n, err = qs.dequeue(ClsEpic, maxSchedSize-read, messageBuffer[read:])
	if err != nil {
		return nil, err
	}
	read = read + n

	n, err = qs.dequeue(ClsBfd, maxSchedSize-read, messageBuffer[read:])
	if err != nil {
		return nil, err
	}
	read = read + n

	n, err = qs.dequeue(ClsScmp, maxSchedSize-read, messageBuffer[read:])
	if err != nil {
		return nil, err
	}
	read = read + n

	n, err = qs.dequeue(ClsScion, maxSchedSize-read, messageBuffer[read:])
	if err != nil {
		return nil, err
	}
	read = read + n

	n, err = qs.dequeue(ClsOthers, maxSchedSize-read, messageBuffer[read:])
	if err != nil {
		return nil, err
	}
	read = read + n

	if read > 0 {
		qs.setToNonempty()
	}

	return messageBuffer[:read], nil
}

// scheduleRoundRobin schedules the packets based on round-robin.
func (qs *Queues) scheduleRoundRobin() ([]ipv4.Message, error) {
	messageBuffer := make([]ipv4.Message, len(qs.mapping))

	read := 0
	for cls := range qs.mapping {
		n, err := qs.dequeue(cls, 1, messageBuffer[read:])
		if err != nil {
			return nil, err
		}
		read = read + n
	}

	if read > 0 {
		qs.setToNonempty()
	}

	return messageBuffer[:read], nil
}

// scheduleColibriPrio gives priority to Colibri packets, but also includes up to one packet of
// each other traffic class.
func (qs *Queues) scheduleColibriPrio() ([]ipv4.Message, error) {
	maxSchedSize := 16
	messageBuffer := make([]ipv4.Message, maxSchedSize)

	// Prioritize Colibri packets
	nrQueues := len(qs.mapping)
	read := 0
	n, err := qs.dequeue(ClsColibri, maxSchedSize-nrQueues, messageBuffer[read:])
	if err != nil {
		return nil, err
	}
	read = read + n

	// Remaining classes are scheduled using round-robin (including one Colibri packet)
	for cls := range qs.mapping {
		n, err := qs.dequeue(cls, 1, messageBuffer[read:])
		if err != nil {
			return nil, err
		}
		read = read + n
	}

	if read > 0 {
		qs.setToNonempty()
	}

	return messageBuffer[:read], nil
}
