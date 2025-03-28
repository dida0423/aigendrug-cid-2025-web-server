package tool

import (
	"time"

	gocql "github.com/gocql/gocql"
)

const (
	ToolRoleUser      = "user"
	ToolRoleAssistant = "assistant"
	ToolRoleSystem    = "system"
)

type Tool struct {
	ID                gocql.UUID     `json:"id"`
	Name              string         `json:"name"`
	Version           string         `json:"version"`
	Description       string         `json:"description"`
	ProviderInterface map[string]any `json:"provider_interface"`
	CreatedAt         time.Time      `json:"created_at"`
}

type CreateToolDTO struct {
	ID                gocql.UUID     `json:"id"`
	Name              string         `json:"name"`
	Version           string         `json:"version"`
	Description       string         `json:"description"`
	ProviderInterface map[string]any `json:"provider_interface"`
}

type ToolMessage struct {
	ID        gocql.UUID     `json:"id"`
	SessionID gocql.UUID     `json:"session_id"`
	ToolID    gocql.UUID     `json:"tool_id"`
	Role      string         `json:"role"`
	Data      map[string]any `json:"data"`
	CreatedAt time.Time      `json:"created_at"`
}

type CreateToolMessageDTO struct {
	SessionID gocql.UUID     `json:"session_id"`
	ToolID    gocql.UUID     `json:"tool_id"`
	Role      string         `json:"role"`
	Data      map[string]any `json:"data"`
}
