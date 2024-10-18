package stdlib

import (
	"os"
	"syscall"

	"github.com/2dprototype/tender"
)

func makeOSProcessState(state *os.ProcessState) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"exited": &tender.UserFunction{
				Name:  "exited",
				Value: FuncARB(state.Exited),
			},
			"pid": &tender.UserFunction{
				Name:  "pid",
				Value: FuncARI(state.Pid),
			},
			"string": &tender.UserFunction{
				Name:  "string",
				Value: FuncARS(state.String),
			},
			"success": &tender.UserFunction{
				Name:  "success",
				Value: FuncARB(state.Success),
			},
		},
	}
}

func makeOSProcess(proc *os.Process) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"kill": &tender.UserFunction{
				Name:  "kill",
				Value: FuncARE(proc.Kill),
			},
			"release": &tender.UserFunction{
				Name:  "release",
				Value: FuncARE(proc.Release),
			},
			"signal": &tender.UserFunction{
				Name: "signal",
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					i1, ok := tender.ToInt64(args[0])
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "int(compatible)",
							Found:    args[0].TypeName(),
						}
					}
					return wrapError(proc.Signal(syscall.Signal(i1))), nil
				},
			},
			"wait": &tender.UserFunction{
				Name: "wait",
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					state, err := proc.Wait()
					if err != nil {
						return wrapError(err), nil
					}
					return makeOSProcessState(state), nil
				},
			},
		},
	}
}
