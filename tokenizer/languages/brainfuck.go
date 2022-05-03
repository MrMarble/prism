package languages

import (
	"strings"

	"github.com/mrmarble/prism/tokenizer"
)

type Brainfuck struct {
	commentLine bool
	commentLoop bool
	left        int
	right       int
}

func (b *Brainfuck) Token(stream *tokenizer.Stream) tokenizer.Kind {
	if stream.EatSpace() {
		return ""
	}
	if stream.Sol() {
		b.commentLine = false
	}
	ch := stream.Next()
	if strings.Contains("><+-.,[]", ch) {
		if b.commentLine {
			if stream.Eol() {
				b.commentLine = false
			}
			return tokenizer.COMMENT
		}
		switch ch {
		case "]":
			b.left++
			return tokenizer.BRACKET
		case "[":
			b.right++
			return tokenizer.BRACKET
		case "+", "-":
			return tokenizer.KEYWORD
		case "<", ">":
			return tokenizer.ATOM
		case ".", ",":
			return tokenizer.DEF
		}
	} else {
		b.commentLine = true
		if stream.Eol() {
			b.commentLine = false
		}
		return tokenizer.COMMENT
	}

	if stream.Eol() {
		b.commentLine = false
	}
	return ""
}
