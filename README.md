# sparsemat

A simple sparse matrix for values over binary Galois Field - GF(2).

Currently implemented:

* [Dictionary of Keys (DOK)](https://en.wikipedia.org/wiki/Sparse_matrix#Dictionary_of_keys_(DOK)) : instead of using (
  row,col) pair in a map it's broken into a map of maps.
* [Compressed Sparse Row (CSR)](https://en.wikipedia.org/wiki/Sparse_matrix#Compressed_sparse_row_(CSR,_CRS_or_Yale_format))
  ,

Future implementations:

* [Coordinate (COO)](https://en.wikipedia.org/wiki/Sparse_matrix#Coordinate_list_(COO))
* [Compressed Sparse Column (CSC)](https://en.wikipedia.org/wiki/Sparse_matrix#Compressed_sparse_column_(CSC_or_CCS))
* [Diagonal (DIA)](https://en.wikipedia.org/wiki/Sparse_matrix#Diagonal)

### Usage

Importing:

```
import mat "github.com/nathanhack/sparsemat"
```

Creating a matrix:

```
m1 := mat.CSRIdentity(3)
m2 := mat.DOKMat(3,3, 1,0,0,0,1,0,0,0,1)
```

Multiplying: `Mul`

```
r := mat.CSRMat(3,3)
r.Mul(m1,m2) // multiplies m1xm2 and stores into r
```

Accessors: `At` and `Set`

```
fmt.Printf("value at (%v,%v) is %v\n",1,1,r.At(1,1))
r.Set(1,1,0)
fmt.Printf("value at (%v,%v) is now %v\n",1,1,r.At(1,1))
```

Slices: `Slice`

```
m := mat.DOKMat(1,4, 1,1,1,1) // creates matrix [1 1 1 1]
s := m.Slice(0,1,1,2) // creates a slice (new matrix) of the two middle 1's [1 1]
```

Transposes: `T`

```
t := m.T() // note this creates a new allocated matrix
```

## Deps

#### Golang

Version 1.14+

#### Libs

[github.com/olekukonko/tablewriter](github.com/olekukonko/tablewriter)



