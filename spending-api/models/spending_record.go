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
	Category     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
