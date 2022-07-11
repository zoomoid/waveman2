package run

import (
	"github.com/zoomoid/waveman/v2/pkg/painter"
	"github.com/zoomoid/waveman/v2/pkg/painter/box"
	"github.com/zoomoid/waveman/v2/pkg/painter/line"
	"github.com/zoomoid/waveman/v2/pkg/svg"
	"github.com/zoomoid/waveman/v2/pkg/transform"
)

// Box painter reference runner
func Box(transformerOptions *transform.ReaderOptions, boxOptions *box.BoxOptions) (string, error) {

	transformer, err := transform.New(transformerOptions)
	if err != nil {
		return "", err
	}
	blocks := transformer.Blocks()

	boxPainter := box.New(&painter.PainterOptions{
		Data: blocks,
	}, boxOptions)

	elements := boxPainter.Draw()

	svg, err := svg.Template(elements, boxPainter.TotalWidth(), boxPainter.TotalHeight(), true)
	if err != nil {
		return "", err
	}

	return svg, nil
}

// Line painter reference runner
func Line(transformerOptions *transform.ReaderOptions, lineOptions *line.LineOptions) (string, error) {

	transformer, err := transform.New(transformerOptions)
	if err != nil {
		return "", err
	}
	blocks := transformer.Blocks()

	linePainter := line.New(&painter.PainterOptions{
		Data: blocks,
	}, lineOptions)

	elements := linePainter.Draw()

	svg, err := svg.Template(elements, linePainter.TotalWidth(), linePainter.TotalHeight(), true)
	if err != nil {
		return "", err
	}

	return svg, nil
}
