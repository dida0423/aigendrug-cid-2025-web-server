package session

import (
	"time"

	gocql "github.com/gocql/gocql"
)

const (
	SessionStatusActive = "active"
)

type Session struct {
	ID             gocql.UUID `json:"id"`
	Name           string     `json:"name"`
	Status         string     `json:"status"`
	ToolStatus     string     `json:"tool_status"`
	AssignedToolID gocql.UUID `json:"assigned_tool_id"`
	CreatedAt      time.Time  `json:"created_at"`
}
