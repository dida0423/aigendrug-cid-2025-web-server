package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionService interface {
	ReadAllSessions(rctx context.Context) ([]*Session, error)
	CreateSession(rctx context.Context, name string) (*Session, error)
	DeleteSession(rctx context.Context, id uuid.UUID) error
}

type sessionService struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewSessionService(c context.Context, db *pgxpool.Pool) SessionService {
	return &sessionService{ctx: c, db: db}
}

func (s *sessionService) ReadAllSessions(rctx context.Context) ([]*Session, error) {
	rows, err := s.db.Query(rctx, `
        SELECT id, name, status, tool_status, assigned_tool_id, created_at
        FROM sessions
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		var session Session
		err := rows.Scan(
			&session.ID,
			&session.Name,
			&session.Status,
			&session.ToolStatus,
			&session.AssignedToolID,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	if len(sessions) == 0 {
		return []*Session{}, nil
	}

	return sessions, nil
}

func (s *sessionService) CreateSession(rctx context.Context, name string) (*Session, error) {
	newUUID := uuid.New()
	createdAt := time.Now()

	_, err := s.db.Exec(rctx, `
        INSERT INTO sessions (id, name, status, created_at)
        VALUES ($1, $2, $3, $4)
    `, newUUID, name, SessionStatusActive, createdAt)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:             newUUID,
		Name:           name,
		Status:         "active",
		ToolStatus:     "",
		AssignedToolID: uuid.UUID{},
		CreatedAt:      time.Now(),
	}, nil
}

func (s *sessionService) DeleteSession(rctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(rctx, "DELETE FROM sessions WHERE id = $1", id)
	return err
}
