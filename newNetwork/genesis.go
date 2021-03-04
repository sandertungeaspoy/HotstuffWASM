package hotstuffwasm

var genesisBlock = Block{
	Cert:     nil,
	View:     0,
	Proposer: 0,
}

// GetGenesis returns a pointer to the genesis block, the starting point for the hotstuff blockchain.
func GetGenesis() *Block {
	return &genesisBlock
}
