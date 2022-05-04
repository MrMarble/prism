package prism_test

import (
	"bytes"
	"testing"

	"github.com/mrmarble/prism"
	"github.com/mrmarble/prism/tokenizer/languages"
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
