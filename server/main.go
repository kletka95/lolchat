package main

import (
	"log"
	"net/http"
	"server/chat/decrypter"
	hub "server/chat/localclientStorage"
	"server/chat/router"
	services "server/chat/usecases"
	"server/env"
	"server/socket"
)

func main() {
	h := hub.NewHub()

	serv := services.NewService(h)
	r := router.NewRouter(serv)
	envel, err := env.Load()
	if err != nil {
		log.Fatal(err)
	}
	handle := socket.NewHandler(r, decrypter.NewBase64DCR())
	http.HandleFunc("/ws", handle.Handler) // WebSocket endpoint
	log.Println("Server started on", envel.GeneralParse())
	err = http.ListenAndServe(envel.GeneralParse(), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
