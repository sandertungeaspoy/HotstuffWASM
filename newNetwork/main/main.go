package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"

	// "syscall/js"

	// "syscall/js"
	"time"

	hotstuff "github.com/HotstuffWASM/newNetwork"
	"github.com/HotstuffWASM/newNetwork/config"
	consensus "github.com/HotstuffWASM/newNetwork/consensus"
	hsecdsa "github.com/HotstuffWASM/newNetwork/crypto/ecdsa"
	"github.com/HotstuffWASM/newNetwork/leaderrotation"
	server "github.com/HotstuffWASM/newNetwork/server"
	synchronizer "github.com/HotstuffWASM/newNetwork/synchronizer"
	"nhooyr.io/websocket"

	"github.com/pion/webrtc/v3"
)

var sendBytes [][]byte
var recvBytes [][]byte
var serverID hotstuff.ID = 0
var recieved chan []byte
var recvLock sync.Mutex
var sendLock sync.Mutex
var srv server.Server
var incomingCmd chan string
var cmdLock sync.Mutex

var peerMap map[hotstuff.ID]*webrtc.DataChannel

// var starter chan struct{}

var start time.Time

var started bool

// var blocks int

func main() {
	// registerCallbacks()
	started = false
	peerMap = make(map[hotstuff.ID]*webrtc.DataChannel)
	serverID = hotstuff.ID(0)
	// value1 := js.Global().Get("document").Call("getElementById", "self-id").Get("innerText").String()
	// selfID, _ := strconv.ParseUint(strings.Split(value1, " ")[1], 10, 32)
	value1 := os.Args[1]
	selfID, _ := strconv.ParseUint(value1, 10, 32)
	serverID = hotstuff.ID(selfID)
	for {
		if serverID != 0 {
			break
		}
		fmt.Print("Sleeping ZzZ ID: ")
		fmt.Println(serverID)
		time.Sleep(1 * time.Second)
	}
	// blockStr := js.Global().Get("document").Call("getElementById", "blocks").Get("value").String()

	// blocks, _ = strconv.Atoi(blockStr)
	// fmt.Println(blocks)
	// if blocks == 0 {
	// 	blocks = 1000
	// }
	// CreateCommandList()

	sendBytes = make([][]byte, 0)
	recvBytes = make([][]byte, 0)
	recieved = make(chan []byte, 32)
	incomingCmd = make(chan string, 16)
	//Public keys
	pubkeyString1 := "-----BEGIN HOTSTUFF PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEyaKwozY7C9LL4CAGyuY3gQvHrysu\nkW2YuGfGHvgumwRANtalltLIWEQ5OS2ewsR2xastcb/gzUBtyj54Mi1saw==\n-----END HOTSTUFF PUBLIC KEY-----"
	pubBlock1, _ := pem.Decode([]byte(pubkeyString1))
	pubkeyBytes1 := pubBlock1.Bytes
	genKey1, _ := x509.ParsePKIXPublicKey(pubkeyBytes1)
	publicKey1 := genKey1.(*ecdsa.PublicKey)

	pubkeyString2 := "-----BEGIN HOTSTUFF PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEitP7/gomqGK/TSgALUpy+MO9N/n1\nvzyHYXvdwPRFPOyS79UEJIYfNCRyex+TRmtB+jwwo1A+x7hCdk2azaF7FA==\n-----END HOTSTUFF PUBLIC KEY-----"
	pubBlock2, _ := pem.Decode([]byte(pubkeyString2))
	pubkeyBytes2 := pubBlock2.Bytes
	genKey2, _ := x509.ParsePKIXPublicKey(pubkeyBytes2)
	publicKey2 := genKey2.(*ecdsa.PublicKey)

	pubkeyString3 := "-----BEGIN HOTSTUFF PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE6RscDCsZrjOUnRuoUrONyPckRVoo\nt+oGPFjNBynLAWtT07yBCPWYUwzM7Zn+IM3KyAxN12UVZd4itCkGfmOlAg==\n-----END HOTSTUFF PUBLIC KEY-----"
	pubBlock3, _ := pem.Decode([]byte(pubkeyString3))
	pubkeyBytes3 := pubBlock3.Bytes
	genKey3, _ := x509.ParsePKIXPublicKey(pubkeyBytes3)
	publicKey3 := genKey3.(*ecdsa.PublicKey)

	pubkeyString4 := "-----BEGIN HOTSTUFF PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEVgiSObm49gLvwQbqrnNO67nqnSD7\nifeYUWR2o3Z+5fLPD1msFn/PouMBfK0Epjr/MBFiFpVBtM8D+D/RJPVKdg==\n-----END HOTSTUFF PUBLIC KEY-----"
	pubBlock4, _ := pem.Decode([]byte(pubkeyString4))
	pubkeyBytes4 := pubBlock4.Bytes
	genKey4, _ := x509.ParsePKIXPublicKey(pubkeyBytes4)
	publicKey4 := genKey4.(*ecdsa.PublicKey)

	//Private Keys
	privkeyString1 := "-----BEGIN HOTSTUFF PRIVATE KEY-----\nMHcCAQEEICYSXzL1em20GwmW5f5u54V8wddf5uZ/FN+3iQPi0OIroAoGCCqGSM49\nAwEHoUQDQgAEyaKwozY7C9LL4CAGyuY3gQvHrysukW2YuGfGHvgumwRANtalltLI\nWEQ5OS2ewsR2xastcb/gzUBtyj54Mi1saw==\n-----END HOTSTUFF PRIVATE KEY-----"
	privBlock1, _ := pem.Decode([]byte(privkeyString1))
	privkeyBytes1 := privBlock1.Bytes
	privkeyPEM1 := []byte(privkeyString1)
	privateKey1, _ := x509.ParseECPrivateKey(privkeyBytes1)

	privkeyString2 := "-----BEGIN HOTSTUFF PRIVATE KEY-----\nMHcCAQEEIDZbG1HNK/NejlozKQHeLYFTFPPi0QYFHNRP/OwlVLB4oAoGCCqGSM49\nAwEHoUQDQgAEitP7/gomqGK/TSgALUpy+MO9N/n1vzyHYXvdwPRFPOyS79UEJIYf\nNCRyex+TRmtB+jwwo1A+x7hCdk2azaF7FA==\n-----END HOTSTUFF PRIVATE KEY-----"
	privBlock2, _ := pem.Decode([]byte(privkeyString2))
	privkeyBytes2 := privBlock2.Bytes
	privkeyPEM2 := []byte(privkeyString2)
	privateKey2, _ := x509.ParseECPrivateKey(privkeyBytes2)

	privkeyString3 := "-----BEGIN HOTSTUFF PRIVATE KEY-----\nMHcCAQEEIKubRjLNMfX+L4dnDSPcyQIZz/DBdPOyURXUyMibFr/LoAoGCCqGSM49\nAwEHoUQDQgAE6RscDCsZrjOUnRuoUrONyPckRVoot+oGPFjNBynLAWtT07yBCPWY\nUwzM7Zn+IM3KyAxN12UVZd4itCkGfmOlAg==\n-----END HOTSTUFF PRIVATE KEY-----"
	privBlock3, _ := pem.Decode([]byte(privkeyString3))
	privkeyBytes3 := privBlock3.Bytes
	privkeyPEM3 := []byte(privkeyString3)
	privateKey3, _ := x509.ParseECPrivateKey(privkeyBytes3)

	privkeyString4 := "-----BEGIN HOTSTUFF PRIVATE KEY-----\nMHcCAQEEIAE318VoY//HCbPSSzeYv69esRxZcoIpu10YYq1+h+FVoAoGCCqGSM49\nAwEHoUQDQgAEVgiSObm49gLvwQbqrnNO67nqnSD7ifeYUWR2o3Z+5fLPD1msFn/P\nouMBfK0Epjr/MBFiFpVBtM8D+D/RJPVKdg==\n-----END HOTSTUFF PRIVATE KEY-----"
	privBlock4, _ := pem.Decode([]byte(privkeyString4))
	privkeyBytes4 := privBlock4.Bytes
	privkeyPEM4 := []byte(privkeyString4)
	privateKey4, _ := x509.ParseECPrivateKey(privkeyBytes4)

	// Certificates
	certString1 := "-----BEGIN CERTIFICATE-----\nMIIBmjCCAUCgAwIBAgIQGdrdEJSbdGkA0Tc1VEgYQTAKBggqhkjOPQQDAjArMSkw\nJwYDVQQDEyBIb3RTdHVmZiBTZWxmLVNpZ25lZCBDZXJ0aWZpY2F0ZTAeFw0yMTAx\nMTMxMTExMzJaFw0zMTAxMTMxMTExMzJaMCsxKTAnBgNVBAMTIEhvdFN0dWZmIFNl\nbGYtU2lnbmVkIENlcnRpZmljYXRlMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\nyaKwozY7C9LL4CAGyuY3gQvHrysukW2YuGfGHvgumwRANtalltLIWEQ5OS2ewsR2\nxastcb/gzUBtyj54Mi1sa6NGMEQwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoG\nCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwDwYDVR0RBAgwBocEfwAAATAKBggqhkjO\nPQQDAgNIADBFAiB2RyAzqIE80aGGtEgEe9k98k7K1x1Q1z41oNEzMHSTmAIhAPJD\nRvgWsBp/hqtV2/PZUL+zoqOAoexkotun/5SV5ZdY\n-----END CERTIFICATE-----"
	// certBlock1, _ := pem.Decode([]byte(certString1))
	// certBytes1 := certBlock1.Bytes
	certPEM1 := []byte(certString1)
	cert1, _ := tls.X509KeyPair(certPEM1, privkeyPEM1)

	certString2 := "-----BEGIN CERTIFICATE-----\nMIIBmjCCAUCgAwIBAgIQC2t9rKAzWVtTDdJnDInLHDAKBggqhkjOPQQDAjArMSkw\nJwYDVQQDEyBIb3RTdHVmZiBTZWxmLVNpZ25lZCBDZXJ0aWZpY2F0ZTAeFw0yMTAx\nMTMxMTExMzJaFw0zMTAxMTMxMTExMzJaMCsxKTAnBgNVBAMTIEhvdFN0dWZmIFNl\nbGYtU2lnbmVkIENlcnRpZmljYXRlMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\nitP7/gomqGK/TSgALUpy+MO9N/n1vzyHYXvdwPRFPOyS79UEJIYfNCRyex+TRmtB\n+jwwo1A+x7hCdk2azaF7FKNGMEQwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoG\nCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwDwYDVR0RBAgwBocEfwAAATAKBggqhkjO\nPQQDAgNIADBFAiAYQV75FDVJJVgjm6WxVV5PhghT4NlF2PRHb4/ATS1QPAIhAOBc\niM4qVWLaB8KlbgMD0pWAFy+l3w0cHPoICEQTySQ+\n-----END CERTIFICATE-----"
	// certBlock2, _ := pem.Decode([]byte(certString2))
	// certBytes2 := certBlock2.Bytes
	certPEM2 := []byte(certString2)
	cert2, _ := tls.X509KeyPair(certPEM2, privkeyPEM2)

	certString3 := "-----BEGIN CERTIFICATE-----\nMIIBmzCCAUGgAwIBAgIRANE1Qm5JZFIqJmOAwduDsiYwCgYIKoZIzj0EAwIwKzEp\nMCcGA1UEAxMgSG90U3R1ZmYgU2VsZi1TaWduZWQgQ2VydGlmaWNhdGUwHhcNMjEw\nMTEzMTExMTMyWhcNMzEwMTEzMTExMTMyWjArMSkwJwYDVQQDEyBIb3RTdHVmZiBT\nZWxmLVNpZ25lZCBDZXJ0aWZpY2F0ZTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA\nBOkbHAwrGa4zlJ0bqFKzjcj3JEVaKLfqBjxYzQcpywFrU9O8gQj1mFMMzO2Z/iDN\nysgMTddlFWXeIrQpBn5jpQKjRjBEMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAK\nBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA8GA1UdEQQIMAaHBH8AAAEwCgYIKoZI\nzj0EAwIDSAAwRQIhAK0xFL0o7gBFstfJmAvt2k2DYICPzI9JAjmBqMle55T5AiAN\nKcfe5MQ7noWfVkyte60WxWU5Lw2pRDOmIOiG/yXPdg==\n-----END CERTIFICATE-----"
	// certBlock3, _ := pem.Decode([]byte(certString3))
	// certBytes3 := certBlock3.Bytes
	certPEM3 := []byte(certString3)
	cert3, _ := tls.X509KeyPair(certPEM3, privkeyPEM3)

	certString4 := "-----BEGIN CERTIFICATE-----\nMIIBmzCCAUGgAwIBAgIRAP6OtVIpSKXwu9dCxSQUBRcwCgYIKoZIzj0EAwIwKzEp\nMCcGA1UEAxMgSG90U3R1ZmYgU2VsZi1TaWduZWQgQ2VydGlmaWNhdGUwHhcNMjEw\nMTEzMTExMTMyWhcNMzEwMTEzMTExMTMyWjArMSkwJwYDVQQDEyBIb3RTdHVmZiBT\nZWxmLVNpZ25lZCBDZXJ0aWZpY2F0ZTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA\nBFYIkjm5uPYC78EG6q5zTuu56p0g+4n3mFFkdqN2fuXyzw9ZrBZ/z6LjAXytBKY6\n/zARYhaVQbTPA/g/0ST1SnajRjBEMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAK\nBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA8GA1UdEQQIMAaHBH8AAAEwCgYIKoZI\nzj0EAwIDSAAwRQIgccoJdDly+VUGqkDU7wLjTpwYZJtiwIH3nkRhoaWtcOUCIQDE\nJ78TmegrP+YshLGWWpifGE6lMsKnVNWrccBy6ZXjQA==\n-----END CERTIFICATE-----"
	// certBlock4, _ := pem.Decode([]byte(certString4))
	// certBytes4 := certBlock4.Bytes
	certPEM4 := []byte(certString4)
	cert4, _ := tls.X509KeyPair(certPEM4, privkeyPEM4)

	var pubKey []*ecdsa.PublicKey
	pubKey = make([]*ecdsa.PublicKey, 5)
	pubKey[1] = publicKey1
	pubKey[2] = publicKey2
	pubKey[3] = publicKey3
	pubKey[4] = publicKey4

	var cert []*tls.Certificate
	cert = make([]*tls.Certificate, 5)
	cert[1] = &cert1
	cert[2] = &cert2
	cert[3] = &cert3
	cert[4] = &cert4

	var privKey []*ecdsa.PrivateKey
	privKey = make([]*ecdsa.PrivateKey, 5)
	privKey[1] = privateKey1
	privKey[2] = privateKey2
	privKey[3] = privateKey3
	privKey[4] = privateKey4

	var certPEM [][]byte
	certPEM = make([][]byte, 5)
	certPEM[1] = certPEM1
	certPEM[2] = certPEM2
	certPEM[3] = certPEM3
	certPEM[4] = certPEM4

	var addr []string
	addr = make([]string, 5)
	addr[1] = "127.0.0.1:13371"
	addr[2] = "127.0.0.1:13372"
	addr[3] = "127.0.0.1:13373"
	addr[4] = "127.0.0.1:13374"

	// leaderRotation := leaderrotation.NewFixed(hotstuff.ID(1))

	// pm := synchronizer.New(leaderRotation, time.Duration(50)*time.Second)
	var cfg *server.Config

	srv = server.Server{
		ID:   serverID,
		Addr: addr[int(serverID)],
		// Pm:        pm,
		Chess:     false,
		CurrCmd:   0,
		Cfg:       cfg,
		PubKey:    pubKey[serverID],
		Cert:      cert[serverID],
		CertPEM:   certPEM[int(serverID)],
		PrivKey:   privKey[int(serverID)],
		SendBytes: sendBytes,
		RecvBytes: recvBytes,
	}

	// blockStr := js.Global().Get("document").Call("getElementById", "blocks").Get("value").String()
	blockStr := os.Args[2]
	srv.MaxCmd, _ = strconv.Atoi(blockStr)
	if srv.MaxCmd == 0 {
		srv.MaxCmd = 1000
	}

	replicaConfig := config.NewConfig(serverID, srv.PrivKey)
	for i, pub := range pubKey {
		info := &config.ReplicaInfo{
			ID:      hotstuff.ID(i),
			Address: addr[i],
			PubKey:  pub,
		}
		replicaConfig.Replicas[hotstuff.ID(i)] = info
	}

	srv.Cfg = server.NewConfig(*replicaConfig)

	// Round robin
	leaderrotation := leaderrotation.NewRoundRobin(srv.Cfg)
	pm := synchronizer.New(leaderrotation, time.Duration(50)*time.Second, time.Duration(10)*time.Second)
	srv.Pm = pm

	hs := consensus.Builder{
		Config:       srv.Cfg,
		Acceptor:     &srv.Cmds,
		Executor:     &srv,
		Synchronizer: srv.Pm,
		CommandQueue: &srv.Cmds,
	}.Build()

	srv.Hs = hs

	srv.Pm.Init(srv.Hs)

	// srv.Pm.Start()

	EstablishConnections()
	srv.StartTime = time.Now()
	srv.TimeSlice = make([]time.Duration, 0)
	// restart:
	for {
		if started {
			fmt.Println("Starting...")
			break
		}
	}
	for {
		// fmt.Println(srv.Chess)
		// if srv.Pm.GetLeader(hs.LastVote()) != srv.Pm.GetLeader(hs.LastVote()+1) {
		// 	if srv.ID == srv.Pm.GetLeader(hs.LastVote()+1) {
		// 		// fmt.Println("I am Leader")
		// 	} else {
		// 		// fmt.Println("I am Normal Replica")
		// 	}
		// }
		if srv.ID == srv.Pm.GetLeader(hs.LastVote()+1) {
			select {
			case msgByte := <-srv.Pm.Proposal:
				if msgByte == nil {
					continue
				}
				senderID, cmd, obj := FormatBytes(msgByte)
				blockString := msgByte
				if senderID != srv.ID && cmd != "Propose" {
					continue
				}
				block := StringToBlock(obj)
				pcString, err := srv.Hs.OnPropose(block)
				if err != nil {
					fmt.Println(err)
					continue
				}
				srv.Hs.Finish(block)
				pc := StringToPartialCert(pcString)
				srv.Hs.OnVote(pc)
				sendLock.Lock()
				SendCommand(blockString)
				sendLock.Unlock()
				srv.Pm.PropDone = true
			case cmd := <-incomingCmd:
				// fmt.Println("Incoming cmd: ")
				// fmt.Print(cmd)
				cmdLock.Lock()
				command := hotstuff.Command(cmd)
				srv.Cmds.Cmds = append(srv.Cmds.Cmds, command)
				cmdLock.Unlock()
			case <-recieved:
				recvLock.Lock()
				newView := strings.Split(string(recvBytes[0]), ":")
				recvLock.Unlock()
				if newView[0] == "NewView" {
					recvLock.Lock()
					msg := StringToNewView(string(recvBytes[0]))
					recvLock.Unlock()
					srv.Hs.OnNewView(msg)
					recvLock.Lock()
					if len(recvBytes) > 1 {
						recvBytes = recvBytes[1:]
					} else {
						recvBytes = make([][]byte, 0)
					}
					recvLock.Unlock()
					continue
				}
				recvLock.Lock()
				pc := StringToPartialCert(string(recvBytes[0]))
				if len(recvBytes) > 1 {
					recvBytes = recvBytes[1:]
				} else {
					recvBytes = make([][]byte, 0)
				}
				recvLock.Unlock()
				srv.Hs.OnVote(pc)
			case <-srv.Pm.NewView:
				msg := srv.Hs.NewView()
				srv.Hs.OnNewView(msg)
				msgString := NewViewToString(msg)
				// fmt.Println("Sending timeout msg to replicas...")
				sendLock.Lock()
				// sendBytes = append(sendBytes, []byte(msgString))
				SendCommand([]byte(msgString))
				sendLock.Unlock()
			}
			// if srv.Hs.BlockChain().Len()%50 == 0 && srv.Hs.BlockChain().Len() != 0 {
			// 	fmt.Printf("%s took %v\n", "50 blocks", time.Since(start))
			// 	start = time.Now()
			// }
			if srv.Pm.PropDone == true && srv.CurrCmd == srv.MaxCmd && srv.Chess == false {
				srv.Pm.Stop()
				fmt.Printf("%v commands took %v\n", srv.MaxCmd, time.Since(start))
				fmt.Println("Pacemaker stopped...")
				fmt.Println(srv.TimeSlice)
				return
			}

		} else {
			select {
			case <-recieved:
				recvLock.Lock()
				newView := strings.Split(string(recvBytes[0]), ":")
				recvLock.Unlock()
				if newView[0] == "NewView" {
					// fmt.Println("Recieved timeout from leader...")
					recvLock.Lock()
					msg := StringToNewView(string(recvBytes[0]))
					recvLock.Unlock()
					srv.Hs.OnNewView(msg)
					recvLock.Lock()
					if len(recvBytes) > 1 {
						recvBytes = recvBytes[1:]
					} else {
						recvBytes = make([][]byte, 0)
					}
					recvLock.Unlock()
					continue
				}
				recvLock.Lock()
				id, cmd, obj := FormatBytes(recvBytes[0])
				if len(recvBytes) > 1 {
					recvBytes = recvBytes[1:]
				} else {
					recvBytes = make([][]byte, 0)
				}
				recvLock.Unlock()
				// fmt.Print("RecvBytes: ")
				// fmt.Println(recvBytes)
				if id == hotstuff.ID(0) || cmd != "Propose" {
					continue
				}
				block := StringToBlock(obj)
				// fmt.Print("Handle propose for view: ")
				// fmt.Println(block.View)
				pcString, err := srv.Hs.OnPropose(block)
				if err != nil {
					fmt.Println(err)
					continue
				}
				srv.Hs.Finish(block)
				// fmt.Print("Next view: ")
				// fmt.Println(hs.Leaf().GetView() + 1)
				// pc := StringToPartialCert(pcString)
				// srv.Hs.OnVote(pc)
				// fmt.Println("Sending PC to leader...")
				sendLock.Lock()
				// sendBytes = append(sendBytes, []byte(pcString))
				SendCommand([]byte(pcString))
				sendLock.Unlock()
			case <-srv.Pm.NewView:
				timeoutview := srv.Hs.NewView()
				srv.Hs.OnNewView(timeoutview)
				msg := NewViewToString(timeoutview)
				// fmt.Println("Sending timeout msg to leader...")
				sendLock.Lock()
				// sendBytes = append(sendBytes, []byte(msg))
				SendCommand([]byte(msg))
				sendLock.Unlock()
			case cmd := <-incomingCmd:
				cmdLock.Lock()
				command := hotstuff.Command(cmd)
				srv.Cmds.Cmds = append(srv.Cmds.Cmds, command)
				cmdLock.Unlock()
			}
			// if srv.Hs.BlockChain().Len()%50 == 0 && srv.Hs.BlockChain().Len() != 0 {
			// 	fmt.Printf("%s took %v\n", "50 blocks", time.Since(start))
			// 	start = time.Now()
			// }
			if srv.CurrCmd == srv.MaxCmd && srv.Chess == false {
				srv.Pm.Stop()
				fmt.Printf("%v commands took %v\n", srv.MaxCmd, time.Since(start))
				fmt.Println("Pacemaker stopped...")
				fmt.Println(srv.TimeSlice)
				return
			}

		}
		if srv.CurrCmd == srv.MaxCmd-1 {
			fmt.Printf("%v commands took %v\n", srv.MaxCmd, time.Since(start))
			fmt.Println(srv.TimeSlice)
		}

	}

	// fmt.Println("Waiting to restart")
	// for {
	// 	select {
	// 	case <-starter:
	// 		fmt.Println("Restarting")
	// 		srv.Pm.Start()
	// 		goto restart
	// 	}
	// }
	// time.Sleep(time.Second * 5)
	// srv.Pm.Start()
	// goto restart
}

