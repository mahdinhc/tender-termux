package stdlib

import (
	"fmt"
	"math/cmplx"

	"github.com/2dprototype/tender"
)

var cmplxModule = map[string]tender.Object{
	"new":         &tender.UserFunction{Name: "new", Value: cmplxNew},
	"conj":        &tender.UserFunction{Name: "conj", Value: cmplxConj},
	"abs":         &tender.UserFunction{Name: "abs", Value: cmplxAbs},
	"arg":         &tender.UserFunction{Name: "arg", Value: cmplxArg},
	"sin":         &tender.UserFunction{Name: "sin", Value: cmplxSin},
	"cos":         &tender.UserFunction{Name: "cos", Value: cmplxCos},
	"acos":        &tender.UserFunction{Name: "acos", Value: cmplxAcos},
	"acosh":       &tender.UserFunction{Name: "acosh", Value: cmplxAcosh},
	"asin":        &tender.UserFunction{Name: "asin", Value: cmplxAsin},
	"asinh":       &tender.UserFunction{Name: "asinh", Value: cmplxAsinh},
	"atan":        &tender.UserFunction{Name: "atan", Value: cmplxAtan},
	"atanh":       &tender.UserFunction{Name: "atanh", Value: cmplxAtanh},
	"cosh":        &tender.UserFunction{Name: "cosh", Value: cmplxCosh},
	"cot":         &tender.UserFunction{Name: "cot", Value: cmplxCot},
	"exp":         &tender.UserFunction{Name: "exp", Value: cmplxExp},
	"inf":         &tender.UserFunction{Name: "inf", Value: cmplxInf},
	"isinf":       &tender.UserFunction{Name: "isinf", Value: cmplxIsInf},
	"isnan":       &tender.UserFunction{Name: "isnan", Value: cmplxIsNaN},
	"log":         &tender.UserFunction{Name: "log", Value: cmplxLog},
	"log10":       &tender.UserFunction{Name: "log10", Value: cmplxLog10},
	"nan":         &tender.UserFunction{Name: "nan", Value: cmplxNaN},
	"phase":       &tender.UserFunction{Name: "phase", Value: cmplxArg}, // alias for arg
	"polar":       &tender.UserFunction{Name: "polar", Value: cmplxPolar},
	"pow":         &tender.UserFunction{Name: "pow", Value: cmplxPow},
	"rect":        &tender.UserFunction{Name: "rect", Value: cmplxRect},
	"sinh":        &tender.UserFunction{Name: "sinh", Value: cmplxSinh},
	"sqrt":        &tender.UserFunction{Name: "sqrt", Value: cmplxSqrt},
	"tan":         &tender.UserFunction{Name: "tan", Value: cmplxTan},
	"tanh":        &tender.UserFunction{Name: "tanh", Value: cmplxTanh},
}

// cmplxNew creates a new complex number from real and imaginary parts.
func cmplxNew(args ...tender.Object) (tender.Object, error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}
	realPart, ok1 := tender.ToFloat64(args[0])
	imagPart, ok2 := tender.ToFloat64(args[1])
	if !ok1 || !ok2 {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "real, imag",
			Expected: "numbers",
			Found:    fmt.Sprintf("%s, %s", args[0].TypeName(), args[1].TypeName()),
		}
	}
	return &tender.Complex{Value: complex(realPart, imagPart)}, nil
}

// cmplxConjugate returns the conjugate of the complex number.
func cmplxConj(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Conj(c.Value)}, nil
}

// cmplxAbs returns the modulus of the complex number.
func cmplxAbs(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Float{Value: cmplx.Abs(c.Value)}, nil
}

// cmplxArg returns the phase (argument) of the complex number.
func cmplxArg(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Float{Value: cmplx.Phase(c.Value)}, nil
}

// cmplxSin returns the sine of the complex number.
func cmplxSin(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Sin(c.Value)}, nil
}

// cmplxCos returns the cosine of the complex number.
func cmplxCos(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Cos(c.Value)}, nil
}

// cmplxAcos returns the arc-cosine of the complex number.
func cmplxAcos(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Acos(c.Value)}, nil
}

// cmplxAcosh returns the hyperbolic arc-cosine of the complex number.
func cmplxAcosh(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Acosh(c.Value)}, nil
}

// cmplxAsin returns the arc-sine of the complex number.
func cmplxAsin(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Asin(c.Value)}, nil
}

