package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"proxy_manager/internal/domain"
	"time"
)

type PostgresProxyRepository struct {
	connPool *pgxpool.Pool
}

func NewPostgresProxyRepository(ctx context.Context, connPool *pgxpool.Pool) PostgresProxyRepository {
	return PostgresProxyRepository{
		connPool: connPool,
	}
}

func (p PostgresProxyRepository) CreateProxy(ctx context.Context, proxy domain.Proxy) (domain.Proxy, error) {
	q := "INSERT INTO proxies(protocol, username, password, host, port) VALUES ($1, $2, $3, $4, $5) RETURNING *"
	rows, _ := p.connPool.Query(ctx, q, proxy.Protocol, proxy.Username, proxy.Password, proxy.Host, proxy.Port)

	createdProxy, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.Proxy])
	if err != nil {
		return domain.Proxy{}, errors.Join(domain.ErrOnCreate, err)
	}
	return *createdProxy, nil
}

func (p PostgresProxyRepository) GetById(ctx context.Context, proxyId int) (domain.Proxy, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresProxyRepository) LockAnyProxy(ctx context.Context, maxClientsCounter int) (domain.Proxy, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresProxyRepository) Update(ctx context.Context, proxy domain.Proxy) (domain.Proxy, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresProxyRepository) Delete(ctx context.Context, proxy domain.Proxy) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresProxyRepository) StartClientCountAutoReducer(ctx context.Context, period time.Duration) {
	//TODO implement me
	panic("implement me")
}
