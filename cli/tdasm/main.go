// cmd/tdasm/main.go
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"encoding/gob"

	"github.com/2dprototype/tender"
	"github.com/2dprototype/tender/parser"
	"github.com/2dprototype/tender/stdlib"
)

var (
	disassembleFlag bool
	assembleFlag    bool
	outputFile      string
	showHelp        bool
)

func init() {
	gob.Register(&parser.SourceFileSet{})
	gob.Register(&parser.SourceFile{})
	gob.Register(&tender.Array{})
	gob.Register(&tender.Bool{})
	gob.Register(&tender.Bytes{})
	gob.Register(&tender.Char{})
	gob.Register(&tender.Function{})
	gob.Register(&tender.Error{})
	gob.Register(&tender.Int{})
	gob.Register(&tender.Float{})
	gob.Register(&tender.BigInt{})
	gob.Register(&tender.BigFloat{})
	gob.Register(&tender.Complex{})
	gob.Register(&tender.ImmutableArray{})
	gob.Register(&tender.ImmutableMap{})
	gob.Register(&tender.Map{})
	gob.Register(&tender.String{})
	gob.Register(&tender.Time{})
	gob.Register(&tender.Null{})
	gob.Register(&tender.NativeFunction{})
	gob.Register(&tender.Tuple{})
	gob.Register(&tender.Struct{})
	gob.Register(&tender.StructType{})
	gob.Register(&tender.BoundMethod{})
	gob.Register(&tender.Matrix[int64]{})
	gob.Register(&tender.Matrix[float64]{})
	gob.Register(&tender.Matrix[complex128]{})
	gob.Register(&tender.ObjectPtr{})
	gob.Register(&tender.Channel{})
	gob.Register(&stdlib.IOWriter{})
	gob.Register(&stdlib.IOReader{})
	flag.BoolVar(&disassembleFlag, "disassemble", false, "Disassemble bytecode to .tasm")
	flag.BoolVar(&assembleFlag, "assemble", false, "Assemble .tasm to bytecode")
	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.BoolVar(&showHelp, "help", false, "Show help")
}

func main() {
	flag.Parse()
	if showHelp || (!disassembleFlag && !assembleFlag) {
		printHelp()
		return
	}
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: input file required")
		printHelp()
		os.Exit(1)
	}
	inputFile := flag.Arg(0)

	if disassembleFlag {
		if err := doDisassemble(inputFile, outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Disassembly error: %v\n", err)
			os.Exit(1)
		}
	} else if assembleFlag {
		if err := doAssemble(inputFile, outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Assembly error: %v\n", err)
			os.Exit(1)
		}
	}
}

func printHelp() {
	fmt.Println(`tdasm - Tender Disassembler / Assembler

Usage:
  tdasm -disassemble <input> [-o <output>]
  tdasm -assemble <input> [-o <output>]

Flags:
  -disassemble   Disassemble a bytecode file (.tdc/.tdo) or source (.td) to .tasm
  -assemble      Assemble a .tasm file back to bytecode (.tdc)
  -o <file>      Output file (default: input with changed extension)
  -help          Show this help

Examples:
  tdasm -disassemble myapp.tdc -o myapp.tasm
  tdasm -assemble myapp.tasm -o myapp.tdc
  tdasm -disassemble myapp.td          # compile first, then disassemble
`)
}

func doDisassemble(inputFile, outFile string) error {
	ext := filepath.Ext(inputFile)
	var bytecode *tender.Bytecode
	var err error

	switch ext {
	case ".tdc", ".tdo":
		bytecode, err = loadBytecode(inputFile)
	case ".td":
		bytecode, err = compileSource(inputFile)
	default:
		bytecode, err = loadBytecode(inputFile)
		if err != nil {
			bytecode, err = compileSource(inputFile)
		}
	}
	if err != nil {
		return err
	}

	asmText := disassembleBytecode(bytecode)
	if outFile == "" {
		outFile = strings.TrimSuffix(inputFile, ext) + ".tasm"
	}
	return os.WriteFile(outFile, []byte(asmText), 0644)
}

func doAssemble(inputFile, outFile string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}
	bytecode, err := assembleBytecode(string(data))
	if err != nil {
		return err
	}
	if outFile == "" {
		outFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".tdc"
	}
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return bytecode.EncodeTDC(f, "v1.0.0")
}

// ---- Disassembly ----