// cmplxAsinh returns the hyperbolic arc-sine of the complex number.
func cmplxAsinh(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Asinh(c.Value)}, nil
}

// cmplxAtan returns the arc-tangent of the complex number.
func cmplxAtan(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Atan(c.Value)}, nil
}

// cmplxAtanh returns the hyperbolic arc-tangent of the complex number.
func cmplxAtanh(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Atanh(c.Value)}, nil
}

// cmplxCosh returns the hyperbolic cosine of the complex number.
func cmplxCosh(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Cosh(c.Value)}, nil
}

// cmplxCot returns the cotangent of the complex number.
// Defined as: cot(x) = 1/tan(x)
func cmplxCot(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	tanVal := cmplx.Tan(c.Value)
	// Note: If tanVal is zero, the result will be Inf.
	return &tender.Complex{Value: 1 / tanVal}, nil
}

// cmplxExp returns the exponential of the complex number.
func cmplxExp(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Exp(c.Value)}, nil
}

// cmplxInf returns the infinite complex number.
func cmplxInf(args ...tender.Object) (tender.Object, error) {
	if len(args) != 0 {
		return nil, tender.ErrWrongNumArguments
	}
	return &tender.Complex{Value: cmplx.Inf()}, nil
}

// cmplxIsInf returns true if the complex number is infinite.
func cmplxIsInf(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	if cmplx.IsInf(c.Value) {
		return tender.TrueValue, nil
	}
	return tender.FalseValue, nil
}

// cmplxIsNaN returns true if the complex number is NaN.
func cmplxIsNaN(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	if cmplx.IsNaN(c.Value) {
		return tender.TrueValue, nil
	}
	return tender.FalseValue, nil
}

// cmplxLog returns the natural logarithm of the complex number.
func cmplxLog(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Log(c.Value)}, nil
}

// cmplxLog10 returns the base-10 logarithm of the complex number.
func cmplxLog10(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Log10(c.Value)}, nil
}

// cmplxNaN returns a NaN complex number.
func cmplxNaN(args ...tender.Object) (tender.Object, error) {
	if len(args) != 0 {
		return nil, tender.ErrWrongNumArguments
	}
	return &tender.Complex{Value: cmplx.NaN()}, nil
}

// cmplxPolar returns the polar coordinates (r, theta) of the complex number.
func cmplxPolar(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	r, theta := cmplx.Polar(c.Value)
	return &tender.ImmutableMap{Value: map[string]tender.Object{
		"r":     &tender.Float{Value: r},
		"theta": &tender.Float{Value: theta},
	}}, nil
}

// cmplxPow returns x raised to the power y.
func cmplxPow(args ...tender.Object) (tender.Object, error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}
	x, ok1 := args[0].(*tender.Complex)
	y, ok2 := args[1].(*tender.Complex)
	if !ok1 || !ok2 {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "x, y",
			Expected: "complex, complex",
			Found:    fmt.Sprintf("%s, %s", args[0].TypeName(), args[1].TypeName()),
		}
	}
	return &tender.Complex{Value: cmplx.Pow(x.Value, y.Value)}, nil
}

// cmplxRect returns the complex number corresponding to the given polar coordinates.
func cmplxRect(args ...tender.Object) (tender.Object, error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}
	r, ok1 := tender.ToFloat64(args[0])
	theta, ok2 := tender.ToFloat64(args[1])
	if !ok1 || !ok2 {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "r, theta",
			Expected: "float, float",
			Found:    fmt.Sprintf("%s, %s", args[0].TypeName(), args[1].TypeName()),
		}
	}
	return &tender.Complex{Value: cmplx.Rect(r, theta)}, nil
}

// cmplxSinh returns the hyperbolic sine of the complex number.
func cmplxSinh(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Sinh(c.Value)}, nil
}

// cmplxSqrt returns the square root of the complex number.
func cmplxSqrt(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Sqrt(c.Value)}, nil
}

// cmplxTan returns the tangent of the complex number.
func cmplxTan(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Tan(c.Value)}, nil
}

// cmplxTanh returns the hyperbolic tangent of the complex number.
func cmplxTanh(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	c, ok := args[0].(*tender.Complex)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "c",
			Expected: "complex",
			Found:    args[0].TypeName(),
		}
	}
	return &tender.Complex{Value: cmplx.Tanh(c.Value)}, nil
}
