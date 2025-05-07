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

// GetTools godoc
// @Summary Get all tools
// @Description Get all tools
// @Tags Tools
// @Accept json
// @Produce json
// @Success 200 {array} Tool
// @Failure 500 {object} map[string]interface{}
// @Router /v1/tool/ [get]
func (sc *ToolController) GetTools(c *gin.Context) {
	tools, err := sc.toolService.ReadAllTools(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tools)
}

// GetTool godoc
// @Summary Get a tool by ID
// @Description Get a tool by ID
// @Tags Tools
// @Accept json
// @Produce json
// @Param id path string true "Tool ID"
// @Success 200 {object} Tool
// @Failure 500 {object} map[string]interface{}
// @Router /v1/tool/{id} [get]
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

// CreateTool godoc
// @Summary Create a new tool
// @Description Create a new tool
// @Tags Tools
// @Accept json
// @Produce json
// @Param dto body CreateToolDTO true "Tool"
// @Success 201 {object} CreateToolDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/tool [post]
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

// DeleteTool godoc
// @Summary Delete a tool
// @Description Delete a tool by ID
// @Tags Tools
// @Accept json
// @Produce json
// @Param id path string true "Tool ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/tool/{id} [delete]
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

// GetSessionToolMessages godoc
// @Summary Get all tool messages for a session
// @Description Get all tool messages for a session
// @Tags Tools
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Success 200 {array} ToolMessage
// @Failure 500 {object} map[string]interface{}
// @Router /v1/tool/messages/{session_id} [get]
func (sc *ToolController) GetSessionToolMessages(c *gin.Context) {
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

// GetToolMessages godoc
// @Summary Get all tool messages
// @Description Get all tool messages
// @Tags Tools
// @Accept json
// @Produce json
// @Success 200 {array} ToolMessage
// @Failure 500 {object} map[string]interface{}
// @Router /v1/tool/messages/ [get]
func (sc *ToolController) GetToolMessages(c *gin.Context) {
	messages, err := sc.toolService.ReadAllToolMessages(c.Request.Context(), gocql.UUID{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// CreateToolMessage godoc
// @Summary Create a new tool message
// @Description Create a new tool message
// @Tags Tools
// @Accept json
// @Produce json
// @Param dto body CreateToolMessageDTO true "tool message"
// @Success 201 {object} CreateToolMessageDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/tool/messages [post]
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

// SendRequestToToolServer godoc
// @Summary Send a request to the tool server
// @Description Send a request to the tool server
// @Tags Tools
// @Accept json
// @Produce json
// @Param id path string true "Tool ID"
// @Param request body []ToolInteractionElement true "Request Body"
// @Success 200 {object} []ToolInteractionElement
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/tool/send_request/{id} [post]
func (sc *ToolController) SendRequestToToolServer(c *gin.Context) {
	toolID, err := gocql.ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var toolRequestDTO []ToolInteractionElement
	if err := c.ShouldBindJSON((&toolRequestDTO)); err != nil {
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	response, err := sc.toolService.SendRequestToToolServer(c.Request.Context(), toolID, toolRequestDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
