package generator

import (
	"testing"
	"testutil"
)

func TestGenerating(t *testing.T) {
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
		t.Fatalf("Board doesn't validate:\n%v", board);
	}
	visitMatrix, error := board.Walk()
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
