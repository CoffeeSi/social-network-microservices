package handler

import (
	"net/http"

	chatpb "github.com/CoffeeSi/social-network-microservices/api-gateway/proto/chat"
	"github.com/gin-gonic/gin"
)

func (g *Gateway) GetChats(c *gin.Context) {
	resp, err := g.chat.GetChats(authContext(c), &chatpb.GetChatsRequest{})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) CreateDirectChat(c *gin.Context) {
	var req struct {
		TargetUserID string `json:"target_user_id"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.chat.CreateDirectChat(authContext(c), &chatpb.CreateDirectChatRequest{TargetUserId: req.TargetUserID})
	writeGRPC(c, resp, err, http.StatusCreated)
}

func (g *Gateway) CreateGroupChat(c *gin.Context) {
	var req struct {
		Name           string   `json:"name"`
		ParticipantIDs []string `json:"participant_ids"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.chat.CreateGroupChat(authContext(c), &chatpb.CreateGroupChatRequest{
		Name:           req.Name,
		ParticipantIds: req.ParticipantIDs,
	})
	writeGRPC(c, resp, err, http.StatusCreated)
}

func (g *Gateway) SendMessage(c *gin.Context) {
	var req struct {
		Content string `json:"content"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.chat.SendMessage(authContext(c), &chatpb.SendMessageRequest{
		ChatId:  c.Param("id"),
		Content: req.Content,
	})
	writeGRPC(c, resp, err, http.StatusCreated)
}

func (g *Gateway) GetMessages(c *gin.Context) {
	resp, err := g.chat.GetMessages(authContext(c), &chatpb.GetMessagesRequest{
		ChatId:   c.Param("id"),
		Page:     int32(queryInt(c, "page", 1)),
		PageSize: int32(queryInt(c, "page_size", 20)),
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) EditMessage(c *gin.Context) {
	var req struct {
		NewContent string `json:"new_content"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.chat.EditMessage(authContext(c), &chatpb.EditMessageRequest{
		ChatId:     c.Param("id"),
		MessageId:  c.Param("message_id"),
		NewContent: req.NewContent,
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) DeleteMessage(c *gin.Context) {
	resp, err := g.chat.DeleteMessage(authContext(c), &chatpb.DeleteMessageRequest{
		ChatId:    c.Param("id"),
		MessageId: c.Param("message_id"),
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) SendReadReceipt(c *gin.Context) {
	resp, err := g.chat.SendReadReceipt(authContext(c), &chatpb.SendReadReceiptRequest{
		ChatId:    c.Param("id"),
		MessageId: c.Param("message_id"),
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) SendTypingStatus(c *gin.Context) {
	var req struct {
		IsTyping bool `json:"is_typing"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.chat.SendTypingStatus(authContext(c), &chatpb.SendTypingStatusRequest{
		ChatId:   c.Param("id"),
		IsTyping: req.IsTyping,
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) AddParticipant(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.chat.AddParticipantToGroupChat(authContext(c), &chatpb.AddParticipantToGroupChatRequest{
		ChatId: c.Param("id"),
		UserId: req.UserID,
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) RemoveParticipant(c *gin.Context) {
	resp, err := g.chat.RemoveParticipantFromGroupChat(authContext(c), &chatpb.RemoveParticipantFromGroupChatRequest{
		ChatId: c.Param("id"),
		UserId: c.Param("user_id"),
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) EditGroupChat(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.chat.EditGroupChat(authContext(c), &chatpb.EditGroupChatRequest{
		ChatId: c.Param("id"),
		Name:   req.Name,
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) DeleteGroupChat(c *gin.Context) {
	resp, err := g.chat.DeleteGroupChat(authContext(c), &chatpb.DeleteGroupChatRequest{ChatId: c.Param("id")})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) LeaveGroupChat(c *gin.Context) {
	resp, err := g.chat.LeaveGroupChat(authContext(c), &chatpb.LeaveGroupChatRequest{ChatId: c.Param("id")})
	writeGRPC(c, resp, err, http.StatusOK)
}
