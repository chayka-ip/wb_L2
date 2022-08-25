package utils

import "unicode"

//IsDigitRune reports whether the rune is a decimal digit.
func IsDigitRune(r rune) bool {
	return unicode.IsDigit(r)
}

// IsSpaceRune reports whether the rune is a space character as defined
// by Unicode's White Space property; in the Latin-1 space
// this is
//	'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
// Other definitions of spacing characters are set by category
// Z and property Pattern_White_Space.
func IsSpaceRune(r rune) bool {
	return unicode.IsSpace(r)
}

//IsBackSlashRune reports whether the rune is backslash symbol
func IsBackSlashRune(r rune) bool {
	return r == 92
}
