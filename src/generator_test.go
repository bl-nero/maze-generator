package generator

import (
	"container/heap"
	"image"
	"rand"
	"testing"
	"testutil"
)

func TestFieldHeap(t *testing.T) {
	numbers := []int{4, 9, 1, 7, 3, 5, 2, 7}
	expected := []int{1, 2, 3, 4, 5, 7, 7, 9}
	h := new(fieldHeap)
	heap.Init(h)
	for _, n := range numbers {
		heap.Push(h, fieldHeapElement{image.Pt(n*2, n*3), n})
	}
	for i, n := range expected {
		actual := heap.Pop(h).(fieldHeapElement)
		expectedCoords := image.Pt(n*2, n*3)
		if !actual.Coords.Eq(expectedCoords) {
			t.Errorf("Coords of element %d are %v, expected %v",
				i, actual.Coords, expectedCoords)
		}
		if actual.Weight != n {
			t.Errorf("Weight of element %d is %d, expected %d",
				i, actual.Weight, n)
		}
	}
}

func TestGenerating(t *testing.T) {
	rand.Seed(0)
	dump := false
	width, height := 10, 5
	board := Generate(width, height)
	if board.Width() != width {
		t.Errorf("Board width is %d, expected %d", board.Width(), width)
	}
	if board.Height() != height {
		t.Errorf("Board height is %d, expected %d", board.Height(), height)
	}
	if !board.Validate() {
		t.Fatalf("Board doesn't validate:\n%v", board)
	}
	visitMatrix, error := board.Walk(false)
	if error != nil {
		t.Fatalf("Unexpected error: %v. Generated board:\n%v", error, board)
	}
	if !testutil.MatricesEqual(trueMatrix(width, height), visitMatrix) {
		t.Errorf("Visit matrix expected to be filled with true, but was:\n%v",
			visitMatrix)
		dump = true
	}
	if board.Complexity() == 0 {
		t.Errorf("Board does not have crossroads")
		dump = true
	}
	exit := *board.Exit()
	exitRoadCount := len(board.At(exit.X, exit.Y).Direction().Decompose())
	if exitRoadCount != 2 {
		t.Errorf("Number of roads coming into exit is %d, expected 2",
			exitRoadCount)
		dump = true
	}

	if dump {
		t.Logf("Generated board:\n%v", board)
	}
}

func trueMatrix(width, height int) [][]bool {
	matrix := make([][]bool, height)
	for y := range matrix {
		matrix[y] = make([]bool, width)
		for x := range matrix[y] {
			matrix[y][x] = true
		}
	}
	return matrix
}
