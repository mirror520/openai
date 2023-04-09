package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/mirror520/openai"
	"github.com/mirror520/openai/conf"
	"github.com/mirror520/openai/persistent/inmem"
	"github.com/mirror520/openai/transport/http"
)

func run() error {
	f, err := os.Open("../config.yaml")
	if err != nil {
		return err
	}
	defer f.Close()

	var cfg *conf.Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return err
	}

	log, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer log.Sync()

	zap.ReplaceGlobals(log)

	repo := inmem.NewChatRepository()
	defer repo.Close()

	r := gin.Default()
	r.Use(cors.Default())

	svc := openai.NewService(repo, cfg)
	svc = openai.LoggingMiddleware(log)(svc)

	apiV1 := r.Group("/openai/v1")
	{
		// POST /chats
		{
			endpoint := openai.CreateChatEndpoint(svc)
			apiV1.POST("/chats", http.CreateChatHandler(endpoint))
		}

		// PATCH /chats/:id
		{
			endpoint := openai.UpdateChatEndpoint(svc)
			apiV1.PATCH("/chats/:id", http.UpdateChatHandler(endpoint))
		}

		// PUT /chats/:id/ask
		{
			endpoint := openai.ChatEndpoint(svc)
			apiV1.PUT("/chats/:id/ask", http.ChatHandler(endpoint))
		}
	}

	// TODO: Service Registration

	return nil
}
