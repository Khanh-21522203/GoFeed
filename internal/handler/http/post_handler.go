package http

import (
	"GoFeed/internal/generated/api/go_feed"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"google.golang.org/grpc/metadata"
)

type postHandler struct {
	clientPool *grpcClientPool
}

func NewpostHandler(clientPool *grpcClientPool) *postHandler {
	return &postHandler{clientPool: clientPool}
}

func (h postHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if body.Content == "" {
		WriteError(w, http.StatusBadRequest, "content is required and must be a non-empty string")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.CreatePost(ctx, &go_feed.CreatePostRequest{
		Content: body.Content,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create post: "+err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, output)
}

func (h *postHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		WriteError(w, http.StatusBadRequest, "post_id is required")
		return
	}
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "post_id is invalid")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.GetPostByID(ctx, &go_feed.GetPostByIDRequest{
		PostId: postID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get post: "+err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, output)
}

func (h postHandler) GetPostOfAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	accountIDStr := r.URL.Query().Get("account_id")
	if accountIDStr == "" {
		WriteError(w, http.StatusBadRequest, "account_id is required")
		return
	}
	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "account_id is invalid")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.GetPostOfAccount(ctx, &go_feed.GetPostOfAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get posts: "+err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, output)
}

func (h postHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	var body struct {
		PostID  uint64 `json:"post_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if body.PostID == 0 {
		WriteError(w, http.StatusBadRequest, "post_id is required and must be a non-zero uint64")
		return
	}
	if body.Content == "" {
		WriteError(w, http.StatusBadRequest, "content is required and must be a non-empty string")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.UpdatePost(ctx, &go_feed.UpdatePostRequest{
		Post: &go_feed.Post{
			Id:      body.PostID,
			Content: body.Content,
		},
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update post: "+err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, output)
}
