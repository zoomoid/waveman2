package options

const (
	DownsamplingMode   string = "downsampling-mode"
	DownsamplingFactor string = "downsampling-factor"
	Aggregator         string = "aggregator"
	Filename           string = "file"
	FilenameShort      string = "f"
	Chunks             string = "chunks"
	ChunksShort        string = "n"
	Output             string = "output"
	OutputShort        string = "o"
)

const (
	DownsamplingModeDescription   string = "Determines the downsampling mode, either by sampling samples from the start, the center, or the end of a chunk"
	DownsamplingFactorDescription string = "Determines the ratio of samples being used for downsampling compared to the full chunk's length. Given in powers of two up two 128"
	AggregatorDescription         string = "Determines the type of aggregator function to use. Chose one of 'max', 'avg', 'rounded-avg', 'mean-square', or 'root-mean-square'"
	FilenameDescription           string = "Determines the file to be sampled, can be relative to the current working directory"
	ChunksDescription             string = "Chunks are the number of samples in the output of a transformation. For the Box painter, this also means the number of blocks, and for the Line painter, the number of root points of the line"
	OutputDescription             string = "Writes the output to a given file. If not specified, writes output to stdout"
)
