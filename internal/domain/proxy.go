package domain

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Proxy struct {
	Id             int64     `json:"proxy_id"           db:"proxy_id"        extensions:"x-order=1"`
	Protocol       string    `json:"protocol"           db:"protocol"        extensions:"x-order=2"`
	Username       string    `json:"username"           db:"username"        extensions:"x-order=3"`
	Password       string    `json:"password"           db:"password"        extensions:"x-order=4"`
	Host           string    `json:"host"               db:"host"            extensions:"x-order=5"`
	Port           int64     `json:"port"               db:"port"            extensions:"x-order=6"`
	OccupiesCount  int64     `json:"occupies_count"     db:"occupies_count"  extensions:"x-order=7"`
	ExpirationDate time.Time `json:"expiration_date"    db:"expiration_date" extensions:"x-order=7"`
	Enabled        bool      `json:"enabled"            db:"enabled"         extensions:"x-order=8"`
}

type ProxyList struct {
	Proxies []Proxy `json:"proxies" extensions:"x-order=1"`
	Offset  int64   `json:"offset"  extensions:"x-order=2"`
	Total   int64   `json:"total"   extensions:"x-order=3"`
}

type ProxyOccupy struct {
	Proxy Proxy  `json:"proxy" extensions:"x-order=1"`
	Key   string `json:"key"   extensions:"x-order=2"`
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
