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

package box

import (
	"errors"
	"fmt"
)

func validateAlignment(alignment string) error {
	a := Alignment(alignment)
	switch a {
	case AlignmentBottom,
		AlignmentCenter,
		AlignmentTop,
		AlignmentEmpty:
		return nil
	}
	return fmt.Errorf("--alignment %s is not supported", alignment)
}

func validateGap(gap float64, width float64) error {
	if gap >= width {
		return errors.New("--gap must not be greater or equal to --width")
	}
	return nil
}
