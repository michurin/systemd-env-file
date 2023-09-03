package sdenv_test

import (
	"fmt"

	"github.com/michurin/systemd-env-file/sdenv"
)

func ExampleParser() {
	const envFileContent = `
; comment
# comment
# empty lines are allowed

by the way, strings without equals sign are considered as comments

simpleKey=simpleValue

spaces = around = are allowed

multiline values = are \
possible \
with joins

true multiline values = "are
possible too"

however = comments # have to start at the beginging of line
`

	keyValues, err := sdenv.Parser([]byte(envFileContent))
	if err != nil {
		panic(err)
	}
	for _, kv := range keyValues {
		fmt.Printf("key=%q, value=%q\n", kv[0], kv[1])
	}
	// output:
	// key="simpleKey", value="simpleValue"
	// key="spaces", value="around = are allowed"
	// key="multiline values", value="are possible with joins"
	// key="true multiline values", value="are\npossible too"
	// key="however", value="comments # have to start at the beginging of line"
}
