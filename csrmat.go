package sparsemat

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type CSRMatrix struct {
	rows, cols int
	data       [][]int
}

type csrMatrix struct {
	Rows, Cols int
	Data       [][]int
}

func (mat *CSRMatrix) MarshalJSON() ([]byte, error) {
	return json.Marshal(csrMatrix{
		Rows: mat.rows,
		Cols: mat.cols,
		Data: mat.data,
	})
}

func (mat *CSRMatrix) UnmarshalJSON(bytes []byte) error {
	var m csrMatrix
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return err
	}

	mat.rows = m.Rows
	mat.cols = m.Cols
	mat.data = m.Data
	return nil
}

//CSRMat creates a new matrix with the specified number of rows and cols.
// If values is empty, the matrix will be zeroized.
// If values are not empty it must have rows*cols items.  The values are expected to
// be 0's or 1's anything else may have unexpected behavior matrix's methods.
func CSRMat(rows, cols int, values ...int) SparseMat {
	return csrMat(rows, cols, values...)
}

func csrMat(rows, cols int, values ...int) *CSRMatrix {
	if len(values) != 0 && len(values) != rows*cols {
		panic(fmt.Sprintf("matrix data length (%v) to length mismatch expected %v", len(values), rows*cols))
	}

	mat := CSRMatrix{
		rows: rows,
		cols: cols,
		data: make([][]int, rows),
	}

	for i := 0; i < rows; i++ {
		mat.data[i] = make([]int, 0)
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

func CSRMatFromVec(vec SparseVector) SparseMat {
	m := CSRMat(1, vec.Len())
	m.SetRow(0, vec)
	return m
}

//Identity create an identity matrix (one's on the diagonal).
func CSRIdentity(size int) SparseMat {
	mat := csrMat(size, size)

	for i := 0; i < size; i++ {
		mat.data[i] = append(mat.data[i], i)
	}

	return mat
}

//Copy will create a NEW matrix that will have all the same values as m.
func CSRMatCopy(m SparseMat) SparseMat {
	mat := csrMat(m.Dims())

	for i := 0; i < mat.rows; i++ {
		mat.SetRow(i, m.Row(i))
	}

	return mat
}

//Slice creates a new matrix containing the slice of data.
func (mat *CSRMatrix) Slice(i, j, rows, cols int) SparseMat {
	if rows <= 0 || cols <= 0 {
		panic("slice rows and cols must >= 1")
	}

	mat.checkRowBounds(i)
	mat.checkColBounds(j)
	mat.checkRowBounds(i + rows - 1)
	mat.checkColBounds(j + cols - 1)

	return mat.slice(i, j, rows, cols)
}

func (mat *CSRMatrix) slice(r, c, rows, cols int) *CSRMatrix {
	m := csrMat(rows, cols)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			m.set(i, j, mat.at(i+r, j+c))
		}
	}

	return m
}

func (mat *CSRMatrix) checkRowBounds(i int) {
	if i < 0 || i >= mat.rows {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, mat.rows-1))
	}
}

func (mat *CSRMatrix) checkColBounds(j int) {
	if j < 0 || j >= mat.cols {
		panic(fmt.Sprintf("%v out of range: [0-%v]", j, mat.cols-1))
	}
}

//Dims returns the dimensions of the matrix.
func (mat *CSRMatrix) Dims() (int, int) {
	return mat.rows, mat.cols
}

//At returns the value at row index i and column index j.
func (mat *CSRMatrix) At(i, j int) int {
	mat.checkRowBounds(i)
	mat.checkColBounds(j)

	return mat.at(i, j)
}

func (mat *CSRMatrix) SwapRows(i1, i2 int) {
	mat.checkRowBounds(i1)
	mat.checkRowBounds(i2)

	if i1 == i2 {
		return
	}

	tmp := mat.data[i1]
	mat.data[i1] = mat.data[i2]
	mat.data[i2] = tmp
}

func (mat *CSRMatrix) SwapColumns(j1, j2 int) {
	mat.checkColBounds(j1)
	mat.checkColBounds(j2)

	if j1 > j2 {
		j1, j2 = j2, j1
	}

	for i := 0; i < mat.rows; i++ {
		row := mat.data[i]
		c1 := findIndex(row, j1)
		c2 := findIndex(row, j2)

		rowLen := len(row)
		j1InRow := c1 < rowLen
		j2InRow := c2 < rowLen
		hasj1 := j1InRow && row[c1] == j1
		hasj2 := j2InRow && row[c2] == j2

		if hasj1 == hasj2 {
			continue
		}

		if hasj1 {
			copy(row[c1:c2], row[c1+1:c2])
			row[c2-1] = j2
			continue
		}

		if c1 < c2 {
			copy(row[c1+1:], row[c1:c2])
		}
		row[c1] = j1
	}
}

