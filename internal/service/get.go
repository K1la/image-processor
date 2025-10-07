package service

import (
	"context"
	"fmt"
	"github.com/K1la/image-processor/internal/model"
	"github.com/google/uuid"
)

func (s *Service) GetImageById(ctx context.Context, id uuid.UUID) (string, error) {
	imageInfo, err := s.db.GetImageInfo(ctx, id)
	if err != nil {
		return "", err
	}

	if imageInfo.Status == InProgress {
		return "", ErrNotProcessed
	}

	fileName := id.String() + "." + imageInfo.Format
	filePath := s.processedDirName + "/" + fileName
	if err = s.file.GetImage(fileName, filePath, Processed); err != nil {
		return "", fmt.Errorf("could not get image from file repo: %w", err)
	}

	return filePath, nil
}

func (s *Service) GetImageStatus(ctx context.Context, id uuid.UUID) (*model.Image, error) {
	return s.db.GetImageInfo(ctx, id)
}
