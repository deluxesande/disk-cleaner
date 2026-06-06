package ui

import (
	"fmt"
	"strings"

	"github.com/deluxesande/disk-cleaner/internal/config"
	"github.com/deluxesande/disk-cleaner/internal/dedupe"
	"github.com/deluxesande/disk-cleaner/internal/scanner"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState uint

const (
	stateMenu sessionState = iota
	stateScanning
	stateResults
)

type scanResultMsg struct {
	report string
}

type mainModel struct {
	state   sessionState
	choices []string
	cursor  int
	action  string
	result  string
	cfg     *config.AppConfig
	spinner spinner.Model // Add the spinner state
}

func InitialModel() mainModel {
	// Initialize a new spinner with a classic dot style
	s := spinner.New()
	s.Spinner = spinner.Dot

	return mainModel{
		state: stateMenu,
		choices: []string{
			"Quick Sweep (Dev Artifacts & Caches)",
			"Deep Scan (Cryptographic Duplicates)",
			"Clear System Temp Files",
			"Configure Exclusion Rules",
			"Exit",
		},
		cursor:  0,
		action:  "",
		cfg:     config.Load(),
		spinner: s,
	}
}

// Init now triggers the spinner to start ticking immediately
func (m mainModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func runSweepCmd(cfg *config.AppConfig) tea.Cmd {
	return func() tea.Msg {
		report := scanner.RunSweep(cfg)
		return scanResultMsg{report: RenderSweepReport(report)}
	}
}

func runDedupeCmd(cfg *config.AppConfig) tea.Cmd {
	return func() tea.Msg {
		duplicates := dedupe.Run(cfg)
		return scanResultMsg{report: RenderDedupeReport(duplicates)}
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle the background scan completing
	case scanResultMsg:
		m.result = msg.report
		m.state = stateResults
		return m, nil

	// Handle the spinner animation frames
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

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
				m.state = stateScanning

				switch m.cursor {
				case 0:
					return m, runSweepCmd(m.cfg)
				case 1:
					return m, runDedupeCmd(m.cfg)
				case 2:
					return m, func() tea.Msg {
						return scanResultMsg{report: "\n[ Clear System Temp Files - Module integration pending ]\n"}
					}
				case 3:
					return m, func() tea.Msg {
						return scanResultMsg{report: "\n[ Configure Exclusion Rules - Module integration pending ]\n"}
					}
				}

			case stateResults:
				m.state = stateMenu
				m.action = ""
				m.result = ""
			}

		case "esc":
			if m.state == stateResults {
				m.state = stateMenu
				m.action = ""
				m.result = ""
			}
		}
	}
	return m, nil
}

func (m mainModel) View() string {
	var b strings.Builder

	switch m.state {
	case stateMenu:
		// Display the current target directory right at the top
		b.WriteString(fmt.Sprintf("\nTarget Directory: %s\n", m.cfg.TargetDir))
		b.WriteString("------------------------------------------------------\n")
		b.WriteString("Select an operation mode:\n\n")

		for i, choice := range m.choices {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}
			b.WriteString(fmt.Sprintf("%s%s\n", cursor, choice))
		}

		b.WriteString("\n------------------------------------------------------\n")
		b.WriteString("Navigation: [j/k] or [Up/Down]  |  Select: [Enter]  |  Quit: [q]\n")

	case stateScanning:
		b.WriteString(fmt.Sprintf("\nExecuting module: %s\n", m.action))
		b.WriteString(fmt.Sprintf("Scanning Target:  %s\n", m.cfg.TargetDir))

		// Render the spinning animation
		b.WriteString(fmt.Sprintf("\n %s Scanning in progress... Please wait.\n", m.spinner.View()))

		b.WriteString("\n------------------------------------------------------\n")

	case stateResults:
		b.WriteString(fmt.Sprintf("\nResults for: %s\n", m.action))
		b.WriteString(m.result)
		b.WriteString("\n------------------------------------------------------\n")
		b.WriteString("Return to Menu: [Enter] or [Esc]  |  Quit: [q]\n")
	}

	return b.String()
}
