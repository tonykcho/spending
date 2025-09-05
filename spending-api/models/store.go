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
