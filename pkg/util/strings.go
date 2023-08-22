package util

import "strings"

// TitleCase converts a string to title case (i.e. "hello world" -> "Hello World").
func TitleCase(s string) string {
	if len(s) == 0 {
		return s
	}

	uppercase := true
	for i, c := range s {
		if c == ' ' {
			uppercase = true
			continue
		}

		if uppercase {
			s = s[:i] + strings.ToUpper(string(c)) + s[i+1:]
			uppercase = false
		}
	}

	return s
}

// ReplaceChars replaces characters in a string according to the given map.
func ReplaceChars(s string, replaceMap map[rune]rune) string {
	runes := []rune(s)
	for i, c := range runes {
		if r, ok := replaceMap[c]; ok {
			runes[i] = r
		} else {
			runes[i] = c
		}
	}

	return string(runes)
}
