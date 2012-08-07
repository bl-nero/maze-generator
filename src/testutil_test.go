// Copyright 2012 Google Inc. All Rights Reserved.
// Author: bleper@google.com (Bartosz Leper)

package testutil

import (
	"testing"
)

type matrixEqualityTest struct {
	M1, M2 [][]bool
	Equal  bool
}

var matrixEqualityTests []matrixEqualityTest = []matrixEqualityTest{
	{
		M1: [][]bool{
			{true, false},
			{false, true},
		},
		M2: [][]bool{
			{true, false},
			{false, true},
		},
		Equal: true,
	},
	{
		M1:    [][]bool{},
		M2:    [][]bool{},
		Equal: true,
	},
	{
		M1: [][]bool{
			{false, true},
			{false, true},
		},
		M2: [][]bool{
			{false, false},
			{false, true},
		},
		Equal: false,
	},
	{
		M1: [][]bool{
			{true, true},
		},
		M2: [][]bool{
			{true, true},
			{true, true},
		},
		Equal: false,
	},
	{
		M1: [][]bool{
			{true},
			{true},
		},
		M2: [][]bool{
			{true, true},
			{true, true},
		},
		Equal: false,
	},
}

func TestMatrixEuqality(t *testing.T) {
	for i, test := range matrixEqualityTests {
		eq := MatricesEqual(test.M1, test.M2)
		if eq != test.Equal {
			t.Errorf("Result for test %d: %v, expected %v",
				i, eq, test.Equal)
		}
	}
}
