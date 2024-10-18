package stdlib

import (
	"math"

	"github.com/2dprototype/tender"
)

var mathModule = map[string]tender.Object{
	"e":       &tender.Float{Value: math.E},
	"pi":      &tender.Float{Value: math.Pi},
	"phi":     &tender.Float{Value: math.Phi},
	"sqrt2":   &tender.Float{Value: math.Sqrt2},
	"sqrtE":   &tender.Float{Value: math.SqrtE},
	"sqrtPi":  &tender.Float{Value: math.SqrtPi},
	"sqrtPhi": &tender.Float{Value: math.SqrtPhi},
	"ln2":     &tender.Float{Value: math.Ln2},
	"log2E":   &tender.Float{Value: math.Log2E},
	"ln10":    &tender.Float{Value: math.Ln10},
	"log10E":  &tender.Float{Value: math.Log10E},
	"abs": &tender.UserFunction{
		Name:  "abs",
		Value: FuncAFRF(math.Abs),
	},
	"acos": &tender.UserFunction{
		Name:  "acos",
		Value: FuncAFRF(math.Acos),
	},
	"acosh": &tender.UserFunction{
		Name:  "acosh",
		Value: FuncAFRF(math.Acosh),
	},
	"asin": &tender.UserFunction{
		Name:  "asin",
		Value: FuncAFRF(math.Asin),
	},
	"asinh": &tender.UserFunction{
		Name:  "asinh",
		Value: FuncAFRF(math.Asinh),
	},
	"atan": &tender.UserFunction{
		Name:  "atan",
		Value: FuncAFRF(math.Atan),
	},
	"atan2": &tender.UserFunction{
		Name:  "atan2",
		Value: FuncAFFRF(math.Atan2),
	},
	"atanh": &tender.UserFunction{
		Name:  "atanh",
		Value: FuncAFRF(math.Atanh),
	},
	"cbrt": &tender.UserFunction{
		Name:  "cbrt",
		Value: FuncAFRF(math.Cbrt),
	},
	"ceil": &tender.UserFunction{
		Name:  "ceil",
		Value: FuncAFRF(math.Ceil),
	},
	"copysign": &tender.UserFunction{
		Name:  "copysign",
		Value: FuncAFFRF(math.Copysign),
	},
	"cos": &tender.UserFunction{
		Name:  "cos",
		Value: FuncAFRF(math.Cos),
	},
	"cosh": &tender.UserFunction{
		Name:  "cosh",
		Value: FuncAFRF(math.Cosh),
	},
	"dim": &tender.UserFunction{
		Name:  "dim",
		Value: FuncAFFRF(math.Dim),
	},
	"erf": &tender.UserFunction{
		Name:  "erf",
		Value: FuncAFRF(math.Erf),
	},
	"erfc": &tender.UserFunction{
		Name:  "erfc",
		Value: FuncAFRF(math.Erfc),
	},
	"exp": &tender.UserFunction{
		Name:  "exp",
		Value: FuncAFRF(math.Exp),
	},
	"exp2": &tender.UserFunction{
		Name:  "exp2",
		Value: FuncAFRF(math.Exp2),
	},
	"expm1": &tender.UserFunction{
		Name:  "expm1",
		Value: FuncAFRF(math.Expm1),
	},
	"floor": &tender.UserFunction{
		Name:  "floor",
		Value: FuncAFRF(math.Floor),
	},
	"gamma": &tender.UserFunction{
		Name:  "gamma",
		Value: FuncAFRF(math.Gamma),
	},
	"hypot": &tender.UserFunction{
		Name:  "hypot",
		Value: FuncAFFRF(math.Hypot),
	},
	"ilogb": &tender.UserFunction{
		Name:  "ilogb",
		Value: FuncAFRI(math.Ilogb),
	},
	"inf": &tender.UserFunction{
		Name:  "inf",
		Value: FuncAIRF(math.Inf),
	},
	"is_inf": &tender.UserFunction{
		Name:  "is_inf",
		Value: FuncAFIRB(math.IsInf),
	},
	"is_nan": &tender.UserFunction{
		Name:  "is_nan",
		Value: FuncAFRB(math.IsNaN),
	},
	"j0": &tender.UserFunction{
		Name:  "j0",
		Value: FuncAFRF(math.J0),
	},
	"j1": &tender.UserFunction{
		Name:  "j1",
		Value: FuncAFRF(math.J1),
	},
	"jn": &tender.UserFunction{
		Name:  "jn",
		Value: FuncAIFRF(math.Jn),
	},
	"ldexp": &tender.UserFunction{
		Name:  "ldexp",
		Value: FuncAFIRF(math.Ldexp),
	},
	"log": &tender.UserFunction{
		Name:  "log",
		Value: FuncAFRF(math.Log),
	},
	"log10": &tender.UserFunction{
		Name:  "log10",
		Value: FuncAFRF(math.Log10),
	},
	"log1p": &tender.UserFunction{
		Name:  "log1p",
		Value: FuncAFRF(math.Log1p),
	},
	"log2": &tender.UserFunction{
		Name:  "log2",
		Value: FuncAFRF(math.Log2),
	},
	"logb": &tender.UserFunction{
		Name:  "logb",
		Value: FuncAFRF(math.Logb),
	},
	"max": &tender.UserFunction{
		Name:  "max",
		Value: FuncAFFRF(math.Max),
	},
	"min": &tender.UserFunction{
		Name:  "min",
		Value: FuncAFFRF(math.Min),
	},
	"mod": &tender.UserFunction{
		Name:  "mod",
		Value: FuncAFFRF(math.Mod),
	},
	"nan": &tender.UserFunction{
		Name:  "nan",
		Value: FuncARF(math.NaN),
	},
	"nextafter": &tender.UserFunction{
		Name:  "nextafter",
		Value: FuncAFFRF(math.Nextafter),
	},
	"pow": &tender.UserFunction{
		Name:  "pow",
		Value: FuncAFFRF(math.Pow),
	},
	"pow10": &tender.UserFunction{
		Name:  "pow10",
		Value: FuncAIRF(math.Pow10),
	},
	"remainder": &tender.UserFunction{
		Name:  "remainder",
		Value: FuncAFFRF(math.Remainder),
	},
	"signbit": &tender.UserFunction{
		Name:  "signbit",
		Value: FuncAFRB(math.Signbit),
	},
	"sin": &tender.UserFunction{
		Name:  "sin",
		Value: FuncAFRF(math.Sin),
	},
	"sinh": &tender.UserFunction{
		Name:  "sinh",
		Value: FuncAFRF(math.Sinh),
	},
	"sqrt": &tender.UserFunction{
		Name:  "sqrt",
		Value: FuncAFRF(math.Sqrt),
	},
	"tan": &tender.UserFunction{
		Name:  "tan",
		Value: FuncAFRF(math.Tan),
	},
	"tanh": &tender.UserFunction{
		Name:  "tanh",
		Value: FuncAFRF(math.Tanh),
	},
	"trunc": &tender.UserFunction{
		Name:  "trunc",
		Value: FuncAFRF(math.Trunc),
	},
	"y0": &tender.UserFunction{
		Name:  "y0",
		Value: FuncAFRF(math.Y0),
	},
	"y1": &tender.UserFunction{
		Name:  "y1",
		Value: FuncAFRF(math.Y1),
	},
	"yn": &tender.UserFunction{
		Name:  "yn",
		Value: FuncAIFRF(math.Yn),
	},
}
