/*
Copyright 2022-2023 zoomoid.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package transform

import (
	"errors"
	"io"
	"math"

	"github.com/hajimehoshi/go-mp3"
)

const (
	DefaultGoMp3Channels   int = 2
	DefaultGoMp3Precision  int = 2
	DefaultGoMp3FrameWidth int = DefaultGoMp3Channels * DefaultGoMp3Precision
)

type Mp3Decoder struct {
	channels  int
	precision int
	width     int
	decoder   *mp3.Decoder
}

func newDecoder(f io.Reader) (*Mp3Decoder, error) {
	d, err := mp3.NewDecoder(f)
	if err != nil {
		return nil, errors.New("failed to construct decoder")
	}

	decoder := &Mp3Decoder{
		channels:  DefaultGoMp3Channels,
		precision: DefaultGoMp3Precision,
		width:     DefaultGoMp3FrameWidth,
		decoder:   d,
	}

	return decoder, nil
}

// Length returns the total size in bytes.
//
// Length returns -1 when the total size is not available e.g. when the given source is not io.Seeker.
//
// Wrapper for (mp3.Decoder).Length
func (d *Mp3Decoder) length() int {
	return int(d.decoder.Length())
}

// Fills the samples slice with len(samples) samples.
//
// Wrapper for (mp3.Decoder).Read
func (d *Mp3Decoder) read(samples [][2]float64) (n int, err error) {
	var tmp [DefaultGoMp3FrameWidth]byte
	for i := range samples {
		dn, err := d.decoder.Read(tmp[:])
		if dn == len(tmp) {
			samples[i], _ = d.decode(tmp[:])
			n++
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// Seek is io.Seeker's Seek.
//
// Seek returns an error when the underlying source is not io.Seeker.
//
// Note that seek uses a byte offset but samples are aligned to 4 bytes (2 channels, 2 bytes each).
// Be careful to seek to an offset that is divisible by 4 if you want to read at full sample boundaries.
// Wrapper for (mp3.Decoder).Seek
func (d *Mp3Decoder) seek(offset int64, whence int) (int64, error) {
	return d.decoder.Seek(offset, whence)
}

func (d *Mp3Decoder) decode(p []byte) (sample [2]float64, n int) {
	for c := range sample {
		x, n := decodeFloat(d.precision, p)
		sample[c] = x
		p = p[n:]
	}
	for c := len(sample); c < 2; c++ {
		_, n := decodeFloat(d.precision, p)
		p = p[n:]
	}
	return sample, d.width
}

func decodeFloat(precision int, p []byte) (x float64, n int) {
	var xUint64 uint64
	for i := precision - 1; i >= 0; i-- {
		xUint64 <<= 8
		xUint64 += uint64(p[i])
	}
	return toFloat(precision, xUint64), precision
}

func toFloat(precision int, xUint64 uint64) float64 {
	if xUint64 >= 1<<uint(precision*8-1) {
		compl := 1<<uint(precision*8) - xUint64
		return -float64(int64(compl)) / (math.Exp2(float64(precision)*8-1) - 1)
	}
	return float64(int64(xUint64)) / (math.Exp2(float64(precision)*8-1) - 1)
}
