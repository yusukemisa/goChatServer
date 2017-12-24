package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// クライアントはchatしている１ユーザーを表現
type cliant struct {
	socket *websocket.Conn
	// ちゃねる
	send chan []byte
	// クライアントの参加してるチャットルーム
	room *room
}

func (c *cliant) read() {
	log.Println("cliant.read")
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *cliant) write() {
	log.Println("cliant.write")
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
