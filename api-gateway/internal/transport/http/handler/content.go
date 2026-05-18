package handler

import (
	"net/http"

	contentpb "github.com/CoffeeSi/social-network-microservices/api-gateway/proto/content"
	"github.com/gin-gonic/gin"
)

func (g *Gateway) CreatePost(c *gin.Context) {
	userID, ok := g.requireUserID(c)
	if !ok {
		return
	}
	var req struct {
		Content   string   `json:"content"`
		MediaURLs []string `json:"media_urls"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.content.CreatePost(c.Request.Context(), &contentpb.CreatePostRequest{
		AuthorId:  userID,
		Content:   req.Content,
		MediaUrls: req.MediaURLs,
	})
	writeGRPC(c, resp, err, http.StatusCreated)
}

func (g *Gateway) ListPosts(c *gin.Context) {
	resp, err := g.content.ListPosts(c.Request.Context(), &contentpb.ListPostsRequest{
		PageSize: int32(queryInt(c, "page_size", 20)),
		Page:     int32(queryInt(c, "page", 1)),
		AuthorId: c.Query("author_id"),
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) GetMyPosts(c *gin.Context) {
	userID, ok := g.requireUserID(c)
	if !ok {
		return
	}
	resp, err := g.content.GetMyPosts(c.Request.Context(), &contentpb.GetMyPostsRequest{
		UserId:   userID,
		PageSize: int32(queryInt(c, "page_size", 20)),
		Page:     int32(queryInt(c, "page", 1)),
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) GetPost(c *gin.Context) {
	resp, err := g.content.GetPost(c.Request.Context(), &contentpb.GetPostRequest{Id: c.Param("id")})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) GetPostStats(c *gin.Context) {
	resp, err := g.content.GetPostStats(c.Request.Context(), &contentpb.GetPostStatsRequest{PostId: c.Param("id")})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) UpdatePost(c *gin.Context) {
	userID, ok := g.requireUserID(c)
	if !ok {
		return
	}
	var req struct {
		Content   string   `json:"content"`
		MediaURLs []string `json:"media_urls"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.content.UpdatePost(c.Request.Context(), &contentpb.UpdatePostRequest{
		Id:        c.Param("id"),
		UserId:    userID,
		Content:   req.Content,
		MediaUrls: req.MediaURLs,
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) DeletePost(c *gin.Context) {
	userID, ok := g.requireUserID(c)
	if !ok {
		return
	}
	_, err := g.content.DeletePost(c.Request.Context(), &contentpb.DeletePostRequest{
		Id:     c.Param("id"),
		UserId: userID,
	})
	writeEmptyGRPC(c, err)
}

func (g *Gateway) CreateComment(c *gin.Context) {
	userID, ok := g.requireUserID(c)
	if !ok {
		return
	}
	var req struct {
		Text string `json:"text"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.content.CreateComment(c.Request.Context(), &contentpb.CreateCommentRequest{
		PostId: c.Param("id"),
		UserId: userID,
		Text:   req.Text,
	})
	writeGRPC(c, resp, err, http.StatusCreated)
}

func (g *Gateway) ListComments(c *gin.Context) {
	resp, err := g.content.ListComments(c.Request.Context(), &contentpb.ListCommentsRequest{
		PostId:   c.Param("id"),
		PageSize: int32(queryInt(c, "page_size", 20)),
		Page:     int32(queryInt(c, "page", 1)),
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) UpdateComment(c *gin.Context) {
	userID, ok := g.requireUserID(c)
	if !ok {
		return
	}
	var req struct {
		Text string `json:"text"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.content.UpdateComment(c.Request.Context(), &contentpb.UpdateCommentRequest{
		Id:     c.Param("id"),
		UserId: userID,
		Text:   req.Text,
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) DeleteComment(c *gin.Context) {
	userID, ok := g.requireUserID(c)
	if !ok {
		return
	}
	_, err := g.content.DeleteComment(c.Request.Context(), &contentpb.DeleteCommentRequest{
		Id:     c.Param("id"),
		UserId: userID,
	})
	writeEmptyGRPC(c, err)
}

func (g *Gateway) ToggleLike(c *gin.Context) {
	userID, ok := g.requireUserID(c)
	if !ok {
		return
	}
	resp, err := g.content.ToggleLike(c.Request.Context(), &contentpb.ToggleLikeRequest{
		PostId: c.Param("id"),
		UserId: userID,
	})
	writeGRPC(c, resp, err, http.StatusOK)
}
