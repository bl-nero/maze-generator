package board

import (
	"image"
	"strings"
	"testing"
	"testutil"
)

func TestDirectionNames(t *testing.T) {
	testCases := map[Direction]string{
		None:                          "None",
		S:                             "S",
		N | W:                         "NW",
		N | E | S | W:                 "NESW",
		Direction(directionMask) << 1: "(illegal)"}
	for dir, name := range testCases {
		if dir.String() != name {
			t.Errorf("Name of direction %d is %q, expected %q",
				uint8(dir), dir.String(), name)
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

var dirNegationTests [][]Direction = [][]Direction{
	{None, N | E | S | W}, {N, E | S | W}, {W | S, N | E},
}

func TestDirectionNegation(t *testing.T) {
	for _, test := range dirNegationTests {
		neg := test[0].Negate()
		if neg != test[1] {
			t.Errorf("Negation of %v is %v, expected %v",
				test[0], neg, test[1])
		}
	}
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
	if !board.Entrance().Eq(image.Pt(0, 0)) {
		t.Errorf("Entrance is %v, expected %v", board.Entrance(), image.Pt(0, 0))
	}
	expectedExit := image.Pt(width-1, height-1)
	if !board.Exit().Eq(expectedExit) {
		t.Errorf("Exit is %v, expected %v", board.Exit(), expectedExit)
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
		if !test.Board.Validate() {
			t.Fatalf("Test %d is broken:\n%v", i, &test.Board)
		}
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
	if !board.Validate() {
		t.Fatal("Test is broken")
	}
	_, error := board.Walk()
	expectedPointStr := image.Pt(1, -1).String()
	if error == nil || !strings.Contains(error.String(), expectedPointStr) {
		t.Errorf("Error is %q, expected to contain %s", error, expectedPointStr)
	}
}

type validationTest struct {
	Fields [][]Field
	Ok     bool
}

var validationTests []validationTest = []validationTest{
	// +-++-+
	// |    |
	// + ++-+
	// + ++-+
	// | || |
	// +-++-+
	{
		Fields: [][]Field{
			{Field(E | S), Field(W)},
			{Field(N), Field(None)},
		},
		Ok: true,
	},

	// +-++-+
	// | |  |
	// +-++-+
	{
		Fields: [][]Field{
			{Field(None), Field(W)},
		},
		Ok: false,
	},

	// +-+
	// | |
	// +-+
	// + +
	// | |
	// +-+
	{
		Fields: [][]Field{
			{Field(None)},
			{Field(N)},
		},
		Ok: false,
	},
}

func TestValidation(t *testing.T) {
	for i, test := range validationTests {
		board := boardImpl{fields: test.Fields}
		validated := board.Validate()
		if validated != test.Ok {
			t.Errorf("Validation %d resulted in %v, expected %v",
				i, validated, test.Ok)
		}
	}
}

type complexityTest struct {
	Fields     [][]Field
	Complexity int
}

var complexityTests []complexityTest = []complexityTest{
	// +-+-+-+
	// |   | |
	// + + + +
	// | |   |
	// +-+-+-+
	{
		Fields: [][]Field{
			{Field(E | S), Field(W | S), Field(S)},
			{Field(N), Field(E | N), Field(W | N)},
		},
		Complexity: 0,
	},

	// +-+-+-+-+-+
	// | | |     |
	// + + + + +-+
	// |     |   |
	// +-+ +-+-+-+
	// |         |
	// +---------+
	{
		Fields: [][]Field{
			{Field(S), Field(S), Field(S | E), Field(S | E | W), Field(W)},
			{Field(N | E), Field(N | E | S | W), Field(N | W), Field(N | E), Field(W)},
			{Field(E), Field(N | E | W), Field(E | W), Field(E | W), Field(W)},
		},
		Complexity: 3,
	},
}

func TestComplexity(t *testing.T) {
	for i, test := range complexityTests {
		board := boardImpl{fields: test.Fields}
		if !board.Validate() {
			t.Fatalf("Test %d is broken:\n%v", i, &board)
		}
		complexity := board.Complexity()
		if complexity != test.Complexity {
			t.Errorf("Complexity of test %d is %d, expected %d",
				i, complexity, test.Complexity)
		}
	}
}
