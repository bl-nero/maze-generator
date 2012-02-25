package board

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

type Field uint8

func (f Field) Direction() Direction {
	return Direction(f & Field(directionMask))
}

func (f Field) isVisited() bool {
	return f&visitedBit != 0
}

type Point struct{ X, Y int }
type Board struct {
	Fields         [][]Field
	Entrance, Exit Point
}

func New(width, height int) Board {
	board := Board{Fields: make([][]Field, height)}
	for i := range board.Fields {
		board.Fields[i] = make([]Field, width)
	}
	return board
}

func (self *Board) walk() {
}
