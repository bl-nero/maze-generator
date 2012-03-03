package generator

import (
	"board"
	"image"
)

func Generate(width, height int) board.Board {
	b := board.New(width, height)
	for y := 1; y < height-1; y++ {
		for x := 0; x < width; x++ {
			b.At(x, y).SetDirection(board.N | board.S)
		}
	}
	for x := 0; x < width; x++ {
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
		b.At(x, height-1).SetDirection(dir)
	}
	*b.Entrance() = image.Pt(0, height-1)
	if (width %2 == 0) {
		*b.Exit() = image.Pt(width-1, height-1)
	} else {
		*b.Exit() = image.Pt(width-1, 0)
	}
	return b
}
