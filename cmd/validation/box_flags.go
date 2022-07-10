package validation

import (
	"errors"
	"fmt"

	"github.com/zoomoid/waveman/v2/pkg/paint/box"
)

func ValidateAlignment(alignment string) error {
	a := box.Alignment(alignment)
	switch a {
	case box.AlignmentBottom:
	case box.AlignmentCenter:
	case box.AlignmentTop:
	case box.AlignmentEmpty:
		return nil
	}
	return fmt.Errorf("--alignment %s is not supported", alignment)
}

func ValidateBoxHeight(height float64) error {
	if height >= 0 {
		return nil
	}
	return errors.New("--height must be non-negative")
}

func ValidateBoxWidth(width float64) error {
	if width >= 0 {
		return nil
	}
	return errors.New("--width must be non-negative")
}

func ValidateGap(gap float64, width float64) error {
	if gap >= width {
		return errors.New("--gap must not be greater or equal to --width")
	}
	return nil
}
