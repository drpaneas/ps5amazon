package util

import (
	"regexp"
	"strings"
	"unicode/utf8"

	strip "github.com/grokify/html-strip-tags-go"
)

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func removeWhitespaceBeginning(s string) string {
	for {
		if strings.Contains(s[0:1], " ") {
			s = trimFirstRune(s)
		} else {
			break
		}
	}
	return s
}

func removeDoubleWhitespace(s string) string {
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")
	return s
}

// ApplyTextFormat clears output
func ApplyTextFormat(s string) string {
	s = strings.Replace(s, "\n", "", -1) // remove line breaks
	s = strip.StripTags(s)               // remove html tags if any
	s = removeWhitespaceBeginning(s)     // remove trailling whitespace at the beginning
	s = removeDoubleWhitespace(s)
	return s
}
