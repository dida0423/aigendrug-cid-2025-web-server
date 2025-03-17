package session

import (
	"context"

	"github.com/gin-gonic/gin"
	gocql "github.com/gocql/gocql"
)

func SetupSessionRoutes(c context.Context, router *gin.Engine, db *gocql.Session) {
	sessionService := NewSessionService(c, db)
	sessionController := NewSessionController(sessionService)

	sessionRoutes := router.Group("/v1/session")
	{
		sessionRoutes.GET("", sessionController.GetSessions)
		sessionRoutes.POST("/:name", sessionController.CreateSession)
		sessionRoutes.DELETE("/:id", sessionController.DeleteSession)
	}
}
