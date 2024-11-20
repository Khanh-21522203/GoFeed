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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateCommentParams struct {
	Token   string
	PostID  uint64
	Content string
}
type CreateCommentOutput struct {
	ID uint64
}
type GetCommentCountOfPostParams struct {
	Token  string
	PostID uint64
}
type GetCommentCountOfPostOutput struct {
	CommentCount int
}
type GetCommentsOfPostParams struct {
	Token  string
	PostID uint64
}
type GetCommentsOfPostOutput struct {
	CommentList []*go_feed.Comment
}
type UpdateCommentParams struct {
	Token   string
	ID      uint64
	Content string
}
type UpdateCommentOutput struct {
}
type DeleteCommentParams struct {
	Token string
	ID    uint64
}
type DeleteCommentOutput struct {
}

type CommentLogic interface {
	CreateComment(ctx context.Context, params CreateCommentParams) (CreateCommentOutput, error)
	GetCommentCountOfPost(ctx context.Context, params GetCommentCountOfPostParams) (GetCommentCountOfPostOutput, error)
	GetCommentsOfPost(ctx context.Context, params GetCommentsOfPostParams) (GetCommentsOfPostOutput, error)
	UpdateComment(ctx context.Context, params UpdateCommentParams) error
	DeleteComment(ctx context.Context, params DeleteCommentParams) error
}

type commentLogic struct {
	goquDatabase        *goqu.Database
	commentDataAccessor database.CommentDataAccessor
	tokenLogic          TokenLogic
	idGenerator         *snowNode
	logger              *zap.Logger
}

func NewCommentLogic(
	goquDatabase *goqu.Database,
	commentDataAccessor database.CommentDataAccessor,
	tokenLogic TokenLogic,
	idGenerator *snowNode,
	logger *zap.Logger,
) CommentLogic {
	return &commentLogic{
		goquDatabase:        goquDatabase,
		commentDataAccessor: commentDataAccessor,
		tokenLogic:          tokenLogic,
		idGenerator:         idGenerator,
		logger:              logger,
	}
}

func (c commentLogic) databaseCommentToProtoComment(comment database.Comment) *go_feed.Comment {
	return &go_feed.Comment{
		CommentId: comment.ID,
		AccountId: comment.AccountID,
		PostId:    comment.PostID,
		Content:   comment.Content,
		CreatedAt: timestamppb.New(comment.CreatedAt),
	}
}

func (c commentLogic) CreateComment(ctx context.Context, params CreateCommentParams) (CreateCommentOutput, error) {
	accountID, _, err := c.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return CreateCommentOutput{}, err
	}
	commentID := c.idGenerator.GenID()
	txErr := c.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		commentID, err = c.commentDataAccessor.WithDatabase(td).CreateComment(ctx, database.Comment{
			ID:        commentID,
			AccountID: accountID,
			PostID:    params.PostID,
			Content:   params.Content,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return CreateCommentOutput{}, txErr
	}
	return CreateCommentOutput{
		ID: commentID,
	}, nil
}
func (c commentLogic) GetCommentCountOfPost(ctx context.Context, params GetCommentCountOfPostParams) (GetCommentCountOfPostOutput, error) {
	_, _, err := c.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetCommentCountOfPostOutput{}, err
	}
	commentCount, err := c.commentDataAccessor.GetCommentCountOfPost(ctx, params.PostID)
	if err != nil {
		return GetCommentCountOfPostOutput{}, err
	}
	return GetCommentCountOfPostOutput{
		CommentCount: commentCount,
	}, nil
}
func (c commentLogic) GetCommentsOfPost(ctx context.Context, params GetCommentsOfPostParams) (GetCommentsOfPostOutput, error) {
	_, _, err := c.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetCommentsOfPostOutput{}, err
	}
	commentList, err := c.commentDataAccessor.GetCommentsOfPost(ctx, params.PostID)
	if err != nil {
		return GetCommentsOfPostOutput{}, err
	}
	return GetCommentsOfPostOutput{
		CommentList: lo.Map(commentList, func(item database.Comment, _ int) *go_feed.Comment {
			return c.databaseCommentToProtoComment(item)
		}),
	}, nil
}
func (c commentLogic) UpdateComment(ctx context.Context, params UpdateCommentParams) error {
	accountID, _, err := c.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return err
	}
	txErr := c.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		comment, err := c.commentDataAccessor.WithDatabase(td).GetCommentByIdWithXLock(ctx, params.ID)
		if err != nil {
			return err
		}
		if comment.AccountID != accountID {
			return status.Error(codes.PermissionDenied, "trying to update a comment the account does not own")
		}
		comment.Content = params.Content
		err = c.commentDataAccessor.WithDatabase(td).UpdateComment(ctx, comment)
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
func (c commentLogic) DeleteComment(ctx context.Context, params DeleteCommentParams) error {
	accountID, _, err := c.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return err
	}
	txErr := c.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		comment, err := c.commentDataAccessor.WithDatabase(td).GetCommentByIdWithXLock(ctx, params.ID)
		if err != nil {
			return err
		}
		if comment.AccountID != accountID {
			return status.Error(codes.PermissionDenied, "trying to delete a comment the account does not own")
		}
		err = c.commentDataAccessor.WithDatabase(td).DeleteComment(ctx, params.ID)
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
