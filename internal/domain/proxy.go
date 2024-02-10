package domain

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type Proxy struct {
	Id            int64  `json:"proxy_id" db:"proxy_id"`
	Protocol      string `json:"protocol" db:"protocol"`
	Username      string `json:"username" db:"username"`
	Password      string `json:"password" db:"password"`
	Host          string `json:"host" db:"host"`
	Port          int64  `json:"port" db:"port"`
	OccupiesCount int64  `json:"occupies_count" db:"occupies_count"`
}

type ProxyList struct {
	Proxies []Proxy `json:"proxies"`
	Offset  int64   `json:"offset"`
	Total   int64   `json:"total"`
}

type ProxyOccupy struct {
	Proxy Proxy  `json:"proxy"`
	Key   string `json:"key"`
}

type ProxyRepository interface {
	CreateProxy(ctx context.Context, proxy Proxy) (Proxy, error)
	GetProxy(ctx context.Context, proxyId int64) (Proxy, error)
	UpdateProxy(ctx context.Context, updatedProxy Proxy) (Proxy, error)
	DeleteProxy(ctx context.Context, proxyId int64) error

	GetProxyList(ctx context.Context, offset int64, limit int64) (ProxyList, error)

	OccupyMostAvailableProxy(ctx context.Context) (ProxyOccupy, error)
	ReleaseProxy(ctx context.Context, key string) error
}

func (p *Proxy) Validate() error {
	if !isValidProtocol(p.Protocol) {
		return fmt.Errorf("invalid protocol, allowed protocols: (%s)", strings.Join(allowedProtocols, ", "))
	}

	if p.Host == "" {
		return errors.New("host can't be empty string")
	}

	if !(p.Port > 0) {
		return errors.New("port must be > 0")
	}
	return nil
}

var allowedProtocols = []string{"http", "https", "socks5"}

func isValidProtocol(protocol string) bool {
	for _, allowedProtocol := range allowedProtocols {
		if protocol == allowedProtocol {
			return true
		}
	}
	return false
}
