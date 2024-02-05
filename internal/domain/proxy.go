package domain

import (
	"context"
	"time"
)

type Proxy struct {
	Id           int    `json:"id" db:"id"`
	Protocol     string `json:"protocol" db:"protocol"`
	Username     string `json:"username" db:"username"`
	Password     string `json:"password" db:"password"`
	Host         string `json:"host" db:"host"`
	Port         int    `json:"port" db:"port"`
	ClientsCount int    `json:"clients_count" db:"clients_count"`
}

type ProxyRepository interface {
	CreateProxy(ctx context.Context, proxy Proxy) (Proxy, error)

	GetById(ctx context.Context, proxyId int) (Proxy, error)
	LockAnyProxy(ctx context.Context, maxClientsCounter int) (Proxy, error)

	Update(context.Context, Proxy) (Proxy, error)
	Delete(context.Context, Proxy) error

	StartClientCountAutoReducer(ctx context.Context, period time.Duration)
}
