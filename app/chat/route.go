package chat

import (
	"context"

	"github.com/gin-gonic/gin"
	gocql "github.com/gocql/gocql"
)

func SetupChatRoutes(c context.Context, router *gin.Engine, db *gocql.Session) {
	chatService := NewChatService(c, db)
	chatController := NewChatController(chatService)

	chatRoutes := router.Group("/v1/chat")
	{
		chatRoutes.GET("/message/:sessionID", chatController.GetChatMessages)
		chatRoutes.POST("/message", chatController.CreateChatMessage)

		chatRoutes.GET("/session/ws", func(ctx *gin.Context) {
			WebSocketHandler(ctx, db)
		})
	}
}
