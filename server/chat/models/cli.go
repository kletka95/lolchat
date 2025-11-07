package models

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Connection *websocket.Conn
	Username   string
	Messages   chan *Message
	Key        *Key
}

func (cli *Client) Close() {
	defer func() {
		recover()
	}()
	close(cli.Messages)

}

type Key struct {
	Keys []string
	Curr int
	//UpdateTime  time.Time
	IdentityKey string
}
