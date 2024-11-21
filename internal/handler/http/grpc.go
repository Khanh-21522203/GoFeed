package http

// import (
// 	"GoFeed/internal/generated/api/go_feed"
// 	"context"
// )

// type grpcClient struct {
// 	gofeedClient *go_feed.GoFeedServiceClient
// }

// func NewGRPCClient(
// 	gofeedClient *go_feed.GoFeedServiceClient,
// ) *grpcClient {
// 	return &grpcClient{
// 		gofeedClient: gofeedClient,
// 	}
// }

// func (g grpcClient) CreateAccount(ctx context.Context, p *go_feed.CreateAccountRequest) (*go_feed.CreateAccountResponse, error) {
// }
// func (g grpcClient) CreateSession(ctx context.Context) {}

// func (g grpcClient) CreatePost(ctx context.Context)       {}
// func (g grpcClient) GetPostByID(ctx context.Context)      {}
// func (g grpcClient) GetPostOfAccount(ctx context.Context) {}
// func (g grpcClient) UpdatePost(ctx context.Context)       {}

// func (g grpcClient) CreateLike(ctx context.Context)            {}
// func (g grpcClient) GetLikeCountOfPost(ctx context.Context)    {}
// func (g grpcClient) GetLikeAccountsOfPost(ctx context.Context) {}
// func (g grpcClient) DeleteLike(ctx context.Context)            {}

// func (g grpcClient) CreateComment(ctx context.Context)         {}
// func (g grpcClient) GetCommentCountOfPost(ctx context.Context) {}
// func (g grpcClient) GetCommentsOfPost(ctx context.Context)     {}
// func (g grpcClient) UpdateComment(ctx context.Context)         {}
// func (g grpcClient) DeleteComment(ctx context.Context)         {}

// func (g grpcClient) CreateFollow(ctx context.Context)               {}
// func (g grpcClient) GetFollowerCountOfAccount(ctx context.Context)  {}
// func (g grpcClient) GetFollowersOfAccount(ctx context.Context)      {}
// func (g grpcClient) GetFollowingCountOfAccount(ctx context.Context) {}
// func (g grpcClient) GetFollowingsOfAccount(ctx context.Context)     {}
// func (g grpcClient) DeleteFollow(ctx context.Context)               {}

// func (g grpcClient) GetNewFeeds(ctx context.Context) {}
