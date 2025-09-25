package service

import (
	"context"

	"github.com/Komilov31/sales-tracker/internal/dto"
)

func (s *Service) UpdateItem(ctx context.Context, id int, item dto.UpdateItem) error {
	return s.storage.UpdateItem(ctx, id, item)
}
