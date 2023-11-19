package fwew_lib

import (
	"fmt"
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

var lenitionable = []string{
	"ts",
	"px", "tx", "kx",
	"p", "t", "k",
	"f", "s", "h",
	"'",
}
var lenition = map[string]string{
	"px": "p",
	"tx": "t",
	"kx": "k",
	"p":  "f",
	"t":  "s",
	"k":  "h",
	"ts": "s",
	"'":  "",
}

var prefixes1Nouns = []string{"fì", "tsa", "fra"}
var prefixes1lenition = []string{"pe", "fay", "tsay", "fìme", "tsame", "fìpxe",
	"tsapxe", "pxe", "pepe", "peme", "pay", "ay", "me"}
var stemPrefixes = []string{"fne", "sna", "munsna"}
var verbPrefixes = []string{"tsuk", "ketsuk"}

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
var verbSuffixes = []string{"tswo", "yu"}

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

func infixError(query string, didYouMean string, ipa string) Word {
	d := Word{}
	d.Navi = query
	d.EN = didYouMean
	d.DE = didYouMean
	d.ET = didYouMean
	d.FR = didYouMean
	d.NL = didYouMean
	d.HU = didYouMean
	d.PL = didYouMean
	d.RU = didYouMean
	d.SV = didYouMean
	d.TR = didYouMean
	d.IPA = ipa
	d.PartOfSpeech = "err."
	return d
}

func deconjugateHelper(input ConjugationCandidate, prefixCheck int, suffixCheck int, unlenite int8, checkInfixes bool) []ConjugationCandidate {
	fmt.Println(input.word + "(" + input.insistPOS + ")")
	fmt.Println(input.prefixes)
	fmt.Println(input.infixes)
	fmt.Println(input.suffixes)
	if !isDuplicate(input) {
		candidates = append(candidates, input)
		newString := ""

		// Make sure that the first set of prefices (a, nì, ke) aren't combined with suffixes
		newPrefixCheck := prefixCheck
		if newPrefixCheck == 0 {
			newPrefixCheck = 1
		}

		switch prefixCheck {
		case 0:
			if strings.HasPrefix(input.word, "a") {
				newCandidate := candidateDupe(input)
				newCandidate.word = input.word[1:]
				newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "a")
				newCandidate.insistPOS = "adj."
				deconjugateHelper(newCandidate, 1, suffixCheck, -1, false)
			} else if strings.HasPrefix(input.word, "nì") {
				newCandidate := candidateDupe(input)
				newCandidate.word = strings.TrimPrefix(input.word, "nì")
				newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "nì")
				newCandidate.insistPOS = "nì."
				// No other affixes allowed
				deconjugateHelper(newCandidate, 10, 10, -1, false) // No other fixes
			} else if strings.HasPrefix(input.word, "ke") {
				// remove it
				newString = strings.TrimPrefix(input.word, "ke")

				newCandidate := candidateDupe(input)
				newCandidate.word = newString
				newCandidate.insistPOS = "any"
				newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "ke")
				deconjugateHelper(newCandidate, 10, 10, -1, false)

				newCandidate.word = "e" + newString
				deconjugateHelper(newCandidate, 10, 10, -1, false)
			}
			fallthrough
		case 1:
			for _, element := range verbPrefixes {
				// If it has a prefix
				if strings.HasPrefix(input.word, element) {
					// remove it
					newCandidate := candidateDupe(input)
					newCandidate.word = strings.TrimPrefix(input.word, element)
					newCandidate.insistPOS = "v."
					newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
					deconjugateHelper(newCandidate, 10, 10, -1, false)

					// check "tsatan", "tan" and "atan"
					newCandidate.word = get_last_rune(element, 1) + newCandidate.word
					deconjugateHelper(newCandidate, 10, 10, -1, false)
				}
			}
			fallthrough
		case 2:
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
					deconjugateHelper(newCandidate, 3, suffixCheck, -1, false)

					// check "tsatan", "tan" and "atan"
					newCandidate.word = get_last_rune(element, 1) + newString
					deconjugateHelper(newCandidate, 3, suffixCheck, -1, false)
				}
			}
			fallthrough
		case 3:
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
						deconjugateHelper(newCandidate, 4, suffixCheck, -1, false)

						// check "pxeylan", "ylan" and "'eylan"
						newCandidate.word = "'" + newCandidate.word
						deconjugateHelper(newCandidate, 4, suffixCheck, -1, false)
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
								deconjugateHelper(newCandidate, 4, suffixCheck, -1, false)
							}
							break // We don't want the "ts" to become "txs"
						}
					}
					if !lenited {
						newCandidate.word = newString
						deconjugateHelper(newCandidate, 3, suffixCheck, -1, false)
					}
				}
			}
			fallthrough
		case 4:
			for _, element := range stemPrefixes {
				// If it has a prefix
				if strings.HasPrefix(input.word, element) {
					// remove it
					newCandidate := candidateDupe(input)
					newCandidate.word = strings.TrimPrefix(input.word, element)
					newCandidate.insistPOS = "n."
					newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
					deconjugateHelper(newCandidate, 5, suffixCheck, -1, false)

					// check "tsatan", "tan" and "atan"
					newCandidate.word = get_last_rune(element, 1) + newCandidate.word
					deconjugateHelper(newCandidate, 5, suffixCheck, -1, false)
				}
			}
			fallthrough
		case 5:
			if strings.HasPrefix(input.word, "tì") {
				if input.insistPOS == "any" || input.insistPOS == "n." {
					// remove it
					newCandidate := candidateDupe(input)
					newCandidate.word = strings.TrimPrefix(input.word, "tì")
					newCandidate.insistPOS = "v."
					newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "tì")
					deconjugateHelper(newCandidate, 10, 10, -1, true) // No other prefixes allowed

					newCandidate.word = "ì" + newCandidate.word
					deconjugateHelper(newCandidate, 10, 10, -1, true) // Or any additional suffixes
				}
			}
		}

		switch suffixCheck {
		case 0:
			// If it has one of them,
			for _, oldSuffix := range verbSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input.word, oldSuffix) {
					newString = strings.TrimSuffix(input.word, oldSuffix)

					newCandidate := candidateDupe(input)
					newCandidate.word = newString
					newCandidate.insistPOS = "v."
					newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
					deconjugateHelper(newCandidate, newPrefixCheck, 1, unlenite, false)

					if oldSuffix == "yu" && strings.HasSuffix(newString, "si") {
						newCandidate.word = strings.TrimSuffix(newString, "si") + " si"
						deconjugateHelper(newCandidate, 10, 10, unlenite, false) // don't allow any other prefixes or suffixes
					}
				}
			}
			fallthrough
		case 1:
			// If it has one of them,
			if strings.HasSuffix(input.word, "a") {
				newString = strings.TrimSuffix(input.word, "a")

				newCandidate := candidateDupe(input)
				newCandidate.word = newString
				newCandidate.insistPOS = "adj."
				newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, "a")
				deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, true)
			}
			for _, oldSuffix := range lastSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input.word, oldSuffix) {
					newString = strings.TrimSuffix(input.word, oldSuffix)

					newCandidate := candidateDupe(input)
					newCandidate.word = newString
					newCandidate.insistPOS = "any"
					newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
					deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, true)
				}
			}
			fallthrough
		case 2: // adpositions, sì, o, case endings
			if input.insistPOS == "any" || input.insistPOS == "n." {
				for _, oldSuffix := range adposuffixes {
					// If it has one of them,
					if strings.HasSuffix(input.word, oldSuffix) {
						newString = strings.TrimSuffix(input.word, oldSuffix)

						newCandidate := candidateDupe(input)
						newCandidate.word = newString
						newCandidate.insistPOS = "n."
						newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
						// all set to 4 to avoid mengeyä -> mengo -> me + 'eng + o
						deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, false)

						if oldSuffix == "ä" {
							// soaiä, tìftiä, etx.
							newString += "a"
							newCandidate.word = newString
							deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, false)
						} else if oldSuffix == "yä" && strings.HasSuffix(newString, "e") {
							// A one-off
							if newString == "tse" {
								newCandidate.word = "tsaw"
								deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, false)
							}
							// ngeyä -> nga
							newCandidate.word = strings.TrimSuffix(newString, "e") + "a"
							deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, false)
							// oengeyä
							newCandidate.word = strings.TrimSuffix(newString, "e")
							if newCandidate.word == "oeng" { //no mengeyä -> meng -> me + 'eng
								deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, false)
							}
							// sneyä -> sno
							newCandidate.word = strings.TrimSuffix(newString, "e") + "o"
							deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, false)
						} else if oldSuffix == "l" && strings.HasSuffix(newString, "a") {
							// oengal
							newCandidate.word = strings.TrimSuffix(newString, "a")
							deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, false)
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
					deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, false)
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
					deconjugateHelper(newCandidate, newPrefixCheck, 5, unlenite, false)
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
						deconjugateHelper(newCandidate, prefixCheck, suffixCheck, -1, false)
					}
					break // We don't want the "ts" to become "txs"
				}
			}
		}

		if checkInfixes {
			// Maybe someone else came in with stripped infixes
			if len(input.word) > 2 && input.word[len(input.word)-3] != ' ' && strings.HasSuffix(input.word, "si") {
				newCandidate := candidateDupe(input)
				newCandidate.word = strings.TrimSuffix(input.word, "si") + " si"
				newCandidate.insistPOS = "v."
				deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, true)
			} else { // If there is a "si", we don't need to check for infixes
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
								deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, true)

								if infix == "ol" {
									newCandidate := candidateDupe(input)
									newCandidate.word = string(runes[:i]) + "ll" + strings.TrimPrefix(shortString, infix)
									newCandidate.infixes = isDuplicateFix(newCandidate.infixes, infix)
									newCandidate.insistPOS = "v."
									deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, true)
								} else if infix == "er" {
									newCandidate := candidateDupe(input)
									newCandidate.word = string(runes[:i]) + "rr" + strings.TrimPrefix(shortString, infix)
									newCandidate.infixes = isDuplicateFix(newCandidate.infixes, infix)
									newCandidate.insistPOS = "v."
									deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, true)
								}
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
	deconjugateHelper(newCandidate, 0, 0, 0, true)
	candidates = candidates[1:]
	return candidates
}