func disassembleBytecode(bc *tender.Bytecode) string {
	var out strings.Builder
	fmt.Fprintf(&out, "; Tender Assembly\n")
	fmt.Fprintf(&out, "; version: %s\n", "v1.0.0")
	fmt.Fprintf(&out, "; compiled: %s\n", time.Now().Format(time.RFC3339))
	out.WriteString("\n")

	out.WriteString("; ---- Constant Pool ----\n")
	for i, c := range bc.Constants {
		fmt.Fprintf(&out, "const %d = %s\n", i, formatConst(c))
	}
	out.WriteString("\n")

	out.WriteString("; ---- Function Definitions ----\n")
	writeFunc(&out, "main", bc.MainFunction)
	for i, c := range bc.Constants {
		if fn, ok := c.(*tender.Function); ok && fn != bc.MainFunction {
			fmt.Fprintf(&out, "\n.func fn_const_%d\n", i)
			writeFuncBody(&out, fn)
			out.WriteString(".endfunc\n")
		}
	}
	return out.String()
}

func writeFunc(out *strings.Builder, label string, fn *tender.Function) {
	fmt.Fprintf(out, ".func %s\n", label)
	writeFuncBody(out, fn)
	out.WriteString(".endfunc\n")
}

func writeFuncBody(out *strings.Builder, fn *tender.Function) {
	lines := tender.FormatInstructions(fn.Instructions, 0)
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		opcode := parts[1]
		operands := parts[2:]
		fmt.Fprintf(out, "    %s", opcode)
		for _, op := range operands {
			fmt.Fprintf(out, " %s", op)
		}
		out.WriteString("\n")
	}
}

func formatConst(c tender.Object) string {
	if c == nil {
		return "null"
	}
	switch v := c.(type) {
	case *tender.Int:
		return fmt.Sprintf("%d", v.Value)
	case *tender.Float:
		return fmt.Sprintf("%g", v.Value)
	case *tender.BigInt:
		return fmt.Sprintf("bigint(%q)", v.Value.String())
	case *tender.BigFloat:
		return fmt.Sprintf("bigfloat(%q)", v.Value.String())
	case *tender.Complex:
		return fmt.Sprintf("complex(%g, %g)", real(v.Value), imag(v.Value))
	case *tender.String:
		return strconv.Quote(v.Value)
	case *tender.Bytes:
		if isPrintable(v.Value) {
			return fmt.Sprintf("bytes(%q)", string(v.Value))
		}
		var parts []string
		for _, b := range v.Value {
			parts = append(parts, fmt.Sprintf("%d", b))
		}
		return fmt.Sprintf("bytes(%s)", strings.Join(parts, ", "))
	case *tender.Char:
		return fmt.Sprintf("'%c'", v.Value)
	case *tender.Bool:
		if v.IsFalsy() {
			return "false"
		}
		return "true"
	case *tender.Null:
		return "null"
	case *tender.Array:
		var elems []string
		for _, e := range v.Value {
			elems = append(elems, formatConst(e))
		}
		return "[" + strings.Join(elems, ", ") + "]"
	case *tender.Map:
		var pairs []string
		for k, val := range v.Value {
			pairs = append(pairs, fmt.Sprintf("%q: %s", k, formatConst(val)))
		}
		return "{" + strings.Join(pairs, ", ") + "}"
	case *tender.ImmutableMap:
		if nameObj, ok := v.Value["__module_name__"]; ok {
			if s, ok := nameObj.(*tender.String); ok {
				return fmt.Sprintf("immutable_map({__module_name__: %q, ...})", s.Value)
			}
		}
		var pairs []string
		for k, val := range v.Value {
			if k == "__module_name__" {
				continue
			}
			pairs = append(pairs, fmt.Sprintf("%q: %s", k, formatConst(val)))
		}
		return "immutable_map({" + strings.Join(pairs, ", ") + "})"
	case *tender.StructType:
		var fields []string
		for _, f := range v.Fields {
			tag := ""
			if f.Tag != "" {
				tag = " " + strconv.Quote(f.Tag)
			}
			fields = append(fields, fmt.Sprintf("%s %s%s", f.Name, f.Type, tag))
		}
		return "struct { " + strings.Join(fields, "; ") + " }"
	case *tender.Function:
		return "<function>"
	default:
		return v.String()
	}
}

func isPrintable(b []byte) bool {
	for _, c := range b {
		if c < 32 || c > 126 {
			return false
		}
	}
	return true
}

// ---- Assembly ----

type asmInstruction struct {
	opcode   parser.Opcode
	operands []int
}

type asmFunction struct {
	label        string
	instructions []asmInstruction
}

