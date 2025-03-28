package tool

import (
	"time"

	gocql "github.com/gocql/gocql"
)

type Tool struct {
	ID                gocql.UUID             `json:"id"`
	Name              string                 `json:"name"`
	Version           string                 `json:"version"`
	Description       string                 `json:"description"`
	ProviderInterface map[string]interface{} `json:"provider_interface"`
	CreatedAt         time.Time              `json:"created_at"`
}

type CreateToolDTO struct {
	ID                gocql.UUID             `json:"id"`
	Name              string                 `json:"name"`
	Version           string                 `json:"version"`
	Description       string                 `json:"description"`
	ProviderInterface map[string]interface{} `json:"provider_interface"`
}
