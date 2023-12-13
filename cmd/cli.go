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

import "github.com/lithammer/dedent"

var (
	WavemanShort string = "waveman generates stylized visual waveforms from mp3 files. Comes with a box painter and a line painter, but can be extended to with other painters easily."

	WavemanLong string = dedent.Dedent(`
		Generate SVG waveforms for one or more mp3 files.

		Prints SVG to stdout when not --output is not specified. When passing in a
		directory, will create SVG files named by the mp3 source files. When the
		--recursive flag is used, *all* mp3 files below the path are used and SVG files
		are colocated with the source mp3 files.
		
		You can configure the sample decoder/transformer in various ways: The number of
		chunks to be passed down to the painter can be set with --chunks (or -n). The
		number must be non-negative. The aggregation function by default uses
		Root-Mean-Square ("rms") for the samples in each chunk. This mimmicks the way
		metering in most DAWs would. Instead, you can also choose "avg", "max",
		"mean-square", or "rounded-avg". The last mode is particularly nice if you don't
		like large floating point numbers in you SVG code, rounding to 3 digits by
		default.
		
		You can improve performance of the waveman by aggressively downsampling the
		audio file. We tested this out and found that using full resolution for the
		aggregation of samples yields minimum visual changes to the audio file, compared
		to the use of high downsampling ratios. We downsample evenly, reducing the
		window of samples per chunks in powers of two. This means that the downsampling
		factor is given as 1/2, 1/4, 1/8, 1/16, etc., up to 1/128.

		To set this factor with flags, use the inverse in --downsampling-factor, e.g.
		"--downsampling-factor 16" for a downsampling factor of 1/16. 
		
		Due to I/O bottlenecks, *this is not done evenly* throughout the file. Instead,
		the downsampling window of a chunk is located either at the start, the middle,
		or the end of a chunk. This behaviour can be set with --downsampling-mode, either
		"head", "center", or "tail".

		--------     ---------------     -----------                  -------
		| File | --> | Transformer | --> | Painter | --> Elements --> | SVG |
		--------     ---------------     -----------                  -------

		The reference implementation brings painters for boxes and lines. Both have 
		multiple configuration options.
	`)

	WavemanExamples string = dedent.Dedent(`
		# Create a black box waveform with 50 blocks for a single mp3
		waveman box --chunks 50 -f audio.mp3

		# Create a line waveform with 32 sample points for a single mp3
		waveman line --chunks 32 -f audio.mp3

		# Create a red box waveform with 50 blocks at 1/8 downsampling factor
		waveman box --chunks 50 --fill-color red --downsampling-factor 8 -f audio.mp3

		# Create a green box waveform for *all* mp3 files in the directory at 
		# 1/4 downsampling
		waveman box --fill-color green --downsampling-factor 4 -f ./

		# Create a closed line waveform with 128 sample points and 1/64 downsampling 
		# from the start of each chunk, spread apart 50 pixels, with a thicker yellow 
		# line and flip the shape horizontally
		waveman line --stroke-color yellow --stroke-width 5px \ 
			--closed --inverted --spread 50 \
			--downsampling-factor 64 --downsampling-mode head \
			-f audio.mp3
	`)
)
