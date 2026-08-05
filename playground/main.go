//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"

	"github.com/2dprototype/tender"
	"github.com/2dprototype/tender/parser"
	"github.com/2dprototype/tender/stdlib"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("tenderRun", js.FuncOf(tenderRun))
	js.Global().Set("tenderParse", js.FuncOf(tenderParse))
	js.Global().Set("tenderGetModules", js.FuncOf(tenderGetModules))

	fmt.Println("[Tender WASM Engine] Ready")
	<-c
}

func tenderRun(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return errorJSON("No source code provided")
	}
	src := args[0].String()
	startTime := time.Now()

	script := tender.NewScript([]byte(src))
	script.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))

	compiled, err := script.Run()
	durationMs := time.Since(startTime).Milliseconds()

	if err != nil {
		resp, _ := json.Marshal(map[string]any{
			"error":      err.Error(),
			"durationMs": durationMs,
		})
		return string(resp)
	}

	var varsMap = make(map[string]string)
	if compiled != nil {
		allVars := compiled.GetAll()
		for _, v := range allVars {
			if v != nil && v.Name() != "" && !v.IsNull() && v.Object() != nil {
				varsMap[v.Name()] = v.Object().String()
			}
		}
	}

	resp, _ := json.Marshal(map[string]any{
		"success":    true,
		"durationMs": durationMs,
		"variables":  varsMap,
	})
	return string(resp)
}

func tenderParse(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return errorJSON("No source code provided")
	}
	src := args[0].String()

	fileSet := parser.NewFileSet()
	srcFile := fileSet.AddFile("playground.td", -1, len(src))

	p := parser.NewParser(srcFile, []byte(src), nil)
	file, err := p.ParseFile()
	if err != nil {
		return errorJSON(err.Error())
	}

	jsonBytes, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return errorJSON(err.Error())
	}

	resp, _ := json.Marshal(map[string]any{
		"ast": string(jsonBytes),
	})
	return string(resp)
}

func tenderGetModules(this js.Value, args []js.Value) any {
	names := stdlib.AllModuleNames()
	jsonBytes, _ := json.Marshal(names)
	return string(jsonBytes)
}

func errorJSON(msg string) string {
	resp, _ := json.Marshal(map[string]any{"error": msg})
	return string(resp)
}
