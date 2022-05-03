package tokenizer

import (
	"strings"
)

type Language interface {
	Token(*Stream) Kind
}

type Kind string

type Token struct {
	Kind    Kind
	Content string
	Col     int
	Line    int
}

const (
	KEYWORD  Kind = "keyword"
	OPERATOR Kind = "operator"
	VARIABLE Kind = "variable"
	STRING   Kind = "string"
	NUMBER   Kind = "number"
	BRACKET  Kind = "bracket"
	ATOM     Kind = "atom"
	DEF      Kind = "def"
	COMMENT  Kind = "comment"
)

func Tokenize(code string, language Language) []Token {
	lines := strings.Split(code, "\n")
	var tokens []Token
	for lineNumber, line := range lines {
		stream := Stream{str: line}
		for {
			if stream.Eol() {
				break
			}
			token := language.Token(&stream)
			if token != "" {
				tokens = append(tokens, Token{Kind: token, Content: stream.Current(), Line: lineNumber, Col: stream.start})
			}
			stream.start = stream.pos
		}
	}
	return tokens
}
