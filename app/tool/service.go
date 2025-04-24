package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	validator "github.com/go-playground/validator/v10"
	gocql "github.com/gocql/gocql"
)

type ToolService interface {
	ReadAllTools(rctx context.Context) ([]*Tool, error)
	ReadTool(rctx context.Context, id gocql.UUID) (*Tool, error)
	CreateTool(rctx context.Context, dto *CreateToolDTO) error
	DeleteTool(rctx context.Context, id gocql.UUID) error
	ReadAllToolMessages(rctx context.Context, sessionID gocql.UUID) ([]*ToolMessage, error)
	CreateToolMessage(rctx context.Context, dto *CreateToolMessageDTO) error
	SendRequestToToolServer(rctx context.Context, id gocql.UUID, userRequestBody map[string]any) (string, error)
}

type toolService struct {
	ctx context.Context
	db  *gocql.Session
}

func NewToolService(c context.Context, db *gocql.Session) ToolService {
	return &toolService{ctx: c, db: db}
}

func (s *toolService) ReadAllTools(rctx context.Context) ([]*Tool, error) {
	var Tools []*Tool
	query := s.db.Query("SELECT id, name, version, description, provider_interface, created_at FROM tools").WithContext(rctx)
	iter := query.Iter()
	defer iter.Close()

	var providerInterfaceStr string

	for {
		var Tool Tool
		if !iter.Scan(
			&Tool.ID,
			&Tool.Name,
			&Tool.Version,
			&Tool.Description,
			&providerInterfaceStr,
			&Tool.CreatedAt,
		) {
			break
		}

		if err := json.Unmarshal([]byte(providerInterfaceStr), &Tool.ProviderInterface); err != nil {
			continue
		}

		Tools = append(Tools, &Tool)
	}

	if len(Tools) == 0 {
		return []*Tool{}, nil
	}

	return Tools, nil
}

func (s *toolService) ReadTool(rctx context.Context, id gocql.UUID) (*Tool, error) {
	var Tool Tool
	query := s.db.Query("SELECT id, name, version, description, provider_interface, created_at FROM tools WHERE id = ?", id).WithContext(rctx)
	iter := query.Iter()
	defer iter.Close()

	var providerInterfaceStr string

	if !iter.Scan(
		&Tool.ID,
		&Tool.Name,
		&Tool.Version,
		&Tool.Description,
		&providerInterfaceStr,
		&Tool.CreatedAt,
	) {
		return nil, nil
	}

	if err := json.Unmarshal([]byte(providerInterfaceStr), &Tool.ProviderInterface); err != nil {
		return nil, err
	}

	return &Tool, nil
}

func (s *toolService) CreateTool(rctx context.Context, dto *CreateToolDTO) error {
	var providerInterfaceStr []byte

	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		return fmt.Errorf("tool validation failed: %w", err)
	}

	if err := validate.Struct(dto.ProviderInterface); err != nil {
		return fmt.Errorf("provider interface validation failed: %w", err)
	}

	providerInterfaceStr, err := json.Marshal(dto.ProviderInterface)
	if err != nil {
		return err
	}

	query := s.db.Query(`
		INSERT INTO tools (id, name, version, description, provider_interface, created_at)
		VALUES (?, ?, ?, ?, ?, toTimestamp(now()))
	`, dto.ID, dto.Name, dto.Version, dto.Description, string(providerInterfaceStr)).WithContext(rctx)

	err = query.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (s *toolService) DeleteTool(rctx context.Context, id gocql.UUID) error {
	query := s.db.Query("DELETE FROM tools WHERE id = ?", id).WithContext(rctx)
	err := query.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (s *toolService) ReadAllToolMessages(rctx context.Context, sessionID gocql.UUID) ([]*ToolMessage, error) {
	var ToolMessages []*ToolMessage
	query := s.db.Query("SELECT id, session_id, tool_id, role, data, created_at FROM tool_messages WHERE session_id = ?", sessionID).WithContext(rctx)
	iter := query.Iter()
	defer iter.Close()

	for {
		var ToolMessage ToolMessage
		var dataStr string

		if !iter.Scan(
			&ToolMessage.ID,
			&ToolMessage.SessionID,
			&ToolMessage.ToolID,
			&ToolMessage.Role,
			&dataStr,
			&ToolMessage.CreatedAt,
		) {
			break
		}

		if err := json.Unmarshal([]byte(dataStr), &ToolMessage.Data); err != nil {
			continue
		}

		ToolMessages = append(ToolMessages, &ToolMessage)
	}

	if len(ToolMessages) == 0 {
		return []*ToolMessage{}, nil
	}

	return ToolMessages, nil
}

func (s *toolService) CreateToolMessage(rctx context.Context, dto *CreateToolMessageDTO) error {
	var dataStr []byte
	dataStr, err := json.Marshal(dto.Data)
	if err != nil {
		return err
	}

	query := s.db.Query(`
		INSERT INTO tool_messages (session_id, tool_id, role, data, created_at)
		VALUES (?, ?, ?, ?, toTimestamp(now()))
	`, dto.SessionID, dto.ToolID, dto.Role, string(dataStr)).WithContext(rctx)

	err = query.Exec()
	if err != nil {
		return err
	}

	return nil
}
