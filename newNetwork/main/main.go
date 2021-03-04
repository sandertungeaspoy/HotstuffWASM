package main

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"math/big"
	"strconv"
	"strings"
	"time"

	hotstuff "github.com/HotstuffWASM/newNetwork"
	"github.com/HotstuffWASM/newNetwork/config"
	consensus "github.com/HotstuffWASM/newNetwork/consensus"
	hsecdsa "github.com/HotstuffWASM/newNetwork/crypto/ecdsa"
	"github.com/HotstuffWASM/newNetwork/leaderrotation"
	server "github.com/HotstuffWASM/newNetwork/server"
	synchronizer "github.com/HotstuffWASM/newNetwork/synchronizer"
)

var serverID hotstuff.ID = 0

func main() {
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
	pubKey[1] = publicKey1
	pubKey[2] = publicKey2
	pubKey[3] = publicKey3
	pubKey[4] = publicKey4

	var cert []*tls.Certificate
	cert[1] = &cert1
	cert[2] = &cert2
	cert[3] = &cert3
	cert[4] = &cert4

	var privKey []*ecdsa.PrivateKey
	privKey[1] = privateKey1
	privKey[2] = privateKey2
	privKey[3] = privateKey3
	privKey[4] = privateKey4

	var certPEM [][]byte
	certPEM[1] = certPEM1
	certPEM[2] = certPEM2
	certPEM[3] = certPEM3
	certPEM[4] = certPEM4

	var addr []string
	addr[1] = "127.0.0.1:13371"
	addr[2] = "127.0.0.1:13372"
	addr[3] = "127.0.0.1:13373"
	addr[4] = "127.0.0.1:13374"

	leaderRotation := leaderrotation.NewFixed(hotstuff.ID(1))
	pm := synchronizer.New(leaderRotation, time.Duration(1000)*time.Millisecond)
	var cfg *server.Config
	var sendBytes [][]byte
	var recvBytes [][]byte

	srv := server.Server{
		ID:        serverID,
		Addr:      addr[int(serverID)],
		Pm:        *pm,
		Cfg:       cfg,
		PubKey:    pubKey[serverID],
		Cert:      cert[serverID],
		CertPEM:   certPEM[int(serverID)],
		PrivKey:   privKey[int(serverID)],
		SendBytes: sendBytes,
		RecvBytes: recvBytes,
	}

	hs := consensus.Builder{
		Config:       srv.Cfg,
		Acceptor:     &srv.Cmds,
		Executor:     &srv,
		Synchronizer: &srv.Pm,
		CommandQueue: &srv.Cmds,
	}.Build()

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

	srv.Hs = hs

	srv.Pm.Init(srv.Hs)

	srv.Pm.Start()

	if srv.ID == srv.Pm.GetLeader(hs.Leaf().GetView()+1) {
		for {
			select {
			case msgByte := <-srv.Pm.Proposal:

				senderID, cmd, obj := FormatBytes(msgByte)
				if senderID != srv.ID && cmd != "Propose" {
					return
				}
				block := StringToBlock(obj)

				pcString, err := srv.Hs.OnPropose(block)
				if err != nil {
					panic(err)
				}
				// fmt.Println(senderID, cmd, pc)
				srv.Hs.Finish(block)
				pc := StringToPartialCert(pcString)
				srv.Hs.OnVote(pc)
				sendBytes = append(sendBytes, msgByte)
			}
			if len(recvBytes) != 0 {
				_, cmd, obj := FormatBytes(recvBytes[0])
				if cmd != " PartialCert " {
					return
				}
				pc := StringToPartialCert(obj)
				srv.Hs.OnVote(pc)
				recvBytes = recvBytes[1:]
			}
		}
	} else {
		for {
			if len(recvBytes) != 0 {
				_, cmd, obj := FormatBytes(recvBytes[0])
				if cmd != "Propose" {
					return
				}
				block := StringToBlock(obj)
				// id, err := srv.GetID()
				// if err != nil {
				// 	panic(err)
				// }
				// block.Proposer = id
				pcString, err := srv.Hs.OnPropose(block)
				if err != nil {
					panic(err)
				}
				// fmt.Println(senderID, cmd, pc)
				srv.Hs.Finish(block)
				sendBytes = append(sendBytes, []byte(pcString))
				recvBytes = recvBytes[1:]
			}
		}
	}
}

// FormatBytes returns the ID of the sender, the command and the block
func FormatBytes(msg []byte) (id hotstuff.ID, cmd string, obj string) {
	msgString := hex.EncodeToString(msg)
	msgStringByte := strings.Split(msgString, " ")

	idString, _ := strconv.ParseUint(msgStringByte[1], 10, 32)
	id = hotstuff.ID(idString)

	cmd = msgStringByte[2]

	obj = msgStringByte[3]

	return id, cmd, obj
}

// StringToBlock returns a block from a given string
func StringToBlock(s string) *hotstuff.Block {
	strByte := strings.Split(s, ":")
	parent, _ := base64.RawStdEncoding.DecodeString(strByte[1])
	var p [32]byte
	copy(p[:], parent)
	parent2 := hotstuff.Hash(p)
	proposer, _ := strconv.ParseUint(strByte[2], 10, 32)
	cmd := hotstuff.Command(strByte[3])
	certHash := []byte(strByte[5])
	var c [32]byte
	copy(c[:], certHash)
	certHash2 := hotstuff.Hash(c)
	var sig map[hotstuff.ID]*hsecdsa.Signature
	sigString := strByte[4]

	sigBytes := strings.Split(sigString, "\n")
	for i := 0; i < len(sigBytes); i++ {
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

	hash, _ := base64.RawStdEncoding.DecodeString(strByte[2])
	var h [32]byte
	copy(h[:], hash)
	hash2 := hotstuff.Hash(h)
	var pc hotstuff.PartialCert = hsecdsa.NewPartialCert(&sign, hash2)
	return pc
}
