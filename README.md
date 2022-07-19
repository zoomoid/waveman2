# waveman2

![zoomoid - Morgendämmerung waveform](hack/Morgend%C3%A4mmerung.svg)

`waveman2` is the Golang successor the first waveman, built to convert audio files
(later, only mp3), into cool-looking, more abstract audio waveforms, by reducing
the samples to a smaller-sized slice and rendering out SVGs.

See <https://github.com/zoomoid/wave-man> for the original project written in
Python. Be warned, decoding is slow, even though the audio file is resampled to
16 times lower sampling rate before processing, and the project depends on
several utilities to decode mp3 files and convert them to PCM before
tranformation.

`waveman2` is designed to be both extensible with other painters than the default
ones, and also to be imported into any other Golang project as a dependency.
This way, you can easily realize the original idea of waveman to be combined
with a web server that processes audio files sent to the server with a set of
defaults and/or user-defined properties, without having to implement web server
functionality in the waveman codebase itself.

`waveman2` comes with a CLI to use for processing audio files. Its usage is similar to
the older waveman, though several flag names changed:

```bash
$ waveman --help

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

The box painter is used by setting the --box flag. The box color can be set
with --color. The alignment axis can be either "top", "center", or "bottom",
and set with --alignment. --height (or -h) sets the height of highest box, thus also the
height of the entire canvas. --width (or -w) sets the width of each box's bounding box.
--gap sets the space left between each box. Boxes are painted centered inside
their bounding box:

|-------------------------------------------|
|<- 0.5 * gap ->|----BOX----|<- 0.5 * gap ->|
|<----------------- width ----------------->|

Lastly, the --rounded (or -r) parameter controls the rounding of the rectangles.
Notably, rounding requires the boxes to have a minimum height, namely at least
the width of the box, to look aesthetically pleasing. When using --rounded,
each box's height will have its width as a lower bound.

The line painter is used by setting the --line flag. A line's path can be closed
by setting the --closed (or -c) flag. This will close the <path> by appending "Z"
at the end of the data points. The resulting shape can be horizontally mirrored by
setting --inverted (-i). This uses CSS transforms as linear transformation, rather
than computing the data points with offset.

When the path is closed, the color of the enclosed shape can be set with
--fill-color. The color of the line is set with --stroke-color, and the width of
the line with --stroke-width. All those require SVG/CSS-compliant values for the
attributes.

Similarly to the box painter, the --height (or -h) flag controls the shape's overall
height. --spread (or -s) controls the horizontal spacing between each of the sample
points.

To make the line appear smoothly from a discrete set of points, we interpolate
control points for each sample point using cubic hermetic interpolation to fit
cubic polynomials. Namely, we implement 2 interpolation schemes: Fritsch-Carlson
and Steffen. Details can be seen here:
http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
This way, the shape appears smooth. Interpolation can also be controlled with
flags: By default, the Frisch-Carlson scheme is used; setting "--interpolation steffen"
uses the Steffen scheme. If you want to disable interpolation entirely, set
"--interpolation none".

Usage:
  waveman [flags]

Examples:

# Create a black box waveform with 50 blocks for a single mp3
waveman --box --chunks 50 -f audio.mp3

# Create a line waveform with 32 sample points for a single mp3
waveman --line --chunks 32 -f audio.mp3

# Create a red box waveform with 50 blocks at 1/8 downsampling factor
waveman --box --chunks 50 --fill-color red --downsampling-factor 8 -f audio.mp3

# Create a green box waveform for *all* mp3 files in the directory at
# 1/4 downsampling
waveman --box --fill-color green --downsampling-factor 4 -f ./

# Create a closed line waveform with 128 sample points and 1/64 downsampling
# from the start of each chunk, spread apart 50 pixels, with a thicker yellow
# line and flip the shape horizontally
waveman --line --stroke-color yellow --stroke-width 5px \
        --closed --inverted --spread 50 \
        --downsampling-factor 64 --downsampling-mode head \
        -f audio.mp3


Flags:
      --aggregator string          Determines the type of aggregator function to use. Chose one of 'max', 'avg', 'rounded-avg', 'mean-square', or 'root-mean-square' (default "rms")
      --alignment string           Alignment of the shapes, chose one of 'top', 'center', or 'bottom' (default "center")
      --box                        Create a box waveform
  -n, --chunks int                 Chunks are the number of samples in the output of a transformation. For the Box painter, this also means the number of blocks, and for the Line painter, the number of root points of the line (default 64)
  -c, --closed                     Whether the SVG path should be closed or left open
      --color string               Fill color of each box (default "black")
      --downsampling-factor int    Determines the ratio of samples being used for downsampling compared to the full chunk's length. Given in powers of two up two 128 (default 1)
      --downsampling-mode string   Determines the downsampling mode, either by sampling samples from the start, the center, or the end of a chunk
  -f, --file strings               Determines the file to be sampled, can be relative to the current working directory
      --fill-color string          Color for the area enclosed by the line (default "rgba(0 0 0 / 0.5)")
      --gap float                  Gap is the spacing left between each box. Boxes are centered horizonally, so half of gap is subtracted from the box's width (default 5)
  -y, --height float               Height of the shape (default 200)
  -h, --help                       help for waveman
      --interpolation string       Interpolation mechanism to be used for smoothing the curve. Choose one of 'none', 'fritsch-carlson', or 'steffen' (default "fritsch-carlson")
  -i, --inverted                   Whether the shape should be inverted horizontally, i.e., switch the vertical alignment from top to bottom
      --line                       Create a line waveform
  -o, --output string              Writes the output to a given file. If not specified, writes output to stdout
  -r, --recursive                  Searches for all mp3 files in the directory below the specified file
      --rounded float              Rounding factor of each box. Given in pixels. See SVG <rect> rx/ry attributes for details (default 10)
      --stroke-color string        Color of the line's stroke (default "black")
      --stroke-width string        Width of the line's stroke (default "5px")
  -w, --width float                Width of each element (default 10)
```

## Building

You can build the project from source by cloning the repository and then running

```bash
# Downloads all the go dependencies needed to build
$ go get -d -v ./...
# Builds a binary named "waveman" in the local directory
$ go build -o waveman
# Run waveman locally
$ ./waveman
```
