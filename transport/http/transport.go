package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"

	"github.com/mirror520/openai"
	"github.com/mirror520/openai/chat"
	"github.com/mirror520/openai/model"
)

func Router(route *gin.RouterGroup, endpoints *openai.ChatEndpoints) {
	// POST /chats
	route.POST("/chats", CreateChatHandler(endpoints.CreateChatEndpoint))

	// PATCH /chats/:id
	route.PATCH("/chats/:id", UpdateChatHandler(endpoints.UpdateChatEndpoint))

	// POST /chats/:id/messages
	route.POST("/chats/:id/messages", ChatHandler(
		endpoints.ChatEndpoint,
		endpoints.ChatStreamEndpoint,
	))
}

func CreateChatHandler(endpoint endpoint.Endpoint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req *openai.CreateChatRequest
		if err := ctx.ShouldBind(&req); err != nil {
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, result)
			return
		}

		resp, err := endpoint(ctx, req)
		if err != nil {
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
			return
		}

		result := model.SuccessResult("chat created")
		result.Data = resp
		ctx.JSON(http.StatusOK, result)
	}
}

func UpdateChatHandler(endpoint endpoint.Endpoint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req *openai.UpdateChatRequest
		if err := ctx.ShouldBind(&req); err != nil {
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, result)
			return
		}

		id, err := chat.ParseID(ctx.Param("id"))
		if err != nil {
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, result)
			return
		}

		req.ID = id

		resp, err := endpoint(ctx, req)
		if err != nil {
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
			return
		}

		result := model.SuccessResult("chat updated")
		result.Data = resp
		ctx.JSON(http.StatusOK, result)
	}
}

func ChatHandler(chatEndpoint endpoint.Endpoint, chatStreamEndpoint endpoint.Endpoint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req *openai.ChatRequest
		if err := ctx.ShouldBind(&req); err != nil {
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, result)
			return
		}

		id, err := chat.ParseID(ctx.Param("id"))
		if err != nil {
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, result)
			return
		}

		req.ID = id

		stream := false
		if s := ctx.Query("stream"); s != "" {
			b, err := strconv.ParseBool(s)
			if err == nil {
				stream = b
			}
		}

		if !stream {
			resp, err := chatEndpoint(ctx, req)
			if err != nil {
				result := model.FailureResult(err)
				ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
				return
			}

			result := model.SuccessResult("chat answered")
			result.Data = resp
			ctx.JSON(http.StatusOK, result)

		} else {
			resp, err := chatStreamEndpoint(ctx, req)
			if err != nil {
				result := model.FailureResult(err)
				ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
				return
			}

			stream, ok := resp.(<-chan string)
			if !ok {
				err := errors.New("invalid stream")
				result := model.FailureResult(err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, result)
				return
			}

			w := ctx.Writer
			w.WriteHeader(http.StatusOK)

			for content := range stream {
				w.WriteString(content)
				w.Flush()
			}
		}
	}
}
