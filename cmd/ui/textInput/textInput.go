package textInput

import (
	"fmt"
	"robbykansas/another-novel-scraper/cmd/novel"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
	// errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8700")).Bold(true).Padding(0, 0, 0)
)

type Output struct {
	Output string
}

func (o *Output) Update(val string) {
	o.Output = val
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	err       error
	output    *Output
	header    string
}

func InitialModel(output *Output, header string, novel *novel.Novel) model {
	ti := textinput.New()
	ti.Placeholder = "Let This Grieving Soul Retire - Woe Is the Weakling Who Leads the Strongest Party"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 150

	return model{
		textInput: ti,
		err:       nil,
		output:    output,
		header:    titleStyle.Render(header),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if len(m.textInput.Value()) > 1 {
				m.output.Update(m.textInput.Value())
				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.header,
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
