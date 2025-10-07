package main

import (
	"github.com/K1la/image-processor/internal/api/handler"
	"github.com/K1la/image-processor/internal/api/router"
	"github.com/K1la/image-processor/internal/api/server"
	"github.com/K1la/image-processor/internal/config"
	"github.com/K1la/image-processor/internal/repository"
	"github.com/K1la/image-processor/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/wb-go/wbf/zlog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	zlog.Init()

	cfg := config.Init()

	db := repository.NewDB(cfg)
	repo := repository.New(db)
	srvc := service.New(repo)

	// sig channel to handle SIGINT and SIGTERM for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	hndlr := handler.New(srvc)
	r := router.New(hndlr)
	s := server.New(cfg.HTTPServer.Address, r)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to start server")
		}
		zlog.Logger.Info().Msg("successfully started server on " + cfg.HTTPServer.Address)
	}()

	// Блокируем main горутину до получения сигнала завершения
	<-sigChan
	zlog.Logger.Info().Msg("Shutting down gracefully...")
}
