package main

import (
	"log"
	"net/http"
	"server/chat/decrypter"
	hub "server/chat/localclientStorage"
	"server/chat/router"
	services "server/chat/usecases"
	"server/socket"
)

func main() {
	h := hub.NewHub()
	serv := services.NewService(h)
	r := router.NewRouter(serv)
	handle := socket.NewHandler(r, decrypter.NewBase64DCR())
	http.HandleFunc("/ws", handle.Handler) // WebSocket endpoint
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
