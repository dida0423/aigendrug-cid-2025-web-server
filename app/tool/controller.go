package tool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type ToolController struct {
	toolService ToolService
}

func NewToolController(toolService ToolService) *ToolController {
	return &ToolController{toolService: toolService}
}

func (sc *ToolController) GetTools(c *gin.Context) {
	tools, err := sc.toolService.ReadAllTools(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tools)
}

func (sc *ToolController) GetTool(c *gin.Context) {
	toolID, err := gocql.ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tool, err := sc.toolService.ReadTool(c.Request.Context(), toolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tool)
}

func (sc *ToolController) CreateTool(c *gin.Context) {
	var dto CreateToolDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := sc.toolService.CreateTool(c.Request.Context(), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func (sc *ToolController) DeleteTool(c *gin.Context) {
	toolID, err := gocql.ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = sc.toolService.DeleteTool(c.Request.Context(), toolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func (sc *ToolController) GetToolMessages(c *gin.Context) {
	sessionID, err := gocql.ParseUUID(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messages, err := sc.toolService.ReadAllToolMessages(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (sc *ToolController) CreateToolMessage(c *gin.Context) {
	var dto CreateToolMessageDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := sc.toolService.CreateToolMessage(c.Request.Context(), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func (sc *ToolController) SendRequestToToolServer(c *gin.Context) {
	toolID, err := gocql.ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody map[string]any
	if err := c.ShouldBindJSON((&reqBody)); err != nil {
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	response, err := sc.toolService.SendRequestToToolServer(c.Request.Context(), toolID, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
