package tender

import (
	"fmt"
	"strings"
	"math"
	"math/cmplx"
	"sync"
	"runtime"
	
	"github.com/2dprototype/tender/token"
)

// MatrixElement constrains allowed numeric types for high-performance matrices.
type MatrixElement interface {
	~int64 | ~float64 | ~complex128
}

// Matrix represents a generic 2D dense matrix backed by a contiguous 1D slice.
type Matrix[T MatrixElement] struct {
	ObjectImpl
	Rows int
	Cols int
	Data []T
}

// Ensure generic Matrix types satisfy the Tender Object interface
var (
	_ Object = (*Matrix[int64])(nil)
	_ Object = (*Matrix[float64])(nil)
	_ Object = (*Matrix[complex128])(nil)
)

// TypeName returns the dynamic runtime object type name.
func (m *Matrix[T]) TypeName() string {
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

// Helper to format float with 2 decimal places, no trailing zeros
func formatFloat(f float64) string {
	// Handle special cases
	if math.IsNaN(f) {
		return "NaN"
	} else if math.IsInf(f, 1) {
		return "Inf"
	} else if math.IsInf(f, -1) {
		return "-Inf"
	}
	
	// Check if it's a whole number
	if math.Abs(f-math.Round(f)) < 1e-9 {
		return fmt.Sprintf("%.0f", f)
	}
	
	// Format with 2 decimal places, but trim trailing zeros
	s := fmt.Sprintf("%.2f", f)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

// formatMatrixValue formats a single matrix element as a string
func formatMatrixValue[T MatrixElement](val T) string {
	var zero T
	var isFloat bool
	var isComplex bool
	
	switch any(zero).(type) {
	case float64:
		isFloat = true
	case complex128:
		isComplex = true
	}
	
	if isComplex {
		c := any(val).(complex128)
		
		// Handle NaN and Inf
		if cmplx.IsNaN(c) {
			return "NaN"
		} else if cmplx.IsInf(c) {
			return "Inf"
		}
		
		realPart := real(c)
		imagPart := imag(c)
		
		// Format real part with 2 decimal places if needed
		realStr := formatFloat(realPart)
		
		// If imaginary part is zero, just return real part
		if math.Abs(imagPart) < 1e-9 {
			return realStr
		}
		
		// Format imaginary part with 2 decimal places if needed
		imagStr := formatFloat(math.Abs(imagPart))
		
		if imagPart < 0 {
			return fmt.Sprintf("%s-%si", realStr, imagStr)
		}
		return fmt.Sprintf("%s+%si", realStr, imagStr)
	}
	
	if isFloat {
		return formatFloat(any(val).(float64))
	}
	
	// int64
	return fmt.Sprintf("%d", any(val).(int64))
}

// String returns a compact representation with angle brackets
func (m *Matrix[T]) String() string {
	if m.Rows == 0 || m.Cols == 0 {
		return "<matrix>"
	}
	
	const maxDisplay = 1000
	var parts []string
	
	if len(m.Data) <= maxDisplay {
		for _, val := range m.Data {
			parts = append(parts, formatMatrixValue(val))
		}
		return "<matrix " + strings.Join(parts, " ") + ">"
	}
	
	// Show first 100 values
	for i := 0; i < maxDisplay; i++ {
		parts = append(parts, formatMatrixValue(m.Data[i]))
	}
	parts = append(parts, "...")
	return "<matrix " + strings.Join(parts, " ") + ">"
}

// PrettyMatrix returns a multi‑line table representation of matrix m,
// optionally with ANSI color codes (if colored is true).
// If the matrix has more than 10000 elements, it returns m.String() instead.
func PrettyMatrix[T MatrixElement](m *Matrix[T], colored bool) string {
	if m.Rows == 0 || m.Cols == 0 {
		return "| |"
	}

	// If matrix is too large, return the compact string representation
	const maxElements = 1000
	if len(m.Data) > maxElements {
		return m.String()
	}

	// Compute column widths
	colWidths := make([]int, m.Cols)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			valStr := formatMatrixValue(m.Data[i*m.Cols+j])
			if len(valStr) > colWidths[j] {
				colWidths[j] = len(valStr)
			}
		}
	}

	var sb strings.Builder
	for i := 0; i < m.Rows; i++ {
		sb.WriteString("|")
		for j := 0; j < m.Cols; j++ {
			val := m.Data[i*m.Cols+j]
			valStr := formatMatrixValue(val)
			if colored {
				sb.WriteString("\033[33m") // yellow for numeric values
				sb.WriteString(valStr)
				sb.WriteString("\033[0m")
			} else {
				sb.WriteString(valStr)
			}
			padding := colWidths[j] - len(valStr)
			if padding > 0 {
				sb.WriteString(strings.Repeat(" ", padding))
			}
			if j < m.Cols-1 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("|")
		if i < m.Rows-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}


// Copy performs a deep copy of the underlying flat slice.
func (m *Matrix[T]) Copy() Object {
	newData := make([]T, len(m.Data))
	copy(newData, m.Data)
	return &Matrix[T]{
		Rows: m.Rows,
		Cols: m.Cols,
		Data: newData,
	}
}

// IsFalsy returns false if the matrix is empty (0 rows or 0 cols).
func (m *Matrix[T]) IsFalsy() bool {
	return len(m.Data) == 0
}

func toFloat64[T MatrixElement](val T) float64 {
	switch v := any(val).(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	case complex128:
		return real(v)
	}
	return 0
}

func toComplex128[T MatrixElement](val T) complex128 {
	switch v := any(val).(type) {
	case int64:
		return complex(float64(v), 0)
	case float64:
		return complex(v, 0)
	case complex128:
		return v
	}
	return 0
}

func scalarToT[T MatrixElement](obj Object) (T, bool) {
	var zero T
	switch any(zero).(type) {
	case int64:
		if v, ok := obj.(*Int); ok {
			return any(v.Value).(T), true
		}
		if v, ok := obj.(*Float); ok {
			return any(int64(v.Value)).(T), true
		}
	case float64:
		if v, ok := obj.(*Int); ok {
			return any(float64(v.Value)).(T), true
		}
		if v, ok := obj.(*Float); ok {
			return any(v.Value).(T), true
		}
	case complex128:
		if v, ok := obj.(*Int); ok {
			return any(complex(float64(v.Value), 0)).(T), true
		}
		if v, ok := obj.(*Float); ok {
			return any(complex(v.Value, 0)).(T), true
		}
		if v, ok := obj.(*Complex); ok {
			return any(v.Value).(T), true
		}
	}
	return zero, false
}

// BinaryOp implements matrix + matrix, - matrix, * matrix, and scalar +, -, *, /.
func (m *Matrix[T]) BinaryOp(op token.Token, rhs Object) (Object, error) {
	// Matrix op Matrix
	if rhsMat, ok := rhs.(*Matrix[T]); ok {
		switch op {
		case token.Add:
			if m.Rows != rhsMat.Rows || m.Cols != rhsMat.Cols {
				return nil, fmt.Errorf("dimension mismatch: cannot add %dx%d and %dx%d", m.Rows, m.Cols, rhsMat.Rows, rhsMat.Cols)
			}
			res := &Matrix[T]{Rows: m.Rows, Cols: m.Cols, Data: make([]T, len(m.Data))}
			for i := range m.Data {
				res.Data[i] = m.Data[i] + rhsMat.Data[i]
			}
			return res, nil

		case token.Sub:
			if m.Rows != rhsMat.Rows || m.Cols != rhsMat.Cols {
				return nil, fmt.Errorf("dimension mismatch: cannot subtract %dx%d and %dx%d", m.Rows, m.Cols, rhsMat.Rows, rhsMat.Cols)
			}
			res := &Matrix[T]{Rows: m.Rows, Cols: m.Cols, Data: make([]T, len(m.Data))}
			for i := range m.Data {
				res.Data[i] = m.Data[i] - rhsMat.Data[i]
			}
			return res, nil

		case token.Mul:
			return mulMatrix(m, rhsMat)

		case token.Quo:
			if m.Rows != rhsMat.Rows || m.Cols != rhsMat.Cols {
				return nil, fmt.Errorf("dimension mismatch: cannot element‑wise divide %dx%d and %dx%d", m.Rows, m.Cols, rhsMat.Rows, rhsMat.Cols)
			}
			res := &Matrix[T]{Rows: m.Rows, Cols: m.Cols, Data: make([]T, len(m.Data))}
			var zero T
			var isFloat bool
			switch any(zero).(type) {
			case float64:
				isFloat = true
			}
			for i := range m.Data {
				if rhsMat.Data[i] == zero {
					if isFloat {
						res.Data[i] = any(math.Inf(1)).(T)
					} else {
						return nil, fmt.Errorf("division by zero matrix element")
					}
				} else {
					res.Data[i] = m.Data[i] / rhsMat.Data[i]
				}
			}
			return res, nil
		}
	}

	// Matrix op Scalar
	scalar, ok := scalarToT[T](rhs)
	if !ok {
		return nil, ErrInvalidOperator
	}

	switch op {
	case token.Add:
		res := &Matrix[T]{Rows: m.Rows, Cols: m.Cols, Data: make([]T, len(m.Data))}
		for i := range m.Data {
			res.Data[i] = m.Data[i] + scalar
		}
		return res, nil
	case token.Sub:
		res := &Matrix[T]{Rows: m.Rows, Cols: m.Cols, Data: make([]T, len(m.Data))}
		for i := range m.Data {
			res.Data[i] = m.Data[i] - scalar
		}
		return res, nil
	case token.Mul:
		res := &Matrix[T]{Rows: m.Rows, Cols: m.Cols, Data: make([]T, len(m.Data))}
		for i := range m.Data {
			res.Data[i] = m.Data[i] * scalar
		}
		return res, nil
	case token.Quo:
		var zero T
		if scalar == zero {
			return nil, fmt.Errorf("division by zero scalar")
		}
		res := &Matrix[T]{Rows: m.Rows, Cols: m.Cols, Data: make([]T, len(m.Data))}
		for i := range m.Data {
			res.Data[i] = m.Data[i] / scalar
		}
		return res, nil
	}

	return nil, ErrInvalidOperator
}

func formatTAsObject[T MatrixElement](val T) Object {
	switch v := any(val).(type) {
	case int64:
		return &Int{Value: v}
	case float64:
		return &Float{Value: v}
	case complex128:
		return &Complex{Value: v}
	}
	return nil
}

// IndexGet handles properties, methods, and row view retrieval.
func (m *Matrix[T]) IndexGet(index Object) (Object, error) {
	if strIdx, ok := index.(*String); ok {
		switch strIdx.Value {
		case "rows":
			return &Int{Value: int64(m.Rows)}, nil
		case "cols":
			return &Int{Value: int64(m.Cols)}, nil
		case "shape":
			return &Tuple{Value: []Object{&Int{Value: int64(m.Rows)}, &Int{Value: int64(m.Cols)}}}, nil
		case "is_square":
			return FromBool(m.Rows == m.Cols), nil
		case "diag":
			if m.Rows != m.Cols {
				return nil, fmt.Errorf("diagonal is only defined for square matrices")
			}
			diag := make([]Object, m.Rows)
			for i := 0; i < m.Rows; i++ {
				diag[i] = formatTAsObject(m.Data[i*m.Cols+i])
			}
			return &Array{Value: diag}, nil
		case "flatten":
			arr := make([]Object, len(m.Data))
			for i, v := range m.Data {
				arr[i] = formatTAsObject(v)
			}
			return &Array{Value: arr}, nil
		case "sum":
			var sum T
			for _, v := range m.Data {
				sum += v
			}
			return formatTAsObject(sum), nil
		case "mean":
			if len(m.Data) == 0 {
				var zero T
				return formatTAsObject(zero), nil
			}
			var sum T
			for _, v := range m.Data {
				sum += v
			}
			var l T
			switch any(l).(type) {
			case int64:
				l = any(int64(len(m.Data))).(T)
			case float64:
				l = any(float64(len(m.Data))).(T)
			case complex128:
				l = any(complex(float64(len(m.Data)), 0)).(T)
			}
			return formatTAsObject(sum / l), nil
		case "min":
			if len(m.Data) == 0 {
				var zero T
				return formatTAsObject(zero), nil
			}
			// Special handling for min/max to avoid complex '<' compile error
			var isComplex bool
			var zero T
			switch any(zero).(type) {
			case complex128:
				isComplex = true
			}
			if isComplex {
				return nil, fmt.Errorf("min is not supported for complex matrices")
			}
			
			// We can't use T with < directly if T includes complex.
			// Workaround: cast to float64 for comparison if we know it's not complex.
			minIdx := 0
			for i, v := range m.Data[1:] {
				if toFloat64(v) < toFloat64(m.Data[minIdx]) {
					minIdx = i + 1
				}
			}
			return formatTAsObject(m.Data[minIdx]), nil

		case "max":
			if len(m.Data) == 0 {
				var zero T
				return formatTAsObject(zero), nil
			}
			var isComplex bool
			var zero T
			switch any(zero).(type) {
			case complex128:
				isComplex = true
			}
			if isComplex {
				return nil, fmt.Errorf("max is not supported for complex matrices")
			}
			
			maxIdx := 0
			for i, v := range m.Data[1:] {
				if toFloat64(v) > toFloat64(m.Data[maxIdx]) {
					maxIdx = i + 1
				}
			}
			return formatTAsObject(m.Data[maxIdx]), nil

		case "rank":
			return &Float{Value: float64(m.rank())}, nil
		case "det":
			if m.Rows != m.Cols {
				return nil, fmt.Errorf("determinant is only defined for square matrices")
			}
			det, err := m.luDet()
			if err != nil {
				return nil, err
			}
			return formatTAsObject(det), nil
		case "trace":
			if m.Rows != m.Cols {
				return nil, fmt.Errorf("trace is only defined for square matrices")
			}
			var sum T
			for i := 0; i < m.Rows; i++ {
				sum += m.Data[i*m.Cols+i]
			}
			return formatTAsObject(sum), nil
		case "T":
			return m.Transpose(), nil
		case "row":
			return &NativeFunction{
				Name: "row",
				Value: func(args ...Object) (Object, error) {
					if len(args) != 1 {
						return nil, ErrWrongNumArguments
					}
					i, ok := ToInt(args[0])
					if !ok {
						return nil, fmt.Errorf("row index must be an integer")
					}
					if i < 0 || i >= m.Rows {
						return nil, fmt.Errorf("row index out of bounds")
					}
					row := make([]Object, m.Cols)
					for j := 0; j < m.Cols; j++ {
						row[j] = formatTAsObject(m.Data[i*m.Cols+j])
					}
					return &Array{Value: row}, nil
				},
			}, nil
		case "col":
			return &NativeFunction{
				Name: "col",
				Value: func(args ...Object) (Object, error) {
					if len(args) != 1 {
						return nil, ErrWrongNumArguments
					}
					j, ok := ToInt(args[0])
					if !ok {
						return nil, fmt.Errorf("column index must be an integer")
					}
					if j < 0 || j >= m.Cols {
						return nil, fmt.Errorf("column index out of bounds")
					}
					col := make([]Object, m.Rows)
					for i := 0; i < m.Rows; i++ {
						col[i] = formatTAsObject(m.Data[i*m.Cols+j])
					}
					return &Array{Value: col}, nil
				},
			}, nil
		}
		return nil, nil
	}

	if intIdx, ok := index.(*Int); ok {
		row := int(intIdx.Value)
		if row < 0 || row >= m.Rows {
			return nil, fmt.Errorf("row index out of bounds")
		}
		return &rowView[T]{matrix: m, row: row}, nil
	}
	return nil, ErrInvalidIndexType
}

func (m *Matrix[T]) IndexSet(index, value Object) error {
	return fmt.Errorf("matrix element assignment must use m[i][j] = val")
}

func mulMatrix[T MatrixElement](m *Matrix[T], rhs *Matrix[T]) (*Matrix[T], error) {
	if m.Cols != rhs.Rows {
		return nil, fmt.Errorf("dimension mismatch: cannot multiply %dx%d and %dx%d", m.Rows, m.Cols, rhs.Rows, rhs.Cols)
	}
	res := &Matrix[T]{Rows: m.Rows, Cols: rhs.Cols, Data: make([]T, m.Rows*rhs.Cols)}

	const blockSize = 64
	numWorkers := runtime.NumCPU()
	if numWorkers < 1 || m.Rows*rhs.Cols < 4096 {
		numWorkers = 1
	}

	if numWorkers == 1 {
		for i0 := 0; i0 < m.Rows; i0 += blockSize {
			iEnd := i0 + blockSize
			if iEnd > m.Rows {
				iEnd = m.Rows
			}
			for k0 := 0; k0 < m.Cols; k0 += blockSize {
				kEnd := k0 + blockSize
				if kEnd > m.Cols {
					kEnd = m.Cols
				}
				for j0 := 0; j0 < rhs.Cols; j0 += blockSize {
					jEnd := j0 + blockSize
					if jEnd > rhs.Cols {
						jEnd = rhs.Cols
					}

					for i := i0; i < iEnd; i++ {
						for k := k0; k < kEnd; k++ {
							temp := m.Data[i*m.Cols+k]
							for j := j0; j < jEnd; j++ {
								res.Data[i*rhs.Cols+j] += temp * rhs.Data[k*rhs.Cols+j]
							}
						}
					}
				}
			}
		}
		return res, nil
	}

	var wg sync.WaitGroup
	rowsPerWorker := (m.Rows + numWorkers - 1) / numWorkers

	for w := 0; w < numWorkers; w++ {
		rStart := w * rowsPerWorker
		rEnd := rStart + rowsPerWorker
		if rStart >= m.Rows {
			break
		}
		if rEnd > m.Rows {
			rEnd = m.Rows
		}

		wg.Add(1)
		go func(rStart, rEnd int) {
			defer wg.Done()
			for i0 := rStart; i0 < rEnd; i0 += blockSize {
				iSubEnd := i0 + blockSize
				if iSubEnd > rEnd {
					iSubEnd = rEnd
				}
				for k0 := 0; k0 < m.Cols; k0 += blockSize {
					kEnd := k0 + blockSize
					if kEnd > m.Cols {
						kEnd = m.Cols
					}
					for j0 := 0; j0 < rhs.Cols; j0 += blockSize {
						jEnd := j0 + blockSize
						if jEnd > rhs.Cols {
							jEnd = rhs.Cols
						}

						for i := i0; i < iSubEnd; i++ {
							for k := k0; k < kEnd; k++ {
								temp := m.Data[i*m.Cols+k]
								for j := j0; j < jEnd; j++ {
									res.Data[i*rhs.Cols+j] += temp * rhs.Data[k*rhs.Cols+j]
								}
							}
						}
					}
				}
			}
		}(rStart, rEnd)
	}
	wg.Wait()
	return res, nil
}

func (m *Matrix[T]) Transpose() *Matrix[T] {
	res := &Matrix[T]{Rows: m.Cols, Cols: m.Rows, Data: make([]T, len(m.Data))}
	const blockSize = 32
	if m.Rows < blockSize || m.Cols < blockSize {
		for i := 0; i < m.Rows; i++ {
			for j := 0; j < m.Cols; j++ {
				res.Data[j*m.Rows+i] = m.Data[i*m.Cols+j]
			}
		}
		return res
	}

	for i0 := 0; i0 < m.Rows; i0 += blockSize {
		iEnd := i0 + blockSize
		if iEnd > m.Rows {
			iEnd = m.Rows
		}
		for j0 := 0; j0 < m.Cols; j0 += blockSize {
			jEnd := j0 + blockSize
			if jEnd > m.Cols {
				jEnd = m.Cols
			}

			for i := i0; i < iEnd; i++ {
				for j := j0; j < jEnd; j++ {
					res.Data[j*m.Rows+i] = m.Data[i*m.Cols+j]
				}
			}
		}
	}
	return res
}

func (m *Matrix[T]) luDet() (T, error) {
	n := m.Rows
	var zero T
	if n == 0 {
		return zero, nil
	}

	var isComplex bool
	switch any(zero).(type) {
	case complex128:
		isComplex = true
	}

	if !isComplex {
		a := make([]float64, n*n)
		for i := 0; i < len(m.Data); i++ {
			a[i] = toFloat64(m.Data[i])
		}
		det := 1.0
		for i := 0; i < n; i++ {
			pivot := i
			maxVal := math.Abs(a[i*n+i])
			for k := i + 1; k < n; k++ {
				if abs := math.Abs(a[k*n+i]); abs > maxVal {
					maxVal = abs
					pivot = k
				}
			}
			if maxVal < 1e-15 {
				return zero, nil
			}
			if pivot != i {
				for col := 0; col < n; col++ {
					a[i*n+col], a[pivot*n+col] = a[pivot*n+col], a[i*n+col]
				}
				det = -det
			}
			for j := i + 1; j < n; j++ {
				factor := a[j*n+i] / a[i*n+i]
				for k := i; k < n; k++ {
					a[j*n+k] -= factor * a[i*n+k]
				}
			}
			det *= a[i*n+i]
		}
		var typedDet T
		switch any(zero).(type) {
		case int64:
			typedDet = any(int64(math.Round(det))).(T)
		case float64:
			typedDet = any(det).(T)
		}
		return typedDet, nil
	}

	a := make([]complex128, n*n)
	for i := 0; i < len(m.Data); i++ {
		a[i] = toComplex128(m.Data[i])
	}
	det := complex128(1.0)
	for i := 0; i < n; i++ {
		pivot := i
		maxVal := cmplx.Abs(a[i*n+i])
		for k := i + 1; k < n; k++ {
			if abs := cmplx.Abs(a[k*n+i]); abs > maxVal {
				maxVal = abs
				pivot = k
			}
		}
		if maxVal < 1e-15 {
			return zero, nil
		}
		if pivot != i {
			for col := 0; col < n; col++ {
				a[i*n+col], a[pivot*n+col] = a[pivot*n+col], a[i*n+col]
			}
			det = -det
		}
		for j := i + 1; j < n; j++ {
			factor := a[j*n+i] / a[i*n+i]
			for k := i; k < n; k++ {
				a[j*n+k] -= factor * a[i*n+k]
			}
		}
		det *= a[i*n+i]
	}
	return any(det).(T), nil
}

func (m *Matrix[T]) rank() int {
	if m.Rows == 0 || m.Cols == 0 {
		return 0
	}
	rows, cols := m.Rows, m.Cols
	var zero T
	var isComplex bool
	switch any(zero).(type) {
	case complex128:
		isComplex = true
	}

	if !isComplex {
		a := make([]float64, rows*cols)
		for i := 0; i < len(m.Data); i++ {
			a[i] = toFloat64(m.Data[i])
		}
		rank := 0
		for col := 0; col < cols; col++ {
			pivot := -1
			for row := rank; row < rows; row++ {
				if math.Abs(a[row*cols+col]) > 1e-12 {
					pivot = row
					break
				}
			}
			if pivot == -1 {
				continue
			}
			if pivot != rank {
				for c := 0; c < cols; c++ {
					a[rank*cols+c], a[pivot*cols+c] = a[pivot*cols+c], a[rank*cols+c]
				}
			}
			for row := rank + 1; row < rows; row++ {
				factor := a[row*cols+col] / a[rank*cols+col]
				for c := col; c < cols; c++ {
					a[row*cols+c] -= factor * a[rank*cols+c]
				}
			}
			rank++
			if rank == rows {
				break
			}
		}
		return rank
	}

	a := make([]complex128, rows*cols)
	for i := 0; i < len(m.Data); i++ {
		a[i] = toComplex128(m.Data[i])
	}
	rank := 0
	for col := 0; col < cols; col++ {
		pivot := -1
		for row := rank; row < rows; row++ {
			if cmplx.Abs(a[row*cols+col]) > 1e-12 {
				pivot = row
				break
			}
		}
		if pivot == -1 {
			continue
		}
		if pivot != rank {
			for c := 0; c < cols; c++ {
				a[rank*cols+c], a[pivot*cols+c] = a[pivot*cols+c], a[rank*cols+c]
			}
		}
		for row := rank + 1; row < rows; row++ {
			factor := a[row*cols+col] / a[rank*cols+col]
			for c := col; c < cols; c++ {
				a[row*cols+c] -= factor * a[rank*cols+c]
			}
		}
		rank++
		if rank == rows {
			break
		}
	}
	return rank
}

func (m *Matrix[T]) Determinant() (Object, error) {
	return m.detFromIndex()
}

func (m *Matrix[T]) detFromIndex() (Object, error) {
	if m.Rows != m.Cols {
		return nil, fmt.Errorf("determinant is only defined for square matrices")
	}
	det, err := m.luDet()
	if err != nil {
		return nil, err
	}
	return formatTAsObject(det), nil
}

func (m *Matrix[T]) Trace() (Object, error) {
	if m.Rows != m.Cols {
		return nil, fmt.Errorf("trace is only defined for square matrices")
	}
	var sum T
	for i := 0; i < m.Rows; i++ {
		sum += m.Data[i*m.Cols+i]
	}
	return formatTAsObject(sum), nil
}

type rowView[T MatrixElement] struct {
	ObjectImpl
	matrix *Matrix[T]
	row    int
}

func (rv *rowView[T]) TypeName() string { return "row-view" }

func (rv *rowView[T]) String() string {
	rowData := make([]string, rv.matrix.Cols)
	for j := 0; j < rv.matrix.Cols; j++ {
		rowData[j] = fmt.Sprintf("%v", rv.matrix.Data[rv.row*rv.matrix.Cols+j])
	}
	return "[" + strings.Join(rowData, ", ") + "]"
}

func (rv *rowView[T]) Copy() Object {
	return &rowView[T]{matrix: rv.matrix, row: rv.row}
}

func (rv *rowView[T]) Equals(another Object) bool {
	return false
}

func (rv *rowView[T]) IndexGet(index Object) (Object, error) {
	if intIdx, ok := index.(*Int); ok {
		col := int(intIdx.Value)
		if col < 0 || col >= rv.matrix.Cols {
			return nil, fmt.Errorf("column index out of bounds")
		}
		return formatTAsObject(rv.matrix.Data[rv.row*rv.matrix.Cols+col]), nil
	}
	return nil, ErrInvalidIndexType
}

func (rv *rowView[T]) IndexSet(index, value Object) error {
	intIdx, ok := index.(*Int)
	if !ok {
		return ErrInvalidIndexType
	}
	col := int(intIdx.Value)
	if col < 0 || col >= rv.matrix.Cols {
		return fmt.Errorf("column index out of bounds")
	}
	val, ok := scalarToT[T](value)
	if !ok {
		return fmt.Errorf("can only assign numeric values to matrix elements")
	}
	rv.matrix.Data[rv.row*rv.matrix.Cols+col] = val
	return nil
}