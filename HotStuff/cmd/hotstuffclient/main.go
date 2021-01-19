package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/relab/gorums/benchmark"
	"github.com/relab/hotstuff/client"
	"github.com/relab/hotstuff/config"
	"github.com/spf13/pflag"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type options struct {
	SelfID      config.ReplicaID `mapstructure:"self-id"`
	RateLimit   int              `mapstructure:"rate-limit"`
	PayloadSize int              `mapstructure:"payload-size"`
	MaxInflight uint64           `mapstructure:"max-inflight"`
	DataSource  string           `mapstructure:"input"`
	Benchmark   bool             `mapstructure:"benchmark"`
	ExitAfter   int              `mapstructure:"exit-after"`
	TLS         bool
	Replicas    []struct {
		ID         config.ReplicaID
		ClientAddr string `mapstructure:"client-address"`
		Pubkey     string
		Cert       string
	}
}

type manualConfig struct {
	SelfID   config.ReplicaID
	Replicas map[uint32]*Replicas
}

// Replicas struct
type Replicas struct {
	ID         config.ReplicaID
	ClientAddr string
	PubKey     *ecdsa.PublicKey
	Cert       *tls.Certificate
	CertPEM    []byte
}

func usage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println()
	fmt.Println("Loads configuration from ./hotstuff.toml")
	fmt.Println()
	fmt.Println("Options:")
	pflag.PrintDefaults()
}

