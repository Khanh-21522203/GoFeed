package http

import "net/http"

type newFeedHandler struct {
	clientPool *grpcClientPool
}

func NewNewFeedHandler(clientPool *grpcClientPool) *newFeedHandler {
	return &newFeedHandler{clientPool: clientPool}
}

func (h newFeedHandler) GetNewFeeds(w http.ResponseWriter, r *http.Request) {}
