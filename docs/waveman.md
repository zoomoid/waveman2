## waveman

waveman generates stylized visual waveforms from mp3 files. Comes with a box painter and a line painter, but can be extended to with other painters easily.

### Synopsis


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


### Examples

```

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

```

### Options

```
      --aggregator string          Determines the type of aggregator function to use. Chose one of 'max', 'avg', 'rounded-avg', 'mean-square', or 'root-mean-square' (default "rms")
  -n, --chunks int                 Chunks are the number of samples in the output of a transformation. For the Box painter, this also means the number of blocks, and for the Line painter, the number of root points of the line (default 64)
      --clamp-high float           Upper clipping of samples (default 1)
      --clamp-low float            Lower clipping of samples
      --downsampling-factor int    Determines the ratio of samples being used for downsampling compared to the full chunk's length. Given in powers of two up two 128 (default 1)
      --downsampling-mode string   Determines the downsampling mode, either by sampling samples from the start, the center, or the end of a chunk
  -f, --file strings               Determines the file to be sampled, can be relative to the current working directory
  -y, --height float               Height of the shape (default 200)
  -h, --help                       help for waveman
      --normalize                  Whether or not to normalize samples to [0,1]. When running in batch mode, this loses overall levels information, as each track is normalized individually
  -o, --output string              Writes the output to a given file. If not specified, writes output to stdout
  -r, --recursive                  Searches for all mp3 files in the directory below the specified file
  -w, --width float                Width of each element (default 10)
      --window string              Window algorithm. Defaults to rectangular, which is equivalent to no windowing. Can be used with other windowing algorithms to filter high sample values at the start and end of tracks. (default "rectangular")
      --window-p float             Window algorithm parameter. For most algorithms, this determines the steepness of the slope of the window
```

### SEE ALSO

* [waveman box](waveman_box.md)	 - 
* [waveman completion](waveman_completion.md)	 - Generate completion script
* [waveman line](waveman_line.md)	 - 

###### Auto generated by spf13/cobra on 13-Dec-2023
