package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/simonhylander/diskotective/internal/color"
	"github.com/simonhylander/diskotective/internal/directory"
	"github.com/simonhylander/diskotective/internal/scan"
	"github.com/simonhylander/diskotective/internal/view"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

var (
	titleStyle = lipgloss.NewStyle().
		Foreground(color.Yellow)

	scannedFileStyle = lipgloss.NewStyle().
		Foreground(color.Indigo)

	spinnerStyle = lipgloss.NewStyle().
		Foreground(color.Indigo)
)

type model struct {
	keys          *KeyMap
	help          help.Model
	spinner       spinner.Model
	scanner       chan scan.ScanEvent
	directory     directory.Directory
	width         int
	title         string
	scannedFile   string
	deletedAmount int64
	confirmDelete bool
	showSpinner   bool
	quitting      bool
}

func newModel() model {
	var (
		keyMap = NewKeyMap()
	)

	var cwd string

	if len(os.Args) == 2 {
		cwd = os.Args[1]
		if _, err := os.Stat(cwd); os.IsNotExist(err) {
			panic("The provided directory does not exist")
		}
	} else {
		cwd, _ = os.Getwd()
	}

	scanner := make(chan scan.ScanEvent)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	d := directory.Directory{
		Cwd:   cwd,
		Items: nil,
		Sizes: nil,
	}

	m := model{
		keys:      keyMap,
		help:      help.New(),
		scanner:   scanner,
		spinner:   s,
		directory: d,
	}

	return m
}

func listenForFiles(m model) tea.Cmd {
	return func() tea.Msg {
		if m.directory.Cwd == "" {
			return nil
		}

		directoryEntries, err := scan.ListFiles(m.directory.Cwd)

		if err != nil {
			panic(err)
		}

		if m.directory.Sizes == nil {
			m.directory.Sizes = make(map[string]int64)
		}

		var items = make([]directory.Item, len(directoryEntries))

		if m.directory.Items == nil {
			m.directory.Items = items
		}

		for i, directoryEntry := range directoryEntries {
			info, err := directoryEntry.Info()

			if err != nil {
				continue
			}

			itemType := directory.ItemTypeFile

			if info.IsDir() {
				itemType = directory.ItemTypeDirectory
				m.directory.Sizes[filepath.Join(m.directory.Cwd, info.Name())] = 0
			} else {
				m.directory.Sizes[filepath.Join(m.directory.Cwd, info.Name())] = info.Size()
			}

			item := directory.Item{
				Name: info.Name(),
				Path: filepath.Join(m.directory.Cwd, info.Name()),
				Type: itemType,
				Size: info.Size(),
			}

			items[i] = item
		}

		m.scanner <- scan.ScanEvent{
			Type:  scan.InitializedDiskScanEventType,
			Items: items,
		}

		var mutex = &sync.Mutex{}
		var wg sync.WaitGroup
		wg.Add(len(items))

		for _, item := range items {
			go func(item directory.Item) {
				defer wg.Done()

				_ = filepath.Walk(item.Path, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						panic(err)
					}

					// Skip directories
					if !info.IsDir() {
						mutex.Lock()

						m.directory.Sizes[item.Path] += info.Size()
						item.Size += info.Size()
						mutex.Unlock()

						m.scanner <- scan.ScanEvent{
							Type: scan.FileDiskScanEvenType,
							ScannedFile: scan.ScannedFile{
								Path: path,
								Size: info.Size(),
							},
						}
					}

					return nil
				})
			}(item)
		}

		wg.Wait()

		for i, item := range items {
			items[i].Size = m.directory.Sizes[item.Path]
		}

		sort.SliceStable(items, func(a, b int) bool {
			return items[a].Size > items[b].Size
		})

		m.scanner <- scan.ScanEvent{
			Type:  scan.CompletedDiskScanEventType,
			Items: items,
		}

		return nil
	}
}

func waitForFileScan(scanner chan scan.ScanEvent) tea.Cmd {
	return func() tea.Msg {
		return scan.ScanEvent(<-scanner)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		listenForFiles(m),
		waitForFileScan(m.scanner),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.CursorUp):
			if m.confirmDelete {
				return m, nil
			}

			m.directory.CursorUp()
			return m, nil
		case key.Matches(msg, m.keys.CursorDown):
			if m.confirmDelete {
				return m, nil
			}

			m.directory.CursorDown()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			if m.confirmDelete {
				return m, nil
			}

			m.directory.NavigateForward()

			return m, tea.Batch(listenForFiles(m), waitForFileScan(m.scanner))
		case key.Matches(msg, m.keys.Backspace):
			if m.confirmDelete {
				return m, nil
			}

			m.directory.NavigateBackward()

			return m, tea.Batch(listenForFiles(m), waitForFileScan(m.scanner))
		case key.Matches(msg, m.keys.Delete):
			if m.confirmDelete {
				return m, nil
			}

			m.confirmDelete = true
			return m, nil

		case key.Matches(msg, m.keys.Yes):
			if m.confirmDelete {
				selected := m.directory.Items[m.directory.Cursor]
				err := os.RemoveAll(selected.Path)
				if err != nil {
					log.Println(err)
				}

				m.directory.Items = append(m.directory.Items[:m.directory.Cursor], m.directory.Items[m.directory.Cursor+1:]...)
				delete(m.directory.Sizes, selected.Path)

				m.deletedAmount += selected.Size
				m.confirmDelete = false

				return m, nil
			}

			return m, nil

		case key.Matches(msg, m.keys.No):
			m.confirmDelete = false
			return m, nil

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}

	case scan.ScanEvent:
		switch msg.Type {
		case scan.InitializedDiskScanEventType:
			m.directory.Items = msg.Items
			m.directory.Cursor = 0
			m.showSpinner = true

			return m, waitForFileScan(m.scanner)
		
		case scan.FileDiskScanEvenType:
			cwd := titleStyle.Render(m.directory.Cwd)
			scannedFile := scannedFileStyle.Render(msg.ScannedFile.Path)

			m.title = cwd
			m.scannedFile = scannedFile

			return m, waitForFileScan(m.scanner)

		case scan.CompletedDiskScanEventType:
			m.title = titleStyle.Render(m.directory.Cwd)
			m.scannedFile = ""
			m.showSpinner = false
			m.directory.Items = msg.Items

			return m, nil
		}
	}

	sa, cmd := m.spinner.Update(msg)
	m.spinner = sa
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {

	if m.quitting && m.deletedAmount > 0 {
		return fmt.Sprintf("Deleted %d bytes", m.deletedAmount)
	}

	output := view.Taskbar(m.width)
	output += m.title + " "

	if m.confirmDelete {
		selected := m.directory.Items[m.directory.Cursor]
		output += view.ConfirmDelete(selected.Name)
	} else {
		if m.showSpinner {
			output += m.spinner.View()
		}

		output += m.scannedFile
	}

	output += "\n"

	output += directory.RenderList(m.directory, m.width)

	helpView := m.help.View(m.keys)
	output += fmt.Sprintf("%s\n", helpView)

	return output
}

func Execute() {
	// os.Setenv("DEBUG", "1")

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
