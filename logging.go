package openai

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/mirror520/openai/chat"
)

func LoggingMiddleware(log *zap.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			log.With(zap.String("service", "openai")),
			next,
		}
	}
}

type loggingMiddleware struct {
	log  *zap.Logger
	next Service
}

func (mw *loggingMiddleware) CreateChat(model string, prompt string, rawOpts json.RawMessage) (chat.ChatID, error) {
	log := mw.log.With(
		zap.String("action", "create_chat"),
	)

	id, err := mw.next.CreateChat(model, prompt, rawOpts)
	if err != nil {
		log.Error(err.Error())
		return chat.ChatID{}, err
	}

	log.Info("done", zap.String("chat_id", id.String()))
	return id, nil
}

func (mw *loggingMiddleware) UpdateChat(model string, prompt string, rawOpts json.RawMessage, id chat.ChatID) error {
	log := mw.log.With(
		zap.String("action", "update_chat"),
		zap.String("chat_id", id.String()),
	)

	err := mw.next.UpdateChat(model, prompt, rawOpts, id)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("done")
	return nil
}

func (mw *loggingMiddleware) Chat(content string, id chat.ChatID) (string, error) {
	log := mw.log.With(
		zap.String("action", "chat"),
		zap.String("chat_id", id.String()),
		zap.String("ask", content),
	)

	answer, err := mw.next.Chat(content, id)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	return answer, nil
}

func (mw *loggingMiddleware) ChatStream(content string, id chat.ChatID) (<-chan string, error) {
	log := mw.log.With(
		zap.String("action", "chat_stream"),
		zap.String("chat_id", id.String()),
		zap.String("ask", content),
	)

	stream, err := mw.next.ChatStream(content, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("get stream")
	return stream, nil
}
