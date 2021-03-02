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
}

var connections connMap

func main() {
	connections.connections = make(map[string]net.Conn)
	connections.answer = make(map[string]string)
	connections.offer = make(map[string]string)

	flag.Parse()
	log.Printf("listening on %q...", *listen)
	go func() {
		err := http.ListenAndServe("127.0.0.1:8080", http.FileServer(http.Dir(*dir)))
		log.Fatalln(err)
	}()

	go func() {
		l, err := net.Listen("tcp", "127.0.0.1:13371")
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
	}()

	go func() {
		l3, err3 := net.Listen("tcp", "127.0.0.1:13373")
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
		l4, err4 := net.Listen("tcp", "127.0.0.1:13374")
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
	opts := &websocket.AcceptOptions{OriginPatterns: []string{"*"}, Subprotocols: []string{"*"}}
	c, err := websocket.Accept(w, r, opts)
	if err != nil {
		s.logf("%v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	conn := websocket.NetConn(ctx, c, 1)
	for {
		msg, err := bufio.NewReader(conn).ReadString('&')
		fmt.Println(msg)
		fmt.Println("Request \"Addr\": ")
		fmt.Println(r.Host)
		msg = strings.Split(msg, "&")[0]
		msgType := strings.Split(strings.Split(msg, "setup:")[1], "\n")[0]
		msgType = strings.TrimSpace(msgType)
		fmt.Println(msgType)

		// connections.connections[r.Host] = conn

		if msgType == "actpass" {
			connections.offer[r.Host] = msg
		} else if msgType == "active" {
			connections.answer[r.Host] = msg
		} else if msgType == "recvOffer" {
			for {
				time.Sleep(time.Millisecond * 2000)
				fmt.Println("inn loop")
				// fmt.Println(connections.offer)
				connections.mux.Lock()
				if _, ok := connections.offer[r.Host]; ok {
					fmt.Println("ok")
					// _, err = bufio.NewWriter(conn).WriteString(connections.offer["localhost:13371"])
					conn.Write([]byte(connections.offer[r.Host]))
					if err != nil {
						fmt.Println(err)
					}
					// conn.Write([]byte(connections.offer["localhost:13371"]))
					connections.mux.Unlock()
					fmt.Println("breaking")
					break
				}
				connections.mux.Unlock()
			}
		} else if msgType == "recvAnswer" {
			for {
				time.Sleep(time.Millisecond * 500)
				fmt.Println("inn loop")
				// fmt.Println(connections.offer)
				connections.mux.Lock()
				if _, ok := connections.answer[r.Host]; ok {
					fmt.Println("ok")
					// _, err = bufio.NewWriter(conn).WriteString(connections.offer["localhost:13371"])
					conn.Write([]byte(connections.answer[r.Host]))
					if err != nil {
						fmt.Println(err)
					}
					// conn.Write([]byte(connections.offer["localhost:13371"]))
					connections.mux.Unlock()
					fmt.Println("breaking")
					break
				}
				connections.mux.Unlock()
			}
			break
		}
	}

	// for {
	// 	time.Sleep(time.Second * 1)
	// 	buff := make([]byte, 4096)
	// 	n, err := conn.Read(buff)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		break
	// 	}
	// 	res := make([]byte, n)
	// 	copy(res, buff[:n])
	// 	fmt.Println(res)
	// }

	defer c.Close(websocket.StatusInternalError, "the sky is falling")
	defer cancel()
	fmt.Println("Accepted")
}
