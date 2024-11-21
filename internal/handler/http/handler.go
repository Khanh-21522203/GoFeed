package http

import (
	"GoFeed/internal/generated/api/go_feed"
	grpc_handle "GoFeed/internal/handler/grpc"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type HttpHandler interface {
	RegisterRoutes(mux *http.ServeMux)
}

type httpHandler struct {
	grpcClient go_feed.GoFeedServiceClient
}

func NewHttpHandler() HttpHandler {
	return &httpHandler{}
}

func (h *httpHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/account", h.CreateAccount)
	mux.HandleFunc("/api/session", h.CreateSession)

	mux.HandleFunc("/api/post", h.CreatePost)
	mux.HandleFunc("/api/post/{post_id}", h.GetPostByID)
	mux.HandleFunc("/api/post/of_account/{account_id}", h.GetPostOfAccount)
	mux.HandleFunc("/api/post", h.UpdatePost)

	mux.HandleFunc("/api/like", h.CreateLike)
	mux.HandleFunc("/api/like/count/of_post/{post_id}", h.GetLikeCountOfPost)
	mux.HandleFunc("/api/like/account/of_post/{post_id}", h.GetLikeAccountsOfPost)
	mux.HandleFunc("/api/like", h.DeleteLike)

	mux.HandleFunc("/api/comment", h.CreateComment)
	mux.HandleFunc("/api/comment/count/of_post/{post_id}", h.GetCommentCountOfPost)
	mux.HandleFunc("/api/comment/of_post/{post_id}", h.GetCommentsOfPost)
	mux.HandleFunc("/api/comment", h.UpdateComment)
	mux.HandleFunc("/api/comment", h.DeleteComment)

	mux.HandleFunc("/api/follow", h.CreateFollow)
	mux.HandleFunc("/api/follow/follower_count/account", h.GetFollowerCountOfAccount)
	mux.HandleFunc("/api/follow/follower/account", h.GetFollowersOfAccount)
	mux.HandleFunc("/api/follow/following_count/account", h.GetFollowingCountOfAccount)
	mux.HandleFunc("/api/follow/following/account", h.GetFollowingsOfAccount)
	mux.HandleFunc("/api/follow", h.DeleteFollow)

	mux.HandleFunc("/api/new_feed", h.GetNewFeeds)
}

func (h httpHandler) httpToGRPCMetadata(r *http.Request) metadata.MD {
	md := metadata.MD{}
	for key, values := range r.Header {
		lowerKey := strings.ToLower(key)
		for _, value := range values {
			md.Append(lowerKey, value)
		}
	}
	return md
}

func (h httpHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	accountName, ok := body["account_name"].(string)
	if !ok || accountName == "" {
		WriteError(w, http.StatusBadRequest, "account_name is required and must be a uint64")
		return
	}
	password, ok := body["password"].(string)
	if !ok || password == "" {
		WriteError(w, http.StatusBadRequest, "password is required and must be a uint64")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.CreateAccount(ctx, &go_feed.CreateAccountRequest{
		AccountName: accountName,
		Password:    password,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	accountName, ok := body["account_name"].(string)
	if !ok || accountName == "" {
		WriteError(w, http.StatusBadRequest, "account_name is required and must be a uint64")
		return
	}
	password, ok := body["password"].(string)
	if !ok || password == "" {
		WriteError(w, http.StatusBadRequest, "password is required and must be a uint64")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	var header metadata.MD

	output, err := h.grpcClient.CreateSession(ctx, &go_feed.CreateSessionRequest{
		AccountName: accountName,
		Password:    password,
	}, grpc.Header(&header))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token := header.Get(grpc_handle.AuthTokenMetadataName)
	if len(token) == 0 {
		WriteError(w, http.StatusInternalServerError, "token not found in gRPC metadata")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"accountId": output.AccountId,
		"token":     token[0],
	})
}

func (h httpHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	content, ok := body["content"].(string)
	if !ok || content == "" {
		WriteError(w, http.StatusBadRequest, "account_name is required and must be a uint64")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.CreatePost(ctx, &go_feed.CreatePostRequest{
		Content: content,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
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

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.GetPostByID(ctx, &go_feed.GetPostByIDRequest{
		PostId: postID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) GetPostOfAccount(w http.ResponseWriter, r *http.Request) {
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

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.GetPostOfAccount(ctx, &go_feed.GetPostOfAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	postId, ok := body["post_id"].(uint64)
	if !ok || postId == 0 {
		WriteError(w, http.StatusBadRequest, "post_id is required and must be a uint64")
		return
	}
	content, ok := body["content"].(string)
	if !ok || content == "" {
		WriteError(w, http.StatusBadRequest, "content is required and must be a uint64")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.UpdatePost(ctx, &go_feed.UpdatePostRequest{
		Post: &go_feed.Post{
			Id:      postId,
			Content: content,
		},
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, output)
}

func (h httpHandler) CreateLike(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	postId, ok := body["post_id"].(uint64)
	if !ok || postId == 0 {
		WriteError(w, http.StatusBadRequest, "post_id is required and must be a uint64")
		return
	}
	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.CreateLike(ctx, &go_feed.CreateLikeRequest{
		PostId: postId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) GetLikeCountOfPost(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("post_id")
	if postIdStr == "" {
		WriteError(w, http.StatusBadRequest, "post_id is required")
		return
	}
	postId, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "post_id is invalid")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.GetLikeCountOfPost(ctx, &go_feed.GetLikeCountOfPostRequest{
		PostId: postId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) GetLikeAccountsOfPost(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("post_id")
	if postIdStr == "" {
		WriteError(w, http.StatusBadRequest, "post_id is required")
		return
	}
	postId, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "post_id is invalid")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.GetLikeAccountsOfPost(ctx, &go_feed.GetLikeAccountsOfPostRequest{
		PostId: postId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) DeleteLike(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	postId, ok := body["post_id"].(uint64)
	if !ok || postId == 0 {
		WriteError(w, http.StatusBadRequest, "post_id is required and must be a uint64")
		return
	}
	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.DeleteLike(ctx, &go_feed.DeleteLikeRequest{
		PostId: postId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h httpHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	postId, ok := body["post_id"].(uint64)
	if !ok || postId == 0 {
		WriteError(w, http.StatusBadRequest, "post_id is required and must be a uint64")
		return
	}
	content, ok := body["content"].(string)
	if !ok || content == "" {
		WriteError(w, http.StatusBadRequest, "content is required and must be a string")
		return
	}
	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.CreateComment(ctx, &go_feed.CreateCommentRequest{
		PostId:  postId,
		Content: content,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) GetCommentCountOfPost(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("post_id")
	if postIdStr == "" {
		WriteError(w, http.StatusBadRequest, "post_id is required")
		return
	}
	postId, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "post_id is invalid")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.GetCommentCountOfPost(ctx, &go_feed.GetCommentCountOfPostRequest{
		PostId: postId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) GetCommentsOfPost(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("post_id")
	if postIdStr == "" {
		WriteError(w, http.StatusBadRequest, "post_id is required")
		return
	}
	postId, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "post_id is invalid")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.GetCommentsOfPost(ctx, &go_feed.GetCommentsOfPostRequest{
		PostId: postId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	commentId, ok := body["comment_id"].(uint64)
	if !ok || commentId == 0 {
		WriteError(w, http.StatusBadRequest, "comment_id is required and must be a uint64")
		return
	}
	content, ok := body["content"].(string)
	if !ok || content == "" {
		WriteError(w, http.StatusBadRequest, "content is required and must be a string")
		return
	}
	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.UpdateComment(ctx, &go_feed.UpdateCommentRequest{
		Comment: &go_feed.Comment{
			CommentId: commentId,
			Content:   content,
		},
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	commentId, ok := body["comment_id"].(uint64)
	if !ok || commentId == 0 {
		WriteError(w, http.StatusBadRequest, "comment_id is required and must be a uint64")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.DeleteComment(ctx, &go_feed.DeleteCommentRequest{
		CommentId: commentId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h httpHandler) CreateFollow(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	followingId, ok := body["following_id"].(uint64)
	if !ok || followingId == 0 {
		WriteError(w, http.StatusBadRequest, "following_id is required and must be a uint64")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.CreateFollow(ctx, &go_feed.CreateFollowRequest{
		FollowingId: followingId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) GetFollowerCountOfAccount(w http.ResponseWriter, r *http.Request) {
	accountIdStr := r.URL.Query().Get("account_id")
	if accountIdStr == "" {
		WriteError(w, http.StatusBadRequest, "account_id is required")
		return
	}
	accountId, err := strconv.ParseUint(accountIdStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "account_id is invalid")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.GetFollowerCountOfAccount(ctx, &go_feed.GetFollowerCountOfAccountRequest{
		AccountId: accountId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) GetFollowersOfAccount(w http.ResponseWriter, r *http.Request) {
	accountIdStr := r.URL.Query().Get("account_id")
	if accountIdStr == "" {
		WriteError(w, http.StatusBadRequest, "account_id is required")
		return
	}
	accountId, err := strconv.ParseUint(accountIdStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "account_id is invalid")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.GetFollowersOfAccount(ctx, &go_feed.GetFollowersOfAccountRequest{
		AccountId: accountId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) GetFollowingCountOfAccount(w http.ResponseWriter, r *http.Request) {
	accountIdStr := r.URL.Query().Get("account_id")
	if accountIdStr == "" {
		WriteError(w, http.StatusBadRequest, "account_id is required")
		return
	}
	accountId, err := strconv.ParseUint(accountIdStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "account_id is invalid")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.GetFollowingCountOfAccount(ctx, &go_feed.GetFollowingCountOfAccountRequest{
		AccountId: accountId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) GetFollowingsOfAccount(w http.ResponseWriter, r *http.Request) {
	accountIdStr := r.URL.Query().Get("account_id")
	if accountIdStr == "" {
		WriteError(w, http.StatusBadRequest, "account_id is required")
		return
	}
	accountId, err := strconv.ParseUint(accountIdStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "account_id is invalid")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.GetFollowingsOfAccount(ctx, &go_feed.GetFollowingsOfAccountRequest{
		AccountId: accountId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}
func (h httpHandler) DeleteFollow(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	followingId, ok := body["following_id"].(uint64)
	if !ok || followingId == 0 {
		WriteError(w, http.StatusBadRequest, "following_id is required and must be a uint64")
		return
	}

	mData := h.httpToGRPCMetadata(r)
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := h.grpcClient.DeleteFollow(ctx, &go_feed.DeleteFollowRequest{
		FollowingId: followingId,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, output)
}

func (h httpHandler) GetNewFeeds(w http.ResponseWriter, r *http.Request) {}
