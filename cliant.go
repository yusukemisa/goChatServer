package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// クライアントはchatしている１ユーザーを表現
type cliant struct {
	socket *websocket.Conn
	// ちゃねる
	send chan *message
	// クライアントの参加してるチャットルーム
	room *room
	// userDataはユーザーに関する情報を保持
	userData map[string]interface{}
}

func (c *cliant) read() {
	log.Println("cliant.read")
	for {
		var msg *message
		// JSONを送受信
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Time = msg.When.Format("2006-01-02 15:04:05")
			msg.Name = c.userData["name"].(string)
			// プロフィール画像URL
			if avatarURL, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}

			c.room.forward <- msg
		} else {
			break
		}
		// バイト配列のみ送受信
		// if _, msg, err := c.socket.ReadMessage(); err == nil {
		// 	c.room.forward <- msg
		// } else {
		// 	break
		// }
	}
	c.socket.Close()
}

func (c *cliant) write() {
	log.Println("cliant.write")
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
		// if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
		// 	break
		// }
	}
	c.socket.Close()
}
