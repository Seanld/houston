package houston


import (
	"text/template"
	"bytes"
)


func Template(file string, data interface{}) (string, error) {
	// What is the point of having to name templates?
	// Just use their instance to refer to them...
	tp, err := template.ParseFiles(file)
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err = tp.Execute(&rendered, data); err != nil {
		return "", err
	}

	return rendered.String(), nil
}
