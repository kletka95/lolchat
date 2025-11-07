package socket

import (
	"log"
	"net/http"
	"server/chat/decrypter"
	"server/chat/models"
	"server/chat/router"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем все origin, иначе браузер блокирует
	},
}

type handler struct {
	r router.Router
	t decrypter.Transport
}

func NewHandler(r router.Router, t decrypter.Transport) *handler {
	return &handler{
		r: r,
		t: t,
	}
}
func (h *handler) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte("input your login"))
	_, p, _ := conn.ReadMessage()
	if string(p) == "null" {
		msg := &models.Message{
			Action:  models.Disconnected,
			Message: "name coudnt be null",
		}
		b, _ := h.t.CryptedMsg(msg)
		conn.WriteMessage(websocket.TextMessage, b)
		return
	}
	defer conn.Close()
	cli, err := h.r.CreateCliAndReturn(conn, string(p))
	if err != nil {
		msg := &models.Message{
			Action:  models.Disconnected,
			Message: string(p) + " client already exists",
		}
		b, _ := h.t.CryptedMsg(msg)
		conn.WriteMessage(websocket.TextMessage, b)
		return
	}
	defer h.r.PopCli(cli.Username)
	defer h.r.BroadCast(models.Disconnected, cli.Username)

	go func() {
		for msg := range cli.Messages {
			b, err := h.t.CryptedMsg(msg)
			if err != nil {
				log.Println(err)
				h.r.SystemMessage(models.Disconnected, cli.Username)
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
				log.Println(err)
				h.r.SystemMessage(models.Disconnected, cli.Username)
				return
			}
		}
	}()
	h.r.SystemMessage(models.NeedKey, cli.Username)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		msg, err := h.t.MessageDecrypt(p)
		if err != nil {
			log.Println(err)
			h.r.SystemMessage(models.Disconnected, cli.Username)
			return
		}
		h.r.Rout(msg)
	}
}
