package main

// import (
// 	// "github.com/relab/gorums"
// 	"context"
// 	"crypto/ecdsa"
// 	"crypto/tls"
// 	"crypto/x509"
// 	"encoding/pem"
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"
// 	"os/signal"
// 	"strconv"
// 	"sync"
// 	"syscall/js"
// 	"time"

// 	"github.com/HotstuffWASM/newNetwork/config"
// 	gorumsrv "github.com/HotstuffWASM/newNetwork/gorumSrv"
// 	"github.com/HotstuffWASM/newNetwork/logging"
// 	"github.com/hotstuff/client"
// 	"github.com/hotstuff/consensus/chainedhotstuff"
// 	"github.com/hotstuff/leaderrotation"
// 	"github.com/hotstuff/synchronizer"
// 	"github.com/relab/gorums"
// 	"github.com/relab/hotstuff"
// 	"google.golang.org/grpc/credentials"
// )

// // cmdID is a unique identifier for a command
// type cmdID struct {
// 	clientID    uint32
// 	sequenceNum uint64
// }

// type options struct {
// 	Privkey         string
// 	Cert            string
// 	SelfID          hotstuff.ID   `mapstructure:"self-id"`
// 	PmType          string        `mapstructure:"pacemaker"`
// 	LeaderID        hotstuff.ID   `mapstructure:"leader-id"`
// 	Schedule        []hotstuff.ID `mapstructure:"leader-schedule"`
// 	ViewChange      int           `mapstructure:"view-change"`
// 	ViewTimeout     int           `mapstructure:"view-timeout"`
// 	BatchSize       int           `mapstructure:"batch-size"`
// 	PrintThroughput bool          `mapstructure:"print-throughput"`
// 	PrintCommands   bool          `mapstructure:"print-commands"`
// 	ClientAddr      string        `mapstructure:"client-listen"`
// 	PeerAddr        string        `mapstructure:"peer-listen"`
// 	TLS             bool
// 	Interval        int
// 	Output          string
// 	Replicas        []struct {
// 		ID         hotstuff.ID
// 		PeerAddr   string `mapstructure:"peer-address"`
// 		ClientAddr string `mapstructure:"client-address"`
// 		Pubkey     string
// 		Cert       string
// 	}
// }

// type srv struct {
// 	ctx       context.Context
// 	cancel    context.CancelFunc
// 	conf      *options
// 	gorumsSrv *gorums.Server
// 	hsSrv     *gorumsrv.Server
// 	cfg       *gorumsrv.Config
// 	hs        hotstuff.Consensus
// 	pm        hotstuff.ViewSynchronizer
// 	cmdCache  *cmdCache

// 	mut          sync.Mutex
// 	finishedCmds map[cmdID]chan struct{}

// 	lastExecTime int64
// }

// type manualConfig struct {
// 	SelfID  uint32
// 	servers map[uint32]*servers
// }
// type servers struct {
// 	ID         uint32
// 	PeerAddr   string
// 	ClientAddr string
// 	PubKey     *ecdsa.PublicKey
// 	Cert       *tls.Certificate
// 	CertPEM    []byte
// 	PrivKey    *ecdsa.PrivateKey
// }

// var serverID uint32
// var setup manualConfig

// func main() {

// 	fmt.Println("Initializing")
// 	registerCallbacks()

// 	serverID = uint32(0)
// 	for {
// 		if serverID != 0 {
// 			break
// 		}
// 		fmt.Print("Sleeping ZzZ ID: ")
// 		fmt.Println(serverID)
// 		time.Sleep(1 * time.Second)
// 	}

// 	var conf options

// 	conf.SelfID = hotstuff.ID(serverID)
// 	conf.PmType = "fixed"
// 	conf.LeaderID = 1
// 	conf.Schedule = make([]hotstuff.ID, 4)
// 	conf.Schedule[0] = hotstuff.ID(1)
// 	conf.Schedule[1] = hotstuff.ID(2)
// 	conf.Schedule[2] = hotstuff.ID(3)
// 	conf.Schedule[3] = hotstuff.ID(4)
// 	conf.ViewChange = 1

