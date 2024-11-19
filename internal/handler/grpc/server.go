package grpc

import (
	"GoFeed/internal/configs"
	"GoFeed/internal/utils"
	"context"
	"net"

	go_feed "GoFeed/internal/generated/api/go_feed"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	handler    go_feed.GoFeedServiceServer
	grpcConfig configs.GRPC
	logger     *zap.Logger
}

func NewServer(logger *zap.Logger) Server {
	return &server{
		logger: logger,
	}
}

func (s server) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, s.logger)

	listener, err := net.Listen("tcp", s.grpcConfig.Address)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to open tcp listener")
		return err
	}
	defer listener.Close()

	server := grpc.NewServer()
	go_feed.RegisterGoFeedServiceServer(server, s.handler)

	logger.With(zap.String("address", s.grpcConfig.Address)).Info("starting grpc server")
	return server.Serve(listener)
}
