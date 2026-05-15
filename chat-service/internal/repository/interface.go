package repository

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/model"
)

type ChatRepositoryInterface interface {
	GetChats(ctx context.Context, userId string) ([]model.Chat, error)
	GetChatByID(ctx context.Context, chatID string) (*model.Chat, error)
	GetMessageByID(ctx context.Context, chatID, messageID string) (*model.ChatMessage, error)
	CreateDirectChat(ctx context.Context, chat *model.Chat) error
	CreateGroupChat(ctx context.Context, chat *model.Chat) error
	SendMessage(ctx context.Context, message *model.ChatMessage) error
	GetMessages(ctx context.Context, chatID string, limit int) ([]model.ChatMessage, error)
	EditMessage(ctx context.Context, message *model.ChatMessage) error
	DeleteMessage(ctx context.Context, messageID string) error
	SendReadReceipt(ctx context.Context, chatID, userID, messageID string) error
	SendTypingStatus(ctx context.Context, chatID, userID string, isTyping bool) error
	AddParticipantToGroupChat(ctx context.Context, chatID, userID string) error
	RemoveParticipantFromGroupChat(ctx context.Context, chatID, userID string) error
	EditGroupChat(ctx context.Context, chat *model.Chat) error
	DeleteGroupChat(ctx context.Context, chatID string) error
	LeaveGroupChat(ctx context.Context, chatID, userID string) error
}
