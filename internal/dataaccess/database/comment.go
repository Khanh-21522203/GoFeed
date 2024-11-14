package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameComments = goqu.T("comments")
)

const (
	ColNameCommentsID        = "id"
	ColNameCommentsAccountID = "account_id"
	ColNameCommentsPostID    = "post_id"
	ColNameCommentsContent   = "content"
	ColNameCommentsCreatedAt = "created_at"
)

type Comment struct {
	ID        uint64
	AccountID uint64
	PostID    string
	Content   string
	CreatedAt time.Time
}

type CommentDataAccessor interface {
	CreateComment(ctx context.Context, account_id uint64, post_id uint64, content string, created_at time.Time) (uint64, error)
	GetCommentCountOfPost(ctx context.Context, post_id uint64) (int, error)
	GetCommentsOfPost(ctx context.Context, post_id uint64) ([]uint64, error)
	UpdateComment(ctx context.Context, comment Comment) error
	DeleteComment(ctx context.Context, id uint64) error
	WithDatabase(database Database) CommentDataAccessor
}

type commentDataAccessor struct {
	database Database
}

func NewCommentDataAccessor(database *goqu.Database) CommentDataAccessor {
	return &commentDataAccessor{
		database: database,
	}
}

func (c *commentDataAccessor) CreateComment(ctx context.Context, account_id uint64, post_id uint64, content string, created_at time.Time) (uint64, error) {
	return 0, nil
}

func (c *commentDataAccessor) GetCommentCountOfPost(ctx context.Context, post_id uint64) (int, error) {
	return 0, nil
}

func (c *commentDataAccessor) GetCommentsOfPost(ctx context.Context, post_id uint64) ([]uint64, error) {
	return nil, nil
}

func (c *commentDataAccessor) UpdateComment(ctx context.Context, comment Comment) error {
	return nil
}

func (c *commentDataAccessor) DeleteComment(ctx context.Context, id uint64) error {
	return nil
}

func (c *commentDataAccessor) WithDatabase(database Database) CommentDataAccessor {
	c.database = database
	return c
}
