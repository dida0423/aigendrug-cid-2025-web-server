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
		toolRoutes.GET("/", toolController.GetTools)
		toolRoutes.POST("/", toolController.CreateTool)
	}
}
