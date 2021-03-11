package main

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"

	// "syscall/js"
	"time"

	hotstuff "github.com/HotstuffWASM/newNetwork"
	"github.com/HotstuffWASM/newNetwork/config"
	consensus "github.com/HotstuffWASM/newNetwork/consensus"
	hsecdsa "github.com/HotstuffWASM/newNetwork/crypto/ecdsa"
	"github.com/HotstuffWASM/newNetwork/leaderrotation"
	server "github.com/HotstuffWASM/newNetwork/server"
	synchronizer "github.com/HotstuffWASM/newNetwork/synchronizer"
)

var sendBytes [][]byte
var recvBytes [][]byte
var serverID hotstuff.ID = 0
var recieved chan []byte

var conn net.Conn
var conn2 net.Conn
var conn3 net.Conn
var conn4 net.Conn

func main() {
	// registerCallbacks()

	idString := os.Args[1]
	idUint, _ := strconv.ParseUint(idString, 10, 32)
	serverID = hotstuff.ID(idUint)

	// serverID = hotstuff.ID(0)
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
	recieved = make(chan []byte, 64)
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

	if idUint == 1 {
		lis, errL := net.Listen("tcp", addr[1])
		if errL != nil {
			fmt.Println(errL)
		}
		fmt.Println("Waiting for conn2")
		var err error
		conn2, err = lis.Accept()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Waiting for conn3")
		conn3, err = lis.Accept()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Waiting for conn4")
		conn4, err = lis.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go handleConn()
	} else {
		var err error
		conn, err = net.Dial("tcp", addr[1])
		if err != nil {
			fmt.Println(err)
		}
		go handleConn()
	}

	leaderRotation := leaderrotation.NewFixed(hotstuff.ID(1))
	pm := synchronizer.New(leaderRotation, time.Duration(2)*time.Second)
	var cfg *server.Config

	srv := server.Server{
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

	srv.Pm.Start()

	if srv.ID == srv.Pm.GetLeader(hs.Leaf().GetView()+1) {
		fmt.Println("I am Leader")
		for {
			// time.Sleep(time.Millisecond * 10)
			fmt.Println("Waiting for reply from replicas or for new proposal to be made...")
			select {
			case msgByte := <-srv.Pm.Proposal:
				fmt.Println("Read channel msg...")
				// fmt.Println(msgByte)
				senderID, cmd, obj := FormatBytes(msgByte)
				fmt.Println("FormatBytes result: ")
				fmt.Print("SenderID: ")
				fmt.Println(senderID)
				fmt.Print("cmd: ")
				fmt.Println(cmd)
				fmt.Print("obj: ")
				fmt.Println(obj)
				if senderID != srv.ID && cmd != "Propose" {
					continue
				}
				fmt.Println("Formating string to block...")
				block := StringToBlock(obj)
				fmt.Print("Formated block: ")
				fmt.Println(block)

				fmt.Println("OnPropose...")
				// fmt.Print(block.Parent)
				pcString, err := srv.Hs.OnPropose(block)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print("SenderID: ")
				fmt.Println(senderID)
				fmt.Print("cmd: ")
				fmt.Println(cmd)
				fmt.Print("pcString: ")
				fmt.Println(pcString)
				// fmt.Println("Finish")
				srv.Hs.Finish(block)
				// fmt.Println("Finish done")
				pc := StringToPartialCert(pcString)
				fmt.Println("OnVote...")
				srv.Hs.OnVote(pc)
				fmt.Println("Sending byte...")
				fmt.Println(runtime.NumGoroutine())
				sendBytes = append(sendBytes, msgByte)
				fmt.Println("Bytes sent...")
			case <-recieved:
				fmt.Println("Recieved byte...")
				// _, cmd, obj := FormatBytes(recvBytes[0])
				// if cmd != " PartialCert " {
				// 	continue
				// }
				if len(recvBytes) == 0 {
					continue
				}
				pc := StringToPartialCert(string(recvBytes[0]))
				if len(recvBytes) > 1 {
					recvBytes = recvBytes[1:]
				} else {
					recvBytes = make([][]byte, 0)
				}
				srv.Hs.OnVote(pc)
			case ok := <-srv.Pm.NewView:
				fmt.Println("Timeout -> start new view")
				if ok {
					srv.Hs.NewView()
				}
			}
		}
	} else {
		fmt.Println("I am normal replica")
		for {
			time.Sleep(time.Millisecond * 100)
			// fmt.Println("Waiting for proposal from leader...")
			select {
			case <-recieved:
				fmt.Println("Recieved byte from leader...")
				id, cmd, obj := FormatBytes(recvBytes[0])
				if len(recvBytes) > 1 {
					recvBytes = recvBytes[1:]
				} else {
					recvBytes = make([][]byte, 0)
				}
				if id == hotstuff.ID(0) || cmd != "Propose" {
					continue
				}
				block := StringToBlock(obj)
				pcString, err := srv.Hs.OnPropose(block)
				if err != nil {
					fmt.Println(err)
				}
				srv.Hs.Finish(block)
				sendBytes = append(sendBytes, []byte(pcString))
			}
		}
	}
}

// FormatBytes returns the ID of the sender, the command and the block
func FormatBytes(msg []byte) (id hotstuff.ID, cmd string, obj string) {
	if len(msg) != 0 {
		msgString := string(msg)
		msgStringByte := strings.Split(msgString, " ")
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

// GetSelfID gets the ID of the server
// func GetSelfID(this js.Value, i []js.Value) interface{} {
// 	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()

// 	selfID, _ := strconv.ParseUint(value1, 10, 32)
// 	serverID = hotstuff.ID(selfID)
// 	fmt.Println(serverID)
// 	return nil
// }

// // PassUint8ArrayToGo passes array
// func PassUint8ArrayToGo(this js.Value, args []js.Value) interface{} {

// 	recv := make([]byte, args[0].Get("length").Int())

// 	_ = js.CopyBytesToGo(recv, args[0])

// 	recvBytes = append(recvBytes, recv)
// 	recieved <- recv

// 	return nil
// }

// // SetUint8ArrayInGo sets array
// func SetUint8ArrayInGo(this js.Value, args []js.Value) interface{} {

// 	if len(sendBytes) == 0 {
// 		return nil
// 	}

// 	var msg []byte

// 	if len(sendBytes) > 1 {
// 		msg, sendBytes = sendBytes[0], sendBytes[1:]
// 	} else {
// 		msg, sendBytes = sendBytes[0], make([][]byte, 0)
// 	}
// 	if msg == nil {
// 		return nil
// 	}
// 	_ = js.CopyBytesToJS(args[0], msg)

// 	return nil
// }

// // GetArraySize gets the array size
// func GetArraySize(this js.Value, args []js.Value) interface{} {

// 	if len(sendBytes) == 0 {
// 		return nil
// 	}
// 	size := make([]byte, 10)

// 	msgSize := []byte(strconv.Itoa(len(sendBytes[0])))

// 	copy(size, msgSize)

// 	_ = js.CopyBytesToJS(args[0], size)

// 	return nil
// }

// func registerCallbacks() {
// 	js.Global().Set("GetSelfID", js.FuncOf(GetSelfID))

// 	js.Global().Set("PassUint8ArrayToGo", js.FuncOf(PassUint8ArrayToGo))
// 	js.Global().Set("SetUint8ArrayInGo", js.FuncOf(SetUint8ArrayInGo))
// 	js.Global().Set("GetArraySize", js.FuncOf(GetArraySize))
// }

func handleConn() {
	if int(serverID) == 1 {
		go func() {
			var msg []byte
			for {
				if len(sendBytes) > 1 {
					msg, sendBytes = sendBytes[0], sendBytes[1:]
					_, err := conn2.Write(msg)
					if err != nil {
						fmt.Println(err)
					}
					_, err = conn3.Write(msg)
					if err != nil {
						fmt.Println(err)
					}
					_, err = conn4.Write(msg)
					if err != nil {
						fmt.Println(err)
					}
				} else if len(sendBytes) == 1 {
					msg, sendBytes = sendBytes[0], make([][]byte, 0)
					_, err := conn2.Write(msg)
					if err != nil {
						fmt.Println(err)
					}
					_, err = conn3.Write(msg)
					if err != nil {
						fmt.Println(err)
					}
					_, err = conn4.Write(msg)
					if err != nil {
						fmt.Println(err)
					}
				}
			}

		}()

		go func() {
			for {
				buff := make([]byte, 4096)
				n, err := conn2.Read(buff)
				if err != nil {
					fmt.Println(err)
				}
				if n > 0 {
					res := make([]byte, n)
					copy(res, buff[:n])
					recvBytes = append(recvBytes, res)
					recieved <- res
				}
			}
		}()
		go func() {
			for {
				buff := make([]byte, 4096)
				n, err := conn3.Read(buff)
				if err != nil {
					fmt.Println(err)
				}
				if n > 0 {
					res := make([]byte, n)
					copy(res, buff[:n])
					recvBytes = append(recvBytes, res)
					recieved <- res
				}
			}
		}()
		go func() {
			for {
				buff := make([]byte, 4096)
				n, err := conn4.Read(buff)
				if err != nil {
					fmt.Println(err)
				}
				if n > 0 {
					res := make([]byte, n)
					copy(res, buff[:n])
					recvBytes = append(recvBytes, res)
					recieved <- res
				}
			}
		}()
	} else {
		go func() {
			var msg []byte
			for {
				if len(sendBytes) > 1 {
					msg, sendBytes = sendBytes[0], sendBytes[1:]
					_, err := conn.Write(msg)
					if err != nil {
						fmt.Println(err)
					}
				} else if len(sendBytes) == 1 {
					msg, sendBytes = sendBytes[0], make([][]byte, 0)
					_, err := conn.Write(msg)
					if err != nil {
						fmt.Println(err)
					}
				}
			}

		}()

		go func() {
			for {
				buff := make([]byte, 4096)
				n, err := conn.Read(buff)
				if err != nil {
					fmt.Println(err)
					return
				}
				if n > 0 {
					res := make([]byte, n)
					copy(res, buff[:n])
					recvBytes = append(recvBytes, res)
					recieved <- res
				}
			}
		}()
	}
}
