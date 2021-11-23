package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

// available commands:
const (
	majorCmd   = "major"
	minorCmd   = "minor"
	patchCmd   = "patch"
	currentCmd = "current"
)

// available flags:
var (
	debug  bool
	prefix string
)

func main() {
	flag.Usage = usage
	flag.BoolVar(&debug, "debug", false, "enable debug logs")
	flag.StringVar(&prefix, "prefix", "", "consider only prefixed tags. also, will be used to print the result")
	flag.Parse()

	log.SetFlags(0)
	log.SetOutput(io.Discard)
	log.SetPrefix(os.Args[0] + ": ")
	if debug {
		log.SetOutput(os.Stderr)
	}

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if errors.As(err, new(usageError)) {
			usage()
			os.Exit(2)
		}
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [flags] <command>\n", os.Args[0])

	// print commands in the same style flag.PrintDefaults prints flags.
	fmt.Fprintf(os.Stderr, "\navailable commands:\n")
	for _, cmd := range []string{majorCmd, minorCmd, patchCmd} {
		fmt.Fprintf(os.Stderr, "  %s\n", cmd)
		fmt.Fprintf(os.Stderr, "\tprint the next %s version\n", cmd)
	}
	fmt.Fprintf(os.Stderr, "  %s\n", currentCmd)
	fmt.Fprintf(os.Stderr, "\tprint the current version\n")

	fmt.Fprintf(os.Stderr, "\navailable flags:\n")
	flag.PrintDefaults()
}

// usageError is a wrapper for error to indicate the need for usage printing.
type usageError struct{ err error }

func (e usageError) Error() string { return e.err.Error() }
func (e usageError) Unwrap() error { return e.err }

func run() error {
	if flag.NArg() == 0 {
		return usageError{errors.New("no command has been provided")}
	}

	var printVersion func(version)
	switch cmd := flag.Arg(0); cmd {
	case majorCmd:
		printVersion = func(v version) { fmt.Print(prefix, v.nextMajor()) }
	case minorCmd:
		printVersion = func(v version) { fmt.Print(prefix, v.nextMinor()) }
	case patchCmd:
		printVersion = func(v version) { fmt.Print(prefix, v.nextPatch()) }
	case currentCmd:
		printVersion = func(v version) { fmt.Print(prefix, v) }
	default:
		return usageError{fmt.Errorf("unknown command %q", cmd)}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var buf bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", "tag", "--sort=-version:refname")
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running %q: %w", "git tag", err)
	}

	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, prefix) {
			log.Printf("skipping %q: missing %q prefix", line, prefix)
			continue
		}

		line = strings.TrimPrefix(line, prefix)
		v, err := parseVersion(line)
		if err != nil {
			log.Printf("skipping %q: %s", line, err)
			continue
		}

		printVersion(v)
		break
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading %q output: %w", "git tag", err)
	}

	return nil
}
