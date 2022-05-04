package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/mrmarble/prism"
	"github.com/mrmarble/prism/themer"
	"github.com/mrmarble/prism/tokenizer/languages"
)

var CLI struct {
	File     string `arg:"" type:"existingfile" help:"File with code"`
	Language string `name:"lang" short:"l" help:"Language to parse." required:""`
	Output   string `name:"output" short:"o" help:"output image" type:"path" default:"prism.png"`

	Numbers bool `short:"n" help:"display line numbers"`
	Header  bool `help:"display header"`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.UsageOnError(),
		kong.Description(fmt.Sprintf(`Create beautiful images of your source code from your terminal.

Supported languages: %s`, strings.Join(languages.List(), ", "))),
	)

	ctx.FatalIfErrorf(run(ctx))
}

func run(ctx *kong.Context) error {
	pr := prism.NewContext()

	pr.SetTheme(themer.Dark())

	if lang, ok := languages.Get(CLI.Language); ok {
		pr.SetLanguage(lang)
	}

	code, err := ioutil.ReadFile(CLI.File)
	if err != nil {
		return err
	}

	options := prism.Options{}

	if CLI.Numbers {
		options.LineNumbers = true
	}

	if CLI.Header {
		options.Header = true
	}

	err = pr.SavePNG(string(code), CLI.Output, options)
	if err != nil {
		return err
	}

	fmt.Printf("Image created at %s\n", CLI.Output)

	return nil
}
