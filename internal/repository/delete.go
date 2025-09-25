package repository

import (
	"context"
	"fmt"
)

func (r *Repository) DeleteItem(ctx context.Context, id int) error {
	query := "DELETE FROM items WHERE id = $1"

	result, err := r.db.Master.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("could not delete item from db: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not delete item from db: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNoSuchItem
	}

	return nil
}
