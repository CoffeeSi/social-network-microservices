package grpc

import (
	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/model"
	pb "github.com/CoffeeSi/social-network-microservices/chat-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapChatToPb(chat model.Chat) *pb.Chat {
	pbChat := &pb.Chat{
		Id:             chat.Id,
		Name:           chat.Name,
		IsGroup:        chat.IsGroup,
		ParticipantIds: chat.ParticipantIds,
		CreatorId:      chat.CreatorId,
		CreatedAt:      timestamppb.New(chat.CreatedAt),
	}

	if chat.LastMessage != nil {
		pbChat.LastMessage = mapMessageToPb(*chat.LastMessage)
	}

	return pbChat
}

func mapMessageToPb(msg model.ChatMessage) *pb.ChatMessage {
	return &pb.ChatMessage{
		Id:        msg.Id,
		ChatId:    msg.ChatId,
		SenderId:  msg.SenderId,
		Content:   msg.Content,
		IsEdited:  msg.IsEdited,
		CreatedAt: timestamppb.New(msg.CreatedAt),
		UpdatedAt: timestamppb.New(msg.UpdatedAt),
	}
}
