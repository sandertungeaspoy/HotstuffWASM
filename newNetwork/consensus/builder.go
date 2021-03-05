package consensus

import (
	"fmt"

	hotstuff "github.com/HotstuffWASM/newNetwork"
	blockchain "github.com/HotstuffWASM/newNetwork/blockchain"
	ecdsa "github.com/HotstuffWASM/newNetwork/crypto/ecdsa"
)

// Builder is used to set up a HotStuff instance
type Builder struct {
	Config       hotstuff.Config
	BlockChain   hotstuff.BlockChain
	Signer       hotstuff.Signer
	Verifier     hotstuff.Verifier
	CommandQueue hotstuff.CommandQueue
	Executor     hotstuff.Executor
	Acceptor     hotstuff.Acceptor
	Synchronizer hotstuff.ViewSynchronizer
}

// Build returns a new chained HotStuff instance
func (b Builder) Build() hotstuff.Consensus {
	hs := &chainedhotstuff{
		cfg:          b.Config,
		blocks:       b.BlockChain,
		signer:       b.Signer,
		verifier:     b.Verifier,
		commands:     b.CommandQueue,
		executor:     b.Executor,
		acceptor:     b.Acceptor,
		synchronizer: b.Synchronizer,
	}
	signer, verifier := ecdsa.New(b.Config)
	if b.BlockChain == nil {
		fmt.Println("Creating new BlockCHain")
		hs.blocks = blockchain.New(100)
	}
	if b.Signer == nil {
		hs.signer = signer
	}
	if b.Verifier == nil {
		hs.verifier = verifier
	}
	hs.init()
	return hs
}
