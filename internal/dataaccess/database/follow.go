package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameFollows = goqu.T("follows")
)

const (
	ColNameFollowsAccountID   = "account_id"
	ColNameFollowsFollowingID = "following_id"
)

type Follow struct {
	AccountID   uint64
	FollowingID string
}

type FollowDataAccessor interface {
	CreateFollow(ctx context.Context, account_id uint64, following_id uint64) error
	GetFollowerCountOfAccount(ctx context.Context, account_id uint64) (int, error)
	GetFollowersOfAccount(ctx context.Context, account_id uint64) ([]uint64, error)
	GetFollowingCountOfAccount(ctx context.Context, account_id uint64) (int, error)
	GetFollowingsOfAccount(ctx context.Context, account_id uint64) ([]uint64, error)
	DeleteFollow(ctx context.Context, account_id uint64, following_id uint64) error
	WithDatabase(database Database) FollowDataAccessor
}

type followDataAccessor struct {
	database Database
}

func NewFollowDataAccessor(database *goqu.Database) FollowDataAccessor {
	return &followDataAccessor{
		database: database,
	}
}

func (f *followDataAccessor) CreateFollow(ctx context.Context, account_id uint64, following_id uint64) error {
	return nil
}

func (f *followDataAccessor) GetFollowerCountOfAccount(ctx context.Context, account_id uint64) (int, error) {
	return 0, nil
}

func (f *followDataAccessor) GetFollowersOfAccount(ctx context.Context, account_id uint64) ([]uint64, error) {
	return nil, nil
}

func (f *followDataAccessor) GetFollowingCountOfAccount(ctx context.Context, account_id uint64) (int, error) {
	return 0, nil
}

func (f *followDataAccessor) GetFollowingsOfAccount(ctx context.Context, account_id uint64) ([]uint64, error) {
	return nil, nil
}

func (f *followDataAccessor) DeleteFollow(ctx context.Context, account_id uint64, following_id uint64) error {
	return nil
}

func (f *followDataAccessor) WithDatabase(database Database) FollowDataAccessor {
	f.database = database
	return f
}
