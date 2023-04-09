package http

import (
	"encoding/json"

	"github.com/mirror520/openai"
	"github.com/mirror520/openai/chat"
)

// TODO
func ProxyMiddleware() openai.ServiceMiddleware {
	return func(next openai.Service) openai.Service {
		return &proxyMiddleware{next}
	}
}

type proxyMiddleware struct {
	next openai.Service
}

func (mw *proxyMiddleware) CreateChat(model string, prompt string, rawOpts json.RawMessage) (chat.ChatID, error) {
	return mw.next.CreateChat(model, prompt, rawOpts)
}

func (mw *proxyMiddleware) UpdateChat(model string, prompt string, rawOpts json.RawMessage, id chat.ChatID) error {
	return mw.next.UpdateChat(model, prompt, rawOpts, id)
}

func (mw *proxyMiddleware) Chat(content string, id chat.ChatID) (string, error) {
	return mw.next.Chat(content, id)
}

func (mw *proxyMiddleware) ChatStream(content string, id chat.ChatID) (<-chan string, error) {
	return mw.next.ChatStream(content, id)
}
