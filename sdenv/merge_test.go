package sdenv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michurin/systemd-env-file/sdenv"
)

func TestEnvironment(t *testing.T) {
	t.Run("nil_input", func(t *testing.T) {
		x, err := sdenv.Environ(nil)
		require.NoError(t, err)
		assert.Nil(t, x)
	})
	t.Run("no_files", func(t *testing.T) {
		a := []string{"x=1"}
		x, err := sdenv.Environ(a)
		require.NoError(t, err)
		assert.Equal(t, []string{"x=1"}, x)
	})
	t.Run("read_files", func(t *testing.T) {
		a := []string{"x=1", "y=1", "z=1"} // a[:1] will be passed to highlight possible side effects
		x, err := sdenv.Environ(a[:1], "testdata/a.env", "testdata/b.env")
		require.NoError(t, err)
		assert.Equal(t, []string{"x=1", "a=1", "b=1"}, x)
		assert.Equal(t, []string{"x=1", "y=1", "z=1"}, a) // no side effects
	})
	t.Run("file_not_found", func(t *testing.T) {
		x, err := sdenv.Environ(nil, "testdata/x.env")
		require.ErrorContains(t, err, "no such file")
		assert.Nil(t, x)
	})
	t.Run("corrupted_file", func(t *testing.T) {
		x, err := sdenv.Environ(nil, "testdata/c.env")
		require.ErrorContains(t, err, "unexpected end of file")
		assert.Nil(t, x)
	})
}
