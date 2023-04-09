package chat

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTypicalRequest(t *testing.T) {
	assert := assert.New(t)

	jsonStr := `
		{
		  "model": "gpt-3.5-turbo",
		  "messages": [
			{
			  "role": "user", 
			  "content": "Count to 100, with a comma between each number and no newlines. E.g., 1, 2, 3, ..."
			}
		  ],
		  "temperature": 0
		}`

	var req *Request
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		assert.Fail(err.Error())
		return
	}

	assert.Equal("gpt-3.5-turbo", req.Model)

	assert.Len(req.Messages, 1)
	assert.Equal(User, req.Messages[0].Role)

	assert.Equal(0.0, *req.Temperature)
}

func TestTypicalResponse(t *testing.T) {
	assert := assert.New(t)

	jsonStr := `
		{
		  "choices": [
		    {
		      "finish_reason": "stop",
		      "index": 0,
		      "message": {
		        "content": "\n\n1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100.",
		        "role": "assistant"
		      }
		    }
		  ],
		  "created": 1677825456,
		  "id": "chatcmpl-6ptKqrhgRoVchm58Bby0UvJzq2ZuQ",
		  "model": "gpt-3.5-turbo-0301",
		  "object": "chat.completion",
		  "usage": {
		    "completion_tokens": 301,
		    "prompt_tokens": 36,
		    "total_tokens": 337
		  }
		}`

	var resp *Response
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		assert.Fail(err.Error())
		return
	}

	assert.Equal("chatcmpl-6ptKqrhgRoVchm58Bby0UvJzq2ZuQ", resp.ID)
	assert.Equal("gpt-3.5-turbo-0301", resp.Model)
	assert.Equal("chat.completion", resp.Object)
	assert.Equal(time.Unix(1677825456, 0), resp.Created)

	assert.Len(resp.Choices, 1)
	assert.Equal(0, resp.Choices[0].Index)
	assert.Equal(Assistant, resp.Choices[0].Message.Role)
	assert.Equal(Stop, *resp.Choices[0].FinishReason)

	assert.Equal(36, resp.Usage.PromptTokens)
	assert.Equal(301, resp.Usage.CompletionTokens)
	assert.Equal(337, resp.Usage.TotalTokens)
}

func TestStreamRequest(t *testing.T) {
	assert := assert.New(t)

	jsonStr := `
		{
		  "model": "gpt-3.5-turbo",
		  "messages": [
			{
			  "role": "user", 
			  "content": "What's 1+1? Answer in one word."
			}
		  ],
		  "temperature": 0,
		  "stream": true
		}`

	var req *Request
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		assert.Fail(err.Error())
		return
	}

	assert.Equal("gpt-3.5-turbo", req.Model)

	assert.Len(req.Messages, 1)
	assert.Equal(User, req.Messages[0].Role)

	assert.Equal(0.0, *req.Temperature)
	assert.True(*req.Stream)
}

