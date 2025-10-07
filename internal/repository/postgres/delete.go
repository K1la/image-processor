package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (p *Postgres) DeleteImage(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM images 
	WHERE id = $1`

	_, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("could not delete image info from db: %w", err)
	}

	return nil
}
