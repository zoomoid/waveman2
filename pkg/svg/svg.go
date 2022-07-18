package svg

import (
	"bytes"
	"text/template"

	"github.com/lithammer/dedent"
	"github.com/yosssi/gohtml"
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
    {{ range $_, $el := .Elements -}}
    {{ $el }}
    {{ end -}}
  </svg>
`)

type TemplateBindings struct {
	PreserveAspectRatio bool
	Width               float64
	Height              float64
	Elements            []string
}

// Template executes the default SVG template and writes all previously created SVG elements to the body
// Returns the template as string
// Returns an error if any failures occur.
func Template(elements []string, elWidth float64, elHeight float64, preserveAspectRatio bool) (*bytes.Buffer, error) {
	tmpl, err := template.New("svg").Parse(DefaultSvgTemplate)
	if err != nil {
		return nil, err
	}

	bindings := &TemplateBindings{
		PreserveAspectRatio: preserveAspectRatio,
		Width:               elWidth,
		Height:              elHeight,
		Elements:            elements,
	}
	buffer := &bytes.Buffer{}
	err = tmpl.Execute(gohtml.NewWriter(buffer), bindings)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
