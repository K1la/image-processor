package service

import (
	"context"
	"fmt"
	"github.com/wb-go/wbf/zlog"
	"os"
)

const (
	workersNum = 3
)

func (s *Service) StartWorkers(ctx context.Context) {
	for i := range workersNum {
		zlog.Logger.Info().Msgf("starting worker with index: %d", i)
		go s.worker(ctx)
	}
}

func (s *Service) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			zlog.Logger.Info().Msg("received signal to stop workers...")
			return
		default:
			if err := s.handleMessage(ctx); err != nil {
				zlog.Logger.Warn().Err(err).Msg("failed to handle message")
				continue
			}
			zlog.Logger.Info().Msg("successfully handled message ans save it in file db")
		}
	}
}

func (s *Service) handleMessage(ctx context.Context) error {
	message, err := s.queue.ConsumeMessage()
	if err != nil {
		return fmt.Errorf("failed to consume message from queue: %w", err)
	}

	origPath := s.originDirName + "/" + message.FileName
	procPath := s.processedDirName + "/" + message.FileName
	file, err := os.Create(procPath)
	if err != nil {
		return fmt.Errorf("failed to create file to save image from storage: %w", err)
	}
	file.Close()

	// TODO: проверить с указателем
	if err = s.ProcessImage(message); err != nil {
		return fmt.Errorf("could not process image: %w", err)
	}

	if err = s.file.SaveImage(message.FileName, procPath, Processed); err != nil {
		return fmt.Errorf("could not save processed message to fileStorage: %s", err.Error())
	}

	if err = s.db.UpdateImageStatus(ctx, message.ID, "finished"); err != nil {
		return fmt.Errorf("could not update image processing status in db: %s", err.Error())
	}

	if err = os.Remove(procPath); err != nil {
		return fmt.Errorf("could not delete temp file: %s", err.Error())
	}

	if err = os.Remove(origPath); err != nil {
		return fmt.Errorf("could not delete temp file: %s", err.Error())
	}

	return nil
}
