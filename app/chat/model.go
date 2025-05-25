package chat

import (
	"time"

	"github.com/google/uuid"
)

const (
	ChatMessageTypeNormal          = 0
	ChatMessageTypeToolSelection   = 1
	ChatMessageTypeToolSuggestions = 2
	ChatMessageTypeToolFurtherInfo = 3
)

const (
	ChatRoleUser      = "user"
	ChatRoleAssistant = "assistant"
	ChatRoleSystem    = "system"
)

type ChatMessage struct {
	ID            uuid.UUID   `json:"id"`
	SessionID     uuid.UUID   `json:"session_id"`
	Role          string      `json:"role"`
	Message       string      `json:"message"`
	CreatedAt     time.Time   `json:"created_at"`
	MessageType   int         `json:"message_type"`
	LinkedToolIDs []uuid.UUID `json:"linked_tool_ids"`
}

type CreateChatMessageDTO struct {
	SessionID     uuid.UUID   `json:"session_id"`
	Role          string      `json:"role"`
	Message       string      `json:"message"`
	MessageType   int         `json:"message_type"`
	LinkedToolIDs []uuid.UUID `json:"linked_tool_ids"`
}
