/*
Xenv loads environment variables from file and runs given process.

The format of the file is exactly the same as systemd [EnvironmentFile] command sourced.

Xenv tries to load xenv.env file or first exists file mentioned in XENV variable.

Usage:

	xenv [flags] command [args...]

The flags are:

	-v
		Show version

	-h
		Show help message

[EnvironmentFile]: https://www.freedesktop.org/software/systemd/man/systemd.exec.html#EnvironmentFile=
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/michurin/systemd-env-file/internal/app"
	"github.com/michurin/systemd-env-file/internal/args"
)

//nolint:gochecknoglobals // TODO move flags to internal/app or even internal/args.
var (
	versionFlag = flag.Bool("v", false, "show version")
	helpFlag    = flag.Bool("h", false, "show help message")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: xenv [flags] command [args ...]")
		flag.PrintDefaults()
	}
	flag.Parse()
	if helpFlag != nil && *helpFlag {
		flag.Usage()
		os.Exit(0)
	}
	if versionFlag != nil && *versionFlag {
		bi, _ := debug.ReadBuildInfo()
		fmt.Fprintln(os.Stderr, "build info:\n"+bi.String())
		os.Exit(0)
	}
	exitCode, err := app.App(os.Environ(), flag.Args(), os.Stdout, os.Stderr, args.Split(os.Getenv("XENV")))
	if err != nil {
		log.Println("Error:", err)
		os.Exit(127) //nolint:gomnd
	}
	os.Exit(exitCode)
}
