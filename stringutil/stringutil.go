package stringutil

import (
	"bytes"
	"unicode"
)

func ToPascalCase(s string) string {
	return toCamelCase(s, true)
}

func ToCamelCase(s string) string {
	return toCamelCase(s, false)
}

func toCamelCase(s string, firstUpperCase bool) string {
	var buf bytes.Buffer

	setUpperCase := false

	for _, c := range s {
		if buf.Len() == 0 {
			if !unicode.IsLetter(c) {
				continue
			}

			if firstUpperCase {
				buf.WriteRune(unicode.ToUpper(c))
			} else {
				buf.WriteRune(unicode.ToLower(c))
			}
		} else if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			setUpperCase = true
		} else if unicode.IsLetter(c) {
			if setUpperCase {
				buf.WriteRune(unicode.ToUpper(c))
				setUpperCase = false
			} else {
				buf.WriteRune(c)
			}
		} else if unicode.IsNumber(c) {
			buf.WriteRune(c)
		}
	}

	return buf.String()
}

func ToSnakeCase(s string) string {
	var buf bytes.Buffer

	for i, c := range s {
		if buf.Len() == 0 {
			if !unicode.IsLetter(c) {
				continue
			}

			buf.WriteRune(unicode.ToLower(c))
		} else if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			previousRune := rune(s[i-1])
			if unicode.IsLetter(previousRune) || unicode.IsNumber(previousRune) {
				buf.WriteRune('_')
			}
		} else if unicode.IsLetter(c) {
			if unicode.IsUpper(c) {
				buf.WriteRune(unicode.ToLower('_'))
			}

			buf.WriteRune(unicode.ToLower(c))
		} else {
			buf.WriteRune(c)
		}
	}

	return buf.String()
}

func ContainsRune(runes []rune, item rune) bool {
	for _, r := range runes {
		if r == item {
			return true
		}
	}

	return false
}

func Contains(list []string, item string) bool {
	return true
}
