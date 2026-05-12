package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"izanr.com/chat/internal/testutils"
	"izanr.com/chat/internal/utils"
)

func TestHash(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	passwd := testutils.RandomString(32)
	hash := utils.HashPassword(passwd)

	assert.True(utils.CheckPasswordHash(hash, passwd))
}
