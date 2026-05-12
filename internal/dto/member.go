package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

type MemberFlags int32

const (
	MemberFlagReadMessage MemberFlags = 1 << iota
	MemberFlagSendMessage
	MemberFlagInvite
	MemberFlagAttachImage
	MemberFlagDeleteMessages
	MemberFlagBanMembers
	MemberFlagAdmin

	MemberFlagsAll     = 1<<iota - 1
	MemberFlagsDefault = MemberFlagReadMessage |
		MemberFlagSendMessage |
		MemberFlagInvite |
		MemberFlagAttachImage
)

func (r MemberFlags) Can(flags MemberFlags) bool {
	return r&flags == flags || r&MemberFlagAdmin == MemberFlagAdmin
}

type Member struct {
	ChannelID uuid.UUID `json:"channel_id" db:"channel_id"`
	UserId    uuid.UUID `json:"user_id" db:"user_id"`

	AddedAt   time.Time `json:"added_at" db:"added_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Nickname mo.Option[string] `json:"nickname" db:"nickname"`
	Flags    MemberFlags       `json:"role" db:"role"`
}

type MemberCreateData struct {
	Nickname mo.Option[string] `json:"nickname" db:"nickname"`
	Role     MemberFlags       `json:"role" db:"role"`
}

type UserInChannel struct {
	Member `db:"member"`
	User   User `json:"user" db:"user"`
}
