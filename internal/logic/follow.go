package logic

import (
	"GoFeed/internal/dataaccess/database"
	"GoFeed/internal/generated/api/go_feed"
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type CreateFollowParams struct {
	Token       string
	FollowingID uint64
}
type CreateFollowOutput struct{}
type GetFollowerCountOfAccountParams struct {
	Token     string
	AccountID uint64
}
type GetFollowerCountOfAccountOutput struct {
	FollowerCount int
}
type GetFollowersOfAccountParams struct {
	Token     string
	AccountID uint64
}
type GetFollowersOfAccountOutput struct {
	FollowerList []*go_feed.Account
}
type GetFollowingCountOfAccountParams struct {
	Token     string
	AccountID uint64
}
type GetFollowingCountOfAccountOutput struct {
	FollowingCount int
}
type GetFollowingsOfAccountParams struct {
	Token     string
	AccountID uint64
}
type GetFollowingsOfAccountOutput struct {
	FollowingList []*go_feed.Account
}
type DeleteFollowParams struct {
	Token       string
	FollowingID uint64
}
type DeleteFollowOutput struct{}

type FollowLogic interface {
	CreateFollow(ctx context.Context, params CreateFollowParams) error
	GetFollowerCountOfAccount(ctx context.Context, params GetFollowerCountOfAccountParams) (GetFollowerCountOfAccountOutput, error)
	GetFollowersOfAccount(ctx context.Context, params GetFollowersOfAccountParams) (GetFollowersOfAccountOutput, error)
	GetFollowingCountOfAccount(ctx context.Context, params GetFollowingCountOfAccountParams) (GetFollowingCountOfAccountOutput, error)
	GetFollowingsOfAccount(ctx context.Context, params GetFollowingsOfAccountParams) (GetFollowingsOfAccountOutput, error)
	DeleteFollow(ctx context.Context, params DeleteFollowParams) error
}

type followLogic struct {
	goquDatabase        *goqu.Database
	followDataAccessor  database.FollowDataAccessor
	accountDataAccessor database.AccountDataAccessor
	tokenLogic          TokenLogic
	logger              *zap.Logger
}

func NewFollowLogic(
	goquDatabase *goqu.Database,
	followDataAccessor database.FollowDataAccessor,
	accountDataAccessor database.AccountDataAccessor,
	tokenLogic TokenLogic,
	logger *zap.Logger,
) FollowLogic {
	return &followLogic{
		goquDatabase:        goquDatabase,
		followDataAccessor:  followDataAccessor,
		accountDataAccessor: accountDataAccessor,
		tokenLogic:          tokenLogic,
		logger:              logger,
	}
}

func (f followLogic) CreateFollow(ctx context.Context, params CreateFollowParams) error {
	accountID, _, err := f.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return err
	}
	txErr := f.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		err := f.followDataAccessor.WithDatabase(td).CreateFollow(ctx, database.Follow{
			AccountID:   accountID,
			FollowingID: params.FollowingID,
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
func (f followLogic) GetFollowerCountOfAccount(ctx context.Context, params GetFollowerCountOfAccountParams) (GetFollowerCountOfAccountOutput, error) {
	_, _, err := f.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetFollowerCountOfAccountOutput{}, err
	}
	followerCount, err := f.followDataAccessor.GetFollowerCountOfAccount(ctx, params.AccountID)
	if err != nil {
		return GetFollowerCountOfAccountOutput{}, err
	}
	return GetFollowerCountOfAccountOutput{
		FollowerCount: followerCount,
	}, nil
}
func (f followLogic) GetFollowersOfAccount(ctx context.Context, params GetFollowersOfAccountParams) (GetFollowersOfAccountOutput, error) {
	_, _, err := f.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetFollowersOfAccountOutput{}, err
	}

	followerIdList, err := f.followDataAccessor.GetFollowersOfAccount(ctx, params.AccountID)
	if err != nil {
		return GetFollowersOfAccountOutput{}, err
	}

	accountList, err := f.accountDataAccessor.GetAccountByIDs(ctx, followerIdList)
	if err != nil {
		return GetFollowersOfAccountOutput{}, err
	}
	return GetFollowersOfAccountOutput{
		FollowerList: lo.Map(accountList, func(item database.Account, _ int) *go_feed.Account {
			return &go_feed.Account{
				Id:          item.ID,
				AccountName: item.Account_name,
			}
		}),
	}, nil
}
func (f followLogic) GetFollowingCountOfAccount(ctx context.Context, params GetFollowingCountOfAccountParams) (GetFollowingCountOfAccountOutput, error) {
	_, _, err := f.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetFollowingCountOfAccountOutput{}, err
	}
	followingCount, err := f.followDataAccessor.GetFollowingCountOfAccount(ctx, params.AccountID)
	if err != nil {
		return GetFollowingCountOfAccountOutput{}, err
	}
	return GetFollowingCountOfAccountOutput{
		FollowingCount: followingCount,
	}, nil
}
func (f followLogic) GetFollowingsOfAccount(ctx context.Context, params GetFollowingsOfAccountParams) (GetFollowingsOfAccountOutput, error) {
	_, _, err := f.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetFollowingsOfAccountOutput{}, err
	}
	followingIdList, err := f.followDataAccessor.GetFollowingsOfAccount(ctx, params.AccountID)
	if err != nil {
		return GetFollowingsOfAccountOutput{}, err
	}

	accountList, err := f.accountDataAccessor.GetAccountByIDs(ctx, followingIdList)
	if err != nil {
		return GetFollowingsOfAccountOutput{}, err
	}

	return GetFollowingsOfAccountOutput{
		FollowingList: lo.Map(accountList, func(item database.Account, _ int) *go_feed.Account {
			return &go_feed.Account{
				Id:          item.ID,
				AccountName: item.Account_name,
			}
		}),
	}, nil
}
func (f followLogic) DeleteFollow(ctx context.Context, params DeleteFollowParams) error {
	accountID, _, err := f.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return err
	}
	txErr := f.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		err := f.followDataAccessor.WithDatabase(td).DeleteFollow(ctx, database.Follow{
			AccountID:   accountID,
			FollowingID: params.FollowingID,
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
