//go:build !wasm && !js

package stdlib

func init() {
	BuiltinModules["os"] = osModule
}
