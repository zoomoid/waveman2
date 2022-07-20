package validation

import (
	"fmt"

	"github.com/zoomoid/waveman2/pkg/painter/line"
)

func ValidateInterpolation(interpolation string) error {
	i := line.Interpolation(interpolation)
	switch i {
	case line.InterpolationFritschCarlson,
		line.InterpolationSteffen,
		line.InterpolationNone,
		line.InterpolationEmpty:
		return nil
	}
	return fmt.Errorf("interpolation %s is not supported", interpolation)
}
