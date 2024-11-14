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
	TabNameAccounts = goqu.T("accounts")

	ErrAccountNotFound = status.Error(codes.NotFound, "account not found")
)

const (
	ColNameAccountsID          = "id"
	ColNameAccountsAccountName = "account_name"
	ColNameAccountsHasing      = "hasing"
)

type Account struct {
	ID           uint64
	Account_name string
	Hashing      string
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account Account) (uint64, error)
	GetAccountByID(ctx context.Context, id uint64) (Account, error)
	GetAccountByAccountName(ctx context.Context, account_name string) (Account, error)
	WithDatabase(database Database) AccountDataAccessor
}

type accountDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewAccountDataAccessor(database *goqu.Database, logger *zap.Logger) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (a accountDataAccessor) CreateAccount(ctx context.Context, account Account) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	_, err := a.database.
		Insert(TabNameAccounts).
		Rows(goqu.Record{
			ColNameAccountsID:          account.ID,
			ColNameAccountsAccountName: account.Account_name,
			ColNameAccountsHasing:      account.Hashing,
		}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create account")
		return 0, err
	}
	// lastInsertedID, err := result.LastInsertId()
	// if err != nil {
	// 	a.logger.With(zap.Error(err)).Error("failed to get last inserted id")
	// 	return 0, err
	// }
	// return uint64(lastInsertedID), nil
	return account.ID, nil
}

func (a accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	account := Account{}
	found, err := a.database.
		From(TabNameAccounts).
		Where(goqu.C(ColNameAccountsID).Eq(id)).
		ScanStructContext(ctx, &account)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account by id")
		return Account{}, status.Error(codes.Internal, "failed to get account by id")
	}
	if !found {
		logger.Warn("cannot find account by id")
		return Account{}, ErrAccountNotFound
	}
	return account, nil
}

func (a accountDataAccessor) GetAccountByAccountName(ctx context.Context, account_name string) (Account, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	account := Account{}
	found, err := a.database.
		From(TabNameAccounts).
		Where(goqu.C(ColNameAccountsAccountName).Eq(account_name)).
		ScanStructContext(ctx, &account)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account by account name")
		return Account{}, status.Error(codes.Internal, "failed to get account by account name")
	}
	if !found {
		logger.Warn("cannot find account by id")
		return Account{}, ErrAccountNotFound
	}
	return account, nil
}

func (a accountDataAccessor) WithDatabase(database Database) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
	}
}
