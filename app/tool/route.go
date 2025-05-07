package tool

import (
	"context"

	"github.com/gin-gonic/gin"
	gocql "github.com/gocql/gocql"
)

func SetupToolRoutes(c context.Context, router *gin.Engine, db *gocql.Session) {
	toolService := NewToolService(c, db)
	toolController := NewToolController(toolService)

	toolRoutes := router.Group("/v1/tool")
	{
		toolRoutes.GET("", toolController.GetTools)
		toolRoutes.GET("/:id", toolController.GetTool)
		toolRoutes.POST("", toolController.CreateTool)
		toolRoutes.DELETE("/:id", toolController.DeleteTool)
		toolRoutes.POST("/send_request/:id", toolController.SendRequestToToolServer)
		toolRoutes.GET("/messages/:session_id", toolController.GetSessionToolMessages)
		toolRoutes.POST("/messages", toolController.CreateToolMessage)
		toolRoutes.GET("/messages", toolController.GetToolMessages)

		toolRoutes.GET("/session/ws", func(ctx *gin.Context) {
			WebSocketHandler(ctx, db)
		})
	}
}
