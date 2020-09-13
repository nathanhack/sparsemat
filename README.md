# sparsemat
A simple sparse matrix for values over binary Galois Field - GF(2).

The library simple to use, fairly memory efficient, and fast - if the matrices are indeed sparse. Of the commonly sparse matrix library types: [Dictionary of Keys (DOK)](https://en.wikipedia.org/wiki/Sparse_matrix#Dictionary_of_keys_(DOK)), [Coordinate (COO)](https://en.wikipedia.org/wiki/Sparse_matrix#Coordinate_list_(COO)), [Compressed Sparse Row (CSR)](https://en.wikipedia.org/wiki/Sparse_matrix#Compressed_sparse_row_(CSR,_CRS_or_Yale_format)), [Compressed Sparse Column (CSC)](https://en.wikipedia.org/wiki/Sparse_matrix#Compressed_sparse_column_(CSC_or_CCS)), and [Diagonal (DIA)](https://en.wikipedia.org/wiki/Sparse_matrix#Diagonal). This library most closely resembles the DOK, but instead of use (row,col) pair in a map it's broken into a map of maps. This can lead to some faster matrix operation than DOKs but not as fast as CSRs or CSCs.   

Slices of matrices, and Transposes (like slices of arrays in golang) are "backed" by the matrix that was sliced or trasnposed.  Changes to the sliced or transposed matrix will affect the original matrix. Likewise, changes to the original will affect slices or transposes of it (see below in the usage for Slice).  

### Usage
Importing:
```
import "github.com/nathanhack/sparsemat"
```

Creating a matrix:
```
m1 := sparsemat.Identity(3)
m2 := sparsemat.New(3,3, 1,0,0,0,1,0,0,0,1)
```

Multiplying: `Mul`
```
r := New(3,3)
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
m := sparsemat.New(1,4, 1,1,1,1) // creates matrix [1 1 1 1]
s := m.Slice(0,1,1,2) // creates a slice of the two middle 1's [1 1]

//Note: a slice is "backed" by matrix it was sliced from
s.Set(0,0,0) // sets s to [0 1] and sets m to [1 0 1 1]

m.Set(0,2,0) // sets m to [1 0 0 1] and sets s to [0 0]

s2 := sparsemat.Copy(s) // makes a New matrix that will have the same shape and values of s 
```

Transposes: `T`
```
t := m.T() // create a matrix that is the transpose of m but is "backed" by m, so changes to one affect each other

t.Set(0,0,0)
// this sets t to 
//  |0|
//  |0|
//  |0|
//  |1|
//
// and m to [0 0 0 1]

```





## Deps
#### Golang
Version 1.14+
#### Libs
[github.com/olekukonko/tablewriter](github.com/olekukonko/tablewriter)



