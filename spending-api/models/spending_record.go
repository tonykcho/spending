package models

import (
	"time"

	"github.com/google/uuid"
)

type SpendingRecord struct {
	Id           int
	UUId         uuid.UUID
	Amount       float32
	Remark       string
	SpendingDate time.Time
	CategoryId   int
	Category     *Category
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsDeleted    bool
	DeletedAt    time.Time
}

func NewSpendingRecord(amount float32, remark string, spendingDate time.Time, categoryId int) *SpendingRecord {
	return &SpendingRecord{
		UUId:         uuid.New(),
		Amount:       amount,
		Remark:       remark,
		SpendingDate: spendingDate,
		CategoryId:   categoryId,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
}
