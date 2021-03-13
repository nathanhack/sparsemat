package sparsemat

type SparseVector interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(bytes []byte) error
	String() string
	At(i int) int
	Set(i, value int)
	SetVec(a SparseVector, i int)
	Len() int
	Dot(a SparseVector) int
	NonzeroValues() (indexToValues map[int]int)
	Slice(i, length int) SparseVector
	Add(a, b SparseVector)
	Equals(v SparseVector) bool
	MulMat(vec SparseVector, mat SparseMat)
	MatMul(mat SparseMat, vec SparseVector)
	And(a, b SparseVector)
	Or(a, b SparseVector)
	XOr(a, b SparseVector)
	Negate()
}

type SparseMat interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(bytes []byte) error
	Slice(i, j, rows, cols int) SparseMat
	Dims() (int, int)
	At(i, j int) int
	SwapRows(i1, i2 int)
	SwapColumns(j1, j2 int)
	AddRows(i1, i2, dest int)
	Set(i, j, value int)
	T() SparseMat
	Zeroize()
	ZeroizeRange(i, j, rows, cols int)
	Mul(a, b SparseMat)
	Add(a, b SparseMat)
	Column(j int) SparseVector
	SetColumn(j int, vec SparseVector)
	Row(i int) SparseVector
	SetRow(i int, vec SparseVector)
	Equals(m SparseMat) bool
	String() string
	SetMatrix(a SparseMat, iOffset, jOffset int)
	Negate()
	And(a, b SparseMat)
	Or(a, b SparseMat)
	XOr(a, b SparseMat)
}
