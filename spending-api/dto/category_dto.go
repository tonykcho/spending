package dto

import (
	"time"

	"github.com/google/uuid"
)

type CategoryDto struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Stores    []*StoreDto
}
