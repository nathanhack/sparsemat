package sparsemat

type SparseVector interface {
	Add(a, b SparseVector)
	And(a, b SparseVector)
	At(i int) int
	Dot(a SparseVector) int
	Equals(v SparseVector) bool
	HammingDistance(a SparseVector) int
	HammingWeight() int
	IsZero() bool
	Len() int
	MarshalJSON() ([]byte, error)
	MatMul(mat SparseMat, vec SparseVector)
	MulMat(vec SparseVector, mat SparseMat)
	Negate() SparseVector
	NonzeroMap() (indicesToValues map[int]int)
	NonzeroArray() (indices []int)
	NextSet(startingIndex int) (index int, has bool)
	Or(a, b SparseVector)
	Set(i, value int)
	SetVec(a SparseVector, i int)
	Slice(i, length int) SparseVector
	String() string
	XOr(a, b SparseVector)
	UnmarshalJSON(bytes []byte) error
}

type SparseMat interface {
	Add(a, b SparseMat)
	AddRows(i1, i2, dest int)
	And(a, b SparseMat)
	At(i, j int) int
	Column(j int) SparseVector
	Dims() (int, int)
	Equals(m SparseMat) bool
	MarshalJSON() ([]byte, error)
	Mul(a, b SparseMat)
	Negate() SparseMat
	Or(a, b SparseMat)
	Row(i int) SparseVector
	Set(i, j, value int)
	SetColumn(j int, vec SparseVector)
	SetMatrix(a SparseMat, iOffset, jOffset int)
	SetRow(i int, vec SparseVector)
	Slice(i, j, rows, cols int) SparseMat
	String() string
	SwapRows(i1, i2 int)
	SwapColumns(j1, j2 int)
	T() SparseMat
	UnmarshalJSON(bytes []byte) error
	XOr(a, b SparseMat)
	Zeroize()
	ZeroizeRange(i, j, rows, cols int)
}
