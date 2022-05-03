package languages_test

import (
	"testing"

	"github.com/mrmarble/prism/tokenizer"
	"github.com/mrmarble/prism/tokenizer/languages"
	"github.com/stretchr/testify/require"
)

func TestBrainfuck(t *testing.T) {
	const code = `++++++++++
[              Bucle para iniciar las memorias (se repite diez veces)
	>+++++++>++++++++++>+++++++++++>+++>+<<<<<-
		70        100       110      30  10
]
>++.              imprime 'H'   (72) 1
>>+.              imprime 'o'  (111) 3
---.                      'l'  (108) 3
<---.                     'a'   (97) 2
>>++.                   espacio (32) 4
<+.                       'm'  (109) 3
++++++++.                 'u'  (117) 3
-------.                  'n'  (110) 3
<+++.                     'd'  (100) 2
>+.                       'o'  (111) 3
>+.                       '!'   (33) 4
>.                        '\n' (10) 5`

	tokens := tokenizer.Tokenize(code, &languages.Brainfuck{})

	require := require.New(t)
	require.Len(tokens, 290, "missing tokens")
	require.EqualValues(tokens[0], tokenizer.Token{Kind: "keyword", Content: "+"}, "parsing error")
}
