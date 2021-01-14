package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	// keyCode "github.com/bitherhq/go-bither/crypto"
	"github.com/felixge/fgprof"
	"github.com/relab/hotstuff"
	"github.com/relab/hotstuff/client"
	"github.com/relab/hotstuff/config"
	"github.com/relab/hotstuff/data"
	"github.com/relab/hotstuff/pacemaker"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type options struct {
	Privkey         string
	Cert            string
	SelfID          config.ReplicaID   `mapstructure:"self-id"`
	PmType          string             `mapstructure:"pacemaker"`
	LeaderID        config.ReplicaID   `mapstructure:"leader-id"`
	Schedule        []config.ReplicaID `mapstructure:"leader-schedule"`
	ViewChange      int                `mapstructure:"view-change"`
	ViewTimeout     int                `mapstructure:"view-timeout"`
	BatchSize       int                `mapstructure:"batch-size"`
	PrintThroughput bool               `mapstructure:"print-throughput"`
	PrintCommands   bool               `mapstructure:"print-commands"`
	ClientAddr      string             `mapstructure:"client-listen"`
	PeerAddr        string             `mapstructure:"peer-listen"`
	TLS             bool
	Interval        int
	Output          string
	Replicas        []struct {
		ID         config.ReplicaID
		PeerAddr   string `mapstructure:"peer-address"`
		ClientAddr string `mapstructure:"client-address"`
		Pubkey     string
		Cert       string
	}
}

type manualConfig struct {
	ID      int
	servers []struct {
		ID         int
		PeerAddr   string
		ClientAddr string
		PubKey     *ecdsa.PublicKey
		Cert       *tls.Certificate
	}
}

func usage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println()
	fmt.Println("Loads configuration from ./hotstuff.toml and file specified by --config")
	fmt.Println()
	fmt.Println("Options:")
	pflag.PrintDefaults()
}

