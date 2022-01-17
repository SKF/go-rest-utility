package auth

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const defaultGracePeriod = 10 * time.Minute

type CachedTokenProvider struct {
	TokenProvider
	gracePeriod time.Duration

	m           sync.RWMutex
	refreshOnce *sync.Once

	rawToken    RawToken
	ttl         time.Time
	supplyError error
}

func NewCachedTokenProvider(provider TokenProvider) *CachedTokenProvider {
	if cachedProvider, ok := provider.(*CachedTokenProvider); ok {
		return cachedProvider
	}

	return &CachedTokenProvider{
		TokenProvider: provider,
		gracePeriod:   defaultGracePeriod,
		refreshOnce:   new(sync.Once),
	}
}

func (p *CachedTokenProvider) WithGracePeriod(duration time.Duration) *CachedTokenProvider {
	p.gracePeriod = duration
	return p
}

func (p *CachedTokenProvider) GetRawToken(ctx context.Context) (RawToken, error) {
	p.m.RLock()

	// Is the cached token still alive?
	if time.Now().Before(p.ttl) {
		defer p.m.RUnlock()
		return p.rawToken, nil
	}

	// Copy the current refresher syncronizer. Needed to avoid races
	// once we unlock the read mutex and want to replace it.
	refreshOnce := p.refreshOnce

	p.m.RUnlock()

	refreshOnce.Do(func() {
		p.m.Lock()
		defer p.m.Unlock()

		// This check avoid some scenarios where a double refresh is possible.
		if !time.Now().Before(p.ttl) {
			p.rawToken, p.ttl, p.supplyError = p.refreshToken(ctx)
		}

		// Replace the refresh synchronizer, to allow someone else to attempt a refresh.
		p.refreshOnce = new(sync.Once)
	})

	p.m.RLock()
	defer p.m.RUnlock()

	return p.rawToken, p.supplyError
}

func (p *CachedTokenProvider) refreshToken(ctx context.Context) (RawToken, time.Time, error) {
	token, err := p.TokenProvider.GetRawToken(ctx)
	if err != nil {
		return "", time.Time{}, err
	}

	exp, err := token.ParseExpires()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("extracting exp claim: %w", err)
	}

	exp = exp.Add(-p.gracePeriod)

	return token, exp, nil
}
