package sparsemat

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type CSRVector struct {
	length  int
	indices []int
}

type csrVector struct {
	Length  int
	Indices []int
}

func (vec *CSRVector) MarshalJSON() ([]byte, error) {
	return json.Marshal(csrVector{
		Length:  vec.length,
		Indices: vec.indices,
	})
}

func (vec *CSRVector) UnmarshalJSON(bytes []byte) error {
	var v csrVector
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}
	vec.length = v.Length
	vec.indices = v.Indices
	return nil
}

func CSRVec(length int, values ...int) SparseVector {
	if len(values) != 0 {
		if length != len(values) {
			panic("length and number of values must be equal")
		}
	}
	indices := make([]int, 0)

	for i, v := range values {
		if v != 0 {
			indices = append(indices, i)
		}
	}

	return &CSRVector{
		length:  length,
		indices: indices,
	}
}

func CSRVecCopy(a SparseVector) SparseVector {
	return &CSRVector{
		length:  a.Len(),
		indices: a.NonzeroArray(),
	}
}

//String returns a string representation of this vector.
func (vec *CSRVector) String() string {
	buff := &strings.Builder{}
	table := tablewriter.NewWriter(buff)

	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)

	row := make([]string, vec.Len())
	for i := 0; i < vec.Len(); i++ {
		row[i] = fmt.Sprint(vec.at(i))
	}
	table.Append(row)

	table.Render()
	return buff.String()
}

func (vec *CSRVector) checkBounds(i int) {
	if i < 0 || i >= vec.length {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, vec.Len()-1))
	}
}

//At returns the value at index i.
func (vec *CSRVector) At(i int) int {
	vec.checkBounds(i)
	return vec.at(i)
}

func (vec *CSRVector) at(j int) int {
	x := findIndex(vec.indices, j)
	if x == len(vec.indices) || vec.indices[x] != j {
		return 0
	}

	return 1
}

//Set sets the value at row index i and column index j to value.
func (vec *CSRVector) Set(i, value int) {
	vec.checkBounds(i)
	vec.set(i, value%2)
}

func (vec *CSRVector) set(j, value int) {
	x := findIndex(vec.indices, j)
	indiceLen := len(vec.indices)
	// if value is zero we remove it from the structure
	if value == 0 {
		// if there was a value there then remove it
		if x != indiceLen && vec.indices[x] == j {
			for i := x; i < indiceLen-1; i++ {
				vec.indices[i] = vec.indices[i+1]
			}
			vec.indices = vec.indices[:indiceLen-1]
		}
		return
	}
	//else the value should be in the structure
	if 0 == indiceLen || x == indiceLen || vec.indices[x] != j {
		//if there isn't a value already there we add it
		vec.indices = append(vec.indices, 0)
		for i := indiceLen; i > x; i-- {
			vec.indices[i] = vec.indices[i-1]
		}
		vec.indices[x] = j
	}
}

//SetVec replaces the values of this vector with the values of from vector a.
func (vec *CSRVector) SetVec(a SparseVector, i int) {
	vec.checkBounds(i)
	vec.checkBounds(a.Len() - 1 + i)

	for ii := 0; ii < a.Len(); ii++ {
		vec.set(ii+i, a.At(ii))
	}
}

func (vec *CSRVector) Len() int {
	return vec.length
}

func (vec *CSRVector) Dot(a SparseVector) int {
	if vec.length != a.Len() {
		panic("lengths must be equal")
	}

	v := 0
	indices := a.NonzeroArray()
	vecLen := len(vec.indices)
	aLen := len(indices)
	for i, j := 0, 0; i < vecLen && j < aLen; {
		if vec.indices[i] == indices[j] {
			v++
			i++
			j++
		} else if vec.indices[i] < indices[j] {
			i++
		} else {
			j++
		}
	}

	return v % 2
}

func (vec *CSRVector) NonzeroMap() (indexToValues map[int]int) {
	indexToValues = make(map[int]int)

	for _, r := range vec.indices {
		indexToValues[r] = 1
	}
	return
}

func (vec *CSRVector) NonzeroArray() (indices []int) {
	indices = make([]int, len(vec.indices))
	copy(indices, vec.indices)
	return
}

//Slice creates a slice of the Vector.  The slice will be connected to the original Vector, changes to one
// causes changes in the other.
func (vec *CSRVector) Slice(i, length int) SparseVector {
	if length <= 0 {
		panic("slice len must >0")
	}

	vec.checkBounds(i)
	vec.checkBounds(i + length - 1)

	v := &CSRVector{
		length:  length,
		indices: make([]int, 0),
	}

	for ii := 0; ii < length; ii++ {
		v.set(ii, vec.at(i+ii))
	}
	return v
}

func (vec *CSRVector) Add(a, b SparseVector) {
	if a == nil || b == nil {
		panic("addition input was found to be nil")
	}
	if a.Len() != b.Len() {
		panic("adding vectors must have the same length")
	}
	if vec.Len() != a.Len() {
		panic("adding vectors, destination must have the same length")
	}

	for i := 0; i < vec.length; i++ {
		vec.set(i, (a.At(i)+b.At(i))%2)
	}
}

