package repository

import (
	"context"
	"fmt"

	"github.com/Komilov31/sales-tracker/internal/dto"
	"github.com/Komilov31/sales-tracker/internal/model"
)

func (r *Repository) CreateItem(ctx context.Context, item dto.CreateItem) (*model.Item, error) {
	query := `INSERT INTO items(type, amount, date, category)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at;`

	var createdItem model.Item
	err := r.db.Master.QueryRowContext(
		ctx,
		query,
		item.Type,
		item.Amount,
		item.Date,
		item.Category,
	).Scan(&createdItem.ID, &createdItem.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not create item in db: %w", err)
	}

	createdItem.Type = item.Type
	createdItem.Amount = item.Amount
	createdItem.Date = item.Date
	createdItem.Category = item.Category

	return &createdItem, nil
}
