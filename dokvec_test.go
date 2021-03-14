package sparsemat

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
)

func TestDOKVector_Add(t *testing.T) {
	tests := []struct {
		a, b, result SparseVector
		expected     SparseVector
	}{
		{DOKVec(3, 0, 1, 0), DOKVec(3, 1, 0, 0), DOKVec(3), DOKVec(3, 1, 1, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Add(test.a, test.b)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestDOKVector_Dot(t *testing.T) {
	tests := []struct {
		a, b     SparseVector
		expected int
	}{
		{DOKVec(3, 1, 1, 1), DOKVec(3, 1, 1, 1), 1},
		{DOKVec(3, 1, 0, 1), DOKVec(3, 1, 1, 1), 0},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.a.Dot(test.b)
			if test.expected != actual {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestDOKVector_Mul(t *testing.T) {
	tests := []struct {
		a        SparseVector
		b        SparseMat
		result   SparseVector
		expected SparseVector
	}{
		{DOKVec(3, 1, 0, 1), DOKIdentity(3), DOKVec(3), DOKVec(3, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.MulMat(test.a, test.b)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestDOKVector_Equals(t *testing.T) {
	tests := []struct {
		a, b     SparseVector
		expected bool
	}{
		{DOKVec(3), DOKVec(3), true},
		{DOKVec(3), DOKVec(4), false},
		{DOKVec(3, 1, 0, 1), DOKVec(3), false},
		{DOKVec(3, 1, 0, 1), DOKVec(3, 1, 0, 1), true},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if test.a.Equals(test.b) != test.expected {
				t.Fatalf("expected %v for %v == %v", test.expected, test.a, test.b)
			}
		})
	}
}

func TestDOKVector_Set(t *testing.T) {
	tests := []struct {
		source, result SparseVector
	}{
		{DOKVec(5, 1, 0, 1, 0, 1), DOKVec(5)},
		{DOKIdentity(5).Row(2), DOKVec(5)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			for i := 0; i < test.source.Len(); i++ {
				test.result.Set(i, test.source.At(i))
			}

			if !test.source.Equals(test.result) {
				t.Fatalf("expected %v but found %v", test.source, test.result)
			}
		})
	}
}

func TestDOKVector_SetVec(t *testing.T) {
	tests := []struct {
		original         SparseVector
		setToSlice       SparseVector
		index            int
		expectedOriginal SparseVector
	}{
		{DOKVec(5, 1, 1, 1, 1, 1), DOKVec(3, 0, 1, 0), 1, DOKVec(5, 1, 0, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.original.SetVec(test.setToSlice, test.index)

			if !test.original.Equals(test.expectedOriginal) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedOriginal, test.original)
			}
		})
	}
}

func TestDOKVector_NonzeroValues(t *testing.T) {
	tests := []struct {
		input    SparseVector
		expected map[int]int
	}{
		{DOKIdentity(4).Row(2), map[int]int{2: 1}},
		{DOKMat(4, 6, 1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1).Row(0), map[int]int{0: 1, 1: 1, 3: 1}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.input.NonzeroValues()
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestDOKVector_Slice(t *testing.T) {
	tests := []struct {
		original         SparseVector
		i, len           int
		addToSlice       SparseVector
		expectedOriginal SparseVector
		expectedSlice    SparseVector
	}{
		{DOKVec(5, 1, 0, 1, 0, 1), 1, 3, DOKVec(3, 1, 1, 1), DOKVec(5, 1, 1, 0, 1, 1), DOKVec(3, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i, test.len)
			sl.Add(sl, test.addToSlice)
			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}
		})
	}
}

func TestDOKVector_Slice2(t *testing.T) {
	tests := []struct {
		original      SparseVector
		i1, len1      int
		i2, len2      int
		expectedSlice SparseVector
	}{
		{DOKVec(7, 0, 0, 1, 0, 1, 0, 0), 1, 5, 1, 3, DOKVec(3, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i1, test.len1).Slice(test.i2, test.len2)

			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}

		})
	}
}

func TestDOKVector_And(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseVector
	}{
		{DOKVec(4, 0, 1, 0, 1), DOKVec(4, 0, 0, 1, 1), DOKVec(4), DOKVec(4, 0, 0, 0, 1)},
		{DOKVec(4, 0, 0, 1, 1), DOKVec(4, 0, 1, 0, 1), DOKVec(4), DOKVec(4, 0, 0, 0, 1)},
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

func TestDOKVector_Or(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseVector
	}{
		{DOKVec(4, 0, 1, 0, 1), DOKVec(4, 0, 0, 1, 1), DOKVec(4), DOKVec(4, 0, 1, 1, 1)},
		{DOKVec(4, 0, 0, 1, 1), DOKVec(4, 0, 1, 0, 1), DOKVec(4), DOKVec(4, 0, 1, 1, 1)},
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

func TestDOKVector_XOr(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseVector
	}{
		{DOKVec(4, 0, 1, 0, 1), DOKVec(4, 0, 0, 1, 1), DOKVec(4), DOKVec(4, 0, 1, 1, 0)},
		{DOKVec(4, 0, 0, 1, 1), DOKVec(4, 0, 1, 0, 1), DOKVec(4), DOKVec(4, 0, 1, 1, 0)},
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

func TestDOKVector_Negate(t *testing.T) {
	tests := []struct {
		x, expected SparseVector
	}{
		{DOKVec(4, 0, 1, 0, 1), DOKVec(4, 1, 0, 1, 0)},
		{DOKVec(4, 0, 0, 1, 1), DOKVec(4, 1, 1, 0, 0)},
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

func TestDOKVector_JSON(t *testing.T) {
	v := DOKVec(5, 1, 0, 1, 0, 1)

	bs, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}

	var actual DOKVector
	err = json.Unmarshal(bs, &actual)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}
	if !v.Equals(&actual) {
		t.Fatalf("expected %v but found %v", v, actual)
	}
}

func TestDOKVecCopy(t *testing.T) {
	tests := []struct {
		vec SparseVector
	}{
		{DOKVec(5, 1, 0, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := DOKVecCopy(test.vec)
			if !actual.Equals(test.vec) {
				t.Fatalf("expected %v but foudn %v", test.vec, actual)
			}
		})
	}
}