package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var SelectedDBMS string

type dbItem string

func (d dbItem) FilterValue() string { return "" }

type dbItemDelegate struct{}

func (d dbItemDelegate) Height() int                             { return 1 }
func (d dbItemDelegate) Spacing() int                            { return 0 }
func (d dbItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d dbItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(dbItem)
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

type dbModel struct {
	list     list.Model
	choice   string
	quitting bool
	done     bool
}

func (m dbModel) Init() tea.Cmd {
	return nil
}

func (m dbModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(dbItem)
			if ok {
				m.choice = string(i)
				SelectedDBMS = m.choice
				m.done = true
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m dbModel) View() string {
	if m.quitting {
		return quitTextStyle.Render("Quitting...")
	}
	return "\n" + m.list.View()
}

func initialSelectDB() dbModel {
	items := []list.Item{
		dbItem("MySQL"),
		dbItem("PostgreSQL"),
		dbItem("SQLite"),
		dbItem("MongoDB"),
	}

	const defaultWidth = 50

	dbList := list.New(items, dbItemDelegate{}, defaultWidth, 10)
	dbList.Title = "Select Database"
	dbList.SetShowStatusBar(false)
	dbList.SetFilteringEnabled(true)
	dbList.Styles.Title = titleStyle
	dbList.Styles.PaginationStyle = paginationStyle
	dbList.Styles.HelpStyle = helpStyle

	return dbModel{list: dbList}
}
