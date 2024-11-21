package producer

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GoFeed/internal/utils"
)

const (
	MessageQueueNewFeedJob = "new_feed_job"
)

type NewFeedJob struct {
	PostID uint64 `json:"id"`
}

type NewFeedJobProducer interface {
	Produce(ctx context.Context, event NewFeedJob) error
}

type newFeedJobProducer struct {
	client Client
	logger *zap.Logger
}

func NewDownloadTaskCreatedProducer(
	client Client,
	logger *zap.Logger,
) NewFeedJobProducer {
	return &newFeedJobProducer{
		client: client,
		logger: logger,
	}
}

func (n newFeedJobProducer) Produce(ctx context.Context, event NewFeedJob) error {
	logger := utils.LoggerWithContext(ctx, n.logger)

	eventBytes, err := json.Marshal(event)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to marshal new feed job event")
		return status.Error(codes.Internal, "failed to marshal new feed job event")
	}

	err = n.client.Produce(ctx, MessageQueueNewFeedJob, eventBytes)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to produce new feed job event")
		return status.Error(codes.Internal, "failed to produce new feed job event")
	}

	return nil
}
