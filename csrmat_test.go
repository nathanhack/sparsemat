package sparsemat

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"testing"
)

func TestCSRMat(t *testing.T) {
	tests := []struct {
		rows, cols int
		data       []int
		expected   [][]int
	}{
		{1, 1, []int{1}, [][]int{{1}}},
		{2, 2, []int{1, 0, 0, 1}, [][]int{{1, 0}, {0, 1}}},
		{2, 2, []int{}, [][]int{{0, 0}, {0, 0}}},
		{2, 2, nil, [][]int{{0, 0}, {0, 0}}},
		{4, 4, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, [][]int{{1, 1, 1, 1}, {1, 1, 1, 1}, {1, 1, 1, 1}, {1, 1, 1, 1}}},
		{4, 4, []int{1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1}, [][]int{{1, 1, 1, 1}, {1, 0, 0, 1}, {1, 0, 0, 1}, {1, 1, 1, 1}}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var m SparseMat
			if test.data != nil {
				m = CSRMat(test.rows, test.cols, test.data...)
			} else {
				m = CSRMat(test.rows, test.cols)
			}

			for i := 0; i < test.rows; i++ {
				for j := 0; j < test.cols; j++ {
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

func TestCSRMatCopy(t *testing.T) {
	tests := []struct {
		mat SparseMat
	}{
		{CSRMat(5, 5)},
		{CSRIdentity(5)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := CSRMatCopy(test.mat)

			if !actual.Equals(test.mat) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.mat, actual)
			}
		})
	}
}

func TestCSRMatrix_Dim(t *testing.T) {
	tests := []struct {
		m            SparseMat
		expectedRows int
		expectedCols int
	}{
		{CSRMat(5, 5), 5, 5},
		{CSRMat(5, 5).Slice(1, 1, 4, 4), 4, 4},
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

func TestCSRMatrix_Slice(t *testing.T) {
	tests := []struct {
		sliced   SparseMat
		expected SparseMat
	}{
		{CSRMat(2, 2, 1, 0, 0, 1).Slice(0, 0, 2, 1), CSRMat(2, 1, 1, 0)},
		{CSRIdentity(8).Slice(3, 0, 4, 4), CSRMat(4, 4, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)},
		{CSRIdentity(8).Slice(3, 0, 4, 4).T(), CSRMat(4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.sliced.Equals(test.expected) {
				t.Fatalf("expected equality between %v and %v", test.sliced, test.expected)
			}

		})
	}
}

func TestCSRMatrix_Equals(t *testing.T) {
	tests := []struct {
		input1, input2 SparseMat
		expected       bool
	}{
		{CSRIdentity(3), CSRIdentity(3), true},
		{CSRIdentity(3).T(), CSRIdentity(3), true},
		{CSRIdentity(4), CSRIdentity(3), false},
		{CSRIdentity(4), nil, false},
		{CSRMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0).T(), CSRMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), false},
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

func TestCSRIdentity(t *testing.T) {
	tests := []struct {
		ident    SparseMat
		expected SparseMat
	}{
		{CSRIdentity(3), CSRMat(3, 3, 1, 0, 0, 0, 1, 0, 0, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.ident.Equals(test.expected) {
				t.Fatalf("expected equality")
			}
		})
	}
}

func TestCSRMatrix_At(t *testing.T) {
	tests := []struct {
		input    SparseMat
		i, j     int
		expected int
	}{
		{CSRIdentity(3), 0, 0, 1},
		{CSRIdentity(3), 0, 1, 0},
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

func TestCSRMatrix_Mul(t *testing.T) {
	tests := []struct {
		m1, m2, result, expected SparseMat
	}{
		{CSRMat(1, 4, 1, 0, 1, 0), CSRMat(4, 1, 1, 0, 1, 0), CSRMat(1, 1), CSRMat(1, 1, 0)},
		{CSRMat(1, 4, 1, 0, 1, 0), CSRMat(4, 1, 1, 0, 0, 0), CSRMat(1, 1), CSRMat(1, 1, 1)},
		{CSRMat(1, 4, 1, 1, 1, 1), CSRMat(4, 1, 1, 1, 1, 0), CSRMat(1, 1), CSRMat(1, 1, 1)},
		{CSRIdentity(3), CSRIdentity(3), CSRMat(3, 3), CSRIdentity(3)},
		{CSRIdentity(3), CSRMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), CSRMat(3, 3), CSRMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0)},
		{CSRMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), CSRIdentity(3), CSRMat(3, 3), CSRMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0)},
		{CSRMat(4, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1).T(), CSRIdentity(4), CSRMat(3, 4), CSRMat(4, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1).T()},
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

func TestCSRMatrix_Zeroize(t *testing.T) {
	tests := []struct {
		original SparseMat
		expected SparseMat
	}{
		{CSRIdentity(3), CSRMat(3, 3)},
		{CSRMat(3, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1), CSRMat(3, 3)},
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

func TestCSRMatrix_ZeroizeRange(t *testing.T) {
	tests := []struct {
		original         SparseMat
		i, j, rows, cols int
		expected         SparseMat
	}{
		{CSRMat(4, 4, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), 1, 1, 2, 2, CSRMat(4, 4, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1)},
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

func TestCSRMatrix_T(t *testing.T) {
	tests := []struct {
		original SparseMat
		expected SparseMat
	}{
		{CSRMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), CSRMat(3, 3, 0, 0, 0, 1, 1, 0, 1, 1, 0)},
		{CSRMat(4, 2, 0, 1, 0, 0, 0, 0, 1, 0), CSRMat(2, 4, 0, 0, 0, 1, 1, 0, 0, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.original.T().Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.original.T())
			}
		})
	}
}

func TestCSRMatrix_Add(t *testing.T) {
	tests := []struct {
		a, b, result SparseMat
		expected     SparseMat
	}{
		{CSRIdentity(3), CSRIdentity(3), CSRMat(3, 3), CSRMat(3, 3)},
		{CSRIdentity(3), CSRMat(3, 3), CSRMat(3, 3), CSRIdentity(3)},
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

func TestCSRMatrix_Add2(t *testing.T) {
	tests := []struct {
		original         SparseMat
		i, j, rows, cols int
		addToSlice       SparseMat
		expectedOriginal SparseMat
		expectedSlice    SparseMat
	}{
		{
			CSRMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1),
			1, 1, 3, 3,
			CSRIdentity(3),
			CSRMat(5, 5, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1),
			CSRMat(3, 3, 0, 1, 1, 1, 0, 1, 1, 1, 0),
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

func TestCSRMatrix_Column(t *testing.T) {
	tests := []struct {
		m        SparseMat
		j        int //column
		expected SparseVector
	}{
		{CSRIdentity(3), 1, CSRVec(3, 0, 1, 0)},
		{CSRIdentity(3), 0, CSRVec(3, 1, 0, 0)},
		{CSRMat(4, 4, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0).Slice(1, 1, 2, 2).T(), 0, CSRVec(2, 0, 1)},
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

func TestCSRMatrix_Row(t *testing.T) {
	tests := []struct {
		m        SparseMat
		i        int //row index
		expected SparseVector
	}{
		{CSRIdentity(3), 1, CSRVec(3, 0, 1, 0)},
		{CSRIdentity(3), 0, CSRVec(3, 1, 0, 0)},
		{CSRMat(4, 4, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0).Slice(1, 1, 2, 2).T(), 1, CSRVec(2, 1, 0)},
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

func TestCSRMatrix_SetColumn(t *testing.T) {
	tests := []struct {
		m        SparseMat
		j        int //column to change
		vec      SparseVector
		expected SparseMat
	}{
		{CSRIdentity(3), 0, CSRVec(3, 0, 1, 0), CSRMat(3, 3, 0, 0, 0, 1, 1, 0, 0, 0, 1)},
		{CSRIdentity(3), 1, CSRIdentity(3).Column(2), CSRMat(3, 3, 1, 0, 0, 0, 0, 0, 0, 1, 1)},
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

func TestCSRMatrix_SetRow(t *testing.T) {
	tests := []struct {
		m        SparseMat
		i        int //row to change
		vec      SparseVector
		expected SparseMat
	}{
		{CSRIdentity(3), 0, CSRVec(3, 0, 1, 0), CSRMat(3, 3, 0, 1, 0, 0, 1, 0, 0, 0, 1)},
		{CSRIdentity(3), 1, CSRIdentity(3).Row(2), CSRMat(3, 3, 1, 0, 0, 0, 0, 1, 0, 0, 1)},
		{CSRMat(4, 3), 0, CSRVec(3, 1, 1, 1), CSRMat(4, 3, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0)},
		{CSRMat(3, 4), 0, CSRVec(4, 1, 1, 1, 1), CSRMat(3, 4, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0)},
		{CSRMat(4, 3), 3, CSRVec(3, 1, 1, 1), CSRMat(4, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1)},
		{CSRMat(3, 4), 2, CSRVec(4, 1, 1, 1, 1), CSRMat(3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1)},
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

func TestCSRMatrix_SetMatrix(t *testing.T) {
	tests := []struct {
		dest             SparseMat
		source           SparseMat
		iOffset, jOffset int
		expected         SparseMat
	}{
		{CSRMat(3, 3), CSRIdentity(3), 0, 0, CSRIdentity(3)},
		{CSRMat(4, 4), CSRMat(2, 2, 1, 1, 1, 1), 1, 1, CSRMat(4, 4, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0)},
		{CSRMat(4, 4), CSRMat(2, 2, 0, 1, 0, 0).T(), 1, 1, CSRMat(4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0)},
		{CSRIdentity(4), CSRMat(2, 2, 1, 1, 1, 1), 1, 1, CSRMat(4, 4, 1, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 1)},
		{CSRMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), CSRIdentity(3), 1, 1, CSRMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1)},
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

func TestCSRMatrix_SetMatrix2(t *testing.T) {
	tests := []struct {
		original         SparseMat
		i, j, rows, cols int
		source           SparseMat
		iOffset, jOffset int
		expectedOriginal SparseMat
		expectedSlice    SparseMat
	}{
		{CSRMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1),
			1, 1, 3, 3,
			CSRIdentity(2),
			1, 1,
			CSRMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1),
			CSRMat(3, 3, 1, 1, 1, 1, 1, 0, 1, 0, 1)},
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

func TestCSRMatrix_JSON(t *testing.T) {
	m := CSRIdentity(3)

	bs, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}

	var actual CSRMatrix
	err = json.Unmarshal(bs, &actual)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}
	if !m.Equals(&actual) {
		t.Fatalf("expected %v but found %v", m, actual)
	}
}

func TestCSRMatrix_And(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseMat
	}{
		{CSRMat(2, 2, 0, 1, 0, 1), CSRMat(2, 2, 0, 0, 1, 1), CSRMat(2, 2), CSRMat(2, 2, 0, 0, 0, 1)},
		{CSRMat(2, 2, 0, 0, 1, 1), CSRMat(2, 2, 0, 1, 0, 1), CSRMat(2, 2), CSRMat(2, 2, 0, 0, 0, 1)},
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

func TestCSRMatrix_Or(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseMat
	}{
		{CSRMat(2, 2, 0, 1, 0, 1), CSRMat(2, 2, 0, 0, 1, 1), CSRMat(2, 2), CSRMat(2, 2, 0, 1, 1, 1)},
		{CSRMat(2, 2, 0, 0, 1, 1), CSRMat(2, 2, 0, 1, 0, 1), CSRMat(2, 2), CSRMat(2, 2, 0, 1, 1, 1)},
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

func TestCSRMatrix_XOr(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseMat
	}{
		{CSRMat(2, 2, 0, 1, 0, 1), CSRMat(2, 2, 0, 0, 1, 1), CSRMat(2, 2), CSRMat(2, 2, 0, 1, 1, 0)},
		{CSRMat(2, 2, 0, 0, 1, 1), CSRMat(2, 2, 0, 1, 0, 1), CSRMat(2, 2), CSRMat(2, 2, 0, 1, 1, 0)},
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

func TestCSRMatrix_Negate(t *testing.T) {
	tests := []struct {
		x, expected SparseMat
	}{
		{CSRMat(2, 2, 0, 1, 0, 1), CSRMat(2, 2, 1, 0, 1, 0)},
		{CSRMat(2, 2, 0, 1, 1, 0), CSRMat(2, 2, 1, 0, 0, 1)},
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

func TestCSRMatrix_SwapRows(t *testing.T) {
	tests := []struct {
		input    SparseMat
		a, b     int
		expected SparseMat
	}{
		{CSRIdentity(4), 0, 0, CSRMat(4, 4, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1)},
		{CSRIdentity(4), 1, 3, CSRMat(4, 4, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0)},
		{CSRMat(4, 5, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1), 2, 3, CSRMat(4, 5, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0)},
		{CSRMat(5, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1), 0, 1, CSRMat(5, 4, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1)},
		{CSRMat(5, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1), 1, 2, CSRMat(5, 4, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1)},
		{CSRMat(5, 4, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0), 4, 0, CSRMat(5, 4, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1)},
		{CSRMat(5, 4, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0), 0, 1, CSRMat(5, 4, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0)},
		{CSRMat(5, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1), 0, 2, CSRMat(5, 4, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1)},
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

func BenchmarkCSRMatrix_SwapRows(b *testing.B) {
	a := CSRMat(5, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1)
	aa := CSRMat(5, 4, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0)
	aaa := CSRMat(5, 4, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1)
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

func TestCSRMatrix_SwapColumns(t *testing.T) {
	tests := []struct {
		input    SparseMat
		a, b     int
		expected SparseMat
	}{
		{CSRIdentity(4), 1, 3, CSRMat(4, 4, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0)},
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

func TestCSRMatrix_AddRows(t *testing.T) {
	tests := []struct {
		input        SparseMat
		i1, i2, dest int
		expected     SparseMat
	}{
		{CSRIdentity(3), 0, 1, 2, CSRMat(3, 3, 1, 0, 0, 0, 1, 0, 1, 1, 0)},
		{CSRMat(4, 5, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1), 1, 2, 1, CSRMat(4, 5, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1)},
		{CSRMat(4, 5, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), 0, 2, 0, CSRMat(4, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)},
		{CSRMat(4, 5, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), 1, 2, 1, CSRMat(4, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)},
		{CSRMat(4, 5, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), 1, 2, 1, CSRMat(4, 5, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)},
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

func BenchmarkCSRMatrix_AddRows(b *testing.B) {
	a := CSRMat(5, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 1; j < 5; j++ {
			a.AddRows(0, j, j)
			a.AddRows(4, j, 4)
			a.AddRows(4, j, 4)
		}
	}
}

func TestCSRMatFromVec(t *testing.T) {
	tests := []struct {
		vec SparseVector
	}{
		{CSRVec(5, 1, 0, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := CSRMatFromVec(test.vec).Row(0)
			if !actual.Equals(test.vec) {
				t.Fatalf("expected %v but found %v", test.vec, actual)
			}
		})
	}
}

func TestCSRMatrix_Set(t *testing.T) {
	m := CSRMat(5, 5)
	for i := 0; i < 5; i++ {
		m.Set(i, i, 1)
	}

	expected := CSRIdentity(5)

	if !expected.Equals(m) {
		t.Fatalf("expected %v but found %v", expected, m)
	}
}

func BenchmarkCSRMatrix_T(b *testing.B) {
	data := make([]int, 100)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Intn(2)
	}
	m := CSRMat(10, 10, data...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.T()
	}
}
