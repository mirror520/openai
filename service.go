package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/mirror520/openai/chat"
	"github.com/mirror520/openai/conf"
)

type Service interface {
	CreateChat(model string, opts json.RawMessage) (chat.ChatID, error)
	UpdateChat(model string, opts json.RawMessage, id chat.ChatID) error
	Chat(content string, id chat.ChatID) (string, error)
	ChatStream(content string, id chat.ChatID) (<-chan string, error)
}

func NewService(chats chat.Repository, cfg *conf.Config) Service {
	return &service{
		log: zap.L().With(
			zap.String("service", "openai"),
		),
		chats:  chats,
		apiKey: cfg.APIKey,
	}
}

type service struct {
	log    *zap.Logger
	chats  chat.Repository
	apiKey string
}

func (svc *service) CreateChat(model string, opts json.RawMessage) (chat.ChatID, error) {
	ctx, err := chat.NewContext(model, opts)
	if err != nil {
		return chat.ChatID{}, err
	}

	if err := svc.chats.Store(ctx); err != nil {
		return chat.ChatID{}, err
	}

	return ctx.ID, nil
}

func (svc *service) UpdateChat(model string, opts json.RawMessage, id chat.ChatID) error {
	ctx, err := svc.chats.Find(id)
	if err != nil {
		return err
	}

	ctx.Model = model

	if err := ctx.Options.Update(opts); err != nil {
		return err
	}

	if err := svc.chats.Store(ctx); err != nil {
		return err
	}

	return nil
}

func (svc *service) Chat(content string, id chat.ChatID) (string, error) {
	ctx, err := svc.chats.Find(id)
	if err != nil {
		return "", err
	}

	ctx.AddMessage(&chat.Message{
		Role:    chat.User,
		Content: content,
	})

	bs, err := json.Marshal(ctx.Request())
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bs))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+svc.apiKey)

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result *chat.Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", result.Err()
	}

	if len(result.Choices) == 0 {
		return "", errors.New("empty choices")
	}

	for _, choice := range result.Choices {
		ctx.AddMessage(choice.Message)
	}

	if err := svc.chats.Store(ctx); err != nil {
		return "", err
	}

	return result.Choices[0].Message.Content, nil
}

func (svc *service) ChatStream(content string, id chat.ChatID) (<-chan string, error) {
	ctx, err := svc.chats.Find(id)
	if err != nil {
		return nil, err
	}

	ctx.AddMessage(&chat.Message{
		Role:    chat.User,
		Content: content,
	})

	reqMsg := ctx.Request()
	*reqMsg.Options.Stream = true

	bs, err := json.Marshal(reqMsg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+svc.apiKey)

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var failedResult *chat.Response
		if err := json.NewDecoder(resp.Body).Decode(&failedResult); err != nil {
			return nil, err
		}

		return nil, failedResult.Err()
	}

	data := make(chan string, 1)

	go svc.stream(ctx, resp.Body, data)

	return data, nil
}

func (svc *service) stream(ctx *chat.Context, reader io.ReadCloser, data chan<- string) error {
	log := svc.log.With(
		zap.String("action", "chat_stream"),
		zap.String("chat_id", ctx.ID.String()),
	)

	defer reader.Close()
	defer close(data)

	msg := new(chat.Message)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var chunk *chat.Response

		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			log.Error(err.Error())
			return err
		}

		if len(chunk.Choices) < 1 {
			err := errors.New("invalid choices")
			log.Error(err.Error())
			return err
		}

		choice := chunk.Choices[0]

		finish := choice.FinishReason
		if finish != nil && *finish == chat.Stop {
			log.Info("done")
			break
		}

		if choice.Delta.Role != "" {
			msg.Role = choice.Delta.Role

			log.Debug("chunk",
				zap.String("role", string(choice.Delta.Role)),
			)
		}

		if choice.Delta.Content != "" {
			data <- msg.Content

			msg.Content += choice.Delta.Content

			log.Debug("chunk",
				zap.String("content", choice.Delta.Content),
				zap.String("full_content", msg.Content),
			)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error(err.Error())
		return err
	}

	ctx.AddMessage(msg)

	if err := svc.chats.Store(ctx); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
