package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"html/template"
	"syscall/js"
)

var registry = make(map[string]info)

type info struct {
	node  js.Value
	state interface{}
	t     *template.Template
}

func getElement(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}

func init() {
	resolve(
		js.Global().Call("fetch", "https://raw.githubusercontent.com/wingify/dom-comparator/master/dist/dom-comparator.js"),
		func(response js.Value) {
			resolve(
				response.Call("text"),
				func(response js.Value) { fmt.Println("got value:", response) },
				func(jerr js.Value) { fmt.Println("error:", jerr) },
			)
		},
		func(jerr js.Value) { fmt.Println("error:", jerr) },
	)
}

func resolve(promise js.Value, then, catch func(js.Value)) {
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		then(args[0])
		return nil
	}))
	promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		catch(args[0])
		return nil
	}))
}

func Register(divID string, state interface{}) {
	node := getElement(divID)
	templateNode := getElement(divID + "-template")
	stateNode := getElement(divID + "-state")

	t, err := template.New("").Delims("[[", "]]").Parse(templateNode.Get("text").String())
	if err != nil {
		panic(err)
	}

	stateBytes, err := base64.StdEncoding.DecodeString(stateNode.Get("text").String())
	if err != nil {
		panic(err)
	}
	if gob.NewDecoder(bytes.NewReader(stateBytes)).Decode(state); err != nil {
		panic(err)
	}

	registry[divID] = info{
		node:  node,
		state: state,
		t:     t,
	}
}

func Render() error {
	for name, v := range registry {
		var buf bytes.Buffer
		if err := v.t.Execute(&buf, v.state); err != nil {
			return fmt.Errorf("failed to render template for %s: %s", name, err)
		}
		// TODO: do a smarter DOM diff
		v.node.Set("innerHTML", buf.String())
	}
	return nil
}
