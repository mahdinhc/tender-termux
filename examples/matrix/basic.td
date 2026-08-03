// tender@v0.0.8

// Create matrices
m1 := matrix([
	[1, 2, 3], 
	[4, 5, 6],
	[4, 5, 6]
])

m2 := matrix([
	[7,  8,  1], 
	[9,  10, 2], 
	[11, 12, 2]
])

// Matrix multiplication
m3 := m1 * m2
debug(m3)

// Element-wise operations
m4 := m1 + 10  // Add scalar to every element
debug(m4)

// Matrix properties
rows := m1.rows
cols := m1.cols
shape := m1.shape
debug(rows, cols, shape)

// Transpose
m1t := m1.T
debug(m1t)

// For square matrices
det := m1.det      // determinant
trace := m1.trace  // trace
diag := m1.diag    // diagonal elements
debug(det, trace, diag)

// Matrix methods
row0 := m1.row(0)   // Get first row as array
col1 := m1.col(1)   // Get second column as array
flat := m1.flatten  // Flatten to array
debug(row0, col1, flat)

// Type conversion
m_int := matrix(3, 3, "int", [1, 2, 3, 4, 5, 6, 7, 8, 9])

m_float := m_int.to_float()      // Convert to float matrix
m_complex := m_int.to_complex()  // Convert to complex matrix
debug(m_float |> typeof, m_complex |> typeof)