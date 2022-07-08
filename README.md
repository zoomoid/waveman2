# waveman2

`waveman2` is the Golang successor the first waveman, built to convert audio files
(later, only mp3), into cool-looking, more abstract audio waveforms, by reducing
the samples to a smaller-sized slice and rendering out SVGs.

See <https://github.com/zoomoid/wave-man> for the original project written in
Python. Be warned, decoding is slow, even though the audio file is resampled to
16 times lower sampling rate before processing, and the project depends on
several utilities to decode mp3 files and convert them to PCM before
tranformation.

`waveman2` is designed to be both extensible with other painters than the default
ones, and also to be imported into any other Golang project as a dependency.
This way, you can easily realize the original idea of waveman to be combined
with a web server that processes audio files sent to the server with a set of
defaults and/or user-defined properties, without having to implement web server
functionality in the waveman codebase itself.

`waveman2` comes with a CLI to use for processing audio files. Its usage is similar to
the older waveman, though several flag names changed:

> TODO: implement and document CLI flags with Cobra
