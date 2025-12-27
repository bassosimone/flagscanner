// example_test.go - Scanner example tests
// SPDX-License-Identifier: GPL-3.0-or-later

package flagscanner_test

import (
	"fmt"

	"github.com/bassosimone/flagscanner"
)

// ExampleScanner_dig demonstrates dig command-line parsing style.
//
// Dig style:
//
//   - Traditional short options with single dash: -v, -f
//
//   - Long options with double dash: --verbose, --file
//
//   - Plus options for dig-specific features: +trace, +short, +noall
//
//   - Plus options are treated as long options (no bundling)
//
//   - Options with arguments: -f file, +timeout=5, --port=53
//
//   - Double dash separator: -- stops option parsing
//
//   - Mixed prefix styles for different option categories
func ExampleScanner_dig() {
	s := &flagscanner.Scanner{
		Prefixes:  []string{"-", "--", "+"}, // Single dash, double dash, and plus prefixes
		Separator: "--",                     // Only double dash separator supported
	}

	args := []string{
		"program", "-v", "+trace", "--verbose", "+short=yes",
		"-f", "config", "--", "remaining", "-args",
	}

	tokens := s.Scan(args[1:])
	for _, token := range tokens {
		fmt.Printf("%#v\n", token)
	}

	// Output:
	// flagscanner.OptionToken{Idx:0, Prefix:"-", Name:"v"}
	// flagscanner.OptionToken{Idx:1, Prefix:"+", Name:"trace"}
	// flagscanner.OptionToken{Idx:2, Prefix:"--", Name:"verbose"}
	// flagscanner.OptionToken{Idx:3, Prefix:"+", Name:"short=yes"}
	// flagscanner.OptionToken{Idx:4, Prefix:"-", Name:"f"}
	// flagscanner.PositionalArgumentToken{Idx:5, Value:"config"}
	// flagscanner.OptionsArgumentsSeparatorToken{Idx:6, Separator:"--"}
	// flagscanner.PositionalArgumentToken{Idx:7, Value:"remaining"}
	// flagscanner.PositionalArgumentToken{Idx:8, Value:"-args"}
}

// ExampleScanner_gnu demonstrates GNU command-line parsing.
//
// GNU style:
//
//   - Short options with single dash: -v, -f
//
//   - Long options with double dash: --verbose, --file
//
//   - Short options can be bundled: -vf equivalent to -v -f
//
//   - Options with arguments: -f file, -ffile, --file=name, --file name
//
//   - Double dash separator: -- stops option parsing
//
//   - Argument permutation (reordering) typically supported at parser level
func ExampleScanner_gnu() {
	s := &flagscanner.Scanner{
		Prefixes:  []string{"-", "--"}, // Single and double dash prefixes
		Separator: "--",
	}

	args := []string{"program", "-v", "--file=config.txt", "-abc", "--", "--an-option", "input.txt"}

	tokens := s.Scan(args[1:])
	for _, token := range tokens {
		fmt.Printf("%#v\n", token)
	}

	// Output:
	// flagscanner.OptionToken{Idx:0, Prefix:"-", Name:"v"}
	// flagscanner.OptionToken{Idx:1, Prefix:"--", Name:"file=config.txt"}
	// flagscanner.OptionToken{Idx:2, Prefix:"-", Name:"abc"}
	// flagscanner.OptionsArgumentsSeparatorToken{Idx:3, Separator:"--"}
	// flagscanner.PositionalArgumentToken{Idx:4, Value:"--an-option"}
	// flagscanner.PositionalArgumentToken{Idx:5, Value:"input.txt"}
}

// ExampleScanner_go demonstrates Go command-line parsing style.
//
// Go style:
//
//   - Short and long options with single dash: -v, -f, -verbose, -file
//
//   - No short option bundling: -vf is treated as single option "vf"
//
//   - Options with arguments: -file=name, -file name
//
//   - Double dash separator: -- stops option parsing
//
//   - Simple and consistent: all options use single dash prefix
func ExampleScanner_go() {
	s := &flagscanner.Scanner{
		Prefixes:  []string{"-"}, // Go uses single dash for all options
		Separator: "--",
	}

	args := []string{"program", "-v", "-file=config.txt", "-verbose", "-debug", "input.txt", "--", "extra"}

	tokens := s.Scan(args[1:])
	for _, token := range tokens {
		fmt.Printf("%#v\n", token)
	}

	// Output:
	// flagscanner.OptionToken{Idx:0, Prefix:"-", Name:"v"}
	// flagscanner.OptionToken{Idx:1, Prefix:"-", Name:"file=config.txt"}
	// flagscanner.OptionToken{Idx:2, Prefix:"-", Name:"verbose"}
	// flagscanner.OptionToken{Idx:3, Prefix:"-", Name:"debug"}
	// flagscanner.PositionalArgumentToken{Idx:4, Value:"input.txt"}
	// flagscanner.OptionsArgumentsSeparatorToken{Idx:5, Separator:"--"}
	// flagscanner.PositionalArgumentToken{Idx:6, Value:"extra"}
}

// ExampleScanner_unix demonstrates traditional UNIX command-line parsing.
//
// Traditional UNIX style:
//
//   - Only single-dash short options: -v, -f
//
//   - No long options (--verbose not supported)
//
//   - Short options can be bundled: -vf equivalent to -v -f
//
//   - Options with arguments: -f file or -ffile
//
//   - No special separators (-- not recognized)
func ExampleScanner_unix() {
	s := &flagscanner.Scanner{
		Prefixes:  []string{"-"}, // Only single-dash options in traditional UNIX
		Separator: "",            // No separators in traditional UNIX
	}

	args := []string{"program", "-v", "-f", "file.txt", "-abc", "input.txt"}

	tokens := s.Scan(args[1:])
	for _, token := range tokens {
		fmt.Printf("%#v\n", token)
	}

	// Output:
	// flagscanner.OptionToken{Idx:0, Prefix:"-", Name:"v"}
	// flagscanner.OptionToken{Idx:1, Prefix:"-", Name:"f"}
	// flagscanner.PositionalArgumentToken{Idx:2, Value:"file.txt"}
	// flagscanner.OptionToken{Idx:3, Prefix:"-", Name:"abc"}
	// flagscanner.PositionalArgumentToken{Idx:4, Value:"input.txt"}
}
