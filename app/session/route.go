package session

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupSessionRoutes(c context.Context, router *gin.Engine, db *pgxpool.Pool) {
	sessionService := NewSessionService(c, db)
	sessionController := NewSessionController(sessionService)

	sessionRoutes := router.Group("/v1/session")
	{
		sessionRoutes.GET("", sessionController.GetSessions)
		sessionRoutes.POST("/:name", sessionController.CreateSession)
		sessionRoutes.DELETE("/:id", sessionController.DeleteSession)
	}
}
