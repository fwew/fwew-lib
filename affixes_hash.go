package fwew_lib

import (
	"strings"
)

type ConjugationCandidate struct {
	word      string
	lenition  []string
	prefixes  []string
	suffixes  []string
	infixes   []string
	insistPOS string
}

func candidateDupe(candidate ConjugationCandidate) (c ConjugationCandidate) {
	a := ConjugationCandidate{}
	a.word = candidate.word
	a.lenition = candidate.lenition
	a.prefixes = candidate.prefixes
	a.infixes = candidate.infixes
	a.suffixes = candidate.suffixes
	a.insistPOS = candidate.insistPOS
	return a
}

var candidates []ConjugationCandidate
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
	"f":  {"f", "p"},
	"p":  {"px"},
	"h":  {"h", "k"},
	"k":  {"kx"},
	"s":  {"s", "t", "ts"},
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
var prefixes1Nouns = []string{"fì", "tsa", "kaw", "fra"}
var prefixes1Any = []string{"tì", "sä"}

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

var infixes = map[rune][]string{
	rune('a'): {"ay", "asy", "aly", "ary", "am", "alm", "arm", "ats", "awn"},
	rune('ä'): {"äng", "äp"},
	rune('e'): {"er", "ei", "eiy", "eng", "eyk"},
	rune('i'): {"iv", "ilv", "irv", "imv", "iyev"},
	rune('ì'): {"ìy", "ìsy", "ìly", "ìry", "ìm", "ìlm", "ìrm", "ìyev"},
	rune('o'): {"ol"},
	rune('u'): {"us", "uy"},
}

func isDuplicate(input string) bool {
	for _, a := range candidates {
		if input == a.word {
			return true
		}
	}
	return false
}

func isDuplicateFix(fixes []string, fix string) (newFixes []string) {
	for _, a := range fixes {
		if fix == a {
			return fixes
		}
	}
	fixes = append(fixes, fix)
	return fixes
}

