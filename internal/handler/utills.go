package handler

import (
	"fmt"
	"time"

	"github.com/Komilov31/sales-tracker/internal/dto"
	"github.com/Komilov31/sales-tracker/internal/model"
)

func validateGetParams(sortBy []string) error {
	fields := map[string]struct{}{
		"type":       {},
		"amount":     {},
		"date":       {},
		"category":   {},
		"created_at": {},
		"id":         {},
	}

	for _, s := range sortBy {
		if _, ok := fields[s]; !ok {
			return fmt.Errorf("invalid field name for sorting")
		}
	}

	return nil
}

func validateDate(from, to string) error {
	if _, err := time.Parse(time.DateOnly, from); err != nil {
		return fmt.Errorf("invalid date format in query parameter")
	}

	if _, err := time.Parse(time.DateOnly, to); err != nil {
		return fmt.Errorf("invalid date format in query parameter, must be in format 'YY-MM-DD'")
	}

	return nil
}

func convertWithoutAggregated(item *model.Item) dto.ItemWithoutAggregated {
	return dto.ItemWithoutAggregated{
		ID: item.ID, Type: item.Type, Amount: item.Amount,
		Date: item.Date, Category: item.Category,
		CreatedAt: item.CreatedAt,
	}
}

func itemsWithoutAggregated(items []model.Item) []dto.ItemWithoutAggregated {
	withoutAggregated := make([]dto.ItemWithoutAggregated, len(items))

	for i, item := range items {
		withoutAggregated[i] = convertWithoutAggregated(&item)
	}

	return withoutAggregated
}
