package fwew_lib

import (
	"strings"
)

var candidates []string
var unlenitionLetters = []string{
	"ts", "kx", "tx", "px", // traps digraphs because they cannot unlenite
	"f", "p", "h", "k", "s",
	"t", "a", "ä", "e", "i",
	"ì", "o", "u", "ù",
}

// "ts" is there to prevent "ts" from becoming "txs"
var unlenition = map[string][]string{
	// digraphs cannot unlenite
	"ts": {}, // here to trap the "ts" ahead of the "t"
	"px": {}, // here to trap the "px" ahead of the "p"
	"kx": {}, // here to trap the "kx" ahead of the "k"
	"tx": {}, // here to trap the "tx" ahead of the "t"
	"f":  {"p"},
	"p":  {"px"},
	"h":  {"k"},
	"k":  {"kx"},
	"s":  {"t", "ts"},
	"t":  {"tx"},
	"a":  {"'a"},
	"ä":  {"'ä"},
	"e":  {"'e"},
	"i":  {"'i"},
	"ì":  {"'ì"},
	"o":  {"'o"},
	"u":  {"'u"},
	"ù":  {"'ù"},
}

var prefixes1lenition = []string{"fay", "tsay", "fìme", "tsame", "fìpxe", "tsapxe", "pxe", "pepe", "ay", "me", "pe"}
var prefixes1 = []string{"fì", "tsa", "kaw", "ke", "fra", "a"}

var lastSuffixes = []string{"sì", "to", "a"}
var adposuffixes = []string{
	// adpositions that can be mistaken for case endings
	"pxel",                //"agentive"
	"mungwrr",             //"dative"
	"kxamlä", "ìlä", "wä", //"genitive"
	"teri", //"topical"
	// Case endings
	"ìl", "l", "it", "ti", "t", "ur", "ru", "r", "yä", "ä", "ìri", "ri",
	// Alphabetized the reverse of these things with exceptions for mistaken ones
	"nemfa", "rofa", "ka", "fa", "na", "ta",
	"lisre", "pxisre", "sre", "luke", "ne",
	"fpi",
	"mì",
	"lok",
	"mìkam", "kam",
	"sìn",
	"äo", "eo", "io", "uo", "ro",
	"tafkip", "takip", "fkip", "kip",
	"ftu", "hu",
	"pximaw", "maw", "pxaw", "few",
	"vay", "kay",
}
var determinerSuffixes = []string{"pe", "o"}
var stemSuffixes = []string{"tsyìp", "fkeyk"}
var animateSuffixes = []string{"yu", "tu"}

func isDuplicate(input string) bool {
	for _, a := range candidates {
		if input == a {
			return true
		}
	}
	return false
}

func deconjugateHelper(input string, prefixCheck int, suffixCheck int) []string {
	if !isDuplicate(input) {
		candidates = append(candidates, input)
		newString := ""
		switch prefixCheck {
		case 0: // determiner prefix: "fì", "tsa", "pe", "fra"
			for _, element := range prefixes1lenition {
				// If it has a lenition-causing prefix
				if strings.HasPrefix(input, element) {
					// remove it
					deconjugateHelper(strings.TrimPrefix(input, element), 0, suffixCheck)
					// Make sure it's not a duplicate
					if !isDuplicate(newString) {
						candidates = append(candidates, newString)
						// find out the possible unlenited forms
						for _, oldPrefix := range unlenitionLetters {
							// If it has a letter that could have changed for lenition,
							if strings.HasPrefix(newString, oldPrefix) {
								// put all possibilities in the candidates
								for _, newPrefix := range unlenition[oldPrefix] {
									deconjugateHelper(newPrefix+strings.TrimPrefix(newString, oldPrefix), 1, suffixCheck)
								}
								break // We don't want the "ts" to become "txs"
							}
						}
						// Now check all the new candicates
						for _, candidate := range candidates {
							deconjugateHelper(candidate, 0, suffixCheck)
						}
					}
				}
			}

			// Non-lenition prefixes
			for _, element := range prefixes1 {
				// If it has a prefix
				if strings.HasPrefix(input, element) {
					// remove it
					newString = strings.TrimPrefix(input, element)
					// Make sure it's not a duplicate
					if !isDuplicate(newString) {
						//candidates = append(candidates, newString)
						deconjugateHelper(newString, 0, suffixCheck)
					}
				}
			}

			// Short lenition check
			for _, oldPrefix := range unlenitionLetters {
				// If it has a letter that could have changed for lenition,
				if strings.HasPrefix(input, oldPrefix) {
					// put all possibilities in the candidates
					for _, newPrefix := range unlenition[oldPrefix] {
						newString = newPrefix + strings.TrimPrefix(input, oldPrefix)
						deconjugateHelper(newString, 1, suffixCheck)
					}
					break // We don't want the "ts" to become "txs"
				}
			}
			fallthrough
		case 1:
			if strings.HasPrefix(input, "fne") {
				deconjugateHelper(strings.TrimPrefix(input, "fne"), 2, suffixCheck)
			}
		}

		switch suffixCheck {
		case 0: // adpositions, sì, o
			for _, oldSuffix := range lastSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input, oldSuffix) {
					newString = strings.TrimSuffix(input, oldSuffix)
					if !isDuplicate(newString) {
						//candidates = append(candidates, newString)
						deconjugateHelper(newString, prefixCheck, 1)
					}
				}
			}
			fallthrough
		case 1:
			for _, oldSuffix := range adposuffixes {
				// If it has one of them,
				if strings.HasSuffix(input, oldSuffix) {
					newString = strings.TrimSuffix(input, oldSuffix)
					if !isDuplicate(newString) {
						//candidates = append(candidates, newString)
						deconjugateHelper(newString, prefixCheck, 2)
					}
				}
			}
			fallthrough
		case 2:
			for _, oldSuffix := range determinerSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input, oldSuffix) {
					newString = strings.TrimSuffix(input, oldSuffix)
					if !isDuplicate(newString) {
						//candidates = append(candidates, newString)
						deconjugateHelper(newString, prefixCheck, 3)
					}
				}
			}
			fallthrough
		case 3:
			for _, oldSuffix := range stemSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input, oldSuffix) {
					newString = strings.TrimSuffix(input, oldSuffix)
					if !isDuplicate(newString) {
						//candidates = append(candidates, newString)
						deconjugateHelper(newString, prefixCheck, 4)
					}
				}
			}
			fallthrough
		case 4:
			for _, oldSuffix := range animateSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input, oldSuffix) {
					newString = strings.TrimSuffix(input, oldSuffix)
					if !isDuplicate(newString) {
						//candidates = append(candidates, newString)
						deconjugateHelper(newString, prefixCheck, 5)
					}
				}
			}
		}

		return candidates
	}
	return nil
}

func deconjugate(input string) []string {
	candidates = []string{} //empty array of strings
	deconjugateHelper(input, 0, 0)
	candidates = candidates[1:]
	return candidates
}

/*func main() {
	fmt.Println(deconjugate("uturu"))
	fmt.Println(deconjugate("faysawtute"))
	fmt.Println(deconjugate("fayfalulukantsyìpperu"))
	fmt.Println(deconjugate("txaw"))
	fmt.Println(deconjugate("pepefneutraltsyìpftusì"))
	fmt.Println(deconjugate("kawtul"))
	fmt.Println(deconjugate("tsamsiyu"))
	fmt.Println(deconjugate("mesukrur"))
}*/
