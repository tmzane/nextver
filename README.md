# nextver

[![ci-img]][ci]
[![docs-img]][docs]
[![report-img]][report]
[![codecov-img]][codecov]
[![license-img]][license]
[![release-img]][release]

A dead simple CLI tool that prints the next semantic version based on the last tag of your git repository.

## Install

### Go

```bash
go install github.com/junk1tm/nextver@latest
```

### Brew

```bash
brew install junk1tm/tap/nextver
```

### Manual

Download a prebuilt binary from the [Releases][release] page.

## Usage

`nextver` supports `major`, `minor` and `patch` commands:

```bash
> git tag
1.2.3
> nextver major
2.0.0
> nextver minor
1.3.0
> nextver patch
1.2.4
```

`nextver` expects an initiated `git` repository with tags following the [semantic versioning specification][semver],
invalid tags are ignored:

```bash
> git tag
a.b.c
1.2.3.4
> nextver -verbose patch
nextver: skipping "a.b.c": invalid semantic version format
nextver: skipping "1.2.3.4": invalid semantic version format
nextver: no valid semantic version has been found
```

`nextver` is designed to be used with `git tag` command:

```bash
git tag $(nextver patch)
```

so the only thing it writes to `stdout` is the next version, everything else, including usage, errors and additional
information, goes to `stderr`. If anything goes wrong, e.g., no `git` repository is found, `nextver` won't
touch `stdout` at all, so the command above will be resolved to `git tag`, which simply prints the existing `git` tags.
Therefore, accidentally creating an invalid tag is unlikely.

### Prefix

The `-prefix` flag allows defining a prefix for tags. If set, `nextver` ignores tags without this prefix and prints the
next version prefixed.

A common practice is to prefix `semver` tags with `v`:

```bash
> git tag
v1.2.3
> nextver -prefix=v patch
v1.2.4
```

Another use case is working with [go submodules][submodule], which requires submodule's tags to be prefixed:

```bash
> tree mainmodule
mainmodule
├── go.mod
└── submodule
    └── go.mod

1 directory, 2 files
> cd mainmodule && git tag
0.1.0
submodule/0.1.1
> nextver -prefix=submodule/ patch
submodule/0.1.2
```

### Help

```bash
> nextver -help
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
[semver]: https://semver.org/#semantic-versioning-specification-semver
[submodule]: https://github.com/go-modules-by-example/index/blob/master/009_submodules/README.md
