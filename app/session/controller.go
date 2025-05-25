package session

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionController struct {
	sessionService SessionService
}

func NewSessionController(sessionService SessionService) *SessionController {
	return &SessionController{sessionService: sessionService}
}

func (sc *SessionController) GetSessions(c *gin.Context) {
	sessions, err := sc.sessionService.ReadAllSessions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func (sc *SessionController) CreateSession(c *gin.Context) {
	name := c.Param("name")
	session, err := sc.sessionService.CreateSession(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (sc *SessionController) DeleteSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sc.sessionService.DeleteSession(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session deleted"})
}
