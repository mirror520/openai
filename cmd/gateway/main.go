package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mirror520/openai"
	"github.com/mirror520/openai/transport/http"
)

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer log.Sync()

	zap.ReplaceGlobals(log)

	// TODO: Service Discovery

	proxyEndpoints := new(openai.ChatEndpoints)
	{
		factory := http.ChatFactory(http.CreateChatEndpoint, "http")
		endpoint, _, _ := factory("127.0.0.1:80")
		proxyEndpoints.CreateChatEndpoint = endpoint
	}

	// UpdateChat
	{
		factory := http.ChatFactory(http.UpdateChatEndpoint, "http")
		endpoint, _, _ := factory("127.0.0.1:80")
		proxyEndpoints.UpdateChatEndpoint = endpoint
	}

	// Chat
	{
		factory := http.ChatFactory(http.ChatEndpoint, "http")
		endpoint, _, _ := factory("127.0.0.1:80")
		proxyEndpoints.ChatEndpoint = endpoint
	}

	// ChatStream
	{
		factory := http.ChatFactory(http.ChatStreamEndpoint, "http")
		endpoint, _, _ := factory("127.0.0.1:80")
		proxyEndpoints.ChatStreamEndpoint = endpoint
	}

	// service (internal use)
	var svc openai.Service // dummy service
	svc = openai.ProxyingMiddleware(proxyEndpoints)(svc)
	svc = openai.LoggingMiddleware(log)(svc)

	// ---

	// endpoint
	endpoints := &openai.ChatEndpoints{
		CreateChatEndpoint: openai.CreateChatEndpoint(svc),
		UpdateChatEndpoint: openai.UpdateChatEndpoint(svc),
		ChatEndpoint:       openai.ChatEndpoint(svc),
		ChatStreamEndpoint: openai.ChatStreamEndpoint(svc),
	}

	// transport (external use)
	r := gin.Default()
	r.Use(cors.Default())
	http.Router(r.Group("/openai/v1"), endpoints)

	r.Run()
}
