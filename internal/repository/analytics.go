package repository

import (
	"context"
	"fmt"

	"github.com/Komilov31/sales-tracker/internal/model"
)

func (r *Repository) GetAggregated(ctx context.Context, from, to string) ([]model.Item, error) {
	query := `SELECT *,
       SUM(CASE WHEN type = 'расход' THEN -amount ELSE amount END) OVER () AS total_sum,
       AVG(CASE WHEN type = 'расход' THEN -amount ELSE amount END) OVER () AS avg_amount,
       COUNT(*) OVER () AS total_count,
       percentile_cont_window(array_agg(CASE WHEN type = 'расход' THEN -amount ELSE amount END) OVER (), 0.5) AS median_amount,
       percentile_cont_window(array_agg(CASE WHEN type = 'расход' THEN -amount ELSE amount END) OVER (), 0.9) AS percentile_90_amount
	   FROM items
	   WHERE (($1 = '' AND $2 = '') OR (date BETWEEN $1::date AND $2::date));`

	var items []model.Item
	rows, err := r.db.Master.QueryContext(
		ctx,
		query,
		from,
		to,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get aggregated data: %w", err)
	}

	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.Amount,
			&item.Date,
			&item.Category,
			&item.CreatedAt,
			&item.Aggregated.Sum,
			&item.Aggregated.Agvarage,
			&item.Aggregated.Count,
			&item.Aggregated.Median,
			&item.Aggregated.Percentile_90,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan aggregated data to model: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}
