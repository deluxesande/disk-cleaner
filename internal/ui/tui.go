package ui

import (
	"fmt"
	"strings"

	"github.com/deluxesande/disk-cleaner/internal/config"
	"github.com/deluxesande/disk-cleaner/internal/dedupe"
	"github.com/deluxesande/disk-cleaner/internal/models"
	"github.com/deluxesande/disk-cleaner/internal/scanner"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState uint

const (
	stateMenu sessionState = iota
	stateScanning
	stateInteractiveList
	stateDeleting
	stateResultSummary
	stateExclusionsList // The interactive list of custom rules
	stateExclusionsAdd  // The text input mode to add a new rule
)

type selectableItem struct {
	Path     string
	Size     int64
	Selected bool
	IsHeader bool
	Title    string
}

type scanResultMsg struct {
	items []selectableItem
}

type deleteDoneMsg struct {
	freedSpace int64
}

type mainModel struct {
	state       sessionState
	choices     []string
	cursor      int
	listCursor  int
	items       []selectableItem
	action      string
	cfg         *config.AppConfig
	spinner     spinner.Model
	textInput   textinput.Model
	freedSpace  int64
	customRules []string // Holds the loaded exclusion rules
}

func InitialModel() mainModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = "e.g. C:\\Projects\\Secret"
	ti.CharLimit = 256
	ti.Width = 50

	return mainModel{
		state: stateMenu,
		choices: []string{
			"Quick Sweep (Dev Artifacts & Caches)",
			"Deep Scan (Cryptographic Duplicates)",
			"Clear System Temp Files",
			"Configure Exclusion Rules",
			"Exit",
		},
		cursor:    0,
		cfg:       config.Load(),
		spinner:   s,
		textInput: ti,
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, textinput.Blink)
}

func buildSweepItems(report models.DiskReport) []selectableItem {
	var items []selectableItem

	if len(report.DevArtifacts) > 0 {
		items = append(items, selectableItem{IsHeader: true, Title: "Development Artifacts"})
		for _, item := range report.DevArtifacts {
			items = append(items, selectableItem{Path: item.Path, Size: item.Size})
		}
	}
	if len(report.AppCaches) > 0 {
		items = append(items, selectableItem{IsHeader: true, Title: "Application Caches"})
		for _, item := range report.AppCaches {
			items = append(items, selectableItem{Path: item.Path, Size: item.Size})
		}
	}
	if len(report.TempFiles) > 0 {
		items = append(items, selectableItem{IsHeader: true, Title: "Temporary Files"})
		for _, item := range report.TempFiles {
			items = append(items, selectableItem{Path: item.Path, Size: item.Size})
		}
	}
	return items
}

func buildDedupeItems(duplicates []models.DuplicateGroup) []selectableItem {
	var items []selectableItem
	for i, group := range duplicates {
		wasted := int64(len(group.Instances)-1) * group.FileSize
		items = append(items, selectableItem{
			IsHeader: true,
			Title:    fmt.Sprintf("Group %d (Wasted: %.2f MB)", i+1, float64(wasted)/(1024*1024)),
		})

		for j, path := range group.Instances {
			items = append(items, selectableItem{
				Path:     path,
				Size:     group.FileSize,
				Selected: j > 0,
			})
		}
	}
	return items
}

func buildTempItems(tempFiles []models.SpaceWaster) []selectableItem {
	var items []selectableItem
	if len(tempFiles) > 0 {
		items = append(items, selectableItem{IsHeader: true, Title: "System Temporary Files"})
		for _, item := range tempFiles {
			items = append(items, selectableItem{Path: item.Path, Size: item.Size})
		}
	}
	return items
}

func runSweepCmd(cfg *config.AppConfig) tea.Cmd {
	return func() tea.Msg {
		report := scanner.RunSweep(cfg)
		return scanResultMsg{items: buildSweepItems(report)}
	}
}

func runDedupeCmd(cfg *config.AppConfig) tea.Cmd {
	return func() tea.Msg {
		duplicates := dedupe.Run(cfg)
		return scanResultMsg{items: buildDedupeItems(duplicates)}
	}
}

func runTempSweepCmd() tea.Cmd {
	return func() tea.Msg {
		tempFiles := scanner.RunTempSweep()
		return scanResultMsg{items: buildTempItems(tempFiles)}
	}
}

