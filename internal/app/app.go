package app

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"runtime/debug"
	"syscall"

	"github.com/michurin/systemd-env-file/sdenv"
)

const skipFileMode = fs.ModeType ^ fs.ModeSymlink

func App(env, args []string, stderr io.Writer, envFiles []string) (int, error) {
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

	lp, err := exec.LookPath(f.Arg(0))
	if err != nil {
		return 0, fmt.Errorf("lookup executable: %w", err)
	}
	err = syscall.Exec(lp, f.Args(), c.CollectionStd())
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}
	return 0, nil // how can we get here?
}
