package app

import (
	"context"

	"aigendrug.com/aigendrug-cid-2025-server/app/chat"
	"aigendrug.com/aigendrug-cid-2025-server/app/session"
	"aigendrug.com/aigendrug-cid-2025-server/app/tool"
	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(c context.Context, router *gin.Engine, db *pgxpool.Pool) {
	chat.SetupChatRoutes(c, router, db)
	session.SetupSessionRoutes(c, router, db)
	tool.SetupToolRoutes(c, router, db)
}
