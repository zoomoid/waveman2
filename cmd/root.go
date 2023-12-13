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

package cmd

import (
	"github.com/rs/zerolog/log"
	corev1 "github.com/zoomoid/waveman2/pkg/plugins/core/v1"
	"github.com/zoomoid/waveman2/pkg/streams"
)

func Execute(version string) {

	rootCmd := NewWaveman(nil, streams.DefaultStreams).
		V(version).
		Plugin(corev1.Box).
		Plugin(corev1.Line).
		Complete()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}
