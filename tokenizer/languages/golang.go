package languages

import (
	"regexp"

	"github.com/mrmarble/prism/tokenizer"
	"golang.org/x/exp/slices"
)

type Golang struct {
	indented    int
	startOfLine bool
	curPunc     string
	curQuote    string
	tokenize    func(*tokenizer.Stream) tokenizer.Kind
}

func (g *Golang) Token(stream *tokenizer.Stream) tokenizer.Kind {
	if stream.EatSpace() {
		return ""
	}
	if g.tokenize == nil {
		return g.tokenBase(stream)
	}
	return g.tokenize(stream)
}

func (g *Golang) tokenBase(stream *tokenizer.Stream) tokenizer.Kind {
	ch := stream.Next()
	if ch == `"` || ch == "'" || ch == "`" {
		g.curQuote = ch
		g.tokenize = g.tokenString
		return g.tokenString(stream)
	}
	rNumber := regexp.MustCompile(`[\d\.]`)
	if rNumber.MatchString(ch) {
		if ch == "." {
			stream.Match(`^[0-9]+([eE][\-+]?[0-9]+)?`, true)
		} else if ch == "0" {
			stream.Match(`^[xX][0-9a-fA-F]+`, true)
			stream.Match(`^0[0-7]+`, true)
		} else {
			stream.Match(`^[0-9]*\.?[0-9]*([eE][\-+]?[0-9]+)?`, true)
		}
		return tokenizer.NUMBER
	}
	rPunc := regexp.MustCompile(`[\[\]{}\(\),;\:\.]`)
	if rPunc.MatchString(ch) {
		g.curPunc = ch
		return tokenizer.BRACKET
	}
	if ch == "/" {
		if stream.EatString("*") != "" {
			g.tokenize = g.tokenComment
			return g.tokenComment(stream)
		}
		if stream.EatString("/") != "" {
			stream.SkipToEnd()
			return tokenizer.COMMENT
		}
	}
	rOperator := regexp.MustCompile(`[+\-*&^%:=<>!|\/]`)
	if rOperator.MatchString(ch) {
		stream.EatWhile(rOperator)
		return tokenizer.OPERATOR
	}
	regex := regexp.MustCompile(`[\w\$_\xa1-\xff]`)
	stream.EatWhile(regex)
	cur := stream.Current()

	keywords := []string{
		"break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var",
	}

	if slices.Contains(keywords, cur) {
		if cur == "case" || cur == "default" {
			g.curPunc = "case"
		}
		return tokenizer.KEYWORD
	}
	atoms := []string{
		"true", "false", "iota", "nil", "append",
		"cap", "close", "complex", "copy", "delete", "image",
		"len", "make", "new", "panic", "print",
		"println", "real", "recover", "bool", "string",
	}

	if slices.Contains(atoms, cur) {
		return tokenizer.ATOM
	}

	return tokenizer.VARIABLE
}

func (g *Golang) tokenComment(stream *tokenizer.Stream) tokenizer.Kind {
	end := false
	ch := ""
	for {
		ch = stream.Next()
		if ch == "" {
			break
		}
		if ch == "/" && end {
			g.tokenize = nil
			break
		}
		end = ch == "*"
	}
	return tokenizer.COMMENT
}

func (g *Golang) tokenString(stream *tokenizer.Stream) tokenizer.Kind {
	escaped := false
	next := ""
	end := false

	for {
		next = stream.Next()
		if next == "" {
			break
		}
		if next == g.curQuote && !escaped {
			end = true
			break
		}
		escaped = !escaped && g.curQuote != "`" && next == "\\"
	}
	if end || (!escaped || g.curQuote == "`") {
		g.tokenize = nil
	}
	return tokenizer.STRING
}
