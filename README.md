![logo](assets/logo.svg)
# Prism

![GitHub](https://img.shields.io/github/license/mrmarble/prism)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mrmarble/prism)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/mrmarble/prism)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrmarble/prism)](https://goreportcard.com/report/github.com/mrmarble/prism)
[![Go Reference](https://pkg.go.dev/badge/github.com/mrmarble/prism.svg)](https://pkg.go.dev/github.com/mrmarble/prism)

Create beautiful images of your source code directly from your terminal.

![example](assets/example.png)

# Installation

Precompiled `prism` binaries can be found at [releases](https://github.com/mrmarble/prism/releases) page.

Instructions below show how to build `prism` from sources.

```sh
go install github.com/mrmarble/prism/cmd/prism@latest # or target a specific version @v0.1.0
```
# Usage
Be sure `prism` executable is under your `$PATH`.

Usage of **prism**: prism [input file] [args...] Run prism without arguments to get help output.

```
Flags:
  -h, --help                  Show context-sensitive help.
  -l, --lang=STRING           Language to parse.
  -o, --output="prism.png"    output image
      --version               Print version information and quit
  -n, --numbers               display line numbers
      --header                display header
```

## Supported languages

See [languages](tokenizer/languages/) for a list of implemented languages.