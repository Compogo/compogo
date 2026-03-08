package tools

import (
	"math/rand"
)

// LetterRunes contains a set of characters for generating random strings (a-z, A-Z)
var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandString returns a random string of the specified length consisting of Latin letters.
// The length of the string is determined by the length parameter.
func RandString(length uint32) string {
	return RandStringFromRunes(length, LetterRunes...)
}

// RandStringFromRunes returns a random string of the specified length composed of the provided runes.
// If runes is empty, an empty string is returned.
func RandStringFromRunes(length uint32, runes ...rune) string {
	name := make([]rune, length)

	for i := range name {
		name[i] = runes[rand.Intn(len(runes))]
	}

	return string(name)
}
