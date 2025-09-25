package handler

import (
	"context"

	"github.com/Komilov31/sales-tracker/internal/dto"
	"github.com/Komilov31/sales-tracker/internal/model"
)

type TrackerService interface {
	CreateItem(ctx context.Context, item dto.CreateItem) (*model.Item, error)
	GetAllItems(ctx context.Context, params dto.GetItemsParams) ([]model.Item, error)
	GetAggregated(ctx context.Context, from, to string) ([]model.Item, error)
	UpdateItem(ctx context.Context, id int, item dto.UpdateItem) error
	DeleteItem(ctx context.Context, id int) error
	CSVAggregated(ctx context.Context, from, to string) (string, error)
	CSVAllItems(ctx context.Context, params dto.GetItemsParams) (string, error)
}

type Handler struct {
	service TrackerService
	ctx     context.Context
}

func New(ctx context.Context, service TrackerService) *Handler {
	return &Handler{
		service: service,
		ctx:     ctx,
	}
}
