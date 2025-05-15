package game

import (
	"fmt"
	"math/rand"
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
	Width, Height         int
	ViewWidth, ViewHeight int
	Entities              []Entity
	PlayerID              string
	Score                 int
	Lives                 int
	Status                GameStatus
	Mode                  string // "text" or "emoji"
	lastEnemySpawnPt      int
}

// Score: 0
// Hi: 10
// . . . . . . . . . .
// . . . . . . . . . .
// . p . . . x . . . .
// g g g g g g g g g g

const (
	JUMP_HEIGHT  = 2
	GROUND_LEVEL = 3
	WIN_WIDTH    = 100
	WIN_HEIGHT   = 4
	VIEW_WIDTH   = 20
)

func InitGameState() Game {
	player := Entity{
		ID:     "p1",
		Type:   Player,
		Pos:    Position{X: 2, Y: 2},
		Sprite: '😎',
		Text:   'p',
		State:  map[string]any{"jumping": false, "falling": false},
	}

	entities := []Entity{
		player,
		{ID: "spike1", Type: Enemy, Pos: Position{X: 25, Y: 2}, Sprite: '🐍', Text: 'x'},
		{ID: "spike2", Type: Enemy, Pos: Position{X: 35, Y: 2}, Sprite: '🐌', Text: 'x'},
		{ID: "bird1", Type: Enemy, Pos: Position{X: 50, Y: 1}, Sprite: '🦇', Text: 'e'},
		{ID: "spike3", Type: Enemy, Pos: Position{X: 65, Y: 2}, Sprite: '🔥', Text: 'x'},
		{ID: "bird2", Type: Enemy, Pos: Position{X: 80, Y: 0}, Sprite: '🦅', Text: 'e'},
	}

	randSprite := func(sprites []rune) rune {
		return sprites[rand.Intn(len(sprites))]
	}

	gndSprites := []rune{'🧱', '🟫'}
	// Add ground entities
	for i := range WIN_WIDTH {

		entities = append(entities, Entity{
			ID:   fmt.Sprintf("g%d", i),
			Type: Ground,
			Pos:  Position{X: i, Y: GROUND_LEVEL},
			// Select random ground sprite [#, =, -]
			Sprite: randSprite(gndSprites),
			Text:   randSprite([]rune{'░', '▒'}),
		})
	}

	return Game{
		Width:            WIN_WIDTH,
		Height:           WIN_HEIGHT,
		ViewWidth:        VIEW_WIDTH,
		ViewHeight:       WIN_HEIGHT,
		PlayerID:         "p1",
		Entities:         entities,
		Status:           "playing",
		Mode:             "text",
		lastEnemySpawnPt: 80, // ? To prevent overlapping
	}
}

func (game *Game) ToggleMode() {
	if game.Mode == "text" {
		game.Mode = "emoji"
	} else {
		game.Mode = "text"
	}
}

func (game *Game) Jump() {
	player := &game.Entities[0]
	if player.State["jumping"].(bool) || player.State["falling"].(bool) {
		return // already jumping or falling
	}

	player.State["jumping"] = true
	player.State["falling"] = false
}

func (game *Game) Move(dir int) {
	player := &game.Entities[0]

	if dir > 0 {
		// Move right
		if player.Pos.X < game.Width-1 {
			player.Pos.X += 1
		}
	} else if dir < 0 {
		// Move left
		if player.Pos.X > 0 {
			player.Pos.X -= 1
		}
	}
}

func (game *Game) Update() {
	if game.Status == "gameover" {
		return
	}

	// game.Score += 1
	speed := 1

	game.lastEnemySpawnPt -= speed // Move last enemy pos

	for i := range game.Entities {
		if game.Entities[i].ID == game.PlayerID {
			game.updatePlayer()
			continue
		}

		// if game.Entities[i].Type == Enemy {
		// 	game.updateEnemy(&game.Entities[i], speed)
		// 	// Check collision with player
		// 	if game.isCollided(game.Entities[0], game.Entities[i]) {
		// 		game.over()
		// 	}
		// 	continue
		// }
		//
		// if game.Entities[i].Type == Ground {
		// 	game.updateGround(&game.Entities[i], speed)
		// 	continue
		// }
	}
}

// Handle player movement
func (game *Game) updatePlayer() {
	player := &game.Entities[0]

	if player.State["jumping"].(bool) {
		// If jumping, go up till JUMP_HEIGHT
		if player.Pos.Y+1 > GROUND_LEVEL-JUMP_HEIGHT {
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

func (game *Game) updateEnemy(enemy *Entity, speed int) {
	// Move enemy to the left
	newX := enemy.Pos.X - speed
	if newX < 0 {
		// Reset enemy to the right side of the screen + random offset
		newX = game.Width/2 + rand.Intn(game.Width/2) + game.lastEnemySpawnPt
		game.lastEnemySpawnPt = newX
	}

	enemy.Pos.X = newX
}

func (game *Game) isCollided(e1, e2 Entity) bool {
	if e1.Pos.X == e2.Pos.X && e1.Pos.Y == e2.Pos.Y {
		return true
	}
	return false
}

func (game *Game) updateGround(ground *Entity, speed int) {
	// Move ground to the left
	newX := ground.Pos.X - speed
	if newX < 0 {
		newX = game.Width - 1
	}
	ground.Pos.X = newX
}

func (game *Game) over() {
	game.Status = "gameover"
	game.Entities[0].Sprite = '💀'
	game.Entities[0].Pos.Y = GROUND_LEVEL - 1
	game.Entities[0].Pos.X -= 1
}
