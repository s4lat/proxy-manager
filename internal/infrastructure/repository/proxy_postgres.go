package repository

import (
	"context"
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"proxy_manager/internal/domain"
	"proxy_manager/internal/usecase"
	"time"
)

type PostgresProxyRepository struct {
	connPool *pgxpool.Pool
}

func NewPostgresProxyRepository(connPool *pgxpool.Pool) PostgresProxyRepository {
	return PostgresProxyRepository{
		connPool: connPool,
	}
}

func (p PostgresProxyRepository) CreateProxy(ctx context.Context, proxy domain.Proxy) (domain.Proxy, error) {
	q := "INSERT INTO proxy(protocol, username, password, host, port) VALUES ($1, $2, $3, $4, $5) RETURNING *, 0 as occupies_count;"
	rows, _ := p.connPool.Query(ctx, q, proxy.Protocol, proxy.Username, proxy.Password, proxy.Host, proxy.Port)

	createdProxy, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.Proxy])
	if err != nil {
		return domain.Proxy{}, err
	}
	return *createdProxy, nil
}

func (p PostgresProxyRepository) GetProxyList(ctx context.Context, offset int64, limit int64) (domain.ProxyList, error) {
	q := "WITH t AS (SELECT proxy.*, COUNT(proxy_occupy.proxy_id) AS occupies_count FROM proxy LEFT JOIN proxy_occupy ON proxy.proxy_id = proxy_occupy.proxy_id GROUP BY proxy.proxy_id) SELECT * FROM  (TABLE  t OFFSET $1 LIMIT  $2) sub RIGHT JOIN (SELECT count(*) FROM t) AS c(total) ON TRUE;"
	rows, _ := p.connPool.Query(ctx, q, offset, limit)
	rowsAsMap, err := pgx.CollectRows(rows, pgx.RowToMap)
	if err != nil {
		return domain.ProxyList{}, err
	}

	proxyList := domain.ProxyList{
		Total:  rowsAsMap[0]["total"].(int64),
		Offset: offset,
	}

	// Если proxy_id в первой записи = null, значит из БД нам не вернулось ни одной прокси
	// и вернулась только одна строка, в которой все поля, кроме "total", равны null
	if rowsAsMap[0]["proxy_id"] == nil {
		return proxyList, nil
	}

	for _, row := range rowsAsMap {
		proxyList.Proxies = append(proxyList.Proxies, domain.Proxy{
			Id:            row["proxy_id"].(int64),
			Protocol:      row["protocol"].(string),
			Host:          row["host"].(string),
			Port:          row["port"].(int64),
			Username:      row["username"].(string),
			Password:      row["password"].(string),
			OccupiesCount: row["occupies_count"].(int64),
		})
	}
	return proxyList, nil
}

func (p PostgresProxyRepository) GetProxy(ctx context.Context, proxyId int64) (domain.Proxy, error) {
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

func (p PostgresProxyRepository) DeleteProxy(ctx context.Context, proxyId int64) error {
	q := "DELETE FROM proxy WHERE proxy_id=$1;"
	_, err := p.connPool.Query(ctx, q, proxyId)
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresProxyRepository) OccupyMostAvailableProxy(ctx context.Context) (domain.ProxyOccupy, error) {
	selectQuery := "SELECT proxy.*, COUNT(proxy_occupy.proxy_id) AS occupies_count FROM proxy LEFT JOIN proxy_occupy ON proxy.proxy_id = proxy_occupy.proxy_id GROUP BY proxy.proxy_id ORDER BY occupies_count ASC LIMIT 1;"
	occupyQuery := "INSERT INTO proxy_occupy(proxy_id, create_timestamp) VALUES($1, EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)) RETURNING *;"

	tx, err := p.connPool.Begin(ctx)
	if err != nil {
		return domain.ProxyOccupy{}, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "LOCK TABLE proxy_occupy IN EXCLUSIVE MODE;")
	if err != nil {
		return domain.ProxyOccupy{}, err
	}

	rows, _ := tx.Query(ctx, selectQuery)

	proxy, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.Proxy])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ProxyOccupy{}, usecase.ErrNotFound
		}
		return domain.ProxyOccupy{}, err
	}

	rows, _ = tx.Query(ctx, occupyQuery, proxy.Id)
	occupyRowMap, err := pgx.CollectOneRow(rows, pgx.RowToMap)
	if err != nil {
		return domain.ProxyOccupy{}, err
	}

	keyBytesArray, ok := occupyRowMap["key"].([16]byte)
	if !ok {
		return domain.ProxyOccupy{}, errors.New("can't convert occupy.key to [16]byte")
	}

	key, err := uuid.FromBytes(keyBytesArray[:])
	if err != nil {
		return domain.ProxyOccupy{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.ProxyOccupy{}, err
	}

	proxy.OccupiesCount += 1
	return domain.ProxyOccupy{
		Proxy: *proxy,
		Key:   key.String(),
	}, nil
}

func (p PostgresProxyRepository) ReleaseProxy(ctx context.Context, key string) error {
	q := "DELETE FROM proxy_occupy WHERE key=$1;"
	_, err := p.connPool.Query(ctx, q, key)
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresProxyRepository) AutoCleanupOccupiedProxies(ctx context.Context, expireTime time.Duration) {
	//TODO implement me
	// EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) - 5 * 60 > create_timestamp - для сравнения;
	// чистить каждую минуту.
	panic("implement me")
}
