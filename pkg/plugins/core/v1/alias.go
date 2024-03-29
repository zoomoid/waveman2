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

package v1

import (
	"github.com/zoomoid/waveman2/pkg/plugins/core/v1/box"
	"github.com/zoomoid/waveman2/pkg/plugins/core/v1/line"
	"github.com/zoomoid/waveman2/pkg/plugins/core/v1/sweep"
	"github.com/zoomoid/waveman2/pkg/plugins/core/v1/wave"
)

var Line = line.Plugin

var NewLinePainter = line.NewPainter

var Box = box.Plugin

var NewBoxPainter = box.NewPainter

var Wave = wave.Plugin

var NewWavePainter = wave.NewPainter

var Sweep = sweep.Plugin

var NewSweepPainter = sweep.NewPainter
