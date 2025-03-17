package session

import (
	"context"
	"time"

	gocql "github.com/gocql/gocql"
)

type SessionService interface {
	ReadAllSessions(rctx context.Context) ([]*Session, error)
	CreateSession(rctx context.Context, name string) (*Session, error)
	DeleteSession(rctx context.Context, id gocql.UUID) error
}

type sessionService struct {
	ctx context.Context
	db  *gocql.Session
}

func NewSessionService(c context.Context, db *gocql.Session) SessionService {
	return &sessionService{ctx: c, db: db}
}

func (s *sessionService) ReadAllSessions(rctx context.Context) ([]*Session, error) {
	var sessions []*Session
	query := s.db.Query("SELECT id, name, status, tool_status, assigned_tool_id, created_at FROM sessions").WithContext(rctx)

	iter := query.Iter()

	for {
		var session Session
		if !iter.Scan(&session.ID,
			&session.Name,
			&session.Status,
			&session.ToolStatus,
			&session.AssignedToolID,
			&session.CreatedAt,
		) {
			break
		}
		sessions = append(sessions, &session)
	}

	if len(sessions) == 0 {
		return []*Session{}, nil
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (s *sessionService) CreateSession(rctx context.Context, name string) (*Session, error) {
	newUUID, err := gocql.RandomUUID()
	if err != nil {
		return nil, err
	}

	if err := s.db.Query(`
		INSERT INTO 
			sessions
			(id, name, status, created_at) 
		VALUES 
			(?, ?, ?, toTimeStamp(now()))`,
		newUUID, name, SessionStatusActive).WithContext(rctx).Exec(); err != nil {
		return nil, err
	}

	return &Session{
		ID:             newUUID,
		Name:           name,
		Status:         "active",
		ToolStatus:     "",
		AssignedToolID: gocql.UUID{},
		CreatedAt:      time.Now(),
	}, nil
}

func (s *sessionService) DeleteSession(rctx context.Context, id gocql.UUID) error {
	if err := s.db.Query("DELETE FROM sessions WHERE id = ?", id).WithContext(rctx).Exec(); err != nil {
		return err
	}

	return nil
}
