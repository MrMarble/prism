package languages_test

import (
	"testing"

	"github.com/mrmarble/prism/tokenizer"
	"github.com/mrmarble/prism/tokenizer/languages"
	"github.com/stretchr/testify/require"
)

func TestGolang(t *testing.T) {
	const code = `package main

import "fmt"

func main() {
	fmt.Println("Hello world!")
}`

	tokens := tokenizer.Tokenize(code, &languages.Golang{})

	require := require.New(t)
	require.Len(tokens, 16)
	require.EqualValues(tokens[0], tokenizer.Token{Kind: "keyword", Content: "package"})
}
