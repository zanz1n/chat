package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

type Channel struct {
	ID      uuid.UUID `json:"id" db:"id"`
	OwnerID uuid.UUID `json:"owner_id" db:"owner_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Name        string            `json:"name" db:"name"`
	Description mo.Option[string] `json:"description" db:"description"`
	PictureID   mo.Option[int64]  `json:"picture_id" db:"picture_id"`
}

type ChannelCreateData struct {
	//govalid:required
	//govalid:maxlength=32
	//govalid:minlength=3
	Name string `json:"name"`
	//govalid:maxlength=128
	Description mo.Option[string] `json:"description"`
}

type ChannelUpdateData = ChannelCreateData

type DirectChannel struct {
	ID int64 `json:"id" db:"id"`

	MinorID uuid.UUID `json:"minor_id" db:"minor_id"`
	MajorID uuid.UUID `json:"major_id" db:"major_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
