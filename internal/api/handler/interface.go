package handler

import (
	"context"
	"github.com/K1la/image-processor/internal/model"
	"github.com/google/uuid"
)

type ServiceI interface {
	ProcessImage(model.Message) error
	CreateImage(context.Context, []byte, model.Message) (*uuid.UUID, error)
	GetImageStatus(context.Context, uuid.UUID) (*model.Image, error)
	GetImageById(context.Context, uuid.UUID) (string, error)
	DeleteImage(context.Context, uuid.UUID) error
}
