package app

import (
	"context"

	"aigendrug.com/aigendrug-cid-2025-server/app/chat"
	"aigendrug.com/aigendrug-cid-2025-server/app/session"
	"aigendrug.com/aigendrug-cid-2025-server/app/tool"
	"github.com/gin-gonic/gin"

	gocql "github.com/gocql/gocql"
)

func SetupRoutes(c context.Context, router *gin.Engine, db *gocql.Session) {
	chat.SetupChatRoutes(c, router, db)
	session.SetupSessionRoutes(c, router, db)
	tool.SetupToolRoutes(c, router, db)
}
