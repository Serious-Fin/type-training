package gui

import (
	"fmt"
	"os"

	speech "personal/type-training/sentence_gen"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/logrusorgru/aurora/v4"
)

func RunGui() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	text             string
	currentIndex     int
	textCorrectUntil int
}

func initialModel() model {
	return model{
		text:             speech.GenerateSentences(1),
		currentIndex:     0,
		textCorrectUntil: -1,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key_msg, ok := msg.(tea.KeyMsg)

	// If Update was called not for a key stroke - ignore
	if !ok {
		return m, nil
	}

	value := key_msg.String()

	// Quit if ctrl+c was pressed
	if value == "ctrl+c" {
		return m, tea.Quit
	}

	if value == "enter" && m.textCorrectUntil >= len(m.text)-2 {
		m.text = speech.GenerateSentences(1)
		m.currentIndex = 0
		m.textCorrectUntil = -1
		return m, nil
	}

	if value == "backspace" {
		if m.currentIndex > 0 {
			m.currentIndex--
		}

		if m.currentIndex <= m.textCorrectUntil {
			m.textCorrectUntil--
		}

		return m, nil
	}

	if ([]rune(m.text))[m.currentIndex] == ([]rune(value))[0] && m.textCorrectUntil+1 == m.currentIndex {
		m.textCorrectUntil++
	}

	if m.currentIndex < len(m.text)-1 {
		m.currentIndex++
	}

	return m, nil
}

func (m model) View() string {
	s := ""
	runes := []rune(m.text)

	// Add all correct (green) characters
	if m.textCorrectUntil >= 0 {
		s += aurora.Sprintf(aurora.Green(string(runes[0 : m.textCorrectUntil+1])))
	}

	// Add all incorrect (red) characters
	if m.textCorrectUntil < m.currentIndex-1 {
		s += aurora.Sprintf(aurora.Red(string(runes[m.textCorrectUntil+1 : m.currentIndex])))
	}

	// Add current character (white background)
	s += aurora.Sprintf(aurora.BgWhite(aurora.Black(string(runes[m.currentIndex : m.currentIndex+1]))))

	// Add not typed characters (grey)
	s += aurora.Sprintf(aurora.Gray(10, string(runes[m.currentIndex+1:])))

	// The footer
	s += "\nPress \"ctrl+c\" to quit.\n"

	// Send the UI for rendering
	return s
}
