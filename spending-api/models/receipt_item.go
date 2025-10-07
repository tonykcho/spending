package models

import (
	"time"

	"github.com/google/uuid"
)

type ReceiptItem struct {
	Id        int
	UUId      uuid.UUID
	ReceiptId int
	Name      string
	Price     float64
	IsDeleted bool
	DeletedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewReceiptItem(receiptId int, name string, price float64) *ReceiptItem {
	return &ReceiptItem{
		UUId:      uuid.New(),
		ReceiptId: receiptId,
		Name:      name,
		Price:     price,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