func main() {
	pflag.Usage = usage

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// some configuration options can be set using flags
	help := pflag.BoolP("help", "h", false, "Prints this text.")
	configFile := pflag.String("config", "", "The path to the config file")
	cpuprofile := pflag.String("cpuprofile", "", "File to write CPU profile to")
	memprofile := pflag.String("memprofile", "", "File to write memory profile to")
	fullprofile := pflag.String("fullprofile", "", "File to write fgprof profile to")
	traceFile := pflag.String("trace", "", "File to write execution trace to")
	pflag.Uint32("self-id", 0, "The id for this replica.")
	pflag.Int("view-change", 100, "How many views before leader change with round-robin pacemaker")
	pflag.Int("batch-size", 100, "How many commands are batched together for each proposal")
	pflag.Int("view-timeout", 1000, "How many milliseconds before a view is timed out")
	pflag.String("privkey", "", "The path to the private key file")
	pflag.String("cert", "", "Path to the certificate")
	pflag.Bool("print-commands", false, "Commands will be printed to stdout")
	pflag.Bool("print-throughput", false, "Throughput measurements will be printed stdout")
	pflag.Int("interval", 1000, "Throughput measurement interval in milliseconds")
	pflag.Bool("tls", false, "Enable TLS")
	pflag.String("client-listen", "", "Override the listen address for the client server")
	pflag.String("peer-listen", "", "Override the listen address for the replica (peer) server")
	pflag.Parse()

	if *help {
		pflag.Usage()
		os.Exit(0)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *fullprofile != "" {
		f, err := os.Create(*fullprofile)
		if err != nil {
			log.Fatal("Could not create fgprof profile: ", err)
		}
		defer f.Close()
		stop := fgprof.Start(f, fgprof.FormatPprof)

		defer func() {
			err := stop()
			if err != nil {
				log.Fatal("Could not write fgprof profile: ", err)
			}
		}()
	}

	if *traceFile != "" {
		f, err := os.Create(*traceFile)
		if err != nil {
			log.Fatal("Could not create trace file: ", err)
		}
		defer f.Close()
		if err := trace.Start(f); err != nil {
			log.Fatal("Failed to start trace: ", err)
		}
		defer trace.Stop()
	}

	viper.BindPFlags(pflag.CommandLine)

	// read main config file in working dir
	viper.SetConfigName("hotstuff")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
		os.Exit(1)
	}

	// read secondary config file
	if *configFile != "" {
		viper.SetConfigFile(*configFile)
		err = viper.MergeInConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read secondary config file: %v\n", err)
			os.Exit(1)
		}
	}

	var conf options
	err = viper.Unmarshal(&conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal config: %v\n", err)
		os.Exit(1)
	}
	conf.Privkey = "keys/r1.key"
	conf.SelfID = 1

	// // 	keys\r1.key
	// privkeyString1 := "26125f32f57a6db41b0996e5fe6ee7857cc1d75fe6e67f14dfb78903e2d0e22b"
	// privkeyBytes1, _ := hex.DecodeString(privkeyString1)
	// // privkey1 := keyCode.HexToECDSA(privkeyBytes1)
	// // 	keys\r2.key
	// privkeyString2 := "365b1b51cd2bf35e8e5a332901de2d815314f3e2d106051cd44ffcec2554b078"
	// privkeyBytes2, _ := hex.DecodeString(privkeyString2)
	// // 	keys\r3.key
	// privkeyString3 := "ab9b4632cd31f5fe2f87670d23dcc90219cff0c174f3b25115d4c8c89b16bfcb"
	// privkeyBytes3, _ := hex.DecodeString(privkeyString3)
	// // 	keys\r4.key
	// privkeyString4 := "0137d7c56863ffc709b3d24b3798bfaf5eb11c59728229bb5d1862ad7e87e155"
	// privkeyBytes4, _ := hex.DecodeString(privkeyString4)
	// fmt.Println("key")
	// fmt.Println(conf.Privkey)

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
	certBlock1, _ := pem.Decode([]byte(certString1))
	certBytes1 := certBlock1.Bytes
	cert1, _ := tls.X509KeyPair(certBytes1, privkeyPEM1)

	certString2 := "-----BEGIN CERTIFICATE-----\nMIIBmjCCAUCgAwIBAgIQC2t9rKAzWVtTDdJnDInLHDAKBggqhkjOPQQDAjArMSkw\nJwYDVQQDEyBIb3RTdHVmZiBTZWxmLVNpZ25lZCBDZXJ0aWZpY2F0ZTAeFw0yMTAx\nMTMxMTExMzJaFw0zMTAxMTMxMTExMzJaMCsxKTAnBgNVBAMTIEhvdFN0dWZmIFNl\nbGYtU2lnbmVkIENlcnRpZmljYXRlMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\nitP7/gomqGK/TSgALUpy+MO9N/n1vzyHYXvdwPRFPOyS79UEJIYfNCRyex+TRmtB\n+jwwo1A+x7hCdk2azaF7FKNGMEQwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoG\nCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwDwYDVR0RBAgwBocEfwAAATAKBggqhkjO\nPQQDAgNIADBFAiAYQV75FDVJJVgjm6WxVV5PhghT4NlF2PRHb4/ATS1QPAIhAOBc\niM4qVWLaB8KlbgMD0pWAFy+l3w0cHPoICEQTySQ+\n-----END CERTIFICATE-----"
	certBlock2, _ := pem.Decode([]byte(certString2))
	certBytes2 := certBlock2.Bytes
	cert2, _ := tls.X509KeyPair(certBytes2, privkeyPEM2)

	certString3 := "-----BEGIN CERTIFICATE-----\nMIIBmzCCAUGgAwIBAgIRANE1Qm5JZFIqJmOAwduDsiYwCgYIKoZIzj0EAwIwKzEp\nMCcGA1UEAxMgSG90U3R1ZmYgU2VsZi1TaWduZWQgQ2VydGlmaWNhdGUwHhcNMjEw\nMTEzMTExMTMyWhcNMzEwMTEzMTExMTMyWjArMSkwJwYDVQQDEyBIb3RTdHVmZiBT\nZWxmLVNpZ25lZCBDZXJ0aWZpY2F0ZTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA\nBOkbHAwrGa4zlJ0bqFKzjcj3JEVaKLfqBjxYzQcpywFrU9O8gQj1mFMMzO2Z/iDN\nysgMTddlFWXeIrQpBn5jpQKjRjBEMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAK\nBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA8GA1UdEQQIMAaHBH8AAAEwCgYIKoZI\nzj0EAwIDSAAwRQIhAK0xFL0o7gBFstfJmAvt2k2DYICPzI9JAjmBqMle55T5AiAN\nKcfe5MQ7noWfVkyte60WxWU5Lw2pRDOmIOiG/yXPdg==\n-----END CERTIFICATE-----"
	certBlock3, _ := pem.Decode([]byte(certString3))
	certBytes3 := certBlock3.Bytes
	cert3, _ := tls.X509KeyPair(certBytes3, privkeyPEM3)

	certString4 := "-----BEGIN CERTIFICATE-----\nMIIBmzCCAUGgAwIBAgIRAP6OtVIpSKXwu9dCxSQUBRcwCgYIKoZIzj0EAwIwKzEp\nMCcGA1UEAxMgSG90U3R1ZmYgU2VsZi1TaWduZWQgQ2VydGlmaWNhdGUwHhcNMjEw\nMTEzMTExMTMyWhcNMzEwMTEzMTExMTMyWjArMSkwJwYDVQQDEyBIb3RTdHVmZiBT\nZWxmLVNpZ25lZCBDZXJ0aWZpY2F0ZTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA\nBFYIkjm5uPYC78EG6q5zTuu56p0g+4n3mFFkdqN2fuXyzw9ZrBZ/z6LjAXytBKY6\n/zARYhaVQbTPA/g/0ST1SnajRjBEMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAK\nBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA8GA1UdEQQIMAaHBH8AAAEwCgYIKoZI\nzj0EAwIDSAAwRQIgccoJdDly+VUGqkDU7wLjTpwYZJtiwIH3nkRhoaWtcOUCIQDE\nJ78TmegrP+YshLGWWpifGE6lMsKnVNWrccBy6ZXjQA==\n-----END CERTIFICATE-----"
	certBlock4, _ := pem.Decode([]byte(certString4))
	certBytes4 := certBlock4.Bytes
	cert4, _ := tls.X509KeyPair(certBytes4, privkeyPEM4)

	var setup manualConfig

	setup.SelfID = 1
	setup.servers[1].ID = 1
	setup.servers[1].PeerAddr =
		setup.servers[1].ClientAddr
	setup.servers[1].Cert
	setup.servers[1].Pubkey

	privkey, err := data.ReadPrivateKeyFile(conf.Privkey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read private key file: %v\n", err)
		os.Exit(1)
	}

	var cert *tls.Certificate
	if conf.TLS {
		if conf.Cert == "" {
			for _, replica := range conf.Replicas {
				if replica.ID == conf.SelfID {
					conf.Cert = replica.Cert
				}
			}
		}
		certB, err := data.ReadCertFile(conf.Cert)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read certificate: %v\n", err)
			os.Exit(1)
		}
		// read in the private key again, but in PEM format
		pkPEM, err := ioutil.ReadFile(conf.Privkey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read private key: %v\n", err)
			os.Exit(1)
		}

		tlsCert, err := tls.X509KeyPair(certB, pkPEM)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse certificate: %v\n", err)
			os.Exit(1)
		}
		cert = &tlsCert
	}

	var clientAddress string

	replicaConfig := config.NewConfig(conf.SelfID, privkey, cert)
	replicaConfig.BatchSize = 4 //conf.BatchSize
	for _, r := range conf.Replicas {
		key, err := data.ReadPublicKeyFile(r.Pubkey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read public key file '%s': %v\n", r.Pubkey, err)
			os.Exit(1)
		}

		if conf.TLS {
			certB, err := data.ReadCertFile(r.Cert)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read certificate: %v\n", err)
				os.Exit(1)
			}
			if !replicaConfig.CertPool.AppendCertsFromPEM(certB) {
				fmt.Fprintf(os.Stderr, "Failed to parse certificate\n")
			}
		}

		// fmt.Println("Addr")
		// fmt.Println(r.PeerAddr)
		info := &config.ReplicaInfo{
			ID:      r.ID,
			Address: r.PeerAddr,
			PubKey:  key,
		}

		if r.ID == conf.SelfID {
			// override own addresses if set
			if conf.ClientAddr != "" {
				clientAddress = conf.ClientAddr
			} else {
				clientAddress = r.ClientAddr
			}
			if conf.PeerAddr != "" {
				info.Address = conf.PeerAddr
			}
		}

		replicaConfig.Replicas[r.ID] = info
	}
	replicaConfig.QuorumSize = len(replicaConfig.Replicas) - (len(replicaConfig.Replicas)-1)/3

	// fmt.Println(replicaConfig.Replicas[1].Address)
	srv := newHotStuffServer(&conf, replicaConfig)
	err = srv.Start(clientAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start HotStuff: %v\n", err)
		os.Exit(1)
	}

	<-signals
	fmt.Fprintf(os.Stderr, "Exiting...\n")
	srv.Stop()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

