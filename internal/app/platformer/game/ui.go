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
		labelStyle.Render("SCORE"),
		scoreStyle.Render(fmt.Sprintf("%d", game.Score)))

	lives := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("LIVES"),
		livesStyle.Render(fmt.Sprintf("%d", game.Lives)))

	var spacer int
	var hudRow string

	scoreLen := len(fmt.Sprintf("SCORE %d", game.Score))

	if game.Status == "gameover" {
		lives = livesStyle.Render("GAME OVER")
		spacer = 2*int(game.ViewWidth) - scoreLen - len("GAME OVER")
	} else {
		spacer = 2*int(game.ViewWidth) - scoreLen - len("LIVES 0")
	}

	hudRow = lipgloss.JoinHorizontal(lipgloss.Top, lives, spaceStyle.Width(spacer).Render(" "), score)
	return hudContainer.Render(hudRow)
}

func (game Game) renderViewport() string {
	offsetX, offsetY := game.getOffset()
	var sb strings.Builder

	// Store this grid in state
	grid := make([][]string, int(game.ViewHeight))
	for y := range grid {
		grid[y] = make([]string, int(game.ViewWidth))
		for x := range grid[y] {
			if game.Mode == "emoji" {
				grid[y][x] = "▪️" // default empty
			} else {
				grid[y][x] = "." // default empty
			}
		}
	}

	for _, e := range game.Entities {
		x, y := int(e.Pos.X-offsetX), int(e.Pos.Y-offsetY)
		if x < 0 || x >= int(game.ViewWidth) || y < 0 || y >= int(game.ViewHeight) {
			continue // Out of bounds
		}

		if game.Mode == "emoji" {
			grid[y][x] = string(e.Sprite)
		} else {
			grid[y][x] = string(e.Text)
		}
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

		if y != GROUND_LEVEL {
			sb.WriteRune('\n')
		}
	}

	return viewportBorder.Render(sb.String())
}

func (game Game) renderFooter() string {
	if game.Status == "gameover" {
		return commandStyle.Render("▶ [R] to Respawn ✦ [Q] to Quit")
	}

	return commandStyle.Render("▶ [K^] Jump ✦ [Q] Quit")
}

func (game *Game) getOffset() (float32, float32) {
	playerPos := game.Entities[0].Pos
	ox := playerPos.X - game.ViewWidth/3
	var oy float32 = 0

	if ox < 0 {
		ox = 0
	} else if ox > game.Width-game.ViewWidth {
		ox = game.Width - game.ViewWidth
	}

	return ox, oy
}
