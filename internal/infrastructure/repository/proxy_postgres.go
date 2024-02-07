package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"proxy_manager/internal/domain"
	"proxy_manager/internal/usecase"
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
	q := "INSERT INTO proxy(protocol, username, password, host, port) VALUES ($1, $2, $3, $4, $5) RETURNING *;"
	rows, _ := p.connPool.Query(ctx, q, proxy.Protocol, proxy.Username, proxy.Password, proxy.Host, proxy.Port)

	createdProxy, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.Proxy])
	if err != nil {
		return domain.Proxy{}, err
	}
	return *createdProxy, nil
}

func (p PostgresProxyRepository) GetProxyList(ctx context.Context, offset int64, limit int64) (domain.ProxyList, error) {
	q := "WITH t AS (SELECT * FROM proxy) SELECT * FROM  (TABLE  t OFFSET $1 LIMIT  $2) sub RIGHT JOIN (SELECT count(*) FROM t) AS c(total) ON TRUE;"
	rows, _ := p.connPool.Query(ctx, q, offset, limit)
	rowsAsMap, err := pgx.CollectRows(rows, pgx.RowToMap)
	if err != nil {
		return domain.ProxyList{}, err
	}

	proxyList := domain.ProxyList{
		Total:  rowsAsMap[0]["total"].(int64),
		Offset: offset,
	}

	if rowsAsMap[0]["proxy_id"] == nil {
		return proxyList, nil
	}

	for _, row := range rowsAsMap {
		proxyList.Proxies = append(proxyList.Proxies, domain.Proxy{
			Id:       row["proxy_id"].(int64),
			Protocol: row["protocol"].(string),
			Host:     row["host"].(string),
			Port:     row["port"].(int64),
			Username: row["username"].(string),
			Password: row["password"].(string),
		})
	}
	return proxyList, nil
}

func (p PostgresProxyRepository) GetProxyById(ctx context.Context, proxyId int64) (domain.Proxy, error) {
	q := "SELECT proxy.*, COUNT(proxy_occupy.proxy_id) AS occupies_count FROM proxy LEFT JOIN proxy_occupy ON proxy.proxy_id = proxy_occupy.proxy_id WHERE proxy.proxy_id = $1 GROUP BY proxy.proxy_id;"
	rows, _ := p.connPool.Query(ctx, q, proxyId)

	proxy, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.Proxy])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Proxy{}, usecase.ErrNotFound
		}
		return domain.Proxy{}, err
	}
	return *proxy, nil
}

func (p PostgresProxyRepository) UpdateProxy(ctx context.Context, proxy domain.Proxy) (domain.Proxy, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresProxyRepository) DeleteProxy(ctx context.Context, proxy domain.Proxy) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresProxyRepository) OccupyMostAvailableProxy(ctx context.Context) (domain.Proxy, error) {
	//TODO implement me
	// Отдает проксю с наименьшим кол-во occupies, создает с ней occupy;
	// Считаются только occupies с expire_timestamp, которое еще не expire.
	// EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) - для получение настоящего timestamp'.
	panic("implement me")
}

func (p PostgresProxyRepository) ReleaseProxy(ctx context.Context, key string) (domain.Proxy, error) {
	//TODO implement me
	// Удалает occupy с данным key;
	panic("implement me")
}

func (p PostgresProxyRepository) AutoCleanupOccupiedProxies(ctx context.Context, expireTime time.Duration) {
	//TODO implement me
	// EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) - 5 * 60 > create_timestamp - для сравнения;
	// чистить каждую минуту.
	panic("implement me")
}
