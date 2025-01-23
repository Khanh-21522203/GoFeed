package cache

import "context"

type NewFeed interface {
	Get(ctx context.Context, account_id uint64) ([]uint64, error)
}
