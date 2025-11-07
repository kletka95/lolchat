package hub

import (
	"fmt"
	"server/chat/models"
	"sync"

	"github.com/gorilla/websocket"
)

type Keyer interface {
	AddKeys(username string, keys []string) error
	GetKey(username string) (string, error)
	AddIdentityKey(username, key string) error
	KeyObserver(username string) bool
	GetIdentityKey(username string) (string, error)
}
type CliStorage interface {
	AddClient(username string, conn *websocket.Conn) error
	PopClient(username string) error
	GetAllCli() ([]*models.Client, error)
	GetUsernameClie(username string) (*models.Client, error)
}
type Hubber interface {
	CliStorage
	Keyer
}
type hub struct {
	sync       *sync.RWMutex
	hubStorage map[string]*models.Client
}

func NewHub() *hub {
	return &hub{
		sync:       new(sync.RWMutex),
		hubStorage: map[string]*models.Client{},
	}
}
func (h *hub) AddClient(username string, conn *websocket.Conn) error {
	h.sync.Lock()
	defer h.sync.Unlock()
	if _, ok := h.hubStorage[username]; ok {

		return fmt.Errorf("user already exists")
	}

	h.hubStorage[username] = &models.Client{
		Username:   username,
		Connection: conn,
		Messages:   make(chan *models.Message),
	}
	//h.sync.Unlock()
	return nil
}
func (h *hub) PopClient(username string) error {
	h.sync.Lock()
	defer h.sync.Unlock()
	if v, ok := h.hubStorage[username]; ok {
		close(v.Messages)
		delete(h.hubStorage, username)
	} else {
		return fmt.Errorf("user does not exist")
	}
	return nil
}
func (h *hub) GetAllCli() ([]*models.Client, error) {
	h.sync.RLock()
	defer h.sync.RUnlock()
	clies := make([]*models.Client, 0, len(h.hubStorage))
	for _, c := range h.hubStorage {
		clies = append(clies, c)
	}
	if len(clies) < 1 {
		return nil, fmt.Errorf("client hub is null")
	}
	return clies, nil
}
func (h *hub) GetUsernameClie(username string) (*models.Client, error) {
	h.sync.RLock()
	defer h.sync.RUnlock()
	if v, ok := h.hubStorage[username]; !ok {
		return nil, fmt.Errorf("client does not exists")
	} else {
		return v, nil
	}
}

// keys logic
func (h *hub) GetKey(username string) (string, error) {
	h.sync.Lock()
	defer h.sync.Unlock()
	v, ok := h.hubStorage[username]
	if !ok {
		return "", fmt.Errorf("user not found")
	}
	if v.Key == nil || len(v.Key.Keys) < 1 {
		return "", fmt.Errorf("keys not found")
	}
	result := h.hubStorage[username].Key.Keys[v.Key.Curr]
	h.hubStorage[username].Key.Curr = (v.Key.Curr + 1) % len(v.Key.Keys)
	return result, nil
}
func (h *hub) GetIdentityKey(username string) (string, error) {
	h.sync.RLock()
	defer h.sync.RUnlock()
	if h.hubStorage[username].Key == nil {
		return "", fmt.Errorf("no identity key in stor")
	} else {
		return h.hubStorage[username].Key.IdentityKey, nil
	}
}
func (h *hub) AddKeys(username string, keys []string) error {
	h.sync.Lock()
	defer h.sync.Unlock()
	v, ok := h.hubStorage[username]
	if !ok {
		return fmt.Errorf("user not found")
	}
	var identityKey string
	if v.Key != nil && v.Key.IdentityKey != "" {
		identityKey = v.Key.IdentityKey
	}

	h.hubStorage[username].Key = &models.Key{
		Curr: 0,
		Keys: keys,
		//UpdateTime:  time.Now(),
		IdentityKey: identityKey,
	}

	return nil
}
func (h *hub) KeyObserver(username string) bool {
	h.sync.RLock()
	defer h.sync.RUnlock()
	v, ok := h.hubStorage[username]
	if !ok || v.Key == nil {
		return false
	}

	total := len(v.Key.Keys)
	curr := v.Key.Curr

	// если текущий индекс указывает на последний ключ — дальше нечего брать
	if curr >= total-1 {
		return false
	}

	return true
}
func (h *hub) AddIdentityKey(username, key string) error {
	h.sync.Lock()
	defer h.sync.Unlock()
	if h.hubStorage[username].Key == nil {
		return fmt.Errorf("you must add keys, before add identity key")
	}
	if h.hubStorage[username].Key.IdentityKey != "" {
		return fmt.Errorf("identity key already user for user")
	}

	h.hubStorage[username].Key.IdentityKey = key
	return nil
}
