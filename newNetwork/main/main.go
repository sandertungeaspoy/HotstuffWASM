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
	"strconv"
	"strings"
	"sync"
	"syscall/js"

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

func main() {
	registerCallbacks()

	peerMap = make(map[hotstuff.ID]*webrtc.DataChannel)

	serverID = hotstuff.ID(0)
	for {
		if serverID != 0 {
			break
		}
		fmt.Print("Sleeping ZzZ ID: ")
		fmt.Println(serverID)
		time.Sleep(1 * time.Second)
	}

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

	leaderRotation := leaderrotation.NewFixed(hotstuff.ID(1))
	pm := synchronizer.New(leaderRotation, time.Duration(50)*time.Second)
	var cfg *server.Config

	srv = server.Server{
		ID:        serverID,
		Addr:      addr[int(serverID)],
		Pm:        pm,
		Cfg:       cfg,
		PubKey:    pubKey[serverID],
		Cert:      cert[serverID],
		CertPEM:   certPEM[int(serverID)],
		PrivKey:   privKey[int(serverID)],
		SendBytes: sendBytes,
		RecvBytes: recvBytes,
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

	go EstablishConnections()

	if srv.ID == srv.Pm.GetLeader(hs.Leaf().GetView()+1) {
		fmt.Println("I am Leader")
		for {

			time.Sleep(time.Millisecond * 100)
			fmt.Println("Waiting for reply from replicas or for new proposal to be made...")
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
				fmt.Println("Formating string to block...")
				block := StringToBlock(obj)
				// fmt.Print("Formated block: ")
				// fmt.Println(block)

				fmt.Println("OnPropose...")
				// fmt.Print(block.Parent)
				pcString, err := srv.Hs.OnPropose(block)
				if err != nil {
					fmt.Println(err)
					continue
				}
				srv.Hs.Finish(block)
				pc := StringToPartialCert(pcString)
				fmt.Println("OnVote...")
				srv.Hs.OnVote(pc)
				fmt.Println("Sending byte...")
				sendLock.Lock()
				SendCommand(blockString)
				sendLock.Unlock()
				fmt.Println("Bytes sent...")
			case <-recieved:
				fmt.Println("Recieved byte...")
				recvLock.Lock()
				newView := strings.Split(string(recvBytes[0]), ":")
				recvLock.Unlock()
				if newView[0] == "NewView" {
					fmt.Println("Recieved timeout from replica")
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
					// fmt.Print("RecvBytes: ")
					// fmt.Println(recvBytes)
					continue
				}
				recvLock.Lock()
				cmd := strings.Split(string(recvBytes[0]), ":")
				if cmd[0] == "Command" {
					fmt.Println("Recieved command: " + cmd[1])
					cmdString := cmd[1]
					if len(recvBytes) > 1 {
						recvBytes = recvBytes[1:]
					} else {
						recvBytes = make([][]byte, 0)
					}
					srv.Cmds.Cmds = append(srv.Cmds.Cmds, hotstuff.Command(cmdString))
					srv.Pm.Proposal <- srv.Hs.Propose()
					recvLock.Unlock()
					continue
				}
				recvLock.Unlock()
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
				fmt.Println("Sending timeout msg to replicas...")
				sendLock.Lock()
				// sendBytes = append(sendBytes, []byte(msgString))
				SendCommand([]byte(msgString))
				sendLock.Unlock()
			}
		}
	} else {
		fmt.Println("I am normal replica")
		for {
			time.Sleep(time.Millisecond * 100)
			fmt.Println("Waiting for proposal from leader...")
			select {
			case <-recieved:
				fmt.Println("Recieved byte from leader...")
				recvLock.Lock()
				newView := strings.Split(string(recvBytes[0]), ":")
				recvLock.Unlock()
				if newView[0] == "NewView" {
					fmt.Println("Recieved timeout from leader...")
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
				fmt.Print("Handle propose for view: ")
				fmt.Println(block.View)
				pcString, err := srv.Hs.OnPropose(block)
				if err != nil {
					fmt.Println(err)
					continue
				}
				srv.Hs.Finish(block)
				pc := StringToPartialCert(pcString)
				srv.Hs.OnVote(pc)
				fmt.Println("Sending PC to leader...")
				sendLock.Lock()
				// sendBytes = append(sendBytes, []byte(pcString))
				SendCommand([]byte(pcString))
				sendLock.Unlock()
			case <-srv.Pm.NewView:
				timeoutview := srv.Hs.NewView()
				srv.Hs.OnNewView(timeoutview)
				msg := NewViewToString(timeoutview)
				fmt.Println("Sending timeout msg to leader...")
				sendLock.Lock()
				// sendBytes = append(sendBytes, []byte(msg))
				SendCommand([]byte(msg))
				sendLock.Unlock()
			case cmd := <-incomingCmd:
				cmdString := "Command:" + cmd

				fmt.Println("Sending command to leader...")
				// sendBytes = append(sendBytes, []byte(cmdString))
				SendCommand([]byte(cmdString))
			}
		}
	}
}

// FormatBytes returns the ID of the sender, the command and the block
func FormatBytes(msg []byte) (id hotstuff.ID, cmd string, obj string) {
	if len(msg) != 0 {
		msgString := string(msg)
		msgStringByte := strings.Split(msgString, ";")
		// fmt.Print("FormatBytes string: ")
		// fmt.Println(msgString)
		// fmt.Print("Byte of msg: ")
		// fmt.Println(msgStringByte[1])

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
	// fmt.Println(s)
	strByte := strings.Split(s, ":")

	signString := strings.Split(strByte[0], "-")
	// fmt.Println(signString)
	rInt := new(big.Int)
	rInt.SetString(signString[0], 0)
	sInt := new(big.Int)
	sInt.SetString(signString[1], 0)
	signer, _ := strconv.ParseUint(signString[2], 10, 32)
	sign := *hsecdsa.NewSignature(rInt, sInt, hotstuff.ID(signer))

	// hash, _ := base64.RawStdEncoding.DecodeString(strByte[2])
	// hash := []byte(strByte[1])
	hash, _ := hex.DecodeString(strByte[1])
	var h [32]byte
	copy(h[:], hash)
	hash2 := hotstuff.Hash(h)
	var pc hotstuff.PartialCert = hsecdsa.NewPartialCert(&sign, hash2)
	// fmt.Print("Pc created: ")
	// fmt.Println(pc)
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
create:
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

	// Register data channel creation handling

	var dc *webrtc.DataChannel

	waiter := make(chan struct{})

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())
			// Register text message handling
			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				// fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
				// err := d.SendText("Message Received")
				recvLock.Lock()
				recvBytes = append(recvBytes, msg.Data)
				recvLock.Unlock()
				recieved <- msg.Data
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

		if attempts > 10 {
			return nil, "error"
		}
		attempts++
		time.Sleep(time.Millisecond * 5000)
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

	if peerConnection.ICEGatheringState() == webrtc.ICEGatheringStateComplete && strings.Contains(peerConnection.LocalDescription().SDP, "c=IN IP4 0.0.0.0") {
		fmt.Println(peerConnection.LocalDescription().SDP)
		DeliverAnswer(peerConnection.LocalDescription().SDP, senderID)

	} else if peerConnection.ICEGatheringState() == webrtc.ICEGatheringStateComplete && !strings.Contains(peerConnection.LocalDescription().SDP, "c=IN IP4 0.0.0.0") {
		goto create
	}

	<-waiter

	removeAnswer(senderID)
	fmt.Println("Returning")
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

	waiter := make(chan struct{})

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", dataChannel.Label(), dataChannel.ID())
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
			RemoveOffer()
		}

	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		// fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
		recvLock.Lock()
		recvBytes = append(recvBytes, msg.Data)
		recvLock.Unlock()
		recieved <- msg.Data
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

	if peerConnection.ICEGatheringState() == webrtc.ICEGatheringStateComplete {
		fmt.Println(peerConnection.LocalDescription().SDP)
		DeliverOffer(peerConnection.LocalDescription().SDP)
	}

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
		if attempts > 10 {
			return nil, "error"
		}
		attempts++
		time.Sleep(time.Millisecond * 5000)
	}

	fmt.Println(answer)

	answersdp := webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: answer}

	// Apply the answer as the remote description
	err = peerConnection.SetRemoteDescription(answersdp)
	if err != nil {
		panic(err)
	}

	fmt.Println("Remote desc set")

	<-waiter
	RemoveOffer()
	fmt.Println("Returning")
	return dataChannel, senderID
}

