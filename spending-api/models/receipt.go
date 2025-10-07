package models

import (
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	Id        int
	UUId      uuid.UUID
	StoreName string
	Total     float64
	Date      time.Time
	IsDeleted bool
	DeletedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	Items []*ReceiptItem
}

func NewReceipt(storeName string, total float64, date time.Time) *Receipt {
	return &Receipt{
		UUId:      uuid.New(),
		StoreName: storeName,
		Total:     total,
		Date:      date,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
