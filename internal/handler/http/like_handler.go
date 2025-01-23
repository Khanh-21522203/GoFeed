package http

import (
	"GoFeed/internal/generated/api/go_feed"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/grpc/metadata"
)

type likeHandler struct {
	clientPool *grpcClientPool
}

func NewLikeHandler(clientPool *grpcClientPool) *likeHandler {
	return &likeHandler{clientPool: clientPool}
}

func (h *likeHandler) CreateLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	var body struct {
		PostID uint64 `json:"post_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.PostID == 0 {
		WriteError(w, http.StatusBadRequest, "post_id is required and must be a uint64")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.CreateLike(ctx, &go_feed.CreateLikeRequest{
		PostId: body.PostID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h *likeHandler) GetLikeCountOfPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	postID, err := h.parseQueryParamUint64(r, "post_id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.GetLikeCountOfPost(ctx, &go_feed.GetLikeCountOfPostRequest{
		PostId: postID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h *likeHandler) GetLikeAccountsOfPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	postID, err := h.parseQueryParamUint64(r, "post_id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.GetLikeAccountsOfPost(ctx, &go_feed.GetLikeAccountsOfPostRequest{
		PostId: postID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h *likeHandler) DeleteLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	var body struct {
		PostID uint64 `json:"post_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.PostID == 0 {
		WriteError(w, http.StatusBadRequest, "post_id is required and must be a uint64")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.DeleteLike(ctx, &go_feed.DeleteLikeRequest{
		PostId: body.PostID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

// Helper method to parse a uint64 query parameter
func (h *likeHandler) parseQueryParamUint64(r *http.Request, param string) (uint64, error) {
	paramValue := r.URL.Query().Get(param)
	if paramValue == "" {
		return 0, fmt.Errorf("%s is required", param)
	}
	return strconv.ParseUint(paramValue, 10, 64)
}
