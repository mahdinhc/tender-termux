package tender

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"
)

func init() {
	addBuiltinFunction("govm", builtinGovm, true)
	addBuiltinFunction("abort", builtinAbort, true)
	addBuiltinFunction("makechan", builtinMakechan, false)
}

type ret struct {
	val Object
	err error
}

type goroutineVM struct {
	*VM      // if not nil, run Function in VM
	ret      // return value
	waitChan chan ret
	done     int64
}

// Starts a independent concurrent goroutine which runs fn(arg1, arg2, ...)
func builtinGovm(args ...Object) (Object, error) {
	vm := args[0].(*VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) == 0 {
		return nil, ErrWrongNumArguments
	}
	return spawnGoroutineInternal(vm, args[0], args[1:]...)
}

func spawnGoroutineInternal(vm *VM, fn Object, args ...Object) (Object, error) {
	if !fn.CanCall() {
		return nil, ErrInvalidArgumentType{
			Name:     "first",
			Expected: "callable function",
			Found:    fn.TypeName(),
		}
	}
	
	gvm := &goroutineVM{
		waitChan: make(chan ret, 1),
	}
	
	var callers []frame
	cfn, compiled := fn.(*Function)
	if compiled {
		gvm.VM = vm.ShallowClone()
	} else {
		callers = vm.callers()
	}
	
	if err := vm.addChild(gvm.VM); err != nil {
		return nil, err
	}
	go func() {
		var val Object
		var err error
		defer func() {
			if perr := recover(); perr != nil {
				if callers == nil {
					panic("callers not saved")
				}
				err = fmt.Errorf("\nRuntime Panic: %v%s\n%s", perr, vm.callStack(callers), debug.Stack())
			}
			if err != nil {
				vm.addError(err)
			}
			gvm.waitChan <- ret{val, err}
			vm.delChild(gvm.VM)
			gvm.VM = nil
		}()
		
		if cfn != nil {
			val, err = gvm.RunCompiled(cfn, args...)
		} else {
			var nargs []Object
			if bltnfn, ok := fn.(*NativeFunction); ok {
				if bltnfn.NeedVMObj {
					// pass VM as the first para to builtin functions
					nargs = append(nargs, vm.selfObject())
				}
			}
			nargs = append(nargs, args...)
			val, err = fn.Call(nargs...)
		}
	}()
	
	return &Goroutine{gvm: gvm}, nil
}

// Triggers the termination process of the current VM and all its descendant VMs.
func builtinAbort(args ...Object) (Object, error) {
	vm := args[0].(*VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}
	vm.Abort() // aborts self and all descendant VMs
	return nil, nil
}

// Returns true if the goroutineVM is done
func (gvm *goroutineVM) wait(seconds int64) bool {
	if atomic.LoadInt64(&gvm.done) == 1 {
		return true
	}
	
	if seconds < 0 {
		seconds = 3153600000 // 100 years
	}
	
	select {
	case gvm.ret = <-gvm.waitChan:
		atomic.StoreInt64(&gvm.done, 1)
	case <-time.After(time.Duration(seconds) * time.Second):
		return false
	}
	
	return true
}

// Waits for the goroutineVM to complete up to timeout seconds.
func (gvm *goroutineVM) waitTimeout(args ...Object) (Object, error) {
	if len(args) > 1 {
		return nil, ErrWrongNumArguments
	}
	timeOut := -1
	if len(args) == 1 {
		t, ok := ToInt(args[0])
		if !ok {
			return nil, ErrInvalidArgumentType{
				Name:     "first",
				Expected: "int(compatible)",
				Found:    args[0].TypeName(),
			}
		}
		timeOut = t
	}
	
	if gvm.wait(int64(timeOut)) {
		return TrueValue, nil
	}
	return FalseValue, nil
}

// Triggers the termination process of the goroutineVM and all its descendant VMs.
func (gvm *goroutineVM) abort(args ...Object) (Object, error) {
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}
	if gvm.VM != nil {
		gvm.Abort()
	}
	return nil, nil
}

// Waits the goroutineVM to complete, return Error object if any runtime error occurred.
func (gvm *goroutineVM) getRet(args ...Object) (Object, error) {
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}
	
	gvm.wait(-1)
	if gvm.ret.err != nil {
		return &Error{Value: &String{Value: gvm.ret.err.Error()}}, nil
	}
	
	return gvm.ret.val, nil
}

type objchan chan Object

// Channel represents a first-class channel object in Tender.
type Channel struct {
	ObjectImpl
	Value chan Object
	size  int
}

func (c *Channel) TypeName() string {
	return "channel"
}

func (c *Channel) String() string {
	return fmt.Sprintf("<channel cap=%d>", cap(c.Value))
}

func (c *Channel) Copy() Object {
	return c
}

func (c *Channel) IsFalsy() bool {
	return c.Value == nil
}

func (c *Channel) Equals(another Object) bool {
	o, ok := another.(*Channel)
	if !ok {
		return false
	}
	return c.Value == o.Value
}

