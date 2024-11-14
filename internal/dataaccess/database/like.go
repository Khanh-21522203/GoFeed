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
	TabNameLikes = goqu.T("likes")
)

const (
	ColNameLikesAccountID = "account_id"
	ColNameLikesPostID    = "post_id"
)

type Like struct {
	AccountID uint64
	PostID    string
}

type LikeDataAccessor interface {
	CreateLike(ctx context.Context, like Like) error
	GetLikeCountOfPost(ctx context.Context, post_id uint64) (int, error)
	GetLikeAccountsOfPost(ctx context.Context, post_id uint64) ([]uint64, error)
	DeleteLike(ctx context.Context, account_id uint64, post_id uint64) error
	WithDatabase(database Database) LikeDataAccessor
}

type likeDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewLikeDataAccessor(database *goqu.Database, logger *zap.Logger) LikeDataAccessor {
	return &likeDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (l likeDataAccessor) CreateLike(ctx context.Context, like Like) error {
	logger := utils.LoggerWithContext(ctx, l.logger)

	_, err := l.database.
		Insert(TabNameLikes).
		Rows(goqu.Record{
			ColNameLikesAccountID: like.AccountID,
			ColNameLikesPostID:    like.PostID,
		}).
		Executor().Exec()

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create like")
		return status.Error(codes.Internal, "failed to create like")
	}
	return nil
}

func (l likeDataAccessor) GetLikeCountOfPost(ctx context.Context, post_id uint64) (int, error) {
	logger := utils.LoggerWithContext(ctx, l.logger)

	var accounts []uint64
	err := l.database.
		Select(ColNameLikesAccountID).
		From(TabNameLikes).
		Where(goqu.C(ColNameLikesPostID).Eq(post_id)).
		ScanValsContext(ctx, &accounts)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get like count of post")
		return 0, status.Error(codes.Internal, "failed to get like count of post")
	}
	return len(accounts), nil
}

func (l likeDataAccessor) GetLikeAccountsOfPost(ctx context.Context, post_id uint64) ([]uint64, error) {
	logger := utils.LoggerWithContext(ctx, l.logger)

	var accounts []uint64
	err := l.database.
		From(TabNameLikes).
		Select(ColNameLikesAccountID).
		Where(goqu.C(ColNameLikesPostID).Eq(post_id)).
		ScanValsContext(ctx, &accounts)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get accounts who like post")
		return nil, status.Error(codes.Internal, "failed to get accounts who like post")
	}
	return accounts, nil
}

func (l likeDataAccessor) DeleteLike(ctx context.Context, account_id uint64, post_id uint64) error {
	logger := utils.LoggerWithContext(ctx, l.logger)

	_, err := l.database.
		Delete(TabNameLikes).
		Where(goqu.C(ColNameLikesAccountID).Eq(account_id), goqu.C(ColNameLikesPostID).Eq(post_id)).
		Executor().Exec()

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to unlike post")
		return status.Error(codes.Internal, "failed to unlike post")
	}
	return nil
}

func (l likeDataAccessor) WithDatabase(database Database) LikeDataAccessor {
	return &likeDataAccessor{
		database: database,
	}
}
