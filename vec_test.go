package mat

import (
	"strconv"
	"testing"
)

func TestVector_Add(t *testing.T) {
	tests := []struct {
		a, b, result *Vector
		expected     *Vector
	}{
		{NewVec(3, 0, 1, 0), NewVec(3, 1, 0, 0), NewVec(3), NewVec(3, 1, 1, 0)},
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

func TestVector_Dot(t *testing.T) {
	tests := []struct {
		a, b     *Vector
		expected int
	}{
		{NewVec(3, 1, 1, 1), NewVec(3, 1, 1, 1), 1},
		{NewVec(3, 1, 0, 1), NewVec(3, 1, 1, 1), 0},
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

func TestVector_Mul(t *testing.T) {
	tests := []struct {
		a        *Vector
		b        *Matrix
		result   *Vector
		expected *Vector
	}{
		{NewVec(3, 1, 0, 1), Identity(3), NewVec(3), NewVec(3, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Mul(test.a, test.b)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestVector_Equals(t *testing.T) {
	tests := []struct {
		a, b     *Vector
		expected bool
	}{
		{NewVec(3), NewVec(3), true},
		{NewVec(3), NewVec(4), false},
		{NewVec(3, 1, 0, 1), NewVec(3), false},
		{NewVec(3, 1, 0, 1), NewVec(3, 1, 0, 1), true},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if test.a.Equals(test.b) != test.expected {
				t.Fatalf("expected %v for %v == %v", test.expected, test.a, test.b)
			}
		})
	}
}

func TestVector_Set(t *testing.T) {
	tests := []struct {
		source, result *Vector
	}{
		{NewVec(5, 1, 0, 1, 0, 1), NewVec(5)},
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

func TestTransposedVector_Set(t *testing.T) {
	tests := []struct {
		source, result *TransposedVector
	}{
		{NewTVec(5, 1, 0, 1, 0, 1), NewTVec(5)},
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

func TestTransposedVector_MulVec(t *testing.T) {
	tests := []struct {
		a        *Matrix
		b        *Vector
		result   *TransposedVector
		expected *Vector
	}{
		{Identity(3), NewVec(3, 0, 1, 0), NewTVec(3), NewVec(3, 0, 1, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.MulVec(test.a, test.b.T())
			if !test.result.T().Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result.T())
			}
		})
	}
}
