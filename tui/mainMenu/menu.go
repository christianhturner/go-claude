package mainmenu

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/tui/constants"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type MenuItem struct {
	title       string
	description string
}

func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.description }
func (i MenuItem) FilterValue() string { return i.title }

type MainMenuModel struct {
	list list.Model
}

func InitialMainMenu() MainMenuModel {
	menuItems := []list.Item{
		MenuItem{title: "Conversations", description: "Create and open conversations with Claude"},
		MenuItem{title: "Import/Export", description: "Import or Export existing go-claude data"},
		MenuItem{title: "Configure", description: "Configure your go-claude experience"},
	}
	mainMenu := list.New(menuItems, list.NewDefaultDelegate(), 0, 0)
	mainMenu.Title = "Go-Claude"
	mainMenu.Styles.Title = titleStyle

	return MainMenuModel{
		list: mainMenu,
	}
}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(MenuItem)
			if ok && i.title == "Conversations" {
				logger.Info("Selecting Converation")
				return m, selectMenuItem(constants.ConversationMenuView)
			}
		}
		if msg.String() == "ctr+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MainMenuModel) View() string {
	return appStyle.Render(m.list.View())
}
