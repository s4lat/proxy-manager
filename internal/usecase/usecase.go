package usecase

import (
	"context"
	"errors"
	"proxy_manager/internal/domain"
)

type UseCase struct {
	proxyRepo domain.ProxyRepository
}

func New(proxyRepo domain.ProxyRepository) UseCase {
	return UseCase{proxyRepo: proxyRepo}
}

func (u *UseCase) CreateProxy(ctx context.Context, proxy domain.Proxy) (domain.Proxy, error) {
	if err := proxy.Validate(); err != nil {
		return domain.Proxy{}, errors.Join(ErrInvalidData, err)
	}

	proxy, err := u.proxyRepo.CreateProxy(ctx, proxy)
	if err != nil {
		return domain.Proxy{}, errors.Join(ErrInRepo, err)
	}
	return proxy, nil
}

func (u *UseCase) GetProxyList(ctx context.Context, offset int64, limit int64) (domain.ProxyList, error) {
	if offset < 0 || limit < 0 {
		return domain.ProxyList{}, errors.Join(ErrInvalidData, errors.New("offset and limit must be non negative"))
	}

	proxyList, err := u.proxyRepo.GetProxyList(ctx, offset, limit)
	if err != nil {
		return domain.ProxyList{}, errors.Join(ErrInRepo, err)
	}
	return proxyList, nil
}

func (u *UseCase) GetProxyById(ctx context.Context, proxyId int64) (domain.Proxy, error) {
	proxy, err := u.proxyRepo.GetProxyById(ctx, proxyId)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return domain.Proxy{}, err
		}
		return domain.Proxy{}, errors.Join(ErrInRepo, err)
	}

	return proxy, nil
}