func deconjugateHelper(input ConjugationCandidate, prefixCheck int, suffixCheck int, unlenite int8) []ConjugationCandidate {
	if !isDuplicate(input.word) {
		candidates = append(candidates, input)
		newString := ""
		switch prefixCheck {
		// This one will demand this makes it use lenition
		case 0:
			if strings.HasPrefix(input.word, "a") {
				newCandidate := candidateDupe(input)
				newCandidate.word = input.word[1:]
				newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "a")
				newCandidate.insistPOS = "adj."
				deconjugateHelper(newCandidate, 1, suffixCheck, -1)
			}
			// determiner prefix: "fì", "tsa", "pe", "fra"
			for _, element := range prefixes1lenition {
				// If it has a lenition-causing prefix
				if strings.HasPrefix(input.word, element) {
					lenited := false
					newString = strings.TrimPrefix(input.word, element)
					// find out the possible unlenited forms
					for _, oldPrefix := range unlenitionLetters {
						// If it has a letter that could have changed for lenition,
						if strings.HasPrefix(newString, oldPrefix) {
							// put all possibilities in the candidates
							lenited = true
							for _, newPrefix := range unlenition[oldPrefix] {
								newCandidate := candidateDupe(input)
								newCandidate.word = newPrefix + strings.TrimPrefix(newString, oldPrefix)
								newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
								if oldPrefix != newPrefix {
									newCandidate.lenition = []string{oldPrefix + "→" + newPrefix}
								}
								newCandidate.insistPOS = "n."
								deconjugateHelper(newCandidate, 1, suffixCheck, -1)
							}
							break // We don't want the "ts" to become "txs"
						}
					}
					if !lenited {
						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
						newCandidate.insistPOS = "n."
						deconjugateHelper(newCandidate, 1, suffixCheck, -1)
					}
				}
			}
			fallthrough
		case 1:
			// Non-lenition prefixes for nouns only
			for _, element := range prefixes1Nouns {
				// If it has a prefix
				if strings.HasPrefix(input.word, element) {
					// remove it
					newString = strings.TrimPrefix(input.word, element)
					// Make sure it's not a duplicate
					if !isDuplicate(newString) {
						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.insistPOS = "n."
						newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
						deconjugateHelper(newCandidate, 2, suffixCheck, -1)
					}
				}
			}

			// Non-lenition prefixes for anything
			for _, element := range prefixes1Any {
				// If it has a prefix
				if strings.HasPrefix(input.word, element) {
					// remove it
					newString = strings.TrimPrefix(input.word, element)
					// Make sure it's not a duplicate
					if !isDuplicate(newString) {
						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
						deconjugateHelper(newCandidate, 2, suffixCheck, -1)
					}
				}
			}

			// Short lenition check
			if unlenite != -1 {
				lenited := false
				for _, oldPrefix := range unlenitionLetters {
					// If it has a letter that could have changed for lenition,
					if strings.HasPrefix(input.word, oldPrefix) {
						// put all possibilities in the candidates
						for _, newPrefix := range unlenition[oldPrefix] {
							lenited = true
							newCandidate := candidateDupe(input)
							newString = newPrefix + strings.TrimPrefix(input.word, oldPrefix)
							newCandidate.word = newString
							if oldPrefix != newPrefix {
								newCandidate.lenition = []string{oldPrefix + "→" + newPrefix}
							}
							deconjugateHelper(newCandidate, 2, suffixCheck, -1)
						}
						break // We don't want the "ts" to become "txs"
					}
				}
				if !lenited {
					deconjugateHelper(input, 1, suffixCheck, -1)
				}
			}

			fallthrough
		case 2:
			if strings.HasPrefix(input.word, "fne") {
				newString = strings.TrimPrefix(input.word, "fne")
				newCandidate := candidateDupe(input)
				newCandidate.word = newString
				newCandidate.insistPOS = "n."
				newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "fne")
				deconjugateHelper(newCandidate, 3, suffixCheck, -1)
			}
		}

		switch suffixCheck {
		case 0: // adpositions, sì, o
			for _, oldSuffix := range lastSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input.word, oldSuffix) {
					newString = strings.TrimSuffix(input.word, oldSuffix)
					if !isDuplicate(newString) {
						//candidates = append(candidates, newString)
						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
						deconjugateHelper(newCandidate, prefixCheck, 1, unlenite)
					}
				}
			}
			fallthrough
		case 1:
			for _, oldSuffix := range adposuffixes {
				// If it has one of them,
				if strings.HasSuffix(input.word, oldSuffix) {
					newString = strings.TrimSuffix(input.word, oldSuffix)
					if !isDuplicate(newString) {
						//candidates = append(candidates, newString)
						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.insistPOS = "n."
						newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
						deconjugateHelper(newCandidate, prefixCheck, 2, unlenite)
					}
				}
			}
			fallthrough
		case 2:
			for _, oldSuffix := range determinerSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input.word, oldSuffix) {
					newString = strings.TrimSuffix(input.word, oldSuffix)
					if !isDuplicate(newString) {
						//candidates = append(candidates, newString)
						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.insistPOS = "n."
						newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
						deconjugateHelper(newCandidate, prefixCheck, 3, unlenite)
					}
				}
			}
			fallthrough
		case 3:
			for _, oldSuffix := range stemSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input.word, oldSuffix) {
					newString = strings.TrimSuffix(input.word, oldSuffix)
					if !isDuplicate(newString) {
						//candidates = append(candidates, newString)
						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
						deconjugateHelper(newCandidate, prefixCheck, 4, unlenite)
					}
				}
			}
			fallthrough
		case 4:
			// If it has one of them,
			if strings.HasSuffix(input.word, "yu") {
				newString = strings.TrimSuffix(input.word, "yu")
				if !isDuplicate(newString) {
					//candidates = append(candidates, newString)
					newCandidate := candidateDupe(input)
					newCandidate.word = newString
					newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, "yu")
					deconjugateHelper(newCandidate, prefixCheck, 5, unlenite)
				}
			}
		}

		// Check for infixes
		runes := []rune(input.word)
		for i, c := range runes {
			// Infixes can only begin with vowels
			if has("aäeiìou", string(c)) {
				shortString := string(runes[i:])
				for _, infix := range infixes[c] {
					if strings.HasPrefix(shortString, infix) {
						newCandidate := candidateDupe(input)
						newCandidate.word = string(runes[:i]) + strings.TrimPrefix(shortString, infix)
						newCandidate.infixes = isDuplicateFix(newCandidate.infixes, infix)
						newCandidate.insistPOS = "v."
						deconjugateHelper(newCandidate, prefixCheck, 5, unlenite)

						if infix == "ol" {
							newCandidate := candidateDupe(input)
							newCandidate.word = string(runes[:i]) + "ll" + strings.TrimPrefix(shortString, infix)
							newCandidate.infixes = isDuplicateFix(newCandidate.infixes, infix)
							newCandidate.insistPOS = "v."
							deconjugateHelper(newCandidate, prefixCheck, 5, unlenite)
						} else if infix == "er" {
							newCandidate := candidateDupe(input)
							newCandidate.word = string(runes[:i]) + "rr" + strings.TrimPrefix(shortString, infix)
							newCandidate.infixes = isDuplicateFix(newCandidate.infixes, infix)
							newCandidate.insistPOS = "v."
							deconjugateHelper(newCandidate, prefixCheck, 5, unlenite)
						}
					}
				}
			}
		}

		return candidates
	}
	return nil
}

func deconjugate(input string) []ConjugationCandidate {
	candidates = []ConjugationCandidate{} //empty array of strings
	newCandidate := ConjugationCandidate{}
	newCandidate.word = input
	deconjugateHelper(newCandidate, 0, 0, 0)
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
