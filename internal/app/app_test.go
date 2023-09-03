package app_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michurin/systemd-env-file/internal/app"
)

func TestApp_ok(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	err := app.App(nil, []string{"testdata/script.sh", "1", "2"}, stdout, stderr, []string{
		"NOT_EXISTS.env", // will be skipped
		"/dev/null",      // will be skipped due to node mode
		"testdata/env.env",
	})

	require.NoError(t, err)
	assert.Equal(t, "args=1 2\nTEST_VAR=ok\n", stdout.String())
	assert.Equal(t, "", stderr.String())
}

func TestApp_errReadingEnvFile(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	err := app.App(nil, []string{"testdata/script.sh", "1", "2"}, stdout, stderr, nil)

	assert.EqualError(t, err, "cannot open env file: open xenv.env: no such file or directory")
	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "", stderr.String())
}

func TestApp_errNoCmd(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	err := app.App(nil, nil, stdout, stderr, nil)

	assert.EqualError(t, err, "you are to specify command")
	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "", stderr.String())
}

func TestApp_errInvalidCmd(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	err := app.App(nil, []string{"NOT_EXISTS.exe"}, stdout, stderr, []string{"testdata/env.env"})

	assert.EqualError(t, err, `cannot run command: exec: "NOT_EXISTS.exe": executable file not found in $PATH`)
	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "", stderr.String())
}
