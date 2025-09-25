package service

import (
	"context"
	"log"
	"os"

	"github.com/Komilov31/sales-tracker/internal/dto"
	"github.com/Komilov31/sales-tracker/internal/model"
)

type Storage interface {
	CreateItem(ctx context.Context, item dto.CreateItem) (*model.Item, error)
	GetAllItems(ctx context.Context, params dto.GetItemsParams) ([]model.Item, error)
	UpdateItem(ctx context.Context, id int, item dto.UpdateItem) error
	DeleteItem(ctx context.Context, id int) error
	GetAggregated(ctx context.Context, from, to string) ([]model.Item, error)
}

type Service struct {
	storage    Storage
	folderName string
}

func New(storage Storage) *Service {
	folderName, err := os.MkdirTemp(".", "csv")
	if err != nil {
		log.Fatal("could not create folder to store csv files: ", err)
	}

	return &Service{
		storage:    storage,
		folderName: folderName,
	}
}
