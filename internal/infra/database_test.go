package infra_test

import (
	"context"
	"testing"
	"time"

	"github.com/duddy57/cotrip-api/internal/infra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNewDB_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = pgContainer.Terminate(ctx) })

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	pool, err := infra.NewDB(ctx, "testuser", "testpass", host, mappedPort.Port(), "testdb")
	require.NoError(t, err)
	require.NotNil(t, pool)
	defer pool.Close()

	err = pool.Ping(ctx)
	assert.NoError(t, err)
}

func TestNewDB_InvalidDSN(t *testing.T) {
	ctx := context.Background()

	pool, err := infra.NewDB(ctx, "user", "pass", "invalid-host-that-does-not-exist", "5432", "db")

	require.NoError(t, err)
	require.NotNil(t, pool)
	defer pool.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err = pool.Ping(pingCtx)
	assert.Error(t, err, "esperado erro de conexão para host inválido")
}

func TestNewDB_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	pool, err := infra.NewDB(ctx, "user", "pass", "localhost", "5432", "db")
	if err != nil {
		assert.Nil(t, pool)
		return
	}

	defer pool.Close()
	err = pool.Ping(ctx)
	assert.Error(t, err)
}
