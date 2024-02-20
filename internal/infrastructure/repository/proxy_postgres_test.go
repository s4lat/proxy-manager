package repository_test

import (
	"context"
	"proxy_manager/internal/infrastructure/repository"
	"proxy_manager/pkg/logger"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const testPostgresURL = "postgres://taskManager:taskManager@localhost:5433/taskManager"

func TestPostgresProxyRepository_OccupyMostAvailableProxy(t *testing.T) {
	ctx := context.Background()

	pgxPool, err := pgxpool.New(ctx, testPostgresURL+"?pool_max_conns=1")
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewPostgresProxyRepository(context.Background(), pgxPool, time.Minute*3, logger.NewTestLogger(t))
	proxyOccupy, err := repo.OccupyMostAvailableProxy(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proxyOccupy)
}
