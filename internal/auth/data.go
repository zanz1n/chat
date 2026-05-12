package auth

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
	"izanr.com/chat/internal/dto"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	*t = Timestamp{time.Unix(i, 0)}
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(t.Unix()))), nil
}

type Token struct {
	UserID    uuid.UUID    `json:"user_id"`
	CreatedAt Timestamp    `json:"iat"`
	ExpiresAt Timestamp    `json:"exp"`
	Issuer    string       `json:"iss"`
	Username  string       `json:"username"`
	Role      dto.UserRole `json:"role"`
}
