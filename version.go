package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var errInvalidFormat = errors.New("invalid semantic version format")

// https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
var semver = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

// version is a parsed semantic version.
type version struct {
	major int
	minor int
	patch int
}

// parseVersion parses a semantic version from the provided string.
// The string must meet the semantic versioning specification requirements,
// see https://semver.org/#semantic-versioning-specification-semver for details.
func parseVersion(s string) (*version, error) {
	const minMatches = 1 + 3 // 1 full match + at least 3 subexpression matches.

	matches := semver.FindStringSubmatch(s)
	if len(matches) < minMatches {
		return nil, errInvalidFormat
	}

	// the first match must be equal to the provided string.
	if matches[0] != s {
		return nil, errInvalidFormat
	}

	var numbers []int
	for _, match := range matches[1:minMatches] {
		n, err := strconv.Atoi(match)
		if err != nil {
			// unreachable due to the semver regexp.
			panic(err)
		}
		numbers = append(numbers, n)
	}

	return &version{
		major: numbers[0],
		minor: numbers[1],
		patch: numbers[2],
	}, nil
}

const format = "%d.%d.%d"

func (v *version) String() string    { return fmt.Sprintf(format, v.major, v.minor, v.patch) }
func (v *version) nextMajor() string { return fmt.Sprintf(format, v.major+1, 0, 0) }
func (v *version) nextMinor() string { return fmt.Sprintf(format, v.major, v.minor+1, 0) }
func (v *version) nextPatch() string { return fmt.Sprintf(format, v.major, v.minor, v.patch+1) }
