//go:build wasm
// +build wasm

package main

import (
	"syscall/js"

	. "github.com/cdvelop/tinystring"
)

func main() {
	// Your WebAssembly code here ok

	// Crear el elemento div
	dom := js.Global().Get("document").Call("createElement", "div")

	buf := Convert().
		Write("<h1>Hello, TinyString</h1>").
		Write("<ul>").
		Write("<li>First item</li>").
		Write("<li>Second item</li>").
		Write("<li>Third item</li>").
		Write("</ul>")

	dom.Set("innerHTML", buf.String())

	// Obtener el body del documento y agregar el elemento
	body := js.Global().Get("document").Get("body")
	body.Call("appendChild", dom)

	logger := func(msg ...any) {
		js.Global().Get("console").Call("log", Translate(msg...).String())
	}

	logger("hello tinystring:", 123, 45.67, true, []string{"a", "b", "c"})

	select {}
}
