// scanner_test.go - Tests for command line scanner.
// SPDX-License-Identifier: GPL-3.0-or-later

package flagscanner

import "testing"

// This test ensures that the [Token.Index] method is working as
// intended for each available token type.
func TestTokenIndex(t *testing.T) {
	tests := []struct {
		name     string
		token    Token
		expected int
	}{
		{
			name:     "OptionToken",
			token:    OptionToken{Idx: 1},
			expected: 1,
		},
		{
			name:     "ArgumentToken",
			token:    PositionalArgumentToken{Idx: 1},
			expected: 1,
		},
		{
			name:     "SeparatorToken",
			token:    OptionsArgumentsSeparatorToken{Idx: 1},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.token.Index()
			if got != tt.expected {
				t.Errorf("Token.Index() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// This test ensures that [Token.String] round trips the original
// token value for each available token type.
func TestTokenString(t *testing.T) {
	tests := []struct {
		name     string
		token    Token
		expected string
	}{
		{
			name:     "OptionToken with single dash",
			token:    OptionToken{Prefix: "-", Name: "v"},
			expected: "-v",
		},
		{
			name:     "OptionToken with double dash",
			token:    OptionToken{Prefix: "--", Name: "verbose"},
			expected: "--verbose",
		},
		{
			name:     "ArgumentToken",
			token:    PositionalArgumentToken{Value: "file.txt"},
			expected: "file.txt",
		},
		{
			name:     "SeparatorToken",
			token:    OptionsArgumentsSeparatorToken{Separator: "--"},
			expected: "--",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.token.String()
			if got != tt.expected {
				t.Errorf("Token.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// This test ensures that we can use `-` to indicate stdout and it is
// recognized as a positional argument rather than as a flag.
func TestScannerZeroLengthOption(t *testing.T) {
	scanner := &Scanner{
		Prefixes:  []string{"-"},
		Separator: "",
	}

	args := []string{"prog", "-"}
	tokens := scanner.Scan(args[1:])
	if len(tokens) != 1 {
		t.Errorf("Expected 1 token, got %d", len(tokens))
	}

	if _, ok := tokens[0].(PositionalArgumentToken); !ok {
		t.Errorf("Expected PositionalArgumentToken, got %T", tokens[0])
	}
}

// This test ensures that the separator stops option parsing and the
// remaining arguments are treated as positional.
func TestScannerSeparatorStopsParsing(t *testing.T) {
	scanner := &Scanner{
		Prefixes:  []string{"-", "--"},
		Separator: "--",
	}

	args := []string{"prog", "--", "-v", "--trace", "file.txt"}
	tokens := scanner.Scan(args[1:])

	if len(tokens) != 4 {
		t.Fatalf("Expected 4 tokens, got %d", len(tokens))
	}

	if _, ok := tokens[0].(OptionsArgumentsSeparatorToken); !ok {
		t.Fatalf("Expected OptionsArgumentsSeparatorToken, got %T", tokens[0])
	}

	for idx := 1; idx < len(tokens); idx++ {
		if _, ok := tokens[idx].(PositionalArgumentToken); !ok {
			t.Errorf("Expected PositionalArgumentToken, got %T", tokens[idx])
		}
	}
}
