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

type followHandler struct {
	clientPool *grpcClientPool
}

func NewFollowHandler(clientPool *grpcClientPool) *followHandler {
	return &followHandler{clientPool: clientPool}
}

func (h followHandler) CreateFollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	var body struct {
		FollowingID uint64 `json:"following_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.FollowingID == 0 {
		WriteError(w, http.StatusBadRequest, "following_id is required and must be a uint64")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.CreateFollow(ctx, &go_feed.CreateFollowRequest{
		FollowingId: body.FollowingID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h followHandler) GetFollowerCountOfAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	accountID, err := h.parseQueryParamUint64(r, "account_id")
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

	output, err := client.GetFollowerCountOfAccount(ctx, &go_feed.GetFollowerCountOfAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h followHandler) GetFollowersOfAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	accountID, err := h.parseQueryParamUint64(r, "account_id")
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

	output, err := client.GetFollowersOfAccount(ctx, &go_feed.GetFollowersOfAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h followHandler) GetFollowingCountOfAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	accountID, err := h.parseQueryParamUint64(r, "account_id")
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

	output, err := client.GetFollowingCountOfAccount(ctx, &go_feed.GetFollowingCountOfAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h followHandler) GetFollowingsOfAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	accountID, err := h.parseQueryParamUint64(r, "account_id")
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

	output, err := client.GetFollowingsOfAccount(ctx, &go_feed.GetFollowingsOfAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h followHandler) DeleteFollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	var body struct {
		FollowingID uint64 `json:"following_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.FollowingID == 0 {
		WriteError(w, http.StatusBadRequest, "following_id is required and must be a uint64")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.DeleteFollow(ctx, &go_feed.DeleteFollowRequest{
		FollowingId: body.FollowingID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

// Helper method to parse a uint64 query parameter
func (h followHandler) parseQueryParamUint64(r *http.Request, param string) (uint64, error) {
	paramValue := r.URL.Query().Get(param)
	if paramValue == "" {
		return 0, fmt.Errorf("%s is required", param)
	}
	return strconv.ParseUint(paramValue, 10, 64)
}
