package testutil

func MatricesEqual(m1, m2 [][]bool) bool {
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
