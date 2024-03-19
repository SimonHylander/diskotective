package view

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/simonhylander/diskotective/internal/color"
	"github.com/simonhylander/diskotective/internal/version"
)

func Taskbar(width int) string {
	taskBarStyle := lipgloss.NewStyle().
		Background(color.Indigo).
		Foreground(color.White)

	text := fmt.Sprintf("Diskotective: %s", version.Version)

	return fmt.Sprintf("%s\n", taskBarStyle.Width(width).Render(text))
}

func ConfirmDelete(name string) string {
	prompt := lipgloss.NewStyle().Foreground(color.Red).Render("Are you sure you want to delete")
	name = lipgloss.NewStyle().Foreground(color.Indigo).Render(name)
	return fmt.Sprintf("%s \"%s\"? (y/n)", prompt, name)
}
