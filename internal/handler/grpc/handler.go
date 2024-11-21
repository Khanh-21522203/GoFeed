package grpc

import (
	"GoFeed/internal/generated/api/go_feed"
	"GoFeed/internal/logic"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	//nolint:gosec // This is just to specify the metadata name
	AuthTokenMetadataName = "GOFEED_AUTH"
)

type grpcHandler struct {
	go_feed.UnimplementedGoFeedServiceServer

	accountLogic logic.AccountLogic
	postLogic    logic.PostLogic
	commentLogic logic.CommentLogic
	followLogic  logic.FollowLogic
	likeLogic    logic.LikeLogic
}

func NewHandler(
	accountLogic logic.AccountLogic,
	postLogic logic.PostLogic,
	commentLogic logic.CommentLogic,
	followLogic logic.FollowLogic,
	likeLogic logic.LikeLogic,
) go_feed.GoFeedServiceServer {
	return &grpcHandler{
		accountLogic: accountLogic,
		postLogic:    postLogic,
		commentLogic: commentLogic,
		followLogic:  followLogic,
		likeLogic:    likeLogic,
	}
}

func (g grpcHandler) getAuthTokenMetadata(ctx context.Context) string {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	metadataValues := metadata.Get(AuthTokenMetadataName)
	if len(metadataValues) == 0 {
		return ""
	}

	return metadataValues[0]
}

