package service

import "context"

func (s *Service) DeleteItem(ctx context.Context, id int) error {
	return s.storage.DeleteItem(ctx, id)
}
