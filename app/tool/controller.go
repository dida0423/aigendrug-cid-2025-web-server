package tool

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func (sc *ToolController) CreateTool(c *gin.Context) {
	var dto CreateToolDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := sc.toolService.CreateTool(c.Request.Context(), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"tool_id": id})
}
