package toolrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type ToolRouterService interface {
	SelectTool(prompt string) (*SelectedTool, error)
}

type toolRouterService struct {
	ctx  context.Context
	host string
}

func NewToolRouterService(c context.Context) ToolRouterService {
	return &toolRouterService{ctx: c, host: os.Getenv("TOOL_ROUTER_HOST")}
}

func (trs *toolRouterService) SelectTool(prompt string) (*SelectedTool, error) {
	req := SelectToolRequestDTO{
		UserPrompt: prompt,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %s", err)
	}

	res, err := http.Post(trs.host+"/select", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to select tool: %s", res.Status)
	}

	var response SelectToolResponseDTO
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	toolID, err := uuid.Parse(response.SelectedToolID)
	if err != nil {
		return nil, fmt.Errorf("invalid tool ID format: %s", err)
	}

	return &SelectedTool{
		ToolName: response.SelectedToolName,
		ToolID:   toolID,
	}, nil
}
