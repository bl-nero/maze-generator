// Copyright 2012 Google Inc. All Rights Reserved.
// Author: bleper@google.com (Bartosz Leper)

package board

import (
	"testing"
)

func TestCreatingBoard(t *testing.T) {
	const width, height = 4, 5
	board := New(width, height)
	if board == nil {
		t.Errorf("New board is nil")
	}
	if len(board) != 5 {
		t.Errorf("Length of board is %d, expected 5", len(board))
	}
	for i, row := range board {
		if len(row) != 4 {
			t.Errorf("Length of row %d is %d, expected 4", i, len(row))
		}
	}
}
