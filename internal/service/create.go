package service

import (
	"context"

	"github.com/Komilov31/sales-tracker/internal/dto"
	"github.com/Komilov31/sales-tracker/internal/model"
)

func (r *Service) CreateItem(ctx context.Context, item dto.CreateItem) (*model.Item, error) {
	return r.storage.CreateItem(ctx, item)
}
