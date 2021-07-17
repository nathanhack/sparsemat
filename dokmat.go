package sparsemat

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type DOKMatrix struct {
	rowValues map[int]map[int]int //hold rowValues for (X,Y)
	colValues map[int]map[int]int //easy access to (Y,X)
	rows      int                 // total number rows available to this matrix
	cols      int                 // total number cols available to this matrix
}

func (mat *DOKMatrix) MarshalJSON() ([]byte, error) {
	m := CSRMatCopy(mat)
	return json.Marshal(m)
}

func (mat *DOKMatrix) UnmarshalJSON(bytes []byte) error {
	var csr CSRMatrix
	err := json.Unmarshal(bytes, &csr)
	if err != nil {
		return err
	}

	rows, cols := csr.Dims()
	mat.rows = rows
	mat.cols = cols
	mat.rowValues = map[int]map[int]int{}
	mat.colValues = map[int]map[int]int{}

	for i := 0; i < rows; i++ {
		mat.rowValues[i] = make(map[int]int)
	}
	for j := 0; j < cols; j++ {
		mat.colValues[j] = make(map[int]int)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			mat.set(i, j, csr.At(i, j))
		}
	}
	return nil
}

//NewMat creates a new matrix with the specified number of rows and cols.
// If values is empty, the matrix will be zeroized.
// If values are not empty it must have rows*cols items.  The values are expected to
// be 0's or 1's anything else may have unexpected behavior matrix's methods.
func DOKMat(rows, cols int, values ...int) SparseMat {
	if len(values) != 0 && len(values) != rows*cols {
		panic(fmt.Sprintf("matrix data length (%v) to length mismatch expected %v", len(values), rows*cols))
	}

	return dokMat(rows, cols, values...)
}

func dokMat(rows, cols int, values ...int) *DOKMatrix {
	mat := DOKMatrix{
		rowValues: map[int]map[int]int{},
		colValues: map[int]map[int]int{},
		rows:      rows,
		cols:      cols,
	}

	for i := 0; i < rows; i++ {
		mat.rowValues[i] = make(map[int]int)
	}
	for j := 0; j < cols; j++ {
		mat.colValues[j] = make(map[int]int)
	}

	if len(values) > 0 {
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				index := i*cols + j
				mat.set(i, j, values[index]%2)
			}
		}
	}

	return &mat
}

func DOKMatFromVec(vec SparseVector) SparseMat {
	m := dokMat(1, vec.Len())
	m.SetRow(0, vec)
	return m
}

//Identity create an identity matrix (one's on the diagonal).
func DOKIdentity(size int) SparseMat {
	mat := dokMat(size, size)

	for i := 0; i < size; i++ {
		mat.set(i, i, 1)
	}

	return mat
}

//Copy will create a NEW matrix that will have all the same values as m.
func DOKMatCopy(m SparseMat) SparseMat {
	mRows, mCols := m.Dims()

	mat := dokMat(mRows, mCols)

	for i := 0; i < mRows; i++ {
		for j := 0; j < mCols; j++ {
			mat.set(i, j, m.At(i, j))
		}
	}

	return mat
}

//Slice creates a slice of the matrix.  The slice will be connected to the original matrix, changes to one
// causes changes in the other.
func (mat *DOKMatrix) Slice(i, j, rows, cols int) SparseMat {
	if rows <= 0 || cols <= 0 {
		panic("slice rows and cols must >= 1")
	}

	mat.checkRowBounds(i)
	mat.checkColBounds(j)

	if i+rows-1 > mat.rows || j+cols-1 > mat.cols {
		panic("slice rows and cols must be in bounds of matrix")
	}
	mat.checkRowBounds(i + rows - 1)
	mat.checkColBounds(j + cols - 1)

	return mat.slice(i, j, rows, cols)
}

func (mat *DOKMatrix) slice(r, c, rows, cols int) SparseMat {
	m := dokMat(rows, cols)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			m.set(i, j, mat.at(i+r, j+c))
		}
	}

	return m
}

func (mat *DOKMatrix) checkRowBounds(i int) {
	if i < 0 || i >= mat.rows {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, mat.rows-1))
	}
}

func (mat *DOKMatrix) checkColBounds(j int) {
	if j < 0 || j >= mat.cols {
		panic(fmt.Sprintf("%v out of range: [0-%v]", j, mat.cols-1))
	}
}

//Dims returns the dimensions of the matrix.
func (mat *DOKMatrix) Dims() (int, int) {
	return mat.rows, mat.cols
}

//At returns the value at row index i and column index j.
func (mat *DOKMatrix) At(i, j int) int {
	mat.checkRowBounds(i)
	mat.checkColBounds(j)
	return mat.at(i, j)
}

