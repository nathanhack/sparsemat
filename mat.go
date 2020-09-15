package mat

import (
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Matrix struct {
	rowValues map[int]map[int]int //hold rowValues for (X,Y)
	colValues map[int]map[int]int //easy access to (Y,X)
	rows      int                 // total number rows available to this matrix
	rowStart  int                 // [rowStart,rowEnd)
	cols      int                 // total number cols available to this matrix
	colStart  int                 // [colStart,colEnd)

}

//NewMat creates a new matrix with the specified number of rows and cols.
// If values is empty, the matrix will be zeroized.
// If values are not empty it must have rows*cols items.  The values are expected to
// be 0's or 1's anything else may have unexpected behavior matrix's methods.
func NewMat(rows, cols int, values ...int) *Matrix {
	if len(values) != 0 && len(values) != rows*cols {
		panic(fmt.Sprintf("matrix data length (%v) to size mismatch expected %v", len(values), rows*cols))
	}

	mat := Matrix{
		rowValues: map[int]map[int]int{},
		colValues: map[int]map[int]int{},
		rows:      rows,
		rowStart:  0,
		cols:      cols,
		colStart:  0,
	}

	if len(values) > 0 {
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				index := i*cols + j
				mat.set(i, j, values[index])
			}
		}
	}

	return &mat
}

func NewMatFromVec(vec *Vector) *Matrix {
	return Copy(vec.mat)
}

//Identity create an identity matrix (one's on the diagonal).
func Identity(size int) *Matrix {
	mat := Matrix{
		rowValues: map[int]map[int]int{},
		colValues: map[int]map[int]int{},
		rows:      size,
		rowStart:  0,
		cols:      size,
		colStart:  0,
	}

	for i := 0; i < size; i++ {
		mat.set(i, i, 1)
	}

	return &mat
}

//Copy will create a NEW matrix that will have all the same values as m.
func Copy(m *Matrix) *Matrix {
	mat := Matrix{
		rowValues: make(map[int]map[int]int),
		colValues: make(map[int]map[int]int),
		rows:      m.rows,
		rowStart:  0,
		cols:      m.cols,
		colStart:  0,
	}

	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			mat.set(i, j, m.At(i, j))
		}
	}

	return &mat
}

//Slice creates a slice of the matrix.  The slice will be connected to the original matrix, changes work to one
// causes changes in the other
func (mat *Matrix) Slice(i, j, rows, cols int) *Matrix {
	if rows <= 0 || cols <= 0 {
		panic("slice rows and cols must >= 1")
	}

	mat.checkRowBounds(i)
	mat.checkColBounds(j)
	r := i + mat.rowStart
	c := j + mat.colStart

	if r+rows-1 > mat.rows || c+cols-1 > mat.cols {
		panic("slice rows and cols must be in bounds of matrix")
	}
	mat.checkRowBounds(i + rows - 1)
	mat.checkColBounds(j + cols - 1)

	return &Matrix{
		rowValues: mat.rowValues,
		rows:      rows,
		rowStart:  r,
		colValues: mat.colValues,
		cols:      cols,
		colStart:  c,
	}
}

func (mat *Matrix) checkRowBounds(i int) {
	if i < 0 || i >= mat.rows {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, mat.rows-1))
	}
}

func (mat *Matrix) checkColBounds(j int) {
	if j < 0 || j >= mat.cols {
		panic(fmt.Sprintf("%v out of range: [0-%v]", j, mat.cols-1))
	}
}

//Dims returns the dimensions of the matrix.
func (mat *Matrix) Dims() (int, int) {
	return mat.rows, mat.cols
}

//At returns the value at row index i and column index j.
func (mat *Matrix) At(i, j int) int {
	mat.checkRowBounds(i)
	mat.checkColBounds(j)
	r := i + mat.rowStart
	c := j + mat.colStart

	return mat.at(r, c)
}

