package service

import (
	"context"
	"github.com/K1la/image-processor/internal/model"
	"github.com/google/uuid"
)

type DBRepo interface {
	CreateImage(context.Context, model.Image) error
	GetImageInfo(context.Context, uuid.UUID) (*model.Image, error)
	UpdateImageStatus(context.Context, uuid.UUID, string) error
	DeleteImage(context.Context, uuid.UUID) error
}

type FileRepo interface {
	SaveImage(string, string, string) error
	GetImage(string, string, string) error
	DeleteImages(string, ...string) error
}

type Queue interface {
	ProduceMessage(model.Message) error
	ConsumeMessage() (model.Message, error)
}
