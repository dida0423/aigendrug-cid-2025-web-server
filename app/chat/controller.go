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