// FormatBytes returns the ID of the sender, the command and the block
func FormatBytes(msg []byte) (id hotstuff.ID, cmd string, obj string) {
	if len(msg) != 0 {
		msgString := string(msg)
		msgStringByte := strings.Split(msgString, ";")
		if len(msgStringByte) == 1 {
			return hotstuff.ID(0), "", ""
		}

		idString, _ := strconv.ParseUint(msgStringByte[1], 10, 32)
		id = hotstuff.ID(idString)

		cmd = msgStringByte[2]

		obj = msgStringByte[3]

		return id, cmd, obj
	}
	return hotstuff.ID(0), "", ""
}

// StringToBlock returns a block from a given string
func StringToBlock(s string) *hotstuff.Block {
	strByte := strings.Split(s, ":")
	// fmt.Print("ParentString: ")
	// // fmt.Println(s)
	// fmt.Print("Converted string")
	// fmt.Println(string())
	// parent, _ := base64.RawStdEncoding.DecodeString(strByte[1])
	// parent := []byte(strByte[1])
	parent, _ := hex.DecodeString(strByte[1])
	// fmt.Println(parent)
	var p [32]byte
	copy(p[:], parent)
	// fmt.Print("p: ")
	// fmt.Println(p)
	parent2 := hotstuff.Hash(p)
	// fmt.Print("Parent hash: ")
	// fmt.Println(parent2)
	proposer, _ := strconv.ParseUint(strByte[2], 10, 32)
	cmd := hotstuff.Command(strByte[3])
	// certHash := []byte(strByte[5])
	certHash, _ := hex.DecodeString(strByte[5])
	var c [32]byte
	copy(c[:], certHash)
	certHash2 := hotstuff.Hash(c)
	var sig map[hotstuff.ID]*hsecdsa.Signature
	sig = make(map[hotstuff.ID]*hsecdsa.Signature)
	sigString := strByte[4]

	sigBytes := strings.Split(sigString, "\n")
	// fmt.Print("SigString:")
	// fmt.Println(sigString)
	if sigString != "" {
		for i := 0; i < len(sigBytes)-1; i++ {
			m := strings.Split(sigBytes[i], "=")
			id, _ := strconv.ParseUint(m[0], 10, 32)
			signString := strings.Split(m[1], "-")
			rInt := new(big.Int)
			rInt.SetString(signString[0], 0)
			sInt := new(big.Int)
			sInt.SetString(signString[1], 0)
			signer, _ := strconv.ParseUint(signString[2], 10, 32)
			sign := *hsecdsa.NewSignature(rInt, sInt, hotstuff.ID(signer))
			sig[hotstuff.ID(id)] = &sign
		}
	}

	var cert hotstuff.QuorumCert = hsecdsa.NewQuorumCert(sig, certHash2)
	view, _ := strconv.ParseUint(strByte[6], 10, 64)

	b := &hotstuff.Block{
		Parent:   parent2,
		Proposer: hotstuff.ID(proposer),
		Cmd:      cmd,
		Cert:     cert,
		View:     hotstuff.View(view),
	}
	b.Hash()
	return b
}

