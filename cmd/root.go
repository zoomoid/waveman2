package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman/v2/cmd/box"
	"github.com/zoomoid/waveman/v2/cmd/line"
	"github.com/zoomoid/waveman/v2/cmd/options"
	"github.com/zoomoid/waveman/v2/pkg/transform"
)

func Execute() {
	rootCmd := NewRootCommand(nil)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

type rootData struct {
	downsamplingMode   string
	downsamplingFactor int
	aggregator         string
	filename           string
	chunks             int
}

func newRootData() *rootData {
	return &rootData{
		downsamplingMode:   string(transform.DefaultDownsamplingMode),
		downsamplingFactor: int(transform.DefaultPrecision),
		aggregator:         string(transform.DefaultTransformerMode),
		filename:           "",
		chunks:             transform.DefaultChunks,
	}
}

func NewRootCommand(data *rootData) *cobra.Command {
	if data == nil {
		data = newRootData()
	}

	cmd := &cobra.Command{
		Use: "waveman",
		// TODO: add short and long description
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	addTransformerOptions(cmd.Flags(), data)

	cmd.AddCommand(line.NewCommand(nil))
	cmd.AddCommand(box.NewCommand(nil))

	return cmd
}

func addTransformerOptions(flags *pflag.FlagSet, data *rootData) {
	flags.StringVar(&data.downsamplingMode, options.DownsamplingMode, "", options.DownsamplingModeDescription)
	flags.IntVar(&data.downsamplingFactor, options.DownsamplingFactor, 1, options.DownsamplingFactorDescription)
	flags.StringVar(&data.aggregator, options.Aggregator, string(transform.DefaultTransformerMode), options.AggregatorDescription)
	flags.StringVarP(&data.filename, options.Filename, options.FilenameShort, "", options.FilenameDescription)
	flags.IntVarP(&data.chunks, options.Chunks, options.ChunksShort, transform.DefaultChunks, options.ChunksDescription)
}
