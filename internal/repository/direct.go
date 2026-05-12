package repository

import (
	"bytes"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"izanr.com/chat/internal/dto"
	"izanr.com/chat/internal/utils"
)

type DirectStorer interface {
	Insert(context.Context, uuid.UUID, uuid.UUID) (dto.DirectChannel, error)
	GetById(context.Context, uuid.UUID, uuid.UUID) (dto.DirectChannel, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) (dto.DirectChannel, error)
}

type direct_channels struct {
	q utils.Querier
}

func NewPgDirects(q utils.Querier) DirectStorer {
	return &direct_channels{q: q}
}

func (r *direct_channels) Insert(
	ctx context.Context,
	user1 uuid.UUID,
	user2 uuid.UUID,
) (dto.DirectChannel, error) {
	// TODO: check for conflict
	return r.selectOne(ctx, query_directs_insert, user1, user2)
}

func (r *direct_channels) GetById(
	ctx context.Context,
	user1 uuid.UUID,
	user2 uuid.UUID,
) (dto.DirectChannel, error) {
	return r.selectOne(ctx, query_directs_get_by_id, user1, user2)
}

func (r *direct_channels) Delete(
	ctx context.Context,
	user1 uuid.UUID,
	user2 uuid.UUID,
) (dto.DirectChannel, error) {
	return r.selectOne(ctx, query_directs_delete, user1, user2)
}

func (r *direct_channels) selectOne(
	ctx context.Context,
	query string,
	minor uuid.UUID,
	major uuid.UUID,
) (dto.DirectChannel, error) {
	if bytes.Compare(minor[:], major[:]) == 1 {
		tmp := minor
		minor = major
		major = tmp
	}

	row, err := r.q.Query(ctx, query, minor, major)
	if err != nil {
		return dto.DirectChannel{}, err
	}

	channel, err := pgx.CollectExactlyOneRow(
		row,
		pgx.RowToStructByName[dto.DirectChannel],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.DirectChannel{}, ErrChannelNotFound
		}
		return dto.DirectChannel{}, err
	}
	return channel, nil
}
