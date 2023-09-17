package main

import (
	"log"
	"os"

	"github.com/michurin/systemd-env-file/internal/app"
	"github.com/michurin/systemd-env-file/internal/args"
)

func main() {
	exitCode, err := app.App(os.Environ(), os.Args[1:], os.Stdout, os.Stderr, args.Split(os.Getenv("XENV")))
	if err != nil {
		log.Println("Error:", err)
		os.Exit(127) //nolint: gomnd
	}
	os.Exit(exitCode)
}
