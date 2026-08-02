package tender

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"math/big"
	"math/cmplx"

	"github.com/2dprototype/tender/parser"
	"github.com/2dprototype/tender/token"
)

var (
	// TrueValue represents a true value.
	TrueValue Object = &Bool{value: true}

	// FalseValue represents a false value.
	FalseValue Object = &Bool{value: false}

	// NullValue represents an null value.
	NullValue Object = &Null{}
)

// Object represents an object in the VM.
type Object interface {
	// TypeName should return the name of the type.
	TypeName() string

	// String should return a string representation of the type's value.
	String() string

	// BinaryOp should return another object that is the result of a given
	// binary operator and a right-hand side object. If BinaryOp returns an
	// error, the VM will treat it as a run-time error.
	BinaryOp(op token.Token, rhs Object) (Object, error)

	// IsFalsy should return true if the value of the type should be considered
	// as falsy.
	IsFalsy() bool

	// Equals should return true if the value of the type should be considered
	// as equal to the value of another object.
	Equals(another Object) bool

	// Copy should return a copy of the type (and its value). Copy function
	// will be used for copy() builtin function which is expected to deep-copy
	// the values generally.
	Copy() Object

	// IndexGet should take an index Object and return a result Object or an
	// error for indexable objects. Indexable is an object that can take an
	// index and return an object. If error is returned, the runtime will treat
	// it as a run-time error and ignore returned value. If Object is not
	// indexable, ErrNotIndexable should be returned as error. If nil is
	// returned as value, it will be converted to NullToken value by the
	// runtime.
	IndexGet(index Object) (value Object, err error)

	// IndexSet should take an index Object and a value Object for index
	// assignable objects. Index assignable is an object that can take an index
	// and a value on the left-hand side of the assignment statement. If Object
	// is not index assignable, ErrNotIndexAssignable should be returned as
	// error. If an error is returned, it will be treated as a run-time error.
	IndexSet(index, value Object) error

	// Iterate should return an Iterator for the type.
	Iterate() Iterator

	// CanIterate should return whether the Object can be Iterated.
	CanIterate() bool

	// Call should take an arbitrary number of arguments and returns a return
	// value and/or an error, which the VM will consider as a run-time error.
	Call(args ...Object) (ret Object, err error)

	// CanCall should return whether the Object can be Called.
	CanCall() bool
}

// ObjectImpl represents a default Object Implementation. To defined a new
// value type, one can embed ObjectImpl in their type declarations to avoid
// implementing all non-significant methods. TypeName() and String() methods
// still need to be implemented.
type ObjectImpl struct {
}

// TypeName returns the name of the type.
func (o *ObjectImpl) TypeName() string {
	panic(ErrNotImplemented)
}

func (o *ObjectImpl) String() string {
	panic(ErrNotImplemented)
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *ObjectImpl) BinaryOp(_ token.Token, _ Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *ObjectImpl) Copy() Object {
	return nil
}

// IsFalsy returns true if the value of the type is falsy.
func (o *ObjectImpl) IsFalsy() bool {
	return false
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *ObjectImpl) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *ObjectImpl) IndexGet(_ Object) (res Object, err error) {
	return nil, ErrNotIndexable
}

// IndexSet sets an element at a given index.
func (o *ObjectImpl) IndexSet(_, _ Object) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an iterator.
func (o *ObjectImpl) Iterate() Iterator {
	return nil
}

// CanIterate returns whether the Object can be Iterated.
func (o *ObjectImpl) CanIterate() bool {
	return false
}

// Call takes an arbitrary number of arguments and returns a return value
// and/or an error.
func (o *ObjectImpl) Call(_ ...Object) (ret Object, err error) {
	return nil, nil
}

// CanCall returns whether the Object can be Called.
func (o *ObjectImpl) CanCall() bool {
	return false
}

// Array represents an array of objects.
type Array struct {
	ObjectImpl
	Value []Object
}

// TypeName returns the name of the type.
func (o *Array) TypeName() string {
	return "array"
}