func (vec *CSRVector) Equals(v SparseVector) bool {
	if vec.length != v.Len() {
		return false
	}

	for i := 0; i < vec.length; i++ {
		if vec.at(i) != v.At(i) {
			return false
		}
	}
	return true
}

func (vec *CSRVector) MulMat(vec2 SparseVector, mat SparseMat) {
	if vec == nil || vec2 == nil || mat == nil {
		panic("vector multiply input was found to be nil")
	}
	matRows, matCols := mat.Dims()
	if vec2.Len() != matRows {
		panic(fmt.Sprintf("multiply shape misalignment can't vector-matrix multiply dims: (%v)x(%v,%v)", vec2.Len(), matRows, matCols))
	}

	if vec.length != matCols {
		panic(fmt.Sprintf("vector not long enough to hold result, actual length:%v required:%v", vec.Len(), matCols))
	}

	for c := 0; c < matCols; c++ {
		vec.set(c, vec2.Dot(mat.Column(c)))
	}
}

func (vec *CSRVector) MatMul(mat SparseMat, vec2 SparseVector) {
	if vec == nil || vec2 == nil || mat == nil {
		panic("vector multiply input was found to be nil")
	}
	matRows, matCols := mat.Dims()
	if vec2.Len() != matCols {
		panic(fmt.Sprintf("multiply shape misalignment can't vector-matrix multiply dims: (%v,%v)x(%v)", matRows, matCols, vec2.Len()))
	}

	if vec.length != matRows {
		panic(fmt.Sprintf("vector not long enough to hold result, actual length:%v required:%v", vec.Len(), matRows))
	}

	for r := 0; r < matRows; r++ {
		vec.set(r, vec2.Dot(mat.Row(r)))
	}
}

func (vec *CSRVector) And(a, b SparseVector) {
	if a == nil || b == nil {
		panic("AND input was found to be nil")
	}

	if vec == a || vec == b {
		panic("AND self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic(fmt.Sprintf("AND shape misalignment both inputs must be equal length found  %v and %v", a.Len(), b.Len()))
	}

	if vec.Len() != a.Len() {
		panic(fmt.Sprintf("vec len:%v does not match expected %v", vec.Len(), a.Len()))
	}

	for i := 0; i < vec.length; i++ {
		vec.set(i, a.At(i)&b.At(i))
	}
}

func (vec *CSRVector) Or(a, b SparseVector) {
	if a == nil || b == nil {
		panic("OR input was found to be nil")
	}

	if vec == a || vec == b {
		panic("OR self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic(fmt.Sprintf("OR shape misalignment both inputs must be equal length found  %v and %v", a.Len(), b.Len()))
	}

	if vec.Len() != a.Len() {
		panic(fmt.Sprintf("vec len:%v does not match expected %v", vec.Len(), a.Len()))
	}

	for i := 0; i < vec.length; i++ {
		vec.set(i, a.At(i)|b.At(i))
	}
}

func (vec *CSRVector) XOr(a, b SparseVector) {
	if a == nil || b == nil {
		panic("XOR input was found to be nil")
	}

	if vec == a || vec == b {
		panic("XOR self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic(fmt.Sprintf("XOR shape misalignment both inputs must be equal length found  %v and %v", a.Len(), b.Len()))
	}

	if vec.Len() != a.Len() {
		panic(fmt.Sprintf("vec len:%v does not match expected %v", vec.Len(), a.Len()))
	}

	for i := 0; i < vec.length; i++ {
		vec.set(i, a.At(i)^b.At(i))
	}
}

func (vec *CSRVector) Negate() {
	indices := make([]int, 0)

	in := make(map[int]bool)
	for i := 0; i < len(vec.indices); i++ {
		in[vec.indices[i]] = true
	}
	for i := 0; i < vec.length; i++ {
		_, has := in[i]
		if !has {
			indices = append(indices, i)
		}
	}
	vec.indices = indices
}

func (vec *CSRVector) IsZero() bool {
	return len(vec.indices) == 0
}

func (vec *CSRVector) HammingWeight() int {
	return len(vec.indices)
}

func (vec *CSRVector) HammingDistance(a SparseVector) int {
	if vec.length != a.Len() {
		panic("HammingDistance() vectors must be same length")
	}

	count := 0
	for i := 0; i < vec.length; i++ {
		if vec.at(i) != a.At(i) {
			count++
		}
	}
	return count
}

//NextSet returns the next bit which is set starting from startingIndex, so if
// the startingIndex is set it will be returned, if not it will be the next bit.
// If no bits are found has bool will be set to false.
func (vec *CSRVector) NextSet(startingIndex int) (index int, has bool) {
	vec.checkBounds(startingIndex)

	index = findIndex(vec.indices, startingIndex)
	if index == len(vec.indices) {
		return -1, false
	}

	return vec.indices[index], true
}
