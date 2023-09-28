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

package validation

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman2/cmd/options"
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

// func ValidateFilenames(filenames []string, output string) error {
// 	if options.OutputType(output) == options.OutputTypeEmpty && len(filenames) > 1 {
// 		return fmt.Errorf("cannot use multiple files with stdout target, use --output file")
// 	}
// 	return nil
// }
