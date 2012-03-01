package generator

import (
	"board"
)

func Generate(width, height int) board.Board {
	b := board.New(width, height)
	for y := 1; y < b.Height()-1; y++ {
		for x := 0; x < b.Width(); x++ {
			b.At(x, y).SetDirection(board.N | board.S)
		}
	}
	for x := 0; x < b.Width(); x++ {
		var dir board.Direction
		if x%2 == 0 {
			dir = board.S | board.E
		} else {
			dir = board.S | board.W
		}
		b.At(x, 0).SetDirection(dir)
		if x%2 == 0 {
			dir = board.N | board.W
		} else {
			dir = board.N | board.E
		}
		b.At(x, b.Height()-1).SetDirection(dir)
	}
	return b
}
