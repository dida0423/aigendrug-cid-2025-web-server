package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gocql "github.com/gocql/gocql"
)

type ChatController struct {
	chatService ChatService
}

func NewChatController(ChatService ChatService) *ChatController {
	return &ChatController{chatService: ChatService}
}

// GetChatMessages godoc
// @Summary Get all chat messages for a session
// @Description Get all chat messages for a session
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionID path string true "Session ID"
// @Success 200 {array} ChatMessage
// @Failure 500 {object} map[string]interface{}
// @Router /v1/chat/message/{sessionID} [get]
func (cc *ChatController) GetChatMessages(c *gin.Context) {
	sessionID, err := gocql.ParseUUID(c.Param("sessionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatMessages, err := cc.chatService.ReadAllChatMessages(c.Request.Context(), sessionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chatMessages)
}

// CreateChatMessage godoc
// @Summary Create a new chat message
// @Description Create a new chat message
// @Tags Chat
// @Accept json
// @Produce json
// @Param chatMessage body CreateChatMessageDTO true "Create Chat Message DTO"
// @Success 200 {object} CreateChatMessageDTO
// @Failure 500 {object} map[string]interface{}
// @Router /v1/chat/message [post]
func (cc *ChatController) CreateChatMessage(c *gin.Context) {
	var chatMessage CreateChatMessageDTO
	if err := c.ShouldBindJSON(&chatMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := cc.chatService.CreateChatMessage(c.Request.Context(), &chatMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chatMessage)
}
