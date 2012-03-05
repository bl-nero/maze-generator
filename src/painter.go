// Copyright 2012 Google Inc. All Rights Reserved.
// Author: bleper@google.com (Bartosz Leper)

package painter

import (
	"board"
	"image"
)

var (
	wallColor  = image.RGBAColor{0, 0, 0, 0xff}
	boardColor = image.RGBAColor{0xff, 0xff, 0xff, 0xff}
	pathColor  = image.RGBAColor{0, 0xff, 0, 0xff}
)

func DrawRect(img *image.RGBA, rect image.Rectangle, color image.RGBAColor) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			img.SetRGBA(x, y, color)
		}
	}
}

func Paint(b board.Board, visitMatrix [][]bool, cellSize, wallThickness int) image.Image {
	width := b.Width()*(cellSize) + wallThickness
	height := b.Height()*(cellSize) + wallThickness
	img := image.NewRGBA(width, height)
	bounds := img.Bounds()

	DrawRect(img, bounds, boardColor)

	for y := 0; y < b.Height(); y++ {
		yBase := y*cellSize + bounds.Min.Y
		for x := 0; x < b.Width(); x++ {
			xBase := x*cellSize + bounds.Min.X
			if visitMatrix != nil && visitMatrix[y][x] {
				DrawRect(img, image.Rect(
					xBase+wallThickness,
					yBase+wallThickness,
					xBase+cellSize,
					yBase+cellSize), pathColor)
			}
			DrawRect(img, image.Rect(
				xBase,
				yBase,
				xBase+wallThickness,
				yBase+wallThickness), wallColor)
			dir := b.At(x, y).Direction()
			if dir&board.N == 0 {
				DrawRect(img, image.Rect(
					xBase+wallThickness,
					yBase,
					xBase+cellSize,
					yBase+wallThickness), wallColor)
			}
			if dir&board.W == 0 {
				DrawRect(img, image.Rect(
					xBase,
					yBase+wallThickness,
					xBase+wallThickness,
					yBase+cellSize), wallColor)
			}
		}
		var lastHeight int
		if b.At(b.Width()-1, y).Direction()&board.E == 0 {
			lastHeight = cellSize
		} else {
			lastHeight = wallThickness
		}
		xBase := b.Width()*cellSize + bounds.Min.X
		DrawRect(img, image.Rect(
			xBase,
			yBase,
			xBase+wallThickness,
			yBase+lastHeight), wallColor)
	}

	yBase := b.Height()*cellSize + bounds.Min.Y
	for x := 0; x < b.Width(); x++ {
		var lastWidth int
		if b.At(x, b.Height()-1).Direction()&board.S == 0 {
			lastWidth = cellSize
		} else {
			lastWidth = wallThickness
		}
		xBase := x*cellSize + bounds.Min.X
		DrawRect(img, image.Rect(
			xBase,
			yBase,
			xBase+lastWidth,
			yBase+wallThickness), wallColor)
	}

	DrawRect(img, image.Rectangle{
		bounds.Max.Add(image.Pt(-wallThickness, -wallThickness)),
		bounds.Max,
	}, wallColor)

	return img
}
