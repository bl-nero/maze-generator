package board

type Direction uint8

const (
	None Direction = 0
	N    Direction = 1 << iota
	E
	S
	W
	maxDirection  = W
	directionMask = maxDirection<<1 - 1
)

type state uint8

const (
	normal state = state(maxDirection) << (1 + iota)
	visited
)

type Field uint8

type Board struct {
	fields [][]Direction
}

func New(width, height int) Board {
	board := Board{make([][]Direction, height)}
	for i := range board.fields {
		board.fields[i] = make([]Direction, width)
	}
	return board
}
