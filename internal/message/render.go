package message

import (
	"bytes"
	"text/template"
)

// Render renders a Go text template with the provided data.
func Render(messageTemplate string, data any) (string, error) {
	tmpl, err := template.New("message").Option("missingkey=error").Parse(messageTemplate)
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", err
	}

	return rendered.String(), nil
}
