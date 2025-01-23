package http

import (
	"net/http"
)

type HttpHandler interface {
	RegisterRoutes(mux *http.ServeMux)
}

type httpHandler struct {
	accountHandler
	postHandler
	likeHandler
	commentHandler
	followHandler
	newFeedHandler
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
