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
	Star
)

type Coor struct {
	X, Y float32
}

type Entity struct {
	ID    string
	Type  EntityType
	Pos   Coor
	Text  rune
	State map[string]any
}

type GameStatus string

type Game struct {
	Width, Height         float32
	ViewWidth, ViewHeight float32
	Entities              []Entity
	PlayerID              string
	Score                 int
	Lives                 int
	Status                GameStatus
	lastEnemySpawnPt      float32
	tick                  int
	tickRate              float32
}

// Score: 0
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
		ID:    "p1",
		Type:  Player,
		Pos:   Coor{X: 3, Y: 2},
		Text:  'O',
		State: map[string]any{"jumping": false, "falling": false},
	}

	entities := []Entity{
		player,
	}

	// Add stars
	for i := range 10 {
		entities = append(entities, Entity{
			ID:   fmt.Sprintf("star%d", i),
			Type: Star,
			Pos:  Coor{X: float32(rand.Intn(2 * WIN_WIDTH)), Y: GROUND_LEVEL - 2},
			Text: '*',
		})
	}

	randSprite := func(sprites []rune) rune {
		return sprites[rand.Intn(len(sprites))]
	}

	// Add ground entities
	for i := range WIN_WIDTH {

		entities = append(entities, Entity{
			ID:   fmt.Sprintf("g%d", i),
			Type: Ground,
			Pos:  Coor{X: float32(i), Y: GROUND_LEVEL},
			Text: randSprite([]rune{'░', '▒', '▒', '▒'}),
		})
	}

	enemySprites := []rune{'@', '#', 'x'}
	// Add enemies
	for i := range 5 {
		// Random enemy spawn point
		enemyPos := Coor{
			X: float32(randRange(VIEW_WIDTH, WIN_WIDTH)),
			Y: float32(randRange(1, GROUND_LEVEL)),
		}

		if rand.Int()%3 != 0 {
			enemyPos.Y = GROUND_LEVEL - 1
		}

		entities = append(entities, Entity{
			ID:    fmt.Sprintf("e%d", i),
			Type:  Enemy,
			Pos:   enemyPos,
			Text:  randSprite(enemySprites),
			State: map[string]any{"alive": true, "moving": true},
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
		Lives:            3,
		lastEnemySpawnPt: 80, // ? To prevent overlapping
		tick:             0,
		tickRate:         4, // Higher is slower
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

// func (game *Game) Move(dir int) {
// 	player := &game.Entities[0]
//
// 	if dir > 0 {
// 		// Move right
// 		if player.Pos.X < game.Width-1 {
// 			player.Pos.X += 1
// 		}
// 	} else if dir < 0 {
// 		// Move left
// 		if player.Pos.X > 0 {
// 			player.Pos.X -= 1
// 		}
// 	}
// }

func (game *Game) Update() {
	if game.Status == "gameover" {
		return
	}

	game.tick++

	if game.tick%int(game.tickRate) != 0 {
		return
	}

	var speed float32 = 1

	for i := range game.Entities {
		if game.Entities[i].ID == game.PlayerID {
			game.updatePlayer()
			continue
		}

		if game.Entities[i].Type == Enemy {
			game.updateEnemy(&game.Entities[i], speed)
			// Check collision with player
			if game.isCollided(game.Entities[0], game.Entities[i]) {
				game.Lives -= 1

				if game.Lives == 0 {
					game.over()
				}
			}
			continue
		}

		if game.Entities[i].Type == Ground {
			game.updateGround(&game.Entities[i], speed)
			continue
		}

		if game.Entities[i].Type == Star {
			game.Entities[i].Pos.X -= speed
			if game.isCollided(game.Entities[0], game.Entities[i]) {
				game.Entities[i].Pos.X += game.ViewWidth + float32(rand.Intn(WIN_WIDTH))
				// Update score & increase speed
				game.Score += 1
				if game.Score%10 == 0 && game.Score > 0 && game.tickRate > 1 {
					game.tickRate -= 0.25
				}
			}

			if game.Entities[i].Pos.X < 0 {
				// Relocate star to right
				game.Entities[i].Pos.X += game.ViewWidth + float32(rand.Intn(WIN_WIDTH))
			}
		}
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

func (game *Game) updateEnemy(enemy *Entity, speed float32) {
	if enemy.State["alive"] == false {
		return
	}

	// Check if enemy is moving
	if enemy.State["moving"] == false {
		playerPos := game.Entities[0].Pos
		// Check enemy is within viewport
		if enemy.Pos.X < playerPos.X+game.ViewWidth {
			// Move enemy to left
			enemy.State["moving"] = true
		}
		return
	}

	// Move enemy if moving
	newX := enemy.Pos.X - speed
	if newX < 0 {
		// Random enemy spawn point
		newX = game.Entities[0].Pos.X + game.lastEnemySpawnPt + float32(rand.Intn(int(game.ViewWidth/2)))
		// Unalive enemy
		// enemy.State["alive"] = false
	}

	enemy.Pos.X = newX
}

func (game *Game) isCollided(e1, e2 Entity) bool {
	if int(e1.Pos.X) == int(e2.Pos.X) && int(e1.Pos.Y) == int(e2.Pos.Y) {
		return true
	}
	return false
}

func (game *Game) updateGround(ground *Entity, speed float32) {
	playerPos := game.Entities[0].Pos
	// Check if ground is within viewport
	ground.Pos.X -= speed

	if ground.Pos.X < playerPos.X-game.ViewWidth-3 {
		// Relocate ground to right
		ground.Pos.X += WIN_WIDTH
		game.Width += 1
	}
}

func (game *Game) over() {
	game.Status = "gameover"
	game.Entities[0].Text = '0'
	game.Entities[0].Pos.Y = GROUND_LEVEL - 1
	game.Entities[0].Pos.X -= 1
}

func randRange(left, right int) int {
	return left + rand.Intn(right-left)
}
