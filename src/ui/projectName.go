package ui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var InsertedProjectName string

type projectNameModel struct {
	projectName textinput.Model
	err         error
	done        bool
}

func initialProjectName() projectNameModel {
	projectName := textinput.New()
	projectName.Placeholder = "golang-project"
	projectName.Focus()
	projectName.CharLimit = 156
	projectName.Width = 20

	return projectNameModel{
		projectName: projectName,
		err:         nil,
	}
}

func (pn projectNameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (pn projectNameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			pn.done = true
			InsertedProjectName = pn.projectName.Value()
			return pn, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return pn, tea.Quit
		}
	case error:
		pn.err = msg
		return pn, nil
	}

	pn.projectName, cmd = pn.projectName.Update(msg)
	return pn, cmd
}

func (pn projectNameModel) View() string {
	return fmt.Sprintf(
		"Project Name: \n%s\n\n\n%s",
		pn.projectName.View(),
		"(ESC to quit)",
	)
}

func ProjectName() {
	p := tea.NewProgram(initialProjectName())
	model, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	pnModel := model.(projectNameModel)
	fmt.Println("Project Name: ", pnModel.projectName.Value())
}
