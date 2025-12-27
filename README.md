# Golang CLI Flags Scanner

[![GoDoc](https://pkg.go.dev/badge/github.com/bassosimone/flagscanner)](https://pkg.go.dev/github.com/bassosimone/flagscanner) [![Build Status](https://github.com/bassosimone/flagscanner/actions/workflows/go.yml/badge.svg)](https://github.com/bassosimone/flagscanner/actions) [![codecov](https://codecov.io/gh/bassosimone/flagscanner/branch/main/graph/badge.svg)](https://codecov.io/gh/bassosimone/flagscanner)

The `flagscanner` Go package contains a scanner for lexing and
classifying command line arguments. It is a building block that
enables building command-line-flags parsers.

For example:

```Go
import (
	"os"

	"github.com/bassosimone/flagscanner"
)

// Construct a scanner recognizing GNU style options.
scanner := &flagscanner.Scanner{
	Prefixes:  []string{"-", "--"},
	Separator: "--",
}

// Lex the command line arguments using the given prefixes and separator
tokens := scanner.Scan(os.Args[1:])
```

The above example configures GNU style options but we support a
wide variety of command-line-flags styles including Go, dig, Windows,
and traditional Unix. See [example_test.go](example_test.go).

## Installation

To add this package as a dependency to your module:

```sh
go get github.com/bassosimone/flagscanner
```

## Development

To run the tests:
```sh
go test -v .
```

To measure test coverage:
```sh
go test -v -cover .
```

## License

```
SPDX-License-Identifier: GPL-3.0-or-later
```

## History

Adapted from [bassosimone/clip](https://github.com/bassosimone/clip/tree/v0.8.0).
