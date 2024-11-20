package logic

import (
	"GoFeed/internal/dataaccess/database"
	"GoFeed/internal/generated/api/go_feed"
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreatePostParams struct {
	Token   string
	Content string
}
type CreatePostOutput struct {
	ID uint64
}
type GetPostByIDParams struct {
	Token string
	ID    uint64
}
type GetPostByIDOutput struct {
	Post *go_feed.Post
}
type GetPostOfAccountParams struct {
	Token      string
	Of_account uint64
}
type GetPostOfAccountOutput struct {
	PostList []*go_feed.Post
}
type UpdatePostParams struct {
	Token   string
	ID      uint64
	Content string
}
type UpdatePostOutput struct {
}
type DeletePostParams struct {
	Token string
	ID    uint64
}
type DeletePostOutput struct {
}

type PostLogic interface {
	CreatePost(ctx context.Context, params CreatePostParams) (CreatePostOutput, error)
	GetPostByID(ctx context.Context, params GetPostByIDParams) (GetPostByIDOutput, error)
	GetPostOfAccount(ctx context.Context, params GetPostOfAccountParams) (GetPostOfAccountOutput, error)
	UpdatePost(ctx context.Context, params UpdatePostParams) (UpdatePostOutput, error)
	DeletePost(ctx context.Context, params DeletePostParams) error
}

type postLogic struct {
	goquDatabase        *goqu.Database
	postDataAccessor    database.PostDataAccessor
	commentDataAccessor database.CommentDataAccessor
	idGenerator         *snowNode
	tokenLogic          TokenLogic
	logger              *zap.Logger
}

func NewPostLogic(
	goquDatabase *goqu.Database,
	postDataAccessor database.PostDataAccessor,
	commentDataAccessor database.CommentDataAccessor,
	idGenerator *snowNode,
	tokenLogic TokenLogic,
	logger *zap.Logger,
) PostLogic {
	return &postLogic{
		goquDatabase:        goquDatabase,
		postDataAccessor:    postDataAccessor,
		commentDataAccessor: commentDataAccessor,
		idGenerator:         idGenerator,
		tokenLogic:          tokenLogic,
		logger:              logger,
	}
}

func (p postLogic) databasePostToProtoPost(post database.Post) *go_feed.Post {
	return &go_feed.Post{
		Id:        post.ID,
		AccountId: post.AccountID,
		Content:   post.Content,
	}
}
func (p postLogic) CreatePost(ctx context.Context, params CreatePostParams) (CreatePostOutput, error) {
	// Authorization -> Insert DB
	accountID, _, err := p.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return CreatePostOutput{}, err
	}
	postID := p.idGenerator.GenID()
	txErr := p.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		postID, err = p.postDataAccessor.WithDatabase(td).CreatePost(ctx, database.Post{
			ID:        postID,
			AccountID: accountID,
			Content:   params.Content,
		})
		if err != nil {
			return err
		}
		// Producer to Kafka
		return nil
	})
	if txErr != nil {
		return CreatePostOutput{}, txErr
	}
	return CreatePostOutput{
		ID: postID,
	}, nil

}
func (p postLogic) GetPostByID(ctx context.Context, params GetPostByIDParams) (GetPostByIDOutput, error) {
	_, _, err := p.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetPostByIDOutput{}, err
	}
	post, err := p.postDataAccessor.GetPostByID(ctx, params.ID)
	if err != nil {
		return GetPostByIDOutput{}, err
	}
	return GetPostByIDOutput{
		p.databasePostToProtoPost(post),
	}, nil
}

func (p postLogic) GetPostOfAccount(ctx context.Context, params GetPostOfAccountParams) (GetPostOfAccountOutput, error) {
	_, _, err := p.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetPostOfAccountOutput{}, err
	}
	postList, err := p.postDataAccessor.GetPostsOfAccount(ctx, params.Of_account)
	if err != nil {
		return GetPostOfAccountOutput{}, err
	}
	return GetPostOfAccountOutput{
		PostList: lo.Map(postList, func(item database.Post, _ int) *go_feed.Post {
			return p.databasePostToProtoPost(item)
		}),
	}, nil
}
func (p postLogic) UpdatePost(ctx context.Context, params UpdatePostParams) (UpdatePostOutput, error) {
	account_id, _, err := p.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return UpdatePostOutput{}, err
	}
	txErr := p.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		post, err := p.postDataAccessor.WithDatabase(td).GetPostByIDWithXLock(ctx, params.ID)
		if err != nil {
			return err
		}
		if account_id != post.AccountID {
			return status.Error(codes.PermissionDenied, "trying to update a post the account does not own")
		}

		post.Content = params.Content
		err = p.postDataAccessor.WithDatabase(td).UpdatePost(ctx, post)
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return UpdatePostOutput{}, txErr
	}
	return UpdatePostOutput{}, nil
}
func (p postLogic) DeletePost(ctx context.Context, params DeletePostParams) error {
	account_id, _, err := p.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return err
	}
	txErr := p.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		post, err := p.postDataAccessor.WithDatabase(td).GetPostByIDWithXLock(ctx, params.ID)
		if err != nil {
			return err
		}
		if account_id != post.AccountID {
			return status.Error(codes.PermissionDenied, "trying to delete a post the account does not own")
		}
		err = p.commentDataAccessor.WithDatabase(td).DeleteCommentOfPost(ctx, params.ID)
		if err != nil {
			return err
		}
		err = p.postDataAccessor.WithDatabase(td).DeletePost(ctx, params.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}
	return nil
}
