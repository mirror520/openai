package openai

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/mirror520/openai/chat"
)

type ChatEndpoints struct {
	CreateChatEndpoint endpoint.Endpoint
	UpdateChatEndpoint endpoint.Endpoint
	ChatEndpoint       endpoint.Endpoint
	ChatStreamEndpoint endpoint.Endpoint
}

type CreateChatRequest struct {
	Model   string          `json:"model"`
	Prompt  string          `json:"prompt"`
	Options json.RawMessage `json:"options"`
}

func CreateChatEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(*CreateChatRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		id, err := svc.CreateChat(req.Model, req.Prompt, req.Options)
		if err != nil {
			return nil, err
		}

		return &id, nil
	}
}

type UpdateChatRequest struct {
	ID      chat.ChatID     `json:"-"`
	Model   string          `json:"model"`
	Prompt  string          `json:"prompt"`
	Options json.RawMessage `json:"options"`
}

func UpdateChatEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(*UpdateChatRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := svc.UpdateChat(req.Model, req.Prompt, req.Options, req.ID); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

type ChatRequest struct {
	ID      chat.ChatID `json:"-"`
	Content string      `json:"content"`
}

func ChatEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(*ChatRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.Chat(req.Content, req.ID)
	}
}

func ChatStreamEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(*ChatRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.ChatStream(req.Content, req.ID)
	}
}
