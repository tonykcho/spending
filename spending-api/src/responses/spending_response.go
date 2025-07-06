package responses

import (
	"time"

	"github.com/google/uuid"
)

type SpendingResponse struct {
	UUId         uuid.UUID
	Amount       float32
	Remark       string
	SpendingDate time.Time
	Category     string
}
