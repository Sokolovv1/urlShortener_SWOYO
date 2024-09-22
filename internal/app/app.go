package app

import (
	"context"
	"go.uber.org/zap"
	"urlShortener/internal/controller"
	"urlShortener/internal/initialize"
	"urlShortener/internal/repository"
	http "urlShortener/internal/server_http"
	"urlShortener/internal/service"
)

func Run(ctx context.Context, config *initialize.Config, logger *zap.Logger, use bool) error {
	var err error
	var pgDb *initialize.DB
	var shortenerRepository service.SwapRepository

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if !use {
		pgDb, err = initialize.NewClient(ctx, config.PGMaxAttemption, config)
		if err != nil {
			logger.Error("error initializing pgDB", zap.Error(err))
			return err
		}
		shortenerRepository, err = repository.NewShortenerRepository(pgDb, logger)

		if err != nil {
			logger.Error("error creating shortener repository", zap.Error(err))
			return err
		}
		logger.Info("successfully connected to pgDB")

	} else {
		shortenerRepository = repository.NewURLStorage(logger)
		logger.Info("initializing shortener repository with local database")
	}

	//shortenerRepository, err := repository.NewShortenerRepository(pgDb.Pool)

	shortenerService := service.NewShortenerService(service.Deps{
		Repository: shortenerRepository,
		Config:     config,
		Logger:     logger,
	})

	shortenerController := controller.NewShortenerController(shortenerService, logger)

	server := http.NewServer(http.ServerConfig{
		Controllers: []http.Controller{shortenerController},
		Logger:      logger,
	})

	go func() {
		if err := server.Start(config.HTTPHost + ":" + config.HTTPPort); err != nil {
			logger.Error("Server startup error", zap.Error(err))
			cancel()
		}
	}()

	server.ListRoutes()

	<-ctx.Done()

	if !use {
		pgDb.Pool.Close()
	}

	server.Shutdown()

	//timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer timeoutCancel()
	logger.Info("shortener service shortener shutdown")
	return nil
}
