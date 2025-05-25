package chat

import (
	"context"

	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatService interface {
	ReadAllChatMessages(rctx context.Context, sessionID uuid.UUID) ([]*ChatMessage, error)
	CreateChatMessage(rctx context.Context, chatMessage *CreateChatMessageDTO) error
}

type chatService struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewChatService(c context.Context, db *pgxpool.Pool) ChatService {
	return &chatService{ctx: c, db: db}
}

func (s *chatService) ReadAllChatMessages(rctx context.Context, sessionID uuid.UUID) ([]*ChatMessage, error) {
	rows, err := s.db.Query(rctx, `
        SELECT
            id,
            session_id,
            role,
            message,
            created_at,
            message_type,
            linked_tool_ids
        FROM
            chat_messages
        WHERE
            session_id = $1
        ORDER BY created_at ASC
    `, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatMessages []*ChatMessage
	for rows.Next() {
		var chatMessage ChatMessage
		err := rows.Scan(
			&chatMessage.ID,
			&chatMessage.SessionID,
			&chatMessage.Role,
			&chatMessage.Message,
			&chatMessage.CreatedAt,
			&chatMessage.MessageType,
			&chatMessage.LinkedToolIDs,
		)
		if err != nil {
			return nil, err
		}
		chatMessages = append(chatMessages, &chatMessage)
	}
	return chatMessages, nil
}

func (s *chatService) CreateChatMessage(rctx context.Context, chatMessage *CreateChatMessageDTO) error {
	newUUID := uuid.New()
	_, err := s.db.Exec(rctx, `
        INSERT INTO chat_messages
            (id, session_id, role, message, created_at, message_type, linked_tool_ids)
        VALUES
            ($1, $2, $3, $4, $5, $6, $7)
    `,
		newUUID,
		chatMessage.SessionID,
		chatMessage.Role,
		chatMessage.Message,
		time.Now(),
		chatMessage.MessageType,
		chatMessage.LinkedToolIDs,
	)
	return err
}
