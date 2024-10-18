package main

import (
	"bytes"
	_ "embed"

	"github.com/2dprototype/tender"
	"github.com/2dprototype/tender/stdlib"
)

//go:embed source.tdo
var inputData []byte

func main() {
	
	modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	bytecode := &tender.Bytecode{}
	err := bytecode.Decode(bytes.NewReader(inputData), modules)
	if err != nil {
		return
	}
	machine := tender.NewVM(bytecode, nil, -1)
	err = machine.Run()
	
	if err == nil {
		return
	}
	
}

