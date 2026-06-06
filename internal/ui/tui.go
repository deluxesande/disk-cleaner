package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type sessionState uint

const (
	stateMenu sessionState = iota
	stateAction
)

type mainModel struct {
	state   sessionState
	choices []string
	cursor  int
	action  string
}

func InitialModel() mainModel {
	return mainModel{
		state: stateMenu,
		choices: []string{
			"Quick Sweep (Dev Artifacts & Caches)",
			"Deep Scan (Cryptographic Duplicates)",
			"Clear System Temp Files",
			"Configure Exclusion Rules",
			"Exit",
		},
		cursor: 0,
		action: "",
	}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.state == stateMenu {
				if m.cursor > 0 {
					m.cursor--
				} else {
					m.cursor = len(m.choices) - 1
				}
			}

		case "down", "j":
			if m.state == stateMenu {
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			}

		case "enter", " ":
			switch m.state {
			case stateMenu:
				if m.choices[m.cursor] == "Exit" {
					return m, tea.Quit
				}
				m.action = m.choices[m.cursor]
				m.state = stateAction
			case stateAction:
				m.state = stateMenu
				m.action = ""
			}

		case "esc":
			if m.state == stateAction {
				m.state = stateMenu
				m.action = ""
			}
		}
	}
	return m, nil
}

func (m mainModel) View() string {
	var b strings.Builder

	switch m.state {
	case stateMenu:
		b.WriteString("\nSelect an operation mode:\n\n")

		for i, choice := range m.choices {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}
			b.WriteString(fmt.Sprintf("%s%s\n", cursor, choice))
		}

		b.WriteString("\n------------------------------------------------------\n")
		b.WriteString("Navigation: [j/k] or [Up/Down]  |  Select: [Enter]  |  Quit: [q]\n")
	case stateAction:
		b.WriteString(fmt.Sprintf("\nExecuting module: %s\n", m.action))
		b.WriteString("\n[ Function Not Implemented Yet ]\n")
		b.WriteString("\n------------------------------------------------------\n")
		b.WriteString("Return to Menu: [Enter] or [Esc]  |  Quit: [q]\n")
	}

	return b.String()
}
