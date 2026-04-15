package tui

import (
	"fmt"
)

type Model struct {
	// TODO: Bubbletea model
}

func New() *Model {
	return &Model{}
}

func (m *Model) Start() error {
	// TODO: Launch bubbletea TUI
	fmt.Println("(TUI not yet implemented)")
	return nil
}