// StringToPartialCert returns a PartialCert from a string
func StringToPartialCert(s string) hotstuff.PartialCert {
	strByte := strings.Split(s, ":")
	signString := strings.Split(strByte[0], "-")

	rInt := new(big.Int)
	rInt.SetString(signString[0], 0)
	sInt := new(big.Int)
	sInt.SetString(signString[1], 0)
	signer, _ := strconv.ParseUint(signString[2], 10, 32)
	sign := *hsecdsa.NewSignature(rInt, sInt, hotstuff.ID(signer))

	hash, _ := hex.DecodeString(strByte[1])
	var h [32]byte
	copy(h[:], hash)
	hash2 := hotstuff.Hash(h)
	var pc hotstuff.PartialCert = hsecdsa.NewPartialCert(&sign, hash2)

	return pc
}

// StringToNewView converts the string to a NewView message
func StringToNewView(s string) hotstuff.NewView {
	stringByte := strings.Split(s, ":")
	viewID, _ := strconv.ParseUint(stringByte[1], 10, 32)

	view, _ := strconv.ParseUint(stringByte[2], 10, 32)

	certHash, _ := hex.DecodeString(stringByte[4])
	var c [32]byte
	copy(c[:], certHash)
	certHash2 := hotstuff.Hash(c)
	var sig map[hotstuff.ID]*hsecdsa.Signature
	sig = make(map[hotstuff.ID]*hsecdsa.Signature)
	sigString := stringByte[3]

	sigBytes := strings.Split(sigString, "\n")
	// fmt.Print("SigString:")
	// fmt.Println(sigString)
	if sigString != "" {
		for i := 0; i < len(sigBytes)-1; i++ {
			m := strings.Split(sigBytes[i], "=")
			id, _ := strconv.ParseUint(m[0], 10, 32)
			signString := strings.Split(m[1], "-")
			rInt := new(big.Int)
			rInt.SetString(signString[0], 0)
			sInt := new(big.Int)
			sInt.SetString(signString[1], 0)
			signer, _ := strconv.ParseUint(signString[2], 10, 32)
			sign := *hsecdsa.NewSignature(rInt, sInt, hotstuff.ID(signer))
			sig[hotstuff.ID(id)] = &sign
		}
	}

	var cert hotstuff.QuorumCert = hsecdsa.NewQuorumCert(sig, certHash2)

	newView := hotstuff.NewView{ID: hotstuff.ID(viewID), View: hotstuff.View(view), QC: cert}
	return newView
}

