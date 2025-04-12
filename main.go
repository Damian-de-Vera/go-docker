package main

import (
	"fmt"
	"go-docker-cli/commands"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	leftContainerStyle = lipgloss.NewStyle().
				Width(20).
				Height(20).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("21"))

	rightContainerStyle = lipgloss.NewStyle().
				Width(80).
				Height(20).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("21"))

	commandFuncs = map[int]func() (*string, error){
		0: commands.DockerPs,
		1: commands.DockerImages,
		2: commands.DockerPull,
	}
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	choices               []string // items on the to-do list
	cursor                int      // which to-do list item our cursor is pointing at
	rightContainerContent string
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices:               []string{"Docker PS", "Docker Images", "Docker Pull"},
		rightContainerContent: "",
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.

		case "x":
			command, err := executeCommand(m)
			if err != nil {
				fmt.Println(err)
				return m, tea.Quit
			}
			m.rightContainerContent = *command

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := ""
	leftContainerContent := strings.Builder{}

	// Iterate over our choices
	for i, choice := range m.choices {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// // Is this choice selected?
		// checked := " " // not selected
		// if _, ok := m.selected[i]; ok {
		// 	checked = "x" // selected!
		// }

		// Render the row
		leftContainerContent.WriteString(fmt.Sprintf("%s  %s\n", cursor, choice))
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, leftContainerStyle.Render(leftContainerContent.String()), rightContainerStyle.Render(m.rightContainerContent))

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func executeCommand(m model) (*string, error) {
	fn, ok := commandFuncs[m.cursor]
	if !ok {
		return nil, fmt.Errorf("comando no encontrado para índice %d", m.cursor)
	}

	result, err := fn()
	if err != nil {
		return nil, err
	}

	return result, nil
}
