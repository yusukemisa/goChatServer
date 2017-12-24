package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// 他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte
	// joinはチャットルームに参加しようとするクライアントのためのチャネル
	join chan *cliant
	// leaveはチャットルームから退出しようとするクライアントのためのチャネル
	leave chan *cliant
	// roomに在室するすべてのクライアントを保持
	cliants map[*cliant]bool
}

/*
Channel とは
- Channel は goroutine 間でのメッセージパッシングをするためのもの
- メッセージの型を指定できる
- first class value であり、引数や戻り値にも使える
- send/receive でブロックする
- buffer で、一度に扱えるメッセージ量を指定できる
*/

func (r *room) run() {
	log.Println("room.run")
	for {
		select {
		case cliant := <-r.join:
			//参加
			r.cliants[cliant] = true
		case cliant := <-r.leave:
			//退出
			delete(r.cliants, cliant)
			close(cliant.send)
		case msg := <-r.forward:
			// すべてのクライアントにメッセージを転送
			for cliant := range r.cliants {
				select {
				case cliant.send <- msg:
					// メッセージを送信
				default:
					// 送信に失敗
					delete(r.cliants, cliant)
					close(cliant.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("ServeHTTP")
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServerHTTP:", err)
		return
	}
	cliant := &cliant{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- cliant
	// 無名関数の即時実行？
	defer func() { r.leave <- cliant }()
	go cliant.write()
	cliant.read()
}

func newRoom() *room {
	log.Println("newRoom")
	return &room{
		forward: make(chan []byte),
		join:    make(chan *cliant),
		leave:   make(chan *cliant),
		cliants: make(map[*cliant]bool),
	}
}