func main() {
	// pflag.Usage = usage

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("This is a test")

	// help := pflag.BoolP("help", "h", false, "Prints this text.")
	// cpuprofile := pflag.String("cpuprofile", "", "File to write CPU profile to")
	// memprofile := pflag.String("memprofile", "", "File to write memory profile to")
	// pflag.Uint32("self-id", 0, "The id for this replica.")
	// pflag.Int("rate-limit", 0, "Limit the request-rate to approximately (in requests per second).")
	// pflag.Int("payload-size", 0, "The size of the payload in bytes")
	// pflag.Uint64("max-inflight", 10000, "The maximum number of messages that the client can wait for at once")
	// pflag.String("input", "", "Optional file to use for payload data")
	// pflag.Bool("benchmark", false, "If enabled, a BenchmarkData protobuf will be written to stdout.")
	// pflag.Int("exit-after", 0, "Number of seconds after which the program should exit.")
	// pflag.Bool("tls", false, "Enable TLS")
	// pflag.Parse()

	// if *help {
	// 	pflag.Usage()
	// 	os.Exit(0)
	// }

	// if *cpuprofile != "" {
	// 	f, err := os.Create(*cpuprofile)
	// 	if err != nil {
	// 		log.Fatal("Could not create CPU profile: ", err)
	// 	}
	// 	defer f.Close()
	// 	if err := pprof.StartCPUProfile(f); err != nil {
	// 		log.Fatal("Could not start CPU profile: ", err)
	// 	}
	// 	defer pprof.StopCPUProfile()
	// }

	// viper.BindPFlags(pflag.CommandLine)

	// read main config file in working dir
	// viper.SetConfigName("hotstuff")
	// viper.AddConfigPath(".")
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
	// 	os.Exit(1)
	// }

	var conf options
	// err = viper.Unmarshal(&conf)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to unmarshal config: %v\n", err)
	// 	os.Exit(1)
	// }

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
	// privBlock1, _ := pem.Decode([]byte(privkeyString1))
	// privkeyBytes1 := privBlock1.Bytes
	privkeyPEM1 := []byte(privkeyString1)
	// privateKey1, _ := x509.ParseECPrivateKey(privkeyBytes1)

	privkeyString2 := "-----BEGIN HOTSTUFF PRIVATE KEY-----\nMHcCAQEEIDZbG1HNK/NejlozKQHeLYFTFPPi0QYFHNRP/OwlVLB4oAoGCCqGSM49\nAwEHoUQDQgAEitP7/gomqGK/TSgALUpy+MO9N/n1vzyHYXvdwPRFPOyS79UEJIYf\nNCRyex+TRmtB+jwwo1A+x7hCdk2azaF7FA==\n-----END HOTSTUFF PRIVATE KEY-----"
	// privBlock2, _ := pem.Decode([]byte(privkeyString2))
	// privkeyBytes2 := privBlock2.Bytes
	privkeyPEM2 := []byte(privkeyString2)
	// privateKey2, _ := x509.ParseECPrivateKey(privkeyBytes2)

	privkeyString3 := "-----BEGIN HOTSTUFF PRIVATE KEY-----\nMHcCAQEEIKubRjLNMfX+L4dnDSPcyQIZz/DBdPOyURXUyMibFr/LoAoGCCqGSM49\nAwEHoUQDQgAE6RscDCsZrjOUnRuoUrONyPckRVoot+oGPFjNBynLAWtT07yBCPWY\nUwzM7Zn+IM3KyAxN12UVZd4itCkGfmOlAg==\n-----END HOTSTUFF PRIVATE KEY-----"
	// privBlock3, _ := pem.Decode([]byte(privkeyString3))
	// privkeyBytes3 := privBlock3.Bytes
	privkeyPEM3 := []byte(privkeyString3)
	// privateKey3, _ := x509.ParseECPrivateKey(privkeyBytes3)

	privkeyString4 := "-----BEGIN HOTSTUFF PRIVATE KEY-----\nMHcCAQEEIAE318VoY//HCbPSSzeYv69esRxZcoIpu10YYq1+h+FVoAoGCCqGSM49\nAwEHoUQDQgAEVgiSObm49gLvwQbqrnNO67nqnSD7ifeYUWR2o3Z+5fLPD1msFn/P\nouMBfK0Epjr/MBFiFpVBtM8D+D/RJPVKdg==\n-----END HOTSTUFF PRIVATE KEY-----"
	// privBlock4, _ := pem.Decode([]byte(privkeyString4))
	// privkeyBytes4 := privBlock4.Bytes
	privkeyPEM4 := []byte(privkeyString4)
	// privateKey4, _ := x509.ParseECPrivateKey(privkeyBytes4)

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

	var setup manualConfig

	setup.SelfID = conf.SelfID
	setup.Replicas = make(map[uint32]*Replicas)
	setup.Replicas[1] = &Replicas{}
	setup.Replicas[1].ID = 1
	setup.Replicas[1].ClientAddr = "127.0.0.1:23371"
	setup.Replicas[1].Cert = &cert1
	setup.Replicas[1].PubKey = publicKey1
	setup.Replicas[1].CertPEM = certPEM1

	setup.Replicas[2] = &Replicas{}
	setup.Replicas[2].ID = 2
	setup.Replicas[2].ClientAddr = "127.0.0.1:23372"
	setup.Replicas[2].Cert = &cert2
	setup.Replicas[2].PubKey = publicKey2
	setup.Replicas[2].CertPEM = certPEM2

	setup.Replicas[3] = &Replicas{}
	setup.Replicas[3].ID = 3
	setup.Replicas[3].ClientAddr = "127.0.0.1:23373"
	setup.Replicas[3].Cert = &cert3
	setup.Replicas[3].PubKey = publicKey3
	setup.Replicas[3].CertPEM = certPEM3

	setup.Replicas[4] = &Replicas{}
	setup.Replicas[4].ID = 4
	setup.Replicas[4].ClientAddr = "127.0.0.1:23374"
	setup.Replicas[4].Cert = &cert4
	setup.Replicas[4].PubKey = publicKey4
	setup.Replicas[4].CertPEM = certPEM4

	replicaConfig := config.NewConfig(0, nil, nil)
	for _, r := range setup.Replicas {
		if conf.TLS {
			if !replicaConfig.CertPool.AppendCertsFromPEM(r.CertPEM) {
				fmt.Fprintf(os.Stderr, "Failed to parse certificate\n")
			}
		}
		info := &config.ReplicaInfo{
			ID:      r.ID,
			Address: r.ClientAddr,
			PubKey:  r.PubKey,
		}

		replicaConfig.Replicas[config.ReplicaID(r.ID)] = info
	}

	client, err := newHotStuffClient(&conf, replicaConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start client: %v\n", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if conf.ExitAfter > 0 {
			select {
			case <-time.After(time.Duration(conf.ExitAfter) * time.Second):
			case <-signals:
				fmt.Fprintf(os.Stderr, "Exiting...\n")
			}
		} else {
			<-signals
			fmt.Fprintf(os.Stderr, "Exiting...\n")
		}
		cancel()
	}()

	err = client.SendCommands(ctx)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "Failed to send commands: %v\n", err)
		client.Close()
		os.Exit(1)
	}
	client.Close()

	stats := client.GetStats()
	throughput := stats.Throughput
	latency := stats.LatencyAvg / float64(time.Millisecond)
	latencySD := math.Sqrt(stats.LatencyVar) / float64(time.Millisecond)

	if !conf.Benchmark {
		fmt.Printf("Throughput (ops/sec): %.2f, Latency (ms): %.2f, Latency Std.dev (ms): %.2f\n",
			throughput,
			latency,
			latencySD,
		)
	} else {
		client.data.MeasuredThroughput = throughput
		client.data.MeasuredLatency = latency
		client.data.LatencyVariance = math.Pow(latencySD, 2) // variance in ms^2
		b, err := proto.Marshal(client.data)
		if err != nil {
			log.Fatalf("Could not marshal benchmarkdata: %v\n", err)
		}
		_, err = os.Stdout.Write(b)
		if err != nil {
			log.Fatalf("Could not write data: %v\n", err)
		}
	}

	// if *memprofile != "" {
	// 	f, err := os.Create(*memprofile)
	// 	if err != nil {
	// 		log.Fatal("could not create memory profile: ", err)
	// 	}
	// 	defer f.Close() // error handling omitted for example
	// 	runtime.GC()    // get up-to-date statistics
	// 	if err := pprof.WriteHeapProfile(f); err != nil {
	// 		log.Fatal("could not write memory profile: ", err)
	// 	}
	// }
}

