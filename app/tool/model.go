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
	ID                gocql.UUID        `json:"id"`
	Name              string            `json:"name"`
	Version           string            `json:"version"`
	Description       string            `json:"description"`
	ProviderInterface ProviderInterface `json:"provider_interface"`
	CreatedAt         time.Time         `json:"created_at"`
}

type CreateToolDTO struct {
	ID                gocql.UUID        `json:"id" validate:"required"`
	Name              string            `json:"name" validate:"required"`
	Version           string            `json:"version"`
	Description       string            `json:"description" validate:"required"`
	ProviderInterface ProviderInterface `json:"provider_interface" validate:"required"`
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

type ProviderInterface struct {
	URL                 string             `json:"url" valdate:"required,url"`
	AuthStrategy        string             `json:"authStrategy" validate:"required"`
	RequestMethod       string             `json:"requestMethod" validate:"required,oneof=GET POST PUT DELETE"`
	RequestContentType  string             `json:"requestContentType" validate:"required"`
	ResponseContentType string             `json:"responseContentType" validate:"required"`
	RequestInterface    []InterfaceElement `json:"requestInterface" validate:"required,min=1,dive"`
	ResponseInterface   []InterfaceElement `json:"responseInterface" validate:"required,min=1,dive"`
}

type InterfaceElement struct {
	ID                string            `json:"id" validate:"required"`
	Type              string            `json:"type" validate:"required,oneof=body query header"`
	Required          bool              `json:"required"`
	Key               string            `json:"key" validate:"required"`
	ValueType         string            `json:"valueType" validate:"required,oneof=string number boolean"`
	BindedElementType BindedElementType `json:"bindedElementType" validate:"required"`
}

type BindedElementType struct {
	Label           string `json:"label" validate:"required"`
	HTMLElementType string `json:"htmlElementType" validate:"required"`
	ValueType       string `json:"valueType" validate:"required,oneof=string number boolean"`
}
