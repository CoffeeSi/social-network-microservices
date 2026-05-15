package transport

import (
	"context"
	"errors"
	"log"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/usecase/dto"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/usecase/interfaces"
	pb "github.com/CoffeeSi/social-network-microservices/content-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ContentHandler struct {
	pb.UnimplementedContentServiceServer
	postUC    interfaces.PostUseCase
	commentUC interfaces.CommentUseCase
	likeUC    interfaces.LikeUseCase
}

func NewContentHandler(
	postUC interfaces.PostUseCase,
	commentUC interfaces.CommentUseCase,
	likeUC interfaces.LikeUseCase,
) *ContentHandler {
	return &ContentHandler{
		postUC:    postUC,
		commentUC: commentUC,
		likeUC:    likeUC,
	}
}
func (ch *ContentHandler) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	res, err := ch.postUC.CreatePost(ctx, dto.CreatePostDTO{
		AuthorID:  req.AuthorId,
		Content:   req.Content,
		MediaURLs: req.MediaUrls,
	})
	if err != nil {
		if errors.Is(err, model.ErrEmptyPost) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		log.Printf("[gRPC Error Trace] CreatePost failed with internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return postToProto(&res), nil
}

func (ch *ContentHandler) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	res, err := ch.postUC.GetPostByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, model.ErrPostNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return postToProto(&res), nil
}
func (ch *ContentHandler) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	var posts []model.Post
	var total int32
	var err error

	feedDTO := dto.GetFeedDTO{
		PageSize: req.PageSize,
		Page:     req.Page,
	}

	if req.AuthorId != "" {
		posts, total, err = ch.postUC.GetUserPosts(ctx, req.AuthorId, feedDTO)
	} else {
		posts, total, err = ch.postUC.GetFeed(ctx, feedDTO)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	pbPosts := make([]*pb.Post, 0, len(posts))
	for _, p := range posts {
		pbPosts = append(pbPosts, postToProto(&p))
	}

	return &pb.ListPostsResponse{
		Posts:      pbPosts,
		TotalCount: total,
	}, nil
}
func (ch *ContentHandler) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*emptypb.Empty, error) {
	err := ch.postUC.DeletePost(ctx, req.Id, req.UserId)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, model.ErrPostNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, model.ErrPermissionDenied) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &emptypb.Empty{}, nil
}

func (ch *ContentHandler) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	res, err := ch.postUC.UpdatePost(ctx, req.Id, req.UserId, req.Content, req.MediaUrls)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, model.ErrEmptyPost) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, model.ErrPostNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, model.ErrPermissionDenied) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return postToProto(&res), nil
}

func (ch *ContentHandler) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.Comment, error) {
	res, err := ch.commentUC.AddComment(ctx, dto.CreateCommentDTO{
		PostID: req.PostId,
		UserID: req.UserId,
		Text:   req.Text,
	})
	if err != nil {
		if errors.Is(err, model.ErrEmptyComment) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, model.ErrPostNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return commentToProto(&res), nil
}

func (ch *ContentHandler) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	comments, total, err := ch.commentUC.ListComments(ctx, dto.ListCommentsDTO{
		PostID:   req.PostId,
		PageSize: req.PageSize,
		Page:     req.Page,
	})
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, model.ErrPostNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	pbComments := make([]*pb.Comment, 0, len(comments))
	for _, c := range comments {
		pbComments = append(pbComments, commentToProto(&c))
	}

	return &pb.ListCommentsResponse{
		Comments:   pbComments,
		TotalCount: total,
	}, nil
}

func (ch *ContentHandler) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*emptypb.Empty, error) {
	err := ch.commentUC.DeleteComment(ctx, req.Id, req.UserId)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, model.ErrCommentNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, model.ErrPermissionDenied) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &emptypb.Empty{}, nil
}

func (ch *ContentHandler) ToggleLike(ctx context.Context, req *pb.ToggleLikeRequest) (*pb.LikeResponse, error) {
	err := ch.likeUC.LikePost(ctx, req.PostId, req.UserId)

	if err == nil {
		post, fetchErr := ch.postUC.GetPostByID(ctx, req.PostId)
		if fetchErr != nil {
			return nil, status.Error(codes.Internal, "like completed but status sync failed")
		}
		return &pb.LikeResponse{NewLikeCount: post.LikeCount, IsLiked: true}, nil
	}

	if errors.Is(err, model.ErrAlreadyLiked) {
		err = ch.likeUC.UnlikePost(ctx, req.PostId, req.UserId)
		if err != nil {
			if errors.Is(err, model.ErrLikeNotFound) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			return nil, status.Error(codes.Internal, "internal server error")
		}

		post, fetchErr := ch.postUC.GetPostByID(ctx, req.PostId)
		if fetchErr != nil {
			return nil, status.Error(codes.Internal, "unlike completed but status sync failed")
		}
		return &pb.LikeResponse{NewLikeCount: post.LikeCount, IsLiked: false}, nil
	}

	if errors.Is(err, model.ErrPostNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, model.ErrInvalidID) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return nil, status.Error(codes.Internal, "internal server error")
}
