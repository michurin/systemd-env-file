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
	"log"
	"os"

	"github.com/michurin/systemd-env-file/internal/app"
	"github.com/michurin/systemd-env-file/internal/args"
)

func main() {
	exitCode, err := app.App(os.Environ(), os.Args, os.Stdout, os.Stderr, args.Split(os.Getenv("XENV")))
	if err != nil {
		log.Println("Error:", err)
		os.Exit(127) //nolint:gomnd
	}
	os.Exit(exitCode)
}
