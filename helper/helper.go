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
