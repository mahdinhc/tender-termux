package stdlib

import (
	"encoding/hex"
	"github.com/2dprototype/tender"
)

var hexModule = map[string]tender.Object{
	"encode": &tender.UserFunction{Value: FuncAYRS(hex.EncodeToString)},
	"decode": &tender.UserFunction{Value: FuncASRYE(hex.DecodeString)},
	"dump": &tender.UserFunction{Value: FuncAYRS(hex.Dump)},
}
