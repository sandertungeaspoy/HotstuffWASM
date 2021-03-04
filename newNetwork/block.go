package hotstuffwasm

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strconv"
)

// Block contains a propsed "command", metadata for the protocol, and a link to the "parent" block.
type Block struct {
	// keep a copy of the hash to avoid hashing multiple times
	hash     *Hash
	Parent   Hash
	Proposer ID
	Cmd      Command
	Cert     QuorumCert
	View     View
}

// NewBlock creates a new Block
func NewBlock(parent Hash, cert QuorumCert, cmd Command, view View, proposer ID) *Block {
	return &Block{
		Parent:   parent,
		Cert:     cert,
		Cmd:      cmd,
		View:     view,
		Proposer: proposer,
	}
}

func (b *Block) String() string {
	return fmt.Sprintf(
		"Block{ hash: %.6s parent: %.6s, proposer: %d, view: %d , cert: %v }",
		b.Hash().String(),
		b.Parent.String(),
		b.Proposer,
		b.View,
		b.Cert,
	)
}

func (b *Block) hashSlow() Hash {
	return sha256.Sum256(b.ToBytes())
}

// Hash returns the hash of the Block
func (b *Block) Hash() Hash {
	if b.hash == nil {
		b.hash = new(Hash)
		*b.hash = b.hashSlow()
	}
	return *b.hash
}

// GetProposer returns the id of the replica who proposed the block.
func (b *Block) GetProposer() ID {
	return b.Proposer
}

// GetParent returns the hash of the parent Block
func (b *Block) GetParent() Hash {
	return b.Parent
}

// GetCommand returns the command
func (b *Block) GetCommand() Command {
	return b.Cmd
}

// QuorumCert returns the quorum certificate in the block
func (b *Block) QuorumCert() QuorumCert {
	return b.Cert
}

// GetView returns the view in which the Block was proposed
func (b *Block) GetView() View {
	return b.View
}

// ToBytes returns the raw byte form of the Block, to be used for hashing, etc.
func (b *Block) ToBytes() []byte {
	buf := b.Parent[:]
	var proposerBuf [4]byte
	binary.LittleEndian.PutUint32(proposerBuf[:], uint32(b.Proposer))
	buf = append(buf, proposerBuf[:]...)
	var viewBuf [8]byte
	binary.LittleEndian.PutUint64(viewBuf[:], uint64(b.View))
	buf = append(buf, viewBuf[:]...)
	buf = append(buf, []byte(b.Cmd)...)
	// genesis and dummy nodes have no certificates
	if b.Cert != nil {
		buf = append(buf, b.Cert.ToBytes()...)
	}
	return buf
}

// ToString returns the Block in a string format to be sent to the other replicas
func (b *Block) ToString() string {
	block := b.hash.String() + ":" + b.Parent.String() + ":" + strconv.FormatUint(uint64(b.Proposer), 10) + ":" + string(b.Cmd) + ":" + b.Cert.GetStringSignatures() + ":" + b.Cert.BlockHash().String() + ":" + strconv.FormatUint(uint64(b.View), 10)
	return block
}