// NewViewToString returns the NewView message as a string
func NewViewToString(view hotstuff.NewView) string {
	msg := "NewView:" + strconv.FormatUint(uint64(view.ID), 10) + ":" + strconv.FormatUint(uint64(view.View), 10) + ":" + view.QC.GetStringSignatures() + ":" + view.QC.BlockHash().String()
	return msg
}

func ConnectToPeer() (*webrtc.DataChannel, string) {
	// create:
	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// peerConnection.OnICEGatheringStateChange(func() {
	// 	fmt.Println(peerConnection.ICEGatheringState())
	// })

	// Register data channel creation handling

	var dc *webrtc.DataChannel

	waiter := make(chan struct{})

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open.\n", d.Label(), d.ID())
			// Register text message handling
			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				// fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
				// err := d.SendText("Message Received")
				// if strings.TrimSpace(string(msg.Data)) == "StartConnectionLeader" {
				// go ConnectionLeader()
				// fmt.Println("Starting connection leader")
				if msg.IsString {
					if strings.TrimSpace(string(msg.Data)) == "StartConnectionLeader" {
						go ConnectionLeader()
					} else if strings.TrimSpace(string(msg.Data)) == "StartWasmStuff" {
						if srv.ID == srv.Pm.GetLeader(srv.Hs.LastVote()+1) {
							srv.Pm.Start()
							// start = time.Now()
							started = true
							// srv.Pm.Proposal <- srv.Hs.Propose()
							srv.Pm.PropDone = false
						} else {
							srv.Pm.Start()
							started = true
							// start = time.Now()
						}
					} else if strings.TrimSpace(string(msg.Data)) == "startChessWhite" {
						// fmt.Println("Starting chess")
						// CreateChessBoard("white")
						srv.Chess = true
					} else if strings.TrimSpace(string(msg.Data)) == "startChessSpectate" {
						// fmt.Println("Starting chess")
						// CreateChessBoard("spectate")
						srv.Chess = true
					}
				} else {
					recvLock.Lock()
					recvBytes = append(recvBytes, msg.Data)
					recvLock.Unlock()
					recieved <- msg.Data
				}
				// if err != nil {
				// 	fmt.Println(err)
				// }
			})
			dc = d
			// fmt.Println(dc)
			close(waiter)

			d.OnClose(func() {
				fmt.Printf("Data channel '%s'-'%d' has been closed\n", d.Label(), d.ID())

				peerKey, ok := mapkeyDataChannel(peerMap, d)

				if ok {
					delete(peerMap, peerKey)
					id := strconv.FormatUint(uint64(peerKey), 10)
					removeAnswer(id)
				}

			})

		})

	})

	offer := "empty"
	senderID := ""
	attempts := 0
	for {
		offer, senderID = ReceiveOffer()
		if offer == "error" {
			return nil, "error"
		}
		if offer != "empty" {
			break
		}

		if attempts > 50 {
			return nil, "error"
		}
		attempts++
		time.Sleep(time.Millisecond * 500)
	}

	offersdp := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offer}

	err = peerConnection.SetRemoteDescription(offersdp)
	if err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// for {
	// 	fmt.Println(peerConnection.ICEGatheringState())
	// 	if peerConnection.ICEGatheringState() == webrtc.ICEGatheringStateComplete && strings.Contains(peerConnection.LocalDescription().SDP, "c=IN IP4 0.0.0.0") {
	// 		fmt.Println(peerConnection.LocalDescription().SDP)
	// 		DeliverAnswer(peerConnection.LocalDescription().SDP, senderID)
	// 		break
	// 	} else if peerConnection.ICEGatheringState() == webrtc.ICEGatheringStateComplete && !strings.Contains(peerConnection.LocalDescription().SDP, "c=IN IP4 0.0.0.0") {
	// 		goto create
	// 	}
	// 	// time.Sleep(time.Second)

	// }

	<-gatherComplete

	DeliverAnswer(peerConnection.LocalDescription().SDP, senderID)

	// if peerConnection.ICEGatheringState() == webrtc.ICEGatheringStateComplete && strings.Contains(peerConnection.LocalDescription().SDP, "c=IN IP4 0.0.0.0") {
	// 	// fmt.Println(peerConnection.LocalDescription().SDP)
	// 	DeliverAnswer(peerConnection.LocalDescription().SDP, senderID)

	// } else if peerConnection.ICEGatheringState() == webrtc.ICEGatheringStateComplete && !strings.Contains(peerConnection.LocalDescription().SDP, "c=IN IP4 0.0.0.0") {
	// 	goto create
	// }

	<-waiter

	removeAnswer(senderID)
	// fmt.Println("Returning")
	return dc, senderID
}

