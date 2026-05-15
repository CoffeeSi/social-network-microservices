package usecase

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/repository"
)

type ChatUsecase struct {
	repo       repository.ChatRepositoryInterface
	userClient UserClientInterface
}

func NewChatUsecase(repo repository.ChatRepositoryInterface, userClient UserClientInterface) *ChatUsecase {
	return &ChatUsecase{
		repo:       repo,
		userClient: userClient,
	}
}

func (uc *ChatUsecase) GetChats(ctx context.Context, userId string) ([]model.Chat, error) {
	return uc.repo.GetChats(ctx, userId)
}

func (uc *ChatUsecase) CreateDirectChat(ctx context.Context, chat *model.Chat) error {
	if len(chat.ParticipantIds) < 2 {
		return errors.New("direct chat requires 2 participants")
	}

	targetUser, err := uc.userClient.GetUserByID(ctx, chat.ParticipantIds[1])
	if err != nil {
		return err
	}

	if targetUser.ID == "" {
		return errors.New("target user not found")
	}
	if chat.ParticipantIds[0] == chat.ParticipantIds[1] {
		return errors.New("user cannot create chat with themselves")
	}

	if chat.ParticipantIds[0] == "" || chat.ParticipantIds[1] == "" {
		return errors.New("user cannot create chat with empty id")
	}
	return uc.repo.CreateDirectChat(ctx, chat)
}

func (uc *ChatUsecase) CreateGroupChat(ctx context.Context, chat *model.Chat) error {
	if len(chat.ParticipantIds) == 0 {
		return errors.New("chat must have at least one participant")
	}

	creatorFound := false
	for _, participantId := range chat.ParticipantIds {
		if participantId == "" {
			return errors.New("user cannot create chat with empty id")
		}
		if participantId == chat.CreatorId {
			creatorFound = true
		}

		_, err := uc.userClient.GetUserByID(ctx, participantId)
		if err != nil {
			return err
		}
	}
	if !creatorFound {
		chat.ParticipantIds = append(chat.ParticipantIds, chat.CreatorId)
	}

	if chat.Name == "" {
		return errors.New("chat name cannot be empty")
	}
	return uc.repo.CreateGroupChat(ctx, chat)
}

func (uc *ChatUsecase) SendMessage(ctx context.Context, message *model.ChatMessage) error {
	if message.SenderId == "" {
		return errors.New("message sender cannot be empty")
	}
	if message.Content == "" {
		return errors.New("message content cannot be empty")
	}
	if message.ChatId == "" {
		return errors.New("message chat id cannot be empty")
	}

	chat, err := uc.repo.GetChatByID(ctx, message.ChatId)
	if err != nil {
		return err
	}
	if chat == nil {
		return errors.New("chat not found")
	}

	isParticipant := slices.Contains(chat.ParticipantIds, message.SenderId)
	if !isParticipant {
		return errors.New("user is not a participant of this chat")
	}

	return uc.repo.SendMessage(ctx, message)
}

func (uc *ChatUsecase) GetMessages(ctx context.Context, chatID string, limit int) ([]model.ChatMessage, error) {
	if chatID == "" {
		return nil, errors.New("chat id cannot be empty")
	}
	if limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}
	return uc.repo.GetMessages(ctx, chatID, limit)
}

func (uc *ChatUsecase) EditMessage(ctx context.Context, message *model.ChatMessage) error {
	if message.Content == "" {
		return errors.New("message content cannot be empty")
	}
	if message.ChatId == "" {
		return errors.New("message chat id cannot be empty")
	}
	findMessage, err := uc.repo.GetMessageByID(ctx, message.ChatId, message.Id)
	if err != nil || findMessage == nil {
		return errors.New("message not found")
	}
	if findMessage.SenderId != message.SenderId {
		return errors.New("user is not the sender of this message")
	}
	message.UpdatedAt = time.Now()
	return uc.repo.EditMessage(ctx, message)
}

