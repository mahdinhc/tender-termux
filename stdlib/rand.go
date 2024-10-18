package stdlib

import (
	"math/rand"

	"github.com/2dprototype/tender"
)

var randModule = map[string]tender.Object{
	"int": &tender.UserFunction{
		Name:  "int",
		Value: FuncARI64(rand.Int63),
	},
	"float": &tender.UserFunction{
		Name:  "float",
		Value: FuncARF(rand.Float64),
	},
	"intn": &tender.UserFunction{
		Name:  "intn",
		Value: FuncAI64RI64(rand.Int63n),
	},
	"exp_float": &tender.UserFunction{
		Name:  "exp_float",
		Value: FuncARF(rand.ExpFloat64),
	},
	"norm_float": &tender.UserFunction{
		Name:  "norm_float",
		Value: FuncARF(rand.NormFloat64),
	},
	"perm": &tender.UserFunction{
		Name:  "perm",
		Value: FuncAIRIs(rand.Perm),
	},
	"seed": &tender.UserFunction{
		Name:  "seed",
		Value: FuncAI64R(rand.Seed),
	},
	"read": &tender.UserFunction{
		Name: "read",
		Value: func(args ...tender.Object) (ret tender.Object, err error) {
			if len(args) != 1 {
				return nil, tender.ErrWrongNumArguments
			}
			y1, ok := args[0].(*tender.Bytes)
			if !ok {
				return nil, tender.ErrInvalidArgumentType{
					Name:     "first",
					Expected: "bytes",
					Found:    args[0].TypeName(),
				}
			}
			res, err := rand.Read(y1.Value)
			if err != nil {
				ret = wrapError(err)
				return
			}
			return &tender.Int{Value: int64(res)}, nil
		},
	},
	"rand": &tender.UserFunction{
		Name: "rand",
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
			src := rand.NewSource(i1)
			return randRand(rand.New(src)), nil
		},
	},
}

func randRand(r *rand.Rand) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"int": &tender.UserFunction{
				Name:  "int",
				Value: FuncARI64(r.Int63),
			},
			"float": &tender.UserFunction{
				Name:  "float",
				Value: FuncARF(r.Float64),
			},
			"intn": &tender.UserFunction{
				Name:  "intn",
				Value: FuncAI64RI64(r.Int63n),
			},
			"exp_float": &tender.UserFunction{
				Name:  "exp_float",
				Value: FuncARF(r.ExpFloat64),
			},
			"norm_float": &tender.UserFunction{
				Name:  "norm_float",
				Value: FuncARF(r.NormFloat64),
			},
			"perm": &tender.UserFunction{
				Name:  "perm",
				Value: FuncAIRIs(r.Perm),
			},
			"seed": &tender.UserFunction{
				Name:  "seed",
				Value: FuncAI64R(r.Seed),
			},
			"read": &tender.UserFunction{
				Name: "read",
				Value: func(args ...tender.Object) (
					ret tender.Object,
					err error,
				) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					y1, ok := args[0].(*tender.Bytes)
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "bytes",
							Found:    args[0].TypeName(),
						}
					}
					res, err := r.Read(y1.Value)
					if err != nil {
						ret = wrapError(err)
						return
					}
					return &tender.Int{Value: int64(res)}, nil
				},
			},
		},
	}
}
