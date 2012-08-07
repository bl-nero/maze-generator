package generator

import (
	"board"
	"container/heap"
	"image"
	"rand"
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
	*b.Entrance() = image.Pt(rand.Intn(width), 0)
	*b.Exit() = image.Pt(rand.Intn(width), height-1)
	boardRectangle := image.Rect(0, 0, width, height)
	fieldQueue := new(fieldHeap)
	heap.Init(fieldQueue)
	heap.Push(fieldQueue, fieldHeapElement{*b.Entrance(), rand.Int()})

	for fieldQueue.Len() > 0 {
		coords := heap.Pop(fieldQueue).(fieldHeapElement).Coords
		field := b.At(coords.X, coords.Y)
		possibleDirections := field.Direction().Negate().Decompose()
		shuffleDirections(possibleDirections)
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
			if !b.Exit().Eq(nextCoords) {
				heap.Push(fieldQueue,
					fieldHeapElement{nextCoords, rand.Int()})
			}
			heap.Push(fieldQueue, fieldHeapElement{coords, rand.Int()})
		}
	}

	b.At(b.Entrance().X, b.Entrance().Y).AddDirection(board.N)
	b.At(b.Exit().X, b.Exit().Y).AddDirection(board.S)

	return b
}

func shuffleDirections(directions []board.Direction) {
	for i := 0; i < len(directions); i++ {
		j := rand.Intn(i + 1)
		directions[i], directions[j] = directions[j], directions[i]
	}
}
