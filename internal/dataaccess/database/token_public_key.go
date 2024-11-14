package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameTokenPublicKeys = goqu.T("token_public_keys")
)

const (
	ColNameTokenPublicKeysID        = "id"
	ColNameTokenPublicKeysPublicKey = "public_key"
)

type TokenPublicKey struct {
	ID        uint64
	PublicKey string
}

type TokenPublicKeyDataAccessor interface {
	CreatePublicKey(ctx context.Context, tokenPublicKey TokenPublicKey) (uint64, error)
	GetPublicKey(ctx context.Context, id uint64) (TokenPublicKey, error)
	WithDatabase(database Database) TokenPublicKeyDataAccessor
}

type tokenPublicKeyDataAccessor struct {
	database Database
}

func NewTokenPublicKeyDataAccessor(database *goqu.Database) TokenPublicKeyDataAccessor {
	return &tokenPublicKeyDataAccessor{
		database: database,
	}
}

func (t *tokenPublicKeyDataAccessor) CreatePublicKey(ctx context.Context, tokenPublicKey TokenPublicKey) (uint64, error) {
	return 0, nil
}

func (t *tokenPublicKeyDataAccessor) GetPublicKey(ctx context.Context, id uint64) (TokenPublicKey, error) {
	return TokenPublicKey{}, nil
}

func (t *tokenPublicKeyDataAccessor) WithDatabase(database Database) TokenPublicKeyDataAccessor {
	t.database = database
	return t
}
