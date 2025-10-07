package dto

import (
	"time"

	"github.com/google/uuid"
)

type ReceiptItemDto struct {
	Id        uuid.UUID
	Name      string
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
