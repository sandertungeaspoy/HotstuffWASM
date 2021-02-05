package main

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	"nhooyr.io/websocket"
)

var serverID uint32
var msgChan chan string

func main() {
	fmt.Println("TEST!")
	fmt.Println("This is running in Webassembly!")

	fmt.Println("Initializing")
	registerCallbacks()

	msgChan = make(chan string)
	serverID = uint32(0)
	for {
		if serverID != 0 {
			break
		}
		fmt.Print("Sleeping ZzZ ID: ")
		fmt.Println(serverID)
		time.Sleep(1 * time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stringAddr := "ws://localhost:1337" + fmt.Sprint(serverID)
	// opts := &websocket.DialOptions{Subprotocols: []string{"echo"}}
	// fmt.Println(opts.Subprotocols)
	c, _, err := websocket.Dial(ctx, stringAddr, nil)
	if err != nil {
		fmt.Println("Dial Failed")
		fmt.Println(err)
	} else {
		conn := websocket.NetConn(ctx, c, 1)
		go func() {
			for {
				msg := <-msgChan
				fmt.Fprintf(conn, msg)
			}
		}()
		for {
			fmt.Println("Trying to read")
			status, err := bufio.NewReader(conn).ReadBytes('\n')
			fmt.Println("Read")
			if err != nil {
				fmt.Println("Read Failed")
				fmt.Println(err)
				return
			}
			fmt.Println(status)
		}
	}
}

func GetSelfID(this js.Value, i []js.Value) interface{} {
	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()

	selfID, _ := strconv.ParseUint(value1, 10, 32)
	serverID = uint32(selfID)
	fmt.Println(serverID)
	return nil
}

func SendChat(this js.Value, i []js.Value) interface{} {
	msg := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()

	msgChan <- msg
	return nil
}

func registerCallbacks() {
	js.Global().Set("GetSelfID", js.FuncOf(GetSelfID))
	js.Global().Set("SendChat", js.FuncOf(SendChat))
}
