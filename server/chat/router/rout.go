package router

import (
	"server/chat/models"
	services "server/chat/usecases"

	"github.com/gorilla/websocket"
)

type rout struct {
	servs services.Servicer
}

func NewRouter(servs services.Servicer) *rout {
	return &rout{
		servs: servs,
	}
}

type Router interface {
	Rout(msg *models.Message)
	CreateCliAndReturn(conn *websocket.Conn, username string) (*models.Client, error)
	PopCli(username string) error
	BroadCast(action models.Action, message string) error
	SystemMessage(action models.Action, username string) error
}

func (r *rout) Rout(msg *models.Message) {
	switch msg.Action {
	case models.GetAllCli:
		r.servs.GetAllCli(msg)
	case models.SendAllCli:
		r.servs.SendAllCli(msg)
	case models.SendMessage:
		r.servs.SendMessage(msg)
	case models.RegistryKey:
		r.servs.RegistryKey(msg)
		r.servs.RegistryIdentityKey(msg)
		r.BroadCast(models.UpdateKey, msg.From)
	case models.GetKey:
		r.servs.GetKey(msg)
		if !r.servs.KeyObserver(msg) {
			r.SystemMessage(models.NeedKey, msg.To)
		}
	}
}
func (r *rout) SystemMessage(action models.Action, username string) error {
	msg := &models.Message{
		To:     username,
		From:   "system",
		Action: action,
	}

	return r.servs.SendMessage(msg)
}
func (r *rout) BroadCast(action models.Action, message string) error {
	msg := &models.Message{
		Action:  action,
		Message: message,
		From:    "system",
		To:      "all",
	}
	return r.servs.SendAllCli(msg)
}
func (r *rout) CreateCliAndReturn(conn *websocket.Conn, username string) (*models.Client, error) {
	err := r.servs.AddClient(username, conn)
	if err != nil {
		return nil, err
	}
	return r.servs.GetOneClie(username)
}
func (r *rout) PopCli(username string) error {
	return r.servs.PopClient(username)
}
