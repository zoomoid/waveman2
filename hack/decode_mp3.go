package main

import (
	"log"
	"math"
	"os"

	"github.com/faiface/beep/mp3"
)

func main() {
	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	stream, _, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// d.Length() is the length in bytes. Each sample is encoded in 4 bytes
	// as two channels, each with 16 bits

	chunks := 50
	chunkSize := stream.Len() / chunks
	blocks := make([]float64, chunks)

	log.Printf("%d chunks, %d samples per chunk, %d block", chunks, chunkSize, len(blocks))

	for i := range blocks {
		b := make([][2]float64, chunkSize)
		_, ok := stream.Stream(b)
		if !ok {
			log.Fatalf("failed to stream block %d", i)
		}

		mono := toMono(b)

		block := rootMeanSquare(mono)
		blocks[i] = block
	}
	blocks = normalize(blocks)
	for _, block := range blocks {
		log.Printf("%g", block)
	}
}

func toMono(samples [][2]float64) []float64 {
	o := make([]float64, len(samples))
	for i := 0; i < len(samples); i++ {
		o[i] = (math.Abs(samples[i][0] + samples[i][1]))
		// log.Printf("%g", o[i])
	}
	return o
}

func average(samples []float64) float64 {
	sum := float64(0)
	for _, sample := range samples {
		sum += sample
	}
	o := sum / float64(len(samples))
	return o
}

func meanSquare(samples []float64) float64 {
	sum := float64(0)
	for _, sample := range samples {
		sum += sample * sample
	}
	o := sum / float64(len(samples))
	return o
}

func rootMeanSquare(samples []float64) float64 {
	sum := float64(0)
	for _, sample := range samples {
		sum += sample * sample
	}
	mean := sum / float64(len(samples))
	return math.Sqrt(mean)
}

func normalize(samples []float64) []float64 {
	max := float64(0)
	for _, sample := range samples {
		if sample > max {
			max = sample
		}
	}
	for i, sample := range samples {
		samples[i] = sample / max
	}
	return samples
}
