package repository

import (
	"context"
	"fmt"

	"github.com/Komilov31/sales-tracker/internal/dto"
)

func (r *Repository) UpdateItem(ctx context.Context, id int, item dto.UpdateItem) error {
	query := `UPDATE items
	SET type = COALESCE($1, type),
		amount = COALESCE($2, amount),
		date = COALESCE($3, date),
		category = COALESCE($4, category)
	WHERE id = $5`

	result, err := r.db.Master.ExecContext(
		ctx,
		query,
		item.Type,
		item.Amount,
		item.Date,
		item.Category,
		id,
	)
	if err != nil {
		return fmt.Errorf("could not update item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update item: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNoSuchItem
	}

	return nil
}
