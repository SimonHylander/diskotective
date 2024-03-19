package directory

import (
	"path/filepath"
)

type Directory struct {
	Cwd    string
	Items  []Item
	Cursor int
	Sizes  map[string]int64
}

func (d *Directory) NavigateForward() {
	item := d.Items[d.Cursor]

	if item.Type == ItemTypeFile {
		return
	}

	d.Cwd = item.Path
	d.Items = nil
	d.Sizes = nil
}

func (d *Directory) NavigateBackward() {
	d.Cwd = filepath.Dir(d.Cwd)
	d.Items = nil
	d.Sizes = nil
}

func (d *Directory) CursorUp() {
	d.Cursor--

	// Prevent the cursor from going into negative
	if d.Cursor < 0 {
		d.Cursor = 0
		return
	}
}

func (d *Directory) CursorDown() {
	d.Cursor++

	// Prevent the cursor from going out of bounds
	if d.Cursor > len(d.Items)-1 {
		d.Cursor = len(d.Items) - 1
	}
}
