package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTitleCase(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "empty string",
			in:   "",
			out:  "",
		},
		{
			name: "single word",
			in:   "hello",
			out:  "Hello",
		},
		{
			name: "two words",
			in:   "hello world",
			out:  "Hello World",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.out, TitleCase(test.in))
		})
	}
}

func TestReplaceChars(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		replaceMap map[rune]rune
		out        string
	}{
		{
			name:       "empty string",
			in:         "",
			replaceMap: map[rune]rune{},
			out:        "",
		},
		{
			name:       "no replacements",
			in:         "hello world",
			replaceMap: map[rune]rune{'a': 'b'},
			out:        "hello world",
		},
		{
			name:       "dashes and underscores",
			in:         "part1-part2_part3",
			replaceMap: map[rune]rune{'-': ' ', '_': ' '},
			out:        "part1 part2 part3",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.out, ReplaceChars(test.in, test.replaceMap))
		})
	}
}
