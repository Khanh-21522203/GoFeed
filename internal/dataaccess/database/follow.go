package database

import (
	"GoFeed/internal/utils"
	"context"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	CreateFollow(ctx context.Context, follow Follow) error
	GetFollowerCountOfAccount(ctx context.Context, account_id uint64) (int, error)
	GetFollowersOfAccount(ctx context.Context, account_id uint64) ([]uint64, error)
	GetFollowingCountOfAccount(ctx context.Context, account_id uint64) (int, error)
	GetFollowingsOfAccount(ctx context.Context, account_id uint64) ([]uint64, error)
	DeleteFollow(ctx context.Context, follow Follow) error
	WithDatabase(database Database) FollowDataAccessor
}

type followDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewFollowDataAccessor(database *goqu.Database, logger *zap.Logger) FollowDataAccessor {
	return &followDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (f followDataAccessor) CreateFollow(ctx context.Context, follow Follow) error {
	logger := utils.LoggerWithContext(ctx, f.logger)

	_, err := f.database.
		Insert(TabNameFollows).
		Rows(goqu.Record{
			ColNameFollowsAccountID:   follow.AccountID,
			ColNameFollowsFollowingID: follow.FollowingID,
		}).
		Executor().Exec()

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create follow")
		return status.Error(codes.Internal, "failed to create follow")
	}
	return nil
}

func (f followDataAccessor) GetFollowerCountOfAccount(ctx context.Context, account_id uint64) (int, error) {
	logger := utils.LoggerWithContext(ctx, f.logger)

	var followers []uint64
	err := f.database.
		Select(ColNameFollowsAccountID).
		From(TabNameFollows).
		Where(goqu.C(ColNameFollowsFollowingID).Eq(account_id)).
		ScanValsContext(ctx, &followers)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get followers of account")
		return 0, status.Error(codes.Internal, "failed to get followers of account")
	}
	return len(followers), nil
}

func (f followDataAccessor) GetFollowersOfAccount(ctx context.Context, account_id uint64) ([]uint64, error) {
	logger := utils.LoggerWithContext(ctx, f.logger)

	var followers []uint64
	err := f.database.
		Select(ColNameFollowsAccountID).
		From(TabNameFollows).
		Where(goqu.C(ColNameFollowsFollowingID).Eq(account_id)).
		ScanValsContext(ctx, &followers)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get followers of account")
		return nil, status.Error(codes.Internal, "failed to get followers of account")
	}
	return followers, nil
}

func (f followDataAccessor) GetFollowingCountOfAccount(ctx context.Context, account_id uint64) (int, error) {
	logger := utils.LoggerWithContext(ctx, f.logger)

	var followings []uint64
	err := f.database.
		Select(ColNameFollowsFollowingID).
		From(TabNameFollows).
		Where(goqu.C(ColNameFollowsAccountID).Eq(account_id)).
		ScanValsContext(ctx, &followings)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get followers of account")
		return 0, status.Error(codes.Internal, "failed to get followers of account")
	}
	return len(followings), nil
}

func (f followDataAccessor) GetFollowingsOfAccount(ctx context.Context, account_id uint64) ([]uint64, error) {
	logger := utils.LoggerWithContext(ctx, f.logger)

	var followings []uint64
	err := f.database.
		Select(ColNameFollowsFollowingID).
		From(TabNameFollows).
		Where(goqu.C(ColNameFollowsAccountID).Eq(account_id)).
		ScanValsContext(ctx, &followings)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get followers of account")
		return nil, status.Error(codes.Internal, "failed to get followers of account")
	}
	return followings, nil
}

func (f followDataAccessor) DeleteFollow(ctx context.Context, follow Follow) error {
	logger := utils.LoggerWithContext(ctx, f.logger)

	_, err := f.database.
		Delete(TabNameFollows).
		Where(
			goqu.C(ColNameFollowsAccountID).Eq(follow.AccountID),
			goqu.C(ColNameFollowsFollowingID).Eq(follow.FollowingID),
		).Executor().Exec()

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to delete follow")
		return status.Error(codes.Internal, "failed to delete follow")
	}
	return nil
}

func (f followDataAccessor) WithDatabase(database Database) FollowDataAccessor {
	return &followDataAccessor{
		database: database,
	}
}
