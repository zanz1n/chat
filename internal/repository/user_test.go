package repository_test

import (
	"fmt"
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
	"izanr.com/chat/internal/dto"
	"izanr.com/chat/internal/repository"
	"izanr.com/chat/internal/testutils"
	"izanr.com/chat/internal/utils"
)

var randstr = testutils.RandomString

func randUserData() dto.UserCreateData {
	return dto.UserCreateData{
		Username:    randstr(12),
		DisplayName: mo.None[string](),
		Email:       mo.None[string](),
		Pronouns:    mo.Some(fmt.Sprintf("%s/%s", randstr(3), randstr(4))),
		Bio:         mo.None[string](),
		Password:    randstr(32),
	}
}

func userRepo(t *testing.T) repository.UserStorer {
	db, cleanup := testutils.CreateDatabase(t)
	t.Cleanup(cleanup)

	return repository.NewPgUsers(db)
}

func TestInsertGet(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	assert := require.New(t)

	users := userRepo(t)

	_, err := users.GetById(ctx, utils.NewUUID())
	assert.ErrorIs(err, repository.ErrUserNotFound)

	user, err := users.Insert(ctx, randUserData())
	assert.NoError(err)

	user1, err := users.GetById(ctx, user.ID)
	assert.NoError(err)

	assert.Equal(user, user1)
}

func TestUpdatePassword(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	assert := require.New(t)

	users := userRepo(t)

	_, err := users.UpdatePassword(ctx,
		utils.NewUUID(),
		dto.UserUpdatePasswordData{},
	)
	assert.ErrorIs(err, repository.ErrUserNotFound)

	data := randUserData()
	user, err := users.Insert(ctx, data)
	assert.NoError(err)

	udata := dto.UserUpdatePasswordData{
		OldPassword: data.Password,
		NewPassword: testutils.RandomString(24),
	}

	user2, err := users.UpdatePassword(ctx, user.ID, udata)
	assert.NoError(err)

	assert.True(user.PasswordMatches(data.Password))
	assert.True(user2.PasswordMatches(udata.NewPassword))

	user3, err := users.GetById(ctx, user.ID)
	assert.NoError(err)

	assert.Equal(user3.UpdatedAt, user2.UpdatedAt)
	assert.NotEqual(user2.UpdatedAt, user.UpdatedAt)

	assert.True(user3.PasswordMatches(udata.NewPassword))
}
