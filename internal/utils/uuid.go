package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func NewUUID() uuid.UUID {
	u, err := uuid.NewV7()
	if err != nil {
		panic(fmt.Errorf("generate UUID failed: %w", err))
	}

	return u
}
