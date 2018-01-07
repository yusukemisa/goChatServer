package main

import (
	"log"
	"net/http"
	"trace"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	// 他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan *message
	// joinはチャットルームに参加しようとするクライアントのためのチャネル
	join chan *cliant
	// leaveはチャットルームから退出しようとするクライアントのためのチャネル
	leave chan *cliant
	// roomに在室するすべてのクライアントを保持
	cliants map[*cliant]bool
	// tracerはチャットルーム上で行われた操作ログを受け取ります
	tracer trace.Tracer
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
			r.tracer.Trace("新しいクライアントが参加しました")
		case cliant := <-r.leave:
			//退出
			delete(r.cliants, cliant)
			close(cliant.send)
			r.tracer.Trace("クライアントが退出しました")
		case msg := <-r.forward:
			// すべてのクライアントにメッセージを転送
			for cliant := range r.cliants {
				select {
				case cliant.send <- msg:
					// メッセージを送信
					r.tracer.Trace(" -- クライアントに送信しました")
				default:
					// 送信に失敗
					delete(r.cliants, cliant)
					close(cliant.send)
					r.tracer.Trace(" -- 送信に失敗しました。クライアントをクリーンアップします")
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
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました:", err)
		return
	}
	log.Println(objx.MustFromBase64(authCookie.Value))
	cliant := &cliant{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
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
		forward: make(chan *message),
		join:    make(chan *cliant),
		leave:   make(chan *cliant),
		cliants: make(map[*cliant]bool),
		tracer:  trace.Off(),
	}
}
