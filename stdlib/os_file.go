package stdlib

import (
	"os"

	"github.com/2dprototype/tender"
)

func makeOSFile(file *os.File) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			// "writer": &IOWriter{Value: file},
			// "reader": &IOReader{Value: file},
			"writer": &tender.UserFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					return &IOWriter{Value: file}, nil
				},
			},	
			"reader": &tender.UserFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					return &IOWriter{Value:file}, nil
				},
			},
			"chdir": &tender.UserFunction{
				Value: FuncARE(file.Chdir),
			},
			"chown": &tender.UserFunction{
				Value: FuncAIIRE(file.Chown),
			},
			"close": &tender.UserFunction{
				Value: FuncARE(file.Close),
			},
			"name": &tender.UserFunction{
				Value: FuncARS(file.Name),
			},
			"readdirnames": &tender.UserFunction{
				Value: FuncAIRSsE(file.Readdirnames),
			}, 
			"sync": &tender.UserFunction{
				Value: FuncARE(file.Sync),
			}, 
			"write": &tender.UserFunction{
				Value: FuncAYRIE(file.Write),
			}, 
			"write_string": &tender.UserFunction{
				Value: FuncASRIE(file.WriteString),
			}, 
			"read": &tender.UserFunction{
				Value: FuncAYRIE(file.Read),
			}, 	
			"read_at": &tender.UserFunction{
				Value: FuncAYI64RIE(file.ReadAt),
			}, 	
			// "truncate": &tender.UserFunction{
				// Name:  "truncate",
				// Value: FuncAI64RE(file.Truncate),
			// }, 
			"chmod": &tender.UserFunction{
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
					return wrapError(file.Chmod(os.FileMode(i1))), nil
				},
			},
			// seek(offset int, whence int) => int/error
			"seek": &tender.UserFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 2 {
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
					i2, ok := tender.ToInt(args[1])
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     "second",
							Expected: "int(compatible)",
							Found:    args[1].TypeName(),
						}
					}
					res, err := file.Seek(i1, i2)
					if err != nil {
						return wrapError(err), nil
					}
					return &tender.Int{Value: res}, nil
				},
			},
			// stat() => imap(fileinfo)/error
			"stat": &tender.UserFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					return osStat(&tender.String{Value: file.Name()})
				},
			},
		},
	}
}
