package stdlib

import (
	"fmt"
	"math/big"

	"github.com/2dprototype/tender"
)

var bigModule = map[string]tender.Object{
	// Constructors
	"int":   &tender.NativeFunction{Name: "int", Value: bigNewInt},
	"float": &tender.NativeFunction{Name: "float", Value: bigNewFloat},

	// BigInt Math
	"iadd": &tender.NativeFunction{Name: "iadd", Value: bigIntAdd},
	"isub": &tender.NativeFunction{Name: "isub", Value: bigIntSub},
	"imul": &tender.NativeFunction{Name: "imul", Value: bigIntMul},
	"idiv": &tender.NativeFunction{Name: "idiv", Value: bigIntDiv},
	"imod": &tender.NativeFunction{Name: "imod", Value: bigIntMod},
	"iexp": &tender.NativeFunction{Name: "iexp", Value: bigIntExp},
	"igcd": &tender.NativeFunction{Name: "igcd", Value: bigIntGCD},
	"icmp": &tender.NativeFunction{Name: "icmp", Value: bigIntCmp},
	"iabs": &tender.NativeFunction{Name: "iabs", Value: bigIntAbs},
	"ineg": &tender.NativeFunction{Name: "ineg", Value: bigIntNeg},

	// BigInt Bitwise Operations
	"iand": &tender.NativeFunction{Name: "iand", Value: bigIntAnd},
	"ior":  &tender.NativeFunction{Name: "ior", Value: bigIntOr},
	"ixor": &tender.NativeFunction{Name: "ixor", Value: bigIntXor},
	"ilsh": &tender.NativeFunction{Name: "ilsh", Value: bigIntLsh},
	"irsh": &tender.NativeFunction{Name: "irsh", Value: bigIntRsh},

	// BigFloat Math
	"fadd":  &tender.NativeFunction{Name: "fadd", Value: bigFloatAdd},
	"fsub":  &tender.NativeFunction{Name: "fsub", Value: bigFloatSub},
	"fmul":  &tender.NativeFunction{Name: "fmul", Value: bigFloatMul},
	"fdiv":  &tender.NativeFunction{Name: "fdiv", Value: bigFloatDiv},
	"fcmp":  &tender.NativeFunction{Name: "fcmp", Value: bigFloatCmp},
	"fsqrt": &tender.NativeFunction{Name: "fsqrt", Value: bigFloatSqrt},
	"fabs":  &tender.NativeFunction{Name: "fabs", Value: bigFloatAbs},
	"fneg":  &tender.NativeFunction{Name: "fneg", Value: bigFloatNeg},
}

// --- Helper Functions ---

func checkArgCount(args []tender.Object, expected int) error {
	if len(args) != expected {
		return tender.ErrWrongNumArguments
	}
	return nil
}

func typeError(name, expected string, found tender.Object) error {
	return tender.ErrInvalidArgumentType{
		Name:     name,
		Expected: expected,
		Found:    found.TypeName(),
	}
}

// --- Constructors ---

func bigNewInt(args ...tender.Object) (tender.Object, error) {
	if len(args) == 0 {
		return &tender.BigInt{Value: new(big.Int)}, nil
	}
	val, ok := tender.ToBigInt(args[0])
	if !ok {
		return nil, typeError("value", "convertible to BigInt", args[0])
	}
	return &tender.BigInt{Value: val}, nil
}

func bigNewFloat(args ...tender.Object) (tender.Object, error) {
	if len(args) == 0 {
		return &tender.BigFloat{Value: new(big.Float)}, nil
	}
	val, ok := tender.ToBigFloat(args[0])
	if !ok {
		return nil, typeError("value", "convertible to BigFloat", args[0])
	}
	return &tender.BigFloat{Value: val}, nil
}

// --- BigInt Operations ---

