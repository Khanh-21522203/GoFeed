package grpc

import (
	"GoFeed/internal/generated/api/go_feed"
)

type handler struct {
	go_feed.UnimplementedGoFeedServiceServer
}

func NewHandler() go_feed.GoFeedServiceServer {
	return &handler{}
}

// func (h handler) CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error) {
// }
// func (h handler) CreateSession(context.Context, *CreateSessionRequest) (*CreateSessionResponse, error)
// func (h handler) CreatePost(context.Context, *CreatePostRequest) (*CreatePostResponse, error)
// func (h handler) GetPostByID(context.Context, *GetPostByIDRequest) (*GetPostByIDResponse, error)
// func (h handler) GetPostOfAccount(context.Context, *GetPostOfAccountRequest) (*GetPostOfAccountResponse, error)
// func (h handler) UpdatePost(context.Context, *UpdatePostRequest) (*UpdatePostResponse, error)
// func (h handler) CreateLike(context.Context, *CreateLikeRequest) (*CreateLikeResponse, error)
// func (h handler) GetLikeCountOfPost(context.Context, *GetLikeCountOfPostRequest) (*GetLikeCountOfPostResponse, error)
// func (h handler) GetLikeAccountsOfPost(context.Context, *GetLikeAccountsOfPostRequest) (*GetLikeAccountsOfPostResponse, error)
// func (h handler) DeleteLike(context.Context, *DeleteLikeRequest) (*DeleteLikeResponse, error)
// func (h handler) CreateComment(context.Context, *CreateCommentRequest) (*CreateCommentResponse, error)
// func (h handler) GetCommentCountOfPost(context.Context, *GetCommentCountOfPostRequest) (*GetCommentCountOfPostResponse, error)
// func (h handler) GetCommentsOfPost(context.Context, *GetCommentsOfPostRequest) (*GetCommentsOfPostResponse, error)
// func (h handler) UpdateComment(context.Context, *UpdateCommentRequest) (*UpdateCommentResponse, error)
// func (h handler) DeleteComment(context.Context, *DeleteCommentRequest) (*DeleteCommentResponse, error)
// func (h handler) CreateFollow(context.Context, *CreateFollowRequest) (*CreateFollowResponse, error)
// func (h handler) GetFollowerCountOfAccount(context.Context, *GetFollowerCountOfAccountRequest) (*GetFollowerCountOfAccountResponse, error)
// func (h handler) GetFollowersOfAccount(context.Context, *GetFollowersOfAccountRequest) (*GetFollowersOfAccountResponse, error)
// func (h handler) GetFollowingCountOfAccount(context.Context, *GetFollowingCountOfAccountRequest) (*GetFollowingCountOfAccountResponse, error)
// func (h handler) GetFollowingsOfAccount(context.Context, *GetFollowingsOfAccountRequest) (*GetFollowingsOfAccountResponse, error)
// func (h handler) DeleteFollow(context.Context, *DeleteFollowRequest) (*DeleteFollowResponse, error)
// func (h handler) GetNewFeeds(context.Context, *GetNewFeedsRequest) (*GetNewFeedsResponse, error)
