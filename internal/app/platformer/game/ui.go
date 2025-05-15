package game

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	borderColor = lipgloss.Color("#5f574f")
	bgColor     = lipgloss.Color("#1a1c2c")
	fgColor     = lipgloss.Color("#c2c3c7")
	yellow      = lipgloss.Color("#ffec27")
	red         = lipgloss.Color("#ff004d")
	blue        = lipgloss.Color("#29adff")

	hudContainer = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(fgColor).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(borderColor).
			Padding(0, 1)

	scoreStyle = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(yellow).
			Bold(true)

	livesStyle = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(red).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(fgColor).
			MarginRight(1)

	viewportBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(borderColor).
			Padding(0, 1)

	commandStyle = lipgloss.NewStyle().
			Foreground(blue)

	spaceStyle = lipgloss.NewStyle().
			Background(bgColor)

	playerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFD0")).Bold(true)

	enemyStyles = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#77DD77")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB347")),
	}

	groundStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
)

func (game Game) Draw() string {
	var sb strings.Builder

	sb.WriteString(game.renderHUD())
	sb.WriteString("\n")
	sb.WriteString(game.renderViewport())
	sb.WriteString("\n")
	sb.WriteString(game.renderFooter())

	return sb.String()
}

func (game Game) renderHUD() string {
	score := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("SCORE:"),
		scoreStyle.Render(fmt.Sprintf("%d", game.Score)))

	lives := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("LIVES:"),
		livesStyle.Render(fmt.Sprintf("%d", game.Lives)))

	spacer := game.ViewWidth*2 - len("SCORE: 000") - len("LIVES: 0") + 2

	hudRow := lipgloss.JoinHorizontal(lipgloss.Top, score, spaceStyle.Width(spacer).Render(" "), lives)
	return hudContainer.Render(hudRow)
}

func (game Game) renderViewport() string {
	offsetX, offsetY := game.getOffset()
	var sb strings.Builder

	// Store this grid in state
	grid := make([][]string, game.ViewHeight)
	for y := range grid {
		grid[y] = make([]string, game.ViewWidth)
		for x := range grid[y] {
			grid[y][x] = "." // default empty
		}
	}

	for _, e := range game.Entities {
		x, y := e.Pos.X-offsetX, e.Pos.Y-offsetY
		if x < 0 || x >= game.ViewWidth || y < 0 || y >= game.ViewHeight {
			continue // Out of bounds
		}
		grid[y][x] = string(e.Text)
	}

	for y := range grid {
		for x := range grid[y] {
			// TODO: Render with colors
			sb.WriteString(grid[y][x])

			if y == GROUND_LEVEL {
				sb.WriteRune('▒')
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}

	return viewportBorder.Render(sb.String())
}

func (game Game) renderFooter() string {
	if game.Status == "gameover" {
		return commandStyle.Render("✖ GAME OVER ✖   ⟶ [R] to Respawn  ⟶ [Q] to Quit")
	}

	return commandStyle.Render("▶ [←J][L→] Move ✦ [K^] Jump ✦ [M] Mode ✦ [Q] Quit")
}

func (game *Game) getOffset() (int, int) {
	playerPos := game.Entities[0].Pos
	ox := playerPos.X - game.ViewWidth/3
	oy := 0

	if ox < 0 {
		ox = 0
	} else if ox > game.Width-game.ViewWidth {
		ox = game.Width - game.ViewWidth
	}

	return ox, oy
}
