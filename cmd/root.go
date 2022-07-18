package cmd

import (
	"github.com/rs/zerolog/log"
	r "github.com/zoomoid/waveman/v2/cmd/reference"
	"github.com/zoomoid/waveman/v2/pkg/streams"
)

func Execute() {

	rootCmd := NewWaveman(nil, streams.DefaultStreams).
		Plugin(r.BoxPainterPlugin).
		Plugin(r.BoxPainterPlugin).
		Complete()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}