func assembleBytecode(source string) (*tender.Bytecode, error) {
	scanner := bufio.NewScanner(strings.NewReader(source))
	var constants []tender.Object
	var functions []*asmFunction
	var currentFunc *asmFunction
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if idx := strings.Index(line, ";"); idx >= 0 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "const ") {
			rest := strings.TrimPrefix(line, "const ")
			parts := strings.SplitN(rest, "=", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("line %d: invalid const definition", lineNum)
			}
			idxStr := strings.TrimSpace(parts[0])
			valStr := strings.TrimSpace(parts[1])
			idx, err := strconv.Atoi(idxStr)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid constant index: %v", lineNum, err)
			}
			val, err := parseConstant(valStr)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid constant value: %v", lineNum, err)
			}
			if idx != len(constants) {
				return nil, fmt.Errorf("line %d: constant index %d out of order (expected %d)", lineNum, idx, len(constants))
			}
			constants = append(constants, val)
		} else if strings.HasPrefix(line, ".func ") {
			label := strings.TrimPrefix(line, ".func ")
			if currentFunc != nil {
				return nil, fmt.Errorf("line %d: nested .func without .endfunc", lineNum)
			}
			currentFunc = &asmFunction{label: label}
		} else if line == ".endfunc" {
			if currentFunc == nil {
				return nil, fmt.Errorf("line %d: .endfunc without .func", lineNum)
			}
			functions = append(functions, currentFunc)
			currentFunc = nil
		} else {
			if currentFunc == nil {
				return nil, fmt.Errorf("line %d: instruction outside function", lineNum)
			}
			fields := strings.Fields(line)
			if len(fields) == 0 {
				continue
			}
			opName := fields[0]
			opcode, err := lookupOpcode(opName)
			if err != nil {
				return nil, fmt.Errorf("line %d: unknown opcode %q", lineNum, opName)
			}
			var operands []int
			for _, f := range fields[1:] {
				op, err := strconv.Atoi(f)
				if err != nil {
					return nil, fmt.Errorf("line %d: invalid operand %q", lineNum, f)
				}
				operands = append(operands, op)
			}
			expected := parser.OpcodeOperands[opcode]
			if len(operands) != len(expected) {
				return nil, fmt.Errorf("line %d: operand count mismatch for %s: got %d, want %d", lineNum, opName, len(operands), len(expected))
			}
			currentFunc.instructions = append(currentFunc.instructions, asmInstruction{opcode: opcode, operands: operands})
		}
	}
	if currentFunc != nil {
		return nil, fmt.Errorf("unclosed .func at EOF")
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var mainFunc *tender.Function
	for _, f := range functions {
		if f.label == "main" {
			mainFunc = buildFunction(f, constants)
		} else {
			idx := -1
			if strings.HasPrefix(f.label, "fn_const_") {
				s := strings.TrimPrefix(f.label, "fn_const_")
				if i, err := strconv.Atoi(s); err == nil {
					idx = i
				}
			}
			if idx < 0 {
				return nil, fmt.Errorf("function label %q not recognized as constant index", f.label)
			}
			fn := buildFunction(f, constants)
			if idx >= len(constants) {
				for len(constants) <= idx {
					constants = append(constants, tender.NullValue)
				}
			}
			constants[idx] = fn
		}
	}
	if mainFunc == nil {
		return nil, fmt.Errorf("no .func main found")
	}

	return &tender.Bytecode{
		MainFunction: mainFunc,
		Constants:    constants,
		FileSet:      parser.NewFileSet(),
	}, nil
}

func buildFunction(f *asmFunction, constants []tender.Object) *tender.Function {
	var insts []byte
	for _, inst := range f.instructions {
		insts = append(insts, tender.MakeInstruction(inst.opcode, inst.operands...)...)
	}
	return &tender.Function{
		Instructions:  insts,
		NumLocals:     0,
		NumParameters: 0,
		VarArgs:       false,
		SourceMap:     nil,
	}
}

func lookupOpcode(name string) (parser.Opcode, error) {
	for i, n := range parser.OpcodeNames {
		if n == name {
			return parser.Opcode(i), nil
		}
	}
	return 0, fmt.Errorf("unknown opcode %q", name)
}

