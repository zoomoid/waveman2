package validation

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman/v2/cmd/options"
)

func ValidatePainterModes(flags *pflag.FlagSet, modes []string) error {
	foundFlag := ""
	for _, mode := range modes {
		m, err := flags.GetBool(mode)
		if err != nil {
			return err
		}
		if m {
			if foundFlag != "" {
				return fmt.Errorf("--%s and --%s are mutually exclusive", foundFlag, mode)
			} else {
				foundFlag = mode
			}
		}
	}
	if foundFlag == "" {
		return fmt.Errorf("painter mode must be specified by its flag")
	}
	return nil
}

func ValidateFilename(filename string) error {
	if filename == "" {
		return errors.New("filename needs to be specified")
	}
	return nil
}

func ValidateHeight(height float64) error {
	if height >= 0 {
		return nil
	}
	return errors.New("--height must be non-negative")
}

func ValidateWidth(width float64) error {
	if width >= 0 {
		return nil
	}
	return errors.New("--width must be non-negative")
}

func ValidateOutput(output string) error {
	o := options.OutputType(output)
	switch o {
	case options.OutputTypeFile, options.OutputTypeEmpty:
		return nil
	}
	return fmt.Errorf("--output does not support type %s, only supported types are %v", output, options.SupportedOutputs)
}
