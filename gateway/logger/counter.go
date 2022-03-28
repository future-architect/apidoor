package logger

import (
	"context"
	"fmt"
	"github.com/Songmu/flextime"
	"github.com/future-architect/apidoor/gateway/model"
	"sync"
	"time"
)

var (
	APICounter = APICallCounter{}

	DefaultCountValidSpan = 30 * time.Second
	DefaultCountSpanDays  = 30
)

type APICallCounter struct {
	sync.Map
}

func (ac *APICallCounter) GetCount(ctx context.Context, contractID int, path model.URITemplate) (int, error) {
	key := counterKey{
		contractID: contractID,
		path:       path.JoinPath(),
	}
	return ac.getCount(ctx, key)
}

func (ac *APICallCounter) getCount(ctx context.Context, key counterKey) (int, error) {
	stat, ok := ac.Load(key)
	if ok {
		stat := stat.(counterStat)
		if flextime.Now().Before(stat.expireAt) {
			return stat.calls, nil
		}
	}
	return ac.updateCount(ctx, key)
}

func (ac *APICallCounter) updateCount(ctx context.Context, key counterKey) (int, error) {
	count64, err := db.countBillingAccessLogDB(ctx, key.contractID, key.path, ac.countStartAt())
	if err != nil {
		return 0, fmt.Errorf("count api call db error: %w", err)
	}
	count := int(count64)

	ac.Store(key, counterStat{
		calls:    count,
		expireAt: ac.countExpireAt(),
	})
	return count, nil
}

func (ac *APICallCounter) countStartAt() time.Time {
	return flextime.Now().AddDate(0, 0, -DefaultCountSpanDays)
}

func (ac *APICallCounter) countExpireAt() time.Time {
	return flextime.Now().Add(DefaultCountValidSpan)
}

type counterKey struct {
	contractID int
	path       string
}

type counterStat struct {
	calls    int
	expireAt time.Time
}
