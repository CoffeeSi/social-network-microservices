package transport

import (
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	pb "github.com/IsFariza/maxat-protobuf/content-service-pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func postToProto(p *model.Post) *pb.Post {
	if p == nil {
		return nil
	}
	return &pb.Post{
		Id:           p.ID,
		AuthorId:     p.AuthorID,
		Content:      p.Content,
		MediaUrls:    p.MediaURLs,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
		CreatedAt:    timestamppb.New(p.CreatedAt),
	}
}
func commentToProto(c *model.Comment) *pb.Comment {
	if c == nil {
		return nil
	}
	return &pb.Comment{
		Id:        c.ID,
		PostId:    c.PostID,
		UserId:    c.UserID,
		Text:      c.Text,
		CreatedAt: timestamppb.New(c.CreatedAt),
	}
}
