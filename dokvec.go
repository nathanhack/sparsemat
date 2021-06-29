package sparsemat

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type DOKVector struct {
	length int
	values map[int]int
}

type dokVector struct {
	Length int
	Values map[int]int
}

func (vec *DOKVector) MarshalJSON() ([]byte, error) {
	return json.Marshal(dokVector{
		Length: vec.length,
		Values: vec.values,
	})
}

func (vec *DOKVector) UnmarshalJSON(bytes []byte) error {
	var v dokVector
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}
	vec.length = v.Length
	vec.values = v.Values
	return nil
}

func DOKVec(length int, values ...int) SparseVector {
	if len(values) != 0 {
		if length != len(values) {
			panic("length and number of values must be equal")
		}
	}

	return dokVec(length, values...)
}

func dokVec(length int, values ...int) *DOKVector {
	valuesMap := make(map[int]int)
	for i := 0; i < len(values); i++ {
		if values[i] == 1 {
			valuesMap[i] = 1
		}
	}
	vec := DOKVector{
		length: length,
		values: valuesMap,
	}

	return &vec
}

func DOKVecCopy(a SparseVector) SparseVector {
	v := dokVec(a.Len())

	for i := range a.NonzeroMap() {
		v.values[i] = 1
	}
	return v
}

//String returns a string representation of this vector.
func (vec *DOKVector) String() string {
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

func (vec *DOKVector) checkBounds(i int) {
	if i < 0 || i >= vec.Len() {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, vec.Len()-1))
	}
}

//At returns the value at index i.
func (vec *DOKVector) At(i int) int {
	vec.checkBounds(i)
	return vec.at(i)
}

func (vec *DOKVector) at(j int) int {
	if _, has := vec.values[j]; has {
		return 1
	}
	return 0
}

//Set sets the value at row index i and column index j to value.
func (vec *DOKVector) Set(i, value int) {
	vec.checkBounds(i)
	vec.set(i, value%2)
}

func (vec *DOKVector) set(j, value int) {
	if value == 0 {
		_, has := vec.values[j]
		if has {
			delete(vec.values, j)
		}
		return
	}
	vec.values[j] = 1
}

//SetVec replaces the values of this vector with the values of from vector a.
func (vec *DOKVector) SetVec(a SparseVector, i int) {
	vec.checkBounds(i + a.Len() - 1)
	for j := 0; j < a.Len(); j++ {
		vec.set(j+i, a.At(j))
	}
}

func (vec *DOKVector) Len() int {
	return vec.length
}

func (vec *DOKVector) Dot(a SparseVector) int {
	if vec.length != a.Len() {
		panic("lengths must be equal")
	}

	v := 0
	for i := range vec.values {
		v += a.At(i)
	}
	return v % 2
}

func (vec *DOKVector) NonzeroMap() (indexToValues map[int]int) {
	indexToValues = make(map[int]int)
	for i := range vec.values {
		indexToValues[i] = 1
	}
	return indexToValues
}

func (vec *DOKVector) NonzeroArray() (indices []int) {
	indices = make([]int, 0, len(vec.values))
	for i := range vec.values {
		indices = append(indices, i)
	}
	sort.Ints(indices)
	return
}

//Slice creates a slice of the Vector.  The slice will be connected to the original Vector, changes to one
// causes changes in the other.
func (vec *DOKVector) Slice(i, length int) SparseVector {
	if length <= 0 {
		panic("slice len must >0")
	}
	vec.checkBounds(i)
	vec.checkBounds(i + length - 1)

	tmp := make(map[int]int)

	if length < len(vec.values) {
		for j := 0; j < length; j++ {
			if vec.values[j+i] > 0 {
				tmp[j] = 1
			}
		}
	} else {
		for j := range vec.values {
			if i <= j && j < i+length {
				tmp[j-i] = 1
			}
		}
	}

	return &DOKVector{
		length: length,
		values: tmp,
	}
}

func (vec *DOKVector) Add(a, b SparseVector) {
	if a == nil || b == nil {
		panic("addition input was found to be nil")
	}

	if a.Len() != b.Len() {
		panic("adding vectors must have the same length")
	}
	if vec.Len() != a.Len() {
		panic("adding vectors, destination must have the same length")
	}

	for i := 0; i < a.Len(); i++ {
		vec.set(i, (a.At(i)+b.At(i))%2)
	}
}

func (vec *DOKVector) Equals(v SparseVector) bool {
	return vec.length == v.Len() && reflect.DeepEqual(vec.values, v.NonzeroMap())
}

func (vec *DOKVector) MulMat(vec2 SparseVector, mat SparseMat) {
	if vec == nil || mat == nil {
		panic("vector multiply input was found to be nil")
	}
	matRows, matCols := mat.Dims()

	if vec2.Len() != matRows {
		panic(fmt.Sprintf("multiply shape misalignment can't vector-matrix multiply dims: (%v)x(%v,%v)", vec2.Len(), matRows, matCols))
	}

	if vec.length != matCols {
		panic(fmt.Sprintf("vector not long enough to hold result, actual length:%v required:%v", vec.length, matCols))
	}

	for i := 0; i < matCols; i++ {
		vec.set(i, vec2.Dot(mat.Column(i)))
	}
}

func (vec *DOKVector) MatMul(mat SparseMat, vec2 SparseVector) {
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

func (vec *DOKVector) And(a, b SparseVector) {
	if a == nil || b == nil {
		panic("AND input was found to be nil")
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

func (vec *DOKVector) Or(a, b SparseVector) {
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

func (vec *DOKVector) XOr(a, b SparseVector) {
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

func (vec *DOKVector) Negate() SparseVector {
	for i := 0; i < vec.length; i++ {
		_, has := vec.values[i]
		if has {
			delete(vec.values, i)
		} else {
			vec.values[i] = 1
		}
	}
	return vec
}

func (vec *DOKVector) IsZero() bool {
	return len(vec.values) == 0
}

func (vec *DOKVector) HammingWeight() int {
	return len(vec.values)
}

func (vec *DOKVector) HammingDistance(a SparseVector) int {
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
func (vec *DOKVector) NextSet(startingIndex int) (index int, has bool) {
	vec.checkBounds(startingIndex)

	for i := startingIndex; i < vec.length; i++ {
		if vec.at(i) > 0 {
			return i, true
		}
	}
	return -1, false
}
