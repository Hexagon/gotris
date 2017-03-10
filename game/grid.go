package game

// Grid --------------------------------------------------------------------------------------------
type GameGrid struct {
	// Standard tetris game field is 10 units wide and 22 units high, with the topmost 2 rows hidden
	Data [10 * 22]rune
}

// Grid methods
func (g *GameGrid) Clear() {

	// Reset game field
	for i, _ := range g.Data {
		g.Data[i] = 0
	}

}

func (g *GameGrid) ApplySprite(s []Vector, p Vector, t rune) int32 {

	var clearedRows int32 = 0

	// Do the applying
	for _, v := range s {

		targetX := v.X + p.X
		targetY := v.Y + p.Y

		g.Data[targetX+targetY*10] = t

		// End condition, return -1!
		if targetY < 2 {
			return -1
		}

	}

	// Check if we shall pop a row
	for y := 2; y < 22; y++ {

		fullRow := true

		for x := 0; x < 10; x++ {
			if g.Data[x+y*10] == 0 {
				fullRow = false
			}
		}

		if fullRow {
			for y2 := y; y2 >= 1; y2-- {
				for x := 0; x < 10; x++ {
					g.Data[x+y2*10] = g.Data[x+(y2-1)*10]
				}
			}
			clearedRows += 1
		}
	}

	return clearedRows
}