func (mat *DOKMatrix) SwapRows(i1, i2 int) SparseMat {
	mat.checkRowBounds(i1)
	mat.checkRowBounds(i2)

	if i1 == i2 {
		return mat
	}

	//make a temp copy of i1 row and zero it out
	tmp := make(map[int]int)
	for c := range mat.rowValues[i1] {
		tmp[c] = 1
		mat.set(i1, c, 0)
	}

	//move the from i2 to i1 and zero out i2's
	for c := range mat.rowValues[i2] {
		mat.set(i1, c, 1)
		mat.set(i2, c, 0)
	}

	//last part of the swap
	for c := range tmp {
		mat.set(i2, c, 1)
	}

	return mat
}

func (mat *DOKMatrix) SwapColumns(j1, j2 int) SparseMat {
	mat.checkColBounds(j1)
	mat.checkColBounds(j2)

	//make a temp copy of i1 row
	tmp := make(map[int]int)
	for r := range mat.colValues[j1] {
		tmp[r] = 1
		mat.set(r, j1, 0)
	}

	//move the from j2 to j1 and clear j2
	for r := range mat.colValues[j2] {
		mat.set(r, j1, 1)
		mat.set(r, j2, 0)
	}

	//populate j2 column
	for r := range tmp {
		mat.set(r, j2, 1)
	}

	return mat
}

//AddRows is fast row operation to add two
// rows and put the result in a destination row.
func (mat *DOKMatrix) AddRows(i1, i2, dest int) SparseMat {
	mat.checkRowBounds(i1)
	mat.checkRowBounds(i2)
	mat.checkRowBounds(dest)

	tmp := make(map[int]int)

	for c, v := range mat.rowValues[i1] {
		tmp[c] += v
	}

	for c, v := range mat.rowValues[i2] {
		tmp[c] += v
	}

	for c := range mat.rowValues[dest] {
		mat.set(dest, c, 0)
	}

	for c, v := range tmp {
		mat.set(dest, c, v%2)
	}

	return mat
}

//AddRows is fast row operation to add two
// rows and put the result in a destination row.
func (mat *DOKMatrix) AddCols(j1, j2, dest int) SparseMat {
	mat.checkColBounds(j1)
	mat.checkColBounds(j2)
	mat.checkColBounds(dest)

	tmp := make(map[int]int)

	for r, v := range mat.colValues[j1] {
		tmp[r] += v
	}

	for r, v := range mat.colValues[j2] {
		tmp[r] += v
	}

	for r := range mat.colValues[dest] {
		mat.set(r, dest, 0)
	}

	for r, v := range tmp {
		mat.set(r, dest, v%2)
	}

	return mat
}

func (mat *DOKMatrix) at(r, c int) int {
	ys, ok := mat.rowValues[r]
	if !ok {
		return 0
	}
	v, ok := ys[c]
	if !ok {
		return 0
	}
	return v
}

//Set sets the value at row index i and column index j to value.
func (mat *DOKMatrix) Set(i, j, value int) SparseMat {
	mat.checkRowBounds(i)
	mat.checkColBounds(j)
	mat.set(i, j, value%2)

	return mat
}

func (mat *DOKMatrix) set(r, c, value int) {
	if value == 0 {
		ys, _ := mat.rowValues[r]
		_, ok := ys[c]
		if !ok {
			return
		}

		delete(ys, c)
		delete(mat.colValues[c], r)

		return
	}

	mat.rowValues[r][c] = value
	mat.colValues[c][r] = value
}

//T returns a matrix that is the transpose of the underlying matrix. Note the transpose
// is connected to matrix it is a transpose of, and changes made to one affect the other.
func (mat *DOKMatrix) T() SparseMat {
	m := dokMat(mat.cols, mat.rows)

	for i, cols := range mat.rowValues {
		for j := range cols {
			m.rowValues[j][i] = 1
			m.colValues[i][j] = 1
		}
	}

	return m
}

//Zeroize take the current matrix sets all values to 0.
func (mat *DOKMatrix) Zeroize() SparseMat {
	mat.rowValues = make(map[int]map[int]int)
	mat.colValues = make(map[int]map[int]int)

	return mat
}

//ZeroizeRange take the current matrix sets values inside the range to zero.
func (mat *DOKMatrix) ZeroizeRange(i, j, rows, cols int) SparseMat {
	if i < 0 || j < 0 || rows < 0 || cols < 0 {
		panic("zeroize must have positive values")
	}
	if mat.rows < i+rows || mat.cols < j+cols {
		panic(fmt.Sprintf("zeroize bounds check failed can't zeroize shape (%v,%v) on a (%v,%v) matrix", i+rows, j+cols, mat.rows, mat.cols))
	}

	mat.zeroize(i, j, rows, cols)

	return mat
}

