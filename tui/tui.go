package tui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/christianhturner/go-claude/config"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/tui/constants"
	conversationmenu "github.com/christianhturner/go-claude/tui/conversationMenu"
	mainmenu "github.com/christianhturner/go-claude/tui/mainMenu"
	"github.com/spf13/viper"
)

var (
	dataDir = viper.GetString(config.DataDirKey)
	cfgFile = viper.GetString(config.CfgFile)
	p       *tea.Program
)

type StateManagerModel struct {
	state                constants.SessionState
	mainMenu             tea.Model
	conversationMenu     tea.Model
	activeConversationId uint
	WindowSize           tea.WindowSizeMsg
}

func StartTea() error {
	m := New()
	p = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.LogError(err, "Alas, there has been an error")
		os.Exit(1)
	}
	return nil
}

func New() StateManagerModel {
	return StateManagerModel{
		state:            constants.MainMenuView,
		mainMenu:         mainmenu.InitialMainMenu(),
		conversationMenu: conversationmenu.InitialConversationMenu(),
	}
}

func (m StateManagerModel) Init() tea.Cmd {
	return nil
}

func (m StateManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.WindowSize = msg // Pass with a view to set the window size ??
	case conversationmenu.BackMsg:
		m.state = constants.MainMenuView
	}

	switch m.state {
	case constants.MainMenuView:
		newMainMenuModel, newCmd := m.mainMenu.Update(msg)
		if newMainMenu, ok := newMainMenuModel.(mainmenu.MainMenuModel); ok {
			m.mainMenu = newMainMenu
		}
		cmd = newCmd
	case constants.ConversationMenuView:
		newConversationMenuModel, newCmd := m.conversationMenu.Update(msg)
		if newConversationMenu, ok := newConversationMenuModel.(conversationmenu.ConversationMenuModel); ok {
			m.conversationMenu = newConversationMenu
		}
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m StateManagerModel) View() string {
	switch m.state {
	case constants.ConversationMenuView:
		return m.conversationMenu.View()
	default:
		return m.mainMenu.View()
	}
}

// https://github.com/bashbunni/pjs/blob/1469435d4fc09561e872b92c018be0f55b74f4bf/tui/tui.go
//
//
// func StartTea() error {
// 	configPath := filepath.Join(dataDir, cfgFile, ".json")
// 	if f, err := tea.LogToFile(configPath, "DEBUG"); err != nil {
// 		fmt.Println("Couldn't open file for logging:", err)
// 		os.Exit(1)
// 	} else {
// 		defer func() {
// 			err = f.Close()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}()
// 	}
// 	m := InitialMainMenu()
// 	c := conversationMenuModel{}
// 	logger.Info("conversation model established")
// 	p := tea.NewProgram(viewManagerModel{
// 		currentView: "main",
// 		mainMenu:    m,
// 		convList:    c,
// 	}, tea.WithAltScreen())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Println("Alas, there has been an error", err)
// 		os.Exit(1)
// 	}
// 	return nil
// }
//
// type viewManagerModel struct {
// 	currentView string
// 	mainMenu    mainMenuModel
// 	convList    conversationMenuModel
// }
//
// func (m viewManagerModel) Init() tea.Cmd {
// 	return nil
// }
//
// func (m viewManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmd tea.Cmd
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		if msg.String() == "ctrl+c" {
// 			return m, tea.Quit
// 		}
// 	case switchViewMsg:
// 		m.currentView = msg.view
// 		if msg.view == "conversations" {
// 			logger.Info("Initializing conversation")
// 			convModel, cmd := InitialConversationMenu()
// 			if cm, ok := convModel.(conversationMenuModel); ok {
// 				m.convList = cm
// 				logger.Info(fmt.Sprintf("Switch to conversation view with %d items", len(m.convList.list.Items())))
// 			}
// 			return m, cmd
// 		}
// 	case tea.WindowSizeMsg:
// 		if m.currentView == "conversations" {
// 			logger.Info(fmt.Sprintf("Setting conversation menu size to %dx%d", msg.Width-2, msg.Height-2))
// 			m.convList.SetSize(msg.Width-2, msg.Height-2)
// 		}
// 	}
//
// 	switch m.currentView {
// 	case "main":
// 		var updatedModel tea.Model
// 		updatedModel, cmd = m.mainMenu.Update(msg)
// 		if updatedMainMenu, ok := updatedModel.(mainMenuModel); ok {
// 			m.mainMenu = updatedMainMenu
// 		}
// 	case "conversations":
// 		var updatedModel tea.Model
// 		updatedModel, cmd = m.convList.Update(msg)
// 		if updatedConvModel, ok := updatedModel.(conversationMenuModel); ok {
// 			m.convList = updatedConvModel
// 		}
// 	}
// 	return m, cmd
// }
//
// func (m viewManagerModel) View() string {
// 	switch m.currentView {
// 	case "main":
// 		return m.mainMenu.View()
// 	case "conversations":
// 		return m.convList.View()
// 	default:
// 		return "Unknown view"
// 	}
// }
