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
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/zoomoid/waveman2/pkg/streams"
)

// Path takes a set of paths (recursive or not) and turns them into a list of FileVisitors,
// adding them to the visitorList wrapper type and records any errors encountered.
func (vl *VisitorList) File(recursive bool, useStdout bool, paths ...string) *VisitorList {
	for _, p := range paths {
		_, err := os.Stat(p)
		if errors.Is(err, fs.ErrNotExist) {
			vl.errors = append(vl.errors, fmt.Errorf(pathNotExistError, p))
		}
		if err != nil {
			vl.errors = append(vl.errors, fmt.Errorf("path %q cannot be accessed: %w", p, err))
			continue
		}

		visitors, err := expandPathsToFileVisitors(p, recursive, SupportedFileExtensions)
		if err != nil {
			vl.errors = append(vl.errors, fmt.Errorf("error reading %q: %w", p, err))
		}
		vl.visitors = append(vl.visitors, visitors...)
	}

	if len(vl.visitors) == 0 && len(vl.errors) == 0 {
		vl.errors = append(vl.errors, fmt.Errorf("error reading %v: supported file extensions are %v", paths, SupportedFileExtensions))
	}

	return vl
}

// VisitorFunc is the function type used for all Visitor implementations
type VisitorFunc func(*File) error

// fileVisitor is instantiated with a source path and whether to use
// Stdout as writer
type fileVisitor struct {
	path string
}

type VisitorList struct {
	visitors        []fileVisitor
	continueOnError bool
	errors          []error
	useStdout       bool
	io              *streams.IO
}

func NewVisitorList(visitors []fileVisitor, io *streams.IO) *VisitorList {
	if io == nil {
		io = streams.DefaultStreams
	}
	return &VisitorList{visitors: visitors, io: io}
}

// ContinueOnError sets the continueOnError flag to true, meaning
// the Visit function of a VisitorList will not return early on error,
// and instead record the error, but finish visiting all Visitors instantiated.
func (v *VisitorList) ContinueOnError() *VisitorList {
	v.continueOnError = true
	return v
}

// UseStdout sets the stdout property on a visitor list
func (v *VisitorList) UseStdout(useStdout bool) *VisitorList {
	v.useStdout = useStdout
	return v
}

// Visit is the canonic Visit implementation for a list of Visitors
// Returns an error on the first error when ContinueOnError is not
// called beforehand, otherwise aggregates all errors in the list of errors
// and returns nil
func (v *VisitorList) Visit(fn VisitorFunc) error {
	for _, visitor := range v.visitors {
		err := visitor.visit(v.useStdout, v.io, fn)
		if err != nil {
			if !v.continueOnError {
				return err
			}
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

// Errors returns a list of errors kept by the internal field on a VisitorList
func (v *VisitorList) Errors() []error {
	return v.errors
}

// Len is the canonic len function on the internal list of visitors of a VisitorList
func (v *VisitorList) Len() int {
	return len(v.visitors)
}

// Visit implements the Visitor interface for fileVisitors by instantiating a File
// struct with all the required data from the source filename and whether to use
// stdout as a writer
func (v *fileVisitor) visit(useStdout bool, streams *streams.IO, fn VisitorFunc) error {
	var f *os.File
	var err error
	f, err = os.Open(v.path)
	if err != nil {
		return err
	}
	defer f.Close()

	p := filepath.Base(v.path)
	dir := filepath.Dir(v.path)
	ext := filepath.Ext(p)
	r, err := regexp.Compile(fmt.Sprintf("%s$", ext))
	if err != nil {
		return err
	}
	bare := r.ReplaceAllString(p, "")
	svgFile := r.ReplaceAllString(p, DefaultSVGExtension)
	svgPath := filepath.Join(dir, svgFile)

	// open writer to stdout, and only create file if stdout is not selected
	// for a given fileVisitor
	var writer io.Writer
	if !useStdout {
		w, err := os.Create(svgPath)
		if err != nil {
			return err
		}
		writer = w
		defer w.Close()
	} else {
		w := streams.Out
		writer = w
	}

	file := &File{
		source:    v.path,
		dir:       dir,
		filename:  bare,
		extension: ext,
		output:    svgPath,
		reader:    f,
		writer:    writer,
	}

	return fn(file)
}
