package model

import (
	"github.com/google/uuid"
	"time"
)

type Image struct {
	ID        uuid.UUID `json:"id"`
	Format    string    `json:"-"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID            uuid.UUID `json:"id"`
	Task          string    `json:"task"`
	FileName      string    `json:"file_name"`
	ContentType   string    `json:"content_type"`
	WatermarkText string    `json:"watermark_string"`
	Resize        struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"resize"`
}