func bigIntAdd(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	b, ok2 := tender.ToBigInt(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigInt", args[0]) }
	return &tender.BigInt{Value: new(big.Int).Add(a, b)}, nil
}

func bigIntSub(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	b, ok2 := tender.ToBigInt(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigInt", args[0]) }
	return &tender.BigInt{Value: new(big.Int).Sub(a, b)}, nil
}

func bigIntMul(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	b, ok2 := tender.ToBigInt(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigInt", args[0]) }
	return &tender.BigInt{Value: new(big.Int).Mul(a, b)}, nil
}

func bigIntDiv(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	b, ok2 := tender.ToBigInt(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigInt", args[0]) }
	if b.Sign() == 0 { return nil, fmt.Errorf("division by zero") }
	return &tender.BigInt{Value: new(big.Int).Quo(a, b)}, nil
}

func bigIntMod(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	b, ok2 := tender.ToBigInt(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigInt", args[0]) }
	if b.Sign() == 0 { return nil, fmt.Errorf("modulo by zero") }
	return &tender.BigInt{Value: new(big.Int).Rem(a, b)}, nil
}

// iexp takes 2 or 3 arguments: base, exponent, and an optional modulo
func bigIntExp(args ...tender.Object) (tender.Object, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, tender.ErrWrongNumArguments
	}
	base, ok1 := tender.ToBigInt(args[0])
	exp, ok2 := tender.ToBigInt(args[1])
	var mod *big.Int
	if len(args) == 3 {
		m, ok3 := tender.ToBigInt(args[2])
		if !ok3 { return nil, typeError("m", "BigInt", args[2]) }
		mod = m
	}
	if !ok1 || !ok2 { return nil, typeError("base, exp", "BigInt", args[0]) }
	return &tender.BigInt{Value: new(big.Int).Exp(base, exp, mod)}, nil
}

func bigIntGCD(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	b, ok2 := tender.ToBigInt(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigInt", args[0]) }
	res := new(big.Int)
	res.GCD(nil, nil, a, b)
	return &tender.BigInt{Value: res}, nil
}

func bigIntCmp(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	b, ok2 := tender.ToBigInt(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigInt", args[0]) }
	return &tender.Int{Value: int64(a.Cmp(b))}, nil
}

func bigIntAbs(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 1); err != nil { return nil, err }
	a, ok := tender.ToBigInt(args[0])
	if !ok { return nil, typeError("a", "BigInt", args[0]) }
	return &tender.BigInt{Value: new(big.Int).Abs(a)}, nil
}

func bigIntNeg(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 1); err != nil { return nil, err }
	a, ok := tender.ToBigInt(args[0])
	if !ok { return nil, typeError("a", "BigInt", args[0]) }
	return &tender.BigInt{Value: new(big.Int).Neg(a)}, nil
}

// --- BigInt Bitwise Operations ---

func bigIntAnd(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	b, ok2 := tender.ToBigInt(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigInt", args[0]) }
	return &tender.BigInt{Value: new(big.Int).And(a, b)}, nil
}

func bigIntOr(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	b, ok2 := tender.ToBigInt(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigInt", args[0]) }
	return &tender.BigInt{Value: new(big.Int).Or(a, b)}, nil
}

func bigIntXor(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	b, ok2 := tender.ToBigInt(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigInt", args[0]) }
	return &tender.BigInt{Value: new(big.Int).Xor(a, b)}, nil
}

func bigIntLsh(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	n, ok2 := tender.ToUint(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, n", "BigInt and uint", args[0]) }
	return &tender.BigInt{Value: new(big.Int).Lsh(a, n)}, nil
}

func bigIntRsh(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigInt(args[0])
	n, ok2 := tender.ToUint(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, n", "BigInt and uint", args[0]) }
	return &tender.BigInt{Value: new(big.Int).Rsh(a, n)}, nil
}

// --- BigFloat Operations ---

func bigFloatAdd(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigFloat(args[0])
	b, ok2 := tender.ToBigFloat(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigFloat", args[0]) }
	return &tender.BigFloat{Value: new(big.Float).Add(a, b)}, nil
}

func bigFloatSub(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigFloat(args[0])
	b, ok2 := tender.ToBigFloat(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigFloat", args[0]) }
	return &tender.BigFloat{Value: new(big.Float).Sub(a, b)}, nil
}

func bigFloatMul(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigFloat(args[0])
	b, ok2 := tender.ToBigFloat(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigFloat", args[0]) }
	return &tender.BigFloat{Value: new(big.Float).Mul(a, b)}, nil
}

func bigFloatDiv(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigFloat(args[0])
	b, ok2 := tender.ToBigFloat(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigFloat", args[0]) }
	if b.Sign() == 0 { return nil, fmt.Errorf("division by zero") }
	return &tender.BigFloat{Value: new(big.Float).Quo(a, b)}, nil
}

func bigFloatCmp(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 2); err != nil { return nil, err }
	a, ok1 := tender.ToBigFloat(args[0])
	b, ok2 := tender.ToBigFloat(args[1])
	if !ok1 || !ok2 { return nil, typeError("a, b", "BigFloat", args[0]) }
	return &tender.Int{Value: int64(a.Cmp(b))}, nil
}

func bigFloatSqrt(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 1); err != nil { return nil, err }
	a, ok := tender.ToBigFloat(args[0])
	if !ok { return nil, typeError("a", "BigFloat", args[0]) }
	return &tender.BigFloat{Value: new(big.Float).Sqrt(a)}, nil
}

func bigFloatAbs(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 1); err != nil { return nil, err }
	a, ok := tender.ToBigFloat(args[0])
	if !ok { return nil, typeError("a", "BigFloat", args[0]) }
	return &tender.BigFloat{Value: new(big.Float).Abs(a)}, nil
}

func bigFloatNeg(args ...tender.Object) (tender.Object, error) {
	if err := checkArgCount(args, 1); err != nil { return nil, err }
	a, ok := tender.ToBigFloat(args[0])
	if !ok { return nil, typeError("a", "BigFloat", args[0]) }
	return &tender.BigFloat{Value: new(big.Float).Neg(a)}, nil
}