package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gesangwidigdo/gostarter/src/program"
)

type Page int

const (
	PageProjectName Page = iota
	PageSelectFramework
	PageSelectDB
)

type appModel struct {
	CurrentPage Page
	ProjectName projectNameModel
	Framework   frameworkModel
	DB          dbModel
	Quitting    bool
}

func initialAppModel() appModel {
	return appModel{
		CurrentPage: PageProjectName,
		ProjectName: initialProjectName(),
		Framework:   initialSelectFramework(),
		DB:          initialSelectDB(),
		Quitting:    false,
	}
}

func (m appModel) Init() tea.Cmd {
	switch m.CurrentPage {
	case PageProjectName:
		return m.ProjectName.Init()

	case PageSelectFramework:
		return m.Framework.Init()

	case PageSelectDB:
		return m.DB.Init()

	default:
		return nil
	}
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.CurrentPage {
	case PageProjectName:
		newModel, cmd := m.ProjectName.Update(msg)
		m.ProjectName = newModel.(projectNameModel)
		if m.ProjectName.done {
			m.CurrentPage = PageSelectFramework
		}
		return m, cmd

	case PageSelectFramework:
		newModel, cmd := m.Framework.Update(msg)
		m.Framework = newModel.(frameworkModel)
		if m.Framework.done {
			m.CurrentPage = PageSelectDB
		}
		return m, cmd

	case PageSelectDB:
		newModel, cmd := m.DB.Update(msg)
		m.DB = newModel.(dbModel)
		if m.DB.done {
			m.Quitting = true
			return m, tea.Quit
		}
		return m, cmd
	}

	return m, nil
}

func (m appModel) View() string {
	switch m.CurrentPage {
	case PageProjectName:
		return m.ProjectName.View()

	case PageSelectFramework:
		return m.Framework.View()

	case PageSelectDB:
		return m.DB.View()
	}

	if m.Quitting {
		return "\nQuitting...\n"
	}
	return ""
}

func RunApp() {
	p := tea.NewProgram(initialAppModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}

	// After the user has selected all the options, we can now initialize the project
	program.ProjectInitialization(InsertedProjectName, InsertedModuleURL)
	program.InstallDependencies(selectedFramework, selectedDBMS)
}
