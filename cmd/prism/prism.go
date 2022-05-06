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

type VersionFlag string

var (
	// Populated by goreleaser during build
	version = "master"
	commit  = "?"
	date    = ""
)

var CLI struct {
	File     string      `arg:"" type:"existingfile" help:"File with code"`
	Language string      `name:"lang" short:"l" help:"Language to parse." required:""`
	Output   string      `name:"output" short:"o" help:"Output file name" type:"path" default:"prism.png"`
	Version  VersionFlag `name:"version" help:"Print version information and quit"`

	Header   bool   `help:"Display header"`
	Lines    string `help:"Specify a range of lines instead of reading the whole file. Ex: 10-20"`
	Numbers  bool   `short:"n" help:"Display line numbers"`
	Relative bool   `short:"r" help:"Use relative numbers. Needs --lines and --numbers"`
}

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong) error {
	fmt.Printf("Prism has version %s built from %s on %s\n", version, commit, date)
	app.Exit(0)

	return nil
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("prism"),
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
	} else {
		return fmt.Errorf("Invalid language")
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

	if CLI.Lines != "" {
		r := prism.Range{}

		err = r.Parse(CLI.Lines)
		if err != nil {
			return err
		}

		options.Range = r
	}

	if CLI.Relative {
		options.Relative = true
	}

	err = pr.SavePNG(string(code), CLI.Output, options)
	if err != nil {
		return err
	}

	fmt.Printf("Image created at %s\n", CLI.Output)

	return nil
}
