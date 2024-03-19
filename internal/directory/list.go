package directory

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/simonhylander/diskotective/internal/color"
	"math"
)

func RenderList(d Directory, width int) string {
	items := d.Items
	cursor := d.Cursor

	view := ""

	for i, item := range items {
		textStyle := lipgloss.NewStyle().Padding(0, 1)

		isSelected := cursor == i

		if isSelected {
			textStyle = textStyle.
				Background(color.Indigo).
				Foreground(color.White)
		}

		size := item.Size
		readableSize := bytesToHumanReadable(size, isSelected)

		name := fmt.Sprintf("%s", item.Name)

		if item.Type == ItemTypeDirectory {
			name += "/"
		}

		name = textStyle.Render(name)

		if isSelected {
			selectedItem := lipgloss.NewStyle().
				Background(color.Indigo).
				Width(width)

			view += selectedItem.Render(fmt.Sprintf("%s%s",
				textStyle.Render(readableSize),
				name,
			))

			view += "\n"
		} else {
			view += textStyle.Render(readableSize)
			view += name
			view += "\n"
		}
	}

	return lipgloss.NewStyle().Render(view)
}

func bytesToHumanReadable(bytes int64, isSelected bool) string {
	suffixes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	bg := lipgloss.NewStyle()

	sizeStyle := lipgloss.NewStyle().
		Width(8).
		AlignHorizontal(lipgloss.Right).
		Bold(true).
		Foreground(color.Yellow)

	suffixStyle := lipgloss.NewStyle().
		Width(2).
		Bold(true)

	if isSelected {
		bg = bg.Background(color.Indigo)

		sizeStyle = sizeStyle.
			Background(color.Indigo)

		suffixStyle = suffixStyle.
			Background(color.Indigo).
			Foreground(color.White)
	}

	exp := int(math.Log(float64(bytes)) / math.Log(1024))
	size := float64(bytes) / math.Pow(1024, float64(exp))
	var suffix string

	if bytes < 1024 {
		size = float64(bytes)
		suffix = suffixStyle.Render(suffixes[0])
	} else {
		suffix = suffixes[exp]
	}

	sizeText := sizeStyle.Render(fmt.Sprintf("%.1f", size))

	return bg.Render(
		sizeText +
			bg.Render(" ") +
			suffixStyle.Render(suffix),
	)
}
