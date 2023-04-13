package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/mirror520/openai"
	"github.com/mirror520/openai/conf"
	"github.com/mirror520/openai/persistent/inmem"
	"github.com/mirror520/openai/transport/http"
)

func main() {
	app := &cli.App{
		Name:  "openai",
		Usage: "OpenAI proxy service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Usage:   "work directory",
				EnvVars: []string{"OPENAI_PATH"},
			},
			&cli.IntFlag{
				Name:    "port",
				Usage:   "service port",
				Value:   8080,
				EnvVars: []string{"OPENAI_PORT"},
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cli *cli.Context) error {
	path := cli.String("path")
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		path = homeDir + "/.openai"
	}

	f, err := os.Open(path + "/config.yaml")
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

	// service
	svc := openai.NewService(repo, cfg)
	svc = openai.LoggingMiddleware(log)(svc)

	// endpoint
	endpoints := &openai.ChatEndpoints{
		CreateChatEndpoint: openai.CreateChatEndpoint(svc),
		UpdateChatEndpoint: openai.UpdateChatEndpoint(svc),
		ChatEndpoint:       openai.ChatEndpoint(svc),
		ChatStreamEndpoint: openai.ChatStreamEndpoint(svc),
	}

	// transport
	r := gin.Default()
	r.Use(cors.Default())
	http.Router(r.Group("/openai/v1"), endpoints)

	port := cli.Int("port")
	go r.Run(":" + strconv.Itoa(port))

	// TODO: Service Registration

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit
	log.Info(sign.String())

	return nil
}
