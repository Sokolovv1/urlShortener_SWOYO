package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"urlShortener/internal/app"
	"urlShortener/internal/initialize"
)

func main() {
	use := flag.Bool("d", false, "use local storage")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	config, error := initialize.Load()
	if error != nil {
		logger.Error("Failed to initialize config", zap.Error(error))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, config, logger, *use); err != nil {
		logger.Error("Failed to run server", zap.Error(err))
	}
}
