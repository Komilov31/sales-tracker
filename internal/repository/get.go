package repository

import (
	"context"
	"fmt"

	"github.com/Komilov31/sales-tracker/internal/dto"
	"github.com/Komilov31/sales-tracker/internal/model"
)

func (r *Repository) GetAllItems(ctx context.Context, params dto.GetItemsParams) ([]model.Item, error) {
	query := "SELECT * FROM items"

	orderBy := prepareParams(params)

	rows, err := r.db.Master.QueryContext(ctx, query+orderBy)
	if err != nil {
		return nil, fmt.Errorf("could not get items from db: %w", err)
	}

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.Amount,
			&item.Date,
			&item.Category,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan row result to model: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}
