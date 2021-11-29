# nextver

[![ci-img]][ci]
[![docs-img]][docs]
[![report-img]][report]
[![codecov-img]][codecov]
[![license-img]][license]
[![release-img]][release]

A dead simple CLI tool that prints the next semantic version based on the last tag of your git repository.

## Install

```bash
go install github.com/junk1tm/nextver@latest
```

## Usage

```bash
usage: nextver [flags] <command>

available commands:
  major
        print the next major version
  minor
        print the next minor version
  patch
        print the next patch version

available flags:
  -p string
        shorthand for -prefix
  -prefix string
        consider only prefixed tags (also, will be used to print the result)
  -v    
        shorthand for -verbose
  -verbose
        print additional information to stderr
  -version
        print the app version
```

[ci]: https://github.com/junk1tm/nextver/actions/workflows/go.yml
[ci-img]: https://github.com/junk1tm/nextver/actions/workflows/go.yml/badge.svg
[docs]: https://pkg.go.dev/github.com/junk1tm/nextver
[docs-img]: https://pkg.go.dev/badge/github.com/junk1tm/nextver.svg
[report]: https://goreportcard.com/report/github.com/junk1tm/nextver
[report-img]: https://goreportcard.com/badge/github.com/junk1tm/nextver
[codecov]: https://codecov.io/gh/junk1tm/nextver
[codecov-img]: https://codecov.io/gh/junk1tm/nextver/branch/main/graph/badge.svg
[license]: https://github.com/junk1tm/nextver/blob/main/LICENSE
[license-img]: https://img.shields.io/github/license/junk1tm/nextver
[release]: https://github.com/junk1tm/nextver/releases
[release-img]: https://img.shields.io/github/v/release/junk1tm/nextver
