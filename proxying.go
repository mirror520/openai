package openai

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mirror520/openai/chat"
)

func ProxyingMiddleware(endpoints *ChatEndpoints) ServiceMiddleware {
	return func(next Service) Service {
		return &proxyingMiddleware{endpoints}
	}
}

type proxyingMiddleware struct {
	*ChatEndpoints
}

func (mw *proxyingMiddleware) CreateChat(model string, prompt string, rawOpts json.RawMessage) (chat.ChatID, error) {
	req := &CreateChatRequest{
		Model:   model,
		Prompt:  prompt,
		Options: rawOpts,
	}

	resp, err := mw.CreateChatEndpoint(context.Background(), req)
	if err != nil {
		return chat.ChatID{}, err
	}

	id, ok := resp.(chat.ChatID)
	if !ok {
		return chat.ChatID{}, errors.New("invalid response")
	}

	return id, nil
}

func (mw *proxyingMiddleware) UpdateChat(model string, prompt string, rawOpts json.RawMessage, id chat.ChatID) error {
	req := &UpdateChatRequest{
		ID:      id,
		Model:   model,
		Prompt:  prompt,
		Options: rawOpts,
	}

	_, err := mw.UpdateChatEndpoint(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

func (mw *proxyingMiddleware) Chat(content string, id chat.ChatID) (string, error) {
	req := &ChatRequest{
		ID:      id,
		Content: content,
	}

	resp, err := mw.ChatEndpoint(context.Background(), req)
	if err != nil {
		return "", err
	}

	answer, ok := resp.(string)
	if !ok {
		return "", errors.New("invalid response")
	}

	return answer, nil
}

func (mw *proxyingMiddleware) ChatStream(content string, id chat.ChatID) (<-chan string, error) {
	req := &ChatRequest{
		ID:      id,
		Content: content,
	}

	resp, err := mw.ChatStreamEndpoint(context.Background(), req)
	if err != nil {
		return nil, err
	}

	stream, ok := resp.(<-chan string)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return stream, nil
}
