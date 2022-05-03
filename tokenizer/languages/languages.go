package languages

import "github.com/mrmarble/prism/tokenizer"

type LanguageList map[string]func() tokenizer.Language

var languages = map[string]func() tokenizer.Language{
	"brainfuck": func() tokenizer.Language { return &Brainfuck{} },
	"golang":    func() tokenizer.Language { return &Golang{} },
}

func Get(name string) (lang tokenizer.Language, ok bool) {
	if language, ok := languages[name]; ok {
		return language(), ok
	}
	return nil, false
}

func GetAll() LanguageList {
	return languages
}

func List() []string {
	keys := make([]string, 0, len(languages))
	for k := range languages {
		keys = append(keys, k)
	}
	return keys
}
