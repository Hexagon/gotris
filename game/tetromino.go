package game

import "math/rand"

type Vector struct {
	X int32
	Y int32
}

type Sprite struct {
	Data []Vector
}

var implementedTetrominos = []rune{'I', 'J', 'L', 'O', 'S', 'T', 'Z'}

type Tetromino struct {
	Type    rune
	Sprites []Sprite
}

func NewTetrominoFactory() TetrominoFactory {
	return TetrominoFactory{make([]Tetromino, 7)}
}

type TetrominoFactory struct {
	Inventory []Tetromino
}

func (tf *TetrominoFactory) getTetromino(tt rune) Tetromino {
	switch tt {
	case 'I':
		return Tetromino{tt, []Sprite{
			Sprite{[]Vector{{0, 1}, {1, 1}, {2, 1}, {3, 1}}},
			Sprite{[]Vector{{2, 0}, {2, 1}, {2, 2}, {2, 3}}},
			Sprite{[]Vector{{0, 2}, {1, 2}, {2, 2}, {3, 2}}},
			Sprite{[]Vector{{1, 0}, {1, 1}, {1, 2}, {1, 3}}},
		}}
	case 'J':
		return Tetromino{tt, []Sprite{
			Sprite{[]Vector{{0, 0}, {0, 1}, {1, 1}, {2, 1}}},
			Sprite{[]Vector{{1, 0}, {2, 0}, {1, 1}, {1, 2}}},
			Sprite{[]Vector{{0, 1}, {1, 1}, {2, 1}, {2, 2}}},
			Sprite{[]Vector{{1, 0}, {1, 1}, {1, 2}, {0, 2}}},
		}}
	case 'L':
		return Tetromino{tt, []Sprite{
			Sprite{[]Vector{{2, 0}, {0, 1}, {1, 1}, {2, 1}}},
			Sprite{[]Vector{{1, 0}, {2, 2}, {1, 1}, {1, 2}}},
			Sprite{[]Vector{{0, 1}, {1, 1}, {2, 1}, {0, 2}}},
			Sprite{[]Vector{{1, 0}, {1, 1}, {1, 2}, {0, 0}}},
		}}
	case 'O':
		return Tetromino{tt, []Sprite{
			Sprite{[]Vector{{1, 0}, {2, 0}, {1, 1}, {2, 1}}},
			Sprite{[]Vector{{1, 0}, {2, 0}, {1, 1}, {2, 1}}},
			Sprite{[]Vector{{1, 0}, {2, 0}, {1, 1}, {2, 1}}},
			Sprite{[]Vector{{1, 0}, {2, 0}, {1, 1}, {2, 1}}},
		}}
	case 'S':
		return Tetromino{tt, []Sprite{
			Sprite{[]Vector{{0, 1}, {1, 1}, {1, 0}, {2, 0}}},
			Sprite{[]Vector{{1, 0}, {1, 1}, {2, 1}, {2, 2}}},
			Sprite{[]Vector{{0, 2}, {1, 2}, {1, 1}, {2, 1}}},
			Sprite{[]Vector{{0, 0}, {0, 1}, {1, 1}, {1, 2}}},
		}}
	case 'T':
		return Tetromino{tt, []Sprite{
			Sprite{[]Vector{{0, 1}, {1, 1}, {1, 0}, {2, 1}}},
			Sprite{[]Vector{{1, 0}, {1, 1}, {2, 1}, {1, 2}}},
			Sprite{[]Vector{{0, 1}, {1, 2}, {1, 1}, {2, 1}}},
			Sprite{[]Vector{{1, 0}, {0, 1}, {1, 1}, {1, 2}}},
		}}
	case 'Z':
		return Tetromino{tt, []Sprite{
			Sprite{[]Vector{{0, 0}, {1, 1}, {1, 0}, {2, 1}}},
			Sprite{[]Vector{{2, 0}, {1, 1}, {2, 1}, {1, 2}}},
			Sprite{[]Vector{{0, 1}, {1, 2}, {1, 1}, {2, 2}}},
			Sprite{[]Vector{{1, 0}, {0, 1}, {1, 1}, {0, 2}}},
		}}
	default:
		panic("Tried to initialize unknown tetromino type")
	}
}

func (tf *TetrominoFactory) Next() Tetromino {

	// Refill factory ?
	if len(tf.Inventory) == 0 {
		for i := 0; i < 7; i++ {
			tf.Inventory = append(tf.Inventory, tf.getTetromino(implementedTetrominos[i]))
		}
	}

	// Draw next

	// - Get a random index between 0 and <no of tetrominoes left>
	ri := rand.Intn(len(tf.Inventory))

	// - Store reference to found tetromino before removing it from factory
	next := tf.Inventory[ri]

	// - Remove from factory
	tf.Inventory = append(tf.Inventory[:ri], tf.Inventory[ri+1:]...)

	return next

}
