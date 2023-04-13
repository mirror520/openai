package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-resty/resty/v2"

	"github.com/mirror520/openai"
	"github.com/mirror520/openai/chat"
	"github.com/mirror520/openai/model"
)

type MakeEndpoint func(baseURL string) endpoint.Endpoint

func ChatFactory(makeEndpoint MakeEndpoint, scheme string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		baseURL := fmt.Sprintf("%s://%s/openai/v1", scheme, instance)
		return makeEndpoint(baseURL), nil, nil
	}
}

func CreateChatEndpoint(baseURL string) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		var result model.Result

		client := resty.New().
			SetBaseURL(baseURL)

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(&request).
			SetResult(&result).
			Post("/chats")

		if err != nil {
			return nil, err
		}

		if resp.StatusCode() != http.StatusOK {
			if result.Status == model.FAILURE {
				return nil, errors.New(result.Msg)
			}

			return nil, errors.New(resp.Status())
		}

		idStr, ok := result.Data.(string)
		if !ok {
			return nil, errors.New("invalid response")
		}

		id, err := chat.ParseID(idStr)
		if err != nil {
			return nil, err
		}

		return id, nil
	}
}

func UpdateChatEndpoint(baseURL string) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		var result model.Result

		client := resty.New().
			SetBaseURL(baseURL)

		req, ok := request.(*openai.UpdateChatRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(&request).
			SetResult(&result).
			Patch("/chats/" + req.ID.String())

		if err != nil {
			return nil, err
		}

		if resp.StatusCode() != http.StatusOK {
			if result.Status == model.FAILURE {
				return nil, errors.New(result.Msg)
			}

			return nil, errors.New(resp.Status())
		}

		return nil, nil
	}
}

func ChatEndpoint(baseURL string) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		// TODO
		return nil, nil
	}
}

func ChatStreamEndpoint(baseURL string) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		// TODO
		return nil, nil
	}
}