func (o *Array) String() string {
	var elements []string
	for _, e := range o.Value {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Array) BinaryOp(op token.Token, rhs Object) (Object, error) {
	if rhs, ok := rhs.(*Array); ok {
		switch op {
		case token.Add:
			if len(rhs.Value) == 0 {
				return o, nil
			}
			return &Array{Value: append(o.Value, rhs.Value...)}, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Array) Copy() Object {
	var c []Object
	for _, elem := range o.Value {
		c = append(c, elem.Copy())
	}
	return &Array{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Array) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Array) Equals(x Object) bool {
	var xVal []Object
	switch x := x.(type) {
	case *Array:
		xVal = x.Value
	case *ImmutableArray:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for i, e := range o.Value {
		if !e.Equals(xVal[i]) {
			return false
		}
	}
	return true
}

// IndexGet returns an element at a given index.
func (o *Array) IndexGet(index Object) (res Object, err error) {
	if strIdx, ok := index.(*String); ok {
		if strIdx.Value == "length" {
			return &Int{Value: int64(len(o.Value))}, nil
		}
		if method, exists := arrayMethods[strIdx.Value]; exists {
			return &NativeFunction{
				Name: strIdx.Value,
				Value: func(args ...Object) (Object, error) {
					return method(append([]Object{o}, args...)...)
				},
			}, nil
		}
		return nil, nil
	}

	intIdx, ok := index.(*Int)
	if !ok { return nil, ErrInvalidIndexType }
	idxVal := int(intIdx.Value)
	if idxVal < 0 || idxVal >= len(o.Value) { return NullValue, nil }
	return o.Value[idxVal], nil
}

// IndexSet sets an element at a given index.
func (o *Array) IndexSet(index, value Object) (err error) {
	intIdx, ok := ToInt(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	if intIdx < 0 || intIdx >= len(o.Value) {
		err = ErrIndexOutOfBounds
		return
	}
	o.Value[intIdx] = value
	return nil
}

// Iterate creates an array iterator.
func (o *Array) Iterate() Iterator {
	return &ArrayIterator{
		v: o.Value,
		l: len(o.Value),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *Array) CanIterate() bool {
	return true
}

// Bool represents a boolean value.
type Bool struct {
	ObjectImpl

	// this is intentionally non-public to force using objects.TrueValue and
	// FalseValue always
	value bool
}

func (o *Bool) String() string {
	if o.value {
		return "true"
	}

	return "false"
}

// TypeName returns the name of the type.
func (o *Bool) TypeName() string {
	return "bool"
}

// Copy returns a copy of the type.
func (o *Bool) Copy() Object {
	return o
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Bool) IsFalsy() bool {
	return !o.value
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Bool) Equals(x Object) bool {
	return o == x
}

// GobDecode decodes bool value from input bytes.
func (o *Bool) GobDecode(b []byte) (err error) {
	o.value = b[0] == 1
	return
}

// GobEncode encodes bool values into bytes.
func (o *Bool) GobEncode() (b []byte, err error) {
	if o.value {
		b = []byte{1}
	} else {
		b = []byte{0}
	}
	return
}

// NativeFunction represents a native function.
type NativeFunction struct {
	ObjectImpl
	Name      string
	Value     CallableFunc
	NeedVMObj bool
}

// TypeName returns the name of the type.
func (o *NativeFunction) TypeName() string {
	return "native-function"
}

func (o *NativeFunction) String() string {
	if o.Name == "" {
		return "<native-function>"
	}
	return "<native-function:" + o.Name + ">"
}

// Copy returns a copy of the type.
func (o *NativeFunction) Copy() Object {
	return &NativeFunction{Value: o.Value, NeedVMObj: o.NeedVMObj}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *NativeFunction) Equals(_ Object) bool {
	return false
}

// Call executes a native function.
func (o *NativeFunction) Call(args ...Object) (Object, error) {
	return o.Value(args...)
}

// CanCall returns whether the Object can be Called.
func (o *NativeFunction) CanCall() bool {
	return true
}

// BuiltinModule is an importable module that's written in Go.
type BuiltinModule struct {
	Attrs map[string]Object
}

// Import returns an immutable map for the module.
func (m *BuiltinModule) Import(moduleName string) (interface{}, error) {
	return m.AsImmutableMap(moduleName), nil
}

// AsImmutableMap converts builtin module into an immutable map.
func (m *BuiltinModule) AsImmutableMap(moduleName string) *ImmutableMap {
	attrs := make(map[string]Object, len(m.Attrs))
	for k, v := range m.Attrs {
		attrs[k] = v.Copy()
	}
	attrs["__module_name__"] = &String{Value: moduleName}
	return &ImmutableMap{Value: attrs}
}

// Bytes represents a byte array.
type Bytes struct {
	ObjectImpl
	Value []byte
}

func (o *Bytes) String() string {
	return fmt.Sprintf("%v", o.Value)
}

// TypeName returns the name of the type.
func (o *Bytes) TypeName() string {
	return "bytes"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Bytes) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch op {
	case token.Add:
		switch rhs := rhs.(type) {
		case *Bytes:
			if len(o.Value)+len(rhs.Value) > MaxBytesLen {
				return nil, ErrBytesLimit
			}
			return &Bytes{Value: append(o.Value, rhs.Value...)}, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Bytes) Copy() Object {
	return &Bytes{Value: append([]byte{}, o.Value...)}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Bytes) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Bytes) Equals(x Object) bool {
	t, ok := x.(*Bytes)
	if !ok {
		return false
	}
	return bytes.Equal(o.Value, t.Value)
}


func (o *Bytes) IndexSet(index, value Object) (err error) {
    intIdx, ok := ToInt(index)
    if !ok {
        return ErrInvalidIndexType
    }
    if intIdx < 0 || intIdx >= len(o.Value) {
        return ErrIndexOutOfBounds
    }
    
    // Ensure the value is an integer
    byteValue, ok := ToByte(value)
    if !ok {
        return ErrInvalidValueType
    }
    
    // Ensure the integer value is within byte range
    if byteValue < 0 || byteValue > 255 {
        return ErrByteValueOutOfRange
    }
    
    // Set the value at the specified index
    o.Value[intIdx] = byteValue
    return nil
}


// IndexGet returns an element (as Int) at a given index.
func (o *Bytes) IndexGet(index Object) (res Object, err error) {
	strIdx, ok := index.(*String) 
	if ok {
		if strIdx.Value == "length" {
			return &Int{Value: int64(len(o.Value))}, nil
		}
		return nil, nil
	}
	intIdx, ok := index.(*Int)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.Value)
	if idxVal < 0 || idxVal >= len(o.Value) {
		res = NullValue
		return
	}
	res = &Int{Value: int64(o.Value[idxVal])}
	return
}

// Iterate creates a bytes iterator.
func (o *Bytes) Iterate() Iterator {
	return &BytesIterator{
		v: o.Value,
		l: len(o.Value),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *Bytes) CanIterate() bool {
	return true
}

// Char represents a character value.
type Char struct {
	ObjectImpl
	Value rune
}

func (o *Char) String() string {
	return string(o.Value)
}

// TypeName returns the name of the type.
func (o *Char) TypeName() string {
	return "char"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Char) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Char:
		switch op {
		case token.Add:
			r := o.Value + rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case token.Sub:
			r := o.Value - rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case token.Less:
			if o.Value < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Spaceship:
			if o.Value < rhs.Value {
				return &Int{Value: -1}, nil
			} else if o.Value > rhs.Value {
				return &Int{Value: 1}, nil
			}
			return &Int{Value: 0}, nil
		}
	case *Int:
		switch op {
		case token.Add:
			r := o.Value + rune(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case token.Sub:
			r := o.Value - rune(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case token.Less:
			if int64(o.Value) < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if int64(o.Value) > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if int64(o.Value) <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if int64(o.Value) >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Spaceship:
			if int64(o.Value) < rhs.Value {
				return &Int{Value: -1}, nil
			} else if int64(o.Value) > rhs.Value {
				return &Int{Value: 1}, nil
			}
			return &Int{Value: 0}, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Char) Copy() Object {
	return &Char{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Char) IsFalsy() bool {
	return o.Value == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Char) Equals(x Object) bool {
	t, ok := x.(*Char)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

// Function represents a compiled function.
type Function struct {
	ObjectImpl
	Instructions  []byte
	NumLocals     int // number of local variables (including function parameters)
	NumParameters int
	VarArgs       bool
	SourceMap     map[int]parser.Pos
	Free          []*ObjectPtr
}

// TypeName returns the name of the type.
func (o *Function) TypeName() string {
	return "function"
}

func (o *Function) String() string {
	return "<function>"
}

// Copy returns a copy of the type.
func (o *Function) Copy() Object {
	return &Function{
		Instructions:  append([]byte{}, o.Instructions...),
		NumLocals:     o.NumLocals,
		NumParameters: o.NumParameters,
		VarArgs:       o.VarArgs,
		SourceMap:     o.SourceMap,
		Free:          append([]*ObjectPtr{}, o.Free...), // DO NOT Copy() of elements; these are variable pointers
	}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Function) Equals(_ Object) bool {
	return false
}

// SourcePos returns the source position of the instruction at ip.
func (o *Function) SourcePos(ip int) parser.Pos {
	for ip >= 0 {
		if p, ok := o.SourceMap[ip]; ok {
			return p
		}
		ip--
	}
	return parser.NoPos
}

// CanCall returns whether the Object can be Called.
func (o *Function) CanCall() bool {
	return true
}

// Error represents an error value.
type Error struct {
	ObjectImpl
	Value Object
}

// TypeName returns the name of the type.
func (o *Error) TypeName() string {
	return "error"
}

func (o *Error) String() string {
	if o.Value != nil {
		return fmt.Sprintf("error: %s", o.Value.String())
	}
	return "error"
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Error) IsFalsy() bool {
	return true // error is always false.
}

// Copy returns a copy of the type.
func (o *Error) Copy() Object {
	return &Error{Value: o.Value.Copy()}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Error) Equals(x Object) bool {
	return o == x // pointer equality
}

// IndexGet returns an element at a given index.
func (o *Error) IndexGet(index Object) (res Object, err error) {
	if strIdx, _ := ToString(index); strIdx != "value" {
		err = ErrInvalidIndexOnError
		return
	}
	res = o.Value
	return
}

// ImmutableArray represents an immutable array of objects.
type ImmutableArray struct {
	ObjectImpl
	Value []Object
}

// TypeName returns the name of the type.
func (o *ImmutableArray) TypeName() string {
	return "immutable-array"
}

func (o *ImmutableArray) String() string {
	var elements []string
	for _, e := range o.Value {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *ImmutableArray) BinaryOp(op token.Token, rhs Object) (Object, error) {
	if rhs, ok := rhs.(*ImmutableArray); ok {
		switch op {
		case token.Add:
			return &Array{Value: append(o.Value, rhs.Value...)}, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *ImmutableArray) Copy() Object {
	var c []Object
	for _, elem := range o.Value {
		c = append(c, elem.Copy())
	}
	return &Array{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *ImmutableArray) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *ImmutableArray) Equals(x Object) bool {
	var xVal []Object
	switch x := x.(type) {
	case *Array:
		xVal = x.Value
	case *ImmutableArray:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for i, e := range o.Value {
		if !e.Equals(xVal[i]) {
			return false
		}
	}
	return true
}

// IndexGet returns an element at a given index.
func (o *ImmutableArray) IndexGet(index Object) (res Object, err error) {
	if strIdx, ok := index.(*String); ok {
		if strIdx.Value == "length" {
			return &Int{Value: int64(len(o.Value))}, nil
		}
		if method, exists := immutableArrayMethods[strIdx.Value]; exists {
			return &NativeFunction{
				Name: strIdx.Value,
				Value: func(args ...Object) (Object, error) {
					return method(append([]Object{o}, args...)...)
				},
			}, nil
		}
		return nil, nil
	}

	intIdx, ok := index.(*Int)
	if !ok { return nil, ErrInvalidIndexType }
	idxVal := int(intIdx.Value)
	if idxVal < 0 || idxVal >= len(o.Value) { return NullValue, nil }
	return o.Value[idxVal], nil
}

// Iterate creates an array iterator.
func (o *ImmutableArray) Iterate() Iterator {
	return &ArrayIterator{
		v: o.Value,
		l: len(o.Value),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *ImmutableArray) CanIterate() bool {
	return true
}

// Tuple represents an immutable array of objects.
type Tuple struct {
	ObjectImpl
	Value []Object
}

// TypeName returns the name of the type.
func (o *Tuple) TypeName() string {
	return "tuple"
}

func (o *Tuple) String() string {
	var elements []string
	for _, e := range o.Value {
		elements = append(elements, e.String())
	}
	if len(elements) == 1 {
		return "(" + elements[0] + ",)"
	}
	return "(" + strings.Join(elements, ", ") + ")"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Tuple) BinaryOp(op token.Token, rhs Object) (Object, error) {
	if rhs, ok := rhs.(*Tuple); ok {
		switch op {
		case token.Add:
			return &Tuple{Value: append(o.Value, rhs.Value...)}, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Tuple) Copy() Object {
	var c []Object
	for _, elem := range o.Value {
		c = append(c, elem.Copy())
	}
	return &Tuple{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Tuple) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Tuple) Equals(x Object) bool {
	var xVal []Object
	switch x := x.(type) {
	case *Array:
		xVal = x.Value
	case *ImmutableArray:
		xVal = x.Value
	case *Tuple:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for i, e := range o.Value {
		if !e.Equals(xVal[i]) {
			return false
		}
	}
	return true
}

// IndexGet returns an element at a given index.
func (o *Tuple) IndexGet(index Object) (res Object, err error) {
	strIdx, ok := index.(*String)
	if ok {
		if strIdx.Value == "length" {
			return &Int{Value: int64(len(o.Value))}, nil
		}
		return nil, nil
	}
	intIdx, ok := index.(*Int)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.Value)
	if idxVal < 0 || idxVal >= len(o.Value) {
		res = NullValue
		return
	}
	res = o.Value[idxVal]
	return
}

// Iterate creates an array iterator.
func (o *Tuple) Iterate() Iterator {
	return &ArrayIterator{
		v: o.Value,
		l: len(o.Value),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *Tuple) CanIterate() bool {
	return true
}

// ImmutableMap represents an immutable map object.
type ImmutableMap struct {
	ObjectImpl
	Value map[string]Object
}

// TypeName returns the name of the type.
func (o *ImmutableMap) TypeName() string {
	return "immutable-map"
}

func (o *ImmutableMap) String() string {
	var pairs []string
	for k, v := range o.Value {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Copy returns a copy of the type.
func (o *ImmutableMap) Copy() Object {
	c := make(map[string]Object)
	for k, v := range o.Value {
		c[k] = v.Copy()
	}
	return &Map{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *ImmutableMap) IsFalsy() bool {
	return len(o.Value) == 0
}

// IndexGet returns the value for the given key.
func (o *ImmutableMap) IndexGet(index Object) (res Object, err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	res, ok = o.Value[strIdx]
	if !ok {
		res = NullValue
	}
	return
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *ImmutableMap) Equals(x Object) bool {
	var xVal map[string]Object
	switch x := x.(type) {
	case *Map:
		xVal = x.Value
	case *ImmutableMap:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for k, v := range o.Value {
		tv := xVal[k]
		if !v.Equals(tv) {
			return false
		}
	}
	return true
}

// Iterate creates an immutable map iterator.
func (o *ImmutableMap) Iterate() Iterator {
	var keys []string
	for k := range o.Value {
		keys = append(keys, k)
	}
	return &MapIterator{
		v: o.Value,
		k: keys,
		l: len(keys),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *ImmutableMap) CanIterate() bool {
	return true
}


// Float represents a floating point number value.
type Float struct {
	ObjectImpl
	Value float64
}

func (o *Float) String() string {
	return strconv.FormatFloat(o.Value, 'f', -1, 64)
}

// TypeName returns the name of the type.
func (o *Float) TypeName() string {
	return "float"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Float) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Float:
		switch op {
		case token.Add:
			r := o.Value + rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Sub:
			r := o.Value - rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Mul:
			r := o.Value * rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Quo:
			if rhs.Value == 0 {
				return &Float{Value: math.Inf(1)}, nil
			}
			r := o.Value / rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Less:
			if o.Value < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Spaceship:
			if o.Value < rhs.Value {
				return &Int{Value: -1}, nil
			} else if o.Value > rhs.Value {
				return &Int{Value: 1}, nil
			}
			return &Int{Value: 0}, nil
		}
	case *Int:
		switch op {
		case token.Add:
			r := o.Value + float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Sub:
			r := o.Value - float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Mul:
			r := o.Value * float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Quo:
			if rhs.Value == 0 {
				return &Float{Value: math.Inf(1)}, nil
			}
			r := o.Value / float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Less:
			if o.Value < float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value > float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value <= float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value >= float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Spaceship:
			if o.Value < float64(rhs.Value) {
				return &Int{Value: -1}, nil
			} else if o.Value > float64(rhs.Value) {
				return &Int{Value: 1}, nil
			}
			return &Int{Value: 0}, nil
		}
		case *BigInt:
			bi := new(big.Int)
			new(big.Float).SetFloat64(o.Value).Int(bi)
			return binaryOpBigInt(op, bi, rhs.Value), nil
		case *BigFloat:
			return binaryOpBigFloat(op, new(big.Float).SetFloat64(o.Value), rhs.Value), nil
		case *Complex:
			return binaryOpComplex(op, complex(o.Value, 0), rhs.Value), nil
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Float) Copy() Object {
	return &Float{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Float) IsFalsy() bool {
	return math.IsNaN(o.Value)
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Float) Equals(x Object) bool {
	t, ok := x.(*Float)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

// Int represents an integer value.
type Int struct {
	ObjectImpl
	Value int64
}

func (o *Int) String() string {
	return strconv.FormatInt(o.Value, 10)
}

// TypeName returns the name of the type.
func (o *Int) TypeName() string {
	return "int"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Int) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Int:
		switch op {
		case token.Add:
			r := o.Value + rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Sub:
			r := o.Value - rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Mul:
			r := o.Value * rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Quo:
			if rhs.Value == 0 {
				return &Float{Value: math.Inf(1)}, nil
			}
			r := o.Value / rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Rem:
			r := o.Value % rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.And:
			r := o.Value & rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Or:
			r := o.Value | rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Xor:
			r := o.Value ^ rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.AndNot:
			r := o.Value &^ rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Shl:
			r := o.Value << uint64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Shr:
			r := o.Value >> uint64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Less:
			if o.Value < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Spaceship:
			if o.Value < rhs.Value {
				return &Int{Value: -1}, nil
			} else if o.Value > rhs.Value {
				return &Int{Value: 1}, nil
			}
			return &Int{Value: 0}, nil
		}
	case *Float:
		switch op {
		case token.Add:
			return &Float{Value: float64(o.Value) + rhs.Value}, nil
		case token.Sub:
			return &Float{Value: float64(o.Value) - rhs.Value}, nil
		case token.Mul:
			return &Float{Value: float64(o.Value) * rhs.Value}, nil
		case token.Quo:
			if rhs.Value == 0 {
				return &Float{Value: math.Inf(1)}, nil
			}
			return &Float{Value: float64(o.Value) / rhs.Value}, nil
		case token.Less:
			if float64(o.Value) < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if float64(o.Value) > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if float64(o.Value) <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if float64(o.Value) >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Spaceship:
			if float64(o.Value) < rhs.Value {
				return &Int{Value: -1}, nil
			} else if float64(o.Value) > rhs.Value {
				return &Int{Value: 1}, nil
			}
			return &Int{Value: 0}, nil
		}
	case *Char:
		switch op {
		case token.Add:
			return &Char{Value: rune(o.Value) + rhs.Value}, nil
		case token.Sub:
			return &Char{Value: rune(o.Value) - rhs.Value}, nil
		case token.Less:
			if o.Value < int64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value > int64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value <= int64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value >= int64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Spaceship:
			if o.Value < int64(rhs.Value) {
				return &Int{Value: -1}, nil
			} else if o.Value > int64(rhs.Value) {
				return &Int{Value: 1}, nil
			}
			return &Int{Value: 0}, nil
		}
		case *BigInt:
			return binaryOpBigInt(op, new(big.Int).SetInt64(o.Value), rhs.Value), nil
		case *BigFloat:
			return binaryOpBigFloat(op, new(big.Float).SetInt64(o.Value), rhs.Value), nil
		case *Complex:
			return binaryOpComplex(op, complex(float64(o.Value), 0), rhs.Value), nil
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Int) Copy() Object {
	return &Int{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Int) IsFalsy() bool {
	return o.Value == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Int) Equals(x Object) bool {
	t, ok := x.(*Int)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

// BigInt represents an arbitrary-precision integer value.
type BigInt struct {
    ObjectImpl
    Value *big.Int
}

func (b *BigInt) String() string {
    return b.Value.String()
}

// TypeName returns the name of the type.
func (b *BigInt) TypeName() string {
    return "bigint"
}


func binaryOpBigInt (op token.Token, lhs *big.Int, rhs *big.Int) Object {
	switch op {
        case token.Add:
            r := new(big.Int).Set(lhs)
            r.Add(r, rhs)
            return &BigInt{Value: r}
        case token.Sub:
            r := new(big.Int).Set(lhs)
            r.Sub(r, rhs)
            return &BigInt{Value: r}
        case token.Mul:
            r := new(big.Int).Set(lhs)
            r.Mul(r, rhs)
            return &BigInt{Value: r}
        case token.Quo:
			if rhs.Int64() == 0 {
				return &Float{Value: math.Inf(1)}
			}
            r := new(big.Int).Set(lhs)
            r.Div(r, rhs)
            return &BigInt{Value: r}
        case token.Rem:
            r := new(big.Int).Set(lhs)
            r.Mod(r, rhs)
            return &BigInt{Value: r}
        case token.And:
            r := new(big.Int).Set(lhs)
            r.And(r, rhs)
            return &BigInt{Value: r}
        case token.Or:
            r := new(big.Int).Set(lhs)
            r.Or(r, rhs)
            return &BigInt{Value: r}
        case token.Xor:
            r := new(big.Int).Set(lhs)
            r.Xor(r, rhs)
            return &BigInt{Value: r}
        case token.AndNot:
            r := new(big.Int).Set(lhs)
            r.AndNot(r, rhs)
            return &BigInt{Value: r}
        case token.Shl:
            r := new(big.Int).Set(lhs)
            r.Lsh(r, uint(rhs.Int64()))
            return &BigInt{Value: r}
        case token.Shr:
            r := new(big.Int).Set(lhs)
            r.Rsh(r, uint(rhs.Int64()))
            return &BigInt{Value: r}
        case token.Less:
            if lhs.Cmp(rhs) < 0 {
                return TrueValue
            }
            return FalseValue
        case token.Greater:
            if lhs.Cmp(rhs) > 0 {
                return TrueValue
            }
            return FalseValue
        case token.LessEq:
            if lhs.Cmp(rhs) <= 0 {
                return TrueValue
            }
            return FalseValue
        case token.GreaterEq:
            if lhs.Cmp(rhs) >= 0 {
                return TrueValue
            }
            return FalseValue
        case token.Spaceship:
			cmp := lhs.Cmp(rhs)
			if cmp < 0 {
				return &Int{Value: -1}
			} else if cmp > 0 {
				return &Int{Value: 1}
			}
			return &Int{Value: 0}
	}
	
	return nil
}

// BinaryOp performs binary operations with another Object.
func (b *BigInt) BinaryOp(op token.Token, rhs Object) (Object, error) {
    switch rhs := rhs.(type) {
		case *BigInt:
			return binaryOpBigInt(op, b.Value, rhs.Value), nil
		case *Float:
			return binaryOpBigInt(op, b.Value, new(big.Int).SetInt64(int64(rhs.Value))), nil
		case *Int:
			return binaryOpBigInt(op, b.Value, new(big.Int).SetInt64(rhs.Value)), nil
		case *BigFloat:
			return binaryOpBigFloat(op, new(big.Float).SetInt(b.Value), rhs.Value), nil
    }
    return nil, ErrInvalidOperator
}

// Copy returns a copy of the BigInt.
func (b *BigInt) Copy() Object {
    return &BigInt{Value: new(big.Int).Set(b.Value)}
}

// IsFalsy returns true if the value of the BigInt is falsy (i.e., zero).
func (b *BigInt) IsFalsy() bool {
    return b.Value.Sign() == 0
}

// Equals checks if the BigInt is equal to another Object.
func (b *BigInt) Equals(x Object) bool {
    t, ok := x.(*BigInt)
    if !ok {
        return false
    }
    return b.Value.Cmp(t.Value) == 0
}


// BigFloat represents an arbitrary-precision floating-point value.
type BigFloat struct {
    ObjectImpl
    Value *big.Float
}

// String returns the string representation of the BigFloat.
func (b *BigFloat) String() string {
    return b.Value.String()
}

// TypeName returns the name of the type.
func (b *BigFloat) TypeName() string {
    return "bigfloat"
}

func binaryOpBigFloat (op token.Token, lhs *big.Float, rhs *big.Float) Object {
	switch op {
        case token.Add:
			r := new(big.Float).Set(lhs)
            r.Add(r, rhs)
            return &BigFloat{Value: r}
        case token.Sub:
			r := new(big.Float).Set(lhs)
            r.Sub(r, rhs)
            return &BigFloat{Value: r}
        case token.Mul:
			r := new(big.Float).Set(lhs)
            r.Mul(r, rhs)
            return &BigFloat{Value: r}
        case token.Quo:
			rv, _ := rhs.Int64()
			if rv == 0 {
				return &Float{Value: math.Inf(1)}
			}
			r := new(big.Float).Set(lhs)
            r.Quo(r, rhs)
            return &BigFloat{Value: r}
        case token.Less:
            if lhs.Cmp(rhs) < 0 {
                return TrueValue
            }
            return FalseValue
        case token.Greater:
            if lhs.Cmp(rhs) > 0 {
                return TrueValue
            }
            return FalseValue
        case token.LessEq:
            if lhs.Cmp(rhs) <= 0 {
                return TrueValue
            }
            return FalseValue
        case token.GreaterEq:
            if lhs.Cmp(rhs) >= 0 {
                return TrueValue
            }
            return FalseValue
        case token.Spaceship:
			cmp := lhs.Cmp(rhs)
			if cmp < 0 {
				return &Int{Value: -1}
			} else if cmp > 0 {
				return &Int{Value: 1}
			}
			return &Int{Value: 0}
	}
	return nil
}

// BinaryOp performs binary operations with another Object.
func (b *BigFloat) BinaryOp(op token.Token, rhs Object) (Object, error) {
    switch rhs := rhs.(type) {
		case *BigFloat:
			return binaryOpBigFloat(op, b.Value, rhs.Value), nil
		case *Float:
			return binaryOpBigFloat(op, b.Value, new(big.Float).SetFloat64(rhs.Value)), nil
		case *Int:
			return binaryOpBigFloat(op, b.Value, new(big.Float).SetInt64(rhs.Value)), nil
		case *BigInt:
			return binaryOpBigFloat(op, b.Value, new(big.Float).SetInt(rhs.Value)), nil
    }
    return nil, ErrInvalidOperator
}

// Copy returns a copy of the BigFloat.
func (b *BigFloat) Copy() Object {
    return &BigFloat{Value: new(big.Float).Set(b.Value)}
}

// IsFalsy returns true if the value of the BigFloat is falsy (i.e., zero).
func (b *BigFloat) IsFalsy() bool {
    return b.Value.Cmp(big.NewFloat(0)) == 0
}

// Equals checks if the BigFloat is equal to another Object.
func (b *BigFloat) Equals(x Object) bool {
    t, ok := x.(*BigFloat)
    if !ok {
        return false
    }
    return b.Value.Cmp(t.Value) == 0
}

// Complex represents a complex number
type Complex struct {
	ObjectImpl
	Value complex128
}

// String returns the string representation of the Complex number.
func (c *Complex) String() string {
	if cmplx.IsNaN(c.Value) {
		return "NaN"
	} else if cmplx.IsInf(c.Value) {
		return "Inf"
	}
	if imag(c.Value) < 0 {
		return fmt.Sprintf("%g%gi", real(c.Value), imag(c.Value))
	}
	return fmt.Sprintf("%g+%gi", real(c.Value), imag(c.Value))
}


// IndexGet returns a character at a given index.
func (c *Complex) IndexGet(index Object) (res Object, err error) {
	strIdx, ok := index.(*String) 
	if ok {
		if strIdx.Value == "real" {
			return &Float{Value: real(c.Value)}, nil
		} else if strIdx.Value == "imag" {
			return &Float{Value: imag(c.Value)}, nil
		}
		return nil, nil
	}
	return nil, nil
}

// TypeName returns the name of the type.
func (c *Complex) TypeName() string {
	return "complex"
}

// BinaryOp performs binary operations with another Object.
func (c *Complex) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Complex:
		return binaryOpComplex(op, c.Value, rhs.Value), nil
	case *Float:
		return binaryOpComplex(op, c.Value, complex(rhs.Value, 0)), nil
	case *Int:
		return binaryOpComplex(op, c.Value, complex(float64(rhs.Value), 0)), nil
	}
	return nil, ErrInvalidOperator
}

// binaryOpComplex handles binary operations for Complex numbers.
func binaryOpComplex(op token.Token, lhs, rhs complex128) Object {
	switch op {
	case token.Add:
		return &Complex{Value: lhs + rhs}
	case token.Sub:
		return &Complex{Value: lhs - rhs}
	case token.Mul:
		return &Complex{Value: lhs * rhs}
	case token.Quo:
		if rhs == 0 {
			// return &Complex{Value: cmplx.Inf()}
			return &Complex{Value: lhs / rhs}
		}
		return &Complex{Value: lhs / rhs}
	}
	return nil
} 

// Copy returns a copy of the Complex number.
func (c *Complex) Copy() Object {
	return &Complex{Value: c.Value}
}

// IsFalsy returns true if the Complex number is zero.
func (c *Complex) IsFalsy() bool {
	return c.Value == 0
}

// Equals checks if the Complex number is equal to another Object.
func (c *Complex) Equals(x Object) bool {
	t, ok := x.(*Complex)
	if !ok {
		return false
	}
	return c.Value == t.Value
}



// Map represents a map of objects.
type Map struct {
	ObjectImpl
	Value map[string]Object
}

// TypeName returns the name of the type.
func (o *Map) TypeName() string {
	return "map"
}

func (o *Map) String() string {
	var pairs []string
	for k, v := range o.Value {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Copy returns a copy of the type.
func (o *Map) Copy() Object {
	c := make(map[string]Object)
	for k, v := range o.Value {
		c[k] = v.Copy()
	}
	return &Map{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Map) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Map) Equals(x Object) bool {
	var xVal map[string]Object
	switch x := x.(type) {
	case *Map:
		xVal = x.Value
	case *ImmutableMap:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for k, v := range o.Value {
		tv := xVal[k]
		if !v.Equals(tv) {
			return false
		}
	}
	return true
}

// IndexGet returns the value for the given key.
func (o *Map) IndexGet(index Object) (res Object, err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	res, ok = o.Value[strIdx]
	if !ok {
		res = NullValue
	}
	return
}

// IndexSet sets the value for the given key.
func (o *Map) IndexSet(index, value Object) (err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	o.Value[strIdx] = value
	return nil
}

// Iterate creates a map iterator.
func (o *Map) Iterate() Iterator {
	var keys []string
	for k := range o.Value {
		keys = append(keys, k)
	}
	return &MapIterator{
		v: o.Value,
		k: keys,
		l: len(keys),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *Map) CanIterate() bool {
	return true
}

// ObjectPtr represents a free variable.
type ObjectPtr struct {
	ObjectImpl
	Value *Object
}

func (o *ObjectPtr) String() string {
	return "free-var"
}

// TypeName returns the name of the type.
func (o *ObjectPtr) TypeName() string {
	return "<free-var>"
}

// Copy returns a copy of the type.
func (o *ObjectPtr) Copy() Object {
	return o
}

// IsFalsy returns true if the value of the type is falsy.
func (o *ObjectPtr) IsFalsy() bool {
	return o.Value == nil
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *ObjectPtr) Equals(x Object) bool {
	return o == x
}

// String represents a string value.
type String struct {
	ObjectImpl
	Value   string
	runeStr []rune
}

// TypeName returns the name of the type.
func (o *String) TypeName() string {
	return "string"
}

func (o *String) String() string {
	return strconv.Quote(o.Value)
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *String) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch op {
	case token.Add:
		switch rhs := rhs.(type) {
		case *String:
			if len(o.Value)+len(rhs.Value) > MaxStringLen {
				return nil, ErrStringLimit
			}
			return &String{Value: o.Value + rhs.Value}, nil
		default:
			rhsStr := rhs.String()
			if len(o.Value)+len(rhsStr) > MaxStringLen {
				return nil, ErrStringLimit
			}
			return &String{Value: o.Value + rhsStr}, nil
		}
	case token.Less:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case token.LessEq:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case token.Greater:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case token.GreaterEq:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case token.Spaceship:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value < rhs.Value {
				return &Int{Value: -1}, nil
			} else if o.Value > rhs.Value {
				return &Int{Value: 1}, nil
			}
			return &Int{Value: 0}, nil
		}
	}
	return nil, ErrInvalidOperator
}

// IsFalsy returns true if the value of the type is falsy.
func (o *String) IsFalsy() bool {
	return len(o.Value) == 0
}

// Copy returns a copy of the type.
func (o *String) Copy() Object {
	return &String{Value: o.Value}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *String) Equals(x Object) bool {
	t, ok := x.(*String)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

var ansiColorMap = map[string]string {
	"reset" : "\033[0m",           // Text Reset
	// Regular Colors
	"black" : "\033[0;30m",        // black
	"red" : "\033[0;31m",          // red
	"green" : "\033[0;32m",        // green
	"yellow" : "\033[0;33m",       // yellow
	"blue" : "\033[0;34m",         // blue
	"purple" : "\033[0;35m",       // purple
	"cyan" : "\033[0;36m",         // cyan
	"white" : "\033[0;37m",        // white
	// Bold
	"bblack" : "\033[1;30m",       // black
	"bred" : "\033[1;31m",         // red
	"bgreen" : "\033[1;32m",       // green
	"byellow" : "\033[1;33m",      // yellow
	"bblue" : "\033[1;34m",        // blue
	"bpurple" : "\033[1;35m",      // purple
	"bcyan" : "\033[1;36m",        // cyan
	"bwhite" : "\033[1;37m",       // white
	// Underline
	"ublack" : "\033[4;30m",       // black
	"ured" : "\033[4;31m",         // red
	"ugreen" : "\033[4;32m",       // green
	"uyellow" : "\033[4;33m",      // yellow
	"ublue" : "\033[4;34m",        // blue
	"upurple" : "\033[4;35m",      // purple
	"ucyan" : "\033[4;36m",        // cyan
	"uwhite" : "\033[4;37m",       // white
	// Background
	"on_black" : "\033[40m",       // black
	"on_red" : "\033[41m",         // red
	"on_green" : "\033[42m",       // green
	"on_yellow" : "\033[43m",      // yellow
	"on_blue" : "\033[44m",        // blue
	"on_purple" : "\033[45m",      // purple
	"on_cyan" : "\033[46m",        // cyan
	"on_white" : "\033[47m",       // white
	// High Intensty
	"iblack" : "\033[0;90m",       // black
	"ired" : "\033[0;91m",         // red
	"igreen" : "\033[0;92m",       // green
	"iyellow" : "\033[0;93m",      // yellow
	"iblue" : "\033[0;94m",        // blue
	"ipurple" : "\033[0;95m",      // purple
	"icyan" : "\033[0;96m",        // cyan
	"iwhite" : "\033[0;97m",       // white
	// Bold High Intensty
	"biblack" : "\033[1;90m",      // black
	"bired" : "\033[1;91m",        // red
	"bigreen" : "\033[1;92m",      // green
	"biyellow" : "\033[1;93m",     // yellow
	"biblue" : "\033[1;94m",       // blue
	"bipurple" : "\033[1;95m",     // purple
	"bicyan" : "\033[1;96m",       // cyan
	"biwhite" : "\033[1;97m",      // white
	// High Intensty backgrounds
	"on_iblack" : "\033[0;100m",   // black
	"on_ired" : "\033[0;101m",     // red
	"on_igreen" : "\033[0;102m",   // green
	"on_iyellow" : "\033[0;103m",  // yellow
	"on_iblue" : "\033[0;104m",    // blue
	"on_ipurple" : "\033[10;95m",  // purple
	"on_icyan" : "\033[0;106m",    // cyan
	"on_iwhite" : "\033[0;107m",   // white
}

// IndexGet returns a character at a given index.
func (o *String) IndexGet(index Object) (res Object, err error) {
	strIdx, ok := index.(*String) 
	if ok {
		if strIdx.Value == "length" {
			return &Int{Value: int64(len(o.Value))}, nil
		} else if v, ok := ansiColorMap[strIdx.Value]; ok {
			return &String{Value: (v + o.Value + "\033[0m")}, nil
		}
		return nil, nil
	}
	intIdx, ok := index.(*Int)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.Value)
	if o.runeStr == nil {
		o.runeStr = []rune(o.Value)
	}
	if idxVal < 0 || idxVal >= len(o.runeStr) {
		res = NullValue
		return
	}
	res = &Char{Value: o.runeStr[idxVal]}
	return
}

func (o *String) IndexSet(index, value Object) error {
	char, ok := value.(*Char)
	if !ok {
		return ErrInvalidIndexValueType
	}

    intIdx, ok := index.(*Int)
    if ok {
        if intIdx.Value >= 0 && intIdx.Value < int64(len(o.Value)) {
			runes := []rune(o.Value)
            runes[intIdx.Value] = char.Value
			o.Value = string(runes)
            return nil
        }

        return ErrIndexOutOfBounds
    }

    return ErrInvalidIndexType
}

// Iterate creates a string iterator.
func (o *String) Iterate() Iterator {
	if o.runeStr == nil {
		o.runeStr = []rune(o.Value)
	}
	return &StringIterator{
		v: o.runeStr,
		l: len(o.runeStr),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *String) CanIterate() bool {
	return true
}

// Time represents a time value.
type Time struct {
	ObjectImpl
	Value time.Time
}

func (o *Time) String() string {
	// return o.Value.String()
	return o.Value.Format("02/01/2006 3:04PM")
}

// TypeName returns the name of the type.
func (o *Time) TypeName() string {
	return "time"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Time) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Int:
		switch op {
		case token.Add: // time + int => time
			if rhs.Value == 0 {
				return o, nil
			}
			return &Time{Value: o.Value.Add(time.Duration(rhs.Value))}, nil
		case token.Sub: // time - int => time
			if rhs.Value == 0 {
				return o, nil
			}
			return &Time{Value: o.Value.Add(time.Duration(-rhs.Value))}, nil
		}
	case *Time:
		switch op {
		case token.Sub: // time - time => int (duration)
			return &Int{Value: int64(o.Value.Sub(rhs.Value))}, nil
		case token.Less: // time < time => bool
			if o.Value.Before(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value.After(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value.Equal(rhs.Value) || o.Value.Before(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value.Equal(rhs.Value) || o.Value.After(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Spaceship:
			if o.Value.Before(rhs.Value) {
				return &Int{Value: -1}, nil
			} else if o.Value.After(rhs.Value) {
				return &Int{Value: 1}, nil
			}
			return &Int{Value: 0}, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Time) Copy() Object {
	return &Time{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Time) IsFalsy() bool {
	return o.Value.IsZero()
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Time) Equals(x Object) bool {
	t, ok := x.(*Time)
	if !ok {
		return false
	}
	return o.Value.Equal(t.Value)
}

// Null represents an null value.
type Null struct {
	ObjectImpl
}

// TypeName returns the name of the type.
func (o *Null) TypeName() string {
	return "null"
}

func (o *Null) String() string {
	return "null"
}

// Copy returns a copy of the type.
func (o *Null) Copy() Object {
	return o
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Null) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Null) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *Null) IndexGet(_ Object) (Object, error) {
	return NullValue, nil
}

// Iterate creates a map iterator.
func (o *Null) Iterate() Iterator {
	return o
}

// CanIterate returns whether the Object can be Iterated.
func (o *Null) CanIterate() bool {
	return true
}

// Next returns true if there are more elements to iterate.
func (o *Null) Next() bool {
	return false
}

// Key returns the key or index value of the current element.
func (o *Null) Key() Object {
	return o
}

// Value returns the value of the current element.
func (o *Null) Value() Object {
	return o
}

var StructTypes = make(map[string]*StructType)

func zeroValue(typeName string) Object {
	switch typeName {
	case "int":
		return &Int{Value: 0}
	case "float":
		return &Float{Value: 0.0}
	case "string":
		return &String{Value: ""}
	case "bool":
		return FalseValue
	case "char":
		return &Char{Value: 0}
	case "bytes":
		return &Bytes{Value: nil}
	case "array":
		return &Array{Value: nil}
	case "map":
		return &Map{Value: make(map[string]Object)}
	default:
		if st, ok := StructTypes[typeName]; ok {
			return st.NewInstance()
		}
		return NullValue
	}
}

func checkType(val Object, expectedType string) bool {
	if expectedType == "" {
		return true
	}
	if val == NullValue {
		return true
	}
	return val.TypeName() == expectedType
}

type StructField struct {
	Name string
	Type string
	Tag  string
}

type StructType struct {
	ObjectImpl
	Name    string
	Fields  []StructField
	Methods map[string]Object
}

func (o *StructType) TypeName() string {
	return "struct-type"
}

func (o *StructType) String() string {
	if o.Name != "" {
		return fmt.Sprintf("struct %s", o.Name)
	}
	return "struct"
}

func (o *StructType) Copy() Object {
	return o
}

func (o *StructType) Equals(another Object) bool {
	t, ok := another.(*StructType)
	if !ok {
		return false
	}
	if o.Name != "" || t.Name != "" {
		return o == t
	}
	if len(o.Fields) != len(t.Fields) {
		return false
	}
	for i := range o.Fields {
		if o.Fields[i].Name != t.Fields[i].Name || o.Fields[i].Type != t.Fields[i].Type {
			return false
		}
	}
	return true
}

func (o *StructType) CanCall() bool {
	return true
}

func (o *StructType) Call(args ...Object) (Object, error) {
	if len(args) > len(o.Fields) {
		return nil, fmt.Errorf("too many arguments for struct %s initialization (expected <= %d, got %d)", o.Name, len(o.Fields), len(args))
	}
	inst := o.NewInstance()
	for i, val := range args {
		fieldName := o.Fields[i].Name
		if err := inst.IndexSet(&String{Value: fieldName}, val); err != nil {
			return nil, err
		}
	}
	return inst, nil
}

func (o *StructType) NewInstance() *Struct {
	fields := make(map[string]Object)
	for _, f := range o.Fields {
		fields[f.Name] = zeroValue(f.Type)
	}
	if o.Methods == nil {
		o.Methods = make(map[string]Object)
	}
	if o.Name != "" {
		StructTypes[o.Name] = o
	}
	return &Struct{
		Type:   o,
		Fields: fields,
	}
}

func (o *StructType) GobEncode() (b []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(o.Name); err != nil {
		return nil, err
	}
	if err := enc.Encode(o.Fields); err != nil {
		return nil, err
	}
	if err := enc.Encode(o.Methods); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (o *StructType) GobDecode(b []byte) (err error) {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&o.Name); err != nil {
		return err
	}
	if err := dec.Decode(&o.Fields); err != nil {
		return err
	}
	if err := dec.Decode(&o.Methods); err != nil {
		return err
	}
	if o.Name != "" {
		StructTypes[o.Name] = o
	}
	return nil
}

type Struct struct {
	ObjectImpl
	Type   *StructType
	Fields map[string]Object
}

func (o *Struct) TypeName() string {
	if o.Type.Name != "" {
		return o.Type.Name
	}
	return "struct"
}

func (o *Struct) String() string {
	var parts []string
	for _, f := range o.Type.Fields {
		val := o.Fields[f.Name]
		if val == nil {
			val = NullValue
		}
		parts = append(parts, fmt.Sprintf("%s:%s", f.Name, val.String()))
	}
	if o.Type.Name != "" {
		return fmt.Sprintf("%s{%s}", o.Type.Name, strings.Join(parts, " "))
	}
	return fmt.Sprintf("struct{%s}", strings.Join(parts, " "))
}

func (o *Struct) Copy() Object {
	newFields := make(map[string]Object)
	for k, v := range o.Fields {
		if v == nil {
			newFields[k] = NullValue
		} else {
			newFields[k] = v.Copy()
		}
	}
	return &Struct{
		Type:   o.Type,
		Fields: newFields,
	}
}

func (o *Struct) Equals(another Object) bool {
	t, ok := another.(*Struct)
	if !ok {
		return false
	}
	if !o.Type.Equals(t.Type) {
		return false
	}
	for k, v := range o.Fields {
		tv, ok := t.Fields[k]
		if !ok {
			return false
		}
		if v == nil {
			if tv != nil && tv != NullValue {
				return false
			}
		} else if !v.Equals(tv) {
			return false
		}
	}
	return true
}

func (o *Struct) IndexGet(index Object) (Object, error) {
	strIdx, ok := index.(*String)
	if !ok {
		return nil, ErrInvalidIndexType
	}
	found := false
	for _, f := range o.Type.Fields {
		if f.Name == strIdx.Value {
			found = true
			break
		}
	}
	if found {
		val := o.Fields[strIdx.Value]
		if val == nil {
			return NullValue, nil
		}
		return val, nil
	}

	if o.Type.Methods != nil {
		if method, exists := o.Type.Methods[strIdx.Value]; exists {
			return &BoundMethod{Receiver: o, Func: method}, nil
		}
	}

	return nil, fmt.Errorf("struct %s has no field or method %s", o.TypeName(), strIdx.Value)
}

func (o *Struct) IndexSet(index, value Object) error {
	strIdx, ok := index.(*String)
	if !ok {
		return ErrInvalidIndexType
	}
	var fieldField *StructField
	for i := range o.Type.Fields {
		if o.Type.Fields[i].Name == strIdx.Value {
			fieldField = &o.Type.Fields[i]
			break
		}
	}
	if fieldField == nil {
		return fmt.Errorf("struct %s has no field %s", o.TypeName(), strIdx.Value)
	}
	if fieldField.Type != "" {
		if !checkType(value, fieldField.Type) {
			return fmt.Errorf("cannot assign %s to field %s of type %s in struct %s", value.TypeName(), strIdx.Value, fieldField.Type, o.TypeName())
		}
	}
	o.Fields[strIdx.Value] = value
	return nil
}

func (o *Struct) GobEncode() (b []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(o.Type); err != nil {
		return nil, err
	}
	if err := enc.Encode(o.Fields); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (o *Struct) GobDecode(b []byte) (err error) {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&o.Type); err != nil {
		return err
	}
	if err := dec.Decode(&o.Fields); err != nil {
		return err
	}
	return nil
}

type BoundMethod struct {
	ObjectImpl
	Receiver Object
	Func     Object
}

func (o *BoundMethod) TypeName() string {
	return "bound-method"
}

func (o *BoundMethod) String() string {
	return "<bound-method>"
}

func (o *BoundMethod) Copy() Object {
	return &BoundMethod{Receiver: o.Receiver.Copy(), Func: o.Func.Copy()}
}

func (o *BoundMethod) Equals(another Object) bool {
	t, ok := another.(*BoundMethod)
	if !ok {
		return false
	}
	return o.Receiver.Equals(t.Receiver) && o.Func.Equals(t.Func)
}

func (o *BoundMethod) CanCall() bool {
	return true
}

func (o *BoundMethod) Call(args ...Object) (Object, error) {
	newArgs := make([]Object, len(args)+1)
	newArgs[0] = o.Receiver
	copy(newArgs[1:], args)
	return o.Func.Call(newArgs...)
}

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

// String provides formatted multi-line matrix representation.
func (m *Matrix[T]) String() string {
	if m.Rows == 0 || m.Cols == 0 {
		return "[]"
	}

	var isFloat bool
	var zero T
	switch any(zero).(type) {
	case float64:
		isFloat = true
	}

	formatVal := func(val T) string {
		if isFloat {
			f := any(val).(float64)
			if math.Abs(f-math.Round(f)) < 1e-9 {
				return fmt.Sprintf("%.0f", f)
			}
			return fmt.Sprintf("%.2f", f)
		}
		return fmt.Sprintf("%v", val)
	}

	colWidths := make([]int, m.Cols)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			valStr := formatVal(m.Data[i*m.Cols+j])
			if len(valStr) > colWidths[j] {
				colWidths[j] = len(valStr)
			}
		}
	}

	var sb strings.Builder
	for i := 0; i < m.Rows; i++ {
		sb.WriteString("│")
		for j := 0; j < m.Cols; j++ {
			valStr := formatVal(m.Data[i*m.Cols+j])
			sb.WriteString(valStr)
			padding := colWidths[j] - len(valStr)
			if padding > 0 {
				sb.WriteString(strings.Repeat(" ", padding))
			}
			if j < m.Cols-1 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("│")
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
			if m.Cols != rhsMat.Rows {
				return nil, fmt.Errorf("dimension mismatch: cannot multiply %dx%d and %dx%d", m.Rows, m.Cols, rhsMat.Rows, rhsMat.Cols)
			}
			res := &Matrix[T]{Rows: m.Rows, Cols: rhsMat.Cols, Data: make([]T, m.Rows*rhsMat.Cols)}
			for i := 0; i < m.Rows; i++ {
				for k := 0; k < m.Cols; k++ {
					temp := m.Data[i*m.Cols+k]
					for j := 0; j < rhsMat.Cols; j++ {
						res.Data[i*rhsMat.Cols+j] += temp * rhsMat.Data[k*rhsMat.Cols+j]
					}
				}
			}
			return res, nil

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

func (m *Matrix[T]) Transpose() *Matrix[T] {
	res := &Matrix[T]{Rows: m.Cols, Cols: m.Rows, Data: make([]T, len(m.Data))}
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			res.Data[j*m.Rows+i] = m.Data[i*m.Cols+j]
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
	
	a := make([][]complex128, n)
	for i := 0; i < n; i++ {
		a[i] = make([]complex128, n)
		for j := 0; j < n; j++ {
			a[i][j] = toComplex128(m.Data[i*n+j])
		}
	}
	
	det := complex128(1.0)
	for i := 0; i < n; i++ {
		pivot := i
		maxVal := cmplx.Abs(a[i][i])
		for k := i + 1; k < n; k++ {
			if abs := cmplx.Abs(a[k][i]); abs > maxVal {
				maxVal = abs
				pivot = k
			}
		}
		if maxVal < 1e-15 {
			return zero, nil
		}
		if pivot != i {
			a[i], a[pivot] = a[pivot], a[i]
			det = -det
		}
		for j := i + 1; j < n; j++ {
			factor := a[j][i] / a[i][i]
			for k := i; k < n; k++ {
				a[j][k] -= factor * a[i][k]
			}
		}
		det *= a[i][i]
	}
	
	var typedDet T
	switch any(zero).(type) {
	case int64:
		typedDet = any(int64(math.Round(real(det)))).(T)
	case float64:
		typedDet = any(real(det)).(T)
	case complex128:
		typedDet = any(det).(T)
	}
	return typedDet, nil
}

func (m *Matrix[T]) rank() int {
	if m.Rows == 0 || m.Cols == 0 {
		return 0
	}
	rows := m.Rows
	cols := m.Cols
	
	a := make([][]complex128, rows)
	for i := 0; i < rows; i++ {
		a[i] = make([]complex128, cols)
		for j := 0; j < cols; j++ {
			a[i][j] = toComplex128(m.Data[i*cols+j])
		}
	}
	
	rank := 0
	for col := 0; col < cols; col++ {
		pivot := -1
		for row := rank; row < rows; row++ {
			if cmplx.Abs(a[row][col]) > 1e-12 {
				pivot = row
				break
			}
		}
		if pivot == -1 {
			continue
		}
		a[rank], a[pivot] = a[pivot], a[rank]
		for row := rank + 1; row < rows; row++ {
			factor := a[row][col] / a[rank][col]
			for c := col; c < cols; c++ {
				a[row][c] -= factor * a[rank][c]
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