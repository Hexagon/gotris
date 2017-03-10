package game

import (
	"encoding/json"
	"fmt"
	"time"
)

// Game --------------------------------------------------------------------------------------------
type Game struct {

	// The 10x22 grid of history, the first two lines should be hidden in visualization
	grid GameGrid

	// Five position array of tetrominoes, where index
	//  0: Current dropping tetromino
	//  1: Coming up 1
	//  2: Coming up 2
	//  3: Coming up 3
	//  -- 4: On hold (TODO)
	Tetrominoes [5]Tetromino
	factory     TetrominoFactory

	// Base offset of dropping tetromino
	// Can be outside of grid if current tetromino sprite has blank rows/columns
	Position      Vector
	GhostPosition Vector

	// Rotation of dropping tetromino
	// Should be >0 and <4
	Rotation int8

	Score uint64
	Level uint64
	Lines uint64

	// Internal stuff below
	Nickname    string
	keyStates   *map[string]interface{}
	dataChannel chan string

	timerReal     time.Time
	timerModified time.Time
}

func NewGame(nickname string, dataChannel chan string) Game {

	game := Game{}

	// Set nickname of this game
	game.Nickname = nickname

	// Reset stuff
	game.Score = 0
	game.Level = 0
	game.grid.Clear()

	// Initialize upcoming tetrominos
	game.Tetrominoes[1] = game.factory.Next()
	game.Tetrominoes[2] = game.factory.Next()
	game.Tetrominoes[3] = game.factory.Next()

	// Set current
	game.NextTetromino()

	// Set up outgoing data channel
	game.dataChannel = dataChannel

	// Initial game update
	game.dataChannel <- game.updateGrid()

	return game
}

func (game *Game) updateGame() string {

	// Marshal gamefield and player
	b, err := json.Marshal(game)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)

}

func (game *Game) updateGrid() string {

	// Marshal gamefield and player
	b, err := json.Marshal(game.grid)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)

}

func (game *Game) Changed() {

	currentSprite := game.Tetrominoes[0].Sprites[game.Rotation].Data

	bogusPosition := game.Position
	bogusPosition.Y += 1

	// Drop til it stop
	for game.validMove(currentSprite, bogusPosition) {
		game.GhostPosition = bogusPosition

		bogusPosition.Y += 1
	}

	game.dataChannel <- game.updateGame()

}

func (game *Game) GetGrid() GameGrid {
	return game.grid
}

func (game *Game) MoveX(d int32) bool {

	game.resetTimerIfLanded()

	currentSprite := game.Tetrominoes[0].Sprites[game.Rotation].Data

	bogusPosition := game.Position
	bogusPosition.X += d

	if game.validMove(currentSprite, bogusPosition) {
		// Seem ok, apply position
		game.Position = bogusPosition
		game.Changed()
		return true
	}

	return false

}

func (game *Game) SetKeyStates(kss *map[string]interface{}) {
	game.keyStates = kss
}

func (game *Game) SetKey(event map[string]interface{}) {

	key, keyOk := event["key"]
	if !keyOk {
		// Key "key" did not exist, bail out
		return
	}
	state, stateOk := event["state"]
	if !stateOk {
		// Key "state" did not exist, bail out
		return
	}

	// Check if key was pressed (hence not released)
	if state.(bool) {
		// Do immidiate actions
		switch key {

		case "right":
			game.MoveX(1)

		case "left":
			game.MoveX(-1)

		case "drop":
			game.Drop()

		case "rotCW":
			game.Rotate(1)

		case "rotCCW":
			game.Rotate(-1)

		case "down":
			game.MoveDown()

		}
	}

	// ToDo: Keep states
}

func (game *Game) MoveDown() bool {

	currentSprite := game.Tetrominoes[0].Sprites[game.Rotation].Data

	bogusPosition := game.Position
	bogusPosition.Y += 1

	game.resetTimerIfLanded()

	if game.validMove(currentSprite, bogusPosition) {
		// Seem ok, apply position
		game.Position = bogusPosition
		game.Changed()
		return true
	} else {
		return false
	}

}

func (game *Game) hasLanded() bool {

	currentSprite := game.Tetrominoes[0].Sprites[game.Rotation].Data

	bogusPosition := game.Position
	bogusPosition.Y += 1

	if game.validMove(currentSprite, bogusPosition) {
		return false
	} else {
		return true
	}

}

func (game *Game) Drop() {

	currentSprite := game.Tetrominoes[0].Sprites[game.Rotation].Data

	bogusPosition := game.Position
	bogusPosition.Y += 1

	var dropOffset uint32 = 0

	// Drop til it stops
	for game.validMove(currentSprite, bogusPosition) {
		game.Position = bogusPosition
		dropOffset += 1
		bogusPosition.Y += 1
	}

	game.AddScore(uint64(dropOffset*2), false)

	game.Lockdown()

}

