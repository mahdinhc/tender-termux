package stdlib

import (
	"github.com/2dprototype/tender/v/colorable"
	"github.com/2dprototype/tender"
)

var colorsModule = map[string]tender.Object{
	"stdout": &tender.UserFunction{
		Value: func(args ...tender.Object) (tender.Object, error) {
			if len(args) != 0 {
				return nil, tender.ErrWrongNumArguments
			}
			return &IOWriter{Value: colorable.NewColorableStdout()}, nil
		},
	},	
	"stderr": &tender.UserFunction{
		Value: func(args ...tender.Object) (tender.Object, error) {
			if len(args) != 0 {
				return nil, tender.ErrWrongNumArguments
			}
			return &IOWriter{Value: colorable.NewColorableStderr()}, nil
		},
	},
}