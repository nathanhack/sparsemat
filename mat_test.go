package mat

import (
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		rows, cols int
		data       []int
		expected   [][]int
	}{
		{1, 1, []int{1}, [][]int{{1}}},
		{2, 2, []int{1, 0, 0, 1}, [][]int{{1, 0}, {0, 1}}},
		{2, 2, []int{}, [][]int{{0, 0}, {0, 0}}},
		{2, 2, nil, [][]int{{0, 0}, {0, 0}}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var m *Matrix
			if test.data != nil {
				m = New(test.rows, test.cols, test.data...)
			} else {
				m = New(test.rows, test.cols)
			}

			for i := 0; i < len(test.expected); i++ {
				for j := 0; j < len(test.expected[i]); j++ {
					expected := test.expected[i][j]
					actual := m.At(i, j)
					if actual != expected {
						t.Fatalf("expected %v but found %v", expected, actual)
					}
				}
			}
		})
	}
}

func TestCopy(t *testing.T) {
	tests := []struct {
		rows, cols int
		data       []int
		expected   [][]int
	}{
		{1, 1, []int{1}, [][]int{{1}}},
		{2, 2, []int{1, 0, 0, 1}, [][]int{{1, 0}, {0, 1}}},
		{2, 2, []int{}, [][]int{{0, 0}, {0, 0}}},
		{2, 2, nil, [][]int{{0, 0}, {0, 0}}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var m1 *Matrix
			if test.data != nil {
				m1 = New(test.rows, test.cols, test.data...)
			} else {
				m1 = New(test.rows, test.cols)
			}

			m := Copy(m1)

			for i := 0; i < len(test.expected); i++ {
				for j := 0; j < len(test.expected[i]); j++ {
					expected := test.expected[i][j]
					actual := m.At(i, j)
					if actual != expected {
						t.Fatalf("expected %v but found %v", expected, actual)
					}
				}
			}
		})
	}
}

func TestMatrix_Dim(t *testing.T) {
	tests := []struct {
		m            *Matrix
		expectedRows int
		expectedCols int
	}{
		{New(5, 5), 5, 5},
		{New(5, 5).Slice(1, 1, 4, 4), 4, 4},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			rows, cols := test.m.Dim()

			if rows != test.expectedRows {
				t.Fatalf("expected %v but found %v", test.expectedRows, rows)
			}
			if cols != test.expectedCols {
				t.Fatalf("expected %v but found %v", test.expectedCols, cols)
			}
		})
	}
}

func TestMatrix_Slice(t *testing.T) {
	tests := []struct {
		sliced   *Matrix
		expected *Matrix
	}{
		{New(2, 2, 1, 0, 0, 1).Slice(0, 0, 2, 1), New(2, 1, 1, 0)},
		{Identity(8).Slice(3, 0, 4, 4), New(4, 4, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)},
		{Identity(8).Slice(3, 0, 4, 4).T(), New(4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.sliced.Equals(test.expected) {
				t.Fatalf("expected equality between %v and %v", test.sliced, test.expected)
			}

		})
	}
}

func TestMatrix_Slice2(t *testing.T) {

	orignal := Identity(5)

	slice := orignal.Slice(1, 1, 3, 3).T()

	slice.Zeroize()

	slice.Set(1, 0, 1)
	slice.Set(2, 0, 1)

	expectedSlice := New(3, 3, 0, 0, 0, 1, 0, 0, 1, 0, 0)
	expectedOriginal := New(5, 5, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1)

	if !slice.Equals(expectedSlice) {
		t.Fatalf("expcted \n%v\n but found \n%v\n", expectedSlice, slice)
	}
	if !orignal.Equals(expectedOriginal) {
		t.Fatalf("expcted \n%v\n but found \n%v\n", expectedOriginal, orignal)
	}
}