//AddRows is fast row operation to add two
// rows and put the result in a destination row.
func (mat *CSRMatrix) AddRows(i1, i2, dest int) {
	mat.checkRowBounds(i1)
	mat.checkRowBounds(i2)
	mat.checkRowBounds(dest)

	av := mat.data[i1]
	bv := mat.data[i2]

	mat.data[dest] = addRows(av, bv)
}

func findIndex(indices []int, value int) int {
	il := len(indices)
	return sort.Search(il, func(i int) bool {
		return indices[i] >= value
	})
}

func insertOneElement(s []int, index int, value int) []int {
	s = append(s, 0)
	copy(s[index+1:], s[index:])
	s[index] = value
	return s
}

func cutRange(a []int, start int, end int) []int {
	copy(a[start:], a[end:])
	a = a[:len(a)-(end-start)]
	return a
}

func (mat *CSRMatrix) at(r, c int) int {
	cols := mat.data[r]
	j := findIndex(cols, c)
	if j == len(cols) || cols[j] != c {
		return 0
	}
	return 1
}

//Set sets the value at row index i and column index j to value.
func (mat *CSRMatrix) Set(i, j, value int) {
	mat.checkRowBounds(i)
	mat.checkColBounds(j)

	mat.set(i, j, value%2)
}

func (mat *CSRMatrix) set(r, c, value int) {
	cols := mat.data[r]
	j := findIndex(cols, c)
	if value == 0 {
		if j == len(cols) || cols[j] != c {
			return
		}
		mat.data[r] = cutRange(cols, j, j+1)
		return
	}

	if j < len(cols) && cols[j] == c {
		return
	}

	mat.data[r] = insertOneElement(cols, j, c)
}

//T returns a new matrix that is the transpose of the underlying matrix.
func (mat *CSRMatrix) T() SparseMat {
	m := csrMat(mat.cols, mat.rows)

	for i := 0; i < mat.rows; i++ {
		cols := mat.data[i]
		for _, j := range cols {
			m.set(j, i, 1)
		}
	}

	return m
}

//Zeroize take the current matrix sets all values to 0.
func (mat *CSRMatrix) Zeroize() {
	mat.data = make([][]int, mat.rows)
	for i := 0; i < mat.rows; i++ {
		mat.data[i] = make([]int, 0)
	}
}

//ZeroizeRange take the current matrix sets values inside the range to zero.
func (mat *CSRMatrix) ZeroizeRange(i, j, rows, cols int) {
	if i < 0 || j < 0 || rows < 0 || cols < 0 {
		panic("zeroize must have positive values")
	}
	if mat.rows < i+rows || mat.cols < j+cols {
		panic(fmt.Sprintf("zeroize bounds check failed can't zeroize shape (%v,%v) on a (%v,%v) matrix", i+rows, j+cols, mat.rows, mat.cols))
	}

	mat.zeroize(i, j, rows, cols)
}

func (mat *CSRMatrix) zeroize(r, c, rows, cols int) {
	for i := r; i < r+rows; i++ {
		for j := c; j < c+cols; j++ {
			mat.set(i, j, 0)
		}
	}
}

//Mul multiplies two matrices and stores the values in this matrix.
func (mat *CSRMatrix) Mul(a, b SparseMat) {
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
}

