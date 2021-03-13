package sparsemat

import (
	"encoding/json"
	"strconv"
	"testing"
)

func TestDOKMat(t *testing.T) {
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
			var m SparseMat
			if test.data != nil {
				m = DOKMat(test.rows, test.cols, test.data...)
			} else {
				m = DOKMat(test.rows, test.cols)
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

func TestDOKMatCopy(t *testing.T) {
	tests := []struct {
		mat SparseMat
	}{
		{DOKMat(5, 5)},
		{DOKIdentity(5)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := DOKMatCopy(test.mat)

			if !actual.Equals(test.mat) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.mat, actual)
			}
		})
	}
}

func TestDOKMatrix_Dim(t *testing.T) {
	tests := []struct {
		m            SparseMat
		expectedRows int
		expectedCols int
	}{
		{DOKMat(5, 5), 5, 5},
		{DOKMat(5, 5).Slice(1, 1, 4, 4), 4, 4},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			rows, cols := test.m.Dims()

			if rows != test.expectedRows {
				t.Fatalf("expected %v but found %v", test.expectedRows, rows)
			}
			if cols != test.expectedCols {
				t.Fatalf("expected %v but found %v", test.expectedCols, cols)
			}
		})
	}
}

func TestDOKMatrix_Slice(t *testing.T) {
	tests := []struct {
		sliced   SparseMat
		expected SparseMat
	}{
		{DOKMat(2, 2, 1, 0, 0, 1).Slice(0, 0, 2, 1), DOKMat(2, 1, 1, 0)},
		{DOKIdentity(8).Slice(3, 0, 4, 4), DOKMat(4, 4, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)},
		{DOKIdentity(8).Slice(3, 0, 4, 4).T(), DOKMat(4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.sliced.Equals(test.expected) {
				t.Fatalf("expected equality between %v and %v", test.sliced, test.expected)
			}

		})
	}
}

func TestDOKMatrix_Equals(t *testing.T) {
	tests := []struct {
		input1, input2 SparseMat
		expected       bool
	}{
		{DOKIdentity(3), DOKIdentity(3), true},
		{DOKIdentity(3).T(), DOKIdentity(3), true},
		{DOKIdentity(4), DOKIdentity(3), false},
		{DOKIdentity(4), nil, false},
		{DOKMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0).T(), DOKMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), false},
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

func TestDOKIdentity(t *testing.T) {
	tests := []struct {
		ident    SparseMat
		expected SparseMat
	}{
		{DOKIdentity(3), DOKMat(3, 3, 1, 0, 0, 0, 1, 0, 0, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.ident.Equals(test.expected) {
				t.Fatalf("expected equality")
			}
		})
	}
}

func TestDOKMatrix_At(t *testing.T) {
	tests := []struct {
		input    SparseMat
		i, j     int
		expected int
	}{
		{DOKIdentity(3), 0, 0, 1},
		{DOKIdentity(3), 0, 1, 0},
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

func TestDOKMatrix_Mul(t *testing.T) {
	tests := []struct {
		m1, m2, result, expected SparseMat
	}{
		{DOKMat(1, 4, 1, 0, 1, 0), DOKMat(4, 1, 1, 0, 1, 0), DOKMat(1, 1), DOKMat(1, 1, 0)},
		{DOKMat(1, 4, 1, 0, 1, 0), DOKMat(4, 1, 1, 0, 0, 0), DOKMat(1, 1), DOKMat(1, 1, 1)},
		{DOKMat(1, 4, 1, 1, 1, 1), DOKMat(4, 1, 1, 1, 1, 0), DOKMat(1, 1), DOKMat(1, 1, 1)},
		{DOKIdentity(3), DOKIdentity(3), DOKMat(3, 3), DOKIdentity(3)},
		{DOKIdentity(3), DOKMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), DOKMat(3, 3), DOKMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0)},
		{DOKMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), DOKIdentity(3), DOKMat(3, 3), DOKMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0)},
		{DOKMat(4, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1).T(), DOKIdentity(4), DOKMat(3, 4), DOKMat(4, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1).T()},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Mul(test.m1, test.m2)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestDOKMatrix_Zeroize(t *testing.T) {
	tests := []struct {
		original SparseMat
		expected SparseMat
	}{
		{DOKIdentity(3), DOKMat(3, 3)},
		{DOKMat(3, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1), DOKMat(3, 3)},
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

func TestDOKMatrix_ZeroizeRange(t *testing.T) {
	tests := []struct {
		original         SparseMat
		i, j, rows, cols int
		expected         SparseMat
	}{
		{DOKMat(4, 4, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), 1, 1, 2, 2, DOKMat(4, 4, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.original.ZeroizeRange(test.i, test.j, test.rows, test.cols)
			if !test.original.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.original)
			}
		})
	}
}

