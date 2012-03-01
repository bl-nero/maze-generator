package board

import (
	"image"
	"strings"
	"testing"
	"testutil"
)

func TestDirectionNames(t *testing.T) {
	testCases := map[Direction]string{
		None:          "None",
		S:             "S",
		N | W:         "NW",
		N | E | S | W: "NESW"}
	for dir, name := range testCases {
		if dir.String() != name {
			t.Errorf("Name of direction %d is %q, expected %q",
				uint8(dir), dir.String(), name)
		}
	}
}

func TestDirectionDecomposing(t *testing.T) {
	testCases := map[Direction][]Direction{
		None:          {},
		N:             {N},
		S:             {S},
		E:             {E},
		W:             {W},
		S | E:         {E, S},
		N | E | S | W: {N, E, S, W},
	}
	for dir, expected := range testCases {
		decomposition := dir.Decompose()
		if !dirArrayEquals(decomposition, expected) {
			t.Errorf("Decomposition of %v is %v, expected %v",
				dir, decomposition, expected)
		}
	}
}

func dirArrayEquals(a1, a2 []Direction) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i := range a1 {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

func performFieldValueTests(t *testing.T, f *Field, visited bool) {
	f.setVisited(visited)
	f.SetDirection(N)
	if f.Direction() != N {
		t.Errorf("Direction set to %v, expected N", f.Direction())
	}
	f.SetDirection(S)
	if f.Direction() != S {
		t.Errorf("Direction set to %v, expected S", f.Direction())
	}
	f.AddDirection(W)
	if f.Direction() != S|W {
		t.Errorf("Direction set to %v, expected %v", f.Direction(), S|W)
	}
	if f.visited() != visited {
		t.Errorf("Visited flag is %v, expected %v", f.visited(), visited)
	}
	if f.Direction() != S|W {
		t.Errorf("Direction after changing visited flag is %v, expected %v",
			f.Direction(), S|W)
	}
}

func TestFieldValues(t *testing.T) {
	var f Field
	if f.Direction() != None {
		t.Errorf("Initial direction is %v, expected None", f.Direction())
	}
	if f.visited() {
		t.Errorf("Field is initially visited")
	}
	performFieldValueTests(t, &f, false)
	performFieldValueTests(t, &f, true)
}

func TestCreatingBoard(t *testing.T) {
	const width, height = 3, 2
	board := New(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			field := board.At(x, y)
			if field.Direction() != None {
				t.Errorf("Direction of (%d, %d) is %v", x, y,
					field.Direction())
			}
			if field.visited() {
				t.Errorf("Field (%d, %d) is visited", x, y)
			}
		}
	}
}

type walkingTest struct {
	Board       boardImpl
	VisitMatrix [][]bool
}

var walkingTests []walkingTest = []walkingTest{

	// + +
	// |X|
	// + +
	{
		Board: boardImpl{
			fields:   [][]Field{{Field(N | S)}},
			entrance: image.Pt(0, 0),
			exit:     image.Pt(0, 0),
		},
		VisitMatrix: [][]bool{{true}},
	},

	// +-+-+
	// |   |
	// + + +
	// |*|x|
	// + + +
	{
		Board: boardImpl{
			fields: [][]Field{
				{Field(E | S), Field(W | S)},
				{Field(N | S), Field(N | S)},
			},
			entrance: image.Pt(0, 1),
			exit:     image.Pt(1, 1),
		},
		VisitMatrix: [][]bool{{true, true}, {true, true}},
	},

	// +-+-+
	//  *  |
	// +-+ +
	// |#|x
	// +-+-+
	{
		Board: boardImpl{
			fields: [][]Field{
				{Field(E | W), Field(W | S)},
				{Field(None), Field(N | E)},
			},
			entrance: image.Pt(0, 0),
			exit:     image.Pt(1, 1),
		},
		VisitMatrix: [][]bool{{true, true}, {false, true}},
	},

	// +-+-+
	//  * x
	// +-+-+
	// |# #|
	// +-+-+
	{
		Board: boardImpl{
			fields: [][]Field{
				{Field(E | W), Field(E | W)},
				{Field(E), Field(W)},
			},
			entrance: image.Pt(0, 0),
			exit:     image.Pt(1, 0),
		},
		VisitMatrix: [][]bool{{true, true}, {false, false}},
	},

	// + +
	// |*|
	// + +
	// |x|
	// + +
	{
		Board: boardImpl{
			fields: [][]Field{
				{Field(N | S)},
				{Field(N | S)},
			},
			entrance: image.Pt(0, 0),
			exit:     image.Pt(0, 1),
		},
		VisitMatrix: [][]bool{{true}, {true}},
	},

	// +-+-+
	//  x *
	// +-+-+
	{
		Board: boardImpl{
			fields: [][]Field{
				{Field(E | W), Field(E | W)},
			},
			entrance: image.Pt(1, 0),
			exit:     image.Pt(0, 0),
		},
		VisitMatrix: [][]bool{{true, true}},
	},

	// +-+-+-+
	// |     |
	// + + + +
	// |x|*| |
	// + + +-+
	{
		Board: boardImpl{
			fields: [][]Field{
				{Field(E | S), Field(E | S | W), Field(S | W)},
				{Field(N | S), Field(N | S), Field(N)},
			},
			entrance: image.Pt(1, 1),
			exit:     image.Pt(0, 1),
		},
		VisitMatrix: [][]bool{{true, true, true}, {true, true, true}},
	},
}

func TestWalking(t *testing.T) {
	for i, test := range walkingTests {
		visitMatrix, error := test.Board.Walk()
		if error != nil {
			t.Errorf("Error in test %d: %v", i, error)
			continue
		}
		if visitMatrix == nil {
			t.Errorf("Visit matrix for test %d is nil. "+
				"Something terrible has happened", i)
			continue
		}
		if !testutil.MatricesEqual(visitMatrix, test.VisitMatrix) {
			t.Errorf("Visit matrix for test %d is %v, expected %v",
				i, visitMatrix, test.VisitMatrix)
		}
	}
}

func TestWalkingFallsOffBoard(t *testing.T) {
	// +-+ +-+
	//  *
	// +-+-+-+
	board := boardImpl{
		fields: [][]Field{
			{Field(E | W), Field(N | E | W), Field(E | W)},
		},
		entrance: image.Pt(0, 0),
		exit:     image.Pt(2, 0),
	}
	_, error := board.Walk()
	expectedPointStr := image.Pt(1, -1).String()
	if error == nil || !strings.Contains(error.String(), expectedPointStr) {
		t.Errorf("Error is %q, expected to contain %s", error, expectedPointStr)
	}
}
