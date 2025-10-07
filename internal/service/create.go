package service

import (
	"context"
	"fmt"
	"github.com/K1la/image-processor/internal/model"
	"github.com/google/uuid"
	"os"
)

func (s *Service) CreateImage(ctx context.Context, data []byte, imageData model.Message) (*uuid.UUID, error) {
	format, err := checkFormat(imageData.ContentType)
	if err != nil {
		return nil, err
	}

	if !isCorrectTask(imageData.Task) {
		return nil, ErrInvalidTask
	}

	id := uuid.New()
	image := model.Image{
		ID:     id,
		Format: format,
		Status: InProgress,
	}

	fileName := id.String() + "." + format
	imageData.ID = id
	imageData.FileName = fileName

	filePath := s.originDirName + "/" + fileName
	if err = os.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("could not save image to dir: %w", err)
	}

	if err = s.file.SaveImage(fileName, filePath, Images); err != nil {
		return nil, err
	}

	if err = s.queue.ProduceMessage(imageData); err != nil {
		return nil, err
	}

	if err = s.db.CreateImage(ctx, image); err != nil {
		return nil, err
	}

	return &id, nil
}
