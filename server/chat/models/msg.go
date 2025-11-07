package models

import "time"

type Message struct {
	Action  Action     `json:"Action"`
	Message string     `json:"Message"`
	Date    int        `json:"Datatime"`
	To      string     `json:"To,omitempty"`
	From    string     `json:"From,omitempty"`
	Keys    MessageKey `json:"Keys,omitempty"`
}
type MessageKey struct {
	IdentityKey string    `json:"Identity"`
	Keys        []string  `json:"Keys"`
	Key         string    `json:"Key"`
	CurrTime    time.Time `json:"CurrTime"`
	Signature   string    `json:"Signature"`
	Iv          string    `json:"Iv"`
}

type Action int32

const (
	SendMessage Action = 1
	SendAllCli  Action = 2
	GetAllCli   Action = 3
	CreateCli   Action = 4

	GetKey       Action = 8
	RegistryKey  Action = 9
	UpdateKey    Action = 10
	NeedKey      Action = 12
	GetSenderKey Action = 15

	Disconnected Action = 11
)
