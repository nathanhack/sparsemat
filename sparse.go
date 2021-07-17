package sparsemat

type SparseVector interface {
	Add(a, b SparseVector) SparseVector
	And(a, b SparseVector) SparseVector
	At(i int) int
	Dot(a SparseVector) int
	Equals(v SparseVector) bool
	HammingDistance(a SparseVector) int
	HammingWeight() int
	IsZero() bool
	Len() int
	MarshalJSON() ([]byte, error)
	MatMul(mat SparseMat, vec SparseVector) SparseVector
	MulMat(vec SparseVector, mat SparseMat) SparseVector
	Negate() SparseVector
	NonzeroMap() (indicesToValues map[int]int)
	NonzeroArray() (indices []int)
	NextSet(startingIndex int) (index int, has bool)
	Or(a, b SparseVector) SparseVector
	Set(i, value int) SparseVector
	SetVec(a SparseVector, i int) SparseVector
	Slice(i, length int) SparseVector
	String() string
	XOr(a, b SparseVector) SparseVector
	UnmarshalJSON(bytes []byte) error
}

type SparseMat interface {
	Add(a, b SparseMat) SparseMat
	AddRows(i1, i2, dest int) SparseMat
	And(a, b SparseMat) SparseMat
	At(i, j int) int
	Column(j int) SparseVector
	Dims() (int, int)
	Equals(m SparseMat) bool
	MarshalJSON() ([]byte, error)
	Mul(a, b SparseMat) SparseMat
	Negate() SparseMat
	Or(a, b SparseMat) SparseMat
	Row(i int) SparseVector
	Set(i, j, value int) SparseMat
	SetColumn(j int, vec SparseVector) SparseMat
	SetMatrix(a SparseMat, iOffset, jOffset int) SparseMat
	SetRow(i int, vec SparseVector) SparseMat
	Slice(i, j, rows, cols int) SparseMat
	String() string
	SwapRows(i1, i2 int) SparseMat
	SwapColumns(j1, j2 int) SparseMat
	T() SparseMat
	UnmarshalJSON(bytes []byte) error
	XOr(a, b SparseMat) SparseMat
	Zeroize() SparseMat
	ZeroizeRange(i, j, rows, cols int) SparseMat
}
