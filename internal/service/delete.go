package service

import (
	"context"
	"github.com/google/uuid"
)

func (s *Service) DeleteImage(ctx context.Context, id uuid.UUID) error {
	images := []string{
		id.String() + ".jpg",
		id.String() + ".jpeg",
		id.String() + ".png",
		id.String() + ".gif",
	}

	if err := s.db.DeleteImage(ctx, id); err != nil {
		return err
	}

	if err := s.file.DeleteImages(Images, images...); err != nil {
		return err
	}

	if err := s.file.DeleteImages(Processed, images...); err != nil {
		return err
	}

	return nil
}
