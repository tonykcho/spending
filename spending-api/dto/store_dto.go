package dto

import (
	"time"

	"github.com/google/uuid"
)

type StoreDto struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