func ConnectToLeader() (*webrtc.DataChannel, string) {

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// peerConnection.OnICEGatheringStateChange(func() {
	// 	fmt.Println(peerConnection.ICEGatheringState())
	// })

	waiter := make(chan struct{})

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open.\n", dataChannel.Label(), dataChannel.ID())
		close(waiter)
		// for range time.NewTicker(5 * time.Second).C {
		// 	message := "signal.RandSeq(15)"
		// 	fmt.Printf("Sending '%s'\n", message)

		// 	// Send the message as text
		// 	sendErr := dataChannel.SendText(message)
		// 	if sendErr != nil {
		// 		panic(sendErr)
		// 	}
		// }
	})

	dataChannel.OnClose(func() {
		fmt.Printf("Data channel '%s'-'%d' has been closed\n", dataChannel.Label(), dataChannel.ID())

		peerKey, ok := mapkeyDataChannel(peerMap, dataChannel)

		if ok {
			delete(peerMap, peerKey)
			id := strconv.FormatUint(uint64(peerKey), 10)
			removeAnswer(id)
		}
		purgeWebRTCDatabase()

	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		// fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
		// if strings.TrimSpace(string(msg.Data)) == "StartConnectionLeader" {
		// go ConnectionLeader()
		// fmt.Println("Starting connection leader")
		if msg.IsString {
			// fmt.Println(string(msg.Data))
			if strings.TrimSpace(string(msg.Data)) == "StartConnectionLeader" {
				go ConnectionLeader()
			} else if strings.TrimSpace(string(msg.Data)) == "StartWasmStuff" {
				if srv.ID == srv.Pm.GetLeader(srv.Hs.LastVote()+1) {
					srv.Pm.Start()
					started = true
					// start = time.Now()
					// srv.Pm.Proposal <- srv.Hs.Propose()
					srv.Pm.PropDone = false
				} else {
					srv.Pm.Start()
					started = true
					// start = time.Now()
				}
			} else if strings.TrimSpace(string(msg.Data)) == "startChessWhite" {
				// fmt.Println("Starting chess")
				// CreateChessBoard("white")
				srv.Chess = true
			} else if strings.TrimSpace(string(msg.Data)) == "startChessSpectate" {
				// fmt.Println("Starting chess")
				// CreateChessBoard("spectate")
				srv.Chess = true
			}
		} else {
			recvLock.Lock()
			recvBytes = append(recvBytes, msg.Data)
			recvLock.Unlock()
			recieved <- msg.Data
		}
	})

	// Create an offer to send to the browser
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// for {
	// 	fmt.Println(peerConnection.ICEGatheringState())
	// 	if peerConnection.ICEGatheringState() == webrtc.ICEGatheringStateComplete {
	// 		fmt.Println(peerConnection.LocalDescription().SDP)
	// 		DeliverOffer(peerConnection.LocalDescription().SDP)
	// 		break
	// 	}
	// 	// time.Sleep(time.Second)

	// }

	<-gatherComplete

	DeliverOffer(peerConnection.LocalDescription().SDP)

	answer := "empty"
	senderID := ""
	attempts := 0
	for {
		answer, senderID = ReceiveAnswer()
		if answer == "error" {
			return nil, "error"
		}
		if answer != "empty" {
			break
		}
		if attempts > 50 {
			return nil, "error"
		}
		attempts++
		time.Sleep(time.Millisecond * 500)
	}

	// fmt.Println(answer)

	answersdp := webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: answer}

	// Apply the answer as the remote description
	err = peerConnection.SetRemoteDescription(answersdp)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Remote desc set")

	<-waiter
	RemoveOffer()
	// fmt.Println("Returning")
	return dataChannel, senderID
}

func DeliverOffer(offer string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://85.165.212.251:13372", nil)
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusInternalError, "WebSocket has been closed")

	conn := websocket.NetConn(ctx, c, 1)

	id := strconv.FormatUint(uint64(serverID), 10)

	fmt.Fprintf(conn, offer+"&"+id+"%")

	c.Close(websocket.StatusNormalClosure, "")
}

func DeliverAnswer(answer string, senderID string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://85.165.212.251:13372", nil)
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusInternalError, "WebSocket has been closed")

	conn := websocket.NetConn(ctx, c, 1)

	// id := strconv.FormatUint(uint64(serverID), 10)

	fmt.Fprintf(conn, answer+"&"+senderID+"%")

	c.Close(websocket.StatusNormalClosure, "")
}

