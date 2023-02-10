package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	server := NewServer()

	http.Handle("/", websocket.Handler(server.HandleSocket))
	http.Handle("/broadcast", websocket.Handler(server.HandleBroadcast))

	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
