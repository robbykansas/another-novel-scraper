package listInput

import (
	"fmt"
	"robbykansas/another-novel-scraper/cmd/flags"
	"robbykansas/another-novel-scraper/cmd/novel"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Change this
var (
	focusedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	titleStyle            = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
	selectedItemStyle     = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
	selectedItemDescStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170"))
	descriptionStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#40BDA3"))
)

type sessionState uint

const (
	titleView sessionState = iota
	webView   sessionState = iota
)

// A Selection represents a choice made in a multiInput step
type Selection struct {
	Choice string
}

// Update changes the value of a Selection's Choice
func (s *Selection) Update(value string) {
	s.Choice = value
}

// A multiInput.model contains the data for the multiInput step.
//
// It has the required methods that make it a bubbletea.Model
type model struct {
	cursor    int
	choices   []flags.NovelData
	selected  map[int]struct{}
	choice    *Selection
	header    string
	exit      *bool
	paginator paginator.Model
	state     sessionState
}

func (m model) Init() tea.Cmd {
	return nil
}

var limitPagination = 3

// InitialModelMulti initializes a multiInput step with
// the given data
func InitialModelMulti(choices []flags.NovelData, selection *Selection, header string, novel *novel.Novel) model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = limitPagination
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(choices))

	return model{
		choices:   choices,
		selected:  make(map[int]struct{}),
		choice:    selection,
		header:    titleStyle.Render(header),
		exit:      &novel.Exit,
		paginator: p,
		state:     titleView,
	}
}

// Update is called when "things happen", it checks for
// important keystrokes to signal when to quit, change selection,
// and confirm the selection.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			*m.exit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < limitPagination-1 {
				m.cursor++
			}
		case "enter", " ":
			if len(m.selected) == 1 {
				m.selected = make(map[int]struct{})
			}
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "y":
			if len(m.selected) == 1 {
				for selectedKey := range m.selected {
					m.choice.Update(m.choices[selectedKey].WebName)
					m.cursor = selectedKey
				}
				return m, tea.Quit
			}
		}
	}
	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

// View is called to draw the multiInput step
func (m model) View() string {
	s := m.header + "\n\n"
	// var b strings.Builder
	start, end := m.paginator.GetSliceBounds(len(m.choices))

	if m.state == titleView {
		for i, choice := range m.choices[start:end] {
			cursor := " "
			if m.cursor == i {
				cursor = focusedStyle.Render(">")
				choice.Title = selectedItemDescStyle.Render(choice.Title)
				choice.AvailableChapter = selectedItemDescStyle.Render(choice.AvailableChapter)
			}

			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = focusedStyle.Render("x")
			}

			title := descriptionStyle.Render(choice.Title)
			availableChapter := descriptionStyle.Render(choice.AvailableChapter)

			s += fmt.Sprintf("%s [%s] %s\n%s\n\n", cursor, checked, title, availableChapter)
		}
	} else {
		for i, choice := range m.choices[start:end] {
			cursor := " "
			if m.cursor == i {
				cursor = focusedStyle.Render(">")
				choice.WebName = selectedItemStyle.Render(choice.WebName)
				choice.Title = selectedItemDescStyle.Render(choice.Title)
				choice.AvailableChapter = selectedItemDescStyle.Render(choice.AvailableChapter)
			}

			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = focusedStyle.Render("x")
			}

			webName := focusedStyle.Render(choice.WebName)
			title := descriptionStyle.Render(choice.Title)
			availableChapter := descriptionStyle.Render(choice.AvailableChapter)

			s += fmt.Sprintf("%s [%s] %s\n%s\n%s\n\n", cursor, checked, webName, title, availableChapter)
		}
	}

	s += fmt.Sprintf("Press %s to confirm choice.\n\n", focusedStyle.Render("y"))
	return s
}
