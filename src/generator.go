package generator

import (
	"board"
	"image"
	"fmt"
	"reflect"
)

type fieldHeapElement struct {
	Coords image.Point
	Weight int
}

type fieldHeap []fieldHeapElement

func (self *fieldHeap) Push(x interface{}) {
	elem, ok := x.(fieldHeapElement)
	if !ok {
		panic(fmt.Sprintf("Tried to push %s into a fieldHeap",
			reflect.TypeOf(x).Name()))
	}
	*self = append(*self, elem)
}

func (self *fieldHeap) Pop() interface{} {
	last := len(*self) - 1
	result := (*self)[last]
	*self = (*self)[:last]
	return result
}

func (self fieldHeap) Len() int           { return len(self) }
func (self fieldHeap) Less(i, j int) bool { return self[i].Weight < self[j].Weight }
func (self fieldHeap) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }

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
	if width%2 == 0 {
		*b.Exit() = image.Pt(width-1, height-1)
	} else {
		*b.Exit() = image.Pt(width-1, 0)
	}
	return b
}
