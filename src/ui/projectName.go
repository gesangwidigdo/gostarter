package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	InsertedProjectName string
	InsertedModuleURL   string
)

var (
	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type projectNameModel struct {
	focusIndex int
	inputs     []textinput.Model
	done       bool
}

func initialProjectName() projectNameModel {
	m := projectNameModel{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 50

		switch i {
		case 0:
			fmt.Println("Project Name: ")
			t.Placeholder = "Project Name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle

		case 1:
			fmt.Println("Module URL: ")
			t.Placeholder = "example.com/username/project"
			t.CharLimit = 100
		}

		m.inputs[i] = t
	}

	return m
}

func (pn projectNameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (pn projectNameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return pn, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && pn.focusIndex == len(pn.inputs) {
				InsertedProjectName = pn.inputs[0].Value()
				InsertedModuleURL = pn.inputs[1].Value()
				pn.done = true
				return pn, nil
			}

			if s == "up" || s == "shift+tab" {
				pn.focusIndex--
			} else {
				pn.focusIndex++
			}

			if pn.focusIndex > len(pn.inputs) {
				pn.focusIndex = 0
			} else if pn.focusIndex < 0 {
				pn.focusIndex = len(pn.inputs)
			}

			cmds := make([]tea.Cmd, len(pn.inputs))
			for i := 0; i <= len(pn.inputs)-1; i++ {
				if i == pn.focusIndex {
					cmds[i] = pn.inputs[i].Focus()
					pn.inputs[i].PromptStyle = focusedStyle
					pn.inputs[i].TextStyle = focusedStyle
					continue
				}

				// remove focused state
				pn.inputs[i].Blur()
				pn.inputs[i].PromptStyle = noStyle
				pn.inputs[i].TextStyle = noStyle
			}

			return pn, tea.Batch(cmds...)
		}
	}
	cmd := pn.updateInputs(msg)

	return pn, cmd
}

func (pn *projectNameModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(pn.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range pn.inputs {
		pn.inputs[i], cmds[i] = pn.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (pn projectNameModel) View() string {
	var b strings.Builder

	for i := range pn.inputs {
		b.WriteString(pn.inputs[i].View())
		if i < len(pn.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if pn.focusIndex == len(pn.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}

func ProjectName() {
	p := tea.NewProgram(initialProjectName())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
