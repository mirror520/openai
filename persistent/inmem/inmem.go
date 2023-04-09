package inmem

import (
	"errors"
	"sync"

	"github.com/mirror520/openai/chat"
)

func NewChatRepository() chat.Repository {
	return &chatRepository{
		chats: make(map[chat.ChatID]*chat.Chat),
	}
}

type chatRepository struct {
	chats map[chat.ChatID]*chat.Chat
	sync.RWMutex
}

func (repo *chatRepository) Store(c *chat.Chat) error {
	repo.Lock()
	repo.chats[c.ID] = c
	repo.Unlock()
	return nil
}

func (repo *chatRepository) Find(id chat.ChatID) (*chat.Chat, error) {
	repo.RLock()
	defer repo.RUnlock()

	c, ok := repo.chats[id]
	if !ok {
		return nil, errors.New("chat not found")
	}

	return c, nil
}

func (repo *chatRepository) Close() error {
	repo.chats = nil
	return nil
}
