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
	proxy.ExpirationDate = proxy.ExpirationDate.UTC()

	if err := proxy.Validate(); err != nil {
		return domain.Proxy{}, errors.Join(ErrInvalidData, err)
	}

	proxy, err := u.proxyRepo.CreateProxy(ctx, proxy)
	if err != nil {
		return domain.Proxy{}, errors.Join(ErrInRepo, err)
	}
	return proxy, nil
}

func (u *UseCase) UpdateProxy(ctx context.Context, updatedProxy domain.Proxy) (domain.Proxy, error) {
	updatedProxy.ExpirationDate = updatedProxy.ExpirationDate.UTC()

	if updatedProxy.ID <= 0 {
		return domain.Proxy{}, errors.Join(ErrInvalidData, errors.New("ProxyID must be > 0"))
	}

	if err := updatedProxy.Validate(); err != nil {
		return domain.Proxy{}, errors.Join(ErrInvalidData, err)
	}

	proxy, err := u.proxyRepo.UpdateProxy(ctx, updatedProxy)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return domain.Proxy{}, err
		}
		return domain.Proxy{}, errors.Join(ErrInRepo, err)
	}
	return proxy, nil
}

func (u *UseCase) DeleteProxy(ctx context.Context, proxyID int64) error {
	if err := u.proxyRepo.DeleteProxy(ctx, proxyID); err != nil {
		return errors.Join(ErrInRepo, err)
	}
	return nil
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

func (u *UseCase) GetProxy(ctx context.Context, proxyID int64) (domain.Proxy, error) {
	proxy, err := u.proxyRepo.GetProxy(ctx, proxyID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return domain.Proxy{}, err
		}
		return domain.Proxy{}, errors.Join(ErrInRepo, err)
	}

	return proxy, nil
}

func (u *UseCase) OccupyMostAvailableProxy(ctx context.Context) (domain.ProxyOccupy, error) {
	proxyOccupy, err := u.proxyRepo.OccupyMostAvailableProxy(ctx)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return domain.ProxyOccupy{}, err
		}
		return domain.ProxyOccupy{}, errors.Join(ErrInRepo, err)
	}

	return proxyOccupy, nil
}

func (u *UseCase) ReleaseProxy(ctx context.Context, key string) error {
	if err := u.proxyRepo.ReleaseProxy(ctx, key); err != nil {
		return errors.Join(ErrInRepo, err)
	}
	return nil
}
