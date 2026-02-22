package fwew_lib

import (
	"strings"
)

/* To help deduce phonemes */
var romanization = map[string]string{
	// Vowels
	"a": "a", "i": "i", "ɪ": "ì",
	"o": "o", "ɛ": "e", "u": "u",
	"æ": "ä",
	// Diphthongs
	"aw": "aw", "ɛj": "ey",
	"aj": "ay", "ɛw": "ew",
	// Psuedovowels
	"ṛ": "rr", "ḷ": "ll",
	// Consonants
	"t": "t", "p": "p", "ʔ": "'",
	"n": "n", "k": "k", "l": "l",
	"s": "s", "ɾ": "r", "j": "y",
	"t͡s": "ts", "t'": "tx", "m": "m",
	"v": "v", "w": "w", "h": "h",
	"ŋ": "ng", "z": "z", "k'": "kx",
	"p'": "px", "f": "f", "r": "r",
	// Reef dialect
	"b": "px", "d": "tx", "g": "kx",
	"ʃ": "sy", "tʃ": "tsy", "ʊ": "ù",
	// mistakes and rarities
	"ʒ": "tsy", "": "", " ": "",
}
var romanization2 = map[string]string{
	// Vowels
	"a": "a", "i": "i", "ɪ": "ì",
	"o": "o", "ɛ": "e", "u": "u",
	"æ": "ä", "õ": "õ", //võvä' only
	// Diphthongs
	"aw": "aw", "ɛj": "ey",
	"aj": "ay", "ɛw": "ew",
	// Psuedovowels
	"ṛ": "rr", "ḷ": "ll",
	// Consonants
	"t": "t", "p": "p", "ʔ": "'",
	"n": "n", "k": "k", "l": "l",
	"s": "s", "ɾ": "r", "j": "y",
	"t͡s": "ts", "t'": "tx", "m": "m",
	"v": "v", "w": "w", "h": "h",
	"ŋ": "ng", "z": "z", "k'": "kx",
	"p'": "px", "f": "f", "r": "r",
	// Reef dialect
	"b": "b", "d": "d", "g": "g",
	"ʃ": "sh", "tʃ": "ch", "ʊ": "ù",
	// mistakes and rarities
	"ʒ": "ch", "": "", " ": "",
}

// syllableToRoman converts a single IPA syllable to its Na'vi romanization.
// It assumes `syllable` is already cleaned of separators like "·", "ˈ", "ˌ".
func syllableToRoman(syllable string) string {
	var breakdown string

	// Onset
	if strings.HasPrefix(syllable, "tʃ") {
		breakdown += "ch"
		syllable = strings.TrimPrefix(syllable, "tʃ")
	} else if len(syllable) >= 4 && syllable[0:4] == "t͡s" {
		// ts
		breakdown += "ts"
		//tsp
		if hasAt("ptk", syllable, 3) {
			if nthRune(syllable, 4) == "'" {
				// ts + ejective onset
				breakdown += romanization2[syllable[4:6]]
				syllable = syllable[6:]
			} else {
				// ts + unvoiced plosive
				breakdown += romanization2[string(syllable[4])]
				syllable = syllable[5:]
			}
		} else if hasAt("lɾmnŋwj", syllable, 3) {
			// ts + other consonant
			breakdown += romanization2[nthRune(syllable, 3)]
			syllable = syllable[4+len(nthRune(syllable, 3)):]
		} else {
			// ts without a cluster
			syllable = syllable[4:]
		}
	} else if hasAt("fs", syllable, 0) {
		breakdown += nthRune(syllable, 0)
		if hasAt("ptk", syllable, 1) {
			if nthRune(syllable, 2) == "'" {
				// f/s + ejective onset
				breakdown += romanization2[syllable[1:3]]
				syllable = syllable[3:]
			} else {
				// f/s + unvoiced plosive
				breakdown += romanization2[string(syllable[1])]
				syllable = syllable[2:]
			}
		} else if hasAt("lɾmnŋwj", syllable, 1) {
			// f/s + other consonant
			breakdown += romanization2[nthRune(syllable, 1)]
			syllable = syllable[1+len(nthRune(syllable, 1)):]
		} else {
			// f/s without a cluster
			syllable = syllable[1:]
		}
	} else if hasAt("ptk", syllable, 0) {
		if nthRune(syllable, 1) == "'" {
			// ejective
			breakdown += romanization2[syllable[0:2]]
			syllable = syllable[2:]
		} else {
			// unvoiced plosive
			breakdown += romanization2[string(syllable[0])]
			syllable = syllable[1:]
		}
	} else if hasAt("ʔlɾhmnŋvwjzbdg", syllable, 0) {
		// normal onset
		breakdown += romanization2[nthRune(syllable, 0)]
		syllable = syllable[len(nthRune(syllable, 0)):]
	} else if hasAt("ʃʒ", syllable, 0) {
		// one sound represented as a cluster
		if nthRune(syllable, 0) == "ʃ" {
			breakdown += "sh"
		}
		syllable = syllable[len(nthRune(syllable, 0)):]
	}

	// Nucleus
	if len(syllable) > 1 && hasAt("jw", syllable, 1) {
		//diphthong
		breakdown += romanization2[syllable[0:len(nthRune(syllable, 0))+1]]
		syllable = string([]rune(syllable)[2:])
	} else if len(syllable) > 1 && hasAt("lr", syllable, 0) {
		// pseudovowel
		breakdown += romanization2[syllable[0:3]]
		return breakdown // pseudovowels can't have a coda
	} else if len(syllable) > 0 {
		//vowel
		breakdown += romanization2[nthRune(syllable, 0)]
		syllable = string([]rune(syllable)[1:])
	}

	// Coda
	if len(syllable) > 0 {
		if nthRune(syllable, 0) == "s" {
			breakdown += "sss" // oìsss only
		} else {
			if syllable == "k̚" {
				breakdown += "k"
			} else if syllable == "p̚" {
				breakdown += "p"
			} else if syllable == "t̚" {
				breakdown += "t"
			} else if syllable == "ʔ̚" {
				breakdown += "'"
			} else {
				if syllable[0] == 'k' && len(syllable) > 1 {
					breakdown += "kx"
				} else {
					breakdown += romanization2[syllable]
				}
			}
		}
	}

	return breakdown
}
