package http

import (
	"net/http"
)

type handler struct {
}

func NewHttpHandler() *handler {
	return &handler{}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/account", h.CreateAccount)
	mux.HandleFunc("/api/session", h.CreateSession)

	mux.HandleFunc("/api/post", h.CreatePost)
	mux.HandleFunc("/api/post/{post_id}", h.GetPostByID)
	mux.HandleFunc("/api/post/of_account", h.GetPostOfAccount)
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

func (h *handler) CreateAccount(w http.ResponseWriter, r *http.Request) {}
func (h *handler) CreateSession(w http.ResponseWriter, r *http.Request) {}

func (h *handler) CreatePost(w http.ResponseWriter, r *http.Request)       {}
func (h *handler) GetPostByID(w http.ResponseWriter, r *http.Request)      {}
func (h *handler) GetPostOfAccount(w http.ResponseWriter, r *http.Request) {}
func (h *handler) UpdatePost(w http.ResponseWriter, r *http.Request)       {}

func (h *handler) CreateLike(w http.ResponseWriter, r *http.Request)            {}
func (h *handler) GetLikeCountOfPost(w http.ResponseWriter, r *http.Request)    {}
func (h *handler) GetLikeAccountsOfPost(w http.ResponseWriter, r *http.Request) {}
func (h *handler) DeleteLike(w http.ResponseWriter, r *http.Request)            {}

func (h *handler) CreateComment(w http.ResponseWriter, r *http.Request)         {}
func (h *handler) GetCommentCountOfPost(w http.ResponseWriter, r *http.Request) {}
func (h *handler) GetCommentsOfPost(w http.ResponseWriter, r *http.Request)     {}
func (h *handler) UpdateComment(w http.ResponseWriter, r *http.Request)         {}
func (h *handler) DeleteComment(w http.ResponseWriter, r *http.Request)         {}

func (h *handler) CreateFollow(w http.ResponseWriter, r *http.Request)               {}
func (h *handler) GetFollowerCountOfAccount(w http.ResponseWriter, r *http.Request)  {}
func (h *handler) GetFollowersOfAccount(w http.ResponseWriter, r *http.Request)      {}
func (h *handler) GetFollowingCountOfAccount(w http.ResponseWriter, r *http.Request) {}
func (h *handler) GetFollowingsOfAccount(w http.ResponseWriter, r *http.Request)     {}
func (h *handler) DeleteFollow(w http.ResponseWriter, r *http.Request)               {}

func (h *handler) GetNewFeeds(w http.ResponseWriter, r *http.Request) {}
