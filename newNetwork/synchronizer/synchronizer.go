package synchronizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	hotstuff "github.com/HotstuffWASM/newNetwork"
)

// Synchronizer is a dumb implementation of the hotstuff.ViewSynchronizer interface.
// It does not do anything to ensure synchronization, it simply makes the local replica
// propose at the correct time, and send new view messages in case of a timeout.
type Synchronizer struct {
	hotstuff.LeaderRotation

	mut      sync.Mutex
	lastBeat hotstuff.View
	timeout  time.Duration
	timer    *time.Timer
	stop     context.CancelFunc
	hs       hotstuff.Consensus
	stopped  bool
	Proposal chan []byte
	NewView  chan bool
}

// New creates a new Synchronizer.
func New(leaderRotation hotstuff.LeaderRotation, initialTimeout time.Duration) *Synchronizer {
	return &Synchronizer{
		LeaderRotation: leaderRotation,
		timeout:        initialTimeout,
		Proposal:       make(chan []byte, 16),
		NewView:        make(chan bool, 2),
	}
}

// OnPropose should be called when a replica has received a new valid proposal.
func (s *Synchronizer) OnPropose() {
	s.mut.Lock()
	defer s.mut.Unlock()
	if s.timer != nil {
		s.timer.Reset(s.timeout)
	}
}

// OnFinishQC should be called when a replica has created a new qc.
func (s *Synchronizer) OnFinishQC() {
	s.beat()
}

// OnNewView should be called when a replica receives a valid NewView message.
func (s *Synchronizer) OnNewView() {
	// fmt.Println("Should beat")
	s.beat()
}

// Init initializes the synchronizer with given the hotstuff instance.
func (s *Synchronizer) Init(hs hotstuff.Consensus) {
	s.hs = hs
}

// Start starts the synchronizer.
func (s *Synchronizer) Start() {
	if s.GetLeader(s.hs.Leaf().GetView()+1) == s.hs.Config().ID() {
		// fmt.Println("Proposing")
		s.Proposal <- s.hs.Propose()
		// fmt.Println("Proposed on channel")
	}
	s.timer = time.NewTimer(s.timeout)
	// var ctx context.Context
	// ctx, s.stop = context.WithCancel(context.Background())
	go func() {
		s.newViewTimeout()
	}()
}

// Stop stops the synchronizer.
func (s *Synchronizer) Stop() {
	s.stopped = true
	s.stop()
	s.mut.Lock()
	if s.timer != nil && !s.timer.Stop() {
		<-s.timer.C
	}
	s.mut.Unlock()
}

func (s *Synchronizer) beat() {
	if s.stopped {
		fmt.Println("Stopped")
		return
	}
	view := s.hs.Leaf().GetView()
	s.mut.Lock()
	if view <= s.lastBeat {
		s.mut.Unlock()
		// logger.Debug("Can't beat more than once per view ", s.lastBeat)
		fmt.Println("Can't beat more than once per view: ", s.lastBeat)
		// fmt.Println(view)
		return
	}
	if s.GetLeader(view+1) != s.hs.Config().ID() {
		s.mut.Unlock()
		return
	}
	fmt.Println("Propose again...")
	s.lastBeat = view

	go func() {
		s.Proposal <- s.hs.Propose()
	}()
	s.mut.Unlock()
}

func (s *Synchronizer) newViewTimeout() {
	for {
		// time.Sleep(time.Millisecond * 10)
		select {
		// case <-ctx.Done():
		// 	return
		case <-s.timer.C:
			fmt.Println("Timeout")
			// s.mut.Lock()
			s.hs.CreateDummy()
			// s.mut.Unlock()
			if s.GetLeader(s.hs.Leaf().View) == s.hs.Config().ID() {
				go func() {
					// s.NewView <- true
					s.hs.NewView()
				}()
			}
			s.mut.Lock()
			s.timer.Reset(s.timeout)
			s.mut.Unlock()
		}
	}
}
