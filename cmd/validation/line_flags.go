package validation

import (
	"errors"
	"fmt"

	"github.com/zoomoid/waveman/v2/pkg/paint/line"
)

func ValidateInterpolation(interpolation string) error {
	i := line.Interpolation(interpolation)
	switch i {
	case line.InterpolationFritschCarlson:
	case line.InterpolationSteffen:
	case line.InterpolationNone:
	case line.InterpolationEmpty:
		return nil
	}
	return fmt.Errorf("interpolation %s is not supported", interpolation)
}

func ValidateSpread(spread float64) error {
	if spread >= 0 {
		return nil
	}
	return errors.New("--spread must be non-negative")
}

func ValidateLineHeight(height float64) error {
	if height >= 0 {
		return nil
	}
	return errors.New("--height must be non-negative")
}
