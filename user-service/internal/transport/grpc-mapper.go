package transport

import (
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/model"
	pb "github.com/CoffeeSi/social-network-microservices/user-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func userToProto(user *model.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Dob:       timestamppb.New(user.DOB),
		CreatedAt: timestamppb.New(user.CreatedAt),
		IsActive:  user.IsActive,
	}
}

func userWithPasswordToProto(user *model.User) *pb.UserWithPassword {
	return &pb.UserWithPassword{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Dob:       timestamppb.New(user.DOB),
		CreatedAt: timestamppb.New(user.CreatedAt),
		IsActive:  user.IsActive,
		Password:  user.Password,
	}
}
