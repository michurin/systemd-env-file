package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/michurin/systemd-env-file/sdenv"
)

const skipFileMode = fs.ModeType ^ fs.ModeSymlink

func App(env, args []string, stdout, stderr io.Writer, envFiles []string) (int, error) {
	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	f.SetOutput(stderr)
	f.Usage = func() {
		fmt.Fprintln(f.Output(), "usage: xenv [flags] command [args ...]")
		f.PrintDefaults()
	}
	versionFlag := f.Bool("v", false, "show version")
	helpFlag := f.Bool("h", false, "show help message")
	err := f.Parse(args[1:])
	if err != nil {
		return 2, err //nolint:gomnd
	}

	if helpFlag != nil && *helpFlag {
		f.Usage()
		return 0, nil
	}
	if versionFlag != nil && *versionFlag {
		bi, _ := debug.ReadBuildInfo()
		fmt.Fprintf(f.Output(),
			"version: %s@%s\n\nbuild info:\n%s\n",
			bi.Path,
			bi.Main.Version,
			bi.String())
		return 0, nil
	}
	if f.NArg() < 1 {
		return 0, fmt.Errorf("you are to specify command")
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
		return 0, fmt.Errorf("no env file found")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return 0, fmt.Errorf("readfile: %w", err)
	}

	pairs, err := sdenv.Parser(data)
	if err != nil {
		return 0, fmt.Errorf("parser: %s: %w", file, err)
	}

	c := sdenv.NewCollectsion()
	c.PushStd(env)
	c.Push(pairs)

	cmd := exec.Command(f.Arg(0), f.Args()[1:]...) //nolint:gosec
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = c.CollectionStd()
	err = cmd.Run()
	if err != nil {
		e := new(exec.ExitError)
		if errors.As(err, &e) {
			ec := e.ExitCode()
			if ec >= 0 { // normal exit: not signaled, not coredumped...
				return ec, nil
			}
		}
		return 0, fmt.Errorf("cannot run command: %w", err)
	}
	return 0, nil
}
