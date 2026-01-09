package storage

import (
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/db"
)

// Storage is a storage adapter for a database.
type Storage struct {
	db db.DB
}

// NewStorage creates a new DB storage.
func NewStorage(db db.DB) *Storage {
	return &Storage{
		db: db,
	}
}
