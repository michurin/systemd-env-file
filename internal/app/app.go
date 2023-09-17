package app

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"

	"github.com/michurin/systemd-env-file/sdenv"
)

const skipFileMode = fs.ModeType ^ fs.ModeSymlink

func App(env, args []string, stdout, stderr io.Writer, envFiles []string) error {
	if len(args) < 1 {
		return fmt.Errorf("you are to specify command")
	}

	file := ""
	if len(envFiles) == 0 {
		file = "xenv.env"
	} else {
		for _, f := range envFiles {
			fi, err := os.Stat(f)
			if err != nil {
				continue
			}
			if fi.Mode()&skipFileMode != 0 {
				continue
			}
			file = f
		}
	}
	if file == "" {
		return fmt.Errorf("no env file found")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("readfile: %w", err)
	}

	pairs, err := sdenv.Parser(data)
	if err != nil {
		return fmt.Errorf("parser: %s: %w", file, err)
	}

	c := sdenv.NewCollectsion()
	c.PushStd(env)
	c.Push(pairs)

	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = c.CollectionStd()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("cannot run command: %w", err)
	}
	return nil
}
