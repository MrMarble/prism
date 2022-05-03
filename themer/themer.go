package themer

import (
	_ "embed"

	"github.com/aymerick/douceur/parser"
	"github.com/mrmarble/prism/tokenizer"
)

//go:embed dark.css
var dark string

type Theme map[tokenizer.Kind]string

func Parse(theme string) (Theme, error) {
	th := Theme{}
	stylesheet, err := parser.Parse(theme)
	if err != nil {
		return nil, err
	}
	for _, rule := range stylesheet.Rules {
		for _, selector := range rule.Selectors {
			for _, declaration := range rule.Declarations {
				if declaration.Property == "color" {
					th[tokenizer.Kind(selector)] = declaration.Value
				}
			}
		}
	}

	return th, nil
}

func Dark() Theme {
	th, _ := Parse(dark)
	return th
}

func (t Theme) GetColor(token tokenizer.Token) string {
	if color, ok := t[token.Kind]; ok {
		return color
	}
	return t["base"]
}
