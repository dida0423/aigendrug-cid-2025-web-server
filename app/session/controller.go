package session

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gocql "github.com/gocql/gocql"
)

type SessionController struct {
	sessionService SessionService
}

func NewSessionController(sessionService SessionService) *SessionController {
	return &SessionController{sessionService: sessionService}
}

// GetSessions godoc
// @Summary Get all sessions
// @Description Get all sessions
// @Tags Sessions
// @Accept json
// @Produce json
// @Success 200 {array} Session
// @Failure 500 {object} map[string]interface{}
// @Router /v1/session [get]
func (sc *SessionController) GetSessions(c *gin.Context) {
	sessions, err := sc.sessionService.ReadAllSessions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

// CreateSession godoc
// @Summary Create a new session
// @Description Create a session with the given name
// @Tags Sessions
// @Accept json
// @Produce json
// @Param name path string true "Session Name"
// @Success 200 {object} Session
// @Failure 500 {object} map[string]interface{}
// @Router /v1/session/{name} [post]
func (sc *SessionController) CreateSession(c *gin.Context) {
	name := c.Param("name")
	session, err := sc.sessionService.CreateSession(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

// DeleteSession godoc
// @Summary Delete a session
// @Description Delete a session by its ID
// @Tags Sessions
// @Param id path string true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/session/{id} [delete]
func (sc *SessionController) DeleteSession(c *gin.Context) {
	sessionID, err := gocql.ParseUUID(c.Param("id"))
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
