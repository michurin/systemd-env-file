package sdenv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michurin/systemd-env-file/sdenv"
)

func TestParser_ok(t *testing.T) {
	for _, cs := range []struct {
		name  string
		data  string
		pairs [][2]string
	}{
		{
			name:  "simple",
			data:  "x=ok",
			pairs: [][2]string{{"x", "ok"}},
		},
		{
			name:  "spaces",
			data:  "\ta\tb\t=\to\tk\t",
			pairs: [][2]string{{"a\tb", "o\tk"}},
		},
		{
			name:  "novalue",
			data:  "x=",
			pairs: [][2]string{{"x", ""}},
		},
		{
			name:  "novalue_spaces",
			data:  "\tx\t=\t\n",
			pairs: [][2]string{{"x", ""}},
		},
		{
			name:  "single_quoted",
			data:  "x = ' A \n B '",
			pairs: [][2]string{{"x", " A \n B "}},
		},
		{
			name:  "double_quoted",
			data:  "x = \" A \n B \"",
			pairs: [][2]string{{"x", " A \n B "}},
		},
		{
			name:  "double_quoted_escape", // i.e. it unescape only ", $, ` and \
			data:  "x = \" \\\" \\\n \\$ \\x \"",
			pairs: [][2]string{{"x", ` "  $ \x `}},
		},
		{
			name: "complex",
			data: `
			lines without equal character are skipped
			# comments can be started by # and ;
			# backslash continue comments \
			this=is_still_comment
			; and comments have to start at the start
			; of line
			this = comment ; considering as part of value
			just Key = just Value
			Ex1 = 'single quotes ' works like that
			Ex2 = however, you can't use quotes inline
			Ex3 = "double quotes " are doing the "same" way
			Ex4 = "quotes allows escapes: \n \x"
			Ex5 = "quotes are doing trick for chars: \", \\, \$"
			Ex6 = naked values allows escapes: \a, \b, \c
			Ex7 = and even \$, and even at the beginning
			Ex8 = \x all the rest
			Ex9 = in naked values backshash \
allows splitting values
			Ex10 = "quotes
allow multi lines"
			Ex11 = "escape in double quotes \
join lines"
			`,
			pairs: [][2]string{
				{"this", "comment ; considering as part of value"},
				{"just Key", "just Value"},
				{"Ex1", `single quotes works like that`},
				{"Ex2", `however, you can't use quotes inline`},
				{"Ex3", `double quotes are doing the "same" way`},
				{"Ex4", `quotes allows escapes: \n \x`},
				{"Ex5", `quotes are doing trick for chars: ", \, $`},
				{"Ex6", `naked values allows escapes: a, b, c`},
				{"Ex7", `and even $, and even at the beginning`},
				{"Ex8", `x all the rest`},
				{"Ex9", `in naked values backshash allows splitting values`},
				{"Ex10", "quotes\nallow multi lines"},
				{"Ex11", "escape in double quotes join lines"},
			},
		},
	} {
		cs := cs
		t.Run(cs.name, func(t *testing.T) {
			kv, err := sdenv.Parser([]rune(cs.data))
			require.NoError(t, err)
			assert.Equal(t, cs.pairs, kv)
		})
	}
}

func TestParser_error(t *testing.T) {
	kv, err := sdenv.Parser([]rune("ok='"))
	require.Error(t, err)
	require.Nil(t, kv)
}