func TestStreamResponse(t *testing.T) {
	assert := assert.New(t)

	{
		jsonStr := `
			{
			  "choices": [
			    {
			      "delta": {
			        "role": "assistant"
			      },
			      "finish_reason": null,
			      "index": 0
			    }
			  ],
			  "created": 1677825464,
			  "id": "chatcmpl-6ptKyqKOGXZT6iQnqiXAH8adNLUzD",
			  "model": "gpt-3.5-turbo-0301",
			  "object": "chat.completion.chunk"
			}`

		var resp *Response
		if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
			assert.Fail(err.Error())
			return
		}

		assert.Equal("chatcmpl-6ptKyqKOGXZT6iQnqiXAH8adNLUzD", resp.ID)
		assert.Equal("gpt-3.5-turbo-0301", resp.Model)
		assert.Equal("chat.completion.chunk", resp.Object)
		assert.Equal(time.Unix(1677825464, 0), resp.Created)

		assert.Len(resp.Choices, 1)
		assert.Equal(0, resp.Choices[0].Index)
		assert.Equal(Assistant, resp.Choices[0].Delta.Role)
		assert.Nil(resp.Choices[0].Message)
		assert.Nil(resp.Choices[0].FinishReason)
	}

	{
		jsonStr := `
			{
			  "choices": [
			    {
			      "delta": {
			        "content": "\n\n"
			      },
			      "finish_reason": null,
			      "index": 0
			    }
			  ],
			  "created": 1677825464,
			  "id": "chatcmpl-6ptKyqKOGXZT6iQnqiXAH8adNLUzD",
			  "model": "gpt-3.5-turbo-0301",
			  "object": "chat.completion.chunk"
			}`

		var resp *Response
		if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
			assert.Fail(err.Error())
			return
		}

		assert.Equal("chatcmpl-6ptKyqKOGXZT6iQnqiXAH8adNLUzD", resp.ID)
		assert.Equal("gpt-3.5-turbo-0301", resp.Model)
		assert.Equal("chat.completion.chunk", resp.Object)
		assert.Equal(time.Unix(1677825464, 0), resp.Created)

		assert.Len(resp.Choices, 1)
		assert.Equal(0, resp.Choices[0].Index)
		assert.Equal(Role(""), resp.Choices[0].Delta.Role)
		assert.Equal("\n\n", resp.Choices[0].Delta.Content)
		assert.Nil(resp.Choices[0].Message)
		assert.Nil(resp.Choices[0].FinishReason)
	}

	{
		jsonStr := `
			{
			  "choices": [
			    {
			      "delta": {
			        "content": "2"
			      },
			      "finish_reason": null,
			      "index": 0
			    }
			  ],
			  "created": 1677825464,
			  "id": "chatcmpl-6ptKyqKOGXZT6iQnqiXAH8adNLUzD",
			  "model": "gpt-3.5-turbo-0301",
			  "object": "chat.completion.chunk"
			}`

		var resp *Response
		if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
			assert.Fail(err.Error())
			return
		}

		assert.Equal("chatcmpl-6ptKyqKOGXZT6iQnqiXAH8adNLUzD", resp.ID)
		assert.Equal("gpt-3.5-turbo-0301", resp.Model)
		assert.Equal("chat.completion.chunk", resp.Object)
		assert.Equal(time.Unix(1677825464, 0), resp.Created)

		assert.Len(resp.Choices, 1)
		assert.Equal(0, resp.Choices[0].Index)
		assert.Equal(Role(""), resp.Choices[0].Delta.Role)
		assert.Equal("2", resp.Choices[0].Delta.Content)
		assert.Nil(resp.Choices[0].Message)
		assert.Nil(resp.Choices[0].FinishReason)
	}

	{
		jsonStr := `
			{
			  "choices": [
			    {
			      "delta": {},
			      "finish_reason": "stop",
			      "index": 0
			    }
			  ],
			  "created": 1677825464,
			  "id": "chatcmpl-6ptKyqKOGXZT6iQnqiXAH8adNLUzD",
			  "model": "gpt-3.5-turbo-0301",
			  "object": "chat.completion.chunk"
			}`

		var resp *Response
		if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
			assert.Fail(err.Error())
			return
		}

		assert.Equal("chatcmpl-6ptKyqKOGXZT6iQnqiXAH8adNLUzD", resp.ID)
		assert.Equal("gpt-3.5-turbo-0301", resp.Model)
		assert.Equal("chat.completion.chunk", resp.Object)
		assert.Equal(time.Unix(1677825464, 0), resp.Created)

		assert.Len(resp.Choices, 1)
		assert.Equal(0, resp.Choices[0].Index)
		assert.Equal(Role(""), resp.Choices[0].Delta.Role)
		assert.Equal("", resp.Choices[0].Delta.Content)
		assert.Nil(resp.Choices[0].Message)
		assert.NotNil(resp.Choices[0].FinishReason)
		assert.Equal(Stop, *resp.Choices[0].FinishReason)
	}
}
