package main

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

var serverID uint32
var msgChan chan string
var msgsRecv [][]byte
var msgsSend [][]byte

func main() {
	fmt.Println("TESTING!")
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

	msgsSend = make([][]byte, 5)
	msgsRecv = make([][]byte, 5)

	msgsSend[0] = []byte("Hello World")
	msgsSend[1] = []byte("Hello World 2")
	msgsSend[2] = []byte("Hello World 3")
	msgsSend[3] = []byte("Hello World 4")
	msgsSend[4] = []byte("Hello World 5")

	go func() {
		msgsRecv = append(msgsRecv, nil)
		for {
			var msg []byte
			if len(msgsRecv) > 1 {

				msg, msgsRecv = msgsRecv[0], msgsRecv[1:]
			} else {
				msg, msgsRecv = msgsRecv[0], nil
			}
			if msg != nil {
				msgsSend = append(msgsSend, msg)
			}

		}
	}()

	<-msgChan
	// 	go func() {
	// 		for {
	// 			msg := <-msgChan
	// 			read = msg
	// 			js.Global().Set("read", []byte(msg))
	// 			write = msg
	// 		}
	// 	}()

	// 	for {
	// 		time.Sleep(time.Second * 2)
	// 		fmt.Println("Trying to read")
	// 		fmt.Println(len(read))
	// 		if len(read) > 1 {
	// 			fmt.Println(read)
	// 			read = ""
	// 		} else {
	// 			// write = "Hello"
	// 			// fmt.Println(write)
	// 		}
	// 	}
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
	fmt.Println("Message sent")
	return nil
}

func PassUint8ArrayToGo(this js.Value, args []js.Value) interface{} {

	received := make([]byte, args[0].Get("length").Int())

	_ = js.CopyBytesToGo(received, args[0])

	msgsRecv = append(msgsRecv, received)
	fmt.Println(received)

	return nil
}

func SetUint8ArrayInGo(this js.Value, args []js.Value) interface{} {

	var msg []byte

	if len(msgsSend) > 1 {
		msg, msgsSend = msgsSend[0], msgsSend[1:]
	} else {
		msg, msgsSend = msgsSend[0], nil
	}
	if msg == nil {
		return nil
	}
	_ = js.CopyBytesToJS(args[0], msg)

	return nil
}

func GetArraySize(this js.Value, args []js.Value) interface{} {

	size := make([]byte, 10)

	msgSize := []byte(strconv.Itoa(len(msgsSend[0])))

	copy(size, msgSize)

	_ = js.CopyBytesToJS(args[0], size)

	return nil
}

func registerCallbacks() {
	js.Global().Set("GetSelfID", js.FuncOf(GetSelfID))
	js.Global().Set("SendChat", js.FuncOf(SendChat))

	js.Global().Set("PassUint8ArrayToGo", js.FuncOf(PassUint8ArrayToGo))
	js.Global().Set("SetUint8ArrayInGo", js.FuncOf(SetUint8ArrayInGo))
	js.Global().Set("GetArraySize", js.FuncOf(GetArraySize))
}
