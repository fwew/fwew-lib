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
	"a":  {"a", "'a"},
	"ä":  {"ä", "'ä"},
	"e":  {"e", "'e"},
	"i":  {"i", "'i"},
	"ì":  {"ì", "'ì"},
	"o":  {"o", "'o"},
	"u":  {"u", "'u"},
	"ù":  {"ù", "'ù"},
}

var prefixes1Nouns = []string{"fì", "tsa", "fra"}
var prefixes1lenition = []string{"pe", "fay", "tsay", "fìme", "tsame", "fìpxe",
	"tsapxe", "pxe", "pepe", "peme", "pay", "ay", "me"}
var stemPrefixes = []string{"fne", "sna", "munsna"}

var lastSuffixes = []string{"sì", "to"}
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

var prefirst = []string{"äp", "eyk"}
var first = []string{"ay", "asy", "aly", "ary", "ìy", "ìsy", "ìly", "ìry", "ol", "er", "ìm",
	"ìlm", "ìrm", "am", "alm", "arm", "ìyev", "iyev", "iv", "ilv", "irv", "imv", "us", "awn"}
var second = []string{"ei", "eiy", "äng", "eng", "uy", "ats"}

func isDuplicate(input ConjugationCandidate) bool {
	for _, a := range candidates {
		if input.word == a.word && input.insistPOS == a.insistPOS {
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

func verifyInfixes(original string, bare string, infixes []string) {

}

func deconjugateHelper(input ConjugationCandidate, prefixCheck int, suffixCheck int, unlenite int8) []ConjugationCandidate {
	if !isDuplicate(input) {
		candidates = append(candidates, input)
		newString := ""

		// Make sure that the first set of prefices (a, nì, ke) aren't combined with suffixes
		newPrefixCheck := prefixCheck
		if newPrefixCheck == 0 {
			newPrefixCheck = 1
		}

		if input.insistPOS != "v." {
			switch prefixCheck {
			case 0:
				if strings.HasPrefix(input.word, "a") {
					newCandidate := candidateDupe(input)
					newCandidate.word = input.word[1:]
					newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "a")
					newCandidate.insistPOS = "adj."
					deconjugateHelper(newCandidate, 1, suffixCheck, -1)
				} else if strings.HasPrefix(input.word, "nì") {
					newCandidate := candidateDupe(input)
					newCandidate.word = strings.TrimPrefix(input.word, "nì")
					newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "nì")
					newCandidate.insistPOS = "nì."
					// No other affixes allowed
					deconjugateHelper(newCandidate, 10, 10, -1) // No other fixes
				} else if strings.HasPrefix(input.word, "ke") {
					// remove it
					newString = strings.TrimPrefix(input.word, "ke")

					newCandidate := candidateDupe(input)
					newCandidate.word = newString
					newCandidate.insistPOS = "any"
					newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "ke")
					deconjugateHelper(newCandidate, 10, 10, -1)

					newCandidate.word = "e" + newString
					deconjugateHelper(newCandidate, 10, 10, -1)
				}
				fallthrough
			case 1:
				// Non-lenition prefixes for nouns only
				for _, element := range prefixes1Nouns {
					// If it has a prefix
					if strings.HasPrefix(input.word, element) {
						// remove it
						newString = strings.TrimPrefix(input.word, element)

						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.insistPOS = "n."
						newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
						deconjugateHelper(newCandidate, 2, suffixCheck, -1)

						// check "tsatan", "tan" and "atan"
						newCandidate.word = get_last_rune(element, 1) + newString
						deconjugateHelper(newCandidate, 2, suffixCheck, -1)
					}
				}
				fallthrough
			case 2:
				// This one will demand this makes it use lenition
				for _, element := range prefixes1lenition {
					// If it has a lenition-causing prefix
					if strings.HasPrefix(input.word, element) {
						lenited := false
						newString = strings.TrimPrefix(input.word, element)

						newCandidate := candidateDupe(input)
						newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
						newCandidate.insistPOS = "n."

						// Could it be pekoyu (pe + 'ekoyu, not pe + kxoyu)
						if has("aäeiìou", get_last_rune(element, 1)) {
							// check "pxeyktan", "yktan" and "eyktan"
							newCandidate.word = get_last_rune(element, 1) + newString
							deconjugateHelper(newCandidate, 3, suffixCheck, -1)

							// check "pxeylan", "ylan" and "'eylan"
							newCandidate.word = "'" + newCandidate.word
							deconjugateHelper(newCandidate, 3, suffixCheck, -1)
						}

						// find out the possible unlenited forms
						for _, oldPrefix := range unlenitionLetters {
							// If it has a letter that could have changed for lenition,
							if strings.HasPrefix(newString, oldPrefix) {
								// put all possibilities in the candidates
								lenited = true

								for _, newPrefix := range unlenition[oldPrefix] {
									newCandidate.word = newPrefix + strings.TrimPrefix(newString, oldPrefix)
									if oldPrefix != newPrefix {
										newCandidate.lenition = []string{oldPrefix + "→" + newPrefix}
									}
									deconjugateHelper(newCandidate, 3, suffixCheck, -1)
								}
								break // We don't want the "ts" to become "txs"
							}
						}
						if !lenited {
							newCandidate.word = newString
							deconjugateHelper(newCandidate, 3, suffixCheck, -1)
						}
					}
				}
				fallthrough
			case 3:
				if strings.HasPrefix(input.word, "tsuk") {
					newCandidate := candidateDupe(input)
					newCandidate.word = strings.TrimPrefix(input.word, "tsuk")
					newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "tsuk")
					newCandidate.insistPOS = "v."
					deconjugateHelper(newCandidate, 4, suffixCheck, -1)
				}

				for _, element := range stemPrefixes {
					// If it has a prefix
					if strings.HasPrefix(input.word, element) {
						// remove it
						newCandidate := candidateDupe(input)
						newCandidate.word = strings.TrimPrefix(input.word, element)
						newCandidate.insistPOS = "n."
						newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
						deconjugateHelper(newCandidate, 4, suffixCheck, -1)

						// check "tsatan", "tan" and "atan"
						newCandidate.word = get_last_rune(element, 1) + newCandidate.word
						deconjugateHelper(newCandidate, 4, suffixCheck, -1)
					}
				}
				fallthrough
			case 4:
				if strings.HasPrefix(input.word, "tì") {
					if input.insistPOS == "any" || input.insistPOS == "n." {
						// remove it
						newCandidate := candidateDupe(input)
						newCandidate.word = strings.TrimPrefix(input.word, "tì")
						newCandidate.insistPOS = "v."
						newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "tì")
						deconjugateHelper(newCandidate, 10, 10, -1) // No other prefixes allowed

						newCandidate.word = "ì" + newCandidate.word
						deconjugateHelper(newCandidate, 10, 10, -1) // Or any additional suffixes
					}
				}
			}

			switch suffixCheck {
			case 0:
				// If it has one of them,
				if strings.HasSuffix(input.word, "a") {
					newString = strings.TrimSuffix(input.word, "a")

					newCandidate := candidateDupe(input)
					newCandidate.word = newString
					newCandidate.insistPOS = "adj."
					newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, "a")
					deconjugateHelper(newCandidate, newPrefixCheck, 1, unlenite)
				}
				for _, oldSuffix := range lastSuffixes {
					// If it has one of them,
					if strings.HasSuffix(input.word, oldSuffix) {
						newString = strings.TrimSuffix(input.word, oldSuffix)

						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.insistPOS = "any"
						newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
						deconjugateHelper(newCandidate, newPrefixCheck, 1, unlenite)
					}
				}
				fallthrough
			case 1: // adpositions, sì, o
				// If it has one of them,
				if strings.HasSuffix(input.word, "tswo") {
					newString = strings.TrimSuffix(input.word, "tswo")

					newCandidate := candidateDupe(input)
					newCandidate.word = newString
					newCandidate.insistPOS = "v."
					newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, "tswo")
					deconjugateHelper(newCandidate, newPrefixCheck, 6, unlenite)
				}
				fallthrough
			case 2:
				if input.insistPOS == "any" || input.insistPOS == "n." {
					for _, oldSuffix := range adposuffixes {
						// If it has one of them,
						if strings.HasSuffix(input.word, oldSuffix) {
							newString = strings.TrimSuffix(input.word, oldSuffix)

							newCandidate := candidateDupe(input)
							newCandidate.word = newString
							newCandidate.insistPOS = "n."
							newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
							deconjugateHelper(newCandidate, newPrefixCheck, 3, unlenite)

							if oldSuffix == "ä" {
								// soaiä, tìftiä, etx.
								newString += "a"
								newCandidate.word = newString
								deconjugateHelper(newCandidate, newPrefixCheck, 3, unlenite)
							} else if oldSuffix == "yä" && strings.HasSuffix(newString, "e") {
								// A one-off
								if newString == "tse" {
									newCandidate.word = "tsaw"
									deconjugateHelper(newCandidate, newPrefixCheck, 3, unlenite)
								}
								// oengeyä
								newCandidate.word = strings.TrimSuffix(newString, "e")
								deconjugateHelper(newCandidate, newPrefixCheck, 3, unlenite)
								// sneyä -> sno
								newCandidate.word = newCandidate.word + "o"
								deconjugateHelper(newCandidate, newPrefixCheck, 3, unlenite)
							} else if oldSuffix == "l" && strings.HasSuffix(newString, "a") {
								// oengal
								newCandidate.word = strings.TrimSuffix(newString, "a")
								deconjugateHelper(newCandidate, newPrefixCheck, 3, unlenite)
							}
						}
					}
				}
				fallthrough
			case 3:
				for _, oldSuffix := range determinerSuffixes {
					// If it has one of them,
					if strings.HasSuffix(input.word, oldSuffix) {
						newString = strings.TrimSuffix(input.word, oldSuffix)

						//candidates = append(candidates, newString)
						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.insistPOS = "n."
						newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
						deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite)
					}
				}
				fallthrough
			case 4:
				for _, oldSuffix := range stemSuffixes {
					// If it has one of them,
					if strings.HasSuffix(input.word, oldSuffix) {
						newString = strings.TrimSuffix(input.word, oldSuffix)

						//candidates = append(candidates, newString)
						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.insistPOS = "n."
						newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
						deconjugateHelper(newCandidate, newPrefixCheck, 5, unlenite)
					}
				}
				fallthrough
			case 5:
				// If it has one of them,
				if strings.HasSuffix(input.word, "yu") {
					newString = strings.TrimSuffix(input.word, "yu")

					//candidates = append(candidates, newString)
					newCandidate := candidateDupe(input)
					newCandidate.word = newString
					newCandidate.insistPOS = "v."
					newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, "yu")
					deconjugateHelper(newCandidate, newPrefixCheck, 6, unlenite)

					if strings.HasSuffix(newString, "si") {
						newCandidate.word = strings.TrimSuffix(newString, "si") + " si"
						deconjugateHelper(newCandidate, newPrefixCheck, 6, unlenite)
					}
				}
			}
		}

		// Short lenition check
		if unlenite != -1 {
			for _, oldPrefix := range unlenitionLetters {
				// If it has a letter that could have changed for lenition,
				if strings.HasPrefix(input.word, oldPrefix) {
					// put all possibilities in the candidates
					for _, newPrefix := range unlenition[oldPrefix] {
						newCandidate := candidateDupe(input)
						newString = newPrefix + strings.TrimPrefix(input.word, oldPrefix)
						newCandidate.word = newString
						if oldPrefix != newPrefix {
							newCandidate.lenition = []string{oldPrefix + "→" + newPrefix}
						}
						deconjugateHelper(newCandidate, prefixCheck, suffixCheck, -1)
					}
					break // We don't want the "ts" to become "txs"
				}
			}
		}

		if input.insistPOS == "any" || input.insistPOS == "v." || input.insistPOS == "adj." {
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
							deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite)

							if infix == "ol" {
								newCandidate := candidateDupe(input)
								newCandidate.word = string(runes[:i]) + "ll" + strings.TrimPrefix(shortString, infix)
								newCandidate.infixes = isDuplicateFix(newCandidate.infixes, infix)
								newCandidate.insistPOS = "v."
								deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite)
							} else if infix == "er" {
								newCandidate := candidateDupe(input)
								newCandidate.word = string(runes[:i]) + "rr" + strings.TrimPrefix(shortString, infix)
								newCandidate.infixes = isDuplicateFix(newCandidate.infixes, infix)
								newCandidate.insistPOS = "v."
								deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite)
							}
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
	newCandidate.insistPOS = "any"
	deconjugateHelper(newCandidate, 0, 0, 0)
	candidates = candidates[1:]
	return candidates
}
