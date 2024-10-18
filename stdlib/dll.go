package stdlib

import (
	"syscall"
	"github.com/2dprototype/tender"
	"unsafe"
)

var dllModule = map[string]tender.Object{
	"call_dll": &tender.UserFunction{Name: "call_dll", Value: callDLLWrapper},
	"new": &tender.UserFunction{Name: "new", Value: dllNew},
	"load": &tender.UserFunction{Name: "load", Value: dllLoad},
}


func dllLoad(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	
	dllName, ok := args[0].(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "dll_name",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	
	dll, err := syscall.LoadDLL(dllName.Value)
	
	if err != nil {
		return wrapError(err), nil
	}
	
	return makeFindProc(dll), nil
}

func makeFindProc(dll *syscall.DLL) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"proc": &tender.UserFunction{
				Name:  "proc",
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					functionName, ok := args[0].(*tender.String)
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     "function_name",
							Expected: "string",
							Found:    args[0].TypeName(),
						}
					}
					proc, err := dll.FindProc(functionName.Value)
					if err != nil {
						return wrapError(err), nil
					}
					return makeLoadProcCall(proc), nil
				},
			},
		},
	}
}

func makeLoadProcCall(proc *syscall.Proc) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"call": &tender.UserFunction{
				Name:  "call",
				Value: func(args ...tender.Object) (tender.Object, error) {
					var uintptrArgs []uintptr
					for _, arg := range args {
						switch val := arg.(type) {
							case *tender.Int:
								uintptrArgs = append(uintptrArgs, uintptr(val.Value))
							case *tender.Float:
								uintptrArgs = append(uintptrArgs, uintptr(val.Value))
							case *tender.String:
								uintptrArgs = append(uintptrArgs, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(val.Value))))
							case *tender.Bytes:
								uintptrArgs = append(uintptrArgs, uintptr(unsafe.Pointer(&val.Value[0])))
							case *tender.Null:
								uintptrArgs = append(uintptrArgs, 0)
							default:
								return nil, tender.ErrInvalidArgumentType{
									Name:     "argument",
									Expected: "integer, float, or string",
									Found:    arg.TypeName(),
							}
						}
					}
					
					ret, _, _ := proc.Call(uintptrArgs...)
					if ret == 0 {
						return &tender.Int{Value: 0}, nil
					}
					
					return  &tender.Int{Value: int64(ret)}, nil
				},
			},
		},
	}
}




func dllNew(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	
	dllName, ok := args[0].(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "dll_name",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	
	dll := syscall.NewLazyDLL(dllName.Value)
	
	return makeNewProc(dll), nil
}

func makeNewProc(dll *syscall.LazyDLL) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"proc": &tender.UserFunction{
				Name:  "proc",
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					functionName, ok := args[0].(*tender.String)
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     "function_name",
							Expected: "string",
							Found:    args[0].TypeName(),
						}
					}
					// println(functionName)
					// println(dll)
					proc := dll.NewProc(functionName.Value)
					return makeNewProcCall(proc), nil
				},
			},
		},
	}
}

func makeNewProcCall(proc *syscall.LazyProc) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"call": &tender.UserFunction{
				Name:  "call",
				Value: func(args ...tender.Object) (tender.Object, error) {
					var uintptrArgs []uintptr
					for _, arg := range args {
						switch val := arg.(type) {
							case *tender.Int:
							uintptrArgs = append(uintptrArgs, uintptr(val.Value))
							case *tender.Float:
							uintptrArgs = append(uintptrArgs, uintptr(val.Value))
							case *tender.String:
							uintptrArgs = append(uintptrArgs, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(val.Value))))
							case *tender.Null:
							uintptrArgs = append(uintptrArgs, 0)
							default:
							return nil, tender.ErrInvalidArgumentType{
								Name:     "argument",
								Expected: "integer, float, or string",
								Found:    arg.TypeName(),
							}
						}
					}
					
					ret, _, _ := proc.Call(uintptrArgs...)
					if ret == 0 {
						return &tender.Int{Value: 0}, nil
					}
					
					return  &tender.Int{Value: int64(ret)}, nil
				},
			},
		},
	}
}



// callDLL calls a function from a DLL dynamically
func callDLL(dllName, functionName string, args ...uintptr) (ret uintptr, err error) {
	dll := syscall.NewLazyDLL(dllName)
	proc := dll.NewProc(functionName)
	
	ret, _, err = proc.Call(args...)
	if ret == 0 {
		return 0, err
	}
	return ret, nil
}

// callDLLWrapper is a wrapper function for callDLL to use in the tender environment
func callDLLWrapper(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) < 2 {
		return nil, tender.ErrWrongNumArguments
	}
	
	dllName, ok := args[0].(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "dll_name",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	
	functionName, ok := args[1].(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "function_name",
			Expected: "string",
			Found:    args[1].TypeName(),
		}
	}
	
	var uintptrArgs []uintptr
	for _, arg := range args[2:] {
		switch val := arg.(type) {
			case *tender.Int:
			uintptrArgs = append(uintptrArgs, uintptr(val.Value))
			case *tender.Float:
			uintptrArgs = append(uintptrArgs, uintptr(val.Value))
			case *tender.String:
			uintptrArgs = append(uintptrArgs, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(val.Value))))
			case *tender.Null:
			uintptrArgs = append(uintptrArgs, 0)
			default:
			return nil, tender.ErrInvalidArgumentType{
				Name:     "argument",
				Expected: "integer, float, or string",
				Found:    arg.TypeName(),
			}
		}
	}
	
	result, err := callDLL(dllName.Value, functionName.Value, uintptrArgs...)
	if err != nil {
		return nil, err
	}
	
	return &tender.Int{Value: int64(result)}, nil
}
