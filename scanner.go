// scanner.go - Command line scanner.
// SPDX-License-Identifier: GPL-3.0-or-later

/*
Package flagscanner provides low-level tokenization of command-line arguments.

The [*Scanner.Scan] method breaks command-line arguments into [Token]
based on configurable option prefixes and a separator, allowing higher-level parsers
to implement custom parsing logic on top of the tokenized stream.

# Token Types

[*Scanner.Scan] produces these token types:

 1. [OptionToken]: Options started with any configured prefix (e.g., -v, --verbose, +trace)

 2. [OptionsArgumentsSeparatorToken]: Special separator (e.g., -- to stop parsing)

 3. [PositionalArgumentToken]: Everything else (positional arguments)

# Option Prefixes

The [*Scanner] is configured with the option prefixes to use when tokenizing
command-line arguments. Prefixes are sorted by length (longest first) to ensure
correct tokenization when prefixes overlap (e.g., "-" and "--").

This design allows building parsers for different command-line styles:

 1. GNU-style: "-", "--" (e.g., -v, --verbose)

 2. Dig-style: "-", "--", "+" (e.g., -v, --verbose, +trace)

 3. Windows-style: "/" (e.g., /v, /verbose)

 4. Go-style: "-" (e.g., -v, -verbose)

# Separator

The [*Scanner] can be configured to recognize and emit as a token the separator
to stop parsing options and treat all remaining arguments as positional.

# Example

Given the "--" and "-" option prefixes and the "--" separator, the
following command line arguments:

	--verbose -k4 -- othercommand -v --trace file.txt

produces the following tokens:

 1. [OptionToken] verbose
 2. [OptionToken] -k4
 3. [OptionsArgumentsSeparatorToken] --
 4. [PositionalArgumentToken] othercommand
 5. [PositionalArgumentToken] -v
 6. [PositionalArgumentToken] --trace
 7. [PositionalArgumentToken] file.txt

Note that everything after the separator becomes a positional argument.
*/
package flagscanner

import (
	"sort"
	"strings"
)

// Scanner is a command line scanner.
//
// We check for the separator first. Then for option prefixes
// sorted by length (longest first).
type Scanner struct {
	// Prefixes contains the prefixes delimiting options.
	//
	// If empty, we don't recognize any prefix.
	Prefixes []string

	// Separator contains the separator between options and arguments.
	//
	// If empty, we don't recognize any separator.
	Separator string
}

// Token is a token lexed by [*Scanner.Scan].
type Token interface {
	// Index returns the token index.
	Index() int

	// String returns the string representation of the token.
	String() string
}

// OptionToken is a [Token] containing an option.
type OptionToken struct {
	// Idx is the position in the original command line arguments.
	Idx int

	// Prefix is the scanned prefix.
	Prefix string

	// Name is the parsed name.
	Name string
}

var _ Token = OptionToken{}

// Index implements [Token].
func (tk OptionToken) Index() int {
	return tk.Idx
}

// String implements [Token].
func (tk OptionToken) String() string {
	return tk.Prefix + tk.Name
}

// PositionalArgumentToken is a [Token] containing a positional argument.
type PositionalArgumentToken struct {
	// Idx is the position in the original command line arguments.
	Idx int

	// Value is the parsed value.
	Value string
}

var _ Token = PositionalArgumentToken{}

// Index implements [Token].
func (tk PositionalArgumentToken) Index() int {
	return tk.Idx
}

// String implements [Token].
func (tk PositionalArgumentToken) String() string {
	return tk.Value
}

// OptionsArgumentsSeparatorToken is a [Token] containing the separator between options and arguments.
type OptionsArgumentsSeparatorToken struct {
	// Idx is the position in the original command line arguments.
	Idx int

	// Separator is the parsed separator.
	Separator string
}

var _ Token = OptionsArgumentsSeparatorToken{}

// Index implements [Token].
func (tk OptionsArgumentsSeparatorToken) Index() int {
	return tk.Idx
}

// String implements [Token].
func (tk OptionsArgumentsSeparatorToken) String() string {
	return tk.Separator
}

// Scan scans the command line arguments and returns a list of [Token].
//
// The args MUST NOT include the program name as the first argument.
//
// This method does not mutate the [*Scanner] and is safe to call concurrently.
func (sx *Scanner) Scan(args []string) []Token {
	// Create an empty list of tokens
	tokens := make([]Token, 0, len(args))

	// Create sorted copy of prefixes (longest first)
	prefixes := make([]string, len(sx.Prefixes))
	copy(prefixes, sx.Prefixes)

	// Sort by length descending, then alphabetically for stability
	sort.SliceStable(prefixes, func(i, j int) bool {
		if len(prefixes[i]) == len(prefixes[j]) {
			return prefixes[i] < prefixes[j]
		}
		return len(prefixes[i]) > len(prefixes[j])
	})

	// Cycle through the remaining arguments
loop:
	for idx, arg := range args {
		// Check for separator first
		if sx.Separator != "" && arg == sx.Separator {
			tokens = append(tokens, OptionsArgumentsSeparatorToken{Idx: idx, Separator: arg})
			for tailIdx, tailArg := range args[idx+1:] {
				tokens = append(tokens, PositionalArgumentToken{
					Idx:   idx + 1 + tailIdx,
					Value: tailArg,
				})
			}
			return tokens
		}

		// Then, check for (sorted) prefixes with actual names
		for _, prefix := range prefixes {
			if strings.HasPrefix(arg, prefix) && len(arg) > len(prefix) {
				tokens = append(tokens, OptionToken{Idx: idx, Prefix: prefix, Name: arg[len(prefix):]})
				continue loop
			}
		}

		// Everything else is an argument
		tokens = append(tokens, PositionalArgumentToken{Idx: idx, Value: arg})
	}

	return tokens
}
