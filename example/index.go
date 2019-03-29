package main

import (
	"fmt"
	"syscall/js"

	"viewer/example/state"
	"viewer/helper"
)

var (
	s state.State
)

func main() {
	helper.Register("content", &s)
	js.Global().Set("UpdateName", js.FuncOf(UpdateName))
	fmt.Println("hello from wasm!")

	select {}
}

func UpdateName(this js.Value, args []js.Value) interface{} {
	s.Name = args[0].String()
	helper.Render()
	return nil
}