func TestDeconjugations(searchNaviWord string) (results []Word) {
	conjugations := deconjugate(searchNaviWord)
	for _, candidate := range conjugations {
		for _, c := range dictHash[candidate.word] {
			if _, ok := dictHash[candidate.word]; ok {

				// Find gerunds (tì-v<us>erb, treated like a noun)
				gerund := false
				infixBan := false
				doubleBan := false
				participle := false

				// Find gerunds (tì-v<us>erb, the act of [verb]ing)
				if len(candidate.infixes) == 1 && candidate.infixes[0] == "us" {
					// Reverse search is more likely to find it immediately
					for i := len(candidate.prefixes) - 1; i >= 0; i-- {
						if candidate.prefixes[i] == "tì" {
							gerund = true
							break
						}
					}
					if !gerund {
						participle = true
					}
				} else if len(candidate.infixes) > 0 {
					// Now reverse search is just gratuitous
					for i := len(candidate.infixes) - 1; i >= 0; i-- {
						if candidate.infixes[i] == "us" || candidate.infixes[i] == "awn" {
							participle = true
							break
						}
					}
				}

				// If the insistPOS and found word agree they are nouns
				if len(candidate.suffixes) == 1 && candidate.suffixes[0] == "tswo" {
					siVerb := false
					if len(candidate.prefixes) == 0 && len(candidate.infixes) == 0 {
						if _, ok := multiword_words[candidate.word]; ok {
							for _, b := range multiword_words[candidate.word] {
								if b[0] == "si" {
									siVerb = true
									a := c
									a.Navi = candidate.word + " si"
									a.Affixes.Lenition = candidate.lenition
									a.Affixes.Prefix = candidate.prefixes
									a.Affixes.Infix = candidate.infixes
									a.Affixes.Suffix = candidate.suffixes
									results = append(results, a)
									break
								}
							}
						}
						if !siVerb {
							a := c
							a.Navi = candidate.word
							a.Affixes.Lenition = candidate.lenition
							a.Affixes.Prefix = candidate.prefixes
							a.Affixes.Infix = candidate.infixes
							a.Affixes.Suffix = candidate.suffixes
							results = append(results, a)
						}
					}
				} else if gerund {
					if strings.HasPrefix(c.PartOfSpeech, "v") {
						// Make sure the <us> is in the correct place
						rebuiltVerb := strings.ReplaceAll(c.InfixLocations, "<0>", "")
						rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<1>", "us")
						rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<2>", "")

						// Does the noun actually contain the verb?
						if strings.Contains(searchNaviWord, rebuiltVerb) {
							a := c
							a.Affixes.Lenition = candidate.lenition
							a.Affixes.Prefix = candidate.prefixes
							a.Affixes.Infix = candidate.infixes
							a.Affixes.Suffix = candidate.suffixes
							results = append(results, a)
						} else {
							results = append(results, infixError(searchNaviWord, "Did you mean **tì"+rebuiltVerb+"**?", c.IPA))
						}
					}
				} else if candidate.insistPOS == "n." {
					// n., pn. and Prop.n. (but not vin.)
					if len(candidate.infixes) == 0 && c.PartOfSpeech[0] != 'v' && strings.HasSuffix(c.PartOfSpeech, "n.") {
						a := c
						a.Affixes.Lenition = candidate.lenition
						a.Affixes.Prefix = candidate.prefixes
						a.Affixes.Suffix = candidate.suffixes
						results = append(results, a)
					}
				} else if candidate.insistPOS == "adj." {
					posNoun := c.PartOfSpeech
					if len(candidate.infixes) == 0 && (posNoun == "adj." || posNoun == "num.") {
						a := c
						a.Affixes.Lenition = candidate.lenition
						a.Affixes.Prefix = candidate.prefixes
						a.Affixes.Suffix = candidate.suffixes
						results = append(results, a)
					}
				} else if candidate.insistPOS == "v." {
					posNoun := c.PartOfSpeech
					if strings.HasPrefix(posNoun, "v") {
						// Verbs with -tswo or -yu cannot have infixes
						if len(candidate.suffixes) > 0 {
							// Reverse search is more likely to find it immediately
							for i := len(candidate.suffixes) - 1; i >= 0; i-- {
								for _, j := range verbSuffixes {
									if candidate.suffixes[i] == j {
										infixBan = true
										break
									}
								}
								if infixBan {
									break
								}
							}
						}

						if len(candidate.prefixes) > 0 {
							// Reverse search is more likely to find it immediately
							for i := len(candidate.prefixes) - 1; i >= 0; i-- {
								for _, j := range verbPrefixes {
									if candidate.prefixes[i] == j {
										if infixBan {
											doubleBan = true
											break
										}
										infixBan = true
										break
									}
								}
								if infixBan || doubleBan {
									break
								}
							}
						}

						// Take action on tsuk-verb-yus
						if doubleBan {
							continue
						}

						a := c
						a.Affixes.Lenition = candidate.lenition
						a.Affixes.Prefix = candidate.prefixes
						a.Affixes.Suffix = candidate.suffixes
						a.Affixes.Infix = candidate.infixes

						if infixBan {
							if len(candidate.infixes) > 0 {
								continue // No nonsense here
							} else {
								results = append(results, a)
							}
						}

						// Make it verify the infixes are in the correct place

						// pre-first position infixes
						rebuiltVerb := c.InfixLocations
						firstInfixes := ""
						found := false

						for _, infix := range prefirst {
							for _, newInfix := range candidate.infixes {
								if newInfix == infix {
									firstInfixes += infix
									found = true
								}
							}
						}
						rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<0>", firstInfixes)

						// first position infixes
						found = false
						firstInfixes = ""
						for _, infix := range first {
							for _, newInfix := range candidate.infixes {
								if newInfix == infix {
									rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<1>", infix)
									firstInfixes = infix
									found = true
									break
								}
							}
							if found {
								break
							}
						}
						rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<1>", "")

						// second position infixes
						found = false
						for _, infix := range second {
							for _, newInfix := range candidate.infixes {
								if newInfix == infix {
									rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<2>", infix)
									found = true
									break
								}
							}
							if found {
								break
							}
						}
						rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<2>", "")

						rebuiltVerb = strings.TrimSpace(rebuiltVerb)

						if identicalRunes(rebuiltVerb, searchNaviWord) {
							results = append(results, a)
						} else if participle {
							if identicalRunes("a"+rebuiltVerb, searchNaviWord) {
								// a-v<us>erb and a-v<awn>erb
								results = append(results, a)
							} else if identicalRunes(rebuiltVerb+"a", searchNaviWord) {
								// v<us>erb-a and v<awn>erb-a
								results = append(results, a)
							} else if firstInfixes == "us" {
								results = append(results, infixError(searchNaviWord, "Did you mean **"+rebuiltVerb+"**?", c.IPA))
							}
						} else if gerund { // ti is needed to weed out non-productive tì-verbs
							results = append(results, infixError(searchNaviWord, "Did you mean **"+rebuiltVerb+"**?", c.IPA))
						}
					}
				} else if candidate.insistPOS == "nì." {
					posNoun := c.PartOfSpeech
					if len(candidate.infixes) == 0 && (posNoun == "adj." || posNoun == "pn.") {
						a := c
						a.Affixes.Lenition = candidate.lenition
						a.Affixes.Prefix = candidate.prefixes
						a.Affixes.Suffix = candidate.suffixes
						results = append(results, a)
					}
				} else if len(candidate.infixes) == 0 {
					a := c
					a.Affixes.Lenition = candidate.lenition
					a.Affixes.Prefix = candidate.prefixes
					a.Affixes.Suffix = candidate.suffixes
					results = append(results, a)
				}

			}
		}
	}
	return
}
