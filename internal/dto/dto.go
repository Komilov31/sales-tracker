package dto

import "time"

type CreateItem struct {
	Type     string `json:"type" validate:"required,oneof=доход расход"`
	Amount   int    `json:"amount" validate:"required,gte=0"`
	Date     string `json:"date" validate:"required,datetime=2006-01-02"`
	Category string `json:"category" validate:"required"`
}

type ItemWithoutAggregated struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Amount    int       `json:"amount"`
	Date      string    `json:"date"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateItem struct {
	Type     *string `json:"type"`
	Amount   *int    `json:"amount"`
	Date     *string `json:"date"`
	Category *string `json:"category"`
}

type GetItemsParams struct {
	SortBy []string
}
