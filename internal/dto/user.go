package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
	"izanr.com/chat/internal/utils"
)

type UserRole string

const (
	UserRoleNone      UserRole = "NONE"
	UserRoleDefault   UserRole = "DEFAULT"
	UserRoleModerator UserRole = "MODERATOR"
	UserRoleAdmin     UserRole = "ADMIN"
)

const BcryptCost = 12

type User struct {
	ID uuid.UUID `json:"id" db:"id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Username    string            `json:"username" db:"username"`
	DisplayName mo.Option[string] `json:"name" db:"name"`
	Email       mo.Option[string] `json:"email" db:"email"`

	Pronouns  mo.Option[string] `json:"pronouns" db:"pronouns"`
	PictureID mo.Option[int64]  `json:"picture_id" db:"picture_id"`
	Bio       mo.Option[string] `json:"bio" db:"bio"`

	Role         UserRole `json:"role" db:"role"`
	PasswordHash []byte   `json:"-" db:"password"`
}

func (u *User) PasswordMatches(passwd string) bool {
	return utils.CheckPasswordHash(u.PasswordHash, passwd)
}

type PartialUser struct {
	ID uuid.UUID `json:"id" db:"id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Username    string            `json:"username" db:"username"`
	DisplayName mo.Option[string] `json:"name" db:"name"`
	Picture     mo.Option[int32]  `json:"picture" db:"picture"`

	Role UserRole `json:"role" db:"role"`
}

type UserCreateData struct {
	//govalid:required
	//govalid:maxlength=16
	//govalid:minlength=3
	Username string `json:"username"`
	//govalid:maxlength=32
	//govalid:minlength=3
	DisplayName mo.Option[string] `json:"name"`
	//govalid:email
	//govalid:maxlength=128
	Email mo.Option[string] `json:"email"`
	//govalid:maxlength=32
	Pronouns mo.Option[string] `json:"pronouns"`
	//govalid:maxlength=128
	Bio mo.Option[string] `json:"bio"`
	//govalid:required
	//govalid:minlength=3
	Password string `json:"password"`
}

type UserUpdateData struct {
	//govalid:maxlength=32
	Pronouns mo.Option[string] `json:"pronouns"`
	//govalid:maxlength=32
	//govalid:minlength=3
	DisplayName mo.Option[string] `json:"name"`
	//govalid:email
	//govalid:maxlength=128
	Email mo.Option[string] `json:"email"`
	//govalid:maxlength=128
	Bio mo.Option[string] `json:"bio"`
}

type UserUpdatePasswordData struct {
	//govalid:required
	//govalid:minlength=3
	OldPassword string `json:"old_password"`
	//govalid:required
	//govalid:minlength=3
	NewPassword string `json:"new_password"`
}