func (mat *DOKMatrix) zeroize(r, c, rows, col int) {
	for rv, cs := range mat.rowValues {
		if rv < r || r+rows <= rv {
			continue
		}
		for cv, v := range cs {
			if v == 0 {
				continue
			}
			if cv < c || c+col <= cv {
				continue
			}
			mat.set(rv, cv, 0)
		}
	}
}

//Mul multiplies two matrices and stores the values in this matrix.
func (mat *DOKMatrix) Mul(a, b SparseMat) SparseMat {
	if a == nil || b == nil {
		panic("multiply input was found to be nil")
	}

	if mat == a || mat == b {
		panic("multiply self assignment not allowed")
	}

	aRows, aCols := a.Dims()
	bRows, bCols := b.Dims()

	if aCols != bRows {
		panic(fmt.Sprintf("multiply shape misalignment can't multiply (%v,%v)x(%v,%v)", aRows, aCols, bRows, bCols))
	}

	mRows, mCols := mat.Dims()
	if mRows != aRows || mCols != bCols {
		panic(fmt.Sprintf("mat shape (%v,%v) does not match expected (%v,%v)", mat.rows, mat.cols, aRows, bCols))
	}
	mat.mul(a, b)

	return mat
}

func (mat *DOKMatrix) mul(a, b SparseMat) {
	//first we need to clear mat
	for i := 0; i < mat.rows; i++ {
		r := a.Row(i)
		for j := 0; j < mat.cols; j++ {
			c := b.Column(j)
			d := r.Dot(c)
			mat.set(i, j, d%2)
		}
	}
}

//Add stores the addition of a and b in this matrix.
func (mat *DOKMatrix) Add(a, b SparseMat) SparseMat {
	if a == nil || b == nil {
		panic("addition input was found to be nil")
	}

	aRows, aCols := a.Dims()
	bRows, bCols := b.Dims()

	if aRows != bRows || aCols != bCols {
		panic(fmt.Sprintf("addition input mat shapes do not match a=(%v,%v) b=(%v,%v)", aRows, aCols, bRows, bCols))
	}
	if mat.rows != aRows || mat.cols != aCols {
		panic(fmt.Sprintf("mat shape (%v,%v) does not match expected (%v,%v)", mat.rows, mat.cols, aRows, aCols))
	}
	mat.add(a, b)

	return mat
}

func (mat *DOKMatrix) add(a, b SparseMat) {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			aa := a.At(i, j)
			bb := b.At(i, j)
			mat.set(i, j, (aa+bb)%2)
		}
	}
}

//Column returns a map containing the non zero row indices as the keys and it's associated values.
func (mat *DOKMatrix) Column(j int) SparseVector {
	mat.checkColBounds(j)

	tmp := make(map[int]int)
	for i := range mat.colValues[j] {
		tmp[i] = 1
	}
	return &DOKVector{
		length: mat.rows,
		values: tmp,
	}
}

//SetColumn sets the values in column j. The values' keys are expected to be row indices.
func (mat *DOKMatrix) SetColumn(j int, vec SparseVector) SparseMat {
	mat.checkColBounds(j)

	if mat.rows != vec.Len() {
		panic("matrix number of columns must equal length of vector")
	}

	mat.colValues[j] = vec.NonzeroMap()
	for i := 0; i < mat.rows; i++ {
		_, has := mat.colValues[j][i]
		if has {
			mat.rowValues[i][j] = 1
		} else {
			_, has = mat.rowValues[i][j]
			if has {
				delete(mat.rowValues[i], j)
			}
		}
	}

	return mat
}

//Row returns a map containing the non zero column indices as the keys and it's associated values.
func (mat *DOKMatrix) Row(i int) SparseVector {
	mat.checkRowBounds(i)

	tmp := make(map[int]int)
	for j := range mat.rowValues[i] {
		tmp[j] = 1
	}
	return &DOKVector{
		length: mat.cols,
		values: tmp,
	}
}

//SetRow sets the values in row i. The values' keys are expected to be column indices.
func (mat *DOKMatrix) SetRow(i int, vec SparseVector) SparseMat {
	mat.checkRowBounds(i)

	if mat.cols != vec.Len() {
		panic("matrix number of columns must equal length of vector")
	}

	mat.rowValues[i] = vec.NonzeroMap()
	for j := 0; j < mat.cols; j++ {
		_, has := mat.rowValues[i][j]
		if has {
			mat.colValues[j][i] = 1
		} else {
			_, has = mat.colValues[j][i]
			if has {
				delete(mat.colValues[j], i)
			}
		}
	}

	return mat
}

