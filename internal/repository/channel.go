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
	ErrChannelNotFound      = fmt.Errorf("channel not found")
	ErrChannelIsDirect      = fmt.Errorf("channel is direct: mutations not allowed")
	ErrUserAlreadyInChannel = fmt.Errorf("user is already in the channel")
)

type ChannelStorer interface {
	Insert(context.Context, uuid.UUID, dto.ChannelCreateData) (dto.Channel, error)

	GetById(context.Context, uuid.UUID) (dto.Channel, error)
	GetByUser(context.Context, uuid.UUID, dto.Pagination) ([]dto.Channel, error)

	// TODO: implement full text search
	SearchByUser(
		ctx context.Context,
		id uuid.UUID,
		text string,
		limit int,
	) ([]dto.Channel, error)

	Update(context.Context, uuid.UUID, dto.ChannelUpdateData) (dto.Channel, error)
	UpdatePicture(
		ctx context.Context,
		channelId uuid.UUID,
		pictureId uuid.UUID,
	) error

	Delete(context.Context, uuid.UUID) (dto.Channel, error)
}

type channels struct {
	q utils.Querier
}

func NewPgChannels(q utils.Querier) ChannelStorer {
	return &channels{q: q}
}

func (r *channels) Insert(
	ctx context.Context,
	userId uuid.UUID,
	data dto.ChannelCreateData,
) (dto.Channel, error) {
	return r.selectOne(ctx, query_channels_update,
		userId,
		data.Name,
		data.Description,
	)
}

func (r *channels) GetById(ctx context.Context, id uuid.UUID) (dto.Channel, error) {
	return r.selectOne(ctx, query_channels_get_by_id, id)
}

func (r *channels) GetByUser(
	ctx context.Context,
	userId uuid.UUID,
	pag dto.Pagination,
) ([]dto.Channel, error) {
	rows, err := r.q.Query(ctx, query_channels_get_by_user,
		userId,
		pag.LastSeen,
		pag.Limit,
	)
	if err != nil {
		return nil, err
	}

	channels, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.Channel])
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (r *channels) SearchByUser(
	ctx context.Context,
	id uuid.UUID,
	text string,
	limit int,
) ([]dto.Channel, error) {
	_ = query_channels_search_by_user
	panic("unimplemented")
}

func (r *channels) Update(
	ctx context.Context,
	id uuid.UUID,
	data dto.ChannelUpdateData,
) (dto.Channel, error) {
	return r.selectOne(ctx, query_channels_update,
		id,
		data.Name,
		data.Description,
	)
}

func (r *channels) UpdatePicture(
	ctx context.Context,
	channelId uuid.UUID,
	pictureId uuid.UUID,
) error {
	cmd, err := r.q.Exec(ctx, query_channels_update_picture,
		channelId,
		pictureId,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() <= 0 {
		return ErrChannelNotFound
	}
	return nil
}

func (r *channels) Delete(ctx context.Context, id uuid.UUID) (dto.Channel, error) {
	return r.selectOne(ctx, query_channels_delete, id)
}

func (r *channels) selectOne(ctx context.Context, query string, args ...any) (dto.Channel, error) {
	row, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return dto.Channel{}, err
	}

	channel, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[dto.Channel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Channel{}, ErrChannelNotFound
		}
		return dto.Channel{}, err
	}
	return channel, nil
}
