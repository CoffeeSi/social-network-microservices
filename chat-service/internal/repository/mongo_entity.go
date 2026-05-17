package repository

import (
	"time"

	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type chatMongoEntity struct {
	Id             primitive.ObjectID      `bson:"_id,omitempty"`
	Name           string                  `bson:"name"`
	IsGroup        bool                    `bson:"is_group"`
	ParticipantIds []string                `bson:"participant_ids"`
	CreatorId      string                  `bson:"creator_id"`
	CreatedAt      time.Time               `bson:"created_at"`
	LastMessage    *chatMessageMongoEntity `bson:"last_message,omitempty"`
}

type chatMessageMongoEntity struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	ChatId    string             `bson:"chat_id"`
	SenderId  string             `bson:"sender_id"`
	Content   string             `bson:"content"`
	IsEdited  bool               `bson:"is_edited"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func toDomainChat(entity chatMongoEntity) model.Chat {
	chat := model.Chat{
		Id:             entity.Id.Hex(),
		Name:           entity.Name,
		IsGroup:        entity.IsGroup,
		ParticipantIds: entity.ParticipantIds,
		CreatorId:      entity.CreatorId,
		CreatedAt:      entity.CreatedAt,
	}

	if entity.LastMessage != nil {
		msg := toDomainMessage(*entity.LastMessage)
		chat.LastMessage = &msg
	}

	return chat
}

func toMongoChat(chat *model.Chat) chatMongoEntity {
	entity := chatMongoEntity{
		Name:           chat.Name,
		IsGroup:        chat.IsGroup,
		ParticipantIds: chat.ParticipantIds,
		CreatorId:      chat.CreatorId,
		CreatedAt:      chat.CreatedAt,
	}

	if chat.Id != "" {
		id, err := primitive.ObjectIDFromHex(chat.Id)
		if err == nil {
			entity.Id = id
		}
	}

	if chat.LastMessage != nil {
		msg := toMongoMessage(*chat.LastMessage)
		entity.LastMessage = &msg
	}

	return entity
}

func toDomainMessage(entity chatMessageMongoEntity) model.ChatMessage {
	return model.ChatMessage{
		Id:        entity.Id.Hex(),
		ChatId:    entity.ChatId,
		SenderId:  entity.SenderId,
		Content:   entity.Content,
		IsEdited:  entity.IsEdited,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func toMongoMessage(msg model.ChatMessage) chatMessageMongoEntity {
	entity := chatMessageMongoEntity{
		ChatId:    msg.ChatId,
		SenderId:  msg.SenderId,
		Content:   msg.Content,
		IsEdited:  msg.IsEdited,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	}

	if msg.Id != "" {
		id, err := primitive.ObjectIDFromHex(msg.Id)
		if err == nil {
			entity.Id = id
		}
	}

	return entity
}
