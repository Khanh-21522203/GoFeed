package logic

import (
	"GoFeed/internal/dataaccess/database"
	"GoFeed/internal/generated/api/go_feed"
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type CreateLikeParams struct {
	Token  string
	PostID uint64
}
type CreateLikeOutput struct{}
type GetLikeCountOfPostParams struct {
	Token  string
	PostID uint64
}
type GetLikeCountOfPostOutput struct {
	LikeCount int
}
type GetLikeAccountsOfPostParams struct {
	Token  string
	PostID uint64
}
type GetLikeAccountsOfPostOutput struct {
	AccountList []*go_feed.Account
}
type DeleteLikeParams struct {
	Token  string
	PostID uint64
}
type DeleteLikeOutput struct{}

type LikeLogic interface {
	CreateLike(ctx context.Context, params CreateLikeParams) error
	GetLikeCountOfPost(ctx context.Context, params GetLikeCountOfPostParams) (GetLikeCountOfPostOutput, error)
	GetLikeAccountsOfPost(ctx context.Context, params GetLikeAccountsOfPostParams) (GetLikeAccountsOfPostOutput, error)
	DeleteLike(ctx context.Context, params DeleteLikeParams) error
}

type likeLogic struct {
	goquDatabase        *goqu.Database
	likeDataAccessor    database.LikeDataAccessor
	accountDataAccessor database.AccountDataAccessor
	tokenLogic          TokenLogic
	logger              *zap.Logger
}

func NewLikeLogic(
	goquDatabase *goqu.Database,
	likeDataAccessor database.LikeDataAccessor,
	accountDataAccessor database.AccountDataAccessor,
	tokenLogic TokenLogic,
	logger *zap.Logger,
) LikeLogic {
	return &likeLogic{
		goquDatabase:        goquDatabase,
		likeDataAccessor:    likeDataAccessor,
		accountDataAccessor: accountDataAccessor,
		tokenLogic:          tokenLogic,
		logger:              logger,
	}
}

func (l likeLogic) CreateLike(ctx context.Context, params CreateLikeParams) error {
	accountID, _, err := l.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return err
	}
	txErr := l.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		err := l.likeDataAccessor.WithDatabase(td).CreateLike(ctx, database.Like{
			AccountID: accountID,
			PostID:    params.PostID,
		})
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
func (l likeLogic) GetLikeCountOfPost(ctx context.Context, params GetLikeCountOfPostParams) (GetLikeCountOfPostOutput, error) {
	_, _, err := l.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetLikeCountOfPostOutput{}, err
	}
	likeCount, err := l.likeDataAccessor.GetLikeCountOfPost(ctx, params.PostID)
	if err != nil {
		return GetLikeCountOfPostOutput{}, err
	}
	return GetLikeCountOfPostOutput{
		LikeCount: likeCount,
	}, nil
}
func (l likeLogic) GetLikeAccountsOfPost(ctx context.Context, params GetLikeAccountsOfPostParams) (GetLikeAccountsOfPostOutput, error) {
	_, _, err := l.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetLikeAccountsOfPostOutput{}, err
	}
	accountIdList, err := l.likeDataAccessor.GetLikeAccountsOfPost(ctx, params.PostID)
	if err != nil {
		return GetLikeAccountsOfPostOutput{}, err
	}
	accountList, err := l.accountDataAccessor.GetAccountByIDs(ctx, accountIdList)
	if err != nil {
		return GetLikeAccountsOfPostOutput{}, err
	}
	return GetLikeAccountsOfPostOutput{
		AccountList: lo.Map(accountList, func(item database.Account, _ int) *go_feed.Account {
			return &go_feed.Account{
				Id:          item.ID,
				AccountName: item.Account_name,
			}
		}),
	}, nil
}
func (l likeLogic) DeleteLike(ctx context.Context, params DeleteLikeParams) error {
	accountID, _, err := l.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return err
	}
	txErr := l.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		err := l.likeDataAccessor.WithDatabase(td).DeleteLike(ctx, accountID, params.PostID)
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
