package generator

import (
	"board"
	"container/heap"
	"image"
)

type fieldHeapElement struct {
	Coords image.Point
	Weight int
}

type fieldHeap []fieldHeapElement

func (self *fieldHeap) Push(x interface{}) {
	*self = append(*self, x.(fieldHeapElement))
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
	if width < 1 || height < 1 {
		return nil
	}

	b := board.New(width, height)
	boardRectangle := image.Rect(0, 0, width, height)
	fieldQueue := new(fieldHeap)
	heap.Init(fieldQueue)
	heap.Push(fieldQueue, fieldHeapElement{*b.Entrance(), 0})
	currentWeight := 1

	for fieldQueue.Len() > 0 {
		coords := heap.Pop(fieldQueue).(fieldHeapElement).Coords
		field := b.At(coords.X, coords.Y)
		possibleDirections := field.Direction().Negate().Decompose()
		pickedDirection := board.None
		var nextCoords image.Point
		for _, dir := range possibleDirections {
			delta, _ := dir.Delta()
			possibleCoords := coords.Add(delta)
			if possibleCoords.In(boardRectangle) &&
				b.At(possibleCoords.X, possibleCoords.Y).Direction() ==
					board.None {
				pickedDirection = dir
				nextCoords = possibleCoords
				break
			}
		}
		if pickedDirection != board.None {
			nextField := b.At(nextCoords.X, nextCoords.Y)
			field.AddDirection(pickedDirection)
			nextField.AddDirection(pickedDirection.Opposite())
			heap.Push(fieldQueue, fieldHeapElement{nextCoords, currentWeight})
			heap.Push(fieldQueue, fieldHeapElement{coords, currentWeight})
		}
	}

	return b
}
