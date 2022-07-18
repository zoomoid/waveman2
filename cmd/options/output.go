package options

type OutputType string

const (
	OutputTypeFile  OutputType = "file"
	OutputTypeEmpty OutputType = ""
)

var (
	SupportedOutputs = []OutputType{OutputTypeFile}
)
