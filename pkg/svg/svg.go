/*
Copyright 2022-2023 zoomoid.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
    viewBox="{{ .Viewbox }}"
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
	PreserveAspectRatio string
	Elements            []string
	Viewbox             string
}

// Template executes the default SVG template and writes all previously created SVG elements to the body
// Returns the template as string
// Returns an error if any failures occur.
func Template(elements []string, preserveAspectRatio bool, viewBox string) (*bytes.Buffer, error) {
	tmpl, err := template.New("svg").Parse(DefaultSvgTemplate)
	if err != nil {
		return nil, err
	}

	preservanceKey := "keep"
	if !preserveAspectRatio {
		preservanceKey = "none"
	}

	bindings := &TemplateBindings{
		PreserveAspectRatio: preservanceKey,
		Elements:            elements,
		Viewbox:             viewBox,
	}
	rawBuffer := &bytes.Buffer{}
	err = tmpl.Execute(rawBuffer, bindings)
	if err != nil {
		return nil, err
	}

	outBuf := bytes.NewBuffer(gohtml.FormatBytes(rawBuffer.Bytes()))

	return outBuf, nil
}