// cmdID is a unique identifier for a command
type cmdID struct {
	clientID    uint32
	sequenceNum uint64
}

type hotstuffServer struct {
	ctx       context.Context
	cancel    context.CancelFunc
	conf      *options
	gorumsSrv *client.GorumsServer
	hs        *hotstuff.HotStuff
	pm        interface {
		Run(context.Context)
	}

	mut          sync.Mutex
	finishedCmds map[cmdID]chan struct{}

	lastExecTime int64
}

func newHotStuffServer(conf *options, replicaConfig *config.ReplicaConfig) *hotstuffServer {
	ctx, cancel := context.WithCancel(context.Background())

	serverOpts := []client.ServerOption{}
	grpcServerOpts := []grpc.ServerOption{}

	if conf.TLS {
		grpcServerOpts = append(grpcServerOpts, grpc.Creds(credentials.NewServerTLSFromCert(replicaConfig.Cert)))
	}

	serverOpts = append(serverOpts, client.WithGRPCServerOptions(grpcServerOpts...))

	srv := &hotstuffServer{
		ctx:          ctx,
		cancel:       cancel,
		conf:         conf,
		gorumsSrv:    client.NewGorumsServer(serverOpts...),
		finishedCmds: make(map[cmdID]chan struct{}),
		lastExecTime: time.Now().UnixNano(),
	}
	var pm hotstuff.Pacemaker
	switch conf.PmType {
	case "fixed":
		pm = pacemaker.NewFixedLeader(conf.LeaderID)
	case "round-robin":
		pm = pacemaker.NewRoundRobin(
			conf.ViewChange, conf.Schedule, time.Duration(conf.ViewTimeout)*time.Millisecond,
		)
	default:
		fmt.Fprintf(os.Stderr, "Invalid pacemaker type: '%s'\n", conf.PmType)
		os.Exit(1)
	}
	srv.hs = hotstuff.New(replicaConfig, pm, conf.TLS, time.Minute, time.Duration(conf.ViewTimeout)*time.Millisecond)
	srv.pm = pm.(interface{ Run(context.Context) })
	// Use a custom server instead of the gorums one
	srv.gorumsSrv.RegisterClientServer(srv)
	return srv
}

