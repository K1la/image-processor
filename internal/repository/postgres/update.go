package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (p *Postgres) UpdateImageStatus(ctx context.Context, id uuid.UUID, newStatus string) error {
	query := `
	UPDATE images
	SET status = $1
	WHERE id = $2	`

	result, err := p.db.ExecContext(ctx, query, newStatus, id)
	if err != nil {
		return fmt.Errorf("could not update image status: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update image status: %w", err)
	}

	if affected == 0 {
		return ErrNoSuchImage
	}

	return nil
}
