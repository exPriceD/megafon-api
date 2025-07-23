package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"megafon-buisness-reports/internal/app"
	"megafon-buisness-reports/internal/config"
)

func main() {
	cfg := config.MustLoad("config/config.dev.yaml", ".env")

	application := app.New(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := application.Run(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
