package models

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	Id         int
	UUId       uuid.UUID
	Name       string
	CategoryId int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsDeleted  bool
	DeletedAt  time.Time
}

func NewStore(name string, categoryId int) *Store {
	return &Store{
		UUId:       uuid.New(),
		Name:       name,
		CategoryId: categoryId,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}
