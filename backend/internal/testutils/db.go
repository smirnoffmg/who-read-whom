package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/what-writers-like/backend/internal/infrastructure/config"
	"github.com/what-writers-like/backend/internal/infrastructure/database"
)

func SetupTestDB(t *testing.T) (*database.Database, func()) {
	ctx := context.Background()

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithStartupTimeout(30*time.Second),
		),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)
	require.NoError(t, err)

	// Give the database a moment to fully initialize
	time.Sleep(500 * time.Millisecond)

	connStr, err := postgresContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Ensure sslmode=disable is in the connection string
	if connStr != "" {
		if connStr[len(connStr)-1] != '?' {
			connStr += "?sslmode=disable"
		} else {
			connStr += "sslmode=disable"
		}
	}

	cfg := &config.Config{DatabaseDSN: connStr, ServerPort: "8080"}
	db, err := database.NewDatabase(cfg)
	require.NoError(t, err)

	cleanup := func() {
		_ = postgresContainer.Terminate(ctx)
	}

	return db, cleanup
}
