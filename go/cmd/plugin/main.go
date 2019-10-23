package main

/*
#include <stdlib.h>

typedef char* (*message_callback) (char *msg);

static inline char* call_message_callback(message_callback ptr, char *msg) {
  return (ptr)(msg);
}
*/
import "C"
import (
	"client"
	"fmt"
	"log"
	"unsafe"

	"github.com/gorilla/websocket"
)

var c *client.Client

//export Start
func Start(url string, user string, password string, fn C.message_callback) int {
	c = client.NewClient(url, user, password)
	c.OnMessageCallback = func(message []byte) string {
		cmsg := C.CString(string(message))
		defer C.free(unsafe.Pointer(cmsg))
		resp := C.call_message_callback(fn, cmsg)
		if resp == nil {
			return ""
		}
		return C.GoString(resp)
	}
	if err := c.Start(); err != nil {
		fmt.Println(err.Error())
		c = nil
		return 1
	}
	return 0
}

//export Stop
func Stop() {
	if c != nil {
		c.Stop()
		c = nil
	}
}

//export SendMessage
func SendMessage(msg string) {
	if c == nil {
		return
	}
	fmt.Printf("SendMessage: %s\n", msg)
	if c.WsConn == nil {
		fmt.Println("Connection error")
		return
	}
	err := c.WsConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func main() {}

/*
Build as library:
go build -buildmode=c-shared -o gisquick.so cmd/plugin/main.go
*/