// 	pubkeyString1 := "-----BEGIN HOTSTUFF PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEyaKwozY7C9LL4CAGyuY3gQvHrysu\nkW2YuGfGHvgumwRANtalltLIWEQ5OS2ewsR2xastcb/gzUBtyj54Mi1saw==\n-----END HOTSTUFF PUBLIC KEY-----"
// 	pubBlock1, _ := pem.Decode([]byte(pubkeyString1))
// 	pubkeyBytes1 := pubBlock1.Bytes
// 	genKey1, _ := x509.ParsePKIXPublicKey(pubkeyBytes1)
// 	publicKey1 := genKey1.(*ecdsa.PublicKey)

// 	pubkeyString2 := "-----BEGIN HOTSTUFF PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEitP7/gomqGK/TSgALUpy+MO9N/n1\nvzyHYXvdwPRFPOyS79UEJIYfNCRyex+TRmtB+jwwo1A+x7hCdk2azaF7FA==\n-----END HOTSTUFF PUBLIC KEY-----"
// 	pubBlock2, _ := pem.Decode([]byte(pubkeyString2))
// 	pubkeyBytes2 := pubBlock2.Bytes
// 	genKey2, _ := x509.ParsePKIXPublicKey(pubkeyBytes2)
// 	publicKey2 := genKey2.(*ecdsa.PublicKey)

// 	pubkeyString3 := "-----BEGIN HOTSTUFF PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE6RscDCsZrjOUnRuoUrONyPckRVoo\nt+oGPFjNBynLAWtT07yBCPWYUwzM7Zn+IM3KyAxN12UVZd4itCkGfmOlAg==\n-----END HOTSTUFF PUBLIC KEY-----"
// 	pubBlock3, _ := pem.Decode([]byte(pubkeyString3))
// 	pubkeyBytes3 := pubBlock3.Bytes
// 	genKey3, _ := x509.ParsePKIXPublicKey(pubkeyBytes3)
// 	publicKey3 := genKey3.(*ecdsa.PublicKey)

// 	pubkeyString4 := "-----BEGIN HOTSTUFF PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEVgiSObm49gLvwQbqrnNO67nqnSD7\nifeYUWR2o3Z+5fLPD1msFn/PouMBfK0Epjr/MBFiFpVBtM8D+D/RJPVKdg==\n-----END HOTSTUFF PUBLIC KEY-----"
// 	pubBlock4, _ := pem.Decode([]byte(pubkeyString4))
// 	pubkeyBytes4 := pubBlock4.Bytes
// 	genKey4, _ := x509.ParsePKIXPublicKey(pubkeyBytes4)
// 	publicKey4 := genKey4.(*ecdsa.PublicKey)

// 	//Private Keys
// 	privkeyString1 := "-----BEGIN HOTSTUFF PRIVATE KEY-----\nMHcCAQEEICYSXzL1em20GwmW5f5u54V8wddf5uZ/FN+3iQPi0OIroAoGCCqGSM49\nAwEHoUQDQgAEyaKwozY7C9LL4CAGyuY3gQvHrysukW2YuGfGHvgumwRANtalltLI\nWEQ5OS2ewsR2xastcb/gzUBtyj54Mi1saw==\n-----END HOTSTUFF PRIVATE KEY-----"
// 	privBlock1, _ := pem.Decode([]byte(privkeyString1))
// 	privkeyBytes1 := privBlock1.Bytes
// 	privkeyPEM1 := []byte(privkeyString1)
// 	privateKey1, _ := x509.ParseECPrivateKey(privkeyBytes1)

