package database

import (
	"GoFeed/internal/utils"
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
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
	PostID    uint64
	Content   string
	CreatedAt time.Time
}

type CommentDataAccessor interface {
	CreateComment(ctx context.Context, comment Comment) (uint64, error)
	GetCommentCountOfPost(ctx context.Context, post_id uint64) (int, error)
	GetCommentsOfPost(ctx context.Context, post_id uint64) ([]Comment, error)
	GetCommentByIdWithXLock(ctx context.Context, id uint64) (Comment, error)
	UpdateComment(ctx context.Context, comment Comment) error
	DeleteComment(ctx context.Context, id uint64) error
	DeleteCommentOfPost(ctx context.Context, post_id uint64) error
	WithDatabase(database Database) CommentDataAccessor
}

type commentDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewCommentDataAccessor(database *goqu.Database, logger *zap.Logger) CommentDataAccessor {
	return &commentDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (c commentDataAccessor) CreateComment(ctx context.Context, comment Comment) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, c.logger)

	_, err := c.database.
		Insert(TabNameComments).
		Rows(goqu.Record{
			ColNameCommentsID:        comment.ID,
			ColNameCommentsAccountID: comment.AccountID,
			ColNameCommentsPostID:    comment.PostID,
			ColNameCommentsContent:   comment.Content,
			// ColNameCommentsCreatedAt: comment.CreatedAt,
		}).Executor().Exec()

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create comment")
		return 0, err
	}
	return comment.ID, nil
}

func (c commentDataAccessor) GetCommentCountOfPost(ctx context.Context, post_id uint64) (int, error) {
	logger := utils.LoggerWithContext(ctx, c.logger)

	var comments []uint64
	err := c.database.
		Select(ColNameCommentsPostID).
		From(TabNameComments).
		Where(goqu.C(ColNameCommentsPostID).Eq(post_id)).
		ScanValsContext(ctx, &comments)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get comment count of post")
		return 0, err
	}
	return len(comments), nil
}

func (c commentDataAccessor) GetCommentsOfPost(ctx context.Context, post_id uint64) ([]Comment, error) {
	logger := utils.LoggerWithContext(ctx, c.logger)

	var comments []Comment
	err := c.database.
		// Select(ColNameCommentsPostID).
		From(TabNameComments).
		Where(goqu.C(ColNameCommentsPostID).Eq(post_id)).
		ScanValsContext(ctx, &comments)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get comment count of post")
		return nil, err
	}
	return comments, nil
}

func (c commentDataAccessor) GetCommentByIdWithXLock(ctx context.Context, id uint64) (Comment, error) {
	logger := utils.LoggerWithContext(ctx, c.logger)

	var comment Comment
	err := c.database.
		From(TabNameComments).
		Where(goqu.C(ColNameCommentsID).Eq(id)).
		ForUpdate(goqu.Wait).
		ScanStructsContext(ctx, &comment)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get comment by id with Xlock")
		return Comment{}, err
	}
	return comment, nil
}

func (c commentDataAccessor) UpdateComment(ctx context.Context, comment Comment) error {
	logger := utils.LoggerWithContext(ctx, c.logger)

	_, err := c.database.
		Update(TabNameComments).
		Set(comment).
		Where(goqu.C(ColNameCommentsID).Eq(comment.ID)).
		Executor().Exec()

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to update comment")
		return err
	}
	return nil
}

func (c commentDataAccessor) DeleteComment(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, c.logger)

	_, err := c.database.
		Delete(TabNameComments).
		Where(goqu.C(ColNameCommentsID).Eq(id)).
		Executor().Exec()

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to update comment")
		return err
	}
	return nil
}

func (c commentDataAccessor) DeleteCommentOfPost(ctx context.Context, post_id uint64) error {
	logger := utils.LoggerWithContext(ctx, c.logger)

	_, err := c.database.
		Delete(TabNameComments).
		Where(goqu.C(ColNameCommentsPostID).Eq(post_id)).
		Executor().Exec()

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to update comment")
		return err
	}
	return nil
}

func (c commentDataAccessor) WithDatabase(database Database) CommentDataAccessor {
	return &commentDataAccessor{
		database: database,
	}
}
