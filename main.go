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

// appVersion is injected at build time.
var appVersion = "dev"

func main() {
	var prefix string
	flag.StringVar(&prefix, "p", "", "shorthand for -prefix")
	flag.StringVar(&prefix, "prefix", "", "consider only prefixed tags (also, will be used to print the result)")

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "\nshorthand for -verbose")
	flag.BoolVar(&verbose, "verbose", false, "print additional information to stderr")

	var version bool
	flag.BoolVar(&version, "version", false, "print the app version")

	flag.Usage = usage
	flag.Parse()

	if version {
		fmt.Fprintf(os.Stderr, "%s version %s\n", os.Args[0], appVersion)
		os.Exit(0)
	}

	log.SetFlags(0)
	log.SetOutput(io.Discard)
	log.SetPrefix(os.Args[0] + ": ")
	if verbose {
		log.SetOutput(os.Stderr)
	}

	if err := run(prefix); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		if errors.As(err, new(usageError)) {
			usage()
			os.Exit(2)
		}
		os.Exit(1)
	}
}

const (
	majorCmd = "major"
	minorCmd = "minor"
	patchCmd = "patch"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [flags] <command>\n", os.Args[0])

	// print commands in the same style flag.PrintDefaults prints flags.
	fmt.Fprintf(os.Stderr, "\navailable commands:\n")
	for _, cmd := range []string{majorCmd, minorCmd, patchCmd} {
		fmt.Fprintf(os.Stderr, "  %s\n", cmd)
		fmt.Fprintf(os.Stderr, "\tprint the next %s version\n", cmd)
	}

	fmt.Fprintf(os.Stderr, "\navailable flags:\n")
	flag.PrintDefaults()
}

func run(prefix string) error {
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
	default:
		return usageError{fmt.Errorf("unknown command %q", cmd)}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", "tag", "--sort=-version:refname")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() != 0 {
			return fmt.Errorf("running %q: %s", "git tag", strings.TrimSpace(stderr.String()))
		}
		return fmt.Errorf("running %q: %w", "git tag", err)
	}

	var current version

	scanner := bufio.NewScanner(&stdout)
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

		current = v
		break
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading %q output: %w", "git tag", err)
	}

	if current.isZero() {
		return fmt.Errorf("no version in the form %q has been found", prefix+"x.y.z")
	}

	log.Printf("the current version is %q", prefix+current.String())
	printVersion(current)

	return nil
}

// usageError is a wrapper for error to indicate the need to print usage.
type usageError struct{ err error }

func (e usageError) Error() string { return e.err.Error() }
func (e usageError) Unwrap() error { return e.err }
