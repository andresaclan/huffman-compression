package main

import (
	"errors"
	"fmt"
	"huffman-compression/internal/compress"
	"huffman-compression/internal/decompress"
	"huffman-compression/internal/file"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	filepicker   filepicker.Model
	choices      []string         // compress or decompress
	cursor       int              // which item our cursor is pointing at
	selected     map[int]struct{} // which item is selected
	selectedFile string
	processing   bool
	quitting     bool
	err          error
}

var action string
var filePath string

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
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
		case "enter", " ":
			for i := range 2 {
				delete(m.selected, i)
			}
			m.selected[m.cursor] = struct{}{}
		case "y":
			// begin compression or decompression
			m.processing = true
			if m.cursor == 0 {
				action = "compress"
			} else {
				action = "decompress"
			}
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}
	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
		filePath = path
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd

}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	if m.processing {
		return "doing work..."
	}
	if m.selectedFile != "" {
		s := fmt.Sprintln("\n  You selected: " + m.filepicker.Styles.Selected.Render(m.selectedFile) + "\n")
		s += "Would you like to compress or decompress this file?\n\n"

		// render the choices
		for i, choice := range m.choices {

			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			// Is this choice selected?
			checked := " " // not selected
			if _, ok := m.selected[i]; ok {
				checked = "x" // selected!
			}

			// Render the row
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}

		s += "\nPress y to confirm choice\n"

		return s
	}

	var s strings.Builder
	s.WriteString("\n ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file to Compress or Decompress:")
	} else {
		s.WriteString("Selected File: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}

func main() {
	// bubble tea
	fp := filepicker.New()
	fp.AllowedTypes = []string{".txt"}
	fp.CurrentDirectory, _ = os.UserHomeDir()

	m := model{
		filepicker: fp,
		choices:    []string{"compress", "decompress"},
		selected:   make(map[int]struct{}, 2),
	}
	_, err := tea.NewProgram(&m).Run()
	if err != nil {
		fmt.Printf("There has been an error: %v", err)
		os.Exit(1)
	}

	// compression/decompression
	outputFilePath := "compressed.txt"

	if action == "" {
		fmt.Println("Invalid choice")
		return
	} else if action == "compress" {
		// compress logic
		data, err := file.ReadFile(filePath)
		if err != nil { // handle error
			fmt.Println("Error reading file:", err)
			return
		}
		compressedData, err := compress.Compress(data)
		if err != nil { // handle error
			fmt.Println("Error compressing data:", err)
			return
		}
		err = file.WriteFile(outputFilePath, compressedData)
		if err != nil { // handle error
			fmt.Println("Error writing file:", err)
			return
		}
		fmt.Println("File compressed successfully!")
		return
	} else if action == "decompress" {
		// decompress logic
		outputFilePath = "decompressed.txt"
		data, err := file.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		decompressedData, err := decompress.Decompress(data)
		if err != nil {
			fmt.Println("Error decompressing data:", err)
			return
		}
		err = file.WriteFile(outputFilePath, decompressedData)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}
		fmt.Println("File decompressed successfully!")
		return
	} else {
		fmt.Println("Invalid function")
		return
	}
}
