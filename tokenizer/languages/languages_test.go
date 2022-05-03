package languages_test

import (
	"testing"

	"github.com/mrmarble/prism/tokenizer"
	"github.com/mrmarble/prism/tokenizer/languages"
)

func TestInterface(t *testing.T) {
	langs := languages.GetAll()
	for name, factory := range langs {
		t.Run(name, func(t *testing.T) {
			iface := factory()
			_, ok := iface.(tokenizer.Language)
			if !ok {
				t.Log("does not implement Language interface")
				t.Fail()
			}
		})
	}
}
