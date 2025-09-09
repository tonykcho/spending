package dto

import (
	"time"

	"github.com/google/uuid"
)

type CategoryDto struct {
	UUId      uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Stores    []*StoreDto
}
