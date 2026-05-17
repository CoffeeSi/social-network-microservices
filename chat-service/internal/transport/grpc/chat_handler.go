package grpc

import (
	"context"
	"time"

	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/usecase"
	pb "github.com/CoffeeSi/social-network-microservices/chat-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatHandler struct {
	uc usecase.ChatUsecaseInterface
	pb.UnimplementedChatServiceServer
}

func NewChatHandler(server *grpc.Server, uc usecase.ChatUsecaseInterface) {
	pb.RegisterChatServiceServer(server, &ChatHandler{uc: uc})
}

func (h *ChatHandler) GetChats(ctx context.Context, req *pb.GetChatsRequest) (*pb.GetChatsResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chats, err := h.uc.GetChats(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get chats: %v", err)
	}

	chatsPb := make([]*pb.Chat, len(chats))
	for i, chat := range chats {
		chatsPb[i] = mapChatToPb(chat)
	}

	return &pb.GetChatsResponse{Chats: chatsPb}, nil
}

func (h *ChatHandler) CreateDirectChat(ctx context.Context, req *pb.CreateDirectChatRequest) (*pb.CreateDirectChatResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chat := &model.Chat{
		IsGroup:        false,
		ParticipantIds: []string{userID, req.GetTargetUserId()},
		CreatorId:      userID,
		CreatedAt:      time.Now(),
	}

	err = h.uc.CreateDirectChat(ctx, chat)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create direct chat: %v", err)
	}

	return &pb.CreateDirectChatResponse{Chat: mapChatToPb(*chat)}, nil
}

func (h *ChatHandler) CreateGroupChat(ctx context.Context, req *pb.CreateGroupChatRequest) (*pb.CreateGroupChatResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chat := &model.Chat{
		Name:           req.GetName(),
		IsGroup:        true,
		ParticipantIds: req.GetParticipantIds(),
		CreatorId:      userID,
		CreatedAt:      time.Now(),
	}

	err = h.uc.CreateGroupChat(ctx, chat)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create group chat: %v", err)
	}

	return &pb.CreateGroupChatResponse{Chat: mapChatToPb(*chat)}, nil
}

func (h *ChatHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	msg := &model.ChatMessage{
		ChatId:    req.GetChatId(),
		SenderId:  userID,
		Content:   req.GetContent(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = h.uc.SendMessage(ctx, msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	return &pb.SendMessageResponse{Message: mapMessageToPb(*msg)}, nil
}

func (h *ChatHandler) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	messages, err := h.uc.GetMessages(ctx, req.GetChatId(), int(req.GetPageSize()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get messages: %v", err)
	}

	messagesPb := make([]*pb.ChatMessage, len(messages))
	for i, msg := range messages {
		messagesPb[i] = mapMessageToPb(msg)
	}

	return &pb.GetMessagesResponse{
		Messages:   messagesPb,
		TotalCount: int32(len(messages)),
	}, nil
}

func (h *ChatHandler) EditMessage(ctx context.Context, req *pb.EditMessageRequest) (*pb.EditMessageResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	msg := &model.ChatMessage{
		ChatId:    req.GetChatId(),
		Id:        req.GetMessageId(),
		SenderId:  userID,
		Content:   req.GetNewContent(),
		UpdatedAt: time.Now(),
	}

	err = h.uc.EditMessage(ctx, msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to edit message: %v", err)
	}

	return &pb.EditMessageResponse{Message: mapMessageToPb(*msg)}, nil
}

func (h *ChatHandler) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = h.uc.DeleteMessage(ctx, req.GetChatId(), req.GetMessageId(), userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete message: %v", err)
	}

	return &pb.DeleteMessageResponse{Status: "success"}, nil
}

func (h *ChatHandler) SendReadReceipt(ctx context.Context, req *pb.SendReadReceiptRequest) (*pb.SendReadReceiptResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = h.uc.SendReadReceipt(ctx, req.GetChatId(), userID, req.GetMessageId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send read receipt: %v", err)
	}

	return &pb.SendReadReceiptResponse{Status: "success"}, nil
}

func (h *ChatHandler) SendTypingStatus(ctx context.Context, req *pb.SendTypingStatusRequest) (*pb.SendTypingStatusResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = h.uc.SendTypingStatus(ctx, req.GetChatId(), userID, req.GetIsTyping())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send typing status: %v", err)
	}

	return &pb.SendTypingStatusResponse{Status: "success"}, nil
}

func (h *ChatHandler) AddParticipantToGroupChat(ctx context.Context, req *pb.AddParticipantToGroupChatRequest) (*pb.AddParticipantToGroupChatResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = h.uc.AddParticipantToGroupChat(ctx, req.GetChatId(), req.GetUserId(), userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add participant: %v", err)
	}

	return &pb.AddParticipantToGroupChatResponse{Status: "success"}, nil
}

func (h *ChatHandler) RemoveParticipantFromGroupChat(ctx context.Context, req *pb.RemoveParticipantFromGroupChatRequest) (*pb.RemoveParticipantFromGroupChatResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = h.uc.RemoveParticipantFromGroupChat(ctx, req.GetChatId(), req.GetUserId(), userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove participant: %v", err)
	}

	return &pb.RemoveParticipantFromGroupChatResponse{Status: "success"}, nil
}

func (h *ChatHandler) EditGroupChat(ctx context.Context, req *pb.EditGroupChatRequest) (*pb.EditGroupChatResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	chat := &model.Chat{
		Id:   req.GetChatId(),
		Name: req.GetName(),
	}

	err = h.uc.EditGroupChat(ctx, chat, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to edit group chat: %v", err)
	}

	return &pb.EditGroupChatResponse{Status: "success"}, nil
}

func (h *ChatHandler) DeleteGroupChat(ctx context.Context, req *pb.DeleteGroupChatRequest) (*pb.DeleteGroupChatResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = h.uc.DeleteGroupChat(ctx, req.GetChatId(), userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete group chat: %v", err)
	}

	return &pb.DeleteGroupChatResponse{Status: "success"}, nil
}

func (h *ChatHandler) LeaveGroupChat(ctx context.Context, req *pb.LeaveGroupChatRequest) (*pb.LeaveGroupChatResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = h.uc.LeaveGroupChat(ctx, req.GetChatId(), userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to leave group chat: %v", err)
	}

	return &pb.LeaveGroupChatResponse{Status: "success"}, nil
}

func (h *ChatHandler) StreamMessages(req *pb.StreamMessagesRequest, stream grpc.ServerStreamingServer[pb.StreamMessagesResponse]) error {
	return status.Error(codes.Unimplemented, "method StreamMessages not implemented")
}
