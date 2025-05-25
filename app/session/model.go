package session

import (
	"time"

	"github.com/google/uuid"
)

const (
	SessionStatusActive = "active"
)

type Session struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	ToolStatus     string    `json:"tool_status"`
	AssignedToolID uuid.UUID `json:"assigned_tool_id"`
	CreatedAt      time.Time `json:"created_at"`
}
