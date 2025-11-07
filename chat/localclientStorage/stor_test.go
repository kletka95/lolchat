package hub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStr(t *testing.T) {
	hub := NewHub()
	_, err := hub.GetKey("ok")
	assert.Error(t, err)
	err = hub.AddClient("ok", nil)
	assert.NoError(t, err)
	_, err = hub.GetKey("ok")
	assert.Error(t, err)
	err = hub.AddKeys("ok", []string{"1", "2", "3"})
	assert.NoError(t, err)
	equal := []string{"1", "2", "3", "1"}
	for i := 0; i < 4; i++ {
		key, err := hub.GetKey("ok")
		assert.Equal(t, equal[i], key)
		assert.NoError(t, err)
	}

}
