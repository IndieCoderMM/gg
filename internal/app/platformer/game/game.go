package game

import (
	"fmt"
	"math/rand"
	"strings"
)

type EntityType int

const (
	Player EntityType = iota
	Enemy
	Ground
)

type Position struct {
	X, Y int
}

type Entity struct {
	ID     string
	Type   EntityType
	Pos    Position
	Sprite rune
	Text   rune
	State  map[string]any
}

type GameStatus string

type Game struct {
	Width, Height      int
	Entities           []Entity
	PlayerID           string
	Score              int
	Status             GameStatus
	Mode               string // "text" or "emoji"
	LastEnemySpawnPosX int
}

// Score: 0
// Hi: 10
// . . . . . . . . . .
// . . . . . . . . . .
// . p . . . x . . . .
// g g g g g g g g g g

const (
	JUMP_HEIGHT  = 3
	GROUND_LEVEL = 3
	WIN_WIDTH    = 20
	WIN_HEIGHT   = 4
)

func InitGameState() Game {
	player := Entity{
		ID:     "p1",
		Type:   Player,
		Pos:    Position{X: 1, Y: 2},
		Sprite: '😎',
		Text:   'p',
		State:  map[string]any{"jumping": false, "falling": false},
	}

	entities := []Entity{
		player,
		{ID: "spike1", Type: Enemy, Pos: Position{X: 25, Y: 2}, Sprite: '🐍', Text: 'x'},
		{ID: "spike2", Type: Enemy, Pos: Position{X: 35, Y: 2}, Sprite: '🐌', Text: 'x'},
		{ID: "bird1", Type: Enemy, Pos: Position{X: 50, Y: 1}, Sprite: '🦅', Text: 'e'},
		{ID: "spike3", Type: Enemy, Pos: Position{X: 65, Y: 2}, Sprite: '🗿', Text: 'x'},
		{ID: "bird2", Type: Enemy, Pos: Position{X: 80, Y: 0}, Sprite: '🦅', Text: 'e'},
	}

	randSprite := func(sprites []rune) rune {
		return sprites[rand.Intn(len(sprites))]
	}

	sprites := []rune{'🧱', '🟫'}
	// Add ground entities
	for i := range 20 {

		entities = append(entities, Entity{
			ID:   fmt.Sprintf("g%d", i),
			Type: Ground,
			Pos:  Position{X: i, Y: GROUND_LEVEL},
			// Select random ground sprite [#, =, -]
			Sprite: randSprite(sprites),
			Text:   randSprite([]rune{'n', 'm'}),
		})
	}

	return Game{
		Width:              WIN_WIDTH,
		Height:             WIN_HEIGHT,
		PlayerID:           "p1",
		Entities:           entities,
		Status:             "playing",
		Mode:               "text",
		LastEnemySpawnPosX: 80,
	}
}

func (game *Game) ToggleMode() {
	if game.Mode == "text" {
		game.Mode = "emoji"
	} else {
		game.Mode = "text"
	}
}

func (game Game) Draw() string {
	grid := make([][]rune, game.Height)
	for y := range grid {
		grid[y] = make([]rune, game.Width)
		for x := range grid[y] {
			if game.Mode == "emoji" {
				grid[y][x] = '🟦' // default empty
			} else {
				grid[y][x] = '.' // default empty
			}
		}
	}

	for _, e := range game.Entities {
		if e.Pos.Y >= 0 && e.Pos.Y < game.Height &&
			e.Pos.X >= 0 && e.Pos.X < game.Width {
			if game.Mode == "emoji" {
				grid[e.Pos.Y][e.Pos.X] = e.Sprite
			} else {
				grid[e.Pos.Y][e.Pos.X] = e.Text
			}
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Score: %d\n", game.Score))
	for _, row := range grid {
		for _, cell := range row {
			sb.WriteRune(cell)
			if game.Mode == "text" {
				sb.WriteRune(' ')
			}
		}
		sb.WriteString("\n")
	}
	if game.Status == "gameover" {
		sb.WriteString(">> Game Over! [R]: Restart\n")
	} else {
		sb.WriteString(">> [K]: Jump; [M]: Toggle mode; [Q]: Quit game\n")
	}
	return sb.String()
}

func (game *Game) Jump() {
	player := &game.Entities[0]
	if player.State["jumping"].(bool) || player.State["falling"].(bool) {
		return // already jumping or falling
	}

	player.State["jumping"] = true
	player.State["falling"] = false
}

func (game *Game) UpdatePlayer() {
	player := &game.Entities[0]

	if player.State["jumping"].(bool) {
		// If jumping, go up till JUMP_HEIGHT
		if player.Pos.Y > GROUND_LEVEL-JUMP_HEIGHT {
			player.Pos.Y -= 1
		} else {
			player.State["jumping"] = false
			player.State["falling"] = true
		}
		return
	}

	if player.State["falling"].(bool) {
		// If falling, go down till GROUND_LEVEL
		if player.Pos.Y < GROUND_LEVEL-1 {
			player.Pos.Y += 1
		} else {
			player.State["falling"] = false
		}
		return
	}
}

func (game *Game) Update() {
	if game.Status == "gameover" {
		return
	}

	game.Score += 1
	speed := 1

	game.LastEnemySpawnPosX -= speed // Move last enemy pos

	for i := range game.Entities {
		if game.Entities[i].ID == game.PlayerID {
			game.UpdatePlayer()
			continue
		}

		if game.Entities[i].Type == Enemy {
			game.UpdateEnemy(game.Entities[0], &game.Entities[i], speed)
			if game.CheckCollision(game.Entities[0], game.Entities[i]) {
				game.Status = "gameover"
				game.Entities[0].Sprite = '💀'
				game.Entities[0].Pos.Y = GROUND_LEVEL - 1
				game.Entities[0].Pos.X -= 1
			}
			continue
		}

		if game.Entities[i].Type == Ground {
			game.UpdateGround(&game.Entities[i], speed)
			continue
		}
	}
}

func (game *Game) UpdateEnemy(player Entity, enemy *Entity, speed int) {
	// Move enemy to the left
	newX := enemy.Pos.X - speed
	if newX < 0 {
		// Reset enemy to the right side of the screen + random offset
		newX = game.Width/2 + rand.Intn(game.Width/2) + game.LastEnemySpawnPosX
		game.LastEnemySpawnPosX = newX
	}

	enemy.Pos.X = newX
}

func (game *Game) CheckCollision(e1, e2 Entity) bool {
	if e1.Pos.X == e2.Pos.X && e1.Pos.Y == e2.Pos.Y {
		return true
	}
	return false
}

func (game *Game) UpdateGround(ground *Entity, speed int) {
	// Move ground to the left
	newX := ground.Pos.X - speed
	if newX < 0 {
		newX = game.Width - 1
	}
	ground.Pos.X = newX
}

func (game *Game) MoveEntity(e Entity, dx, dy int, outside bool) {
	newX := e.Pos.X + dx
	newY := e.Pos.Y + dy

	if !outside {
		// Check bounds if not outside
		if newX < 0 || newX >= game.Width || newY < 0 || newY >= game.Height {
			return
		}
	}

	e.Pos.X = newX
	e.Pos.Y = newY
}