// 	privkeyString2 := "-----BEGIN HOTSTUFF PRIVATE KEY-----\nMHcCAQEEIDZbG1HNK/NejlozKQHeLYFTFPPi0QYFHNRP/OwlVLB4oAoGCCqGSM49\nAwEHoUQDQgAEitP7/gomqGK/TSgALUpy+MO9N/n1vzyHYXvdwPRFPOyS79UEJIYf\nNCRyex+TRmtB+jwwo1A+x7hCdk2azaF7FA==\n-----END HOTSTUFF PRIVATE KEY-----"
// 	privBlock2, _ := pem.Decode([]byte(privkeyString2))
// 	privkeyBytes2 := privBlock2.Bytes
// 	privkeyPEM2 := []byte(privkeyString2)
// 	privateKey2, _ := x509.ParseECPrivateKey(privkeyBytes2)

// 	privkeyString3 := "-----BEGIN HOTSTUFF PRIVATE KEY-----\nMHcCAQEEIKubRjLNMfX+L4dnDSPcyQIZz/DBdPOyURXUyMibFr/LoAoGCCqGSM49\nAwEHoUQDQgAE6RscDCsZrjOUnRuoUrONyPckRVoot+oGPFjNBynLAWtT07yBCPWY\nUwzM7Zn+IM3KyAxN12UVZd4itCkGfmOlAg==\n-----END HOTSTUFF PRIVATE KEY-----"
// 	privBlock3, _ := pem.Decode([]byte(privkeyString3))
// 	privkeyBytes3 := privBlock3.Bytes
// 	privkeyPEM3 := []byte(privkeyString3)
// 	privateKey3, _ := x509.ParseECPrivateKey(privkeyBytes3)

// 	privkeyString4 := "-----BEGIN HOTSTUFF PRIVATE KEY-----\nMHcCAQEEIAE318VoY//HCbPSSzeYv69esRxZcoIpu10YYq1+h+FVoAoGCCqGSM49\nAwEHoUQDQgAEVgiSObm49gLvwQbqrnNO67nqnSD7ifeYUWR2o3Z+5fLPD1msFn/P\nouMBfK0Epjr/MBFiFpVBtM8D+D/RJPVKdg==\n-----END HOTSTUFF PRIVATE KEY-----"
// 	privBlock4, _ := pem.Decode([]byte(privkeyString4))
// 	privkeyBytes4 := privBlock4.Bytes
// 	privkeyPEM4 := []byte(privkeyString4)
// 	privateKey4, _ := x509.ParseECPrivateKey(privkeyBytes4)

// 	// Certificates
// 	certString1 := "-----BEGIN CERTIFICATE-----\nMIIBmjCCAUCgAwIBAgIQGdrdEJSbdGkA0Tc1VEgYQTAKBggqhkjOPQQDAjArMSkw\nJwYDVQQDEyBIb3RTdHVmZiBTZWxmLVNpZ25lZCBDZXJ0aWZpY2F0ZTAeFw0yMTAx\nMTMxMTExMzJaFw0zMTAxMTMxMTExMzJaMCsxKTAnBgNVBAMTIEhvdFN0dWZmIFNl\nbGYtU2lnbmVkIENlcnRpZmljYXRlMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\nyaKwozY7C9LL4CAGyuY3gQvHrysukW2YuGfGHvgumwRANtalltLIWEQ5OS2ewsR2\nxastcb/gzUBtyj54Mi1sa6NGMEQwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoG\nCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwDwYDVR0RBAgwBocEfwAAATAKBggqhkjO\nPQQDAgNIADBFAiB2RyAzqIE80aGGtEgEe9k98k7K1x1Q1z41oNEzMHSTmAIhAPJD\nRvgWsBp/hqtV2/PZUL+zoqOAoexkotun/5SV5ZdY\n-----END CERTIFICATE-----"
// 	// certBlock1, _ := pem.Decode([]byte(certString1))
// 	// certBytes1 := certBlock1.Bytes
// 	certPEM1 := []byte(certString1)
// 	cert1, _ := tls.X509KeyPair(certPEM1, privkeyPEM1)