func (c *Channel) IndexGet(index Object) (Object, error) {
	strIdx, ok := index.(*String)
	if !ok {
		return nil, ErrInvalidIndexType
	}
	switch strIdx.Value {
	case "send":
		oc := objchan(c.Value)
		return &NativeFunction{Value: oc.send, NeedVMObj: true}, nil
	case "recv":
		oc := objchan(c.Value)
		return &NativeFunction{Value: oc.recv, NeedVMObj: true}, nil
	case "close":
		oc := objchan(c.Value)
		return &NativeFunction{Value: oc.close}, nil
	case "cap":
		return &Int{Value: int64(cap(c.Value))}, nil
	case "len":
		return &Int{Value: int64(len(c.Value))}, nil
	}
	return nil, fmt.Errorf("channel has no field or method %s", strIdx.Value)
}

// Goroutine represents a handle to a running goroutineVM.
type Goroutine struct {
	ObjectImpl
	gvm *goroutineVM
}

func (g *Goroutine) TypeName() string {
	return "goroutine"
}

func (g *Goroutine) String() string {
	return "<goroutine>"
}

func (g *Goroutine) Copy() Object {
	return g
}

func (g *Goroutine) IsFalsy() bool {
	return g.gvm == nil
}

func (g *Goroutine) Equals(another Object) bool {
	o, ok := another.(*Goroutine)
	if !ok {
		return false
	}
	return g.gvm == o.gvm
}

func (g *Goroutine) IndexGet(index Object) (Object, error) {
	strIdx, ok := index.(*String)
	if !ok {
		return nil, ErrInvalidIndexType
	}
	switch strIdx.Value {
	case "result":
		return &NativeFunction{Value: g.gvm.getRet}, nil
	case "wait":
		return &NativeFunction{Value: g.gvm.waitTimeout}, nil
	case "abort":
		return &NativeFunction{Value: g.gvm.abort}, nil
	}
	return nil, fmt.Errorf("goroutine has no field or method %s", strIdx.Value)
}

// Makes a channel to send/receive object
func builtinMakechan(args ...Object) (Object, error) {
	var size int
	switch len(args) {
	case 0:
	case 1:
		n, ok := ToInt(args[0])
		if !ok {
			return nil, ErrInvalidArgumentType{
				Name:     "first",
				Expected: "int(compatible)",
				Found:    args[0].TypeName(),
			}
		}
		size = n
	default:
		return nil, ErrWrongNumArguments
	}
	
	oc := make(objchan, size)
	return &Channel{Value: oc, size: size}, nil
}

// Sends an obj to the channel, will block if channel is full and the VM has not been aborted.
func (oc objchan) send(args ...Object) (Object, error) {
	vm := args[0].(*VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}
	select {
	case <-vm.AbortChan:
		return nil, ErrVMAborted
	case oc <- args[0]:
	}
	return nil, nil
}

// Receives an obj from the channel, will block if channel is empty and the VM has not been aborted.
func (oc objchan) recv(args ...Object) (Object, error) {
	vm := args[0].(*VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}
	select {
	case <-vm.AbortChan:
		return nil, ErrVMAborted
	case obj, ok := <-oc:
		if ok {
			return obj, nil
		}
	}
	return nil, nil
}

// Closes the channel.
func (oc objchan) close(args ...Object) (Object, error) {
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}
	close(oc)
	return nil, nil
}

// FastShallowClone creates a lightweight shallow VM clone for synchronous call execution
// without childCtl lock contention or extra allocation overhead.
func FastShallowClone(v *VM) *VM {
	vClone := &VM{
		constants:   v.constants,
		sp:          0,
		globals:     v.globals,
		fileSet:     v.fileSet,
		frames:      make([]*frame, 0, initialFrames),
		framesIndex: 1,
		ip:          -1,
		maxAllocs:   v.maxAllocs,
		AbortChan:   v.AbortChan,
		In:          v.In,
		Out:         v.Out,
		Args:        v.Args,
	}
	vClone.frames = append(vClone.frames, &frame{
		fn: emptyEntry,
		ip: -1,
	})
	return vClone
}

// WrapFuncCall synchronously executes a callable object from Go space.
func WrapFuncCall(vm *VM, args ...Object) (retVal Object, err error) {
	if len(args) == 0 {
		return nil, ErrWrongNumArguments
	}
	fn := args[0]
	if !fn.CanCall() {
		return nil, ErrInvalidArgumentType{Name: "first", Expected: "callable function", Found: fn.TypeName()}
	}
	
	defer func() {
		if perr := recover(); perr != nil {
			err = fmt.Errorf("\nRuntime Panic within Native Bridge: %v\n%s", perr, debug.Stack())
		}
	}()

	cfn, compiled := fn.(*Function)
	if compiled {
		clone := FastShallowClone(vm)
		return clone.RunCompiled(cfn, args[1:]...)
	}
	
	var nargs []Object
	if bltnfn, ok := fn.(*NativeFunction); ok {
		if bltnfn.NeedVMObj {
			nargs = append(nargs, vm.selfObject())
		}
	}
	nargs = append(nargs, args[1:]...)
	return fn.Call(nargs...)
}