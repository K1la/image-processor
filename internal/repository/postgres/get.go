package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/K1la/image-processor/internal/model"
	"github.com/google/uuid"
)

func (p *Postgres) GetImageInfo(ctx context.Context, id uuid.UUID) (*model.Image, error) {
	query := `
	SELECT * 
	FROM images 
	WHERE id = $1`

	var image model.Image
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&image.ID,
		&image.Format,
		&image.Status,
		&image.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoSuchImage
		}

		return nil, fmt.Errorf("could not get image from db: %w", err)
	}

	return &image, nil
}