type qspec struct {
	faulty int
}

func (q *qspec) ExecCommandQF(_ *client.Command, signatures map[uint32]*client.Empty) (*client.Empty, bool) {
	if len(signatures) < q.faulty+1 {
		return nil, false
	}
	return &client.Empty{}, true
}

type hotstuffClient struct {
	inflight      uint64
	reader        io.ReadCloser
	conf          *options
	mgr           *client.Manager
	replicaConfig *config.ReplicaConfig
	gorumsConfig  *client.Configuration
	wg            sync.WaitGroup
	stats         benchmark.Stats       // records latency and throughput
	data          *client.BenchmarkData // stores time and duration for each command
}

func newHotStuffClient(conf *options, replicaConfig *config.ReplicaConfig) (*hotstuffClient, error) {
	nodes := make(map[string]uint32, len(replicaConfig.Replicas))
	for _, r := range replicaConfig.Replicas {
		nodes[r.Address] = uint32(r.ID)
	}

	grpcOpts := []grpc.DialOption{grpc.WithBlock()}

	if conf.TLS {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(replicaConfig.CertPool, "")))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}

	mgr, err := client.NewManager(client.WithNodeMap(nodes), client.WithGrpcDialOptions(grpcOpts...),
		client.WithDialTimeout(time.Minute),
	)
	if err != nil {
		return nil, err
	}
	faulty := (len(replicaConfig.Replicas) - 1) / 3
	gorumsConf, err := mgr.NewConfiguration(mgr.NodeIDs(), &qspec{faulty: faulty})
	if err != nil {
		mgr.Close()
		return nil, err
	}
	var reader io.ReadCloser
	if conf.DataSource == "" {
		reader = ioutil.NopCloser(rand.Reader)
	} else {
		f, err := os.Open(conf.DataSource)
		if err != nil {
			mgr.Close()
			return nil, err
		}
		reader = f
	}
	return &hotstuffClient{
		reader:        reader,
		conf:          conf,
		mgr:           mgr,
		replicaConfig: replicaConfig,
		gorumsConfig:  gorumsConf,
		data:          &client.BenchmarkData{},
	}, nil
}

func (c *hotstuffClient) Close() {
	c.mgr.Close()
	c.reader.Close()
}

func (c *hotstuffClient) GetStats() *benchmark.Result {
	return c.stats.GetResult()
}

func (c *hotstuffClient) SendCommands(ctx context.Context) error {
	var num uint64
	var sleeptime time.Duration
	if c.conf.RateLimit > 0 {
		sleeptime = time.Second / time.Duration(c.conf.RateLimit)
	}

	defer c.stats.End()
	defer c.wg.Wait()
	c.stats.Start()

	for {
		if atomic.LoadUint64(&c.inflight) < c.conf.MaxInflight {
			atomic.AddUint64(&c.inflight, 1)
			data := make([]byte, c.conf.PayloadSize)
			n, err := c.reader.Read(data)
			if err != nil {
				return err
			}
			cmd := &client.Command{
				ClientID:       uint32(c.conf.SelfID),
				SequenceNumber: num,
				Data:           data[:n],
			}
			now := time.Now()
			promise := c.gorumsConfig.ExecCommand(ctx, cmd)
			num++

			c.wg.Add(1)
			go func(promise *client.FutureEmpty, sendTime time.Time) {
				_, err := promise.Get()
				atomic.AddUint64(&c.inflight, ^uint64(0))
				if err != nil {
					qcError, ok := err.(client.QuorumCallError)
					if !ok || qcError.Reason != context.Canceled.Error() {
						log.Printf("Did not get enough signatures for command: %v\n", err)
					}
				}
				duration := time.Since(sendTime)
				c.stats.AddLatency(duration)
				if c.conf.Benchmark {
					c.data.Stats = append(c.data.Stats, &client.CommandStats{
						StartTime: timestamppb.New(sendTime),
						Duration:  durationpb.New(duration),
					})
				}
				c.wg.Done()
			}(promise, now)
		}

		if c.conf.RateLimit > 0 {
			time.Sleep(sleeptime)
		}

		err := ctx.Err()
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
