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

	"github.com/jackc/pgx/v5"
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
	basedb   string
}

func (c *pgcontaienr) ConnStr(ctx context.Context, db string) (string, error) {
	endpoint, err := c.PortEndpoint(ctx, "5432/tcp", "")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.user,
		c.password,
		endpoint,
		db,
	), nil
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

	res := new(pgcontaienr{
		PostgresContainer: pg,
		user:              username,
		password:          password,
		basedb:            username,
	})

	log.Printf(
		"DBUtils: Postgres container running after[%v]",
		time.Since(start).Round(time.Millisecond),
	)
	start = time.Now()

	connstr, err := res.ConnStr(ctx, res.basedb)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, connstr)
	if err != nil {
		return nil, err
	}
	defer pool.Close()

	if err = sql1.MigrateSkipKV(ctx, pool); err != nil {
		return nil, err
	}

	log.Printf(
		"DBUtils: Migrated postgres db[%s] took[%v]",
		res.basedb,
		time.Since(start).Round(time.Millisecond),
	)

	return res, nil
})

func CreateDatabase(t *testing.T) (*pgx.Conn, func()) {
	assert := require.New(t)

	container, err := onceDB()
	assert.NoError(err, "Run postgres testcontainer failed")

	start := time.Now()

	dbName := "db_" + RandomStringLower(8)
	ecode, err := container.ExecSQL(t.Context(),
		fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s STRATEGY=FILE_COPY;",
			dbName,
			container.basedb,
		),
	)
	assert.NoErrorf(err, "Create postgres db %s failed", dbName)
	assert.Equal(0, ecode)

	connstr, err := container.ConnStr(t.Context(), dbName)
	assert.NoError(err)

	conn, err := pgx.Connect(t.Context(), connstr)

	log.Printf("DBUtils: Created test db[%s] took[%v]",
		dbName,
		time.Since(start).Round(time.Millisecond),
	)

	return conn, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		defer conn.Close(ctx)

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
