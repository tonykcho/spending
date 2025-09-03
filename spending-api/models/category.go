package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	Id        int
	UUId      uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
	DeletedAt time.Time
}
