package leaderrotation

import hotstuff "github.com/HotstuffWASM/newNetwork"

type roundRobin struct {
	cfg hotstuff.Config
}

// GetLeader returns the id of the leader in the given view
func (rr roundRobin) GetLeader(view hotstuff.View) hotstuff.ID {
	// TODO: does not support reconfiguration
	// assume IDs start at 1
	// if view == 0 {
	// 	return hotstuff.ID(0)
	// }
	// if view < 10 {
	// 	return hotstuff.ID(1)
	// } else if view < 20 {
	// 	return hotstuff.ID(2)
	// } else if view < 30 {
	// 	return hotstuff.ID(3)
	// } else {
	// 	return hotstuff.ID(4)
	// }
	return hotstuff.ID(view%hotstuff.View(rr.cfg.Len()) + 1)
}

// NewRoundRobin returns a new round-robin leader rotation implementation.
func NewRoundRobin(cfg hotstuff.Config) hotstuff.LeaderRotation {
	return roundRobin{cfg}
}