func ReceiveOffer() (string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://85.165.212.251:13372", nil)
	if err != nil {
		return "error", "Websocket"
	}
	defer c.Close(websocket.StatusInternalError, "WebSocket has been closed")

	conn := websocket.NetConn(ctx, c, 1)

	id := strconv.FormatUint(uint64(serverID), 10)

	fmt.Fprintf(conn, "setup:recvOffer\n&"+id+"%")

	offer, err := bufio.NewReader(conn).ReadString('%')
	if err != nil {
		return "empty", "message"
	}
	// fmt.Println(offer)

	msgs := strings.Split(offer, "&")

	c.Close(websocket.StatusNormalClosure, "")
	return msgs[0], strings.Split(msgs[1], "%")[0]
}

func ReceiveAnswer() (string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://85.165.212.251:13372", nil)
	if err != nil {
		return "error", "Websocket"
	}
	defer c.Close(websocket.StatusInternalError, "WebSocket has been closed")

	conn := websocket.NetConn(ctx, c, 1)

	id := strconv.FormatUint(uint64(serverID), 10)

	fmt.Fprintf(conn, "setup:recvAnswer\n&"+id+"%")

	answer, err := bufio.NewReader(conn).ReadString('%')
	if err != nil {
		return "empty", "message"
	}

	// fmt.Println(answer)

	msgs := strings.Split(answer, "&")

	c.Close(websocket.StatusNormalClosure, "")
	return msgs[0], strings.Split(msgs[1], "%")[0]
}

func RemoveOffer() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://85.165.212.251:13372", nil)
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusInternalError, "WebSocket has been closed")

	conn := websocket.NetConn(ctx, c, 1)

	id := strconv.FormatUint(uint64(serverID), 10)

	fmt.Fprintf(conn, "setup:removeOffer\n&"+id+"%")

	c.Close(websocket.StatusNormalClosure, "")
}

func removeAnswer(senderID string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://85.165.212.251:13372", nil)
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusInternalError, "WebSocket has been closed")

	conn := websocket.NetConn(ctx, c, 1)

	fmt.Fprintf(conn, "setup:removeAnswer\n&"+senderID+"%")

	c.Close(websocket.StatusNormalClosure, "")
}

func purgeWebRTCDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://85.165.212.251:13372", nil)
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusInternalError, "WebSocket has been closed")

	conn := websocket.NetConn(ctx, c, 1)

	fmt.Fprintf(conn, "")

	c.Close(websocket.StatusNormalClosure, "setup:purgeDatabase\n&0%")
}

func EstablishConnections() {

	// started := false
	serverConID := 1
	// leader := false

	if srv.ID == hotstuff.ID(1) {
		for {
			if len(peerMap) == 3 {
				time.Sleep(time.Second * 20)
				break
			}

			dc, peerID := ConnectToPeer()
			if peerID == "error" || peerID == "empty" {
				time.Sleep(time.Second * 5)
				continue
			}

			peerIDUint, _ := strconv.ParseUint(peerID, 10, 32)
			peerIDHot := hotstuff.ID(peerIDUint)

			peerMap[peerIDHot] = dc
			if len(peerMap) == 3 && !started {
				peerMap[hotstuff.ID(2)].SendText("StartConnectionLeader")
				purgeWebRTCDatabase()
				// srv.Pm.Start()
				// srv.Pm.Proposal <- srv.Hs.Propose()
				// srv.Pm.PropDone = false
				started = true
				// CreateChessGame()
			}
		}
	} else if srv.ID == hotstuff.ID(2) {
		for {
			if len(peerMap) == 0 {
				dc, leaderID := ConnectToLeader()
				if leaderID == "error" || leaderID == "empty" {
					// time.Sleep(time.Second * 5)
					continue
				}

				// leaderIDUint, _ := strconv.ParseUint(leaderID, 10, 32)
				// leaderIDHot := hotstuff.ID(leaderIDUint)

				peerMap[hotstuff.ID(serverConID)] = dc
				// srv.Pm.Start()
				break
			}
			time.Sleep(time.Second * 5)
		}
	} else if srv.ID == hotstuff.ID(3) {

		for {
			if len(peerMap) < 2 {
				dc, leaderID := ConnectToLeader()
				if leaderID == "error" || leaderID == "empty" {
					// time.Sleep(time.Second * 5)
					continue
				}

				// leaderIDUint, _ := strconv.ParseUint(leaderID, 10, 32)
				// leaderIDHot := hotstuff.ID(leaderIDUint)

				peerMap[hotstuff.ID(serverConID)] = dc
				serverConID++
				// srv.Pm.Start()
			} else {
				break
			}
			time.Sleep(time.Second * 5)
		}
	} else if srv.ID == hotstuff.ID(4) {
		for {
			if len(peerMap) < 3 {
				dc, leaderID := ConnectToLeader()
				if leaderID == "error" || leaderID == "empty" {
					// time.Sleep(time.Second * 5)
					continue
				}

				// leaderIDUint, _ := strconv.ParseUint(leaderID, 10, 32)
				// leaderIDHot := hotstuff.ID(leaderIDUint)

				peerMap[hotstuff.ID(serverConID)] = dc
				serverConID++

				if len(peerMap) == 3 && !started {
					SendStringTo("StartWasmStuff", hotstuff.ID(0))
					srv.Pm.Start()
					started = true
					// CreateChessGame()
					// start = time.Now()
					purgeWebRTCDatabase()
					// started = true
					break
				}
			}
			time.Sleep(time.Second * 5)
		}
	}

	// for {
	// 	if srv.ID == srv.Pm.GetLeader(srv.Hs.Leaf().GetView()+1) {
	// 		if leader == false {
	// 			purgeWebRTCDatabase()
	// 			leader = true
	// 		}

	// 		if len(peerMap) == 3 {
	// 			time.Sleep(time.Second * 20)
	// 			continue
	// 		}

	// 		dc, peerID := ConnectToPeer()
	// 		if peerID == "error" || peerID == "empty" {
	// 			time.Sleep(time.Second * 5)
	// 			continue
	// 		}

	// 		peerIDUint, _ := strconv.ParseUint(peerID, 10, 32)
	// 		peerIDHot := hotstuff.ID(peerIDUint)

	// 		peerMap[peerIDHot] = dc
	// 		if len(peerMap) == 3 && !started {
	// 			srv.Pm.Proposal <- srv.Hs.Propose()
	// 			srv.Pm.PropDone = false
	// 			started = true
	// 		}

	// 	} else {
	// 		if len(peerMap) == 0 {
	// 			dc, leaderID := ConnectToLeader()
	// 			if leaderID == "error" || leaderID == "empty" {
	// 				// time.Sleep(time.Second * 5)
	// 				continue
	// 			}

	// 			leaderIDUint, _ := strconv.ParseUint(leaderID, 10, 32)
	// 			leaderIDHot := hotstuff.ID(leaderIDUint)

	// 			peerMap[leaderIDHot] = dc
	// 			srv.Pm.Start()
	// 		}
	// 		time.Sleep(time.Second * 30)
	// 	}
	// }

}

func ConnectionLeader() {
	for {
		if len(peerMap) == 3 {
			// time.Sleep(time.Second * 20)
			break
		}

		dc, peerID := ConnectToPeer()
		if peerID == "error" || peerID == "empty" {
			time.Sleep(time.Second * 5)
			continue
		}

		peerIDUint, _ := strconv.ParseUint(peerID, 10, 32)
		peerIDHot := hotstuff.ID(peerIDUint)

		peerMap[peerIDHot] = dc
		if len(peerMap) == 3 {
			if srv.ID == hotstuff.ID(2) {
				peerMap[hotstuff.ID(3)].SendText("StartConnectionLeader")
			}
			purgeWebRTCDatabase()
			// srv.Pm.Start()
			// start = time.Now()
			break
		}
	}
}

