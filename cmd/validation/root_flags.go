package validation

import (
	"fmt"

	"github.com/spf13/pflag"
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
