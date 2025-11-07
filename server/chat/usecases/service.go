package services

import (
	"encoding/json"
	hub "server/chat/localclientStorage"
	"server/chat/models"

	"github.com/gorilla/websocket"
)

type service struct {
	hub hub.Hubber
}
type Servicer interface {
	SendAllCli(message *models.Message) error
	SendMessage(message *models.Message) error
	AddClient(username string, conn *websocket.Conn) error
	PopClient(username string) error
	GetAllCli(msg *models.Message) error
	GetOneClie(username string) (*models.Client, error)

	RegistryKey(msg *models.Message) error
	GetKey(msg *models.Message) error
	RegistryIdentityKey(msg *models.Message) error
	KeyObserver(msg *models.Message) bool
}

func NewService(h hub.Hubber) *service {
	return &service{
		hub: h,
	}
}
func (s *service) SendAllCli(message *models.Message) error {
	clies, err := s.hub.GetAllCli()
	if err != nil {
		return err
	}
	for _, c := range clies {
		c.Messages <- message
	}
	return nil
}
func (s *service) SendMessage(message *models.Message) error {
	cli, err := s.hub.GetUsernameClie(message.To)
	if err != nil {
		return err
	}
	cli.Messages <- message
	return nil
}
func (s *service) AddClient(username string, conn *websocket.Conn) error {
	return s.hub.AddClient(username, conn)
}
func (s *service) PopClient(username string) error {
	return s.hub.PopClient(username)
}

func jsoner(s []string) ([]byte, error) {
	return json.Marshal(s)
}
func (s *service) GetAllCli(msg *models.Message) error {
	clis, err := s.hub.GetAllCli()
	if err != nil {
		return err
	}
	usernames := make([]string, 0, len(clis))
	for _, c := range clis {
		usernames = append(usernames, c.Username)
	}
	cli, err := s.hub.GetUsernameClie(msg.From)
	if err != nil {
		return err
	}
	userBytes, err := jsoner(usernames)
	if err != nil {
		return err
	}
	msg.Message = string(userBytes)
	cli.Messages <- msg
	return nil
}
func (s *service) GetOneClie(username string) (*models.Client, error) {
	return s.hub.GetUsernameClie(username)
}

// ////////////////////////////////////////
// //////////////////////////////////////////
func (s *service) RegistryKey(msg *models.Message) error {
	return s.hub.AddKeys(msg.From, msg.Keys.Keys)
}
func (s *service) RegistryIdentityKey(msg *models.Message) error {
	return s.hub.AddIdentityKey(msg.From, msg.Keys.IdentityKey)
}
func (s *service) KeyObserver(msg *models.Message) bool {
	return s.hub.KeyObserver(msg.To)
}
func (s *service) GetKey(msg *models.Message) error {
	peerMessage := &models.Message{
		Date:   msg.Date,
		From:   msg.To,
		To:     msg.From,
		Action: msg.Action,
	}
	senderPublicKey, err := s.hub.GetKey(msg.From)
	if err != nil {
		return err
	}
	senderIndentityKey, err := s.hub.GetIdentityKey(msg.From)
	if err != nil {
		return err
	}
	peerMessage.Keys.IdentityKey = senderIndentityKey
	peerMessage.Keys.Key = senderPublicKey
	peer, err := s.GetOneClie(msg.To)
	if err != nil {
		return err
	}
	peer.Messages <- peerMessage

	senderMessage := &models.Message{
		Date:   msg.Date,
		Action: msg.Action,
		From:   msg.From,
		To:     msg.To,
	}

	peerPublicKey, err := s.hub.GetKey(msg.To)
	if err != nil {
		return err
	}
	peerIdentityKey, err := s.hub.GetIdentityKey(msg.To)
	if err != nil {
		return err
	}
	senderMessage.Keys.IdentityKey = peerIdentityKey
	senderMessage.Keys.Key = peerPublicKey
	sender, err := s.GetOneClie(msg.From)
	if err != nil {
		return err
	}
	sender.Messages <- senderMessage
	return nil

}
