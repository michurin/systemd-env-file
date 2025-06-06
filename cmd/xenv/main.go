// The xenv command loads environment variables from file and runs given process.
//
// The format of the file is exactly the same as systemd [EnvironmentFile] command sourced.
//
// Xenv tries to load xenv.env file or first regular file mentioned in XENV variable if defined.
//
// Usage:
//
//	xenv [flags] command [args...]
//
// Flags:
//
//	-v  Show version
//	-h  Show help message
//
// [EnvironmentFile]: https://www.freedesktop.org/software/systemd/man/systemd.exec.html#EnvironmentFile=
package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/michurin/systemd-env-file/sdenv"
)

// One hundred lines of ugly code here:
// global from flag package, os.Exit (including flag's internals), syscall.Exec, unreachable code...
// I won't overcomplicate it, gave up making this code testable and
// wrote integration tests for compiled binary though.

func main() {
	opts()
	execute(buildEnv(lookupEnvFile()))
}

func opts() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: xenv [flags] command [args ...]")
		fmt.Fprintln(flag.CommandLine.Output(), "  -h\tshow help message") // -h supported under the hood, however doesn't show up in PrintDefaults
		flag.PrintDefaults()
	}
	versionFlag := flag.Bool("v", false, "show version")
	debugginFlag := flag.Bool("d", false, "debug")
	flag.Parse() // we do not need to check error cause ExitOnError

	if versionFlag != nil && *versionFlag {
		bi, _ := debug.ReadBuildInfo()
		fmt.Fprintf(flag.CommandLine.Output(),
			"version: %s@%s\n\nbuild info:\n%s\n",
			bi.Path,
			bi.Main.Version,
			bi.String())
		os.Exit(0)
	}
	if debugginFlag != nil && *debugginFlag {
		log.SetOutput(os.Stderr)
		log.SetFlags(0)
		log.SetPrefix("DEBUG: ")
	} else {
		log.SetOutput(io.Discard) // just mute default logger
	}
	if flag.NArg() < 1 {
		exit("Error: you have to specify command\n")
	}
}

func lookupEnvFile() string {
	const skipFileMode = fs.ModeType ^ fs.ModeSymlink
	s := os.Getenv("XENV")
	if len(s) == 0 {
		log.Println("No $XENV, just taking xenv.env from current directory")
		return "xenv.env"
	}
	log.Println("Considering $XENV:", s)
	for _, f := range strings.Split(s, ":") {
		log.Println("Considering part:", f)
		fi, err := os.Stat(f)
		if err != nil {
			log.Println("Skipping part due to error:", err.Error())
			continue
		}
		mode := fi.Mode()
		if mode&skipFileMode != 0 {
			log.Printf("Skipping part due to mode: %[1]s, skipping reason: %[2]s", mode, mode&skipFileMode)
			continue
		}
		log.Printf("File is taken: %s", f)
		return f
	}
	exit("No env file found\n")
	return "" // unreachable
}

func buildEnv(envFile string) []string {
	data, err := os.ReadFile(envFile)
	if err != nil {
		exit("Cannot read file %s: %v\n", envFile, err)
	}

	pairs, err := sdenv.Parser(data)
	if err != nil {
		exit("Cannot parse: %s: %v\n", envFile, err)
	}

	c := sdenv.NewCollectsion()
	c.PushStd(os.Environ())
	c.Push(pairs)
	return c.CollectionStd()
}

func execute(env []string) {
	cmdArgs := flag.Args()
	lp, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		exit("Lookup executable: %s: %v\n", cmdArgs[0], err)
	}
	err = syscall.Exec(lp, cmdArgs, env)
	if err != nil {
		exit("Exec: %v\n", err)
	}
	exit("Exec: exited without error (looks like error)\n")
}

func exit(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(2) //nolint:mnd
}
