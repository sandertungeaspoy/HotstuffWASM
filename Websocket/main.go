package main

import (
	"flag"
	"log"
	"net/http"
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
	err := http.ListenAndServe("127.0.0.1:8080", http.FileServer(http.Dir(*dir)))
	log.Fatalln(err)

}
