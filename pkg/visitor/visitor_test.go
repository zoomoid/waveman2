package visitor

import (
	"bytes"
	"os"
	"testing"

	"github.com/zoomoid/waveman/v2/pkg/streams"
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
			path: "../hack/Im Schatten Der Nacht.mp3",
		},
		{
			path: "../hack/Morgendämmerung.mp3",
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
