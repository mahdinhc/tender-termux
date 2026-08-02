package tender

import (
	"fmt"
	"strings"
)

// UMatrixElement constrains allowed numeric types for high-performance matrices.
type UMatrixElement interface {
	~int64 | ~float64 | ~complex128
}

// UMatrix represents a generic 2D dense matrix backed by a contiguous 1D slice.
type UMatrix[T UMatrixElement] struct {
	ObjectImpl
	Rows int
	Cols int
	Data []T
}

// Ensure generic UMatrix types satisfy the Tender Object interface
var (
	_ Object = (*UMatrix[int64])(nil)
	_ Object = (*UMatrix[float64])(nil)
	_ Object = (*UMatrix[complex128])(nil)
)

// NewUMatrix allocates a zero-initialized matrix of shape (rows x cols).
func NewUMatrix[T UMatrixElement](rows, cols int) *UMatrix[T] {
	return &UMatrix[T]{
		Rows: rows,
		Cols: cols,
		Data: make([]T, rows*cols),
	}
}

// NewUMatrixWithData constructs a matrix using an existing contiguous slice.
func NewUMatrixWithData[T UMatrixElement](rows, cols int, data []T) (*UMatrix[T], error) {
	if len(data) != rows*cols {
		return nil, fmt.Errorf("matrix: data size %d does not match shape (%dx%d)", len(data), rows, cols)
	}
	return &UMatrix[T]{
		Rows: rows,
		Cols: cols,
		Data: data,
	}, nil
}

// Get returns the element at row r and column c.
func (m *UMatrix[T]) Get(r, c int) T {
	return m.Data[r*m.Cols+c]
}

// Set updates the element at row r and column c.
func (m *UMatrix[T]) Set(r, c int, val T) {
	m.Data[r*m.Cols+c] = val
}

// TypeName returns the dynamic runtime object type name.
func (m *UMatrix[T]) TypeName() string {
	var zero T
	switch any(zero).(type) {
	case int64:
		return "matrix:int"
	case float64:
		return "matrix:float"
	case complex128:
		return "matrix:complex"
	default:
		return "matrix"
	}
}

// String provides formatted multi-line matrix representation.
func (m *UMatrix[T]) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for r := 0; r < m.Rows; r++ {
		if r > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString("[")
		for c := 0; c < m.Cols; c++ {
			sb.WriteString(fmt.Sprintf("%v", m.Get(r, c)))
			if c < m.Cols-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString("]")
		if r < m.Rows-1 {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

// Copy performs a deep copy of the underlying flat slice.
func (m *UMatrix[T]) Copy() Object {
	newData := make([]T, len(m.Data))
	copy(newData, m.Data)
	return &UMatrix[T]{
		Rows: m.Rows,
		Cols: m.Cols,
		Data: newData,
	}
}

// IsFalsy returns false if the matrix is empty (0 rows or 0 cols).
func (m *UMatrix[T]) IsFalsy() bool {
	return m.Rows == 0 || m.Cols == 0
}