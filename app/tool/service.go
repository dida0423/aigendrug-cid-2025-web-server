package tool

import (
	"context"
	"encoding/json"

	gocql "github.com/gocql/gocql"
)

type ToolService interface {
	ReadAllTools(rctx context.Context) ([]*Tool, error)
	CreateTool(rctx context.Context, dto *CreateToolDTO) error
	DeleteTool(rctx context.Context, id string) error
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

		json.Unmarshal([]byte(providerInterfaceStr), &Tool.ProviderInterface)

		Tools = append(Tools, &Tool)
	}

	if len(Tools) == 0 {
		return []*Tool{}, nil
	}

	return Tools, nil
}

func (s *toolService) CreateTool(rctx context.Context, dto *CreateToolDTO) error {
	var providerInterfaceStr []byte
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

func (s *toolService) DeleteTool(rctx context.Context, id string) error {
	query := s.db.Query("DELETE FROM tools WHERE id = ?", id).WithContext(rctx)
	err := query.Exec()
	if err != nil {
		return err
	}

	return nil
}
