// Copyright 2012 Google Inc. All Rights Reserved.
// Author: bleper@google.com (Bartosz Leper)

package board

import (
	"testing"
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

func TestCreatingBoard(t *testing.T) {
	const width, height = 3, 2
	board := New(width, height)
	if len(board.Fields) != height {
		t.Fatalf("Length of board is %d, expected %d", len(board.Fields), height)
	}
	for i, row := range board.Fields {
		if len(row) != width {
			t.Errorf("Length of row %d is %d, expected %d", i, len(row), width)
		}
		for j, field := range row {
			if field.Direction() != None {
				t.Errorf("Direction of (%d, %d) is %v", i, j,
					field.Direction())
			}
			if field.isVisited() {
				t.Errorf("Field (%d, %d) is visited")
			}
		}
	}
}

type walkingTest struct {

}

func TestWalking(t *testing.T) {
	//board := New(
}
