package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	var i int
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
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
		i++
		fmt.Println("アクセス！:", i)
	})
	// webサーバー開始
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
