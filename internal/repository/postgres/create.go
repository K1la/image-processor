package postgres

import (
	"context"
	"fmt"
	"github.com/K1la/image-processor/internal/model"
)

func (p *Postgres) CreateImage(ctx context.Context, image model.Image) error {
	query := `
	INSERT INTO 
	images(id, format, status)
	VALUES ($1, $2, $3)`

	_, err := p.db.ExecContext(ctx, query, image.ID, image.Format, image.Status)
	if err != nil {
		return fmt.Errorf("could not save image info in db: %w", err)
	}

	return nil
}
