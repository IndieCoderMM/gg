package platformer

import (
	"time"

	g "github.com/Kaamkiya/gg/internal/app/platformer/game"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	FPS = 8
)

type model struct {
	game g.Game
}

func InitGame() *tea.Program {
	game := g.InitGameState()
	m := model{game}

	return tea.NewProgram(m)
}

type tickMsg time.Time

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second/FPS, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.game.Jump()
		case "m":
			m.game.ToggleMode()
		}
	case tickMsg:
		m.game.Update()
		return m, tea.Tick(time.Second/FPS, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	return m, nil
}

func (m model) View() string {
	return m.game.Draw()
}

func Run() {
	game := InitGame()

	if _, err := game.Run(); err != nil {
		panic(err)
	}
}
