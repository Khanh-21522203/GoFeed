package logic

import (
	"GoFeed/internal/dataaccess/database"
	"GoFeed/internal/generated/api/go_feed"
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateAccountParams struct {
	AccountName string
	Password    string
}

type CreateAccountOutput struct {
	ID          uint64
	AccountName string
}

type CreateSessionParams struct {
	AccountName string
	Password    string
}

type CreateSessionOutput struct {
	Account *go_feed.Account
	Token   string
}

type AccountLogic interface {
	CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error)
	CreateSession(ctx context.Context, params CreateSessionParams) (CreateSessionOutput, error)
}

type accountLogic struct {
	goquDatabase        *goqu.Database
	accountDataAccessor database.AccountDataAccessor
	hashLogic           HashLogic
	tokenLogic          TokenLogic
	idGenerator         *snowNode
	logger              *zap.Logger
}

func NewAccountLogic(
	goquDatabase *goqu.Database,
	accountDataAccessor database.AccountDataAccessor,
	hashLogic HashLogic,
	tokenLogic TokenLogic,
	idGenerator *snowNode,
	logger *zap.Logger,
) AccountLogic {
	return &accountLogic{
		goquDatabase:        goquDatabase,
		accountDataAccessor: accountDataAccessor,
		hashLogic:           hashLogic,
		tokenLogic:          tokenLogic,
		idGenerator:         idGenerator,
		logger:              logger,
	}
}

func (a accountLogic) isAccountNameExist(ctx context.Context, account_name string) (bool, error) {
	_, err := a.accountDataAccessor.GetAccountByAccountName(ctx, account_name)
	if err != nil {
		if errors.Is(err, database.ErrAccountNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a accountLogic) databaseAccountToProtoAccount(account database.Account) *go_feed.Account {
	return &go_feed.Account{
		Id:          account.ID,
		AccountName: account.Account_name,
	}
}

func (a accountLogic) CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error) {
	// Check if account_name is exist -> hasing password -> add to DB
	accountNameExist, err := a.isAccountNameExist(ctx, params.AccountName)
	if err != nil {
		return CreateAccountOutput{}, status.Error(codes.Internal, "failed to check if account name is exist")
	}
	if accountNameExist {
		return CreateAccountOutput{}, status.Error(codes.AlreadyExists, "account name is already taken")
	}

	accountID := a.idGenerator.GenID()
	txErr := a.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {

		hashedPassword, hashErr := a.hashLogic.Hash(ctx, params.Password)
		if hashErr != nil {
			return hashErr
		}

		accountID, err = a.accountDataAccessor.WithDatabase(td).CreateAccount(ctx, database.Account{
			ID:           accountID,
			Account_name: params.AccountName,
			Hashing:      hashedPassword,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return CreateAccountOutput{}, txErr
	}

	return CreateAccountOutput{
		ID:          accountID,
		AccountName: params.AccountName,
	}, nil
}
func (a accountLogic) CreateSession(ctx context.Context, params CreateSessionParams) (CreateSessionOutput, error) {
	// Check if account name is exist -> Check if password equal hashing -> Get JWT token -> return ouput
	existingAccount, err := a.accountDataAccessor.GetAccountByAccountName(ctx, params.AccountName)
	if err != nil {
		return CreateSessionOutput{}, err
	}

	isHashEqual, err := a.hashLogic.IsHashEqual(ctx, existingAccount.Hashing, params.Password)
	if err != nil {
		return CreateSessionOutput{}, err
	}
	if !isHashEqual {
		return CreateSessionOutput{}, status.Error(codes.Unauthenticated, "incorrect password")
	}

	token, _, err := a.tokenLogic.GetToken(ctx, existingAccount.ID)
	if err != nil {
		return CreateSessionOutput{}, err
	}
	return CreateSessionOutput{
		Account: a.databaseAccountToProtoAccount(existingAccount),
		Token:   token,
	}, nil
}
