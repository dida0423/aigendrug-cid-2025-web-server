package chat

import (
	"context"

	gocql "github.com/gocql/gocql"
)

type ChatService interface {
	ReadAllChatMessages(rctx context.Context, sessionID gocql.UUID) ([]*ChatMessage, error)
	CreateChatMessage(rctx context.Context, chatMessage *CreateChatMessageDTO) error
}

type chatService struct {
	ctx context.Context
	db  *gocql.Session
}

func NewChatService(c context.Context, db *gocql.Session) ChatService {
	return &chatService{ctx: c, db: db}
}

func (s *chatService) ReadAllChatMessages(rctx context.Context, sessionID gocql.UUID) ([]*ChatMessage, error) {
	var chatMessages []*ChatMessage
	query := s.db.Query(`
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
			session_id = ? 
		ORDER BY created_at ASC`,
		sessionID).WithContext(rctx)

	iter := query.Iter()

	defer iter.Close()

	for {
		var chatMessage ChatMessage
		if !iter.Scan(&chatMessage.ID,
			&chatMessage.SessionID,
			&chatMessage.Role,
			&chatMessage.Message,
			&chatMessage.CreatedAt,
			&chatMessage.MessageType,
			&chatMessage.LinkedToolIDs) {
			break
		}
		chatMessages = append(chatMessages, &chatMessage)
	}

	return chatMessages, nil
}

func (s *chatService) CreateChatMessage(rctx context.Context, chatMessage *CreateChatMessageDTO) error {
	newUUID, err := gocql.RandomUUID()
	if err != nil {
		return err
	}

	query := s.db.Query(`
		INSERT INTO 
			chat_messages 
			(id, session_id, role, message, created_at, message_type, linked_tool_ids) 
		VALUES 
			(?, ?, ?, ?, toTimestamp(now()), ?, ?)`,
		newUUID,
		chatMessage.SessionID,
		chatMessage.Role,
		chatMessage.Message,
		chatMessage.MessageType,
		chatMessage.LinkedToolIDs).WithContext(rctx)

	if err := query.Exec(); err != nil {
		return err
	}

	return nil
}
