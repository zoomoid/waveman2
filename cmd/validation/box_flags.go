package validation

import (
	"errors"
	"fmt"

	"github.com/zoomoid/waveman2/pkg/painter/box"
)

func ValidateAlignment(alignment string) error {
	a := box.Alignment(alignment)
	switch a {
	case box.AlignmentBottom,
		box.AlignmentCenter,
		box.AlignmentTop,
		box.AlignmentEmpty:
		return nil
	}
	return fmt.Errorf("--alignment %s is not supported", alignment)
}

func ValidateGap(gap float64, width float64) error {
	if gap >= width {
		return errors.New("--gap must not be greater or equal to --width")
	}
	return nil
}
