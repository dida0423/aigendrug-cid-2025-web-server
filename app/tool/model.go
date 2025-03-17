package tool

import (
	"time"

	gocql "github.com/gocql/gocql"
)

type Tool struct {
	ID                gocql.UUID             `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	ImageURL          string                 `json:"image_url"`
	ProviderInterface map[string]interface{} `json:"provider_interface"`
	CreatedAt         time.Time              `json:"created_at"`
}

type CreateToolDTO struct {
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	ImageURL          string                 `json:"image_url"`
	ProviderInterface map[string]interface{} `json:"provider_interface"`
}
