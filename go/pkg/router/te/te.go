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
	"net"

	"golang.org/x/net/ipv4"

	"github.com/scionproto/scion/go/lib/log"
	"github.com/scionproto/scion/go/lib/serrors"
	"github.com/scionproto/scion/go/lib/underlay/conn"
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
const queueSize = 16

// Queues describes the queues (one for each traffic class) for a certain router interface.
// The 'mapping' translates traffic classes to their respective queue. The 'nonempty' channel is
// used to signal to the scheduler that packets are ready, which is necessary to omit
// computationally expensive spinning over the queues. The 'scheduler' denotes the scheduling
// algorithm that will be used, and 'writeBuffer' is a pre-allocated buffer to speed up the
// WriteBatch() function.
type Queues struct {
	mapping     map[TrafficClass]chan ipv4.Message
	nonempty    chan bool
	scheduler   Scheduler
	writeBuffer conn.Messages
	queueBuffer map[TrafficClass]chan ipv4.Message
	borrowed    map[TrafficClass]int
}

// NewQueues creates new queues.
func NewQueues(bufSize int) *Queues {
	qs := &Queues{}
	qs.nonempty = make(chan bool, 1)
	qs.mapping = make(map[TrafficClass]chan ipv4.Message)
	qs.queueBuffer = make(map[TrafficClass]chan ipv4.Message)
	qs.borrowed = make(map[TrafficClass]int)

	// Create queues that will contain the packets to be sent
	qs.mapping[ClsOthers] = make(chan ipv4.Message, queueSize)
	qs.mapping[ClsColibri] = make(chan ipv4.Message, queueSize)
	qs.mapping[ClsEpic] = make(chan ipv4.Message, queueSize)
	qs.mapping[ClsBfd] = make(chan ipv4.Message, queueSize)
	qs.mapping[ClsScmp] = make(chan ipv4.Message, queueSize)
	qs.mapping[ClsScion] = make(chan ipv4.Message, queueSize)

	// Create queues that will contain buffers to write packet data into
	for cls := range qs.mapping {
		qs.queueBuffer[cls] = make(chan ipv4.Message, queueSize)
		
		msgs := conn.NewReadMessages(queueSize)
		for _, msg := range msgs {
			msg.Buffers[0] = make([]byte, bufSize)
			qs.queueBuffer[cls] <- msg
		}
	}
	return qs
}

// Enqueue adds the message 'm' to the queue corresponding to traffic class 'tc'.
func (qs *Queues) Enqueue(tc TrafficClass, m []byte, outAddr *net.UDPAddr) error {
	// Get the queue containing the preallocated packets
	qbuf, ok := qs.queueBuffer[tc]
	if !ok {
		return serrors.New("unknown traffic class")
	}
	
	// Retrieve free buffer if available
	var p ipv4.Message
	select {
	case p = <-qbuf:
	default:
		// todo: set congestion flag?
		return serrors.New("no free packet buffer", "traffic class", tc)
	}

	// Copy the message into the buffer
	p.Addr = outAddr
	p.Buffers[0] = p.Buffers[0][:len(m)]
	copy(p.Buffers[0], m)

	// Get the queue where we want to put the packet
	q, ok := qs.mapping[tc]
	if !ok {
		return serrors.New("unknown traffic class")
	}

	// Put the packet into the queue
	select {
	case q <- p:
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
	
	qs.borrowed[tc] = qs.borrowed[tc] + counter
	return counter, nil
}

// Schedule dequeues messages from the queues and prioritizes them according the the specified
// Scheduler. It returns the messages that will be forwarded.
func (qs *Queues) Schedule() ([]ipv4.Message, error) {
	switch qs.scheduler {
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

// TODO
func (qs *Queues) SetScheduler(s Scheduler) {
	qs.scheduler = s
}

// TODO: adapt to selected algorithm
func (qs *Queues) GetMaxBatchSize() int {
	return 16
}

func (qs *Queues) ReturnBuffers(ms []ipv4.Message) error {
	counter := 0
	for tc, borrowed := range qs.borrowed {
		if counter + borrowed > len(ms) {
			return serrors.New("too many packets to return")
		}

		for i := 0; i < borrowed; i++ {
			// Reset buffer
			ms[counter + i].Buffers[0] = ms[counter + i].Buffers[0][:0]

			// Return packet to queue buffer
			select {
			case qs.queueBuffer[tc] <- ms[counter + i]:
			default:
				log.Debug("test")
				return serrors.New("queueBuffer full", "traffic class", tc)
			}
		}

		counter += borrowed
		qs.borrowed[tc] = 0
	}
	return nil
}




// TODO: move to separate file


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