func (mat *Matrix) at(r, c int) int {
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
func (mat *Matrix) Set(i, j, value int) {
	mat.checkRowBounds(i)
	mat.checkColBounds(j)
	r := i + mat.rowStart
	c := j + mat.colStart

	mat.set(r, c, value)
}

func (mat *Matrix) set(r, c, value int) {
	if value == 0 {
		ys, ok := mat.rowValues[r]
		if !ok {
			return
		}

		_, ok = ys[c]
		if !ok {
			return
		}

		delete(ys, c)
		if len(mat.rowValues[r]) == 0 {
			delete(mat.rowValues, r)
		}

		delete(mat.colValues[c], r)
		if len(mat.colValues[c]) == 0 {
			delete(mat.colValues, c)
		}

		return
	}

	ys, ok := mat.rowValues[r]
	if !ok {
		ys = make(map[int]int)
		mat.rowValues[r] = ys
	}
	ys[c] = value

	xs, ok := mat.colValues[c]
	if !ok {
		xs = make(map[int]int)
		mat.colValues[c] = xs
	}
	xs[r] = value
}

//T returns a matrix that is the transpose of the underlying matrix. Note the transpose
// is connected to matrix it is a transpose of, and changes made to one affect the other.
func (mat *Matrix) T() *Matrix {
	return &Matrix{
		rowValues: mat.colValues,
		rows:      mat.cols,
		rowStart:  mat.colStart,
		colValues: mat.rowValues,
		cols:      mat.rows,
		colStart:  mat.rowStart,
	}
}

//Zeroize take the current matrix sets all values to 0.
func (mat *Matrix) Zeroize() {
	mat.zeroize(0, 0, mat.rows, mat.cols)
}

//ZeroizeRange take the current matrix sets values inside the range to zero.
func (mat *Matrix) ZeroizeRange(i, j, rows, cols int) {
	if i < 0 || j < 0 || rows < 0 || cols < 0 {
		panic("zeroize must have positive values")
	}
	if mat.rows < i+rows || mat.cols < j+cols {
		panic(fmt.Sprintf("zeroize bounds check failed can't zeroize shape (%v,%v) on a (%v,%v) matrix", i+rows, j+cols, mat.rows, mat.cols))
	}

	r := i + mat.rowStart
	c := j + mat.colStart

	mat.zeroize(r, c, rows, cols)
}

func (mat *Matrix) zeroize(r, c, rows, col int) {
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
func (mat *Matrix) Mul(a, b *Matrix) {
	if a == nil || b == nil {
		panic("multiply input was found to be nil")
	}

	if mat == a || mat == b {
		panic("multiply self assignment not allowed")
	}

	if a.cols != b.rows {
		panic(fmt.Sprintf("multiply shape misalignment can't multiply (%v,%v)x(%v,%v)", a.rows, a.cols, b.rows, b.cols))
	}

	mRows, mCols := mat.Dims()
	aRows, _ := a.Dims()
	_, bCols := b.Dims()
	if mRows != aRows || mCols != bCols {
		panic(fmt.Sprintf("mat shape (%v,%v) does not match expected (%v,%v)", mat.rows, mat.cols, a.rows, b.cols))
	}

	mat.mul(a, b)
}

func (mat *Matrix) mul(a, b *Matrix) {
	//first we need to clear mat
	mat.zeroize(0, 0, mat.rows, mat.cols)

	for r, cs := range a.rowValues {
		i := r - a.rowStart

		for c, rs := range b.colValues {
			j := c - b.colStart
			value := 0
			for ics, v1 := range cs {
				ci := ics - a.colStart

				v2, ok := rs[ci+b.rowStart]
				if ok {
					value += v1 * v2
				}
			}
			value = value % 2

			mat.Set(i, j, value)
		}
	}
}

//Add stores the addition of a and b in this matrix.
func (mat *Matrix) Add(a, b *Matrix) {
	if a == nil || b == nil {
		panic("addition input was found to be nil")
	}
	if mat == a || mat == b {
		panic("addition self assignment not allowed")
	}

	if a.rows != b.rows || a.cols != b.cols {
		panic(fmt.Sprintf("addition input mat shapes do not match a=(%v,%v) b=(%v,%v)", a.rows, a.cols, b.rows, b.cols))
	}
	if mat.rows != a.rows || mat.cols != a.cols {
		panic(fmt.Sprintf("mat shape (%v,%v) does not match expected (%v,%v)", mat.rows, mat.cols, a.rows, a.cols))
	}

	mat.add(a, b)
}

func (mat *Matrix) add(a, b *Matrix) {
	//first we need to clear mat
	mat.setMatrix(a, 0, 0)

	for r, cs := range b.rowValues {
		i := r - b.rowStart
		mr := i + mat.rowStart
		for c, v := range cs {
			j := c - b.colStart
			mc := j + mat.colStart
			mat.set(mr, mc, (mat.at(mr, mc)+v)%2)
		}
	}
}

//Column returns a map containing the non zero row indices as the keys and it's associated values.
func (mat *Matrix) Column(j int) *TransposedVector {
	mat.checkColBounds(j)

	return &TransposedVector{
		mat: mat.Slice(0, j, mat.rows, 1),
	}
}

//SetColumn sets the values in column j. The values' keys are expected to be row indices.
func (mat *Matrix) SetColumn(j int, vec *TransposedVector) {
	mat.checkColBounds(j)

	if mat.rows != vec.Len() {
		panic("matrix number of columns must equal length of vector")
	}

	c := j + mat.colStart

	//first we'll zeroize
	rs := mat.colValues[c]
	for r := range rs {
		mat.set(r, c, 0)
	}

	//now set the new values
	for i, v := range vec.mat.colValues[0] {
		r := i + mat.rowStart
		mat.set(r, c, v)
	}
}

//Row returns a map containing the non zero column indices as the keys and it's associated values.
func (mat *Matrix) Row(i int) *Vector {
	mat.checkRowBounds(i)

	return &Vector{
		mat: mat.Slice(i, 0, 1, mat.cols),
	}
}

//SetRow sets the values in row i. The values' keys are expected to be column indices.
func (mat *Matrix) SetRow(i int, vec *Vector) {
	mat.checkColBounds(i)

	if mat.cols != vec.Len() {
		panic("matrix number of columns must equal length of vector")
	}

	r := i + mat.rowStart

	//first we'll zeroize
	cs := mat.rowValues[r]
	for c := range cs {
		mat.set(r, c, 0)
	}

	//now set the new values
	for j, v := range vec.mat.rowValues[0] {
		c := j + mat.colStart
		mat.set(r, c, v)
	}
}

//Equals return true if the m matrix has the same shape and values as this matrix.
func (mat *Matrix) Equals(m *Matrix) bool {
	if mat == m {
		return true
	}

	if mat == nil || m == nil {
		return false
	}

	if mat.rows != m.rows || mat.cols != m.cols {
		return false
	}

	for i := 0; i < mat.rows; i++ {
		r := i + mat.rowStart
		cs, ok1 := mat.rowValues[r]
		ar := i + m.rowStart
		acs, ok2 := m.rowValues[ar]

		if !ok1 && !ok2 {
			continue
		}

		for j := 0; j < mat.cols; j++ {
			c := j + mat.colStart
			v1, ok1 := cs[c]
			ac := j + m.colStart
			v2, ok2 := acs[ac]

			if !ok1 && !ok2 {
				continue
			}

			if (!ok1 && ok2) || (ok1 && !ok2) {
				return false
			}

			if v1 != v2 {
				return false
			}
		}
	}
	return true
}

//String returns a string representation of this matrix.
func (mat Matrix) String() string {
	buff := &strings.Builder{}
	table := tablewriter.NewWriter(buff)

	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)

	for i := 0; i < mat.rows; i++ {
		row := make([]string, mat.cols)
		for j := 0; j < mat.cols; j++ {
			row[j] = fmt.Sprint(mat.At(i, j))
		}
		table.Append(row)
	}

	table.Render()
	return buff.String()
}

//SetMatrix replaces the values of this matrix with the values of from matrix a. The shape of 'a' must be less than or equal mat.
// If the 'a' shape is less then iOffset and jOffset can be used to place 'a' matrix in a specific location.
func (mat *Matrix) SetMatrix(a *Matrix, iOffset, jOffset int) {
	if iOffset < 0 || jOffset < 0 {
		panic("offsets must be positive values [0,+)")
	}
	if mat.rows < iOffset+a.rows || mat.cols < jOffset+a.cols {
		panic(fmt.Sprintf("set matrix have equal or smaller shape (%v,%v), found a=(%v,%v)", mat.rows, mat.cols, iOffset+a.rows, jOffset+a.cols))
	}

	mat.setMatrix(a, iOffset, jOffset)
}

func (mat *Matrix) setMatrix(a *Matrix, rOffset, cOffset int) {
	mat.zeroize(rOffset, cOffset, a.rows, a.cols)

	for r, cs := range a.rowValues {
		i := r - a.rowStart
		mr := i + mat.rowStart + rOffset
		for c, v := range cs {
			j := c - a.colStart
			mc := j + mat.colStart + cOffset
			mat.set(mr, mc, v)
		}
	}
}
