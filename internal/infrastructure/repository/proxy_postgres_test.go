package repository_test

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"proxy_manager/internal/infrastructure/repository"
	"testing"
)

const testPostgresUrl = "postgres://taskManager:taskManager@localhost:5433/taskManager"

func TestPostgresProxyRepository_OccupyMostAvailableProxy(t *testing.T) {
	ctx := context.Background()

	pgxPool, err := pgxpool.New(ctx, testPostgresUrl+"?pool_max_conns=1")
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewPostgresProxyRepository(context.Background(), pgxPool)
	proxyOccupy, err := repo.OccupyMostAvailableProxy(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proxyOccupy)
}
