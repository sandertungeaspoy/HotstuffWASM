package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

// echoServer is the WebSocket echo server implementation.
// It ensures the client speaks the echo subprotocol and
// only allows one message every 100ms with a 10 message burst.
// wasmServer is the WebSocket wasm server implementation.
type wasmServer struct {
	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
}

var (
	listen = flag.String("listen", "127.0.0.1:8080", "listen address")
	dir    = flag.String("dir", ".", "directory to serve")
)

// var connections map[string][]net.Conn
type connMap struct {
	mux         sync.Mutex
	connections map[string]net.Conn
	answer      map[string]string
	offer       map[string]string
	completed   map[string]bool
}

type system struct {
	inUseIDs []int
	connMaps map[int]connMap
}

var connections connMap
var systemMap system

func main() {
	connections.connections = make(map[string]net.Conn)
	connections.answer = make(map[string]string)
	connections.offer = make(map[string]string)
	connections.completed = make(map[string]bool)
	systemMap.inUseIDs = make([]int, 0)
	systemMap.connMaps = make(map[int]connMap)

	flag.Parse()
	log.Printf("listening on %q...", *listen)
	go func() {
		err := http.ListenAndServe("127.0.0.1:8080", http.FileServer(http.Dir(*dir)))
		log.Fatalln(err)
	}()

	l2, err2 := net.Listen("tcp", "127.0.0.1:13372")
	if err2 != nil {
		fmt.Println(err2)
	}
	log.Printf("listening on http://%v", l2.Addr())

	server2 := &http.Server{
		Handler: wasmServer{
			logf: log.Printf,
		},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc2 := make(chan error, 1)
	go func() {
		// fmt.Println("Failed At serve")
		errc2 <- server2.Serve(l2)
	}()

	sigs2 := make(chan os.Signal, 1)
	signal.Notify(sigs2, os.Interrupt)
	select {
	case err := <-errc2:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs2:
		log.Printf("terminating: %v", sig)
	}

}

func (s wasmServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Failed at ServeHTTP")
	// if _, ok := connections.connections[r.Host]; ok {
	// 	return
	// }
	fmt.Println(r.RemoteAddr)
	opts := &websocket.AcceptOptions{OriginPatterns: []string{"*"}, Subprotocols: []string{"*"}}
	c, err := websocket.Accept(w, r, opts)
	if err != nil {
		s.logf("%v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer c.Close(websocket.StatusInternalError, "WebSocket has been closed")
	defer cancel()
	conn := websocket.NetConn(ctx, c, 1)

	msg, err := bufio.NewReader(conn).ReadString('%')
	fmt.Println(msg)
	fmt.Println("Request \"Addr\": ")
	fmt.Println(r.Host)
	if strings.TrimSpace(msg) == "" {
		fmt.Println(err)
		return
	}
	msgs := strings.Split(msg, "&")
	senderID := strings.Split(msgs[1], "%")
	msgType := strings.Split(strings.Split(msgs[0], "setup:")[1], "\n")[0]
	msgType = strings.TrimSpace(msgType)
	fmt.Println(msgType)

	// connections.connections[r.Host] = conn

	if msgType == "actpass" {
		connections.mux.Lock()
		connections.offer[senderID[0]] = msgs[0]
		connections.completed[senderID[0]] = false
		connections.mux.Unlock()
	} else if msgType == "active" {
		connections.mux.Lock()
		connections.answer[senderID[0]] = msgs[0]
		connections.completed[senderID[0]] = true
		connections.mux.Unlock()
	} else if msgType == "recvOffer" {
		connections.mux.Lock()
		fmt.Println(connections.offer)
		for key, value := range connections.completed {
			if !value {
				conn.Write([]byte(connections.offer[key] + "&" + key + "%"))
				connections.mux.Unlock()
				return
			}
		}
		connections.mux.Unlock()
	} else if msgType == "recvAnswer" {
		connections.mux.Lock()
		fmt.Println(connections.answer)
		if connections.completed[senderID[0]] {
			conn.Write([]byte(connections.answer[senderID[0]] + "&" + senderID[0] + "%"))
			delete(connections.answer, senderID[0])
			connections.mux.Unlock()
			return
		}
		connections.mux.Unlock()
	} else if msgType == "removeOffer" {
		connections.mux.Lock()
		_, ok := connections.offer[senderID[0]]
		if ok {
			delete(connections.offer, senderID[0])
			delete(connections.completed, senderID[0])
			connections.mux.Unlock()
			return
		}
		connections.mux.Unlock()
	} else if msgType == "removeAnswer" {
		connections.mux.Lock()
		_, ok := connections.answer[senderID[0]]
		if ok {
			delete(connections.answer, senderID[0])
			delete(connections.completed, senderID[0])
			connections.mux.Unlock()
			return
		}
		connections.mux.Unlock()
	} else if msgType == "purgeDatabase" {
		connections.mux.Lock()
		connections.connections = make(map[string]net.Conn)
		connections.answer = make(map[string]string)
		connections.offer = make(map[string]string)
		connections.completed = make(map[string]bool)
		connections.mux.Unlock()
		return
	}

	fmt.Println("Accepted")
}
