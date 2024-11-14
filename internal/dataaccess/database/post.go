package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNamePosts = goqu.T("posts")
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
	CreatePost(ctx context.Context, account_id uint64, content string) (uint64, error)
	GetPostByID(ctx context.Context, id uint64) (Post, error)
	GetPostsOfAccount(ctx context.Context, account_id uint64) ([]uint64, error)
	UpdatePost(ctx context.Context, post Post) error
	DeletePost(ctx context.Context, id uint64) error
	WithDatabase(database Database) PostDataAccessor
}

type postDataAccessor struct {
	database Database
}

func NewPostDataAccessor(database *goqu.Database) PostDataAccessor {
	return &postDataAccessor{
		database: database,
	}
}

func (p *postDataAccessor) CreatePost(ctx context.Context, account_id uint64, content string) (uint64, error) {
	return 0, nil
}

func (p *postDataAccessor) GetPostByID(ctx context.Context, id uint64) (Post, error) {
	return Post{}, nil
}

func (p *postDataAccessor) GetPostsOfAccount(ctx context.Context, account_id uint64) ([]uint64, error) {
	return nil, nil
}

func (p *postDataAccessor) UpdatePost(ctx context.Context, post Post) error {
	return nil
}

func (p *postDataAccessor) DeletePost(ctx context.Context, id uint64) error {
	return nil
}

func (p *postDataAccessor) WithDatabase(database Database) PostDataAccessor {
	p.database = database
	return p
}
