package logic

import (
	"GoFeed/internal/configs"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HashLogic interface {
	Hash(ctx context.Context, data string) (string, error)
	IsHashEqual(ctx context.Context, hashed string, password string) (bool, error)
}

type hash struct {
	hashConfig configs.Hash
}

func NewHash(hashConfig configs.Hash) HashLogic {
	return &hash{
		hashConfig: hashConfig,
	}
}

func (h hash) Hash(_ context.Context, data string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(data), h.hashConfig.Cost)
	if err != nil {
		return "", status.Error(codes.Internal, "failed to hash data")
	}

	return string(hashed), nil
}

func (h hash) IsHashEqual(_ context.Context, hashed string, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, status.Error(codes.Internal, "failed to check if data equal hash")
	}

	return true, nil
}
