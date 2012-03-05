// Copyright 2012 Google Inc. All Rights Reserved.
// Author: bleper@google.com (Bartosz Leper)

package painter

import (
	"board"
	"image"
)

var (
	black = image.GrayColor{0}
	white = image.GrayColor{0xff}
)

func DrawRect(img *image.Gray, rect image.Rectangle, color image.GrayColor) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			img.SetGray(x, y, color)
		}
	}
}

func Paint(b board.Board, cellSize, wallThickness int) image.Image {
	width := b.Width()*(cellSize) + wallThickness
	height := b.Height()*(cellSize) + wallThickness
	img := image.NewGray(width, height)
	bounds := img.Bounds()

	DrawRect(img, bounds, white)

	for y := 0; y < b.Height(); y++ {
		yBase := y*cellSize + bounds.Min.Y
		for x := 0; x < b.Width(); x++ {
			xBase := x*cellSize + bounds.Min.X
			DrawRect(img, image.Rect(
				xBase,
				yBase,
				xBase+wallThickness,
				yBase+wallThickness), black)
			dir := b.At(x, y).Direction()
			if dir&board.N == 0 {
				DrawRect(img, image.Rect(
					xBase+wallThickness,
					yBase,
					xBase+cellSize,
					yBase+wallThickness), black)
			}
			if dir&board.W == 0 {
				DrawRect(img, image.Rect(
					xBase,
					yBase+wallThickness,
					xBase+wallThickness,
					yBase+cellSize), black)
			}
		}
	}

	DrawRect(img,
		image.Rectangle{
			image.Point{bounds.Min.X, bounds.Max.Y - wallThickness},
			bounds.Max,
		},
		black)
	DrawRect(img,
		image.Rectangle{
			image.Point{bounds.Max.X - wallThickness, bounds.Min.Y},
			bounds.Max,
		},
		black)

	return img
}
