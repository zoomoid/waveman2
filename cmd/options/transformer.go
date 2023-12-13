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

package options

const (
	DownsamplingMode   string = "downsampling-mode"
	DownsamplingFactor string = "downsampling-factor"
	Aggregator         string = "aggregator"
	Chunks             string = "chunks"
	ChunksShort        string = "n"

	Normalize string = "normalize"

	ClampLow  string = "clamp-low"
	ClampHigh string = "clamp-high"

	WindowP         string = "window-p"
	WindowAlgorithm string = "window"
)

const (
	DownsamplingModeDescription   string = "Determines the downsampling mode, either by sampling samples from the start, the center, or the end of a chunk"
	DownsamplingFactorDescription string = "Determines the ratio of samples being used for downsampling compared to the full chunk's length. Given in powers of two up two 128"
	AggregatorDescription         string = "Determines the type of aggregator function to use. Chose one of 'max', 'avg', 'rounded-avg', 'mean-square', or 'root-mean-square'"
	ChunksDescription             string = "Chunks are the number of samples in the output of a transformation. For the Box painter, this also means the number of blocks, and for the Line painter, the number of root points of the line"

	NormalizeDescription string = "Whether or not to normalize samples to [0,1]. When running in batch mode, this loses overall levels information, as each track is normalized individually"

	ClampLowDescription  string = "Lower clipping of samples"
	ClampHighDescription string = "Upper clipping of samples"

	WindowPDescription         string = "Window algorithm parameter. For most algorithms, this determines the steepness of the slope of the window"
	WindowAlgorithmDescription string = "Window algorithm. Defaults to rectangular, which is equivalent to no windowing. Can be used with other windowing algorithms to filter high sample values at the start and end of tracks."
)
