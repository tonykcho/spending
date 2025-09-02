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
}
