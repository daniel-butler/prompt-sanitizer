package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/openclaw/prompt-sanitizer/pkg/wrapper"
)

// Version is set at build time via ldflags
var Version = "dev"

func main() {
	if err := run(os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(stderr)

	source := fs.String("source", "Unknown", "Source label for the content")
	filePath := fs.String("file", "", "File to wrap (if not reading from stdin)")
	showVersion := fs.Bool("version", false, "Print version and exit")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	if *showVersion {
		fmt.Fprintln(stdout, Version)
		return nil
	}

	var content string
	var err error

	// Check if we have remaining args (command execution mode)
	remainingArgs := fs.Args()
	if len(remainingArgs) > 0 {
		// Command execution mode
		content, err = executeCommand(remainingArgs)
		if err != nil {
			return fmt.Errorf("executing command: %w", err)
		}
	} else if *filePath != "" {
		// File mode
		content, err = readFile(*filePath)
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}
	} else {
		// Stdin mode
		content, err = readFromReader(stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}
	}

	// Wrap and output
	wrapped := wrapper.WrapContent(content, *source)
	fmt.Fprintln(stdout, wrapped)
	return nil
}

func readFromReader(r io.Reader) (string, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func readFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func executeCommand(args []string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}
	return string(output), nil
}
