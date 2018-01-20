package main

import (
	"net/http"
	"log"
)

func HelloResource(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Server\n")
}

//リクエストに応じたfuncを定義
func main() {
	//routingはここ
	http.HandleFunc("/", HelloResource)
	//この下にrouting定義を増やしていけばいい感じな気がする。

	log.Printf("Start Go HTTP Server")

	//port監視と実行
	port := os.Getenv("PORT")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