func DeliverOffer(offer string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:13372", nil)
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

	c, _, err := websocket.Dial(ctx, "ws://localhost:13372", nil)
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

	c, _, err := websocket.Dial(ctx, "ws://localhost:13372", nil)
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
	fmt.Println(offer)

	msgs := strings.Split(offer, "&")

	c.Close(websocket.StatusNormalClosure, "")
	return msgs[0], strings.Split(msgs[1], "%")[0]
}

func ReceiveAnswer() (string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:13372", nil)
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

	fmt.Println(answer)

	msgs := strings.Split(answer, "&")

	c.Close(websocket.StatusNormalClosure, "")
	return msgs[0], strings.Split(msgs[1], "%")[0]
}

func RemoveOffer() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:13372", nil)
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

	c, _, err := websocket.Dial(ctx, "ws://localhost:13372", nil)
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

	c, _, err := websocket.Dial(ctx, "ws://localhost:13371", nil)
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusInternalError, "WebSocket has been closed")

	conn := websocket.NetConn(ctx, c, 1)

	fmt.Fprintf(conn, "")

	c.Close(websocket.StatusNormalClosure, "setup:purgeDatabase\n&0%")
}

func EstablishConnections() {

	started := false

	for {
		if srv.ID == srv.Pm.GetLeader(srv.Hs.Leaf().GetView()+1) {
			if len(peerMap) == 3 {
				time.Sleep(time.Second * 20)
				continue
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
				srv.Pm.Start()
				started = true
			}

		} else {

			if len(peerMap) == 0 {
				dc, leaderID := ConnectToLeader()
				if leaderID == "error" || leaderID == "empty" {
					// time.Sleep(time.Second * 5)
					continue
				}

				leaderIDUint, _ := strconv.ParseUint(leaderID, 10, 32)
				leaderIDHot := hotstuff.ID(leaderIDUint)

				peerMap[leaderIDHot] = dc
				srv.Pm.Start()
			}
			time.Sleep(time.Second * 30)
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
	for _, peer := range peerMap {
		err := peer.Send(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSelfID gets the ID of the server
func GetSelfID(this js.Value, i []js.Value) interface{} {
	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()

	selfID, _ := strconv.ParseUint(value1, 10, 32)
	serverID = hotstuff.ID(selfID)
	fmt.Println(serverID)
	return nil
}

// PassUint8ArrayToGo passes array
func PassUint8ArrayToGo(this js.Value, args []js.Value) interface{} {
	recv := make([]byte, args[0].Get("length").Int())

	_ = js.CopyBytesToGo(recv, args[0])

	recvLock.Lock()
	recvBytes = append(recvBytes, recv)
	recvLock.Unlock()
	recieved <- recv

	return nil
}

// SetUint8ArrayInGo sets array
func SetUint8ArrayInGo(this js.Value, args []js.Value) interface{} {
	sendLock.Lock()
	if len(sendBytes) == 0 {
		sendLock.Unlock()
		return nil
	}

	if args[0].Get("length").Int() == 0 {
		sendLock.Unlock()
		return nil
	}

	var msg []byte
	if len(sendBytes) > 1 {
		msg, sendBytes = sendBytes[0], sendBytes[1:]
	} else {
		msg, sendBytes = sendBytes[0], make([][]byte, 0)
	}
	sendLock.Unlock()
	if msg == nil {
		return nil
	}
	// fmt.Println("Sending bytes to JS")
	_ = js.CopyBytesToJS(args[0], msg)

	return nil
}

// GetArraySize gets the array size
func GetArraySize(this js.Value, args []js.Value) interface{} {

	if len(sendBytes) == 0 {
		_ = js.CopyBytesToJS(args[1], []byte{0})
		return nil
	}
	size := make([]byte, 10)

	msgSize := []byte(strconv.Itoa(len(sendBytes[0])))

	copy(size, msgSize)

	_ = js.CopyBytesToJS(args[0], size)
	_ = js.CopyBytesToJS(args[1], []byte{1})

	return nil
}

// GetCommand gets the ID of the server
func GetCommand(this js.Value, i []js.Value) interface{} {
	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()

	cmd := string(value1)
	cmd = strconv.FormatUint(uint64(serverID), 10) + "cmdID" + cmd
	if serverID == srv.Pm.GetLeader(srv.Hs.Leaf().View+1) {
		cmdLock.Lock()
		command := hotstuff.Command(cmd)
		srv.Cmds.Cmds = append(srv.Cmds.Cmds, command)
		cmdLock.Unlock()
		srv.Pm.Proposal <- srv.Hs.Propose()
	} else {
		incomingCmd <- cmd
	}
	return nil
}

func registerCallbacks() {
	js.Global().Set("GetSelfID", js.FuncOf(GetSelfID))
	js.Global().Set("GetCommand", js.FuncOf(GetCommand))
	js.Global().Set("PassUint8ArrayToGo", js.FuncOf(PassUint8ArrayToGo))
	js.Global().Set("SetUint8ArrayInGo", js.FuncOf(SetUint8ArrayInGo))
	js.Global().Set("GetArraySize", js.FuncOf(GetArraySize))
}

// defer elapsed("GetSelfID")()
func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}
