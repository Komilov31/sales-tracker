package service

import (
	"context"

	"github.com/Komilov31/sales-tracker/internal/dto"
	"github.com/Komilov31/sales-tracker/internal/model"
)

func (s *Service) GetAllItems(ctx context.Context, params dto.GetItemsParams) ([]model.Item, error) {
	return s.storage.GetAllItems(ctx, params)
}

func (s *Service) GetAggregated(ctx context.Context, from, to string) ([]model.Item, error) {
	return s.storage.GetAggregated(ctx, from, to)
}
