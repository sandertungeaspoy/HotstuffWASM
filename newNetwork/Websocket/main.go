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
	listen = flag.String("listen", "localhost:8080", "listen address")
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

var connections connMap

func main() {
	connections.connections = make(map[string]net.Conn)
	connections.answer = make(map[string]string)
	connections.offer = make(map[string]string)
	connections.completed = make(map[string]bool)

	flag.Parse()
	log.Printf("listening on %q...", *listen)
	go func() {
		err := http.ListenAndServe("localhost:8080", http.FileServer(http.Dir(*dir)))
		log.Fatalln(err)
	}()

	go func() {
		l, err := net.Listen("tcp", "localhost:13371")
		if err != nil {
			fmt.Println(err)
		}
		log.Printf("listening on http://%v", l.Addr())

		server1 := &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc <- server1.Serve(l)
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt)
		select {
		case err := <-errc:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		l2, err2 := net.Listen("tcp", "localhost:13372")
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
	}()

	go func() {
		l3, err3 := net.Listen("tcp", "localhost:13373")
		if err3 != nil {
			fmt.Println(err3)
		}
		log.Printf("listening on http://%v", l3.Addr())

		server3 := &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc3 := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc3 <- server3.Serve(l3)
		}()

		sigs3 := make(chan os.Signal, 1)
		signal.Notify(sigs3, os.Interrupt)
		select {
		case err := <-errc3:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs3:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		l4, err4 := net.Listen("tcp", "localhost:13374")
		if err4 != nil {
			fmt.Println(err4)
		}
		log.Printf("listening on http://%v", l4.Addr())

		server4 := &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc4 := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc4 <- server4.Serve(l4)
		}()

		sigs4 := make(chan os.Signal, 1)
		signal.Notify(sigs4, os.Interrupt)
		select {
		case err := <-errc4:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs4:
			log.Printf("terminating: %v", sig)
		}
	}()

	for {
		sigsExit := make(chan os.Signal, 1)
		signal.Notify(sigsExit, os.Interrupt)
		select {
		case <-sigsExit:
			fmt.Println("Break loop close")
			return
		}

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
	defer c.Close(websocket.StatusInternalError, "the sky is falling")
	defer cancel()
	conn := websocket.NetConn(ctx, c, 1)
	if strings.Split(r.Host, ":")[1] == "13371" {
		connections.connections = make(map[string]net.Conn)
		connections.answer = make(map[string]string)
		connections.offer = make(map[string]string)
		fmt.Println("Connection Map reset!")
		conn.Close()
		cancel()
		return
	}
	msg, err := bufio.NewReader(conn).ReadString('%')
	fmt.Println(msg)
	fmt.Println("Request \"Addr\": ")
	fmt.Println(r.Host)
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
	}

	fmt.Println("Accepted")
}