//Equals return true if the m matrix has the same shape and values as this matrix.
func (mat *DOKMatrix) Equals(m SparseMat) bool {
	if mat == m {
		return true
	}

	if mat == nil || m == nil {
		return false
	}

	mRows, mCols := m.Dims()

	if mat.rows != mRows || mat.cols != mCols {
		return false
	}

	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			if mat.at(i, j) != m.At(i, j) {
				return false
			}
		}
	}
	return true
}

//String returns a string representation of this matrix.
func (mat *DOKMatrix) String() string {
	buff := &strings.Builder{}
	table := tablewriter.NewWriter(buff)

	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)

	for i := 0; i < mat.rows; i++ {
		row := make([]string, mat.cols)
		for j := 0; j < mat.cols; j++ {
			row[j] = fmt.Sprint(mat.at(i, j))
		}
		table.Append(row)
	}

	table.Render()
	return buff.String()
}

//SetMatrix replaces the values of this matrix with the values of from matrix a. The shape of 'a' must be less than or equal mat.
// If the 'a' shape is less then iOffset and jOffset can be used to place 'a' matrix in a specific location.
func (mat *DOKMatrix) SetMatrix(a SparseMat, iOffset, jOffset int) SparseMat {
	if iOffset < 0 || jOffset < 0 {
		panic("offsets must be positive values [0,+)")
	}
	aRows, aCols := a.Dims()
	if mat.rows < iOffset+aRows || mat.cols < jOffset+aCols {
		panic(fmt.Sprintf("set matrix have equal or smaller shape (%v,%v), found a=(%v,%v)", mat.rows, mat.cols, iOffset+aRows, jOffset+aCols))
	}

	for i := 0; i < aRows; i++ {
		for j := 0; j < aCols; j++ {
			mat.set(i+iOffset, j+jOffset, a.At(i, j))
		}
	}

	return mat
}

//Negate performs a piecewise logical negation.
func (mat *DOKMatrix) Negate() SparseMat {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {

			v := mat.at(i, j) + 1
			mat.set(i, j, v%2)
		}
	}
	return mat
}

//And executes a piecewise logical AND on the two matrices and stores the values in this matrix.
func (mat *DOKMatrix) And(a, b SparseMat) SparseMat {
	if a == nil || b == nil {
		panic("AND input was found to be nil")
	}

	aRows, aCols := a.Dims()
	bRows, bCols := b.Dims()

	if aRows != bRows || aCols != bCols {
		panic(fmt.Sprintf("AND shape misalignment both inputs must be equal found (%v,%v) and (%v,%v)", aRows, aCols, bRows, bCols))
	}

	if mat.rows != aRows || mat.cols != aCols {
		panic(fmt.Sprintf("mat shape (%v,%v) does not match expected (%v,%v)", mat.rows, mat.cols, aRows, bCols))
	}

	mat.and(a, b)

	return mat
}

func (mat *DOKMatrix) and(a, b SparseMat) {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			aa := a.At(i, j)
			bb := b.At(i, j)
			mat.set(i, j, aa&bb)
		}
	}
}

//Or executes a piecewise logical OR on the two matrices and stores the values in this matrix.
func (mat *DOKMatrix) Or(a, b SparseMat) SparseMat {
	if a == nil || b == nil {
		panic("OR input was found to be nil")
	}

	aRows, aCols := a.Dims()
	bRows, bCols := b.Dims()

	if aRows != bRows || aCols != bCols {
		panic(fmt.Sprintf("OR shape misalignment both inputs must be equal found (%v,%v) and (%v,%v)", aRows, aCols, bRows, bCols))
	}

	if mat.rows != aRows || mat.cols != aCols {
		panic(fmt.Sprintf("mat shape (%v,%v) does not match expected (%v,%v)", mat.rows, mat.cols, aRows, bCols))
	}

	mat.or(a, b)

	return mat
}

func (mat *DOKMatrix) or(a, b SparseMat) {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			mat.set(i, j, a.At(i, j)|b.At(i, j))
		}
	}
}

//XOr executes a piecewise logical XOR on the two matrices and stores the values in this matrix.
func (mat *DOKMatrix) XOr(a, b SparseMat) SparseMat {
	if a == nil || b == nil {
		panic("XOR input was found to be nil")
	}

	aRows, aCols := a.Dims()
	bRows, bCols := b.Dims()

	if aRows != bRows || aCols != bCols {
		panic(fmt.Sprintf("XOR shape misalignment both inputs must be equal found (%v,%v) and (%v,%v)", aRows, aCols, bRows, bCols))
	}

	if mat.rows != aRows || mat.cols != aCols {
		panic(fmt.Sprintf("mat shape (%v,%v) does not match expected (%v,%v)", mat.rows, mat.cols, aRows, bCols))
	}
	mat.xor(a, b)

	return mat
}

func (mat *DOKMatrix) xor(a, b SparseMat) {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			mat.set(i, j, a.At(i, j)^b.At(i, j))
		}
	}
}
