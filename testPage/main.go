package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

// wasmServer is the WebSocket wasm server implementation.
type wasmServer struct {
	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
}

var (
	listen = flag.String("listen", "127.0.0.1:8080", "listen address")
	dir    = flag.String("dir", ".", "directory to serve")
)

var servers []*http.Server

// var connections map[string][]net.Conn
type connMap struct {
	mux         sync.Mutex
	connections map[string]net.Conn
}

var connections connMap

func main() {
	servers = make([]*http.Server, 16)
	// connections = make(map[string][]net.Conn)
	connections.connections = make(map[string]net.Conn)
	flag.Parse()
	log.Printf("listening on %q...", *listen)
	go func() {
		err := http.ListenAndServe("127.0.0.1:8080", http.FileServer(http.Dir(*dir)))
		log.Fatalln(err)
	}()

	go func() {
		l, err := net.Listen("tcp", "127.0.0.1:13711")
		if err != nil {
			fmt.Println(err)
		}
		log.Printf("listening on http://%v", l.Addr())

		servers[0] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc <- servers[0].Serve(l)
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
		l2, err2 := net.Listen("tcp", "127.0.0.1:13721")
		if err2 != nil {
			fmt.Println(err2)
		}
		log.Printf("listening on http://%v", l2.Addr())

		servers[1] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc2 := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc2 <- servers[1].Serve(l2)
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
		l3, err3 := net.Listen("tcp", "127.0.0.1:13731")
		if err3 != nil {
			fmt.Println(err3)
		}
		log.Printf("listening on http://%v", l3.Addr())

		servers[2] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc3 := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc3 <- servers[2].Serve(l3)
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
		l4, err4 := net.Listen("tcp", "127.0.0.1:13741")
		if err4 != nil {
			fmt.Println(err4)
		}
		log.Printf("listening on http://%v", l4.Addr())

		servers[3] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc4 := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc4 <- servers[3].Serve(l4)
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

	go func() {
		l5, err5 := net.Listen("tcp", "127.0.0.1:23711")
		if err5 != nil {
			fmt.Println(err5)
		}
		log.Printf("listening on http://%v", l5.Addr())

		servers[4] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc5 := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc5 <- servers[4].Serve(l5)
		}()

		sigs5 := make(chan os.Signal, 1)
		signal.Notify(sigs5, os.Interrupt)
		select {
		case err := <-errc5:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs5:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		l6, err6 := net.Listen("tcp", "127.0.0.1:23721")
		if err6 != nil {
			fmt.Println(err6)
		}
		log.Printf("listening on http://%v", l6.Addr())

		servers[5] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc6 := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc6 <- servers[5].Serve(l6)
		}()

		sigs6 := make(chan os.Signal, 1)
		signal.Notify(sigs6, os.Interrupt)
		select {
		case err := <-errc6:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs6:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		l7, err7 := net.Listen("tcp", "127.0.0.1:23731")
		if err7 != nil {
			fmt.Println(err7)
		}
		log.Printf("listening on http://%v", l7.Addr())

		servers[6] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc7 := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc7 <- servers[6].Serve(l7)
		}()

		sigs7 := make(chan os.Signal, 1)
		signal.Notify(sigs7, os.Interrupt)
		select {
		case err := <-errc7:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs7:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		l8, err8 := net.Listen("tcp", "127.0.0.1:23741")
		if err8 != nil {
			fmt.Println(err8)
		}
		log.Printf("listening on http://%v", l8.Addr())

		servers[7] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc8 := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc8 <- servers[7].Serve(l8)
		}()

		sigs8 := make(chan os.Signal, 1)
		signal.Notify(sigs8, os.Interrupt)
		select {
		case err := <-errc8:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs8:
			log.Printf("terminating: %v", sig)
		}
	}()
	// Client ADDR

	go func() {
		lc, errcl := net.Listen("tcp", "127.0.0.1:13712")
		if errcl != nil {
			fmt.Println(errcl)
		}
		log.Printf("listening on http://%v", lc.Addr())

		servers[8] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errccl := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errccl <- servers[8].Serve(lc)
		}()

		sigsc := make(chan os.Signal, 1)
		signal.Notify(sigsc, os.Interrupt)
		select {
		case err := <-errccl:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigsc:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		l2c, err2c := net.Listen("tcp", "127.0.0.1:13722")
		if err2c != nil {
			fmt.Println(err2c)
		}
		log.Printf("listening on http://%v", l2c.Addr())

		servers[9] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc2c := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc2c <- servers[9].Serve(l2c)
		}()

		sigs2c := make(chan os.Signal, 1)
		signal.Notify(sigs2c, os.Interrupt)
		select {
		case err := <-errc2c:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs2c:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		l3c, err3c := net.Listen("tcp", "127.0.0.1:13732")
		if err3c != nil {
			fmt.Println(err3c)
		}
		log.Printf("listening on http://%v", l3c.Addr())

		servers[10] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc3c := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc3c <- servers[10].Serve(l3c)
		}()

		sigs3c := make(chan os.Signal, 1)
		signal.Notify(sigs3c, os.Interrupt)
		select {
		case err := <-errc3c:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs3c:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		l4c, err4c := net.Listen("tcp", "127.0.0.1:13742")
		if err4c != nil {
			fmt.Println(err4c)
		}
		log.Printf("listening on http://%v", l4c.Addr())

		servers[11] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc4c := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc4c <- servers[11].Serve(l4c)
		}()

		sigs4c := make(chan os.Signal, 1)
		signal.Notify(sigs4c, os.Interrupt)
		select {
		case err := <-errc4c:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs4c:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		l5c, err5c := net.Listen("tcp", "127.0.0.1:23712")
		if err5c != nil {
			fmt.Println(err5c)
		}
		log.Printf("listening on http://%v", l5c.Addr())

		servers[12] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc5c := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc5c <- servers[12].Serve(l5c)
		}()

		sigs5c := make(chan os.Signal, 1)
		signal.Notify(sigs5c, os.Interrupt)
		select {
		case err := <-errc5c:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs5c:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		l6c, err6c := net.Listen("tcp", "127.0.0.1:23722")
		if err6c != nil {
			fmt.Println(err6c)
		}
		log.Printf("listening on http://%v", l6c.Addr())

		servers[13] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc6c := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc6c <- servers[13].Serve(l6c)
		}()

		sigs6c := make(chan os.Signal, 1)
		signal.Notify(sigs6c, os.Interrupt)
		select {
		case err := <-errc6c:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs6c:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		l7c, err7c := net.Listen("tcp", "127.0.0.1:23732")
		if err7c != nil {
			fmt.Println(err7c)
		}
		log.Printf("listening on http://%v", l7c.Addr())

		servers[14] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc7c := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc7c <- servers[14].Serve(l7c)
		}()

		sigs7c := make(chan os.Signal, 1)
		signal.Notify(sigs7c, os.Interrupt)
		select {
		case err := <-errc7c:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs7c:
			log.Printf("terminating: %v", sig)
		}
	}()
	go func() {
		l8c, err8c := net.Listen("tcp", "127.0.0.1:23742")
		if err8c != nil {
			fmt.Println(err8c)
		}
		log.Printf("listening on http://%v", l8c.Addr())

		servers[15] = &http.Server{
			Handler: wasmServer{
				logf: log.Printf,
			},
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		errc8c := make(chan error, 1)
		go func() {
			// fmt.Println("Failed At serve")
			errc8c <- servers[15].Serve(l8c)
		}()

		sigs8c := make(chan os.Signal, 1)
		signal.Notify(sigs8c, os.Interrupt)
		select {
		case err := <-errc8c:
			log.Printf("failed to serve: %v", err)
		case sig := <-sigs8c:
			log.Printf("terminating: %v", sig)
		}
	}()

	go func() {
		for {
			connections.mux.Lock()
			if _, ok := connections.connections["127.0.0.1:13721"]; ok {
				conn := connections.connections["127.0.0.1:13721"]

				if _, ok := connections.connections["127.0.0.1:13722"]; ok {
					conn1 := connections.connections["127.0.0.1:13722"]
					connections.mux.Unlock()
					io.Copy(conn1, conn)
				} else {
					connections.mux.Unlock()
				}
			} else {
				connections.mux.Unlock()
			}
		}

	}()

	// go func() {
	// 	for {
	// 		// sigsExit := make(chan os.Signal, 1)
	// 		// signal.Notify(sigsExit, os.Interrupt)
	// 		// select {
	// 		// case <-sigsExit:
	// 		// 	break
	// 		// }
	// 		// fmt.Println(connections)
	// 		connections.mux.Lock()
	// 		if _, ok := connections.connections["127.0.0.1:13721"]; ok {
	// 			conn := connections.connections["127.0.0.1:13721"]
	// 			conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
	// 			connections.mux.Unlock()
	// 			buff := make([]byte, 4096)
	// 			n, err := conn.Read(buff)
	// 			res := make([]byte, n)
	// 			copy(res, buff[:n])
	// 			// msg, err := bufio.NewReader(conn).ReadString('\n')
	// 			if n > 0 {
	// 				fmt.Println(res)
	// 			}
	// 			if err != nil {
	// 				// fmt.Println(err)
	// 			} else {
	// 				for {
	// 					connections.mux.Lock()
	// 					if _, ok := connections.connections["127.0.0.1:13722"]; ok {
	// 						conn1 := connections.connections["127.0.0.1:13722"]
	// 						conn1.Write(res)
	// 						connections.mux.Unlock()
	// 						break
	// 						// fmt.Fprintf(connections.connections["127.0.0.1:13721"], msg)
	// 					}
	// 					connections.mux.Unlock()
	// 					// if _, ok := connections.connections["127.0.0.1:13731"]; ok {
	// 					// 	conn2 := connections.connections["127.0.0.1:13731"]
	// 					// 	conn2.Write(res)
	// 					// 	// fmt.Fprintf(connections.connections["127.0.0.1:13731"], msg)
	// 					// }
	// 					// if _, ok := connections.connections["127.0.0.1:13741"]; ok {
	// 					// 	conn3 := connections.connections["127.0.0.1:13741"]
	// 					// 	conn3.Write(res)
	// 					// 	// fmt.Fprintf(connections.connections["127.0.0.1:13741"], msg)
	// 					// }

	// 				}

	// 			}
	// 		} else {
	// 			connections.mux.Unlock()
	// 		}
	// 	}

	// }()

	go func() {
		for {
			// sigsExit := make(chan os.Signal, 1)
			// signal.Notify(sigsExit, os.Interrupt)
			// select {
			// case <-sigsExit:
			// 	break
			// }
			// fmt.Println(connections)
			connections.mux.Lock()
			if _, ok := connections.connections["127.0.0.1:13722"]; ok {
				conn := connections.connections["127.0.0.1:13722"]
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
				connections.mux.Unlock()
				buff := make([]byte, 4096)
				n, err := conn.Read(buff)
				res := make([]byte, n)
				copy(res, buff[:n])
				// msg, err := bufio.NewReader(conn).ReadString('\n')
				if n > 0 {
					fmt.Println(res)
				}
				if err != nil {
					// fmt.Println(err)
				} else {
					for {
						connections.mux.Lock()
						if _, ok := connections.connections["127.0.0.1:13721"]; ok {
							conn1 := connections.connections["127.0.0.1:13721"]
							conn1.Write(res)
							connections.mux.Unlock()
							break
							// fmt.Fprintf(connections.connections["127.0.0.1:13721"], msg)
						}
						connections.mux.Unlock()
						// if _, ok := connections.connections["127.0.0.1:13731"]; ok {
						// 	conn2 := connections.connections["127.0.0.1:13731"]
						// 	conn2.Write(res)
						// 	// fmt.Fprintf(connections.connections["127.0.0.1:13731"], msg)
						// }
						// if _, ok := connections.connections["127.0.0.1:13741"]; ok {
						// 	conn3 := connections.connections["127.0.0.1:13741"]
						// 	conn3.Write(res)
						// 	// fmt.Fprintf(connections.connections["127.0.0.1:13741"], msg)
						// }

					}

				}
			} else {
				connections.mux.Unlock()
			}
		}

	}()

	go func() {
		for {
			// sigsExit := make(chan os.Signal, 1)
			// signal.Notify(sigsExit, os.Interrupt)
			// select {
			// case <-sigsExit:
			// 	break
			// }
			// fmt.Println(connections)
			connections.mux.Lock()
			if _, ok := connections.connections["127.0.0.1:13711"]; ok {
				conn := connections.connections["127.0.0.1:13711"]
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
				connections.mux.Unlock()
				buff := make([]byte, 4096)
				n, err := conn.Read(buff)
				res := make([]byte, n)
				copy(res, buff[:n])
				// msg, err := bufio.NewReader(conn).ReadString('\n')
				if n > 0 {
					fmt.Println(res)
				}
				if err != nil {
					// fmt.Println(err)
				} else {
					for {
						connections.mux.Lock()
						if _, ok := connections.connections["127.0.0.1:13712"]; ok {
							conn1 := connections.connections["127.0.0.1:13712"]
							conn1.Write(res)
							connections.mux.Unlock()
							break
							// fmt.Fprintf(connections.connections["127.0.0.1:13721"], msg)
						}
						connections.mux.Unlock()
						// if _, ok := connections.connections["127.0.0.1:13731"]; ok {
						// 	conn2 := connections.connections["127.0.0.1:13731"]
						// 	conn2.Write(res)
						// 	// fmt.Fprintf(connections.connections["127.0.0.1:13731"], msg)
						// }
						// if _, ok := connections.connections["127.0.0.1:13741"]; ok {
						// 	conn3 := connections.connections["127.0.0.1:13741"]
						// 	conn3.Write(res)
						// 	// fmt.Fprintf(connections.connections["127.0.0.1:13741"], msg)
						// }

					}

				}
			} else {
				connections.mux.Unlock()
			}
		}

	}()

	go func() {
		for {
			// sigsExit := make(chan os.Signal, 1)
			// signal.Notify(sigsExit, os.Interrupt)
			// select {
			// case <-sigsExit:
			// 	break
			// }
			// fmt.Println(connections)
			connections.mux.Lock()
			if _, ok := connections.connections["127.0.0.1:13712"]; ok {
				conn := connections.connections["127.0.0.1:13712"]
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
				connections.mux.Unlock()
				buff := make([]byte, 4096)
				n, err := conn.Read(buff)
				res := make([]byte, n)
				copy(res, buff[:n])
				// msg, err := bufio.NewReader(conn).ReadString('\n')
				if n > 0 {
					fmt.Println(res)
				}
				if err != nil {
					// fmt.Println(err)
				} else {
					for {
						connections.mux.Lock()
						if _, ok := connections.connections["127.0.0.1:13711"]; ok {
							conn1 := connections.connections["127.0.0.1:13711"]
							conn1.Write(res)
							connections.mux.Unlock()
							break
							// fmt.Fprintf(connections.connections["127.0.0.1:13721"], msg)
						}
						connections.mux.Unlock()
						// if _, ok := connections.connections["127.0.0.1:13731"]; ok {
						// 	conn2 := connections.connections["127.0.0.1:13731"]
						// 	conn2.Write(res)
						// 	// fmt.Fprintf(connections.connections["127.0.0.1:13731"], msg)
						// }
						// if _, ok := connections.connections["127.0.0.1:13741"]; ok {
						// 	conn3 := connections.connections["127.0.0.1:13741"]
						// 	conn3.Write(res)
						// 	// fmt.Fprintf(connections.connections["127.0.0.1:13741"], msg)
						// }

					}

				}
			} else {
				connections.mux.Unlock()
			}
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
	connections.mux.Lock()
	defer connections.mux.Unlock()
	// if _, ok := connections.connections[r.Host]; ok {
	// 	return
	// }
	opts := &websocket.AcceptOptions{OriginPatterns: []string{"*"}, Subprotocols: []string{"*"}}
	c, err := websocket.Accept(w, r, opts)
	if err != nil {
		s.logf("%v", err)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	conn := websocket.NetConn(ctx, c, 1)
	fmt.Println("Request \"Addr\": ")
	fmt.Println(r.Host)
	dialing := r.Host

	connections.connections[dialing] = conn
	// connections[dialing] = append(connections[dialing], conn)

	// defer c.Close(websocket.StatusInternalError, "the sky is falling")
	// defer cancel()
	fmt.Println("Accepted")
}

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func echo(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}