func (srv *hotstuffServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	err = srv.hs.Start()
	if err != nil {
		return err
	}

	go srv.gorumsSrv.Serve(lis)
	go srv.pm.Run(srv.ctx)
	go srv.onExec()

	return nil
}

func (srv *hotstuffServer) Stop() {
	srv.gorumsSrv.Stop()
	srv.cancel()
	srv.hs.Close()
}

func (srv *hotstuffServer) ExecCommand(_ context.Context, cmd *client.Command, out func(*client.Empty, error)) {
	finished := make(chan struct{})
	id := cmdID{cmd.ClientID, cmd.SequenceNumber}
	srv.mut.Lock()
	srv.finishedCmds[id] = finished
	srv.mut.Unlock()
	// marshal the message back to a byte so that HotStuff can process it.
	// TODO: think of a better way to do this.
	b, err := proto.MarshalOptions{Deterministic: true}.Marshal(cmd)
	if err != nil {
		log.Fatalf("Failed to marshal command: %v", err)
		out(nil, status.Errorf(codes.InvalidArgument, "Failed to marshal command: %v", err))
	}
	srv.hs.AddCommand(data.Command(b))

	go func(id cmdID, finished chan struct{}) {
		<-finished

		srv.mut.Lock()
		delete(srv.finishedCmds, id)
		srv.mut.Unlock()

		// send response
		out(&client.Empty{}, nil)
	}(id, finished)
}

func (srv *hotstuffServer) onExec() {
	for cmds := range srv.hs.GetExec() {
		if len(cmds) > 0 && srv.conf.PrintThroughput {
			now := time.Now().UnixNano()
			prev := atomic.SwapInt64(&srv.lastExecTime, now)
			fmt.Printf("%d, %d\n", now-prev, len(cmds))
		}

		for _, cmd := range cmds {
			m := new(client.Command)
			err := proto.Unmarshal([]byte(cmd), m)
			if err != nil {
				log.Printf("Failed to unmarshal command: %v\n", err)
			}
			if srv.conf.PrintCommands {
				fmt.Printf("%s", m.Data)
			}
			srv.mut.Lock()
			if c, ok := srv.finishedCmds[cmdID{m.ClientID, m.SequenceNumber}]; ok {
				c <- struct{}{}
			}
			srv.mut.Unlock()
		}
	}
}
