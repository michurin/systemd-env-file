package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/michurin/systemd-env-file/sdenv"
)

func app(args []string, stdout, stderr io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("you are to specify command")
	}
	env, err := sdenv.Environ(os.Environ(), "xenv.env")
	if err != nil {
		return fmt.Errorf("cannot open env file: %w", err)
	}
	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = env
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("cannot run command: %w", err)
	}
	return nil
}

func main() {
	err := app(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		log.Println("Error:", err)
	}
}
