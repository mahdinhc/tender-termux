package stdlib

import (
	"github.com/2dprototype/tender"
)

func wrapError(err error) tender.Object {
	if err == nil {
		return tender.TrueValue
	}
	return &tender.Error{Value: &tender.String{Value: err.Error()}}
}
