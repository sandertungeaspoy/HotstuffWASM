// Package blockchain provides an implementation of the hotstuff.Blockchain interface.
package blockchain

import (
	"container/list"
	"fmt"
	"sync"

	hotstuff "github.com/HotstuffWASM/newNetwork"
)

// blockChain stores a limited amount of blocks in a map.
// blocks are evicted in LRU order.
type blockChain struct {
	mut         sync.Mutex
	maxSize     int
	Blocks      map[hotstuff.Hash]*list.Element
	accessOrder list.List
}

// New creates a new BlockChain with a maximum size.
// Blocks are dropped in least recently used order.
func New(maxSize int) hotstuff.BlockChain {
	return &blockChain{
		maxSize: maxSize,
		Blocks:  make(map[hotstuff.Hash]*list.Element),
	}
}

func (chain *blockChain) dropOldest() {
	elem := chain.accessOrder.Back()
	block := elem.Value.(*hotstuff.Block)
	delete(chain.Blocks, block.Hash())
	chain.accessOrder.Remove(elem)
}

// Store stores a block in the blockchain
func (chain *blockChain) Store(block *hotstuff.Block) {
	chain.mut.Lock()
	defer chain.mut.Unlock()

	if len(chain.Blocks)+1 > chain.maxSize {
		chain.dropOldest()
	}

	elem := chain.accessOrder.PushFront(block)
	chain.Blocks[block.Hash()] = elem
}

// Get retrieves a block given its hash
func (chain *blockChain) Get(hash hotstuff.Hash) (*hotstuff.Block, bool) {
	chain.mut.Lock()
	defer chain.mut.Unlock()

	fmt.Print("Inside get")
	elem, ok := chain.Blocks[hash]
	if !ok {
		return nil, false
	}

	chain.accessOrder.MoveToFront(elem)

	return elem.Value.(*hotstuff.Block), true
}
