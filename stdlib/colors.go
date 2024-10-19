package stdlib

import (
	"github.com/2dprototype/tender/v/colorable"
	"github.com/2dprototype/tender"
)

var colorsModule = map[string]tender.Object{
	"stdout": &IOWriter{Value: colorable.NewColorableStdout()},
	"stderr": &IOWriter{Value: colorable.NewColorableStderr()},
}