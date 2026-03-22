package main

import (
	"capturecraft-api/internal/auth"
	"capturecraft-api/internal/config"
	"capturecraft-api/internal/handlers"
	"capturecraft-api/internal/server"
	"capturecraft-api/internal/storage/memory"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	store := memory.New()
	authSvc := auth.NewService(cfg.JWTSecret, cfg.TokenTTL, store)
	handler := handlers.New(store, authSvc)
	app := server.New(handler, authSvc)

	if cfg.Mode == "lambda" {
		adapter := fiberadapter.New(app)
		lambda.Start(adapter.ProxyWithContext)
		return
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting capturecraft API on %s (mode=%s)", addr, cfg.Mode)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
