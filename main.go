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

	fmt.Printf("websocket serve client on %s at /ws", *clientAddr)

	err := http.ListenAndServe(*clientAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	serveClient()
}
