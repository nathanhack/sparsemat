package csr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Vector struct {
	length  int
	indices []int
}

type vector struct {
	Length  int
	Indices []int
}

func (vec *Vector) MarshalJSON() ([]byte, error) {
	return json.Marshal(vector{
		Length:  vec.length,
		Indices: vec.indices,
	})
}

func (vec *Vector) UnmarshalJSON(bytes []byte) error {
	var v vector
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}
	vec.length = v.Length
	vec.indices = v.Indices
	return nil
}

func NewVec(length int, values ...int) *Vector {
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

	return &Vector{
		length:  length,
		indices: indices,
	}
}

//func CopyVec(a *Vector) *Vector {
//
//	indices := make([]int, len(a.indices))
//	copy(indices, a.indices)
//
//	return &Vector{
//		length:  a.length,
//		indices: indices,
//	}
//}

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
		row[i] = fmt.Sprint(vec.at(i))
	}
	table.Append(row)

	table.Render()
	return buff.String()
}

func (vec *Vector) checkBounds(i int) {
	if i < 0 || i >= vec.length {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, vec.Len()-1))
	}
}

//At returns the value at index i.
func (vec *Vector) At(i int) int {
	vec.checkBounds(i)
	return vec.at(i)
}

func (vec *Vector) at(j int) int {
	x := findIndex(vec.indices, j)
	if x == len(vec.indices) || vec.indices[x] != j {
		return 0
	}

	return 1
}

//Set sets the value at row index i and column index j to value.
func (vec *Vector) Set(i, value int) {
	vec.checkBounds(i)
	vec.set(i, value)
}

func (vec *Vector) set(j, value int) {
	x := findIndex(vec.indices, j)
	// if value is zero we remove it from the structure
	if value%2 == 0 {
		// if there was a value there then remove it
		if x != len(vec.indices) && vec.indices[x] == j {
			vec.indices = cutRange(vec.indices, x, x+1)
		}
		return
	}
	//else the value should be in the structure
	if 0 == len(vec.indices) || x == len(vec.indices) || vec.indices[x] != j {
		//if there isn't a value already there
		vec.indices = insertOneElement(vec.indices, x, j)
	}
}

//SetVec replaces the values of this vector with the values of from vector a.
func (vec *Vector) SetVec(a *Vector, i int) {
	vec.checkBounds(i)
	vec.checkBounds(a.length + i)

	for ii := 0; ii < a.length; ii++ {
		vec.set(ii+i, a.at(ii))
	}
}

func (vec *Vector) Len() int {
	return vec.length
}

func (vec *Vector) Dot(a *Vector) int {
	min := vec.length
	if min > a.length {
		min = a.length
	}
	sum := 0
	for i := 0; i < min; i++ {
		j := vec.at(i)
		k := a.at(i)
		sum += j * k
	}
	return sum % 2
}

func (vec *Vector) NonzeroValues() (indexToValues map[int]int) {
	indexToValues = make(map[int]int)

	for _, r := range vec.indices {
		indexToValues[r] = 1
	}
	return
}

//Slice creates a slice of the Vector.  The slice will be connected to the original Vector, changes to one
// causes changes in the other.
func (vec *Vector) Slice(i, length int) *Vector {
	if length <= 0 {
		panic("slice len must >0")
	}

	vec.checkBounds(i)
	vec.checkBounds(i + length)

	v := &Vector{
		length:  length,
		indices: make([]int, 0),
	}

	for ii := 0; ii < length; ii++ {
		v.set(ii, vec.at(i+ii))
	}
	return v
}

func (vec *Vector) Add(a, b *Vector) {
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
		vec.set(i, a.at(i)+b.at(i))
	}
}

func (vec *Vector) Equals(v *Vector) bool {
	return vec.length == v.length && reflect.DeepEqual(vec.indices, v.indices)
}

func (vec *Vector) Mul(vec2 *Vector, mat *Matrix) {
	if vec == nil || vec2 == nil || mat == nil {
		panic("vector multiply input was found to be nil")
	}

	if vec2.length != mat.rows {
		panic(fmt.Sprintf("multiply shape misalignment can't vector-matrix multiply dims: (%v)x(%v,%v)", vec2.length, mat.rows, mat.cols))
	}

	_, matCols := mat.Dims()
	if vec.length != matCols {
		panic(fmt.Sprintf("vector not long enough to hold result, actual length:%v required:%v", vec.Len(), mat.cols))
	}

	for c := 0; c < matCols; c++ {
		vec.set(c, vec2.Dot(mat.Column(c)))
	}
}

func (vec *Vector) And(a, b *Vector) {
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
		vec.set(i, a.at(i)&b.at(i))
	}
}

func (vec *Vector) Or(a, b *Vector) {
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
		vec.set(i, a.at(i)|b.at(i))
	}
}

func (vec *Vector) XOr(a, b *Vector) {
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
		vec.set(i, a.at(i)^b.at(i))
	}
}

func (vec *Vector) Negate() {
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