func TestMatrix_Equals(t *testing.T) {
	tests := []struct {
		input1, input2 *Matrix
		expected       bool
	}{
		{Identity(3), Identity(3), true},
		{Identity(3).T(), Identity(3), true},
		{Identity(4), Identity(3), false},
		{Identity(4), nil, false},
		{nil, Identity(4), false},
		{nil, nil, true},
		{New(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0).T(), New(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), false},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.input1.Equals(test.input2)
			if actual != test.expected {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestIdentity(t *testing.T) {
	tests := []struct {
		ident    *Matrix
		expected *Matrix
	}{
		{Identity(3), New(3, 3, 1, 0, 0, 0, 1, 0, 0, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.ident.Equals(test.expected) {
				t.Fatalf("expected equality")
			}
		})
	}
}

func TestMatrix_At(t *testing.T) {
	tests := []struct {
		input    *Matrix
		i, j     int
		expected int
	}{
		{Identity(3), 0, 0, 1},
		{Identity(3), 0, 1, 0},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.input.At(test.i, test.j)
			if actual != test.expected {
				t.Fatalf("expected %v at (%v,%v) but found %v", test.expected, test.i, test.j, actual)
			}
		})
	}
}

func TestMatrix_Mul(t *testing.T) {
	tests := []struct {
		m1, m2, result, expected *Matrix
	}{
		{New(1, 4, 1, 0, 1, 0), New(4, 1, 1, 0, 1, 0), New(1, 1), New(1, 1, 0)},
		{New(1, 4, 1, 0, 1, 0), New(4, 1, 1, 0, 0, 0), New(1, 1), New(1, 1, 1)},
		{New(1, 4, 1, 1, 1, 1), New(4, 1, 1, 1, 1, 0), New(1, 1), New(1, 1, 1)},
		{Identity(3), Identity(3), New(3, 3), Identity(3)},
		{Identity(3), New(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), New(3, 3), New(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0)},
		{New(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), Identity(3), New(3, 3), New(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0)},
		{New(4, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1).T(), Identity(4), New(3, 4), New(4, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1).T()},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Mul(test.m1, test.m2)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected equality")
			}
		})
	}
}

func TestMatrix_Zeroize(t *testing.T) {
	tests := []struct {
		original *Matrix
		expected *Matrix
	}{
		{Identity(3), New(3, 3)},
		{New(3, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1), New(3, 3)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.original.Zeroize()
			if !test.original.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.original)
			}
		})
	}
}

func TestMatrix_T(t *testing.T) {
	tests := []struct {
		original *Matrix
		expected *Matrix
	}{
		{New(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), New(3, 3, 0, 0, 0, 1, 1, 0, 1, 1, 0)},
		{New(4, 2, 0, 1, 0, 0, 0, 0, 1, 0), New(2, 4, 0, 0, 0, 1, 1, 0, 0, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.original.T().Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.original.T())
			}
		})
	}
}

func TestMatrix_Add(t *testing.T) {
	tests := []struct {
		a, b, result *Matrix
		expected     *Matrix
	}{
		{Identity(3), Identity(3), New(3, 3), New(3, 3)},
		{Identity(3), New(3, 3), New(3, 3), Identity(3)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Add(test.a, test.b)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.result)
			}
		})
	}
}

func TestMatrix_Column(t *testing.T) {
	tests := []struct {
		m        *Matrix
		j        int //column
		expected map[int]int
	}{
		{Identity(3), 1, map[int]int{1: 1}},
		{Identity(3), 0, map[int]int{0: 1}},
		{New(4, 4, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0).Slice(1, 1, 2, 2).T(), 0, map[int]int{1: 1}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.m.Column(test.j)

			if len(actual) != len(test.expected) {
				t.Fatalf("expected %v values found %v: %v", len(test.expected), len(actual), actual)
			}

			for row, v1 := range test.expected {
				v2, ok := actual[row]
				if !ok {
					t.Fatalf("expected to find row %v", row)
				}
				if v1 != v2 {
					t.Fatalf("expected value at row %v to be %v but found %v", row, v1, v2)
				}
			}
		})
	}
}

func TestMatrix_Row(t *testing.T) {
	tests := []struct {
		m        *Matrix
		i        int //row index
		expected map[int]int
	}{
		{Identity(3), 1, map[int]int{1: 1}},
		{Identity(3), 0, map[int]int{0: 1}},
		{New(4, 4, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0).Slice(1, 1, 2, 2).T(), 1, map[int]int{0: 1}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.m.Row(test.i)

			if len(actual) != len(test.expected) {
				t.Fatalf("expected %v values found %v: %v", len(test.expected), len(actual), actual)
			}

			for col, v1 := range test.expected {
				v2, ok := actual[col]
				if !ok {
					t.Fatalf("expected to find col %v", col)
				}
				if v1 != v2 {
					t.Fatalf("expected value at col %v to be %v but found %v", col, v1, v2)
				}
			}
		})
	}
}

func TestMatrix_SetColumn(t *testing.T) {
	tests := []struct {
		m        *Matrix
		j        int //column to change
		values   map[int]int
		expected *Matrix
	}{
		{Identity(3), 0, map[int]int{1: 1}, New(3, 3, 0, 0, 0, 1, 1, 0, 0, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.m.SetColumn(test.j, test.values)
			if !test.m.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.m)
			}
		})
	}
}

func TestMatrix_SetRow(t *testing.T) {
	tests := []struct {
		m        *Matrix
		i        int //row to change
		values   map[int]int
		expected *Matrix
	}{
		{Identity(3), 0, map[int]int{1: 1}, New(3, 3, 0, 1, 0, 0, 1, 0, 0, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.m.SetRow(test.i, test.values)
			if !test.m.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.m)
			}
		})
	}
}