func TestDOKMatrix_T(t *testing.T) {
	tests := []struct {
		original SparseMat
		expected SparseMat
	}{
		{DOKMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), DOKMat(3, 3, 0, 0, 0, 1, 1, 0, 1, 1, 0)},
		{DOKMat(4, 2, 0, 1, 0, 0, 0, 0, 1, 0), DOKMat(2, 4, 0, 0, 0, 1, 1, 0, 0, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.original.T().Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.original.T())
			}
		})
	}
}

func TestDOKMatrix_Add(t *testing.T) {
	tests := []struct {
		a, b, result SparseMat
		expected     SparseMat
	}{
		{DOKIdentity(3), DOKIdentity(3), DOKMat(3, 3), DOKMat(3, 3)},
		{DOKIdentity(3), DOKMat(3, 3), DOKMat(3, 3), DOKIdentity(3)},
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

func TestDOKMatrix_Add2(t *testing.T) {
	tests := []struct {
		original         SparseMat
		i, j, rows, cols int
		addToSlice       SparseMat
		expectedOriginal SparseMat
		expectedSlice    SparseMat
	}{
		{
			DOKMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1),
			1, 1, 3, 3,
			DOKIdentity(3),
			DOKMat(5, 5, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1),
			DOKMat(3, 3, 0, 1, 1, 1, 0, 1, 1, 1, 0),
		},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i, test.j, test.rows, test.cols)

			sl.Add(sl, test.addToSlice)
			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}
		})
	}
}

func TestDOKMatrix_Column(t *testing.T) {
	tests := []struct {
		m        SparseMat
		j        int //column
		expected SparseVector
	}{
		{DOKIdentity(3), 1, DOKVec(3, 0, 1, 0)},
		{DOKIdentity(3), 0, DOKVec(3, 1, 0, 0)},
		{DOKMat(4, 4, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0).Slice(1, 1, 2, 2).T(), 0, DOKVec(2, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.m.Column(test.j)

			if !actual.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestDOKMatrix_Row(t *testing.T) {
	tests := []struct {
		m        SparseMat
		i        int //row index
		expected SparseVector
	}{
		{DOKIdentity(3), 1, DOKVec(3, 0, 1, 0)},
		{DOKIdentity(3), 0, DOKVec(3, 1, 0, 0)},
		{DOKMat(4, 4, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0).Slice(1, 1, 2, 2).T(), 1, DOKVec(2, 1, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.m.Row(test.i)

			if !actual.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestR(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: test cases
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}

func TestDOKMatrix_SetColumn(t *testing.T) {
	tests := []struct {
		m        SparseMat
		j        int //column to change
		vec      SparseVector
		expected SparseMat
	}{
		{DOKIdentity(3), 0, DOKVec(3, 0, 1, 0), DOKMat(3, 3, 0, 0, 0, 1, 1, 0, 0, 0, 1)},
		{DOKIdentity(3), 1, DOKIdentity(3).Column(2), DOKMat(3, 3, 1, 0, 0, 0, 0, 0, 0, 1, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.m.SetColumn(test.j, test.vec)
			if !test.m.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.m)
			}
		})
	}
}

func TestDOKMatrix_SetRow(t *testing.T) {
	tests := []struct {
		m        SparseMat
		i        int //row to change
		vec      SparseVector
		expected SparseMat
	}{
		{DOKIdentity(3), 0, DOKVec(3, 0, 1, 0), DOKMat(3, 3, 0, 1, 0, 0, 1, 0, 0, 0, 1)},
		{DOKIdentity(3), 1, DOKIdentity(3).Row(2), DOKMat(3, 3, 1, 0, 0, 0, 0, 1, 0, 0, 1)},
		{DOKMat(4, 3), 0, DOKVec(3, 1, 1, 1), DOKMat(4, 3, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0)},
		{DOKMat(3, 4), 0, DOKVec(4, 1, 1, 1, 1), DOKMat(3, 4, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0)},
		{DOKMat(4, 3), 3, DOKVec(3, 1, 1, 1), DOKMat(4, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1)},
		{DOKMat(3, 4), 2, DOKVec(4, 1, 1, 1, 1), DOKMat(3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.m.SetRow(test.i, test.vec)
			if !test.m.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.m)
			}
		})
	}
}

func TestDOKMatrix_SetMatrix(t *testing.T) {
	tests := []struct {
		dest             SparseMat
		source           SparseMat
		iOffset, jOffset int
		expected         SparseMat
	}{
		{DOKMat(3, 3), DOKIdentity(3), 0, 0, DOKIdentity(3)},
		{DOKMat(4, 4), DOKMat(2, 2, 1, 1, 1, 1), 1, 1, DOKMat(4, 4, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0)},
		{DOKMat(4, 4), DOKMat(2, 2, 0, 1, 0, 0).T(), 1, 1, DOKMat(4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0)},
		{DOKIdentity(4), DOKMat(2, 2, 1, 1, 1, 1), 1, 1, DOKMat(4, 4, 1, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 1)},
		{DOKMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), DOKIdentity(3), 1, 1, DOKMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.dest.SetMatrix(test.source, test.iOffset, test.jOffset)
			if !test.dest.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.dest)
			}
		})
	}
}

