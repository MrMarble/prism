package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/mrmarble/prism"
	"github.com/mrmarble/prism/themer"
	"github.com/mrmarble/prism/tokenizer/languages"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var CLI struct {
	File     string `arg:"" type:"existingfile" help:"File with code"`
	Language string `name:"lang" short:"l" help:"Language to parse." required:""`
	Output   string `name:"output" short:"o" help:"output image" type:"path"`

	Numbers bool `short:"n" help:"display line numbers"`
	Header  bool `help:"display header"`

	Debug bool `help:"Debug logging"`
}

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, PartsExclude: []string{"time"}})
}

func main() {
	ctx := kong.Parse(&CLI, kong.UsageOnError(), kong.Description(fmt.Sprintf("Create beautiful images of your source code from your command line.\n\nSupported languages: %s", strings.Join(languages.List(), ", "))))
	if CLI.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	ctx.FatalIfErrorf(run(ctx))
}

func run(ctx *kong.Context) error {
	pr := prism.NewContext()
	log.Debug().Msg("Created context")

	pr.SetTheme(themer.Dark())
	log.Debug().Msg("Theme set")

	if lang, ok := languages.Get(CLI.Language); ok {
		pr.SetLanguage(lang)
	}

	log.Debug().Msg("Language set")

	code, err := ioutil.ReadFile(CLI.File)
	if err != nil {
		return err
	}
	log.Debug().Msg("Code loaded")

	options := prism.Options{}

	if CLI.Numbers {
		options.LineNumbers = true
	}
	if CLI.Header {
		options.Header = true
	}
	if CLI.Output != "" {
		err := pr.SavePNG(string(code), CLI.Output, options)
		if err != nil {
			return err
		}
		log.Info().Str("Output", CLI.Output).Msg("Image saved!")
	} else {
		err := pr.SavePNG(string(code), "prism.png", options)
		if err != nil {
			return err
		}
		log.Info().Str("Output", "prism.png").Msg("Image saved!")

	}

	return nil
}
