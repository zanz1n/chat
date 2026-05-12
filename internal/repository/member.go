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
	ErrMemberNotFound = fmt.Errorf("member not found")
	ErrEmptyChannel   = fmt.Errorf("channel is empty")
)

type MemberStorer interface {
	Insert(context.Context, dto.Member) error

	GetById(
		ctx context.Context,
		channelId uuid.UUID,
		userId uuid.UUID,
	) (dto.Member, error)
	GetWithUser(
		ctx context.Context,
		channelId uuid.UUID,
		userId uuid.UUID,
	) (dto.UserInChannel, error)

	GetByChannel(context.Context, uuid.UUID, dto.Pagination) ([]dto.Member, error)
	GetUsersByChannel(
		context.Context,
		uuid.UUID,
		dto.Pagination,
	) ([]dto.UserInChannel, error)

	UpdateNickname(
		ctx context.Context,
		channelId uuid.UUID,
		userId uuid.UUID,
		nickname string,
	) (dto.Member, error)
	UpdateRole(
		ctx context.Context, channelId uuid.UUID,
		userId uuid.UUID,
		role dto.MemberFlags,
	) (dto.Member, error)

	Delete(
		ctx context.Context,
		channelId uuid.UUID,
		userId uuid.UUID,
	) (dto.Member, error)
}

type members struct {
	q utils.Querier
}

func NewPgMembers(q utils.Querier) MemberStorer {
	return &members{q: q}
}

func (r *members) Insert(ctx context.Context, data dto.Member) error {
	_, err := r.q.Exec(ctx, query_members_insert,
		data.ChannelID,
		data.UserId,
		data.AddedAt,
		data.UpdatedAt,
		data.Nickname,
		data.Flags,
	)
	if err != nil {
		// TODO: check for conflict
		return err
	}
	return nil
}

func (r *members) GetById(
	ctx context.Context,
	channelId uuid.UUID,
	userId uuid.UUID,
) (dto.Member, error) {
	return r.selectOne(ctx, query_members_get_by_id, channelId, userId)
}

func (r *members) GetWithUser(
	ctx context.Context,
	channelId uuid.UUID,
	userId uuid.UUID,
) (dto.UserInChannel, error) {
	row, err := r.q.Query(ctx, query_members_get_with_user, channelId, userId)
	if err != nil {
		return dto.UserInChannel{}, err
	}

	uic, err := pgx.CollectExactlyOneRow(row, collectUic)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.UserInChannel{}, ErrMemberNotFound
		}
		return dto.UserInChannel{}, err
	}
	return uic, nil
}

func (r *members) GetByChannel(
	ctx context.Context,
	channelId uuid.UUID,
	pag dto.Pagination,
) ([]dto.Member, error) {
	rows, err := r.q.Query(ctx, query_members_get_by_channel,
		channelId,
		pag.LastSeen,
		pag.Limit,
	)
	if err != nil {
		return nil, err
	}

	members, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.Member])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmptyChannel
		}
		return nil, err
	}
	return members, nil
}

func (r *members) GetUsersByChannel(
	ctx context.Context,
	channelId uuid.UUID,
	pag dto.Pagination,
) ([]dto.UserInChannel, error) {
	row, err := r.q.Query(ctx, query_members_get_users_by_channel,
		channelId,
		pag.LastSeen,
		pag.Limit,
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectRows(row, collectUic)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmptyChannel
		}
		return nil, err
	}
	return res, nil
}

func (r *members) UpdateNickname(
	ctx context.Context,
	channelId uuid.UUID,
	userId uuid.UUID,
	nickname string,
) (dto.Member, error) {
	return r.selectOne(ctx, query_members_update_nickname,
		channelId,
		userId,
		nickname,
	)
}

func (r *members) UpdateRole(
	ctx context.Context,
	channelId uuid.UUID,
	userId uuid.UUID,
	role dto.MemberFlags,
) (dto.Member, error) {
	return r.selectOne(ctx, query_members_update_role,
		channelId,
		userId,
		role,
	)
}

func (r *members) Delete(
	ctx context.Context,
	channelId uuid.UUID,
	userId uuid.UUID,
) (dto.Member, error) {
	return r.selectOne(ctx, query_members_delete, channelId, userId)
}

func (r *members) selectOne(ctx context.Context, query string, args ...any) (dto.Member, error) {
	row, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return dto.Member{}, err
	}

	member, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[dto.Member])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Member{}, ErrMemberNotFound
		}
		return dto.Member{}, err
	}
	return member, nil
}
