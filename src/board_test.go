// Copyright 2012 Google Inc. All Rights Reserved.
// Author: bleper@google.com (Bartosz Leper)

package board

import (
	"testing"
	"image"
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
	// | |
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
	// | | |
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
	//     |
	// +-+ +
	// |X|
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
	//
	// +-+-+
	// |X X|
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
}

func TestWalking(t *testing.T) {
	for i, test := range walkingTests {
		visitMatrix := test.Board.Walk()
		if !matricesEqual(visitMatrix, test.VisitMatrix) {
			t.Errorf("Visit matrix for test %d is %v, expected %v",
				i, visitMatrix, test.VisitMatrix)
		}
	}
}

func matricesEqual(m1, m2 [][]bool) bool {
	if len(m1) != len(m2) {
		return false
	}
	for i, row1 := range m1 {
		row2 := m2[i]
		if len(row1) != len(row2) {
			return false
		}
		for j, item1 := range row1 {
			if item1 != row2[j] {
				return false
			}
		}
	}
	return true
}
