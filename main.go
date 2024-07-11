package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/ValerySidorin/depman/internal/gui"
	"github.com/ValerySidorin/depman/internal/index"
)

var ErrInvalidFlags = errors.New("invalid flags")

func main() {
	if len(os.Args) < 4 {
		log.Fatal("invalid arguments")
	}
	cmd := os.Args[1 : len(os.Args)-1]
	if err := validateGoCmd(cmd); err != nil {
		log.Fatal(err)
	}

	query := os.Args[len(os.Args)-1]

	index := &index.DepsDev{}
	packages, err := index.Search(query)
	if err != nil {
		log.Fatal(err)
	}

	execCmd := gui.Handle(strings.Join(cmd, " "), packages)
	if execCmd == "" {
		os.Exit(0)
	}
	osCmd := exec.Command(cmd[0], strings.Split(execCmd, " ")[1:]...)
	stdout, err := osCmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := osCmd.Start(); err != nil {
		log.Fatal(err)
	}

	_, _ = io.Copy(os.Stdout, stdout)
}

func validateGoCmd(args []string) error {
	if args[0] != "go" ||
		(args[1] != "get" && args[1] != "install") {
		return errors.New("invalid go command")
	}

	flags := args[2:]
	var validFlags []string

	if args[1] == "get" {
		validFlags = []string{"-u", "-t", "-v"}
	} else if args[1] == "install" {
		validFlags = []string{}
	}

	if err := validateFlags(flags, validFlags); err != nil {
		return fmt.Errorf("validate cmd flags: %w", err)
	}

	return nil
}

func validateFlags(flags []string, validFlags []string) error {
	if len(flags) == 0 {
		return nil
	}

	if len(flags) > len(validFlags) {
		return fmt.Errorf("invalid flag len: %w", ErrInvalidFlags)
	}

	mflags := make(map[string]struct{})
	for _, f := range flags {
		mflags[f] = struct{}{}
	}

	if len(flags) != len(mflags) {
		return fmt.Errorf("invalid flag set len: %w", ErrInvalidFlags)
	}

	for f := range mflags {
		if !slices.Contains(validFlags, f) {
			return fmt.Errorf("invalid flag found: %w", ErrInvalidFlags)
		}
	}

	return nil
}
