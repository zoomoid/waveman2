package svg

import (
	"bytes"
	"text/template"

	"github.com/lithammer/dedent"
)

// DefaultSvgTemplate contains the Golang template to be executed with elements
var DefaultSvgTemplate = dedent.Dedent(`
	<svg 
		baseProfile="tiny"
		preserveAspectRatio="{{ .PreserveAspectRatio }}"
		version="1.2"
		viewBox="0 0 {{ .Width }} {{ .Height }}"
		height="100%" width="100%"
		xmlns="http://www.w3.org/2000/svg"
		xmlns:ev="http://www.w3.org/2001/xml-events"
		xmlns:xlink="http://www.w3.org/1999/xlink"
	>
	{{- range $el, _ := .Elements -}}
		{{ $el }}
	{{- end -}}
	</svg>
`)

type TemplateBindings struct {
	PreserveAspectRatio bool
	Width               float32
	Height              float32
	Elements            []string
}

// Template executes the default SVG template and writes all previously created SVG elements to the body
// Returns the template as string
// Returns an error if any failures occur.
func Template(elements []string, elWidth float32, elHeight float32, preserveAspectRatio bool) (string, error) {
	tmpl, err := template.New("svg").Parse(DefaultSvgTemplate)
	if err != nil {
		return "", err
	}

	bindings := &TemplateBindings{
		PreserveAspectRatio: preserveAspectRatio,
		Width:               elWidth * float32(len(elements)),
		Height:              elHeight,
		Elements:            elements,
	}
	buffer := &bytes.Buffer{}
	err = tmpl.Execute(buffer, bindings)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
