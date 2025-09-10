package store_handlers

import (
	"fmt"

	"github.com/google/uuid"
)

type UpdateStoreRequest struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (request UpdateStoreRequest) Valid() error {
	if request.Id == uuid.Nil {
		return fmt.Errorf("id cannot be empty")
	}

	if request.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}