func parseConstant(s string) (tender.Object, error) {
	s = strings.TrimSpace(s)
	if s == "true" {
		return tender.TrueValue, nil
	}
	if s == "false" {
		return tender.FalseValue, nil
	}
	if s == "null" {
		return tender.NullValue, nil
	}
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		unquoted, err := strconv.Unquote(s)
		if err != nil {
			return nil, err
		}
		return &tender.String{Value: unquoted}, nil
	}
	if len(s) >= 3 && s[0] == '\'' && s[2] == '\'' {
		return &tender.Char{Value: rune(s[1])}, nil
	}
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		inner := strings.TrimSpace(s[1 : len(s)-1])
		var elems []tender.Object
		if inner != "" {
			parts := splitTopLevel(inner, ',')
			for _, p := range parts {
				val, err := parseConstant(p)
				if err != nil {
					return nil, err
				}
				elems = append(elems, val)
			}
		}
		return &tender.Array{Value: elems}, nil
	}
	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		inner := strings.TrimSpace(s[1 : len(s)-1])
		m := make(map[string]tender.Object)
		if inner != "" {
			parts := splitTopLevel(inner, ',')
			for _, p := range parts {
				kv := strings.SplitN(p, ":", 2)
				if len(kv) != 2 {
					return nil, fmt.Errorf("invalid map entry: %q", p)
				}
				keyStr := strings.TrimSpace(kv[0])
				valStr := strings.TrimSpace(kv[1])
				if len(keyStr) < 2 || keyStr[0] != '"' || keyStr[len(keyStr)-1] != '"' {
					return nil, fmt.Errorf("map key must be a string: %q", keyStr)
				}
				key, err := strconv.Unquote(keyStr)
				if err != nil {
					return nil, err
				}
				val, err := parseConstant(valStr)
				if err != nil {
					return nil, err
				}
				m[key] = val
			}
		}
		return &tender.Map{Value: m}, nil
	}
	if strings.HasPrefix(s, "complex(") && strings.HasSuffix(s, ")") {
		inner := s[8 : len(s)-1]
		parts := strings.Split(inner, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid complex: %q", s)
		}
		re, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return nil, err
		}
		im, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return nil, err
		}
		return &tender.Complex{Value: complex(re, im)}, nil
	}
	if strings.HasPrefix(s, "bigint(") && strings.HasSuffix(s, ")") {
		inner := s[7 : len(s)-1]
		if len(inner) >= 2 && inner[0] == '"' && inner[len(inner)-1] == '"' {
			inner, _ = strconv.Unquote(inner)
		}
		bi := new(tender.BigInt)
		bi.Value = new(big.Int)
		_, ok := bi.Value.SetString(inner, 10)
		if !ok {
			return nil, fmt.Errorf("invalid bigint: %q", inner)
		}
		return bi, nil
	}
	if strings.HasPrefix(s, "bigfloat(") && strings.HasSuffix(s, ")") {
		inner := s[9 : len(s)-1]
		if len(inner) >= 2 && inner[0] == '"' && inner[len(inner)-1] == '"' {
			inner, _ = strconv.Unquote(inner)
		}
		bf := new(tender.BigFloat)
		bf.Value = new(big.Float)
		_, ok := bf.Value.SetString(inner)
		if !ok {
			return nil, fmt.Errorf("invalid bigfloat: %q", inner)
		}
		return bf, nil
	}
	if strings.HasPrefix(s, "bytes(") && strings.HasSuffix(s, ")") {
		inner := s[6 : len(s)-1]
		inner = strings.TrimSpace(inner)
		if len(inner) >= 2 && inner[0] == '"' && inner[len(inner)-1] == '"' {
			unquoted, err := strconv.Unquote(inner)
			if err != nil {
				return nil, err
			}
			return &tender.Bytes{Value: []byte(unquoted)}, nil
		}
		var values []byte
		parts := strings.Split(inner, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			v, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			if v < 0 || v > 255 {
				return nil, fmt.Errorf("byte value out of range: %d", v)
			}
			values = append(values, byte(v))
		}
		return &tender.Bytes{Value: values}, nil
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return &tender.Int{Value: i}, nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return &tender.Float{Value: f}, nil
	}
	return nil, fmt.Errorf("unable to parse constant: %q", s)
}

func splitTopLevel(s string, sep rune) []string {
	var parts []string
	var current strings.Builder
	level := 0
	inQuote := false
	for _, ch := range s {
		switch {
		case ch == '"':
			inQuote = !inQuote
			current.WriteRune(ch)
		case inQuote:
			current.WriteRune(ch)
		case ch == '[' || ch == '{' || ch == '(':
			level++
			current.WriteRune(ch)
		case ch == ']' || ch == '}' || ch == ')':
			level--
			current.WriteRune(ch)
		case ch == sep && level == 0:
			parts = append(parts, strings.TrimSpace(current.String()))
			current.Reset()
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}
	return parts
}

// ---- Helpers ----

func loadBytecode(path string) (*tender.Bytecode, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	bc := &tender.Bytecode{}
	if len(data) >= 4 && data[0] == 'T' && data[1] == 'D' && data[2] == 'C' && data[3] == 1 {
		_, _, err = bc.DecodeTDC(bytes.NewReader(data), stdlib.GetModuleMap())
		if err == nil {
			return bc, nil
		}
	}
	err = bc.Decode(bytes.NewReader(data), stdlib.GetModuleMap())
	return bc, err
}

func compileSource(path string) (*tender.Bytecode, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fileSet := parser.NewFileSet()
	srcFile := fileSet.AddFile(filepath.Base(path), -1, len(src))
	p := parser.NewParser(srcFile, src, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, err
	}
	c := tender.NewCompiler(srcFile, nil, nil, stdlib.GetModuleMap(), nil)
	c.EnableFileImport(true)
	c.SetImportDir(filepath.Dir(path))
	if err := c.Compile(file); err != nil {
		return nil, err
	}
	bc := c.Bytecode()
	bc.RemoveDuplicates()
	return bc, nil
}