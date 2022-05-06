package prism_test

import (
	"bytes"
	"testing"

	"github.com/mrmarble/prism"
	"github.com/mrmarble/prism/tokenizer/languages"
	"github.com/sebdah/goldie/v2"
)

const code = `package main

import "fmt"

// main function
func main() {
	fmt.Println("Hello world!")
}`

func BenchmarkEncodePNG(b *testing.B) {
	ctx := prism.NewContext()
	lang, _ := languages.Get("golang")
	ctx.SetLanguage(lang)

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer

		err := ctx.EncodePNG(code, &buf, prism.Options{})
		if err != nil {
			b.Fatal()
		}
	}
}

func TestSavePNG(t *testing.T) {
	options := map[string]prism.Options{
		"default_options":             {},
		"with_header":                 {Header: true},
		"with_numbers":                {LineNumbers: true},
		"with_numbers_header":         {LineNumbers: true, Header: true},
		"with_range":                  {Range: prism.Range{Start: 1, End: 3}},
		"with_range_relative_numbers": {Range: prism.Range{Start: 5, End: 8}, LineNumbers: true, Relative: true},
	}

	for name, option := range options {
		t.Run(name, func(t *testing.T) {
			ctx := prism.NewContext()

			lang, ok := languages.Get("golang")
			if !ok {
				t.Fatal("Unexpected error selecting language")
			}

			ctx.SetLanguage(lang)

			var buf bytes.Buffer

			err := ctx.EncodePNG(code, &buf, option)
			if err != nil {
				t.Fatal("Error generating PNG")
			}
			g := goldie.New(t)
			g.Assert(t, name, buf.Bytes())
		})
	}
}
