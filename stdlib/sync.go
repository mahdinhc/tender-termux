package stdlib

import (
	"sync"

	"github.com/2dprototype/tender"
)

var syncModule = map[string]tender.Object{
	"mutex": &tender.NativeFunction{
		Name:  "mutex",
		Value: syncNewMutex,
	},
	"rwmutex": &tender.NativeFunction{
		Name:  "rwmutex",
		Value: syncNewRWMutex,
	},
	"wait_group": &tender.NativeFunction{
		Name:  "wait_group",
		Value: syncNewWaitGroup,
	},
	"once": &tender.NativeFunction{
		Name:  "once",
		Value: syncNewOnce,
	},
	"cond": &tender.NativeFunction{
		Name:  "cond",
		Value: syncNewCond,
	},
	"map": &tender.NativeFunction{
		Name:  "map",
		Value: syncNewMap,
	},
}

func syncNewMutex(args ...tender.Object) (tender.Object, error) {
	if len(args) != 0 {
		return nil, tender.ErrWrongNumArguments
	}
	return makeMutexObj(&sync.Mutex{}), nil
}

func makeMutexObj(mu *sync.Mutex) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"lock": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					mu.Lock()
					return tender.NullValue, nil
				},
			},
			"unlock": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					mu.Unlock()
					return tender.NullValue, nil
				},
			},
			"try_lock": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					return tender.FromBool(mu.TryLock()), nil
				},
			},
		},
	}
}

func syncNewRWMutex(args ...tender.Object) (tender.Object, error) {
	if len(args) != 0 {
		return nil, tender.ErrWrongNumArguments
	}
	return makeRWMutexObj(&sync.RWMutex{}), nil
}

func makeRWMutexObj(rw *sync.RWMutex) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"lock": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					rw.Lock()
					return tender.NullValue, nil
				},
			},
			"unlock": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					rw.Unlock()
					return tender.NullValue, nil
				},
			},
			"rlock": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					rw.RLock()
					return tender.NullValue, nil
				},
			},
			"runlock": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					rw.RUnlock()
					return tender.NullValue, nil
				},
			},
			"try_lock": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					return tender.FromBool(rw.TryLock()), nil
				},
			},
			"try_rlock": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					return tender.FromBool(rw.TryRLock()), nil
				},
			},
		},
	}
}

func syncNewWaitGroup(args ...tender.Object) (tender.Object, error) {
	if len(args) != 0 {
		return nil, tender.ErrWrongNumArguments
	}
	return makeWaitGroupObj(&sync.WaitGroup{}), nil
}

func makeWaitGroupObj(wg *sync.WaitGroup) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"add": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					delta, ok := tender.ToInt(args[0])
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "int(compatible)",
							Found:    args[0].TypeName(),
						}
					}
					wg.Add(delta)
					return tender.NullValue, nil
				},
			},
			"done": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					wg.Done()
					return tender.NullValue, nil
				},
			},
			"wait": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					wg.Wait()
					return tender.NullValue, nil
				},
			},
		},
	}
}

func syncNewOnce(args ...tender.Object) (tender.Object, error) {
	if len(args) != 0 {
		return nil, tender.ErrWrongNumArguments
	}
	return makeOnceObj(&sync.Once{}), nil
}

func makeOnceObj(once *sync.Once) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"do": &tender.NativeFunction{
				NeedVMObj: true,
				Value: func(args ...tender.Object) (tender.Object, error) {
					vmObj := args[0].(*tender.VMObj).Value
					args = args[1:]
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					fn := args[0]
					if !fn.CanCall() {
						return nil, tender.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "callable function",
							Found:    fn.TypeName(),
						}
					}
					var runErr error
					once.Do(func() {
						_, runErr = tender.WrapFuncCall(vmObj, fn)
					})
					if runErr != nil {
						return nil, runErr
					}
					return tender.NullValue, nil
				},
			},
		},
	}
}

type mapLocker struct {
	m *tender.ImmutableMap
}

func (l *mapLocker) Lock() {
	if lockFn, ok := l.m.Value["lock"]; ok {
		_, _ = lockFn.Call()
	}
}

func (l *mapLocker) Unlock() {
	if unlockFn, ok := l.m.Value["unlock"]; ok {
		_, _ = unlockFn.Call()
	}
}