// 	certString2 := "-----BEGIN CERTIFICATE-----\nMIIBmjCCAUCgAwIBAgIQC2t9rKAzWVtTDdJnDInLHDAKBggqhkjOPQQDAjArMSkw\nJwYDVQQDEyBIb3RTdHVmZiBTZWxmLVNpZ25lZCBDZXJ0aWZpY2F0ZTAeFw0yMTAx\nMTMxMTExMzJaFw0zMTAxMTMxMTExMzJaMCsxKTAnBgNVBAMTIEhvdFN0dWZmIFNl\nbGYtU2lnbmVkIENlcnRpZmljYXRlMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\nitP7/gomqGK/TSgALUpy+MO9N/n1vzyHYXvdwPRFPOyS79UEJIYfNCRyex+TRmtB\n+jwwo1A+x7hCdk2azaF7FKNGMEQwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoG\nCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwDwYDVR0RBAgwBocEfwAAATAKBggqhkjO\nPQQDAgNIADBFAiAYQV75FDVJJVgjm6WxVV5PhghT4NlF2PRHb4/ATS1QPAIhAOBc\niM4qVWLaB8KlbgMD0pWAFy+l3w0cHPoICEQTySQ+\n-----END CERTIFICATE-----"
// 	// certBlock2, _ := pem.Decode([]byte(certString2))
// 	// certBytes2 := certBlock2.Bytes
// 	certPEM2 := []byte(certString2)
// 	cert2, _ := tls.X509KeyPair(certPEM2, privkeyPEM2)

// 	certString3 := "-----BEGIN CERTIFICATE-----\nMIIBmzCCAUGgAwIBAgIRANE1Qm5JZFIqJmOAwduDsiYwCgYIKoZIzj0EAwIwKzEp\nMCcGA1UEAxMgSG90U3R1ZmYgU2VsZi1TaWduZWQgQ2VydGlmaWNhdGUwHhcNMjEw\nMTEzMTExMTMyWhcNMzEwMTEzMTExMTMyWjArMSkwJwYDVQQDEyBIb3RTdHVmZiBT\nZWxmLVNpZ25lZCBDZXJ0aWZpY2F0ZTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA\nBOkbHAwrGa4zlJ0bqFKzjcj3JEVaKLfqBjxYzQcpywFrU9O8gQj1mFMMzO2Z/iDN\nysgMTddlFWXeIrQpBn5jpQKjRjBEMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAK\nBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA8GA1UdEQQIMAaHBH8AAAEwCgYIKoZI\nzj0EAwIDSAAwRQIhAK0xFL0o7gBFstfJmAvt2k2DYICPzI9JAjmBqMle55T5AiAN\nKcfe5MQ7noWfVkyte60WxWU5Lw2pRDOmIOiG/yXPdg==\n-----END CERTIFICATE-----"
// 	// certBlock3, _ := pem.Decode([]byte(certString3))
// 	// certBytes3 := certBlock3.Bytes
// 	certPEM3 := []byte(certString3)
// 	cert3, _ := tls.X509KeyPair(certPEM3, privkeyPEM3)

