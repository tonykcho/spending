package responses

import (
	"time"

	"github.com/google/uuid"
)

type CategoryResponse struct {
	UUId      uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