func (game *Game) resetTimerIfLanded() {
	if game.hasLanded() {
		// Allow stuffing around for 2 seconds
		if time.Since(game.timerReal) < 2*time.Second {
			game.timerModified = time.Now()
		}
	}
}

func (game *Game) Rotate(d int8) {

	game.resetTimerIfLanded()

	bogusRotation := game.Rotation
	bogusRotation += d

	if bogusRotation < 0 {
		bogusRotation = 4 - 1
	} else if bogusRotation > 4-1 {
		bogusRotation = 0
	}

	currentSprite := game.Tetrominoes[0].Sprites[bogusRotation].Data

	// Check if the desired rotation was a valid move
	if game.validMove(currentSprite, game.Position) {
		// Seem ok, apply position
		game.Rotation = bogusRotation
		game.Changed()
		return
	}

	// Try wall kick left
	bogusPosition := game.Position
	bogusPosition.X -= 1
	if game.validMove(currentSprite, bogusPosition) {
		game.Rotation = bogusRotation
		game.Position = bogusPosition
		game.Changed()
		return
	}

	// Try wall kick right
	bogusPosition = game.Position
	bogusPosition.X += 1
	if game.validMove(currentSprite, bogusPosition) {
		game.Rotation = bogusRotation
		game.Position = bogusPosition
		game.Changed()
		return
	}

	// Try double kick left
	bogusPosition = game.Position
	bogusPosition.X -= 2
	if game.validMove(currentSprite, bogusPosition) {
		game.Rotation = bogusRotation
		game.Position = bogusPosition
		game.Changed()
		return
	}

	// Try double kick right
	bogusPosition = game.Position
	bogusPosition.X += 2
	if game.validMove(currentSprite, bogusPosition) {
		game.Rotation = bogusRotation
		game.Position = bogusPosition
		game.Changed()
		return
	}

}

func (game *Game) validMove(s []Vector, p Vector) bool {
	for _, v := range s {

		// Check that there is space below
		targetX := v.X + p.X
		targetY := v.Y + p.Y

		// Check X
		outOfRange := targetX < 0 || targetX > 10-1

		// Check Y
		outOfRange = outOfRange || targetY < 0 || targetY > 22-1 || game.grid.Data[targetX+targetY*10] > 0

		if outOfRange {
			return false
		}
	}

	return true
}

// Game methods
func (game *Game) Iterate() bool {

	if time.Since(game.timerModified) > game.iterateDelayMs()*time.Millisecond {

		// Reset timer
		game.timerReal = time.Now()
		game.timerModified = time.Now()

		if !game.MoveDown() {
			return game.Lockdown()
		}

	}

	return true

}

func (game *Game) iterateDelayMs() time.Duration {

	// Reduce 10 ms for each level
	delay := 500 - 20*game.Level

	if delay < 120 {
		delay = 120
	}

	return time.Duration(delay)

}

func (game *Game) NextTetromino() {
	game.Position = Vector{3, 0}
	game.Rotation = 0
	game.Tetrominoes[0] = game.Tetrominoes[1]
	game.Tetrominoes[1] = game.Tetrominoes[2]
	game.Tetrominoes[2] = game.Tetrominoes[3]
	game.Tetrominoes[3] = game.factory.Next()
}

func (game *Game) AddScore(baseScore uint64, levelBoost bool) {
	if levelBoost {
		game.Score += baseScore * (game.Level + 1)
	} else {
		game.Score += baseScore
	}
}

func (game *Game) Lockdown() bool {

	var clearedRows int32
	clearedRows = game.grid.ApplySprite(game.Tetrominoes[0].Sprites[game.Rotation].Data, game.Position, game.Tetrominoes[0].Type)

	// End condition ?
	if clearedRows == -1 {
		return false
	}

	// Count cleared lines
	game.Lines += uint64(clearedRows)

	// Add score for cleared rows
	switch clearedRows {
	case 1:
		game.AddScore(40, true)
	case 2:
		game.AddScore(40*2*2, true) // x2
	case 3:
		game.AddScore(40*3*4, true) // x4
	case 4:
		game.AddScore(40*4*8, true) // x8
	}

	// Time to level up?
	if game.Lines > (game.Level+1)*5 {
		game.Level += 1
	}

	// Spawn new tetromino!
	game.NextTetromino()

	// Reset drop button (PANIC safe?)
	if game.keyStates != nil {
		(*game.keyStates)["down"] = false
	}

	// Notify that stuf has changed
	game.Changed()
	game.dataChannel <- game.updateGrid()

	return true

}
