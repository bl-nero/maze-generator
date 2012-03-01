package board

import (
	"bytes"
	"image"
	"os"
)

type Direction uint8

const (
	None Direction = 0
	N    Direction = 1 << iota
	E
	S
	W
	visitedBit    Field = 1 << iota
	minDirection        = N
	maxDirection        = W
	directionMask uint8 = uint8(maxDirection)<<1 - 1
)

var dirNames map[Direction]string = map[Direction]string{N: "N", E: "E", S: "S", W: "W"}

func (self Direction) String() string {
	if self == None {
		return "None"
	}
	res := ""
	for dir, name := range dirNames {
		if self&dir != 0 {
			res += name
		}
	}
	return res
}

var dirDeltas map[Direction]image.Point = map[Direction]image.Point{
	N: {0, -1},
	E: {1, 0},
	S: {0, 1},
	W: {-1, 0},
}

func (self Direction) Delta() (delta image.Point, error os.Error) {
	delta, ok := dirDeltas[self]
	if !ok {
		error = os.NewError("Unable to fetch delta of a composite direction " +
			self.String())
	}
	return
}

func (self Direction) Decompose() []Direction {
	result := make([]Direction, 0, 4)
	for d := minDirection; d <= maxDirection; d <<= 1 {
		if self&d != 0 {
			result = append(result, d)
		}
	}
	return result
}

type Field uint8

func (f Field) Direction() Direction {
	return Direction(f & Field(directionMask))
}

func (f *Field) SetDirection(dir Direction) {
	*f = Field(uint8(*f)&^directionMask | uint8(dir))
}

func (f *Field) AddDirection(dir Direction) {
	*f |= Field(dir)
}

func (f Field) visited() bool {
	return f&visitedBit != 0
}

func (f *Field) setVisited(visited bool) {
	if visited {
		*f = Field(visitedBit)
	} else {
		*f = Field(0)
	}
}

type Board interface {
	Width() int
	Height() int
	At(x, y int) *Field
	Entrance() *image.Point
	Exit() *image.Point
	Walk() ([][]bool, os.Error)
	//String() string
}

func New(width, height int) Board {
	board := boardImpl{fields: make([][]Field, height)}
	for i := range board.fields {
		board.fields[i] = make([]Field, width)
	}
	return &board
}

type boardImpl struct {
	fields         [][]Field
	entrance, exit image.Point
}

func (self *boardImpl) Width() int             { return len(self.fields[0]) }
func (self *boardImpl) Height() int            { return len(self.fields) }
func (self *boardImpl) At(x, y int) *Field     { return &self.fields[y][x] }
func (self *boardImpl) Entrance() *image.Point { return &self.entrance }
func (self *boardImpl) Exit() *image.Point     { return &self.exit }

func (self *boardImpl) String() string {
	var buf bytes.Buffer
	for y, row := range self.fields {
		for _, field := range row {
			if field.Direction()&N != 0 {
				buf.WriteString("+ +")
			} else {
				buf.WriteString("+-+")
			}
		}
		buf.WriteString("\n")

		for x, field := range row {
			dir := field.Direction()
			if dir&W != 0 {
				buf.WriteString(" ")
			} else {
				buf.WriteString("|")
			}
			point := image.Pt(x, y)
			switch {
			case point.Eq(*self.Entrance()) && point.Eq(*self.Exit()):
				buf.WriteString("X")
			case point.Eq(*self.Entrance()):
				buf.WriteString("*")
			case point.Eq(*self.Exit()):
				buf.WriteString("x")
			default:
				buf.WriteString(" ")
			}
			if dir&E != 0 {
				buf.WriteString(" ")
			} else {
				buf.WriteString("|")
			}
		}
		buf.WriteString("\n")

		for _, field := range row {
			if field.Direction()&S != 0 {
				buf.WriteString("+ +")
			} else {
				buf.WriteString("+-+")
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func (self *boardImpl) Walk() (visitMatrix [][]bool, error os.Error) {
	visitMatrix = make([][]bool, self.Height())
	for i := range visitMatrix {
		visitMatrix[i] = make([]bool, self.Width())
	}
	error = self.walkInternal(visitMatrix, *self.Entrance())
	return
}

func (self *boardImpl) walkInternal(visitMatrix [][]bool, p image.Point) os.Error {
	if visitMatrix[p.Y][p.X] {
		return nil
	}
	visitMatrix[p.Y][p.X] = true
	boardRectangle := image.Rect(0, 0, self.Width(), self.Height())
	directions := self.At(p.X, p.Y).Direction().Decompose()
	for _, dir := range directions {
		delta, error := dir.Delta()
		if error != nil {
			return error
		}
		p2 := p.Add(delta)
		if p2.In(boardRectangle) {
			error = self.walkInternal(visitMatrix, p2)
			if error != nil {
				return error
			}
		} else if !p.Eq(*self.Entrance()) && !p.Eq(*self.Exit()) {
			return os.NewError("Falling out of the board into " + p2.String())
		}
	}
	return nil
}
