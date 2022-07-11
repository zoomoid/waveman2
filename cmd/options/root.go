package options

const (
	Filename       string = "file"
	FilenameShort  string = "f"
	Output         string = "output"
	OutputShort    string = "o"
	Recursive      string = "recursive"
	RecursiveShort string = "r"
	Width          string = "width"
	WidthShort     string = "w"
	Height         string = "height"
	HeightShort    string = "y"
)

const (
	FilenameDescription  string = "Determines the file to be sampled, can be relative to the current working directory"
	OutputDescription    string = "Writes the output to a given file. If not specified, writes output to stdout"
	RecursiveDescription string = "Searches for all mp3 files in the directory below the specified file"
	HeightDescription    string = "Height of the shape"
	WidthDescription     string = "Width of each element"
)
