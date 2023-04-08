package chat

import (
	"encoding/json"
	"errors"
	"time"
)

type Request struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
	Options
}

type Response struct {
	ID      string
	Model   string
	Object  string
	Created time.Time
	Choices []*Choice
	Usage   *Usage
	Error   *Error
}

func (resp *Response) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID      string    `json:"id"`
		Model   string    `json:"Model"`
		Object  string    `json:"object"`
		Created int64     `json:"created"`
		Choices []*Choice `json:"choices"`
		Usage   *Usage    `json:"usage"`
		Error   *Error    `json:"error"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	resp.ID = raw.ID
	resp.Model = raw.Model
	resp.Object = raw.Object
	resp.Created = time.Unix(raw.Created, 0)
	resp.Choices = raw.Choices
	resp.Usage = raw.Usage
	resp.Error = raw.Error

	return nil
}

func (resp *Response) Err() error {
	errMsg := resp.Error.Type + ": " + resp.Error.Message
	return errors.New(errMsg)
}

type FinishReason string

const (
	// API returned complete model output
	Stop FinishReason = "stop"

	// Incomplete model output due to max_tokens parameter or token limit
	Length FinishReason = "length"

	// Omitted content due to a flag from our content filters
	ContentFilter FinishReason = "content_filter"
)

type Choice struct {
	Index        int           `json:"index"`
	Message      *Message      `json:"message,omitempty"`
	Delta        *Message      `json:"delta,omitempty"`
	FinishReason *FinishReason `json:"finish_reason,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Error struct {
	Message string
	Type    string
	Param   string
	Code    string
}
