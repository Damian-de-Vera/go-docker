package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DockerImage struct {
	Containers   string `json:"Containers"`
	CreatedAt    string `json:"CreatedAt"`
	CreatedSince string `json:"CreatedSince"`
	Digest       string `json:"Digest"`
	ID           string `json:"ID"`
	Repository   string `json:"Repository"`
	SharedSize   string `json:"SharedSize"`
	Size         string `json:"Size"`
	Tag          string `json:"Tag"`
	UniqueSize   string `json:"UniqueSize"`
	VirtualSize  string `json:"VirtualSize"`
}

var (
	customTable        table.Model
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
		0: dockerPs,
		1: dockerImages,
		2: dockerPull,
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

func dockerPs() (*string, error) {
	out, err := exec.Command("docker", "ps", "--format", "{{json .}}").Output()
	if err != nil {

		return nil, err
	}
	result := string(out)

	return &result, nil
}

func dockerImages() (*string, error) {
	out, err := exec.Command("docker", "images", "--format", "{{json .}}").Output()
	if err != nil {
		os.Exit(1)
		return nil, err
	}

	result := string(out)
	linesArray := strings.Split(result, "\n")

	// Creamos un slice de DockerImage
	var images []DockerImage
	for _, line := range linesArray {
		if line == "" {
			continue
		}

		var image DockerImage
		err := json.Unmarshal([]byte(line), &image)
		if err != nil {
			fmt.Println("Error al parsear JSON:", err)
			continue
		}
		images = append(images, image)
	}
	columns := []table.Column{
		{Title: "ID", Width: 30},
		{Title: "Repository", Width: 30},
		{Title: "Tag", Width: 30},
	}
	var rows []table.Row

	// Mostramos las imágenes
	for _, image := range images {
		row := table.Row{image.Repository, image.Tag, image.ID}

		rows = append(rows, row)

	}
	customTable = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(7),
	)
	tableString := customTable.View()
	return &tableString, nil
}

func dockerPull() (*string, error) {

	stringToResponse := "docker pull"
	return &stringToResponse, nil
}
