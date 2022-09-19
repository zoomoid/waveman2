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

package options

const (
	DownsamplingMode   string = "downsampling-mode"
	DownsamplingFactor string = "downsampling-factor"
	Aggregator         string = "aggregator"
	Chunks             string = "chunks"
	ChunksShort        string = "n"
)

const (
	DownsamplingModeDescription   string = "Determines the downsampling mode, either by sampling samples from the start, the center, or the end of a chunk"
	DownsamplingFactorDescription string = "Determines the ratio of samples being used for downsampling compared to the full chunk's length. Given in powers of two up two 128"
	AggregatorDescription         string = "Determines the type of aggregator function to use. Chose one of 'max', 'avg', 'rounded-avg', 'mean-square', or 'root-mean-square'"
	ChunksDescription             string = "Chunks are the number of samples in the output of a transformation. For the Box painter, this also means the number of blocks, and for the Line painter, the number of root points of the line"
)
