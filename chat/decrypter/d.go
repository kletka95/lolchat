package decrypter

import (
	"encoding/base64"
	"encoding/json"
	"server/chat/models"
)

func b64Decoder(msg []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(msg))

}
func b64Encode(msg []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(msg))
}

type Transport interface {
	MessageDecrypt(msg []byte) (*models.Message, error)
	CryptedMsg(msg *models.Message) ([]byte, error)
}
type base64dcr struct {
}

func NewBase64DCR() *base64dcr {
	return &base64dcr{}
}
func (d *base64dcr) MessageDecrypt(msg []byte) (*models.Message, error) {
	enc, err := b64Decoder(msg)
	if err != nil {

		return nil, err
	}
	//fmt.Println(string(enc))
	message := new(models.Message)
	err = json.Unmarshal(enc, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
func (d *base64dcr) CryptedMsg(msg *models.Message) ([]byte, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return b64Encode(b), nil
}
