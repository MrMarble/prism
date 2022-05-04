package themer

import (
	"github.com/aymerick/douceur/parser"
	"github.com/mrmarble/prism/tokenizer"
)

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
	return Theme{
		tokenizer.Kind("body"): "#282c34",

		tokenizer.COMMENT:  "#5c6370",
		tokenizer.KEYWORD:  "#c678dd",
		tokenizer.OPERATOR: "#56b6c2",
		tokenizer.VARIABLE: "#e06c75",
		tokenizer.STRING:   "#98c379",
		tokenizer.NUMBER:   "#d19a66",
		tokenizer.ATOM:     "#d19a66",
		tokenizer.BRACKET:  "#abb2bf",
		tokenizer.DEF:      "#e5c07b",
	}
}

func (t Theme) GetColor(token tokenizer.Token) string {
	if color, ok := t[token.Kind]; ok {
		return color
	}

	return "#000"
}
