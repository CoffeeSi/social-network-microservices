package http

import (
	"github.com/CoffeeSi/social-network-microservices/api-gateway/internal/transport/http/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(g *handler.Gateway) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), cors())

	router.GET("/health", g.Health)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		auth.POST("/register", g.Register)
		auth.POST("/login", g.Login)
		auth.POST("/refresh", g.Refresh)

		users := v1.Group("/users")
		users.POST("", g.CreateUser)
		users.GET("", g.ListUsers)
		users.GET("/:id", g.GetUser)
		users.PATCH("/:id", g.UpdateUser)
		users.DELETE("/:id", g.DeleteUser)
		users.PUT("/:id/password", g.ChangePassword)

		posts := v1.Group("/posts")
		posts.POST("", g.CreatePost)
		posts.GET("", g.ListPosts)
		posts.GET("/:id", g.GetPost)
		posts.PATCH("/:id", g.UpdatePost)
		posts.DELETE("/:id", g.DeletePost)
		posts.POST("/:id/comments", g.CreateComment)
		posts.GET("/:id/comments", g.ListComments)
		posts.POST("/:id/like", g.ToggleLike)

		comments := v1.Group("/comments")
		comments.PATCH("/:id", g.UpdateComment)
		comments.DELETE("/:id", g.DeleteComment)

		chats := v1.Group("/chats")
		chats.GET("", g.GetChats)
		chats.POST("/direct", g.CreateDirectChat)
		chats.POST("/group", g.CreateGroupChat)
		chats.PATCH("/:id", g.EditGroupChat)
		chats.DELETE("/:id", g.DeleteGroupChat)
		chats.POST("/:id/leave", g.LeaveGroupChat)
		chats.POST("/:id/participants", g.AddParticipant)
		chats.DELETE("/:id/participants/:user_id", g.RemoveParticipant)
		chats.POST("/:id/messages", g.SendMessage)
		chats.GET("/:id/messages", g.GetMessages)
		chats.PATCH("/:id/messages/:message_id", g.EditMessage)
		chats.DELETE("/:id/messages/:message_id", g.DeleteMessage)
		chats.POST("/:id/messages/:message_id/read", g.SendReadReceipt)
		chats.POST("/:id/typing", g.SendTypingStatus)
	}

	return router
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
