//go:build X_DISABLED
// +build X_DISABLED

// TODO
// app.App has to be split into Flags, EnvFileLookup, Exec, and tests
// has to be split as well.

package app_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michurin/systemd-env-file/internal/app"
)

func run(t *testing.T, env, args, files []string) (int, string, string, error) {
	t.Helper()
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	ec, err := app.App(env, append([]string{"xenv"}, args...), stdout, stderr, files)
	return ec, stdout.String(), stderr.String(), err
}

func TestApp_ok(t *testing.T) {
	ec, stdout, stderr, err := run(t, nil, []string{"testdata/script.sh", "1", "2"}, []string{
		"NOT_EXISTS.env", // will be skipped
		"/dev/null",      // will be skipped due to node mode
		"testdata/env.env",
	})

	require.NoError(t, err)
	assert.Zero(t, ec)
	assert.Equal(t, "args=1 2\nTEST_VAR=ok\n", stdout)
	assert.Zero(t, stderr)
}

func TestApp_dontOverrideExiting(t *testing.T) {
	ec, stdout, stderr, err := run(t, []string{"TEST_VAR=x"}, []string{"testdata/script.sh"}, []string{"testdata/env.env"})

	require.NoError(t, err)
	assert.Zero(t, ec)
	assert.Equal(t, "args=\nTEST_VAR=x\n", stdout)
	assert.Zero(t, stderr)
}

func TestApp_flagH(t *testing.T) {
	ec, stdout, stderr, err := run(t, nil, []string{"-h"}, nil)

	require.NoError(t, err)
	assert.Zero(t, ec)
	assert.Zero(t, stdout)
	assert.Equal(t,
		"usage: xenv [flags] command [args ...]\n"+
			"  -h\tshow help message\n"+
			"  -v\tshow version\n",
		stderr)
}

func TestApp_flagV(t *testing.T) {
	ec, stdout, stderr, err := run(t, nil, []string{"-v"}, nil)

	require.NoError(t, err)
	assert.Zero(t, ec)
	assert.Zero(t, stdout)
	assert.Contains(t, stderr, "build info:\ngo\tgo1")
}

func TestApp_flagWrong(t *testing.T) {
	ec, stdout, stderr, err := run(t, nil, []string{"-X"}, nil)

	assert.Error(t, err, "flag provided but not defined: -X")
	assert.Equal(t, 2, ec)
	assert.Zero(t, stdout)
	assert.Equal(t,
		"flag provided but not defined: -X\n"+
			"usage: xenv [flags] command [args ...]\n"+
			"  -h\tshow help message\n"+
			"  -v\tshow version\n",
		stderr)
}

func TestApp_errReadingEnvFile(t *testing.T) {
	ec, stdout, stderr, err := run(t, nil, []string{"doesn't matter"}, nil)

	assert.EqualError(t, err, "readfile: open xenv.env: no such file or directory")
	assert.Zero(t, ec)
	assert.Zero(t, stdout)
	assert.Zero(t, stderr)
}

func TestApp_errNoCmd(t *testing.T) {
	ec, stdout, stderr, err := run(t, nil, []string{}, nil)

	assert.EqualError(t, err, "you are to specify command")
	assert.Zero(t, ec)
	assert.Zero(t, stdout)
	assert.Zero(t, stderr)
}

func TestApp_errInvalidCmd(t *testing.T) {
	ec, stdout, stderr, err := run(t, nil, []string{"NOT_EXISTS.exe"}, []string{"testdata/env.env"})

	assert.EqualError(t, err, `cannot run command: exec: "NOT_EXISTS.exe": executable file not found in $PATH`)
	assert.Zero(t, ec)
	assert.Zero(t, stdout)
	assert.Zero(t, stderr)
}

func TestApp_errNotEmptyListButNoMatches(t *testing.T) {
	ec, stdout, stderr, err := run(t, nil, []string{"placeholder"}, []string{"/dev/null"})

	assert.EqualError(t, err, "no env file found")
	assert.Zero(t, ec)
	assert.Zero(t, stdout)
	assert.Zero(t, stderr)
}

func TestApp_errInvalidFileFormat(t *testing.T) {
	ec, stdout, stderr, err := run(t, nil, []string{"placeholder"}, []string{"testdata/invalid.env"})

	assert.EqualError(t, err, "parser: testdata/invalid.env: unexpected end of file")
	assert.Zero(t, ec)
	assert.Zero(t, stdout)
	assert.Zero(t, stderr)
}

func TestApp_errorCode(t *testing.T) {
	ec, stdout, stderr, err := run(t, nil, []string{"testdata/script-exit-code.sh"}, []string{"testdata/env.env"})

	assert.NoError(t, err)
	assert.Equal(t, 7, ec)
	assert.Zero(t, stdout)
	assert.Zero(t, stderr)
}

func TestApp_wrongFinishing(t *testing.T) {
	ec, stdout, stderr, err := run(t, nil, []string{"testdata/script-wrong-finishing.sh"}, []string{"testdata/env.env"})

	assert.EqualError(t, err, "cannot run command: signal: killed")
	assert.Zero(t, ec)
	assert.Zero(t, stdout)
	assert.Zero(t, stderr)
}
