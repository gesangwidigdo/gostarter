package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var SelectedFramework string

type item string

func (i item) FilterValue() string { return "" }

type frameworkItemDelegate struct{}

func (d frameworkItemDelegate) Height() int                             { return 1 }
func (d frameworkItemDelegate) Spacing() int                            { return 0 }
func (d frameworkItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d frameworkItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type frameworkModel struct {
	list     list.Model
	choice   string
	quitting bool
	done     bool
}

func (m frameworkModel) Init() tea.Cmd {
	return nil
}

func (m frameworkModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
				SelectedFramework = m.choice
				m.done = true
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m frameworkModel) View() string {
	if m.quitting {
		return quitTextStyle.Render("Quitting...")
	}
	return "\n" + m.list.View()
}

func initialSelectFramework() frameworkModel {
	items := []list.Item{
		item("Gin"),
		item("Echo"),
		item("Iris"),
	}

	const defaultWidth = 80

	l := list.New(items, frameworkItemDelegate{}, defaultWidth, 10)
	l.Title = titleStyle.Render("Select a framework (template only available for gin)")
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return frameworkModel{list: l}
}
