/*
Copyright 2022 zoomoid.

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

package visitor

import (
	"bytes"
	"os"
	"testing"

	"github.com/zoomoid/waveman2/pkg/streams"
)

// creates a buffer to capture stdout writes while running Visit
func ioFactory() *streams.IO {
	outBuf := bytes.NewBuffer([]byte{})
	return &streams.IO{
		In:     os.Stdin,
		Out:    outBuf,
		ErrOut: os.Stderr,
	}
}

func TestVisitorList(t *testing.T) {
	fw := []fileVisitor{
		{
			path: "../hack/Morgendämmerung.mp3",
		}, {
			path: "../hack/zoomoid - Interstellar.mp3",
		},
	}

	vl := NewVisitorList(fw, ioFactory())

	err := vl.Visit(func(f *File) error {
		t.Log(f.source)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestContinueOnError(t *testing.T) {
	fw := []fileVisitor{
		{
			// path that does not exists
			path: "../../hack/Im Schatten Der Nacht.mp3",
		},
		{
			path: "../../hack/Morgendämmerung.mp3",
		},
	}
	vl := NewVisitorList(fw, ioFactory())

	err := vl.ContinueOnError().Visit(func(f *File) error {
		t.Log(f.source)
		return nil
	})
	if err != nil {
		t.Fatal("returned error from Visit with ContinueOnError")
	}
	errs := vl.Errors()
	// expect exactly one error opening a nonexisting file
	if len(errs) != 1 {
		t.Fatalf("%v", errs)
	}

}

func TestUseStdout(t *testing.T) {
	fw := []fileVisitor{
		{
			// path that does not exists
			path: "../../hack/Morgendämmerung.mp3",
		},
	}
	s := ioFactory()
	vl := NewVisitorList(fw, s)

	err := vl.UseStdout(true).Visit(func(f *File) error {
		var buf bytes.Buffer
		_, err := buf.WriteString(f.source)
		if err != nil {
			t.Fatal(err)
		}
		f.Print(&buf)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	b, _ := s.Out.(*bytes.Buffer)

	if b.Len() == 0 {
		t.Fatalf("expected output to contain data, found %s", b.String())
	}
}
