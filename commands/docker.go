package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
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

var customTable table.Model

func DockerPs() (*string, error) {
	out, err := exec.Command("docker", "ps", "--format", "{{json .}}").Output()
	if err != nil {

		return nil, err
	}
	result := string(out)

	return &result, nil
}

func DockerImages() (*string, error) {
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
		{Title: "Repositoy", Width: 30},
		{Title: "Tag", Width: 30},
		{Title: "ID", Width: 30},
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

func DockerPull() (*string, error) {

	stringToResponse := "docker pull"
	return &stringToResponse, nil
}
