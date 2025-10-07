package main

import (
	"context"
	"github.com/K1la/image-processor/internal/api/handler"
	"github.com/K1la/image-processor/internal/api/router"
	"github.com/K1la/image-processor/internal/api/server"
	"github.com/K1la/image-processor/internal/config"
	"github.com/K1la/image-processor/internal/kafka"
	"github.com/K1la/image-processor/internal/repository/minio"
	"github.com/K1la/image-processor/internal/repository/postgres"
	"github.com/K1la/image-processor/internal/service"
	"github.com/wb-go/wbf/zlog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	zlog.Init()

	cfg := config.Init()

	db := postgres.NewDB(cfg)
	repo := postgres.New(db)
	que := kafka.New(cfg.Kafka)
	fileStor := minio.New(cfg.Minio)

	srvc := service.New(repo, fileStor, que)

	hndlr := handler.New(srvc)
	r := router.New(hndlr)
	s := server.New(cfg.HTTPServer.Address, r)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// sig channel to handle SIGINT and SIGTERM for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	srvc.StartWorkers(ctx)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Msgf("recieved shutting down signal %v. Shutting down...", sig)
		cancel()
	}()

	if err := s.ListenAndServe(); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to start server")
	}
	zlog.Logger.Info().Msg("successfully started server on " + cfg.HTTPServer.Address)
}
