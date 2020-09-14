package mat

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"strings"
)

type Vector struct {
	mat *Matrix
}

func NewVec(length int, values ...int) *Vector {
	if len(values) != 0 {
		if length != len(values) {
			panic("length and number of values must be equal")
		}
	}
	vec := Vector{
		mat: NewMat(1, length, values...),
	}

	return &vec
}

func CopyVec(a *Vector) *Vector {
	return &Vector{
		mat: Copy(a.mat),
	}
}

func (vec *Vector) offset() int {
	return vec.mat.colStart
}

//String returns a string representation of this vector.
func (vec *Vector) String() string {
	buff := &strings.Builder{}
	table := tablewriter.NewWriter(buff)

	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)

	row := make([]string, vec.Len())
	for i := 0; i < vec.Len(); i++ {
		j := i + vec.offset()
		row[i] = fmt.Sprint(vec.at(j))
	}
	table.Append(row)

	table.Render()
	return buff.String()
}

func (vec *Vector) checkBounds(i int) {
	if i < 0 || i >= vec.Len() {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, vec.Len()-1))
	}
}

//At returns the value at index i.
func (vec *Vector) At(i int) int {
	vec.checkBounds(i)
	j := i + vec.offset()
	return vec.at(j)
}

func (vec *Vector) at(j int) int {
	return vec.mat.rowValues[0][j]
}

//Set sets the value at row index i and column index j to value.
func (vec *Vector) Set(i, value int) {
	vec.checkBounds(i)
	j := i + vec.offset()

	vec.set(j, value)
}

func (vec *Vector) set(j, value int) {
	vec.mat.set(0, j, value)
}

func (vec *Vector) Len() int {
	return vec.mat.cols
}

func (vec *Vector) Dot(a *Vector) int {
	m := NewMat(1, 1)
	m.Mul(vec.mat, a.mat.T())
	return m.at(0, 0)
}

func (vec *Vector) T() *TransposedVector {
	return &TransposedVector{
		mat: vec.mat.T(),
	}
}

func (vec *Vector) Add(a, b *Vector) {
	if a == nil || b == nil {
		panic("addition input was found to be nil")
	}
	if vec == a || vec == b {
		panic("addition self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic("adding vectors must have the same length")
	}
	if vec.Len() != a.Len() {
		panic("adding vectors, destination must have the same length")
	}

	vec.mat.add(a.mat, b.mat)
}

func (vec *Vector) Equals(v *Vector) bool {
	return vec.mat.Equals(v.mat)
}

func (vec *Vector) Mul(vec2 *Vector, mat *Matrix) {
	if vec == nil || mat == nil {
		panic("vector multiply input was found to be nil")
	}

	if vec == vec2 || vec.mat == mat {
		panic("vector multiply self assignment not allowed")
	}

	if vec2.mat.cols != mat.rows {
		panic(fmt.Sprintf("multiply shape misalignment can't vector-matrix multiply dims: (%v)x(%v,%v)", vec2.mat.cols, mat.rows, mat.cols))
	}

	_, matCols := mat.Dims()
	if vec.mat.cols != matCols {
		panic(fmt.Sprintf("vector not long enough to hold result, actual length:%v required:%v", vec.Len(), mat.cols))
	}

	vec.mat.mul(vec2.mat, mat)
}

type TransposedVector struct {
	mat *Matrix
}

func NewTVec(length int, values ...int) *TransposedVector {
	if len(values) != 0 {
		if length != len(values) {
			panic("length and number of values must be equal")
		}
	}
	vec := Vector{
		mat: NewMat(1, length, values...),
	}

	return vec.T()
}

func CopyTVec(a *TransposedVector) *TransposedVector {
	return &TransposedVector{
		mat: Copy(a.mat),
	}
}

func (tvec *TransposedVector) checkBounds(i int) {
	if i < 0 || i >= tvec.Len() {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, tvec.Len()-1))
	}
}

func (tvec *TransposedVector) offset() int {
	return tvec.mat.rowStart
}

func (tvec *TransposedVector) T() *Vector {
	return &Vector{
		mat: tvec.mat.T(),
	}
}

func (tvec *TransposedVector) Len() int {
	return tvec.mat.rows
}

func (tvec *TransposedVector) MulVec(a *Matrix, b *TransposedVector) {
	if a == nil || b == nil {
		panic("multiply input was found to be nil")
	}

	if tvec == b || tvec.mat == a {
		panic("multiply self assignment not allowed")
	}

	if a.cols != b.mat.rows {
		panic(fmt.Sprintf("multiply shape misalignment can't matrix-vector multiply (%v,%v)x(%v,1)", a.rows, a.cols, b.mat.rows))
	}

	if tvec.Len() != b.Len() {
		panic(fmt.Sprintf("transposed vector length (%v) does not match expected (%v)", tvec.Len(), b.Len()))
	}

	tvec.mat.mul(a, b.mat)
}

//At returns the value at index i.
func (tvec *TransposedVector) At(j int) int {
	tvec.checkBounds(j)
	i := j + tvec.offset()
	return tvec.at(i)
}

func (tvec *TransposedVector) at(i int) int {
	return tvec.mat.rowValues[i][0]
}

//Set sets the value at row index i and column index j to value.
func (tvec *TransposedVector) Set(j, value int) {
	tvec.checkBounds(j)
	i := j + tvec.offset()

	tvec.set(i, value)
}

func (tvec *TransposedVector) set(i, value int) {
	tvec.mat.set(i, 0, value)
}

func (tvec *TransposedVector) Equals(v *TransposedVector) bool {
	return tvec.mat.Equals(v.mat)
}

func (tvec *TransposedVector) String() string {
	buff := &strings.Builder{}
	table := tablewriter.NewWriter(buff)

	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)

	row := make([]string, tvec.Len())
	for i := 0; i < tvec.Len(); i++ {
		j := i + tvec.offset()
		row[i] = fmt.Sprint(tvec.at(j))
	}
	table.Append(row)

	table.Render()
	return buff.String()
}
