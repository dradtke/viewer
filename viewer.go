package viewer

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"html/template"
)

func Funcs() template.FuncMap {
	return map[string]interface{}{
		"render": render,
	}
}

type templater interface {
	TemplatesAndState() (*template.Template, interface{})
}

type Templater struct {
	T            *template.Template
	InitialState interface{}
}

func (t Templater) TemplatesAndState() (*template.Template, interface{}) {
	return t.T, t.InitialState
}

func render(divID, name string, data templater) (template.HTML, error) {
	t, initialState := data.TemplatesAndState()

	var buf, result bytes.Buffer // TODO: pool?
	if err := gob.NewEncoder(&buf).Encode(initialState); err != nil {
		return template.HTML(""), fmt.Errorf("failed to encode initial state for client template %s: %s", name, err)
	}
	result.WriteString(fmt.Sprintf(`<script type="application/gob" id="%s-state">%s</script>`, divID, base64.StdEncoding.EncodeToString(buf.Bytes())))

	buf.Reset()
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		return template.HTML(""), fmt.Errorf("failed to perform initial server-side render of client template %s: %s", name, err)
	}
	result.WriteString(fmt.Sprintf(`<script type="text/template" id="%s-template">%s</script>`, divID, buf.String()))

	// TODO: can we override the "template" template function? That might allow client templates to reference
	// other client templates without having to pass in more than one name.
	ct, err := template.New("").Delims("[[", "]]").Parse(buf.String())
	if err != nil {
		return template.HTML(""), fmt.Errorf("failed to parse client template %s: %s", name, err)
	}

	buf.Reset()
	if err := ct.Execute(&buf, initialState); err != nil {
		return template.HTML(""), fmt.Errorf("failed to execute client template %s: %s", name, err)
	}
	result.WriteString(fmt.Sprintf(`<div id="%s">%s</div>`, divID, buf.String()))

	return template.HTML(result.String()), nil
}