// 	certString4 := "-----BEGIN CERTIFICATE-----\nMIIBmzCCAUGgAwIBAgIRAP6OtVIpSKXwu9dCxSQUBRcwCgYIKoZIzj0EAwIwKzEp\nMCcGA1UEAxMgSG90U3R1ZmYgU2VsZi1TaWduZWQgQ2VydGlmaWNhdGUwHhcNMjEw\nMTEzMTExMTMyWhcNMzEwMTEzMTExMTMyWjArMSkwJwYDVQQDEyBIb3RTdHVmZiBT\nZWxmLVNpZ25lZCBDZXJ0aWZpY2F0ZTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA\nBFYIkjm5uPYC78EG6q5zTuu56p0g+4n3mFFkdqN2fuXyzw9ZrBZ/z6LjAXytBKY6\n/zARYhaVQbTPA/g/0ST1SnajRjBEMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAK\nBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA8GA1UdEQQIMAaHBH8AAAEwCgYIKoZI\nzj0EAwIDSAAwRQIgccoJdDly+VUGqkDU7wLjTpwYZJtiwIH3nkRhoaWtcOUCIQDE\nJ78TmegrP+YshLGWWpifGE6lMsKnVNWrccBy6ZXjQA==\n-----END CERTIFICATE-----"
// 	// certBlock4, _ := pem.Decode([]byte(certString4))
// 	// certBytes4 := certBlock4.Bytes
// 	certPEM4 := []byte(certString4)
// 	cert4, _ := tls.X509KeyPair(certPEM4, privkeyPEM4)

// 	setup.SelfID = uint32(conf.SelfID)
// 	setup.servers = make(map[uint32]*servers)
// 	setup.servers[1] = &servers{}
// 	setup.servers[1].ID = 1
// 	setup.servers[1].PeerAddr = "127.0.0.1:13371"
// 	setup.servers[1].ClientAddr = "127.0.0.1:23371"
// 	setup.servers[1].Cert = &cert1
// 	setup.servers[1].PubKey = publicKey1
// 	setup.servers[1].PrivKey = privateKey1
// 	setup.servers[1].CertPEM = certPEM1

// 	setup.servers[2] = &servers{}
// 	setup.servers[2].ID = 2
// 	setup.servers[2].PeerAddr = "127.0.0.1:13372"
// 	setup.servers[2].ClientAddr = "127.0.0.1:23372"
// 	setup.servers[2].Cert = &cert2
// 	setup.servers[2].PubKey = publicKey2
// 	setup.servers[2].PrivKey = privateKey2
// 	setup.servers[2].CertPEM = certPEM2

// 	setup.servers[3] = &servers{}
// 	setup.servers[3].ID = 3
// 	setup.servers[3].PeerAddr = "127.0.0.1:13373"
// 	setup.servers[3].ClientAddr = "127.0.0.1:23373"
// 	setup.servers[3].Cert = &cert3
// 	setup.servers[3].PubKey = publicKey3
// 	setup.servers[3].PrivKey = privateKey3
// 	setup.servers[3].CertPEM = certPEM3

// 	setup.servers[4] = &servers{}
// 	setup.servers[4].ID = 4
// 	setup.servers[4].PeerAddr = "127.0.0.1:13374"
// 	setup.servers[4].ClientAddr = "127.0.0.1:23374"
// 	setup.servers[4].Cert = &cert4
// 	setup.servers[4].PubKey = publicKey4
// 	setup.servers[4].PrivKey = privateKey4
// 	setup.servers[4].CertPEM = certPEM4

// 	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
// 	defer cancel()

// 	start(ctx, &conf)
// }

// func start(ctx context.Context, conf *options) {
// 	privkey := setup.servers[setup.SelfID].PrivKey

// 	var creds credentials.TransportCredentials
// 	var tlsCert tls.Certificate
// 	// if conf.TLS {
// 	// 	creds, tlsCert = loadCreds(conf)
// 	// }

// 	var clientAddress string
// 	replicaConfig := config.NewConfig(conf.SelfID, privkey, creds)
// 	for _, r := range setup.servers {

// 		info := &config.ReplicaInfo{
// 			ID:      hotstuff.ID(r.ID),
// 			Address: r.PeerAddr,
// 			PubKey:  r.PubKey,
// 		}

// 		if r.ID == setup.SelfID {
// 			// override own addresses if set
// 			if conf.ClientAddr != "" {
// 				clientAddress = conf.ClientAddr
// 			} else {
// 				clientAddress = r.ClientAddr
// 			}
// 			if conf.PeerAddr != "" {
// 				info.Address = conf.PeerAddr
// 			}
// 		}

