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

package options

const (
	Filename       string = "file"
	FilenameShort  string = "f"
	Output         string = "output"
	OutputShort    string = "o"
	Recursive      string = "recursive"
	RecursiveShort string = "r"
	Width          string = "width"
	WidthShort     string = "w"
	Height         string = "height"
	HeightShort    string = "y"
)

const (
	FilenameDescription  string = "Determines the file to be sampled, can be relative to the current working directory"
	OutputDescription    string = "Writes the output to a given file. If not specified, writes output to stdout"
	RecursiveDescription string = "Searches for all mp3 files in the directory below the specified file"
	HeightDescription    string = "Height of the shape"
	WidthDescription     string = "Width of each element"
)
