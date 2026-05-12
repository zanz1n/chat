package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"izanr.com/chat/internal/dto"
	"izanr.com/chat/internal/utils"
)

var (
	ErrUserNotFound         = fmt.Errorf("user not found")
	ErrUserAlreadyExists    = fmt.Errorf("user already exists")
	ErrUserPasswdMismatches = fmt.Errorf("user password mismatches")
)

type UserStorer interface {
	Insert(context.Context, dto.UserCreateData) (dto.User, error)

	GetById(context.Context, uuid.UUID) (dto.User, error)
	GetByUsername(context.Context, string) (dto.User, error)
	GetMany(context.Context, dto.Pagination) ([]dto.User, error)

	UpdateUsername(context.Context, uuid.UUID, string) (dto.User, error)
	Update(context.Context, uuid.UUID, dto.UserUpdateData) (dto.User, error)
	UpdatePicture(
		ctx context.Context,
		userId uuid.UUID,
		pictureId uuid.UUID,
	) error
	UpdatePassword(context.Context, uuid.UUID, dto.UserUpdatePasswordData) (dto.User, error)

	Delete(context.Context, uuid.UUID) (dto.User, error)
}

type users struct {
	q utils.Querier
}

func NewPgUsers(q utils.Querier) UserStorer {
	return &users{q: q}
}

func (r *users) Insert(ctx context.Context, data dto.UserCreateData) (dto.User, error) {
	password := utils.HashPassword(data.Password)

	u, err := r.selectOne(ctx, query_users_insert,
		data.Username,
		data.DisplayName,
		data.Email,
		data.Pronouns,
		dto.UserRoleDefault,
		password,
	)
	if err != nil {
		if isConflict(err) {
			err = ErrUserAlreadyExists
		}
	}
	return u, err
}

func (r *users) GetById(ctx context.Context, id uuid.UUID) (dto.User, error) {
	return r.selectOne(ctx, query_users_get_by_id, id)
}

func (r *users) GetByUsername(ctx context.Context, username string) (dto.User, error) {
	return r.selectOne(ctx, query_users_get_by_username, username)
}

func (r *users) GetMany(ctx context.Context, pag dto.Pagination) ([]dto.User, error) {
	rows, err := r.q.Query(ctx, query_users_get_many, pag.LastSeen, pag.Limit)
	if err != nil {
		return nil, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.User])
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *users) UpdateUsername(
	ctx context.Context,
	id uuid.UUID,
	username string,
) (dto.User, error) {
	u, err := r.selectOne(ctx, query_users_update_username, id, username)
	if err != nil {
		if isConflict(err) {
			err = ErrUserAlreadyExists
		}
	}
	return u, err
}

func (r *users) Update(
	ctx context.Context,
	id uuid.UUID,
	data dto.UserUpdateData,
) (dto.User, error) {
	// no conflicts possible in this one
	return r.selectOne(ctx, query_users_update,
		id,
		data.Pronouns,
		data.DisplayName,
		data.Email,
		data.Bio,
	)
}

func (r *users) UpdatePicture(
	ctx context.Context,
	userId uuid.UUID,
	pictureId uuid.UUID,
) error {
	_, err := r.q.Exec(ctx, query_users_update_picture, userId, pictureId)
	return err
}

func (r *users) UpdatePassword(
	ctx context.Context,
	id uuid.UUID,
	data dto.UserUpdatePasswordData,
) (dto.User, error) {
	user, err := r.GetById(ctx, id)
	if err != nil {
		return dto.User{}, err
	}

	if !user.PasswordMatches(data.OldPassword) {
		return dto.User{}, ErrUserPasswdMismatches
	}
	newHash := utils.HashPassword(data.NewPassword)

	return r.selectOne(ctx, query_users_update_password, id, newHash)
}

func (r *users) Delete(ctx context.Context, id uuid.UUID) (dto.User, error) {
	return r.selectOne(ctx, query_users_delete, id)
}

func (r *users) selectOne(ctx context.Context, query string, args ...any) (dto.User, error) {
	row, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return dto.User{}, err
	}

	user, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[dto.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.User{}, ErrUserNotFound
		}
		return dto.User{}, err
	}
	return user, nil
}
