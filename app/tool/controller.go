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

	err = sc.toolService.DeleteTool(c.Request.Context(), toolID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
