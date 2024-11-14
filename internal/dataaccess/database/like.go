package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
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
	CreateLike(ctx context.Context, account_id uint64, post_id uint64) error
	GetLikeCountOfPost(ctx context.Context, post_id uint64) (int, error)
	GetLikeAccountsOfPost(ctx context.Context, post_id uint64) ([]uint64, error)
	DeleteLike(ctx context.Context, account_id uint64, post_id uint64) error
	WithDatabase(database Database) LikeDataAccessor
}

type likeDataAccessor struct {
	database Database
}

func NewLikeDataAccessor(database *goqu.Database) LikeDataAccessor {
	return &likeDataAccessor{
		database: database,
	}
}

func (l *likeDataAccessor) CreateLike(ctx context.Context, account_id uint64, post_id uint64) error {
	return nil
}

func (l *likeDataAccessor) GetLikeCountOfPost(ctx context.Context, post_id uint64) (int, error) {
	return 0, nil
}

func (l *likeDataAccessor) GetLikeAccountsOfPost(ctx context.Context, post_id uint64) ([]uint64, error) {
	return nil, nil
}

func (l *likeDataAccessor) DeleteLike(ctx context.Context, account_id uint64, post_id uint64) error {
	return nil
}

func (l *likeDataAccessor) WithDatabase(database Database) LikeDataAccessor {
	l.database = database
	return l
}