func (uc *ChatUsecase) DeleteMessage(ctx context.Context, chatID, messageID, userID string) error {
	if chatID == "" {
		return errors.New("chat id cannot be empty")
	}
	if messageID == "" {
		return errors.New("message id cannot be empty")
	}
	findMessage, err := uc.repo.GetMessageByID(ctx, chatID, messageID)
	if err != nil || findMessage == nil {
		return errors.New("message not found")
	}
	if findMessage.SenderId != userID {
		return errors.New("user is not the sender of this message")
	}
	return uc.repo.DeleteMessage(ctx, messageID)
}

func (uc *ChatUsecase) SendReadReceipt(ctx context.Context, chatID, userID, messageID string) error {
	if chatID == "" {
		return errors.New("chat id cannot be empty")
	}
	if userID == "" {
		return errors.New("user id cannot be empty")
	}
	if messageID == "" {
		return errors.New("message id cannot be empty")
	}
	return uc.repo.SendReadReceipt(ctx, chatID, userID, messageID)
}

func (uc *ChatUsecase) SendTypingStatus(ctx context.Context, chatID, userID string, isTyping bool) error {
	if chatID == "" {
		return errors.New("chat id cannot be empty")
	}
	if userID == "" {
		return errors.New("user id cannot be empty")
	}
	return uc.repo.SendTypingStatus(ctx, chatID, userID, isTyping)
}

func (uc *ChatUsecase) AddParticipantToGroupChat(ctx context.Context, chatID, targetUserID, adminID string) error {
	if chatID == "" {
		return errors.New("chat id cannot be empty")
	}
	if targetUserID == "" {
		return errors.New("target user id cannot be empty")
	}

	_, err := uc.userClient.GetUserByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	chat, err := uc.repo.GetChatByID(ctx, chatID)
	if err != nil {
		return err
	}
	if chat == nil {
		return errors.New("chat not found")
	}
	if chat.CreatorId != adminID {
		return errors.New("only the creator can add participants")
	}

	return uc.repo.AddParticipantToGroupChat(ctx, chatID, targetUserID)
}

func (uc *ChatUsecase) RemoveParticipantFromGroupChat(ctx context.Context, chatID, targetUserID, adminID string) error {
	if chatID == "" {
		return errors.New("chat id cannot be empty")
	}
	if targetUserID == "" {
		return errors.New("target user id cannot be empty")
	}

	chat, err := uc.repo.GetChatByID(ctx, chatID)
	if err != nil {
		return err
	}
	if chat == nil {
		return errors.New("chat not found")
	}
	if chat.CreatorId != adminID {
		return errors.New("only the creator can remove participants")
	}

	return uc.repo.RemoveParticipantFromGroupChat(ctx, chatID, targetUserID)
}

func (uc *ChatUsecase) EditGroupChat(ctx context.Context, chat *model.Chat, adminID string) error {
	if chat.Id == "" {
		return errors.New("chat id cannot be empty")
	}
	if chat.Name == "" {
		return errors.New("chat name cannot be empty")
	}
	if len(chat.ParticipantIds) > 0 && chat.ParticipantIds[0] == "" {
		return errors.New("user cannot edit chat with empty id")
	}

	existingChat, err := uc.repo.GetChatByID(ctx, chat.Id)
	if err != nil {
		return err
	}
	if existingChat == nil {
		return errors.New("chat not found")
	}
	if existingChat.CreatorId != adminID {
		return errors.New("only the creator can edit the chat")
	}

	return uc.repo.EditGroupChat(ctx, chat)
}

func (uc *ChatUsecase) DeleteGroupChat(ctx context.Context, chatID, adminID string) error {
	if chatID == "" {
		return errors.New("chat id cannot be empty")
	}

	chat, err := uc.repo.GetChatByID(ctx, chatID)
	if err != nil {
		return err
	}
	if chat == nil {
		return errors.New("chat not found")
	}
	if chat.CreatorId != adminID {
		return errors.New("only the creator can delete the chat")
	}

	return uc.repo.DeleteGroupChat(ctx, chatID)
}

func (uc *ChatUsecase) LeaveGroupChat(ctx context.Context, chatID, userID string) error {
	if chatID == "" {
		return errors.New("chat id cannot be empty")
	}
	if userID == "" {
		return errors.New("user id cannot be empty")
	}
	return uc.repo.LeaveGroupChat(ctx, chatID, userID)
}
