package sdenv_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/michurin/systemd-env-file/sdenv"
)

func ExampleCollection() {
	c := sdenv.NewCollectsion()
	c.PushStd([]string{ // adding first bunch
		"A=1",
		"B=2",
		"A=3", // redefine A=1
	})
	c.PushStd([]string{ // adding second bunch
		"A=4", // will be skipped as known from previous bunch
		"C=5",
	})
	fmt.Println(strings.Join(c.CollectionStd(), "\n"))
	// output:
	// B=2
	// A=3
	// C=5
}

func TestCollectsion(t *testing.T) {
	t.Parallel()
	c := sdenv.NewCollectsion()
	c.Push([][2]string{
		{"A", "1"},
		{"B", "1"},
		{"C", "1"},
		{"A", "2"},
	})
	c.Push([][2]string{
		{"A", "3"},
		{"D", "1"},
		{"D", "2"},
	})
	assert.Equal(t, [][2]string{
		{"B", "1"}, // from first sequece
		{"C", "1"},
		{"A", "2"},
		{"D", "2"}, // from second sequene
	}, c.Collection())
}

func TestCollectsion_std(t *testing.T) {
	t.Parallel()
	c := sdenv.NewCollectsion()
	c.PushStd([]string{"A=1", "B=", "=3", "=", ""})
	assert.Equal(t, []string{
		"A=1",
		"B=",
		"=", // all empty keys are collapsed, as expected
	}, c.CollectionStd())
}
