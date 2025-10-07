package service

import (
	"errors"
	"log"
	"os"
)

var (
	ErrInvalidImageFormat = errors.New("invalid image format, must be in (jpg, png, gif)")
	ErrInvalidTask        = errors.New("invalid task, must be in(resize, watermark, miniature generating)")
	ErrNotProcessed       = errors.New("image is not ready yet")
)

const (
	Processed  = "processed"
	InProgress = "in progress"
	Images     = "images"
)

type Service struct {
	db               DBRepo
	file             FileRepo
	queue            Queue
	originDirName    string
	processedDirName string
}

func New(d DBRepo, f FileRepo, q Queue) *Service {
	images, err := os.MkdirTemp("./", Images)
	if err != nil {
		log.Fatalf("could not create temporary directory to store images: %v", err)
	}

	processed, err := os.MkdirTemp("./", Processed)
	if err != nil {
		log.Fatalf("could not create temporary directory to store processed images: %v", err)
	}

	return &Service{
		db:               d,
		file:             f,
		queue:            q,
		originDirName:    images,
		processedDirName: processed,
	}
}
