package run

import (
	"io"

	"github.com/zoomoid/waveman/v2/pkg/painter"
	"github.com/zoomoid/waveman/v2/pkg/painter/box"
	"github.com/zoomoid/waveman/v2/pkg/painter/line"
	"github.com/zoomoid/waveman/v2/pkg/svg"
	"github.com/zoomoid/waveman/v2/pkg/transform"
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

	svg, err := svg.Template(elements, boxPainter.Width(), boxPainter.Height(), true)
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

	svg, err := svg.Template(elements, linePainter.Width(), linePainter.Height(), true)
	if err != nil {
		return "", err
	}

	return svg.String(), nil
}
