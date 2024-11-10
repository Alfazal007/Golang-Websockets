package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var wsManager WebSocketManager

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade the websocket connection")
		return
	}
	defer conn.Close()
	for {
		t, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error in the message", err)
			break
		}
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println("Invalid message", err)
			break
		}
		fmt.Println(msg)

		if msg.MessageType == Connect {
			wsManager.Add(msg.Username, conn)
		} else if msg.MessageType == Disconnect {
			wsManager.Remove(msg.Username, conn)
		} else if msg.MessageType == NormalMessage {
			fmt.Println("Normal message is here")
			wsManager.SendMessage(
				msg.Username,
				msg.Message.Receiver,
				conn,
				msg.Message.Message,
				t)
		}
	}
}

func main() {
	wsManager = WebSocketManager{
		connections: make(map[string]*websocket.Conn),
		mutex:       sync.RWMutex{},
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r)
	})
	err := http.ListenAndServe(":3000", nil)
	fmt.Println(err)
}
