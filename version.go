package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	errWrongNumbersCount = errors.New("semver must have exactly 3 numbers")
	errNegativeNumber    = errors.New("all semver numbers must be non-negative")
)

// version is a parsed semantic version.
type version struct {
	major int
	minor int
	patch int
}

// parseVersion parses a semantic version from the provided string.
// The string must take the form "x.y.z", where x, y and z are non-negative integers.
// No "v" or any other prefix is expected.
// See https://semver.org/#semantic-versioning-specification-semver for details.
func parseVersion(s string) (version, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return version{}, errWrongNumbersCount
	}

	var numbs []int
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return version{}, fmt.Errorf("parsing number: %w", err)
		}
		if n < 0 {
			return version{}, errNegativeNumber
		}
		numbs = append(numbs, n)
	}

	return version{
		major: numbs[0],
		minor: numbs[1],
		patch: numbs[2],
	}, nil
}

const format = "%d.%d.%d"

func (v version) String() string    { return fmt.Sprintf(format, v.major, v.minor, v.patch) }
func (v version) nextMajor() string { return fmt.Sprintf(format, v.major+1, 0, 0) }
func (v version) nextMinor() string { return fmt.Sprintf(format, v.major, v.minor+1, 0) }
func (v version) nextPatch() string { return fmt.Sprintf(format, v.major, v.minor, v.patch+1) }
