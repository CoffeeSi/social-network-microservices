package repository

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/chat-service/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository struct {
	db          *mongo.MongoDB
	chatsCol    *mongodb.Collection
	messagesCol *mongodb.Collection
}

func NewChatRepository(db *mongo.MongoDB) *ChatRepository {
	return &ChatRepository{
		db:          db,
		chatsCol:    db.Database().Collection("chats"),
		messagesCol: db.Database().Collection("messages"),
	}
}

func (r *ChatRepository) GetChats(ctx context.Context, userId string) ([]model.Chat, error) {
	cursor, err := r.chatsCol.Find(ctx, bson.M{"participant_ids": userId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entities []chatMongoEntity
	if err := cursor.All(ctx, &entities); err != nil {
		return nil, err
	}

	chats := make([]model.Chat, len(entities))
	for i, entity := range entities {
		chats[i] = toDomainChat(entity)
	}

	return chats, nil
}

func (r *ChatRepository) GetChatByID(ctx context.Context, chatID string) (*model.Chat, error) {
	oid, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return nil, err
	}
	var entity chatMongoEntity
	err = r.chatsCol.FindOne(ctx, bson.M{"_id": oid}).Decode(&entity)
	if err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	chat := toDomainChat(entity)
	return &chat, nil
}

func (r *ChatRepository) GetMessageByID(ctx context.Context, chatID, messageID string) (*model.ChatMessage, error) {
	messageOid, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return nil, err
	}
	var entity chatMessageMongoEntity
	// chat_id is stored as string in messages (see mongo_entity.chatMessageMongoEntity).
	err = r.messagesCol.FindOne(ctx, bson.M{"chat_id": chatID, "_id": messageOid}).Decode(&entity)
	if err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	message := toDomainMessage(entity)
	return &message, nil
}

func (r *ChatRepository) CreateDirectChat(ctx context.Context, chat *model.Chat) error {
	filter := bson.M{
		"is_group": false,
		"participant_ids": bson.M{
			"$all": chat.ParticipantIds,
		},
	}

	var existingEntity chatMongoEntity
	err := r.chatsCol.FindOne(ctx, filter).Decode(&existingEntity)
	if err == nil {
		chat.Id = existingEntity.Id.Hex()
		chat.CreatedAt = existingEntity.CreatedAt
		return nil
	}

	entity := toMongoChat(chat)
	res, err := r.chatsCol.InsertOne(ctx, entity)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		chat.Id = oid.Hex()
	}
	return nil
}

func (r *ChatRepository) CreateGroupChat(ctx context.Context, chat *model.Chat) error {
	entity := toMongoChat(chat)
	res, err := r.chatsCol.InsertOne(ctx, entity)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		chat.Id = oid.Hex()
	}
	return nil
}

func (r *ChatRepository) SendMessage(ctx context.Context, message *model.ChatMessage) error {
	entity := toMongoMessage(*message)
	res, err := r.messagesCol.InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		message.Id = oid.Hex()
	}

	if message.ChatId != "" {
		chatId, err := primitive.ObjectIDFromHex(message.ChatId)
		if err == nil {
			update := bson.M{
				"$set": bson.M{
					"last_message": entity,
				},
			}
			_, err = r.chatsCol.UpdateByID(ctx, chatId, update)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *ChatRepository) GetMessages(ctx context.Context, chatID string, limit int) ([]model.ChatMessage, error) {
	filter := bson.M{"chat_id": chatID}
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.messagesCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entities []chatMessageMongoEntity
	if err := cursor.All(ctx, &entities); err != nil {
		return nil, err
	}

	messages := make([]model.ChatMessage, len(entities))
	for i, entity := range entities {
		messages[i] = toDomainMessage(entity)
	}

	return messages, nil
}

func (r *ChatRepository) EditMessage(ctx context.Context, message *model.ChatMessage) error {
	oid, err := primitive.ObjectIDFromHex(message.Id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"content":    message.Content,
			"is_edited":  true,
			"updated_at": message.UpdatedAt,
		},
	}

	_, err = r.messagesCol.UpdateByID(ctx, oid, update)
	return err
}

func (r *ChatRepository) DeleteMessage(ctx context.Context, messageID string) error {
	oid, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return err
	}
	_, err = r.messagesCol.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *ChatRepository) SendReadReceipt(ctx context.Context, chatID, userID, messageID string) error {
	return nil
}

func (r *ChatRepository) SendTypingStatus(ctx context.Context, chatID, userID string, isTyping bool) error {
	return nil
}

func (r *ChatRepository) AddParticipantToGroupChat(ctx context.Context, chatID, userID string) error {
	oid, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return err
	}
	update := bson.M{
		"$addToSet": bson.M{
			"participant_ids": userID,
		},
	}
	_, err = r.chatsCol.UpdateByID(ctx, oid, update)
	return err
}

func (r *ChatRepository) RemoveParticipantFromGroupChat(ctx context.Context, chatID, userID string) error {
	oid, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return err
	}
	update := bson.M{
		"$pull": bson.M{
			"participant_ids": userID,
		},
	}
	_, err = r.chatsCol.UpdateByID(ctx, oid, update)
	return err
}

func (r *ChatRepository) EditGroupChat(ctx context.Context, chat *model.Chat) error {
	oid, err := primitive.ObjectIDFromHex(chat.Id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name": chat.Name,
		},
	}
	_, err = r.chatsCol.UpdateByID(ctx, oid, update)
	return err
}

func (r *ChatRepository) DeleteGroupChat(ctx context.Context, chatID string) error {
	oid, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return err
	}
	_, err = r.chatsCol.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	_, err = r.messagesCol.DeleteMany(ctx, bson.M{"chat_id": chatID})
	return err
}

func (r *ChatRepository) LeaveGroupChat(ctx context.Context, chatID, userID string) error {
	oid, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return err
	}
	update := bson.M{
		"$pull": bson.M{
			"participant_ids": userID,
		},
	}
	_, err = r.chatsCol.UpdateByID(ctx, oid, update)
	return err
}
