package dto

import (
	"time"

	"github.com/google/uuid"
)

type SpendingDto struct {
	UUId         uuid.UUID
	Amount       float32
	Remark       string
	SpendingDate time.Time
	Category     *CategoryDto
}
