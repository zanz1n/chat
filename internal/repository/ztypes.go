package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/mo"
	"izanr.com/chat/internal/dto"
)

type joinableUser struct {
	ID           uuid.UUID         `db:"user.id"`
	CreatedAt    time.Time         `db:"user.created_at"`
	UpdatedAt    time.Time         `db:"user.updated_at"`
	Username     string            `db:"user.username"`
	DisplayName  mo.Option[string] `db:"user.name"`
	Email        mo.Option[string] `db:"user.email"`
	Pronouns     mo.Option[string] `db:"user.pronouns"`
	PictureID    mo.Option[int32]  `db:"user.picture_id"`
	Bio          mo.Option[string] `db:"user.bio"`
	Role         dto.UserRole      `db:"user.role"`
	PasswordHash []byte            `db:"user.password_hash"`

	dto.User
}

func (m joinableUser) into() dto.User {
	return m.User
}

type joinableMember struct {
	ChannelID uuid.UUID         `db:"member.channel_id"`
	UserId    uuid.UUID         `db:"member.user_id"`
	AddedAt   time.Time         `db:"member.added_at"`
	UpdatedAt time.Time         `db:"member.updated_at"`
	Nickname  mo.Option[string] `db:"member.nickname"`
	Role      dto.MemberFlags   `db:"member.role"`

	dto.Member
}

func (m joinableMember) into() dto.Member {
	return m.Member
}

type joinableUiC struct {
	joinableUser
	joinableMember
}

func (uic joinableUiC) into() dto.UserInChannel {
	return dto.UserInChannel{
		User:   uic.joinableUser.into(),
		Member: uic.joinableMember.into(),
	}
}

func collectUic(row pgx.CollectableRow) (dto.UserInChannel, error) {
	res, err := pgx.RowToStructByName[joinableUiC](row)
	if err != nil {
		return dto.UserInChannel{}, nil
	}
	return res.into(), nil
}
