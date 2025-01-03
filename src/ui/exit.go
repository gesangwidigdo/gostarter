package ui

// A simple program that counts down from 5 and then exits.

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type exitModel int

func initialExit() exitModel {
	return 1
}

func (m exitModel) Init() tea.Cmd {
	return tick
}

func (m exitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		m--
		if m == 0 {
			return m, tea.Batch(tick, func() tea.Msg {
				time.Sleep(2 * time.Second)
				return tea.Quit
			})
		}
		return m, tick
	}
	return m, nil
}

func (m exitModel) View() string {
	return fmt.Sprint("Exiting...\n", m)
}

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}
