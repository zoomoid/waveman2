/*
Copyright 2022 zoomoid.

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

package main

import (
	"log"

	"github.com/spf13/cobra/doc"
	"github.com/zoomoid/waveman2/cmd"
	r "github.com/zoomoid/waveman2/pkg/reference"
	"github.com/zoomoid/waveman2/pkg/streams"
)

func main() {
	waveman := cmd.NewWaveman(nil, streams.DefaultStreams).
		Plugin(r.BoxPainterPlugin).
		Plugin(r.LinePainterPlugin).
		Complete()
	err := doc.GenMarkdownTree(waveman, "./")
	if err != nil {
		log.Fatal(err)
	}
}
