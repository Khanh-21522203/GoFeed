package http

import (
	"GoFeed/internal/configs"
	"GoFeed/internal/utils"
	"context"
	"net/http"

	"go.uber.org/zap"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	grpcConfig configs.GRPC
	httpConfig configs.HTTP
	logger     *zap.Logger
}

func NewServer(grpcConfig configs.GRPC, httpConfig configs.HTTP, logger *zap.Logger) Server {
	return &server{
		grpcConfig: grpcConfig,
		httpConfig: httpConfig,
		logger:     logger,
	}
}

func (s server) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, s.logger)

	mux := http.NewServeMux()

	handler := NewHttpHandler()
	handler.registerRoutes(mux)

	logger.With(zap.String("address", s.httpConfig.Address)).Info("starting http server")
	return http.ListenAndServe(s.httpConfig.Address, mux)
}
