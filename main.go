package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// read template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates",
			t.filename)))
	})
	if err := t.templ.Execute(w, nil); err != nil {
		log.Fatal("ServeHTTP:", err)
	}
}

func main() {
	// チャットルーム作成
	r := newRoom()
	// テンプレ設定
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// チャットルーム開始
	go r.run()

	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<html>
			<head>
			<title>chat!</title>
			</head>
			<body>
			チャットしようぜ！お前chatサーバーな！
			</body>
		</html>
		`))
	} */

	// webサーバー開始
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	log.Println("webサーバー起動完了")
}
