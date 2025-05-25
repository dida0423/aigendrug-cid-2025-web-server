package toolrouter

import "github.com/google/uuid"

type SelectToolRequestDTO struct {
	UserPrompt string `json:"user_prompt"`
}

type SelectToolResponseDTO struct {
	SelectedToolName string `json:"selected_tool_name"`
	SelectedToolID   string `json:"selected_tool_id"`
}

type SelectedTool struct {
	ToolName string    `json:"tool_name"`
	ToolID   uuid.UUID `json:"tool_id"`
}
