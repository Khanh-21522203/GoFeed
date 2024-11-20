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
	logger   *zap.Logger
}

func NewTokenPublicKeyDataAccessor(database *goqu.Database, logger *zap.Logger) TokenPublicKeyDataAccessor {
	return &tokenPublicKeyDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (t tokenPublicKeyDataAccessor) CreatePublicKey(ctx context.Context, tokenPublicKey TokenPublicKey) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	_, err := t.database.
		Insert(TabNameTokenPublicKeys).
		Rows(goqu.Record{
			ColNameTokenPublicKeysID:        tokenPublicKey.ID,
			ColNameTokenPublicKeysPublicKey: tokenPublicKey.PublicKey,
		}).
		Executor().Exec()

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to insert token public key")
		return 0, status.Error(codes.Internal, "failed to insert token public key")
	}
	return 0, nil
}

func (t tokenPublicKeyDataAccessor) GetPublicKey(ctx context.Context, id uint64) (TokenPublicKey, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	var token TokenPublicKey
	err := t.database.
		From(TabNameTokenPublicKeys).
		Where(goqu.C(ColNameTokenPublicKeysID).Eq(id)).
		ScanValsContext(ctx, &token)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get token public key")
		return TokenPublicKey{}, status.Error(codes.Internal, "failed to get token public key")
	}
	return token, nil
}

func (t tokenPublicKeyDataAccessor) WithDatabase(database Database) TokenPublicKeyDataAccessor {
	return &tokenPublicKeyDataAccessor{
		database: database,
	}
}
