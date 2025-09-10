package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	Id        int
	UUId      uuid.UUID
	Name      string
	Stores    []*Store
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
	DeletedAt time.Time
}

func NewCategory(name string) *Category {
	return &Category{
		UUId:      uuid.New(),
		Name:      name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
