## waveman line



### Synopsis


The line painter draws a line that traces each data point.

A line's path can be closed by setting the --closed (or -c) flag.
This will close the <path> by appending "Z" at the end of the data points.

When the path is closed, the color of the enclosed shape can be set with 
--fill-color.

The color of the line is set with --stroke-color, and the width of the line 
with --stroke-width. All those require SVG/CSS-compliant values for the
attributes.

The shape can be horizontally mirrored by setting --inverted (-i).

To create a symmetric shape, similar to Box with alignment = center, but with a 
continuous line, use --mirrored.

Similarly to the box painter, the --height (or -h) flag controls the shape's overall
height. 

--spread (or -s) controls the horizontal spacing between each of the sample
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


```
waveman line [flags]
```

### Options

```
  -c, --closed                 Whether the SVG path should be closed or left open
      --fill-color string      Color for the area enclosed by the line (default "rgba(0 0 0 / 0.5)")
  -h, --help                   help for line
      --interpolation string   Interpolation mechanism to be used for smoothing the curve [none,fritsch-carlson,steffen] (default "fritsch-carlson")
  -i, --inverted               Whether the shape should be inverted horizontally, i.e., switch the vertical alignment from top to bottom
      --mirrored               Whether the shape should be mirrored, creating a symmetric waveform. Note that this is not *realistic*, i.e. it does not represent positive and negative polarity of samples!
      --stroke-color string    Color of the line's stroke (default "none")
      --stroke-width float     Width of the line's stroke
```

### Options inherited from parent commands

```
      --aggregator string          Determines the type of aggregator function to use. Chose one of 'max', 'avg', 'rounded-avg', 'mean-square', or 'root-mean-square' (default "rms")
  -n, --chunks int                 Chunks are the number of samples in the output of a transformation. For the Box painter, this also means the number of blocks, and for the Line painter, the number of root points of the line (default 64)
      --clamp-high float           Upper clipping of samples (default 1)
      --clamp-low float            Lower clipping of samples
      --downsampling-factor int    Determines the ratio of samples being used for downsampling compared to the full chunk's length. Given in powers of two up two 128 (default 1)
      --downsampling-mode string   Determines the downsampling mode, either by sampling samples from the start, the center, or the end of a chunk
  -f, --file strings               Determines the file to be sampled, can be relative to the current working directory
  -y, --height float               Height of the shape (default 200)
      --normalize                  Whether or not to normalize samples to [0,1]. When running in batch mode, this loses overall levels information, as each track is normalized individually
  -o, --output string              Writes the output to a given file. If not specified, writes output to stdout
  -r, --recursive                  Searches for all mp3 files in the directory below the specified file
  -w, --width float                Width of each element (default 10)
      --window string              Window algorithm. Defaults to rectangular, which is equivalent to no windowing. Can be used with other windowing algorithms to filter high sample values at the start and end of tracks. (default "rectangular")
      --window-p float             Window algorithm parameter. For most algorithms, this determines the steepness of the slope of the window
```

### SEE ALSO

* [waveman](waveman.md)	 - waveman generates stylized visual waveforms from mp3 files. Comes with a box painter and a line painter, but can be extended to with other painters easily.

###### Auto generated by spf13/cobra on 19-Dec-2023
