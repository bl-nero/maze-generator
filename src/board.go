package board

import "image"

type Direction uint8

const (
	None Direction = 0
	N    Direction = 1 << iota
	E
	S
	W
	visitedBit    Field = 1 << iota
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

func (self Direction) Delta() (delta image.Point, ok bool) {
	delta, ok = dirDeltas[self]
	return
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
	Walk() [][]bool
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

func (self *boardImpl) Walk() [][]bool {
	visitMatrix := make([][]bool, self.Height())
	for i := range visitMatrix {
		visitMatrix[i] = make([]bool, self.Width())
	}

	p := *self.Entrance();
	for !p.Eq(*self.Exit()) {
		visitMatrix[p.Y][p.X] = true
		delta, ok := self.At(p.X, p.Y).Direction().Delta()
		if (!ok) {
			return nil
		}
		p = p.Add(delta)
	}
	visitMatrix[p.Y][p.X] = true

	/*for y, row := range visitMatrix {
		for x := range row {
			row[x] = self.At(x, y).Direction() != None
		}
	}*/
	return visitMatrix
}
