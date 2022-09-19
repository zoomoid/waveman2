## waveman completion

Generate completion script

### Synopsis

To load completions:
	
	Bash:
	
		$ source <(waveman completion bash)
	
		# To load completions for each session, execute once:
		# Linux:
		$ waveman completion bash > /etc/bash_completion.d/waveman
		# macOS:
		$ waveman completion bash > $(brew --prefix)/etc/bash_completion.d/waveman
	
	Zsh:
	
		# If shell completion is not already enabled in your environment,
		# you will need to enable it.  You can execute the following once:
	
		$ echo "autoload -U compinit; compinit" >> ~/.zshrc
	
		# To load completions for each session, execute once:
		$ waveman completion zsh > "${fpath[1]}/_waveman"
	
		# You will need to start a new shell for this setup to take effect.
	
	fish:
	
		$ waveman completion fish | source
	
		# To load completions for each session, execute once:
		$ waveman completion fish > ~/.config/fish/completions/waveman.fish
	
	PowerShell:
	
		PS> waveman completion powershell | Out-String | Invoke-Expression
	
		# To load completions for every new session, run:
		PS> waveman completion powershell > waveman.ps1
		# and source this file from your PowerShell profile.
	

```
waveman completion [bash|zsh|fish|powershell]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --aggregator string          Determines the type of aggregator function to use. Chose one of 'max', 'avg', 'rounded-avg', 'mean-square', or 'root-mean-square' (default "rms")
      --alignment string           Alignment of the shapes, chose one of 'top', 'center', or 'bottom' (default "center")
      --box                        Create a box waveform
  -n, --chunks int                 Chunks are the number of samples in the output of a transformation. For the Box painter, this also means the number of blocks, and for the Line painter, the number of root points of the line (default 64)
  -c, --closed                     Whether the SVG path should be closed or left open
      --color string               Fill color of each box (default "black")
      --downsampling-factor int    Determines the ratio of samples being used for downsampling compared to the full chunk's length. Given in powers of two up two 128 (default 1)
      --downsampling-mode string   Determines the downsampling mode, either by sampling samples from the start, the center, or the end of a chunk
  -f, --file strings               Determines the file to be sampled, can be relative to the current working directory
      --fill-color string          Color for the area enclosed by the line (default "rgba(0 0 0 / 0.5)")
      --gap float                  Gap is the spacing left between each box. Boxes are centered horizonally, so half of gap is subtracted from the box's width (default 5)
  -y, --height float               Height of the shape (default 200)
      --interpolation string       Interpolation mechanism to be used for smoothing the curve. Choose one of 'none', 'fritsch-carlson', or 'steffen' (default "fritsch-carlson")
  -i, --inverted                   Whether the shape should be inverted horizontally, i.e., switch the vertical alignment from top to bottom
      --line                       Create a line waveform
  -o, --output string              Writes the output to a given file. If not specified, writes output to stdout
  -r, --recursive                  Searches for all mp3 files in the directory below the specified file
      --rounded float              Rounding factor of each box. Given in pixels. See SVG <rect> rx/ry attributes for details (default 10)
      --stroke-color string        Color of the line's stroke (default "black")
      --stroke-width float         Width of the line's stroke (default 5)
  -w, --width float                Width of each element (default 10)
```

### SEE ALSO

* [waveman](waveman.md)	 - waveman generates stylized visual waveforms from mp3 files. Comes with a box painter and a line painter, but can be extended to with other painters easily.

###### Auto generated by spf13/cobra on 19-Sep-2022