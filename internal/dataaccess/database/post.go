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
	TabNamePosts = goqu.T("posts")

	ErrPostNotFound = status.Error(codes.NotFound, "post not found")
	// ErrAccountHasNoPost = status.Error(codes.NotFound, "account has no post")
)

const (
	ColNamePostsID        = "id"
	ColNamePostsAccountID = "account_id"
	ColNamePostContent    = "content"
)

type Post struct {
	ID        uint64
	AccountID uint64
	Content   string
}

type PostDataAccessor interface {
	CreatePost(ctx context.Context, post Post) (uint64, error)
	GetPostByID(ctx context.Context, id uint64) (Post, error)
	GetPostsOfAccount(ctx context.Context, account_id uint64) ([]uint64, error)
	UpdatePost(ctx context.Context, post Post) error
	DeletePost(ctx context.Context, id uint64) error
	WithDatabase(database Database) PostDataAccessor
}

type postDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewPostDataAccessor(database *goqu.Database, logger *zap.Logger) PostDataAccessor {
	return &postDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (p postDataAccessor) CreatePost(ctx context.Context, post Post) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, p.logger)

	_, err := p.database.
		Insert(TabNamePosts).
		Rows(goqu.Record{
			ColNamePostsID:        post.ID,
			ColNamePostsAccountID: post.AccountID,
			ColNamePostContent:    post.Content,
		}).
		Executor().
		ExecContext(ctx)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create post")
		return 0, status.Error(codes.Internal, "failed to create post")
	}
	return 0, nil
}

func (p postDataAccessor) GetPostByID(ctx context.Context, id uint64) (Post, error) {
	logger := utils.LoggerWithContext(ctx, p.logger)

	post := Post{}
	found, err := p.database.
		From(TabNamePosts).
		Where(goqu.C(ColNamePostsID).Eq(id)).
		ScanStructContext(ctx, &post)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get post by id")
		return Post{}, status.Error(codes.Internal, "failed to get post by id")
	}
	if !found {
		logger.Warn("cannot find post by id")
		return Post{}, ErrPostNotFound
	}
	return post, nil
}

func (p postDataAccessor) GetPostsOfAccount(ctx context.Context, account_id uint64) ([]uint64, error) {
	logger := utils.LoggerWithContext(ctx, p.logger)

	var posts []uint64
	err := p.database.
		Select(ColNamePostsID).
		From(TabNamePosts).
		Where(goqu.C(ColNamePostsAccountID).Eq(account_id)).
		ScanValsContext(ctx, &posts)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get posts of account")
		return nil, status.Error(codes.Internal, "failed to get posts of account")
	}

	return posts, nil
}

func (p postDataAccessor) UpdatePost(ctx context.Context, post Post) error {
	logger := utils.LoggerWithContext(ctx, p.logger)

	_, err := p.database.
		Update(TabNamePosts).
		Set(post).
		Where(goqu.Ex{ColNamePostsID: post.ID}).
		Executor().
		ExecContext(ctx)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to update post")
		return status.Error(codes.Internal, "failed to update post")
	}
	return nil
}

func (p postDataAccessor) DeletePost(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, p.logger)

	_, err := p.database.
		Delete(TabNamePosts).
		Where(goqu.Ex{ColNamePostsID: id}).
		Executor().
		ExecContext(ctx)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to delete post")
		return status.Error(codes.Internal, "failed to delete post")
	}
	return nil
}

func (p postDataAccessor) WithDatabase(database Database) PostDataAccessor {
	return &postDataAccessor{
		database: database,
	}
}