// 		replicaConfig.Replicas[hotstuff.ID(r.ID)] = info
// 	}

// 	logging.NameLogger(fmt.Sprintf("hs%d", conf.SelfID))

// 	srv := newServer(conf, replicaConfig, &tlsCert)
// 	err := srv.Start(clientAddress)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Failed to start HotStuff: %v\n", err)
// 		os.Exit(1)
// 	}

// 	<-ctx.Done()
// 	srv.Stop()
// }

// func newServer(conf *options, replicaConfig *config.ReplicaConfig, tlsCert *tls.Certificate) *srv {
// 	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
// 	defer cancel()

// 	serverOpts := []gorums.ServerOption{}

// 	srv := &srv{
// 		ctx:          ctx,
// 		cancel:       cancel,
// 		conf:         conf,
// 		gorumsSrv:    gorums.NewServer(serverOpts...),
// 		cmdCache:     newCmdCache(conf.BatchSize),
// 		finishedCmds: make(map[cmdID]chan struct{}),
// 		lastExecTime: time.Now().UnixNano(),
// 	}

// 	var err error
// 	srv.cfg = gorumsrv.NewConfig(*replicaConfig)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Failed to init gorums backend: %s\n", err)
// 		os.Exit(1)
// 	}

// 	srv.hsSrv = gorumsrv.NewServer(*replicaConfig)

// 	var leaderRotation hotstuff.LeaderRotation
// 	switch conf.PmType {
// 	case "fixed":
// 		leaderRotation = leaderrotation.NewFixed(conf.LeaderID)
// 	case "round-robin":
// 		leaderRotation = leaderrotation.NewRoundRobin(srv.cfg)
// 	default:
// 		fmt.Fprintf(os.Stderr, "Invalid pacemaker type: '%s'\n", conf.PmType)
// 		os.Exit(1)
// 	}
// 	srv.pm = synchronizer.New(leaderRotation, time.Duration(conf.ViewTimeout)*time.Millisecond)
// 	srv.hs = chainedhotstuff.Builder{
// 		Config:       srv.cfg,
// 		Acceptor:     srv.cmdCache,
// 		Executor:     srv,
// 		Synchronizer: srv.pm,
// 		CommandQueue: srv.cmdCache,
// 	}.Build()
// 	// Use a custom server instead of the gorums one
// 	client.RegisterClientServer(srv.gorumsSrv, srv)
// 	return srv
// }

// func (srv *srv) Start(address string) error {
// 	lis, err := net.Listen("tcp", address)
// 	if err != nil {
// 		return err
// 	}

// 	err = srv.hsSrv.Start(srv.hs)
// 	if err != nil {
// 		return err
// 	}

// 	err = srv.cfg.Connect(10 * time.Second)
// 	if err != nil {
// 		return err
// 	}

// 	// sleep so that all replicas can be ready before we start
// 	time.Sleep(time.Duration(srv.conf.ViewTimeout) * time.Millisecond)

// 	srv.pm.Start()

// 	go func() {
// 		err := srv.gorumsSrv.Serve(lis)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}()

// 	return nil
// }

// func (srv *srv) Stop() {
// 	srv.pm.Stop()
// 	srv.cfg.Close()
// 	srv.hsSrv.Stop()
// 	srv.gorumsSrv.Stop()
// 	srv.cancel()
// }

// // GetSelfID returns the ID of the replica
// func GetSelfID(this js.Value, i []js.Value) interface{} {
// 	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()

// 	selfID, _ := strconv.ParseUint(value1, 10, 32)
// 	serverID = uint32(selfID)
// 	fmt.Println(serverID)
// 	return nil
// }

// func registerCallbacks() {
// 	js.Global().Set("GetSelfID", js.FuncOf(GetSelfID))
// }
