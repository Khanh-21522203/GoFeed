package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameAccounts = goqu.T("accounts")
)

const (
	ColNameAccountsID          = "id"
	ColNameAccountsAccountName = "account_name"
	ColNameAccountsHasing      = "hasing"
)

type Account struct {
	ID           uint64
	Account_name string
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account_name string, hashing string) (uint64, error)
	GetAccountByID(ctx context.Context, id uint64) (Account, error)
	GetAccountByAccountName(ctx context.Context, account_name string) (Account, error)
	WithDatabase(database Database) AccountDataAccessor
}

type accountDataAccessor struct {
	database Database
}

func NewAccountDataAccessor(database *goqu.Database) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
	}
}

func (a *accountDataAccessor) CreateAccount(ctx context.Context, account_name string, hashing string) (uint64, error) {
	return 0, nil
}

func (a *accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	return Account{}, nil
}

func (a *accountDataAccessor) GetAccountByAccountName(ctx context.Context, account_name string) (Account, error) {
	return Account{}, nil
}

func (a *accountDataAccessor) WithDatabase(database Database) AccountDataAccessor {
	a.database = database
	return a
}
