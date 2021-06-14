package sparsemat

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func TestCSRVector_Add(t *testing.T) {
	tests := []struct {
		a, b, result SparseVector
		expected     SparseVector
	}{
		{CSRVec(3, 0, 1, 0), CSRVec(3, 1, 0, 0), CSRVec(3), CSRVec(3, 1, 1, 0)},
		{CSRVec(3, 0, 1, 0), CSRVec(3, 0, 0, 1), CSRVec(3), CSRVec(3, 0, 1, 1)},
		{CSRVec(3, 1, 1, 0), CSRVec(3, 0, 1, 1), CSRVec(3), CSRVec(3, 1, 0, 1)},
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

func BenchmarkCSRVector_Add(b *testing.B) {
	benchmarks := []struct {
		size int
	}{
		{10},
		{100},
		{1000},
		{10000},
	}
	for bi, bm := range benchmarks {
		b.Run(strconv.Itoa(bi), func(b *testing.B) {
			av, bv := randomCSRVector(bm.size), randomCSRVector(bm.size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				av.Add(av, bv)
			}
		})
	}
}

func randomCSRVector(size int) SparseVector {
	vec := CSRVec(size)

	for i := 0; i < size; i++ {
		vec.Set(i, rand.Intn(2))
	}
	return vec
}

func TestCSRVector_Dot(t *testing.T) {
	tests := []struct {
		a, b     SparseVector
		expected int
	}{
		{CSRVec(3, 1, 1, 1), CSRVec(3, 1, 1, 1), 1},
		{CSRVec(3, 1, 0, 1), CSRVec(3, 1, 1, 1), 0},
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

func TestCSRVector_Mul(t *testing.T) {
	tests := []struct {
		a        SparseVector
		b        SparseMat
		result   SparseVector
		expected SparseVector
	}{
		{CSRVec(3, 1, 0, 1), CSRIdentity(3), CSRVec(3), CSRVec(3, 1, 0, 1)},
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

func TestCSRVector_Equals(t *testing.T) {
	tests := []struct {
		a, b     SparseVector
		expected bool
	}{
		{CSRVec(3), CSRVec(3), true},
		{CSRVec(3), CSRVec(4), false},
		{CSRVec(3, 1, 0, 1), CSRVec(3), false},
		{CSRVec(3, 1, 0, 1), CSRVec(3, 1, 0, 1), true},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if test.a.Equals(test.b) != test.expected {
				t.Fatalf("expected %v for %v == %v", test.expected, test.a, test.b)
			}
		})
	}
}

func TestCSRVector_Set(t *testing.T) {
	tests := []struct {
		input  SparseVector
		index  int
		value  int
		result SparseVector
	}{
		{CSRVec(5, 1, 0, 1, 0, 1), 0, 0, CSRVec(5, 0, 0, 1, 0, 1)},
		{CSRVec(5, 1, 0, 1, 0, 1), 1, 1, CSRVec(5, 1, 1, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.input.Set(test.index, test.value)

			if !test.input.Equals(test.result) {
				t.Fatalf("expected %v but found %v", test.input, test.result)
			}
		})
	}
}

func TestCSRVector_SetVec(t *testing.T) {
	tests := []struct {
		original         SparseVector
		setToSlice       SparseVector
		index            int
		expectedOriginal SparseVector
	}{
		{CSRVec(5, 1, 1, 1, 1, 1), CSRVec(3, 0, 1, 0), 1, CSRVec(5, 1, 0, 1, 0, 1)},
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

func TestCSRVector_NonzeroMap(t *testing.T) {
	tests := []struct {
		input    SparseVector
		expected map[int]int
	}{
		{CSRIdentity(4).Row(2), map[int]int{2: 1}},
		{CSRMat(4, 6, 1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1).Row(0), map[int]int{0: 1, 1: 1, 3: 1}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.input.NonzeroMap()
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestCSRVector_NonzeroArray(t *testing.T) {
	tests := []struct {
		input    SparseVector
		expected []int
	}{
		{CSRIdentity(4).Row(2), []int{2}},
		{CSRMat(4, 6, 1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1).Row(0), []int{0, 1, 3}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.input.NonzeroArray()
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestCSRVector_Slice(t *testing.T) {
	tests := []struct {
		original      SparseVector
		i, len        int
		expectedSlice SparseVector
	}{
		{CSRVec(5, 1, 0, 1, 0, 1), 1, 3, CSRVec(3, 0, 1, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i, test.len)
			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}
		})
	}
}

func TestCSRVector_Slice2(t *testing.T) {
	tests := []struct {
		original      SparseVector
		i1, len1      int
		i2, len2      int
		expectedSlice SparseVector
	}{
		{CSRVec(7, 0, 0, 1, 0, 1, 0, 0), 1, 5, 1, 3, CSRVec(3, 1, 0, 1)},
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

func TestCSRVector_And(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseVector
	}{
		{CSRVec(4, 0, 1, 0, 1), CSRVec(4, 0, 0, 1, 1), CSRVec(4), CSRVec(4, 0, 0, 0, 1)},
		{CSRVec(4, 0, 0, 1, 1), CSRVec(4, 0, 1, 0, 1), CSRVec(4), CSRVec(4, 0, 0, 0, 1)},
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

func TestCSRVector_Or(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseVector
	}{
		{CSRVec(4, 0, 1, 0, 1), CSRVec(4, 0, 0, 1, 1), CSRVec(4), CSRVec(4, 0, 1, 1, 1)},
		{CSRVec(4, 0, 0, 1, 1), CSRVec(4, 0, 1, 0, 1), CSRVec(4), CSRVec(4, 0, 1, 1, 1)},
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

func TestCSRVector_XOr(t *testing.T) {
	tests := []struct {
		x, y, result, expected SparseVector
	}{
		{CSRVec(4, 0, 1, 0, 1), CSRVec(4, 0, 0, 1, 1), CSRVec(4), CSRVec(4, 0, 1, 1, 0)},
		{CSRVec(4, 0, 0, 1, 1), CSRVec(4, 0, 1, 0, 1), CSRVec(4), CSRVec(4, 0, 1, 1, 0)},
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

func TestCSRVector_Negate(t *testing.T) {
	tests := []struct {
		x, expected SparseVector
	}{
		{CSRVec(4, 0, 1, 0, 1), CSRVec(4, 1, 0, 1, 0)},
		{CSRVec(4, 0, 0, 1, 1), CSRVec(4, 1, 1, 0, 0)},
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

func TestCSRVector_JSON(t *testing.T) {
	v := CSRVec(5, 1, 0, 1, 0, 1)

	bs, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}

	var actual CSRVector
	err = json.Unmarshal(bs, &actual)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}
	if !v.Equals(&actual) {
		t.Fatalf("expected %v but found %v", v, actual)
	}
}

func TestCSRVecCopy(t *testing.T) {
	tests := []struct {
		vec SparseVector
	}{
		{CSRVec(5, 1, 0, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := CSRVecCopy(test.vec)
			if !actual.Equals(test.vec) {
				t.Fatalf("expected %v but foudn %v", test.vec, actual)
			}
		})
	}
}

func BenchmarkCSRVecCopy(b *testing.B) {
	vec := CSRVec(2000)
	for i := 0; i < vec.Len(); i++ {
		vec.Set(i, rand.Intn(2))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CSRVecCopy(vec)
	}
}

func BenchmarkCSRVector_Dot(b *testing.B) {
	benchmarks := []struct {
		size int
	}{
		{0},
		{1},
		{2},
		{3},
		{4},
		{5},
		{1000},
	}
	for ti, bm := range benchmarks {
		b.Run(strconv.Itoa(ti), func(b *testing.B) {
			data1 := make([]int, bm.size)
			data2 := make([]int, bm.size)

			for j := 0; j < bm.size; j++ {
				data1[j] = rand.Intn(2)
				data2[j] = rand.Intn(2)
			}

			a := CSRVec(len(data1), data1...)
			aa := CSRVec(len(data2), data2...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				a.Dot(a)
				a.Dot(aa)
				aa.Dot(a)
			}
		})
	}
}

func BenchmarkCSRVector_MulMat(b *testing.B) {
	data := make([]int, 100)
	data1 := make([]int, 10)

	for i := 0; i < len(data); i++ {
		data[i] = rand.Intn(2)
	}
	for i := 0; i < len(data1); i++ {
		data1[i] = rand.Intn(2)
	}
	m := CSRMat(10, 10, data...)
	v := CSRVec(10, data1...)
	org := CSRVec(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		org.MulMat(v, m)
	}
}

func TestCSRVector_NextSet(t *testing.T) {
	tests := []struct {
		vec           SparseVector
		index         int
		expectedIndex int
		expectedHas   bool
	}{
		{CSRVec(5, 1, 0, 0, 0, 1), 0, 0, true},
		{CSRVec(5, 1, 0, 0, 0, 1), 1, 4, true},
		{CSRVec(5, 1, 0, 0, 0, 0), 1, -1, false},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			index, has := test.vec.NextSet(test.index)

			if index != test.expectedIndex {
				t.Fatalf("expected %v index but found %v", test.expectedIndex, index)
			}

			if has != test.expectedHas {
				t.Fatalf("expected %v bool but found %v", test.expectedHas, has)
			}
		})
	}
}
