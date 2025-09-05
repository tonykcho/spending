package dto

import (
	"time"

	"github.com/google/uuid"
)

type StoreDto struct {
	UUId      uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