// mapkeyDataChannel finds the key for a specific datachannel in the peermap
func mapkeyDataChannel(m map[hotstuff.ID]*webrtc.DataChannel, value *webrtc.DataChannel) (key hotstuff.ID, ok bool) {
	for k, v := range m {
		if v.ID() == value.ID() {
			key = k
			ok = true
			return
		}
	}
	return
}

func SendCommand(cmd []byte) error {
	if srv.ID == srv.Pm.GetLeader(srv.Hs.LastVote()) {
		for _, peer := range peerMap {
			err := peer.Send(cmd)
			if err != nil {
				return err
			}
		}
	} else {
		// fmt.Print("Sending to ")
		// fmt.Println(srv.Pm.GetLeader(srv.Hs.LastVote() + 1))
		// fmt.Println(peerMap)
		if srv.ID == srv.Pm.GetLeader(srv.Hs.LastVote()+1) {
			recvLock.Lock()
			recvBytes = append(recvBytes, cmd)
			recvLock.Unlock()
			recieved <- cmd

		} else {
			conn, ok := peerMap[srv.Pm.GetLeader(srv.Hs.LastVote()+1)]
			if ok {
				err := conn.Send(cmd)
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func SendStringTo(cmd string, srvID hotstuff.ID) error {
	if srvID == hotstuff.ID(0) {
		for _, peer := range peerMap {
			err := peer.SendText(cmd)
			if err != nil {
				return err
			}
		}
	} else {
		err := peerMap[srvID].SendText(cmd)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

// // GetSelfID gets the ID of the server
// func GetSelfID(this js.Value, i []js.Value) interface{} {
// 	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()

// 	selfID, _ := strconv.ParseUint(value1, 10, 32)
// 	serverID = hotstuff.ID(selfID)
// 	fmt.Println(serverID)
// 	return nil
// }

// // GetBlockNumber gets the amount of blocks to run for
// func GetBlockNumber(this js.Value, i []js.Value) interface{} {
// 	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()

// 	srv.MaxCmd, _ = strconv.Atoi(value1)
// 	return nil
// }

// // PassUint8ArrayToGo passes array
// func PassUint8ArrayToGo(this js.Value, args []js.Value) interface{} {
// 	recv := make([]byte, args[0].Get("length").Int())

// 	_ = js.CopyBytesToGo(recv, args[0])

// 	recvLock.Lock()
// 	recvBytes = append(recvBytes, recv)
// 	recvLock.Unlock()
// 	recieved <- recv

// 	return nil
// }

// // SetUint8ArrayInGo sets array
// func SetUint8ArrayInGo(this js.Value, args []js.Value) interface{} {
// 	sendLock.Lock()
// 	if len(sendBytes) == 0 {
// 		sendLock.Unlock()
// 		return nil
// 	}

// 	if args[0].Get("length").Int() == 0 {
// 		sendLock.Unlock()
// 		return nil
// 	}

// 	var msg []byte
// 	if len(sendBytes) > 1 {
// 		msg, sendBytes = sendBytes[0], sendBytes[1:]
// 	} else {
// 		msg, sendBytes = sendBytes[0], make([][]byte, 0)
// 	}
// 	sendLock.Unlock()
// 	if msg == nil {
// 		return nil
// 	}
// 	// fmt.Println("Sending bytes to JS")
// 	_ = js.CopyBytesToJS(args[0], msg)

// 	return nil
// }

// // GetArraySize gets the array size
// func GetArraySize(this js.Value, args []js.Value) interface{} {

// 	if len(sendBytes) == 0 {
// 		_ = js.CopyBytesToJS(args[1], []byte{0})
// 		return nil
// 	}
// 	size := make([]byte, 10)

// 	msgSize := []byte(strconv.Itoa(len(sendBytes[0])))

// 	copy(size, msgSize)

// 	_ = js.CopyBytesToJS(args[0], size)
// 	_ = js.CopyBytesToJS(args[1], []byte{1})

// 	return nil
// }

// // GetCommand gets the ID of the server
// func GetCommand(this js.Value, i []js.Value) interface{} {
// 	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()
// 	fmt.Println(string(value1))

// 	cmd := string(value1)
// 	if cmd != "" {
// 		cmd = strconv.FormatUint(uint64(serverID), 10) + "cmdID" + cmd
// 		incomingCmd <- cmd
// 	}
// 	// cmd = strconv.FormatUint(uint64(serverID), 10) + "cmdID" + cmd
// 	// if serverID == srv.Pm.GetLeader(srv.Hs.LastVote()+1) {
// 	// 	cmdLock.Lock()
// 	// 	command := hotstuff.Command(cmd)
// 	// 	srv.Cmds.Cmds = append(srv.Cmds.Cmds, command)
// 	// 	cmdLock.Unlock()
// 	// 	// srv.Pm.Proposal <- srv.Hs.Propose()
// 	// } else {
// 	// 	incomingCmd <- cmd
// 	// }
// 	return nil
// }

// func StartAgain(this js.Value, args []js.Value) interface{} {
// 	// fmt.Println("Before")
// 	starter <- struct{}{}
// 	// fmt.Println("After")
// 	return nil
// }

// func AppendCmd(document js.Value, cmd string) {

// 	div := js.Global().Get("document").Call("getElementById", "cmdList")

// 	text := document.Call("createElement", "p")

// 	text.Set("innerText", cmd)

// 	div.Call("appendChild", text)

// 	document.Get("body").Call("appendChild", div)
// }

// func CreateCommandList() error {

// 	document := js.Global().Get("document")
// 	div := document.Call("createElement", "div")

// 	div.Call("setAttribute", "style", "overflow:scroll; height:500px; width:500px; float:right; margin-right:150px")
// 	div.Call("setAttribute", "id", "cmdList")

// 	text := document.Call("createElement", "p")

// 	text.Set("innerText", "List of Executed Commands")

// 	div.Call("appendChild", text)

// 	document.Get("body").Call("appendChild", div)

// 	return nil
// }

// func CreateChessGame() error {
// 	document := js.Global().Get("document")
// 	div := document.Call("createElement", "div")
// 	div.Call("setAttribute", "id", "ChessDiv")
// 	textbox := document.Call("createElement", "input")
// 	textbox.Call("setAttribute", "type", "text")
// 	textbox.Call("setAttribute", "id", "ChessVS")
// 	lbl := document.Call("createElement", "label")
// 	lbl.Call("setAttribute", "for", "ChessVS")
// 	lbl.Set("innerText", "Player to invite: ")

// 	buttonGroup := document.Call("createElement", "div")
// 	buttonGroup.Call("setAttribute", "id", "ChessVS")
// 	// buttonGroup.Call("setAttribute", "class", "center")

// 	for id := range peerMap {
// 		fmt.Println("Creating buttons..")
// 		button := document.Call("createElement", "button")
// 		button.Call("setAttribute", "type", "button")
// 		idInt := strconv.FormatUint(uint64(id), 10)
// 		innerText := "Server " + idInt
// 		button.Set("innerText", innerText)
// 		if id == 1 {
// 			button.Call("setAttribute", "class", "button1")
// 		} else if id == 2 {
// 			button.Call("setAttribute", "style", "background-color: #ed7d31")
// 		} else if id == 3 {
// 			button.Call("setAttribute", "style", "background-color: #5b9bd5")
// 		} else if id == 4 {
// 			button.Call("setAttribute", "style", "background-color: #70ad47")
// 		}
// 		button.Call("setAttribute", "onClick", "CreateChess()")
// 		buttonGroup.Call("appendChild", button)
// 	}
// 	// document.getElementById('self-id').style = 'background-color: #ffc000;
// 	// <div id="menu" class="btn-group-vertical">
// 	// 	<button type="button" class="btn btn-md btn-default" onclick="restartWebRTCHandler();" id="restartButton">Restart WebRTC</button>
// 	// 	<a id="checksum" role="button" class="btn btn-md btn-default" href="hash.txt" download>MD5 Checksum</a>
// 	// </div>

// 	// CreatChessBtn := document.Call("createElement", "button")
// 	// CreatChessBtn.Set("innerText", "Invite to Chess")
// 	// CreatChessBtn.Call("setAttribute", "id", "ChessGen")
// 	// CreatChessBtn.Call("setAttribute", "onClick", "CreateChess()")

// 	div.Call("appendChild", lbl)
// 	div.Call("appendChild", buttonGroup)
// 	// div.Call("appendChild", CreatChessBtn)

// 	document.Call("getElementById", "chessGame").Call("appendChild", div)

// 	return nil
// }

// func CreateChessBoard(color string) {
// 	document := js.Global().Get("document")
// 	document.Call("getElementById", "ChessDiv").Call("setAttribute", "style", "display: none")
// 	div := document.Call("createElement", "div")
// 	div.Call("setAttribute", "id", "myBoard")
// 	div.Call("setAttribute", "style", "width: 400px; float:left")
// 	// document.Get("body").Call("appendChild", div)
// 	document.Call("getElementById", "chessGame").Call("appendChild", div)

// 	info := document.Call("createElement", "div")
// 	info.Call("setAttribute", "style", "width: 500px")

// 	fen := document.Call("createElement", "div")
// 	fen.Call("setAttribute", "id", "fen")
// 	// fen.Call("setAttribute", "class", "space")
// 	// document.Get("body").Call("appendChild", fen)
// 	info.Call("appendChild", fen)

// 	space1 := document.Call("createElement", "div")
// 	space1.Call("setAttribute", "class", "moreSpace")
// 	info.Call("appendChild", space1)

// 	status := document.Call("createElement", "div")
// 	status.Call("setAttribute", "id", "status")
// 	// status.Call("setAttribute", "class", "space")
// 	// document.Get("body").Call("appendChild", status)
// 	// document.Call("getElementById", "chessGame")
// 	info.Call("appendChild", status)

// 	space2 := document.Call("createElement", "div")
// 	space2.Call("setAttribute", "class", "moreSpace")
// 	info.Call("appendChild", space2)

// 	pgn := document.Call("createElement", "div")
// 	pgn.Call("setAttribute", "id", "pgn")
// 	// pgn.Call("setAttribute", "style", "float:left")
// 	// document.Get("body").Call("appendChild", pgn)
// 	info.Call("appendChild", pgn)

// 	document.Call("getElementById", "chessGame").Call("appendChild", info)

// 	role := document.Call("createElement", "script")
// 	roleString := ("var role = \"" + color + "\"")
// 	role.Set("innerText", roleString)

// 	chess := document.Call("createElement", "script")
// 	// chess.Call("setAttribute", "id", "chess")
// 	// "var move = game.move({ from: source, to: target, promotion: 'q'});"+
// 	// " game.undo(); "+
// 	// " if (move === null) return 'snapback';"+
// 	chess.Set("innerText", "var board = null; var game = new Chess();"+
// 		"var $status = $('#status');"+
// 		"var $fen = $('#fen');"+
// 		"var $pgn = $('#pgn'); "+
// 		"function onDragStart (source, piece, position, orientation) {"+
// 		" if (game.game_over()) return false;"+
// 		" if ((game.turn() === 'w' && role === 'black') || (game.turn() === 'b' && role === 'white' ) || (role === 'black' && piece.search(/^w/) !== -1) || (role === 'white' && piece.search(/^b/) !== -1) || (role === 'spectate')) {return false}};"+
// 		" function onDrop (source, target) { "+
// 		" var possibleMoves = game.moves({square: source});"+
// 		" var allowed = false;"+
// 		" for (i = 0; i < possibleMoves.length; i++) {"+
// 		" if (possibleMoves[i].includes(target)) { "+
// 		" allowed = true; break;};};"+
// 		" if(allowed === false) { return 'snapback';};"+
// 		" document.getElementById(\"command\").value = \"chess\" + source + \"fromTo\" + target;"+
// 		" GetCommand('command'); };"+
// 		" function updateStatus () {"+
// 		"var status = '';"+
// 		"var moveColor = 'White';"+
// 		"if (game.turn() === 'b') {moveColor = 'Black'};"+
// 		"if (game.in_checkmate()) { "+
// 		"status = 'Game over, ' + moveColor + ' is in checkmate.'}"+
// 		"	else if (game.in_draw()) {"+
// 		"status = 'Game over, drawn position'} "+
// 		"else {status = moveColor + ' to move';"+
// 		" if (game.in_check()) {"+
// 		"status += ', ' + moveColor + ' is in check'		  }		};"+
// 		"$status.html(status);"+
// 		"$fen.html(game.fen());"+
// 		"$pgn.html(game.pgn());"+
// 		"board.position(game.fen())	  }; "+
// 		"function onSnapEnd () {board.position(game.fen())};"+
// 		" var config = {draggable: true,position: 'start',onDragStart: onDragStart,onDrop: onDrop, onSnapEnd: onSnapEnd, orientation: '"+color+"'};"+
// 		"board = Chessboard('myBoard', config);	  "+
// 		"updateStatus()")
// 	document.Get("body").Call("appendChild", chess)
// 	document.Get("body").Call("appendChild", role)
// }

// func CreateChess(this js.Value, args []js.Value) interface{} {
// 	vsID := args[0].String()
// 	chessVS, err := strconv.Atoi(vsID)
// 	if err != nil {
// 		return nil
// 	}

// 	// fmt.Println(hotstuff.ID(chessVS))
// 	SendStringTo("startChessWhite", hotstuff.ID(chessVS))

// 	for id, _ := range peerMap {
// 		if id != serverID && id != hotstuff.ID(chessVS) {
// 			SendStringTo("startChessSpectate", id)
// 		}
// 	}
// 	srv.Chess = true
// 	CreateChessBoard("black")
// 	return nil
// }

// func ChessTest(this js.Value, args []js.Value) interface{} {

// 	from := "e2"
// 	to := "e4"

// 	document := js.Global().Get("document")

// 	game := js.Global().Get("game")
// 	board := js.Global().Get("board")
// 	chessCmd := "ChessCMD = {from: '" + from + "', to: '" + to + "', promotion: 'q'}"
// 	move := document.Call("createElement", "script")
// 	move.Set("innerText", chessCmd)
// 	document.Get("body").Call("appendChild", move)
// 	// AppendCmd(chessCmd)
// 	game.Call("move", js.Global().Get("ChessCMD"))
// 	board.Call("position", game.Call("fen"))
// 	fmt.Println("Chess executed")
// 	document.Get("body").Call("removeChild", move)

// 	return nil
// }

// func registerCallbacks() {
// 	js.Global().Set("GetSelfID", js.FuncOf(GetSelfID))
// 	js.Global().Set("GetBlockNumber", js.FuncOf(GetBlockNumber))
// 	js.Global().Set("GetCommand", js.FuncOf(GetCommand))
// 	js.Global().Set("PassUint8ArrayToGo", js.FuncOf(PassUint8ArrayToGo))
// 	js.Global().Set("SetUint8ArrayInGo", js.FuncOf(SetUint8ArrayInGo))
// 	js.Global().Set("GetArraySize", js.FuncOf(GetArraySize))
// 	js.Global().Set("StartAgain", js.FuncOf(StartAgain))
// 	js.Global().Set("CreateChess", js.FuncOf(CreateChess))
// 	// js.Global().Set("ChessTest", js.FuncOf(ChessTest))
// }

// defer elapsed("GetSelfID")()
func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}
