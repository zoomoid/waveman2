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
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/zoomoid/waveman2/pkg/streams"
)

const (
	pathNotExistError string = "the path %q does not exist"
)

var SupportedFileExtensions = []string{".mp3"}

const (
	DefaultSVGExtension string = ".svg"
)

type File struct {
	source    string
	filename  string
	dir       string
	extension string
	output    string

	reader io.Reader
	writer io.Writer
}

// Print writes a buffer to a file's writer, i.e., either Stdout or
// the fresh SVG created for when more than one mp3 is fed into waveman
func (f *File) Print(data *bytes.Buffer) error {
	_, err := f.writer.Write(data.Bytes())
	return err
}

func (f *File) Reader() io.Reader {
	return f.reader
}

// expandPaths transforms all filename flag arguments into fileVisitors,
// and combines them in a VisitorList wrapper type
//
// Runs validation on the filename options struct and returns a list of errors
// recorded during validation
func ExpandPaths(paths []string, recursive bool, useStdout bool, io *streams.IO) (*VisitorList, []error) {
	var errs []error

	// allocate a new VisitorList
	vl := NewVisitorList(nil, io)

	for _, s := range paths {
		matches, err := expandIfGlob(s)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		vl.File(recursive, useStdout, matches...)
	}
	return vl, errs
}

// ignoreFile is a filter function that skips all files with
// non-matching file extensions (i.e. no mp3 files)
func ignoreFile(path string, extensions []string) bool {
	if len(extensions) == 0 {
		return false
	}
	ext := filepath.Ext(path)
	for _, s := range extensions {
		if s == ext {
			return false
		}
	}
	return true
}

// expandIfGlob attempts to expand a pattern and returns either a list of
// all matches or the original pattern if expansion fails
func expandIfGlob(pattern string) ([]string, error) {
	if _, err := os.Stat(pattern); errors.Is(err, fs.ErrNotExist) {
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) == 0 {
			return nil, fmt.Errorf(pathNotExistError, pattern)
		}
		if err == filepath.ErrBadPattern {
			return nil, fmt.Errorf("patterns %q is not valid: %w", pattern, err)
		}
		return matches, err
	}
	return []string{pattern}, nil
}

// For a given path, create fileVisitors for all files matching one of the extensions in the directory tree
// below.
// Explore subdirectories iff recursive is set to true.
// useStdout sets a flag such that during Visit, Stdout is used as a io.Writer instead of creating a new
// file handle.
func expandPathsToFileVisitors(paths string, recursive bool, extensions []string) ([]fileVisitor, error) {
	var visitors []fileVisitor
	err := filepath.WalkDir(paths, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != paths && !recursive {
				return filepath.SkipDir
			}
			return nil
		}

		if path != paths && ignoreFile(path, extensions) {
			return nil
		}

		v := fileVisitor{
			path: path,
		}

		visitors = append(visitors, v)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return visitors, nil
}
