package agent

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tylerdotai/saia/internal/config"
)

type Agent struct {
	cfg     *config.Config
	db      *sql.DB
}

func New(cfg *config.Config, database *sql.DB) *Agent {
	return &Agent{cfg: cfg, db: database}
}

func (a *Agent) Prompt(ctx context.Context, p Prompt) (*Response, error) {
	// TODO: Core agent loop
	return nil, fmt.Errorf("not implemented")
}

type Prompt struct {
	SessionID  string
	Platform   string
	ChannelID  string
	UserID     string
	Content    string
	Role       string
	Metadata   map[string]any
}

type Response struct {
	Content   string
	SkillUsed string
	ToolCalls []ToolCall
	SessionID string
	Metadata  map[string]any
}

type ToolCall struct {
	Name      string
	Args      map[string]string
	Result    string
	Success   bool
}