func (g grpcHandler) CreateAccount(ctx context.Context, request *go_feed.CreateAccountRequest) (*go_feed.CreateAccountResponse, error) {
	output, err := g.accountLogic.CreateAccount(ctx, logic.CreateAccountParams{
		AccountName: request.GetAccountName(),
		Password:    request.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.CreateAccountResponse{
		AccountId: output.ID,
	}, nil
}
func (g grpcHandler) CreateSession(ctx context.Context, request *go_feed.CreateSessionRequest) (*go_feed.CreateSessionResponse, error) {
	output, err := g.accountLogic.CreateSession(ctx, logic.CreateSessionParams{
		AccountName: request.GetAccountName(),
		Password:    request.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	err = grpc.SetHeader(ctx, metadata.Pairs(AuthTokenMetadataName, output.Token))
	if err != nil {
		return nil, err
	}

	return &go_feed.CreateSessionResponse{
		AccountId: output.Account.Id,
	}, nil
}

func (g grpcHandler) CreatePost(ctx context.Context, request *go_feed.CreatePostRequest) (*go_feed.CreatePostResponse, error) {
	output, err := g.postLogic.CreatePost(ctx, logic.CreatePostParams{
		Token:   g.getAuthTokenMetadata(ctx),
		Content: request.GetContent(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.CreatePostResponse{
		PostId: output.ID,
	}, nil
}
func (g grpcHandler) GetPostByID(ctx context.Context, request *go_feed.GetPostByIDRequest) (*go_feed.GetPostByIDResponse, error) {
	output, err := g.postLogic.GetPostByID(ctx, logic.GetPostByIDParams{
		Token: g.getAuthTokenMetadata(ctx),
		ID:    request.GetPostId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.GetPostByIDResponse{
		Post: output.Post,
	}, nil
}
func (g grpcHandler) GetPostOfAccount(ctx context.Context, request *go_feed.GetPostOfAccountRequest) (*go_feed.GetPostOfAccountResponse, error) {
	output, err := g.postLogic.GetPostOfAccount(ctx, logic.GetPostOfAccountParams{
		Token:      g.getAuthTokenMetadata(ctx),
		Of_account: request.GetAccountId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.GetPostOfAccountResponse{
		PostList: output.PostList,
	}, nil
}
func (g grpcHandler) UpdatePost(ctx context.Context, request *go_feed.UpdatePostRequest) (*go_feed.UpdatePostResponse, error) {
	_, err := g.postLogic.UpdatePost(ctx, logic.UpdatePostParams{
		Token:   g.getAuthTokenMetadata(ctx),
		ID:      request.GetPost().Id,
		Content: request.GetPost().Content,
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.UpdatePostResponse{}, nil
}

func (g grpcHandler) CreateLike(ctx context.Context, request *go_feed.CreateLikeRequest) (*go_feed.CreateLikeResponse, error) {
	err := g.likeLogic.CreateLike(ctx, logic.CreateLikeParams{
		Token:  g.getAuthTokenMetadata(ctx),
		PostID: request.GetPostId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.CreateLikeResponse{}, nil
}
func (g grpcHandler) GetLikeCountOfPost(ctx context.Context, request *go_feed.GetLikeCountOfPostRequest) (*go_feed.GetLikeCountOfPostResponse, error) {
	output, err := g.likeLogic.GetLikeCountOfPost(ctx, logic.GetLikeCountOfPostParams{
		Token:  g.getAuthTokenMetadata(ctx),
		PostID: request.GetPostId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.GetLikeCountOfPostResponse{
		LikeCount: uint64(output.LikeCount),
	}, nil
}
func (g grpcHandler) GetLikeAccountsOfPost(ctx context.Context, request *go_feed.GetLikeAccountsOfPostRequest) (*go_feed.GetLikeAccountsOfPostResponse, error) {
	output, err := g.likeLogic.GetLikeAccountsOfPost(ctx, logic.GetLikeAccountsOfPostParams{
		Token:  g.getAuthTokenMetadata(ctx),
		PostID: request.GetPostId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.GetLikeAccountsOfPostResponse{
		AccountList: output.AccountList,
	}, nil
}
func (g grpcHandler) DeleteLike(ctx context.Context, request *go_feed.DeleteLikeRequest) (*go_feed.DeleteLikeResponse, error) {
	err := g.likeLogic.DeleteLike(ctx, logic.DeleteLikeParams{
		Token:  g.getAuthTokenMetadata(ctx),
		PostID: request.GetPostId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.DeleteLikeResponse{}, nil
}

func (g grpcHandler) CreateComment(ctx context.Context, request *go_feed.CreateCommentRequest) (*go_feed.CreateCommentResponse, error) {
	output, err := g.commentLogic.CreateComment(ctx, logic.CreateCommentParams{
		Token:   g.getAuthTokenMetadata(ctx),
		PostID:  request.GetPostId(),
		Content: request.GetContent(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.CreateCommentResponse{
		CommentId: output.ID,
	}, nil
}
func (g grpcHandler) GetCommentCountOfPost(ctx context.Context, request *go_feed.GetCommentCountOfPostRequest) (*go_feed.GetCommentCountOfPostResponse, error) {
	output, err := g.commentLogic.GetCommentCountOfPost(ctx, logic.GetCommentCountOfPostParams{
		Token:  g.getAuthTokenMetadata(ctx),
		PostID: request.GetPostId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.GetCommentCountOfPostResponse{
		CommentCount: uint64(output.CommentCount),
	}, nil
}
func (g grpcHandler) GetCommentsOfPost(ctx context.Context, request *go_feed.GetCommentsOfPostRequest) (*go_feed.GetCommentsOfPostResponse, error) {
	output, err := g.commentLogic.GetCommentsOfPost(ctx, logic.GetCommentsOfPostParams{
		Token:  g.getAuthTokenMetadata(ctx),
		PostID: request.GetPostId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.GetCommentsOfPostResponse{
		CommentList: output.CommentList,
	}, nil
}
func (g grpcHandler) UpdateComment(ctx context.Context, request *go_feed.UpdateCommentRequest) (*go_feed.UpdateCommentResponse, error) {
	err := g.commentLogic.UpdateComment(ctx, logic.UpdateCommentParams{
		Token:   g.getAuthTokenMetadata(ctx),
		ID:      request.GetComment().CommentId,
		Content: request.GetComment().Content,
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.UpdateCommentResponse{
		CommentId: request.GetComment().CommentId,
	}, nil
}
func (g grpcHandler) DeleteComment(ctx context.Context, request *go_feed.DeleteCommentRequest) (*go_feed.DeleteCommentResponse, error) {
	err := g.commentLogic.DeleteComment(ctx, logic.DeleteCommentParams{
		Token: g.getAuthTokenMetadata(ctx),
		ID:    request.GetCommentId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.DeleteCommentResponse{}, nil
}

func (g grpcHandler) CreateFollow(ctx context.Context, request *go_feed.CreateFollowRequest) (*go_feed.CreateFollowResponse, error) {
	err := g.followLogic.CreateFollow(ctx, logic.CreateFollowParams{
		Token:       g.getAuthTokenMetadata(ctx),
		FollowingID: request.GetFollowingId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.CreateFollowResponse{}, nil
}
func (g grpcHandler) GetFollowerCountOfAccount(ctx context.Context, request *go_feed.GetFollowerCountOfAccountRequest) (*go_feed.GetFollowerCountOfAccountResponse, error) {
	output, err := g.followLogic.GetFollowerCountOfAccount(ctx, logic.GetFollowerCountOfAccountParams{
		Token:     g.getAuthTokenMetadata(ctx),
		AccountID: request.GetAccountId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.GetFollowerCountOfAccountResponse{
		FollowerCount: uint64(output.FollowerCount),
	}, nil
}
func (g grpcHandler) GetFollowersOfAccount(ctx context.Context, request *go_feed.GetFollowersOfAccountRequest) (*go_feed.GetFollowersOfAccountResponse, error) {
	output, err := g.followLogic.GetFollowersOfAccount(ctx, logic.GetFollowersOfAccountParams{
		Token:     g.getAuthTokenMetadata(ctx),
		AccountID: request.GetAccountId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.GetFollowersOfAccountResponse{
		FollowerList: output.FollowerList,
	}, nil
}
func (g grpcHandler) GetFollowingCountOfAccount(ctx context.Context, request *go_feed.GetFollowingCountOfAccountRequest) (*go_feed.GetFollowingCountOfAccountResponse, error) {
	output, err := g.followLogic.GetFollowingCountOfAccount(ctx, logic.GetFollowingCountOfAccountParams{
		Token:     g.getAuthTokenMetadata(ctx),
		AccountID: request.GetAccountId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.GetFollowingCountOfAccountResponse{
		FollowingCount: uint64(output.FollowingCount),
	}, nil
}
func (g grpcHandler) GetFollowingsOfAccount(ctx context.Context, request *go_feed.GetFollowingsOfAccountRequest) (*go_feed.GetFollowingsOfAccountResponse, error) {
	output, err := g.followLogic.GetFollowingsOfAccount(ctx, logic.GetFollowingsOfAccountParams{
		Token:     g.getAuthTokenMetadata(ctx),
		AccountID: request.GetAccountId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.GetFollowingsOfAccountResponse{
		FollowingList: output.FollowingList,
	}, nil
}
func (g grpcHandler) DeleteFollow(ctx context.Context, request *go_feed.DeleteFollowRequest) (*go_feed.DeleteFollowResponse, error) {
	err := g.followLogic.DeleteFollow(ctx, logic.DeleteFollowParams{
		Token:       g.getAuthTokenMetadata(ctx),
		FollowingID: request.GetFollowingId(),
	})
	if err != nil {
		return nil, err
	}

	return &go_feed.DeleteFollowResponse{}, nil
}

// func (g grpcHandler) GetNewFeeds(ctx context.Context, request *go_feed.GetNewFeedsRequest) (*go_feed.GetNewFeedsResponse, error)
