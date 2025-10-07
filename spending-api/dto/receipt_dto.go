package dto

import (
	"time"

	"github.com/google/uuid"
)

type ReceiptDto struct {
	Id        uuid.UUID
	StoreName string
	Date      time.Time
	Total     float64
	Items     []*ReceiptItemDto

	CreatedAt time.Time
	UpdatedAt time.Time
}
