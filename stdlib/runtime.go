package stdlib

import (
	"runtime"

	"github.com/2dprototype/tender"
)

var runtimeModule = map[string]tender.Object{
	"goos":     &tender.String{Value: runtime.GOOS},
	"goarch":   &tender.String{Value: runtime.GOARCH},
	"compiler": &tender.String{Value: runtime.Compiler},
	"version": &tender.NativeFunction{
		Name: "version",
		Value: func(args ...tender.Object) (tender.Object, error) {
			if len(args) != 0 {
				return nil, tender.ErrWrongNumArguments
			}
			return &tender.String{Value: runtime.Version()}, nil
		},
	},
	"num_cpu": &tender.NativeFunction{
		Name: "num_cpu",
		Value: func(args ...tender.Object) (tender.Object, error) {
			if len(args) != 0 {
				return nil, tender.ErrWrongNumArguments
			}
			return &tender.Int{Value: int64(runtime.NumCPU())}, nil
		},
	},
	"num_goroutine": &tender.NativeFunction{
		Name: "num_goroutine",
		Value: func(args ...tender.Object) (tender.Object, error) {
			if len(args) != 0 {
				return nil, tender.ErrWrongNumArguments
			}
			return &tender.Int{Value: int64(runtime.NumGoroutine())}, nil
		},
	},
	"num_cgo_call": &tender.NativeFunction{
		Name: "num_cgo_call",
		Value: func(args ...tender.Object) (tender.Object, error) {
			if len(args) != 0 {
				return nil, tender.ErrWrongNumArguments
			}
			return &tender.Int{Value: runtime.NumCgoCall()}, nil
		},
	},
	"gc": &tender.NativeFunction{
		Name: "gc",
		Value: func(args ...tender.Object) (tender.Object, error) {
			if len(args) != 0 {
				return nil, tender.ErrWrongNumArguments
			}
			runtime.GC()
			return tender.NullValue, nil
		},
	},
	"gosched": &tender.NativeFunction{
		Name: "gosched",
		Value: func(args ...tender.Object) (tender.Object, error) {
			if len(args) != 0 {
				return nil, tender.ErrWrongNumArguments
			}
			runtime.Gosched()
			return tender.NullValue, nil
		},
	},
	"goexit": &tender.NativeFunction{
		Name: "goexit",
		Value: func(args ...tender.Object) (tender.Object, error) {
			if len(args) != 0 {
				return nil, tender.ErrWrongNumArguments
			}
			runtime.Goexit()
			return tender.NullValue, nil
		},
	},
	"keep_alive": &tender.NativeFunction{
		Name: "keep_alive",
		Value: func(args ...tender.Object) (tender.Object, error) {
			if len(args) != 1 {
				return nil, tender.ErrWrongNumArguments
			}
			runtime.KeepAlive(args[0])
			return tender.NullValue, nil
		},
	},
	"gomaxprocs": &tender.NativeFunction{
		Name: "gomaxprocs",
		Value: func(args ...tender.Object) (tender.Object, error) {
			if len(args) != 1 {
				return nil, tender.ErrWrongNumArguments
			}
			n, ok := tender.ToInt(args[0])
			if !ok {
				return nil, tender.ErrInvalidArgumentType{
					Name:     "first",
					Expected: "int(compatible)",
					Found:    args[0].TypeName(),
				}
			}
			prev := runtime.GOMAXPROCS(n)
			return &tender.Int{Value: int64(prev)}, nil
		},
	},
	"lock_os_thread": &tender.NativeFunction{
		Name: "lock_os_thread",
		Value: FuncAR(runtime.LockOSThread),
	},
	"unlock_os_thread": &tender.NativeFunction{
		Name: "unlock_os_thread",
		Value: FuncAR(runtime.UnlockOSThread),
	},
	"read_mem_stats": &tender.NativeFunction{
		Name: "read_mem_stats",
		Value: runtimeReadMemStats,
	},
}

func runtimeReadMemStats(args ...tender.Object) (tender.Object, error) {
	if len(args) != 0 {
		return nil, tender.ErrWrongNumArguments
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"alloc":          &tender.Int{Value: int64(m.Alloc)},
			"total_alloc":    &tender.Int{Value: int64(m.TotalAlloc)},
			"sys":            &tender.Int{Value: int64(m.Sys)},
			"lookups":        &tender.Int{Value: int64(m.Lookups)},
			"mallocs":        &tender.Int{Value: int64(m.Mallocs)},
			"frees":          &tender.Int{Value: int64(m.Frees)},
			"heap_alloc":     &tender.Int{Value: int64(m.HeapAlloc)},
			"heap_sys":       &tender.Int{Value: int64(m.HeapSys)},
			"heap_idle":      &tender.Int{Value: int64(m.HeapIdle)},
			"heap_inuse":     &tender.Int{Value: int64(m.HeapInuse)},
			"heap_released":  &tender.Int{Value: int64(m.HeapReleased)},
			"heap_objects":   &tender.Int{Value: int64(m.HeapObjects)},
			"stack_inuse":    &tender.Int{Value: int64(m.StackInuse)},
			"stack_sys":      &tender.Int{Value: int64(m.StackSys)},
			"mspan_inuse":    &tender.Int{Value: int64(m.MSpanInuse)},
			"mspan_sys":      &tender.Int{Value: int64(m.MSpanSys)},
			"mcache_inuse":   &tender.Int{Value: int64(m.MCacheInuse)},
			"mcache_sys":     &tender.Int{Value: int64(m.MCacheSys)},
			"buck_hash_sys":  &tender.Int{Value: int64(m.BuckHashSys)},
			"gc_sys":         &tender.Int{Value: int64(m.GCSys)},
			"other_sys":      &tender.Int{Value: int64(m.OtherSys)},
			"next_gc":        &tender.Int{Value: int64(m.NextGC)},
			"last_gc":        &tender.Int{Value: int64(m.LastGC)},
			"pause_total_ns": &tender.Int{Value: int64(m.PauseTotalNs)},
			"num_gc":         &tender.Int{Value: int64(m.NumGC)},
			"num_forced_gc":  &tender.Int{Value: int64(m.NumForcedGC)},
			"gccpu_fraction": &tender.Float{Value: m.GCCPUFraction},
		},
	}, nil
}