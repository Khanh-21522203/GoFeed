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

type commentHandler struct {
	clientPool *grpcClientPool
}

func NewCommentHandler(clientPool *grpcClientPool) *commentHandler {
	return &commentHandler{clientPool: clientPool}
}

func (h *commentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	var body struct {
		PostID  uint64 `json:"post_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.PostID == 0 || body.Content == "" {
		WriteError(w, http.StatusBadRequest, "post_id and content are required")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.CreateComment(ctx, &go_feed.CreateCommentRequest{
		PostId:  body.PostID,
		Content: body.Content,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h *commentHandler) GetCommentCountOfPost(w http.ResponseWriter, r *http.Request) {
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

	output, err := client.GetCommentCountOfPost(ctx, &go_feed.GetCommentCountOfPostRequest{
		PostId: postID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h *commentHandler) GetCommentsOfPost(w http.ResponseWriter, r *http.Request) {
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

	output, err := client.GetCommentsOfPost(ctx, &go_feed.GetCommentsOfPostRequest{
		PostId: postID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h *commentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	var body struct {
		CommentID uint64 `json:"comment_id"`
		Content   string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.CommentID == 0 || body.Content == "" {
		WriteError(w, http.StatusBadRequest, "comment_id and content are required")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.UpdateComment(ctx, &go_feed.UpdateCommentRequest{
		Comment: &go_feed.Comment{
			CommentId: body.CommentID,
			Content:   body.Content,
		},
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h *commentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	var body struct {
		CommentID uint64 `json:"comment_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.CommentID == 0 {
		WriteError(w, http.StatusBadRequest, "comment_id is required and must be a uint64")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get gRPC client: "+err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.DeleteComment(ctx, &go_feed.DeleteCommentRequest{
		CommentId: body.CommentID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

// Helper method to parse a uint64 query parameter
func (h *commentHandler) parseQueryParamUint64(r *http.Request, param string) (uint64, error) {
	paramValue := r.URL.Query().Get(param)
	if paramValue == "" {
		return 0, fmt.Errorf("%s is required", param)
	}
	return strconv.ParseUint(paramValue, 10, 64)
}
