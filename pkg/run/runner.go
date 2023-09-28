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

package run

import (
	"io"

	"github.com/zoomoid/waveman2/pkg/painter"
	"github.com/zoomoid/waveman2/pkg/painter/box"
	"github.com/zoomoid/waveman2/pkg/painter/line"
	"github.com/zoomoid/waveman2/pkg/svg"
	"github.com/zoomoid/waveman2/pkg/transform"
)

// Box painter reference runner
func Box(f io.Reader, transformerOptions *transform.ReaderOptions, boxOptions *box.BoxOptions) (string, error) {
	transformer, err := transform.New(transformerOptions, f)
	if err != nil {
		return "", err
	}
	blocks := transformer.Blocks()

	boxPainter := box.New(&painter.PainterOptions{
		Data: blocks,
	}, boxOptions)

	elements := boxPainter.Draw()

	svg, err := svg.Template(elements, true, boxPainter.Viewbox())
	if err != nil {
		return "", err
	}

	return svg.String(), nil
}

// Line painter reference runner
func Line(f io.Reader, transformerOptions *transform.ReaderOptions, lineOptions *line.LineOptions) (string, error) {
	transformer, err := transform.New(transformerOptions, f)
	if err != nil {
		return "", err
	}
	blocks := transformer.Blocks()

	linePainter := line.New(&painter.PainterOptions{
		Data: blocks,
	}, lineOptions)

	elements := linePainter.Draw()

	svg, err := svg.Template(elements, true, linePainter.Viewbox())
	if err != nil {
		return "", err
	}

	return svg.String(), nil
}