func syncNewCond(args ...tender.Object) (tender.Object, error) {
	var locker sync.Locker
	if len(args) == 0 {
		locker = &sync.Mutex{}
	} else if len(args) == 1 {
		if m, ok := args[0].(*tender.ImmutableMap); ok {
			locker = &mapLocker{m: m}
		} else {
			return nil, tender.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "mutex or rwmutex",
				Found:    args[0].TypeName(),
			}
		}
	} else {
		return nil, tender.ErrWrongNumArguments
	}

	cond := sync.NewCond(locker)
	return makeCondObj(cond), nil
}

func makeCondObj(cond *sync.Cond) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"wait": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					cond.Wait()
					return tender.NullValue, nil
				},
			},
			"signal": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					cond.Signal()
					return tender.NullValue, nil
				},
			},
			"broadcast": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					cond.Broadcast()
					return tender.NullValue, nil
				},
			},
			"lock": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					cond.L.Lock()
					return tender.NullValue, nil
				},
			},
			"unlock": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					cond.L.Unlock()
					return tender.NullValue, nil
				},
			},
		},
	}
}

func syncNewMap(args ...tender.Object) (tender.Object, error) {
	if len(args) != 0 {
		return nil, tender.ErrWrongNumArguments
	}
	return makeMapObj(&sync.Map{}), nil
}

func makeMapObj(m *sync.Map) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"store": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 2 {
						return nil, tender.ErrWrongNumArguments
					}
					m.Store(args[0], args[1])
					return tender.NullValue, nil
				},
			},
			"load": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					val, loaded := m.Load(args[0])
					resVal, ok := val.(tender.Object)
					if !ok || !loaded {
						resVal = tender.NullValue
					}
					return &tender.Array{Value: []tender.Object{resVal, tender.FromBool(loaded)}}, nil
				},
			},
			"load_or_store": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 2 {
						return nil, tender.ErrWrongNumArguments
					}
					actual, loaded := m.LoadOrStore(args[0], args[1])
					resVal, ok := actual.(tender.Object)
					if !ok {
						resVal = tender.NullValue
					}
					return &tender.Array{Value: []tender.Object{resVal, tender.FromBool(loaded)}}, nil
				},
			},
			"load_and_delete": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					value, loaded := m.LoadAndDelete(args[0])
					resVal, ok := value.(tender.Object)
					if !ok || !loaded {
						resVal = tender.NullValue
					}
					return &tender.Array{Value: []tender.Object{resVal, tender.FromBool(loaded)}}, nil
				},
			},
			"delete": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					m.Delete(args[0])
					return tender.NullValue, nil
				},
			},
			"swap": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 2 {
						return nil, tender.ErrWrongNumArguments
					}
					previous, loaded := m.Swap(args[0], args[1])
					resVal, ok := previous.(tender.Object)
					if !ok || !loaded {
						resVal = tender.NullValue
					}
					return &tender.Array{Value: []tender.Object{resVal, tender.FromBool(loaded)}}, nil
				},
			},
			"compare_and_swap": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 3 {
						return nil, tender.ErrWrongNumArguments
					}
					swapped := m.CompareAndSwap(args[0], args[1], args[2])
					return tender.FromBool(swapped), nil
				},
			},
			"compare_and_delete": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 2 {
						return nil, tender.ErrWrongNumArguments
					}
					deleted := m.CompareAndDelete(args[0], args[1])
					return tender.FromBool(deleted), nil
				},
			},
			"range": &tender.NativeFunction{
				NeedVMObj: true,
				Value: func(args ...tender.Object) (tender.Object, error) {
					vmObj := args[0].(*tender.VMObj).Value
					args = args[1:]
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					fn := args[0]
					if !fn.CanCall() {
						return nil, tender.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "callable function",
							Found:    fn.TypeName(),
						}
					}
					var rangeErr error
					m.Range(func(key, value any) bool {
						kObj, ok := key.(tender.Object)
						if !ok {
							kObj = tender.NullValue
						}
						vObj, ok := value.(tender.Object)
						if !ok {
							vObj = tender.NullValue
						}
						retVal, err := tender.WrapFuncCall(vmObj, fn, kObj, vObj)
						if err != nil {
							rangeErr = err
							return false
						}
						if retVal != nil && retVal.IsFalsy() {
							return false
						}
						return true
					})
					if rangeErr != nil {
						return nil, rangeErr
					}
					return tender.NullValue, nil
				},
			},
		},
	}
}