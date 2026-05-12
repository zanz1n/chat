package testutils

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	sql1 "izanr.com/chat/sql"
)

var (
	PostgresImage   = "docker.io/library/postgres:18-alpine"
	PostgresTimeout = time.Minute
	EnableTestLogs  = false
)

type pgcontaienr struct {
	*postgres.PostgresContainer
	user     string
	password string
}

func (c *pgcontaienr) ExecSQL(ctx context.Context, sql string) (int, error) {
	cmd := []string{
		"sh", "-c",
		fmt.Sprintf("psql -U %s -c \"%s\"", c.user, sql),
	}
	ecode, rd, err := c.Container.Exec(ctx, cmd)
	if err != nil {
		return ecode, err
	}

	out := io.Discard
	if EnableTestLogs {
		out = os.Stdout
	}

	_, _ = io.Copy(out, rd)
	out.Write(append([]byte(strings.Join(cmd, " ")), '\n'))

	return ecode, nil
}

var onceDB = sync.OnceValues(func() (*pgcontaienr, error) {
	start := time.Now()

	username := "docker"
	password := RandomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), PostgresTimeout)
	defer cancel()

	opts := []testcontainers.ContainerCustomizer{
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		postgres.WithDatabase(username),
		postgres.BasicWaitStrategies(),
	}

	if !EnableTestLogs {
		opts = append(opts,
			testcontainers.WithLogger(log.New(io.Discard, "", log.Flags())),
		)
	}

	pg, err := postgres.Run(ctx, PostgresImage, opts...)
	if err != nil {
		return nil, err
	}

	// necessary for postgres to spin up correctly
	// time.Sleep(3 * time.Second)

	res := new(pgcontaienr{
		PostgresContainer: pg,
		user:              username,
		password:          password,
	})

	log.Printf(
		"DBUtils: Postgres container running after[%v]",
		time.Since(start).Round(time.Millisecond),
	)

	return res, nil
})

func CreateDatabase(t *testing.T) (*pgxpool.Pool, func()) {
	assert := require.New(t)

	container, err := onceDB()
	assert.NoError(err, "Run postgres testcontainer failed")

	start := time.Now()

	endpoint, err := container.PortEndpoint(t.Context(), "5432/tcp", "")
	assert.NoError(err, "Run postgres testcontainer failed")

	dbName := "db_" + RandomStringLower(8)
	ecode, err := container.ExecSQL(t.Context(),
		fmt.Sprintf("CREATE DATABASE %s;", dbName),
	)
	assert.NoErrorf(err, "Create postgres db %s failed", dbName)
	assert.Equal(0, ecode)

	connstr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		container.user,
		container.password,
		endpoint,
		dbName,
	)

	pool, err := pgxpool.New(t.Context(), connstr)
	assert.NoError(err, "Connect to postgres failed")

	err = sql1.MigrateSkipKV(t.Context(), pool)
	assert.NoError(err, "Migrate postgres failed")

	log.Printf("DBUtils: Created test db[%s] took[%v]",
		dbName,
		time.Since(start).Round(time.Millisecond),
	)

	return pool, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		defer pool.Close()

		_ = ctx

		_, err := container.ExecSQL(ctx,
			fmt.Sprintf("DROP DATABASE %s WITH (FORCE);", dbName),
		)
		if err != nil {
			log.Printf("DBUtils: Faield to drop test db[%s]: %s", dbName, err)
		} else {
			log.Printf("DBUtils: Dropped test db[%s] runtime[%v]",
				dbName,
				time.Since(start).Round(10*time.Millisecond),
			)
		}
	}
}
