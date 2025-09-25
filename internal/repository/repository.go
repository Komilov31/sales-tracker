package repository

import (
	"errors"

	"github.com/wb-go/wbf/dbpg"
)

var (
	ErrNoSuchItem = errors.New("there is no item with such id")
)

type Repository struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Repository {
	return &Repository{
		db: db,
	}
}