func (mat *CSRMatrix) mul(a, b SparseMat) {
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
func (mat *CSRMatrix) Add(a, b SparseMat) {
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
}

func (mat *CSRMatrix) add(a, b SparseMat) {
	for i := 0; i < mat.rows; i++ {
		mat.data[i] = addRows(a.Row(i).NonzeroArray(), b.Row(i).NonzeroArray())
	}
}

func addRows(av, bv []int) []int {
	avLen := len(av)
	bvLen := len(bv)
	vec := make([]int, 0, avLen+bvLen)

	ai := 0
	bi := 0
	for ai < avLen && bi < bvLen {
		switch {
		case av[ai] < bv[bi]:
			vec = append(vec, av[ai])
			ai++
		case av[ai] > bv[bi]:
			vec = append(vec, bv[bi])
			bi++
		case av[ai] == bv[bi]:
			ai++
			bi++
		}
	}

	for ; ai < avLen; ai++ {
		vec = append(vec, av[ai])
	}
	for ; bi < bvLen; bi++ {
		vec = append(vec, bv[bi])
	}
	return vec
}

//Column returns a map containing the non zero row indices as the keys and it's associated values.
func (mat *CSRMatrix) Column(j int) SparseVector {
	mat.checkColBounds(j)

	indices := make([]int, 0, mat.rows)

	for i := 0; i < mat.rows; i++ {
		row := mat.data[i]
		c := findIndex(row, j)
		if c < len(row) && row[c] == j {
			indices = append(indices, i)
		}
	}

	return &CSRVector{
		length:  mat.rows,
		indices: indices,
	}
}

//SetColumn sets the values in column j. The values' keys are expected to be row indices.
func (mat *CSRMatrix) SetColumn(j int, vec SparseVector) {
	mat.checkColBounds(j)

	if mat.rows != vec.Len() {
		panic("matrix number of columns must equal length of vector")
	}

	for i := 0; i < mat.rows; i++ {
		ii := vec.At(i)
		mat.set(i, j, ii)
	}
}

//Row returns a map containing the non zero column indices as the keys and it's associated values.
func (mat *CSRMatrix) Row(i int) SparseVector {
	mat.checkRowBounds(i)

	row := mat.data[i]
	vec := make([]int, len(row))
	copy(vec, row)

	return &CSRVector{
		length:  mat.cols,
		indices: vec,
	}
}

//SetRow sets the values in row i. The values' keys are expected to be column indices.
func (mat *CSRMatrix) SetRow(i int, vec SparseVector) {
	mat.checkRowBounds(i)

	if mat.cols != vec.Len() {
		panic("matrix number of columns must equal length of vector")
	}

	mat.data[i] = vec.NonzeroArray()
}

//Equals return true if the m matrix has the same shape and values as this matrix.
func (mat *CSRMatrix) Equals(m SparseMat) bool {
	if mat == m {
		return true
	}

	if mat == nil || m == nil {
		return false
	}

	r, c := m.Dims()

	if mat.rows != r || mat.cols != c {
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
func (mat CSRMatrix) String() string {
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
func (mat *CSRMatrix) SetMatrix(a SparseMat, iOffset, jOffset int) {
	if iOffset < 0 || jOffset < 0 {
		panic("offsets must be positive values [0,+)")
	}
	aRows, aCols := a.Dims()
	if mat.rows < iOffset+aRows || mat.cols < jOffset+aCols {
		panic(fmt.Sprintf("set matrix have equal or smaller shape (%v,%v), found a=(%v,%v)", mat.rows, mat.cols, iOffset+aRows, jOffset+aCols))
	}

	mat.setMatrix(a, iOffset, jOffset)
}

func (mat *CSRMatrix) setMatrix(a SparseMat, rOffset, cOffset int) {
	aRows, aCols := a.Dims()
	for i := 0; i < aRows; i++ {
		for j := 0; j < aCols; j++ {
			mat.set(rOffset+i, cOffset+j, a.At(i, j))
		}
	}
}

//Negate performs a piecewise logical negation.
func (mat *CSRMatrix) Negate() {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			mat.set(i, j, (mat.at(i, j)+1)%2)
		}
	}
}

//And executes a piecewise logical AND on the two matrices and stores the values in this matrix.
func (mat *CSRMatrix) And(a, b SparseMat) {
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
}

func (mat *CSRMatrix) and(a, b SparseMat) {
	//first we need to clear mat
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			mat.set(i, j, a.At(i, j)&b.At(i, j))
		}
	}
}

//Or executes a piecewise logical OR on the two matrices and stores the values in this matrix.
func (mat *CSRMatrix) Or(a, b SparseMat) {
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
}

func (mat *CSRMatrix) or(a, b SparseMat) {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			mat.set(i, j, a.At(i, j)|b.At(i, j))
		}
	}
}

//XOr executes a piecewise logical XOR on the two matrices and stores the values in this matrix.
func (mat *CSRMatrix) XOr(a, b SparseMat) {
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
}

func (mat *CSRMatrix) xor(a, b SparseMat) {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			mat.set(i, j, a.At(i, j)^b.At(i, j))
		}
	}
}
