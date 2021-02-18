package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

// echoServer is the WebSocket echo server implementation.
// It ensures the client speaks the echo subprotocol and
// only allows one message every 100ms with a 10 message burst.
type echoServer struct {
	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
}

var (
	listen = flag.String("listen", "127.0.0.1:8080", "listen address")
	dir    = flag.String("dir", ".", "directory to serve")
)

func main() {
	flag.Parse()
	log.Printf("listening on %q...", *listen)
	go func() {
		err := http.ListenAndServe("127.0.0.1:8080", http.FileServer(http.Dir(*dir)))
		log.Fatalln(err)
	}()

	l, err := net.Listen("tcp", "127.0.0.1:13371")
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("listening on http://%v", l.Addr())

	s := &http.Server{
		Handler: echoServer{
			logf: log.Printf,
		},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		fmt.Println("Failed At serve")
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

}

func (s echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Failed at ServeHTTP")
	opts := &websocket.AcceptOptions{OriginPatterns: []string{"*"}} // Subprotocols: []string{"echo"}}
	c, err := websocket.Accept(w, r, opts)
	if err != nil {
		s.logf("%v", err)
		return
	}

	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	fmt.Println("Accpted")
	fmt.Println(c.Subprotocol())
	// if c.Subprotocol() != "echo" {
	// 	c.Close(websocket.StatusPolicyViolation, "client must speak the echo subprotocol")
	// 	return
	// }
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	conn := websocket.NetConn(ctx, c, 1)
	go func() {
		for {

			buff := make([]byte, 4096)
			n, err := conn.Read(buff)
			res := make([]byte, n)
			copy(res, buff[:n])
			// _, msg, err := c.Read(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res)
		}
	}()
	var (
		scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	)
	for {
		fmt.Print("Enter message to send")
		scanner.Scan()
		msg := scanner.Text()
		msg += "\n"
		// writer, err := c.Writer(ctx, 1)
		// fmt.Println([]byte(msg))
		// writer.Write([]byte(msg))
		// writer.Close()
		// fmt.Println("Wrote")
		conn.Write([]byte(msg))
		if err != nil {
			fmt.Println(err)
		}
	}

	// fmt.Println("Is echo")
	// l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	// for {
	// 	err = echo(r.Context(), c, l)
	// 	if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
	// 		return
	// 	}
	// 	if err != nil {
	// 		s.logf("failed to echo with %v: %v", r.RemoteAddr, err)
	// 		return
	// 	}
	// }
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
