package usecase

import (
	"context"
	"proxy_manager/internal/domain"
)

type UseCase struct {
	proxyRepo domain.ProxyRepository
}

func New(proxyRepo domain.ProxyRepository) UseCase {
	return UseCase{proxyRepo: proxyRepo}
}

func (u *UseCase) CreateProxy(ctx context.Context, proxy domain.Proxy) (domain.Proxy, error) {
	return u.proxyRepo.CreateProxy(ctx, proxy)
}
