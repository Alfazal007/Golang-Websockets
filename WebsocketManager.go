package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	Connect       MessageType = "connect"
	NormalMessage MessageType = "message"
	Disconnect    MessageType = "disconnect"
)

type Message struct {
	Message     MessageInside `json:"message"`
	Username    string        `json:"username"`
	MessageType MessageType   `json:"messageType"`
}

type MessageInside struct {
	Message  string `json:"message"`
	Receiver string `json:"receiver"`
}

type WebSocketManager struct {
	connections map[string]*websocket.Conn
	mutex       sync.RWMutex
}

func (websocketmanager *WebSocketManager) Add(username string, connection *websocket.Conn) {
	websocketmanager.mutex.RLock()
	_, exist := websocketmanager.connections[username]
	websocketmanager.mutex.RUnlock()
	if exist {
		fmt.Println("Already exists")
		return
	} else {
		websocketmanager.mutex.Lock()
		websocketmanager.connections[username] = connection
		websocketmanager.mutex.Unlock()
	}
}

func (websocketmanager *WebSocketManager) Remove(username string, connection *websocket.Conn) {
	websocketmanager.mutex.RLock()
	conn, exist := websocketmanager.connections[username]
	websocketmanager.mutex.RUnlock()

	if exist && conn == connection {
		conn.Close()
		websocketmanager.mutex.Lock()
		delete(websocketmanager.connections, username)
		websocketmanager.mutex.Unlock()
	} else {
		conn.Close()
	}
}

func (websocketmanager *WebSocketManager) SendMessage(
	sender string,
	receiver string,
	senderConn *websocket.Conn,
	message string,
	messageType int,
) {
	fmt.Println("Called the send message")
	websocketmanager.mutex.RLock()
	senderConnExisting, exist := websocketmanager.connections[sender]
	websocketmanager.mutex.RUnlock()
	if !exist {
		fmt.Println("Sender does not exist")
		senderConn.Close()
		return
	} else if senderConnExisting != senderConn {
		senderConn.Close()
		fmt.Println("Not same")
		return
	}
	websocketmanager.mutex.RLock()
	receiverConnExisting, exist := websocketmanager.connections[receiver]
	websocketmanager.mutex.RUnlock()
	if !exist {
		fmt.Println("Receiver does not exist")
		return
	}
	receiverConnExisting.WriteMessage(messageType, []byte(message))
}