func executeDeletionCmd(items []selectableItem) tea.Cmd {
	return func() tea.Msg {
		var freed int64
		for _, item := range items {
			if !item.IsHeader && item.Selected {
				err := DeleteItem(item.Path)
				if err == nil {
					freed += item.Size
				}
			}
		}
		return deleteDoneMsg{freedSpace: freed}
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case scanResultMsg:
		m.items = msg.items
		m.listCursor = 0
		if len(m.items) > 0 {
			if m.items[0].IsHeader && len(m.items) > 1 {
				m.listCursor = 1
			}
			m.state = stateInteractiveList
		} else {
			m.state = stateResultSummary
			m.freedSpace = 0
		}
		return m, nil

	case deleteDoneMsg:
		m.freedSpace = msg.freedSpace
		m.state = stateResultSummary
		return m, nil

	case spinner.TickMsg:
		if m.state == stateScanning || m.state == stateDeleting {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case tea.KeyMsg:
		// ISOLATED EVENT LOOP: Typing a new exclusion rule
		if m.state == stateExclusionsAdd {
			switch msg.String() {
			case "esc":
				m.state = stateExclusionsList
				m.textInput.Reset()
				m.textInput.Blur()
				return m, nil
			case "enter":
				val := strings.TrimSpace(m.textInput.Value())
				if val != "" {
					_ = config.AddCustomExclusion(val)
					m.customRules = config.GetCustomExclusions() // Refresh the list
				}
				m.textInput.Reset()
				m.textInput.Blur()
				m.state = stateExclusionsList
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
		}

		// ISOLATED EVENT LOOP: Managing the custom exclusions list
		if m.state == stateExclusionsList {
			switch msg.String() {
			case "esc":
				m.state = stateMenu
				return m, nil
			case "up", "k":
				if m.listCursor > 0 {
					m.listCursor--
				}
			case "down", "j":
				if m.listCursor < len(m.customRules)-1 {
					m.listCursor++
				}
			case "a", "n":
				m.state = stateExclusionsAdd
				m.textInput.Focus()
				return m, textinput.Blink
			case "x", "delete":
				if len(m.customRules) > 0 {
					target := m.customRules[m.listCursor]
					_ = config.RemoveCustomExclusion(target)
					m.customRules = config.GetCustomExclusions() // Refresh

					// Keep cursor in bounds if we deleted the last item
					if m.listCursor >= len(m.customRules) && m.listCursor > 0 {
						m.listCursor--
					}
				}
			case "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}

		// STANDARD EVENT LOOP: Main navigation and results list
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
			} else if m.state == stateInteractiveList {
				if m.listCursor > 0 {
					m.listCursor--
					if m.items[m.listCursor].IsHeader {
						if m.listCursor > 0 {
							m.listCursor--
						} else {
							m.listCursor++
						}
					}
				}
			}

		case "down", "j":
			if m.state == stateMenu {
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			} else if m.state == stateInteractiveList {
				if m.listCursor < len(m.items)-1 {
					m.listCursor++
					if m.items[m.listCursor].IsHeader && m.listCursor < len(m.items)-1 {
						m.listCursor++
					}
				}
			}

		case " ":
			if m.state == stateInteractiveList && !m.items[m.listCursor].IsHeader {
				m.items[m.listCursor].Selected = !m.items[m.listCursor].Selected
			}

		case "enter":
			switch m.state {
			case stateMenu:
				if m.choices[m.cursor] == "Exit" {
					return m, tea.Quit
				}
				m.action = m.choices[m.cursor]

				switch m.cursor {
				case 0:
					m.state = stateScanning
					return m, runSweepCmd(m.cfg)
				case 1:
					m.state = stateScanning
					return m, runDedupeCmd(m.cfg)
				case 2:
					m.state = stateScanning
					return m, runTempSweepCmd()
				case 3:
					m.state = stateExclusionsList
					m.customRules = config.GetCustomExclusions()
					m.listCursor = 0
					return m, nil
				}

			case stateInteractiveList:
				m.state = stateDeleting
				return m, executeDeletionCmd(m.items)

			case stateResultSummary:
				m.state = stateMenu
				m.action = ""
				m.items = nil
			}

		case "esc":
			if m.state == stateInteractiveList || m.state == stateResultSummary {
				m.state = stateMenu
				m.action = ""
				m.items = nil
			}
		}
	}
	return m, nil
}

func (m mainModel) View() string {
	var b strings.Builder

	switch m.state {
	case stateMenu:
		b.WriteString(fmt.Sprintf("\nTarget Directory: %s\n", PathStyle.Render(m.cfg.TargetDir)))
		b.WriteString(SubtleStyle.Render(strings.Repeat("-", 60)))
		b.WriteString("\n")
		b.WriteString("Select an operation mode:\n\n")

		for i, choice := range m.choices {
			cursor := "  "
			if m.cursor == i {
				cursor = BrandStyle.Render("> ")
			}
			b.WriteString(fmt.Sprintf("%s%s\n", cursor, choice))
		}

		b.WriteString("\n")
		b.WriteString(SubtleStyle.Render("Navigation: [j/k]  |  Select: [Enter]  |  Quit: [q]"))
		b.WriteString("\n")

	case stateScanning:
		b.WriteString(fmt.Sprintf("\nExecuting module: %s\n", CategoryTitleStyle.Render(m.action)))
		b.WriteString(fmt.Sprintf("Scanning Target:  %s\n", PathStyle.Render(m.cfg.TargetDir)))
		b.WriteString(fmt.Sprintf("\n %s Scanning in progress... Please wait.\n", m.spinner.View()))

	case stateInteractiveList:
		b.WriteString("\n")
		b.WriteString(HeaderStyle.Render(" INTERACTIVE CLEANUP "))
		b.WriteString("\n\n")

		maxVisible := 15
		start := 0
		if m.listCursor >= maxVisible {
			start = m.listCursor - maxVisible + 1
		}
		end := start + maxVisible
		if end > len(m.items) {
			end = len(m.items)
		}

		var totalSelected int64
		for i := start; i < end; i++ {
			item := m.items[i]

			if item.IsHeader {
				b.WriteString("\n")
				b.WriteString(CategoryTitleStyle.Render(item.Title))
				b.WriteString("\n")
				continue
			}

			if item.Selected {
				totalSelected += item.Size
			}

			cursor := "   "
			if m.listCursor == i {
				cursor = BrandStyle.Render(" > ")
			}

			checkbox := "[ ]"
			if item.Selected {
				checkbox = BrandStyle.Render("[x]")
			}

			sizeStr := FormatSize(float64(item.Size) / (1024 * 1024))
			pathStr := PathStyle.Render(item.Path)

			b.WriteString(fmt.Sprintf("%s%s %s  %s\n", cursor, checkbox, sizeStr, pathStr))
		}

		b.WriteString(SubtleStyle.Render(strings.Repeat("-", 60)))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("MARKED FOR DELETION: %s\n", FormatSize(float64(totalSelected)/(1024*1024))))
		b.WriteString(SubtleStyle.Render("Toggle: [Space]  |  Execute Deletion: [Enter]  |  Cancel: [Esc]"))
		b.WriteString("\n")

	case stateDeleting:
		b.WriteString("\n")
		b.WriteString(HeaderStyle.Render(" PURGING DATA "))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("\n %s Destroying selected files... Please wait.\n", m.spinner.View()))

	case stateResultSummary:
		b.WriteString("\n")
		b.WriteString(HeaderStyle.Render(" OPERATION COMPLETE "))
		b.WriteString("\n\n")

		if m.freedSpace > 0 {
			b.WriteString(fmt.Sprintf("Successfully recovered %s of storage.\n", FormatSize(float64(m.freedSpace)/(1024*1024))))
		} else {
			b.WriteString("No files were deleted.\n")
		}

		b.WriteString("\n")
		b.WriteString(SubtleStyle.Render("Return to Menu: [Enter]  |  Quit: [q]"))
		b.WriteString("\n")

	case stateExclusionsList:
		b.WriteString("\n")
		b.WriteString(HeaderStyle.Render(" EXCLUSION RULES "))
		b.WriteString("\n\n")

		sysCount := len(config.GetSystemExclusions())
		b.WriteString(CategoryTitleStyle.Render("System Rules (Protected):"))
		b.WriteString("\n  ")
		b.WriteString(SubtleStyle.Render(fmt.Sprintf("%d default OS paths protected", sysCount)))
		b.WriteString("\n\n")

		b.WriteString(CategoryTitleStyle.Render("Custom Rules:"))
		b.WriteString("\n")

		if len(m.customRules) == 0 {
			b.WriteString("  ")
			b.WriteString(SubtleStyle.Render("No custom rules defined."))
			b.WriteString("\n")
		} else {
			for i, rule := range m.customRules {
				cursor := "  "
				if m.listCursor == i {
					cursor = BrandStyle.Render("> ")
				}
				b.WriteString(fmt.Sprintf("%s%s\n", cursor, rule))
			}
		}

		b.WriteString("\n")
		b.WriteString(SubtleStyle.Render(strings.Repeat("-", 60)))
		b.WriteString("\n")
		b.WriteString(SubtleStyle.Render("Add New: [a]  |  Delete Selected: [x]  |  Back: [Esc]"))
		b.WriteString("\n")

	case stateExclusionsAdd:
		b.WriteString("\n")
		b.WriteString(HeaderStyle.Render(" ADD EXCLUSION PATH "))
		b.WriteString("\n\n")

		b.WriteString(m.textInput.View())

		b.WriteString("\n\n")
		b.WriteString(SubtleStyle.Render(strings.Repeat("-", 60)))
		b.WriteString("\n")
		b.WriteString(SubtleStyle.Render("Save: [Enter]  |  Cancel: [Esc]"))
		b.WriteString("\n")
	}

	return b.String()
}