func TestDOKMatrix_SetMatrix2(t *testing.T) {
	tests := []struct {
		original         SparseMat
		i, j, rows, cols int
		source           SparseMat
		iOffset, jOffset int
		expectedOriginal SparseMat
		expectedSlice    SparseMat
	}{
		{DOKMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1),
			1, 1, 3, 3,
			DOKIdentity(2),
			1, 1,
			DOKMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1),
			DOKMat(3, 3, 1, 1, 1, 1, 1, 0, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i, test.j, test.rows, test.cols)
			sl.SetMatrix(test.source, test.iOffset, test.jOffset)
			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}
		})
	}
}

func TestDOKMatrix_JSON(t *testing.T) {
	m := DOKIdentity(3)

	bs, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}

	var actual DOKMatrix
	err = json.Unmarshal(bs, &actual)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}
	if !m.Equals(&actual) {
		t.Fatalf("expected %v but found %v", m, actual)
	}
}

func TestDOKMatrix_And(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseMat
	}{
		{DOKMat(2, 2, 0, 1, 0, 1), DOKMat(2, 2, 0, 0, 1, 1), DOKMat(2, 2), DOKMat(2, 2, 0, 0, 0, 1)},
		{DOKMat(2, 2, 0, 0, 1, 1), DOKMat(2, 2, 0, 1, 0, 1), DOKMat(2, 2), DOKMat(2, 2, 0, 0, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.And(test.x, test.y)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestDOKMatrix_Or(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseMat
	}{
		{DOKMat(2, 2, 0, 1, 0, 1), DOKMat(2, 2, 0, 0, 1, 1), DOKMat(2, 2), DOKMat(2, 2, 0, 1, 1, 1)},
		{DOKMat(2, 2, 0, 0, 1, 1), DOKMat(2, 2, 0, 1, 0, 1), DOKMat(2, 2), DOKMat(2, 2, 0, 1, 1, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Or(test.x, test.y)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestDOKMatrix_XOr(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseMat
	}{
		{DOKMat(2, 2, 0, 1, 0, 1), DOKMat(2, 2, 0, 0, 1, 1), DOKMat(2, 2), DOKMat(2, 2, 0, 1, 1, 0)},
		{DOKMat(2, 2, 0, 0, 1, 1), DOKMat(2, 2, 0, 1, 0, 1), DOKMat(2, 2), DOKMat(2, 2, 0, 1, 1, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.XOr(test.x, test.y)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestDOKMatrix_Negate(t *testing.T) {
	tests := []struct {
		x, expected SparseMat
	}{
		{DOKMat(2, 2, 0, 1, 0, 1), DOKMat(2, 2, 1, 0, 1, 0)},
		{DOKMat(2, 2, 0, 1, 1, 0), DOKMat(2, 2, 1, 0, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.x.Negate()

			if !test.x.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.x)
			}
		})
	}
}

func TestDOKMatrix_SwapRows(t *testing.T) {
	tests := []struct {
		input    SparseMat
		a, b     int
		expected SparseMat
	}{
		{DOKIdentity(4), 0, 0, DOKMat(4, 4, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1)},
		{DOKIdentity(4), 1, 3, DOKMat(4, 4, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0)},
		{DOKMat(4, 5, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1), 2, 3, DOKMat(4, 5, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0)},
		{DOKMat(5, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1), 0, 1, DOKMat(5, 4, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1)},
		{DOKMat(5, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1), 1, 2, DOKMat(5, 4, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1)},
		{DOKMat(5, 4, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0), 4, 0, DOKMat(5, 4, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1)},
		{DOKMat(5, 4, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0), 0, 1, DOKMat(5, 4, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0)},
		{DOKMat(5, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1), 0, 2, DOKMat(5, 4, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.input.SwapRows(test.a, test.b)
			if !test.input.Equals(test.expected) {
				t.Fatalf("after rowswap(%v <> %v) expected \n%v\n but found \n%v\n", test.a, test.b, test.expected, test.input)
			}
		})
	}
}

func BenchmarkDOKMatrix_SwapRows(b *testing.B) {
	a := DOKMat(5, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1)
	aa := DOKMat(5, 4, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0)
	aaa := DOKMat(5, 4, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.SwapRows(0, 1)
		a.SwapRows(0, 1)
		aa.SwapRows(0, 1)
		aa.SwapRows(0, 1)
		aaa.SwapRows(0, 4)
		aaa.SwapRows(0, 4)
	}
}

func TestDOKMatrix_SwapColumns(t *testing.T) {
	tests := []struct {
		input    SparseMat
		a, b     int
		expected SparseMat
	}{
		{DOKIdentity(4), 1, 3, DOKMat(4, 4, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.input.SwapColumns(test.a, test.b)
			if !test.input.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.input)
			}
		})
	}
}

func TestDOKMatrix_AddRows(t *testing.T) {
	tests := []struct {
		input        SparseMat
		i1, i2, dest int
		expected     SparseMat
	}{
		{DOKIdentity(3), 0, 1, 2, DOKMat(3, 3, 1, 0, 0, 0, 1, 0, 1, 1, 0)},
		{DOKMat(4, 5, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1), 1, 2, 1, DOKMat(4, 5, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1)},
		{DOKMat(4, 5, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), 0, 2, 0, DOKMat(4, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)},
		{DOKMat(4, 5, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), 1, 2, 1, DOKMat(4, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)},
		{DOKMat(4, 5, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), 1, 2, 1, DOKMat(4, 5, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.input.AddRows(test.i1, test.i2, test.dest)
			if !test.input.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.input)
			}
		})
	}
}

func BenchmarkAddRow(b *testing.B) {
	a := DOKMat(5, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 1; j < 5; j++ {
			a.AddRows(0, j, j)
			a.AddRows(4, j, 4)
			a.AddRows(4, j, 4)
		}
	}
}

func TestDOKMatFromVec(t *testing.T) {
	tests := []struct {
		vec SparseVector
	}{
		{DOKVec(5, 1, 0, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := DOKMatFromVec(test.vec).Row(0)
			if !actual.Equals(test.vec) {
				t.Fatalf("expected %v but found %v", test.vec, actual)
			}
		})
	}
}

func TestDOKMatrix_Set(t *testing.T) {
	m := DOKMat(5, 5)
	for i := 0; i < 5; i++ {
		m.Set(i, i, 1)
	}

	expected := DOKIdentity(5)

	if !expected.Equals(m) {
		t.Fatalf("expected %v but found %v", expected, m)
	}
}

//func TestDOKMatrix_AddCols(t *testing.T) {
//	tests := []struct {
//		input        SparseMat
//		j1, j2, dest int
//		expected     SparseMat
//	}{
//		{DOKIdentity(3), 0, 1, 2, DOKMat(3, 3, 1, 0, 1, 0, 1, 1, 0, 0, 0)},
//	}
//	for i, test := range tests {
//		t.Run(strconv.Itoa(i), func(t *testing.T) {
//			test.input.AddCols(test.j1, test.j2, test.dest)
//			if !test.input.Equals(test.expected) {
//				t.Fatalf("expected %v but found %v", test.expected, test.input)
//			}
//		})
//	}
//}
