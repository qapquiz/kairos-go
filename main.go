package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"kairos-go/client"
)

var clientAddr = flag.String("clientAddr", ":8080", "http service address")

func serveClient() {
	flag.Parse()

	clientHub := client.NewHub()
	go clientHub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		client.ServeWs(clientHub, w, r)
	})

	fmt.Println("websocket serve client on :8080 at /ws")

	err := http.ListenAndServe(*clientAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	serveClient()
}
