package main

import (
	"flag"
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
	//log.Println(r.Host, r.Header)
	if err := t.templ.Execute(w, r); err != nil {
		log.Fatal("ServeHTTP:", err)
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	// チャットルーム作成
	r := newRoom()
	//r.tracer = trace.New(os.Stdout)
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

	log.Println("webサーバー起動. ポート", *addr)
	// webサーバー開始
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
