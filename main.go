package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	gom "github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

// read template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// setting ////////////////////
type authEnv struct {
	SecurityKey string         `json:"securityKey"`
	Provider    []authProvider `json:"providers"`
}
type authProvider struct {
	Name        string `json:"name"`
	CliantID    string `json:"cliantId"`
	SecretKey   string `json:"secretKey"`
	CallbackURL string `json:"callbackUrl"`
}

// setting ////////////////////

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
	var envFilePath = flag.String("env", "env.json", "環境定義")

	jsonBytes, err := ioutil.ReadFile(*envFilePath)
	if err != nil {
		log.Fatal("ReadFile:", err)
	}

	env := new(authEnv)
	if err := json.Unmarshal(jsonBytes, env); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return
	}

	flag.Parse()
	// 認証処理 ////////////////////
	// セキュリティキー
	gom.SetSecurityKey(env.SecurityKey)
	gom.WithProviders(
		google.New(env.Provider[0].CliantID, env.Provider[0].SecretKey, env.Provider[0].CallbackURL),
		github.New(env.Provider[1].CliantID, env.Provider[1].SecretKey, env.Provider[1].CallbackURL),
		facebook.New(env.Provider[2].CliantID, env.Provider[2].SecretKey, env.Provider[2].CallbackURL),
	)
	// チャットルーム作成
	r := newRoom()
	//r.tracer = trace.New(os.Stdout)
	// テンプレ設定
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
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
