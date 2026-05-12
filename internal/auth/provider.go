package auth

import (
	"context"

	"github.com/google/uuid"
)

type Provider interface {
	GenerateToken(context.Context, Token) (string, error)
	ValidateToken(context.Context, string) (Token, error)

	GenerateRefresh(context.Context, uuid.UUID) (string, error)
	ValidateRefresh(context.Context, string) (uuid.UUID, error)
}
