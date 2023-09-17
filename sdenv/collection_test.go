package sdenv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/michurin/systemd-env-file/sdenv"
)

func TestCollectsion(t *testing.T) {
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
	c := sdenv.NewCollectsion()
	c.PushStd([]string{"A=1", "B=", "=3", "=", ""})
	assert.Equal(t, []string{
		"A=1",
		"B=",
		"=", // all empty keys are collapsed, as expected
	}, c.CollectionStd())
}
