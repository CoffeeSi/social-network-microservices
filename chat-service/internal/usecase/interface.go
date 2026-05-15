package usecase

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/model"
)

type ChatUsecaseInterface interface {
	GetChats(ctx context.Context, userId string) ([]model.Chat, error)
	CreateDirectChat(ctx context.Context, chat *model.Chat) error
	CreateGroupChat(ctx context.Context, req *model.Chat) error
	SendMessage(ctx context.Context, message *model.ChatMessage) error
	GetMessages(ctx context.Context, chatID string, limit int) ([]model.ChatMessage, error)
	EditMessage(ctx context.Context, message *model.ChatMessage) error
	DeleteMessage(ctx context.Context, chatID, messageID, userID string) error
	SendReadReceipt(ctx context.Context, chatID, userID, messageID string) error
	SendTypingStatus(ctx context.Context, chatID, userID string, isTyping bool) error
	AddParticipantToGroupChat(ctx context.Context, chatID, targetUserID, adminID string) error
	RemoveParticipantFromGroupChat(ctx context.Context, chatID, targetUserID, adminID string) error
	EditGroupChat(ctx context.Context, chat *model.Chat, adminID string) error
	DeleteGroupChat(ctx context.Context, chatID, adminID string) error
	LeaveGroupChat(ctx context.Context, chatID, userID string) error
}

type UserClientInterface interface {
	GetUserByID(ctx context.Context, id string) (model.User, error)
}
