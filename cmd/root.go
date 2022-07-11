package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
)

func Execute() {

	streams := &IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	rootCmd := NewWaveman(nil, streams).
		Plugin(LinePainterPlugin).
		Plugin(BoxPainterPlugin).
		Complete()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}
