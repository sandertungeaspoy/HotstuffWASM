package server

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"syscall/js"

	hotstuff "github.com/HotstuffWASM/newNetwork"
	"github.com/HotstuffWASM/newNetwork/config"
	"github.com/HotstuffWASM/newNetwork/synchronizer"
	// "github.com/HotstuffWASM/newNetwork/crypto/ecdsa"
)

// Server is the server-side of the gorums backend.
// It is responsible for calling handler methods on the consensus instance.
type Server struct {
	ID        hotstuff.ID
	Addr      string
	Hs        hotstuff.Consensus
	Pm        *synchronizer.Synchronizer
	Cfg       *Config
	PubKey    *ecdsa.PublicKey
	Cert      *tls.Certificate
	CertPEM   []byte
	PrivKey   *ecdsa.PrivateKey
	Cmds      CmdBuffer
	SendBytes [][]byte
	RecvBytes [][]byte
}

// NewServer creates a new Server.
func NewServer(replicaConfig config.ReplicaConfig) *Server {
	srv := &Server{}
	srv.Addr = replicaConfig.Replicas[replicaConfig.ID].Address
	return srv
}

// // Start creates a listener on the configured address and starts the server.
// func (srv *Server) Start(hs hotstuff.Consensus) error {
// 	lis, err := net.Listen("tcp", srv.addr)
// 	if err != nil {
// 		return fmt.Errorf("failed to listen on %s: %w", srv.addr, err)
// 	}
// 	srv.StartOnListener(hs, lis)
// 	return nil
// }

// // StartOnListener starts the server on the given listener.
// func (srv *Server) StartOnListener(hs hotstuff.Consensus, listener net.Listener) {
// 	srv.hs = hs
// 	go func() {
// 		err := srv.gorumsSrv.Serve(listener)
// 		if err != nil {
// 			logger.Errorf("An error occurred while serving: %v", err)
// 		}
// 	}()
// }

// GetID returns the ID of the sender
func (srv *Server) GetID(msg []byte) (hotstuff.ID, error) {
	msgFromByte := hex.EncodeToString(msg)
	msgStrings := strings.Split(msgFromByte, " ")
	id, err := strconv.ParseUint(msgStrings[1], 10, 32)
	if err != nil {
		return hotstuff.ID(0), err
	}
	return hotstuff.ID(id), err
}

// Propose handles a replica's response to the Propose QC from the leader.
func (srv *Server) Propose(block *hotstuff.Block) string {
	// id, err := srv.GetID()
	// if err != nil {
	// 	panic(err)
	// }
	// // defaults to 0 if error
	// block.Proposer = id
	pc, err := srv.Hs.OnPropose(block)
	if err != nil {
		panic(err)
	}
	// leaderID := srv.hs.Synchronizer().GetLeader(srv.hs.LastVote() + 1)
	// if leaderID == srv.hs.Config().ID() {
	// 	srv.hs.Finish(block)
	// 	srv.hs.OnVote(pc)
	// 	return
	// }

	// leader, ok := srv.hs.Config().Replica(leaderID)
	// if !ok {
	// 	// logger.Warnf("Replica with ID %d was not found!", leaderID)
	// 	return
	// }

	// leader.Vote(pc)
	// srv.hs.Finish(block)
	return pc
}

// Vote handles an incoming vote message.
func (srv *Server) Vote(cert hotstuff.PartialCert) {
	srv.Hs.OnVote(cert)
}

// NewView handles the leader's response to receiving a NewView rpc from a replica.
func (srv *Server) NewView(msg *hotstuff.NewView) {
	// id, err := srv.GetID()
	// if err != nil {
	// 	// logger.Infof("Failed to get client ID: %v", err)
	// 	return
	// }
	// msg.ID = id
	srv.Hs.OnNewView(*msg)
}

// Fetch handles an incoming fetch request.
func (srv *Server) Fetch(hash *hotstuff.Hash) {

	block, ok := srv.Hs.BlockChain().Get(*hash)
	if !ok {
		return
	}

	// logger.Debugf("OnFetch: %.8s", hash)

	id := srv.ID
	// if err != nil {
	// 	// logger.Infof("Fetch: could not get peer id: %v", err)
	// }

	replica, ok := srv.Hs.Config().Replica(id)
	if !ok {
		// logger.Infof("Fetch: could not find replica with id: %d", id)
		return
	}

	replica.Deliver(block)
}

// Deliver handles an incoming deliver message.
func (srv *Server) Deliver(_ context.Context, block *hotstuff.Block) {
	srv.Hs.OnDeliver(block)
}

// Exec executes a command
func (srv *Server) Exec(cmd hotstuff.Command) {
	fmt.Print("Command executed: ")
	fmt.Println(cmd)
	AppendCmd(string(cmd))
	// if cmd == srv.Cmds.Cmds[0] {
	// 	srv.Cmds.Cmds = srv.Cmds.Cmds[1:]
	// }
}

// CmdBuffer is a buffer for the commands
type CmdBuffer struct {
	Cmds          []hotstuff.Command
	mut           sync.Mutex
	serialNumbers map[uint64]int // highest proposed serial number per client ID
}

// Accept accepts incoming comands
func (cmdBuf *CmdBuffer) Accept(cmd hotstuff.Command) bool {
	if len(cmdBuf.serialNumbers) == 0 {
		cmdBuf.serialNumbers = make(map[uint64]int)
	}
	cmdBuf.mut.Lock()
	defer cmdBuf.mut.Unlock()

	cmdString := strings.Split(string(cmd), "sNumber")
	id, _ := strconv.ParseUint(cmdString[0], 10, 32)
	serial, _ := strconv.Atoi(cmdString[1])

	for _, cmds := range cmdBuf.Cmds {
		oldCmdString := strings.Split(string(cmds), "sNumber")
		oldID, _ := strconv.ParseUint(oldCmdString[0], 10, 32)
		if serialNo := cmdBuf.serialNumbers[oldID]; serialNo >= serial {
			// command is too old, can't accept
			return false
		}

	}
	// cmdBuf.Cmds = append(cmdBuf.Cmds, cmd)
	cmdBuf.serialNumbers[id] = serial
	return true
}

// GetCommand returns the front command from the commandbuffer
func (cmdBuf *CmdBuffer) GetCommand() *hotstuff.Command {
	if len(cmdBuf.Cmds) != 0 {
		cmdBuf.mut.Lock()
		cmd := cmdBuf.Cmds[0]
		cmdBuf.Cmds = cmdBuf.Cmds[1:]
		cmdBuf.mut.Unlock()
		return &cmd
	}
	return nil
}

func AppendCmd(cmd string) {

	document := js.Global().Get("document")

	div := document.Call("getElementById", "cmdList")

	// divChild := document.Call("getElementById", "cmdList").Get("childNodes[0]")

	text := document.Call("createElement", "p")

	text.Set("innerText", cmd)

	div.Call("insertBefore", text, div.Get("firstElementChild"))

	// document.Get("body").Call("appendChild", div)
}
