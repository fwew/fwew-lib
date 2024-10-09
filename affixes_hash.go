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
var candidateMap = map[string]ConjugationCandidate{}
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
var prefixes1lenition = []string{"pe", "fay",
	"pxe", "pepe", "pay", "ay", "me"}
var stemPrefixes = []string{"fne", "sna", "munsna"}
var verbPrefixes = []string{"tsuk", "ketsuk"}

var adposuffixes = []string{
	// adpositions that can be mistaken for case endings
	"pxel",                //"agentive"
	"mungwrr",             //"dative"
	"kxamlä", "ìlä", "wä", //"genitive"
	"teri", //"topical"
	// Case endings
	"ìl", "l", "it", "ti", "t", "ur", "ru", "r", "yä", "ä", "e", "ye", "ìri", "ri",
	// Sorted alphabetically by their reverse forms
	"nemfa", "rofa", "ka", "fa", "na", "ta", "ya", //-a
	"lisre", "pxisre", "sre", "luke", "ne", //-e
	"fpi",          //-i
	"mì",           //-ì
	"lok",          //-k
	"mìkam", "kam", //-m
	"ken", "sìn", //-n
	"äo", "eo", "io", "uo", "ro", "to", //-o
	"tafkip", "takip", "fkip", "kip", //-p
	"ftu", "hu", //-u
	"pximaw", "maw", "pxaw", "few", //-w
	"vay", "kay", //-y
}

var vowelSuffixes = map[string][]string{
	"äo":  []string{"ä", "e"},
	"eo":  []string{"e"},
	"io":  []string{"i"},
	"uo":  []string{"u"},
	"ìlä": []string{"ì"},
	"o":   []string{"o"},
}
var stemSuffixes = []string{"tsyìp", "fkeyk"}
var verbSuffixes = []string{"tswo", "yu"}

var infixes = map[rune][]string{
	rune('a'): {"ay", "asy", "aly", "ary", "am", "alm", "arm", "ats", "awn"},
	rune('ä'): {"äng", "äpeyk", "äp"},
	rune('e'): {"epeyk", "ep", "er", "ei", "eiy", "eng", "eyk"},
	rune('i'): {"iv", "ilv", "irv", "imv", "iyev"},
	rune('ì'): {"ìy", "ìsy", "ìly", "ìry", "ìm", "ìlm", "ìrm", "ìyev"},
	rune('o'): {"ol"},
	rune('u'): {"us", "uy"},
}

var prefirst = []string{"äp", "äpeyk", "ep", "epeyk", "eyk"}
var first = []string{"ay", "asy", "aly", "ary", "ìy", "ìsy", "ìly", "ìry", "ol", "er", "ìm",
	"ìlm", "ìrm", "am", "alm", "arm", "ìyev", "iyev", "iv", "ilv", "irv", "imv", "us", "awn"}
var second = []string{"ei", "eiy", "äng", "eng", "uy", "ats"}

var weirdNounSuffixes = map[string]string{
	// For "tsa" with case endings
	// Canonized in:
	// https://naviteri.org/2011/08/new-vocabulary-clothing/comment-page-1/#comment-912
	"tsa": "tsaw",
	// The a re-appears when case endings are added (it uses a instead of ì)
	"oenga": "oeng",
	// Foreign nouns
	"'ìnglìs":      "'ìnglìsì",
	"keln":         "kelnì",
	"kerìsmìs":     "kerìsmìsì",
	"kìreys":       "kìreysì", // https://naviteri.org/2011/09/miscellaneous-vocabulary/
	"tsìräf":       "tsìräfì",
	"nìyu york":    "nìyu yorkì",
	"nu york":      "nu yorkì", // https://naviteri.org/2013/01/awvea-posti-zisita-amip-first-post-of-the-new-year/
	"päts":         "pätsì",
	"post":         "postì",
	"losäntsyeles": "losäntsyelesì",
}

func isDuplicate(input ConjugationCandidate) bool {
	if a, ok := candidateMap[input.word]; ok {
		if input.insistPOS == a.insistPOS {
			if len(input.prefixes) == len(a.prefixes) && len(input.suffixes) == len(a.suffixes) {
				if len(input.infixes) == len(a.infixes) {
					return true
				}
			}
		}
	}
	return false
}

func isDuplicateFix(fixes []string, fix string) (newFixes []string) {
	if fix == "eng" {
		fix = "äng"
	} else if fix == "ep" {
		fix = "äp"
	} else if fix == "epeyk" {
		fix = "äpeyk"
	} else if fix == "ye" {
		fix = "yä"
	} else if fix == "e" {
		fix = "ä"
	}
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
	d.EN = "Did you mean **" + didYouMean + "**?" // English
	// TODO: Translations
	d.DE = d.EN // German (Deutsch)
	d.ES = d.EN // Spanish (Español)
	d.ET = d.EN // Estonian (Eesti)
	d.FR = d.EN // French (Français)
	d.HU = d.EN // Hungarian (Magyar)
	d.KO = d.EN // Korean (한국어)
	d.NL = d.EN // Dutch (Nederlands)
	d.PL = d.EN // Polish (Polski)
	d.PT = d.EN // Portuguese (Português)
	d.RU = d.EN // Russian (Русский)
	d.SV = d.EN // Swedish (Svenska)
	d.TR = d.EN // Turkish (Türkçe)
	d.UK = d.EN // Ukrainian (Українська)
	d.IPA = ipa
	d.PartOfSpeech = "err."
	return d
}

// fuction to check given string is in array or not
// modified from https://www.golinuxcloud.com/golang-array-contains/
func implContainsAny(sl []string, names []string) bool {
	// iterate over the array and compare given string to each element
	for _, value := range sl {
		for _, name := range names {
			if value == name {
				return true
			}
		}
	}
	return false
}

func deconjugateHelper(input ConjugationCandidate, prefixCheck int, suffixCheck int, unlenite int8, checkInfixes bool, lastPrefix string, lastSuffix string) []ConjugationCandidate {
	if isDuplicate(input) {
		return candidates
	}

	// fneu checking for fne-'u
	if len(lastPrefix) > 0 && len(input.word) > 0 && is_vowel(nth_rune(lastPrefix, -1)) && is_vowel(nth_rune(input.word, 0)) {
		if !implContainsAny(prefixes1lenition, []string{lastPrefix}) { // do not do this for leniting prefixes
			newCandidate := candidateDupe(input)
			newCandidate.word = "'" + newCandidate.word
			deconjugateHelper(newCandidate, prefixCheck, suffixCheck, unlenite, checkInfixes, "", "")
		}
	}

	// fea checkeing for fe'a
	if len(lastSuffix) > 0 && len(input.word) > 0 && is_vowel(nth_rune(lastSuffix, 0)) && is_vowel(nth_rune(input.word, -1)) {
		newCandidate := candidateDupe(input)
		newCandidate.word += "'"
		deconjugateHelper(newCandidate, prefixCheck, suffixCheck, unlenite, checkInfixes, "", "")
	}

	// Exceptions for how words conjugate
	if len(input.suffixes) == 1 {
		if validWord, ok := weirdNounSuffixes[input.word]; ok {
			input.word = validWord
			if !isDuplicate(input) {
				candidates = append(candidates, input)
				candidateMap[input.word] = input
			}
			return candidates
		}
	}

	if len(input.infixes) > 0 && implContainsAny(input.infixes, []string{"ats", "uy"}) {
		// for the cases of zen<ats>eke and zen<uy>eke
		// confirmed in here: https://forum.learnnavi.org/index.php?msg=493217
		if input.word == "zeneke" {
			input.word = "zenke"
			if !isDuplicate(input) {
				candidates = append(candidates, input)
				candidateMap[input.word] = input
			}
			return candidates
		}
	}

	candidates = append(candidates, input)
	candidateMap[input.word] = input

	// Add a way for e to become ä again if we're down to 1 syllable
	if len([]rune(input.word)) < 8 && (len(input.prefixes) > 0 || len(input.infixes) > 0 || len(input.suffixes) > 0) && strings.Contains(input.word, "e") {
		// could be tskxäpx (7 letters 1 syllable)
		newCandidate := candidateDupe(input)
		newCandidate.word = strings.ReplaceAll(newCandidate.word, "e", "ä")
		deconjugateHelper(newCandidate, prefixCheck, suffixCheck, unlenite, checkInfixes, "", "")
	}

	newString := ""

	if input.insistPOS == "n." || input.insistPOS == "any" {
		// For [word] si becoming [word]tswo
		if strings.HasSuffix(input.word, "tswo") {
			newCandidate := candidateDupe(input)
			newCandidate.word = strings.TrimSuffix(input.word, "tswo") + " si"
			newCandidate.insistPOS = "v."
			newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, "tswo")
			if !isDuplicate(newCandidate) {
				candidates = append(candidates, newCandidate)
				candidateMap[input.word] = input
			}
		}
	}

	if input.insistPOS == "adj." || input.insistPOS == "any" {
		// For lrrtok-susi and others
		if strings.HasSuffix(input.word, "-susi") || strings.HasSuffix(input.word, "-susia") {
			found := false
			trimmedWord := strings.TrimSuffix(input.word, "-susi")
			aPosition := 0
			if strings.HasSuffix(input.word, "-susia") {
				trimmedWord = strings.TrimSuffix(input.word, "-susia")
				aPosition = 1
			}

			for _, pairWordSet := range multiword_words[trimmedWord] {
				for _, pairWord := range pairWordSet {
					if pairWord == "si" {
						found = true
						break
					}
				}
				if found {
					break
				}
			}

			if !found && aPosition == 0 && strings.HasPrefix(trimmedWord, "a") {
				noA := strings.TrimPrefix(trimmedWord, "a")
				for _, pairWordSet := range multiword_words[noA] {
					for _, pairWord := range pairWordSet {
						if pairWord == "si" {
							found = true
							break
						}
					}
					if found {
						aPosition = -1
						break
					}
				}
			}

			if !isDuplicate(input) {
				candidates = append(candidates, input)
				candidateMap[input.word] = input
			} // to bump the real candidate into recognition

			if found {
				newCandidate := candidateDupe(input)
				newCandidate.word = trimmedWord + " si"
				if aPosition == -1 {
					newCandidate.word = strings.TrimPrefix(trimmedWord, "a") + " si"
					newCandidate.prefixes = append(newCandidate.prefixes, "a")
				}
				newCandidate.infixes = []string{"us"}
				newCandidate.insistPOS = "v."
				if aPosition == 1 {
					newCandidate.suffixes = append(newCandidate.suffixes, "a")
				}
				if !isDuplicate(newCandidate) {
					candidates = append(candidates, newCandidate)
					candidateMap[input.word] = input
				}
			}
			return candidates
		}
	}

	// Make sure that the first set of prefices (a, nì, ke) aren't combined with suffixes
	newPrefixCheck := prefixCheck
	if newPrefixCheck == 0 {
		newPrefixCheck = 1
	}

	switch prefixCheck {
	case 0:
		if strings.HasPrefix(input.word, "a") && input.insistPOS != "n." && !strings.HasPrefix(input.insistPOS, "ad") {
			// No nouns, adpositions or adverbs
			newCandidate := candidateDupe(input)
			newCandidate.word = input.word[1:]
			newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "a")
			newCandidate.insistPOS = "adj."
			deconjugateHelper(newCandidate, 1, suffixCheck, -1, false, "a", "")
			newCandidate.insistPOS = "v."
			deconjugateHelper(newCandidate, 1, suffixCheck, -1, true, "a", "")
		} else if strings.HasPrefix(input.word, "nì") {
			newCandidate := candidateDupe(input)
			newCandidate.word = strings.TrimPrefix(input.word, "nì")
			newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, "nì")
			newCandidate.insistPOS = "nì."
			// No other affixes allowed
			deconjugateHelper(newCandidate, 10, 10, -1, false, "nì", "") // No other fixes
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
				deconjugateHelper(newCandidate, 10, 10, -1, false, element, "")

				// check "tsatan", "tan" and "atan"
				newCandidate.word = get_last_rune(element, 1) + newCandidate.word
				deconjugateHelper(newCandidate, 10, 10, -1, false, element, "")
			}
		}
		fallthrough
	case 2:
		// Non-lenition prefixes for nouns only
		if input.insistPOS == "any" || input.insistPOS == "n." {
			for _, element := range prefixes1Nouns {
				// If it has a prefix
				if strings.HasPrefix(input.word, element) {
					// remove it
					newString = strings.TrimPrefix(input.word, element)

					newCandidate := candidateDupe(input)
					newCandidate.word = newString
					newCandidate.insistPOS = "n."
					newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
					deconjugateHelper(newCandidate, 3, suffixCheck, -1, false, element, "")

					// check "tsatan", "tan" and "atan"
					newCandidate.word = get_last_rune(element, 1) + newString
					deconjugateHelper(newCandidate, 3, suffixCheck, -1, false, element, "")
				}
			}
		}
		fallthrough
	case 3:
		if input.insistPOS == "any" || input.insistPOS == "n." {
			// This one will demand this makes it use lenition
			for _, element := range prefixes1lenition {
				// If it has a lenition-causing prefix
				if strings.HasPrefix(input.word, element) {
					lenited := false
					newString = strings.TrimPrefix(input.word, element)

					newCandidate := candidateDupe(input)
					newCandidate.word = newString
					newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
					newCandidate.insistPOS = "n."

					// Could it be pekoyu (pe + 'ekoyu, not pe + kxoyu)
					if has("aäeiìou", get_last_rune(element, 1)) {
						// check "pxeyktan", "yktan" and "eyktan"
						newCandidate.word = get_last_rune(element, 1) + newString
						deconjugateHelper(newCandidate, 4, suffixCheck, -1, false, element, "")

						// check "pxeylan", "ylan" and "'eylan"
						newCandidate.word = "'" + newCandidate.word
						deconjugateHelper(newCandidate, 4, suffixCheck, -1, false, element, "")
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
									newCandidate.lenition = []string{newPrefix + "→" + oldPrefix}
								}
								deconjugateHelper(newCandidate, 4, suffixCheck, -1, false, oldPrefix, "")
							}
							break // We don't want the "ts" to become "txs"
						}
					}
					if !lenited {
						newCandidate.word = newString
						deconjugateHelper(newCandidate, 3, suffixCheck, -1, false, element, "")
					}
				}
			}
		}
		fallthrough
	case 4:
		if input.insistPOS == "any" || input.insistPOS == "n." {
			for _, element := range stemPrefixes {
				// If it has a prefix
				if strings.HasPrefix(input.word, element) {
					// remove it
					newCandidate := candidateDupe(input)
					newCandidate.word = strings.TrimPrefix(input.word, element)
					newCandidate.insistPOS = "n."
					newCandidate.prefixes = isDuplicateFix(newCandidate.prefixes, element)
					deconjugateHelper(newCandidate, 5, suffixCheck, -1, false, element, "")

					// check "tsatan", "tan" and "atan"
					newCandidate.word = get_last_rune(element, 1) + newCandidate.word
					deconjugateHelper(newCandidate, 5, suffixCheck, -1, false, element, "")
				}
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
				deconjugateHelper(newCandidate, 10, 10, -1, true, "tì", "") // No other prefixes allowed

				newCandidate.word = "ì" + newCandidate.word
				deconjugateHelper(newCandidate, 10, 10, -1, true, "tì", "") // Or any additional suffixes
			}
		}
	}

	switch suffixCheck {
	case 0:
		// Made sì its own suffix and no suffixes can come after it
		if len(input.suffixes) == 0 && strings.HasSuffix(input.word, "sì") {
			newCandidate := candidateDupe(input)
			newCandidate.word = strings.TrimSuffix(newCandidate.word, "sì")
			newCandidate.suffixes = append(newCandidate.suffixes, "sì")
			deconjugateHelper(newCandidate, newPrefixCheck, 1, unlenite, checkInfixes, "", "sì")
		}
		// special case: short genitives of pronouns like "oey" and "ngey"
		if input.insistPOS == "any" || input.insistPOS == "n." {
			if strings.HasSuffix(input.word, "y") {
				// oey to oe
				newCandidate := candidateDupe(input)
				newCandidate.word = strings.TrimSuffix(input.word, "y")
				newCandidate.insistPOS = "pn."
				newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, "y")
				deconjugateHelper(newCandidate, newPrefixCheck, 10, unlenite, false, "", "y")

				// ngey to nga
				if strings.HasSuffix(newCandidate.word, "e") {
					newCandidate.word = strings.TrimSuffix(newCandidate.word, "e") + "a"
					newCandidate.insistPOS = "pn."
					deconjugateHelper(newCandidate, newPrefixCheck, 10, unlenite, false, "", "y")
				}
			}
		}
		fallthrough
	case 1:
		for _, oldSuffix := range adposuffixes {
			// If it has one of them,
			if strings.HasSuffix(input.word, oldSuffix) {
				newString = strings.TrimSuffix(input.word, oldSuffix)

				newCandidate := candidateDupe(input)
				newCandidate.word = newString
				newCandidate.insistPOS = "n."
				newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
				// all set to 2 to avoid mengeyä -> mengo -> me + 'eng + o
				deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", oldSuffix)

				if oldSuffix == "ä" && !strings.HasSuffix(input.word, "yä") && strings.HasSuffix(input.word, "iä") { // Don't make peyä -> yä -> ya (air)
					// soaiä, tìftiä, etx.
					newString += "a"
					newCandidate.word = newString
					deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", oldSuffix)
				} else if oldSuffix == "e" && !strings.HasSuffix(input.word, "ye") && strings.HasSuffix(input.word, "ie") {
					// reef of above
					newString += "a"
					newCandidate.word = newString
					deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", "ä")
				} else if oldSuffix == "yä" && strings.HasSuffix(newString, "e") {
					// A one-off
					if newString == "tse" {
						newCandidate.word = "tsaw"
						deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", oldSuffix)
					}
					// ngeyä -> nga
					newCandidate.word = strings.TrimSuffix(newString, "e") + "a"
					deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", oldSuffix)
					// oengeyä
					newCandidate.word = strings.TrimSuffix(newString, "e")
					if newCandidate.word == "oeng" { //no mengeyä -> meng -> me + 'eng
						deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", oldSuffix)
					}
					// sneyä -> sno
					newCandidate.word = strings.TrimSuffix(newString, "e") + "o"
					deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", oldSuffix)
				} else if oldSuffix == "ye" && strings.HasSuffix(newString, "e") {
					// reef of above
					if newString == "tse" {
						newCandidate.word = "tsaw"
						deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", "yä")
					}
					// ngeye -> nga
					newCandidate.word = strings.TrimSuffix(newString, "e") + "a"
					deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", "yä")
					// oengeye
					newCandidate.word = strings.TrimSuffix(newString, "e")
					if newCandidate.word == "oeng" { //no mengeyä -> meng -> me + 'eng
						deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", "yä")
					}
					// sneye -> sno
					newCandidate.word = strings.TrimSuffix(newString, "e") + "o"
					deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", "yä")
				} else if vowels, ok := vowelSuffixes["yä"]; ok {
					for _, vowel := range vowels {
						// Make sure zekwä-äo is recognized
						if strings.HasSuffix(newString, vowel+"-") {
							newString = strings.TrimSuffix(newString, "-")
							newCandidate.word = newString
							deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", "yä")
						}
					}
				} else {
					newCandidate.word = strings.TrimSuffix(newString, oldSuffix)
					deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, false, "", oldSuffix)
				}
			}
		}
		fallthrough
	case 2:
		if input.insistPOS == "any" || input.insistPOS == "n." {
			if strings.HasSuffix(input.word, "pe") {
				newString = strings.TrimSuffix(input.word, "pe")

				newCandidate := candidateDupe(input)
				newCandidate.word = newString
				newCandidate.insistPOS = "n."
				newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, "pe")
				deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, false, "", "pe")
			}
		}
		fallthrough
	case 3:
		// If it has one of them,
		if strings.HasSuffix(input.word, "a") && input.insistPOS != "n." && !strings.HasPrefix(input.insistPOS, "ad") {
			// No nouns, adpositions or adverbs
			newString = strings.TrimSuffix(input.word, "a")

			newCandidate := candidateDupe(input)
			newCandidate.word = newString
			newCandidate.insistPOS = "adj."
			newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, "a")
			deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, true, "", "a")
			newCandidate.insistPOS = "v."
			deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, true, "", "a")
		}

		fallthrough
	case 4: // -o suffix "some"
		if input.insistPOS == "any" || input.insistPOS == "n." {
			if strings.HasSuffix(input.word, "o") {
				newString = strings.TrimSuffix(input.word, "o")

				newCandidate := candidateDupe(input)
				newCandidate.word = newString
				newCandidate.insistPOS = "n."
				newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, "o")
				deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, false, "", "o")

				// Make sure fya'o-o is recognized
				if vowels, ok := vowelSuffixes["o"]; ok {
					for _, vowel := range vowels {
						// Make sure fya'o-o is recognized
						if strings.HasSuffix(newString, vowel+"-") {
							newString = strings.TrimSuffix(newString, "-")
							newCandidate.word = newString
							deconjugateHelper(newCandidate, newPrefixCheck, 5, unlenite, false, "", "o")
						}
					}
				}
			}
		}
		fallthrough
	case 5:
		if input.insistPOS == "any" || input.insistPOS == "n." {
			for _, oldSuffix := range stemSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input.word, oldSuffix) {
					newString = strings.TrimSuffix(input.word, oldSuffix)

					//candidates = append(candidates, newString)
					newCandidate := candidateDupe(input)
					newCandidate.word = newString
					newCandidate.insistPOS = "n."
					newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
					deconjugateHelper(newCandidate, newPrefixCheck, 6, unlenite, false, "", oldSuffix)
				}
			}
		}
		fallthrough
	case 6:
		// If it has one of them,
		if input.insistPOS == "any" || input.insistPOS == "n." {
			// verb suffixes change things from verbs to nouns, that's why we check for noun status
			for _, oldSuffix := range verbSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input.word, oldSuffix) {
					newString = strings.TrimSuffix(input.word, oldSuffix)
					newCandidate := candidateDupe(input)
					newCandidate.word = newString
					newCandidate.insistPOS = "v."

					newCandidate.suffixes = isDuplicateFix(newCandidate.suffixes, oldSuffix)
					deconjugateHelper(newCandidate, 10, 10, unlenite, false, "", oldSuffix) // Don't allow any other prefixes
					// They may turn the insistPOS back into a noun

					if oldSuffix == "yu" && strings.HasSuffix(newString, "si") {
						newCandidate.word = strings.TrimSuffix(newString, "si") + " si"
						deconjugateHelper(newCandidate, 10, 10, unlenite, false, "", oldSuffix) // don't allow any other prefixes or suffixes
					}
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
						newCandidate.lenition = []string{newPrefix + "→" + oldPrefix}
					}
					deconjugateHelper(newCandidate, prefixCheck, suffixCheck, -1, false, "", "")
				}
				break // We don't want the "ts" to become "txs"
			}
		}
	}

	if checkInfixes && len(input.infixes) < 3 {
		// Maybe someone else came in with stripped infixes
		if len(input.word) > 2 && input.word[len(input.word)-3] != ' ' &&
			strings.HasSuffix(input.word, "si") && !strings.HasSuffix(input.word, "usi") &&
			!strings.HasSuffix(input.word, "atsi") {
			newCandidate := candidateDupe(input)
			newCandidate.word = strings.TrimSuffix(input.word, "si") + " si"
			newCandidate.insistPOS = "v."
			deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, true, "", "")
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
							deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, true, "", "")

							if infix == "ol" {
								newCandidate := candidateDupe(input)
								newCandidate.word = string(runes[:i]) + "ll" + strings.TrimPrefix(shortString, infix)
								newCandidate.infixes = isDuplicateFix(newCandidate.infixes, infix)
								newCandidate.insistPOS = "v."
								deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, true, "", "")
							} else if infix == "er" {
								newCandidate := candidateDupe(input)
								newCandidate.word = string(runes[:i]) + "rr" + strings.TrimPrefix(shortString, infix)
								newCandidate.infixes = isDuplicateFix(newCandidate.infixes, infix)
								newCandidate.insistPOS = "v."
								deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, true, "", "")
							}
						}
					}
				}
			}
		}
	}
	return candidates
}

func deconjugate(input string) []ConjugationCandidate {
	candidates = []ConjugationCandidate{} //empty array of strings
	candidateMap = map[string]ConjugationCandidate{}
	newCandidate := ConjugationCandidate{}
	newCandidate.word = input
	newCandidate.insistPOS = "any"
	deconjugateHelper(newCandidate, 0, 0, 0, true, "", "")
	candidates = candidates[1:]
	return candidates
}

func TestDeconjugations(searchNaviWord string) (results []Word) {
	conjugations := deconjugate(searchNaviWord)
	for _, candidate := range conjugations {
		a := strings.ReplaceAll(candidate.word, "ù", "u")
		standardizedWordArray := dialectCrunch(strings.Split(a, " "), false)
		a = ""
		for i, b := range standardizedWordArray {
			if i != 0 {
				a += " "
			}
			a += b
		}

		for _, c := range dictHash[a] {
			for _, pos := range strings.Split(c.PartOfSpeech, ",") {
				pos = strings.ReplaceAll(pos, " ", "")

				// An inter. can act like a noun or an adjective, so it gets special treatment
				if pos == "inter." && candidate.insistPOS[0] != 'v' && len(candidate.infixes) == 0 {
					dupe := false
					for _, b := range results {
						if b.Navi == c.Navi {
							dupe = true
							break
						}
					}
					if !dupe {
						a := c
						a.Affixes.Lenition = candidate.lenition
						a.Affixes.Prefix = candidate.prefixes
						a.Affixes.Infix = candidate.infixes
						a.Affixes.Suffix = candidate.suffixes
						results = AppendAndAlphabetize(results, a)
						continue
					}
				}

				// Find gerunds (tì-v<us>erb, treated like a noun)
				gerund := false
				infixBan := false
				doubleBan := false
				attributed := false
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
				if len(candidate.suffixes) < 3 && len(candidate.suffixes) > 0 && candidate.suffixes[0] == "tswo" {
					if pos[0] == 'v' {
						siVerb := false
						if len(candidate.infixes) == 0 {
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
										results = AppendAndAlphabetize(results, a)
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
								results = AppendAndAlphabetize(results, a)
							}
						}
					}
				} else if gerund {
					if pos[0] == 'v' {
						// Make sure the <us> is in the correct place
						rebuiltVerb := strings.ReplaceAll(c.InfixLocations, "<0>", "")
						rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<1>", "us")
						rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<2>", "")

						// Does the noun actually contain the verb?
						noTìftang := strings.TrimPrefix(rebuiltVerb, "'")
						if strings.Contains(searchNaviWord, noTìftang) || strings.Contains(searchNaviWord, dialectCrunch([]string{rebuiltVerb}, false)[0]) {
							a := c
							a.Affixes.Lenition = candidate.lenition
							a.Affixes.Prefix = candidate.prefixes
							a.Affixes.Infix = candidate.infixes
							a.Affixes.Suffix = candidate.suffixes
							results = AppendAndAlphabetize(results, a)
						} else if len(results) == 0 {
							results = AppendAndAlphabetize(results, infixError(searchNaviWord, "tì"+rebuiltVerb, c.IPA))
						}
					}
				} else if candidate.insistPOS == "n." {
					// n., pn., Prop.n. and inter. (but not vin.)
					if len(candidate.infixes) == 0 {
						if (pos[0] != 'v' && strings.HasSuffix(pos, "n.")) || pos == "inter." {
							a := c
							a.Affixes.Lenition = candidate.lenition
							a.Affixes.Prefix = candidate.prefixes
							a.Affixes.Suffix = candidate.suffixes
							results = AppendAndAlphabetize(results, a)
						}
					}
				} else if candidate.insistPOS == "pn." {
					// pn.
					if len(candidate.infixes) == 0 && strings.HasSuffix(pos, "pn.") {
						a := c
						a.Affixes.Lenition = candidate.lenition
						a.Affixes.Prefix = candidate.prefixes
						a.Affixes.Suffix = candidate.suffixes
						results = AppendAndAlphabetize(results, a)
					}
				} else if candidate.insistPOS == "adj." {
					posNoun := pos
					if len(candidate.infixes) == 0 && (posNoun == "adj." || posNoun == "num.") {
						a := c
						a.Affixes.Lenition = candidate.lenition
						a.Affixes.Prefix = candidate.prefixes
						a.Affixes.Suffix = candidate.suffixes
						results = AppendAndAlphabetize(results, a)
					}
				} else if candidate.insistPOS == "v." {
					posNoun := pos
					if strings.HasPrefix(posNoun, "v") {
						// Verbs with -tswo or -yu cannot have infixes
						if len(candidate.suffixes) > 0 {
							for i := len(candidate.suffixes) - 1; i >= 0; i-- {
								if candidate.suffixes[i] == "a" {
									attributed = true
									break
								}
							}
							// Forward search fixs the "a" before "yu" and "tswo"
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

						looseTì := false
						tsuk := false

						if len(candidate.prefixes) > 0 {
							// Reverse search is more likely to find it immediately
							for i := len(candidate.prefixes) - 1; i >= 0; i-- {
								if candidate.prefixes[i] == "a" {
									attributed = true
								} else if candidate.prefixes[i] == "tì" {
									// we found gerunds up top, so this isn't needed
									looseTì = true
									break
								} else {
									for _, j := range verbPrefixes {
										if candidate.prefixes[i] == j {
											if infixBan {
												doubleBan = true
												break
											}
											infixBan = true
											tsuk = true
											break
										}
									}
								}

								if infixBan || doubleBan || looseTì {
									break
								}
							}
						}

						// Don't want a[verb] and [verb]a
						if attributed && (len(candidate.infixes) == 0 || infixBan) && !tsuk {
							continue
						}

						// Take action on tsuk-verb-yus and a-verb-tswos
						if doubleBan || (attributed && !tsuk && infixBan) || looseTì {
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
								results = AppendAndAlphabetize(results, a)
							}
						}

						// Make it verify the infixes are in the correct place
						ol := false
						er := false

						// pre-first position infixes
						rebuiltVerb := c.InfixLocations
						if c.InfixLocations == "z<0><1>en<2>ke" && implContainsAny(candidate.infixes, []string{"ats", "uy"}) {
							rebuiltVerb = "z<0><1>en<2>eke"
						}
						firstInfixes := ""

						for _, newInfix := range candidate.infixes {
							if implContainsAny(prefirst, []string{newInfix}) {
								firstInfixes += newInfix
								rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<0>", firstInfixes)
								if newInfix == "epeyk" || newInfix == "äpeyk" {
									newCandidateInfixes := []string{}
									for _, newInfix2 := range candidate.infixes {
										// äpeyk gets split
										if newInfix2 == "epeyk" || newInfix2 == "äpeyk" {
											newCandidateInfixes = append(newCandidateInfixes, "äp")
											newCandidateInfixes = append(newCandidateInfixes, "eyk")
										} else {
											newCandidateInfixes = append(newCandidateInfixes, newInfix2)
										}
									}
									a.Affixes.Infix = newCandidateInfixes
								}
								break
							}
						}
						rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<0>", "")

						// first position infixes
						firstInfixes = ""
						for _, newInfix := range candidate.infixes {
							if implContainsAny(first, []string{newInfix}) {
								rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<1>", newInfix)
								firstInfixes = newInfix
								if newInfix == "ol" {
									ol = true
								} else if newInfix == "er" {
									er = true
								}
								break
							}
						}
						rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<1>", "")

						// second position infixes
						for _, newInfix := range candidate.infixes {
							if newInfix == "eng" {
								rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<2>", "äng")
								break
							} else if implContainsAny(second, []string{newInfix}) {
								rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<2>", newInfix)
								break
							}
						}
						rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<2>", "")

						rebuiltVerb = strings.TrimSpace(rebuiltVerb)

						if ol && strings.Contains(rebuiltVerb, "olll") {
							rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "olll", "ol")
						}
						if er && strings.Contains(rebuiltVerb, "errr") {
							rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "errr", "er")
						}

						rebuiltVerbForest := rebuiltVerb
						rebuiltVerbArray := dialectCrunch(strings.Split(rebuiltVerb, " "), false)
						rebuiltVerb = ""
						for k, x := range rebuiltVerbArray {
							if k != 0 {
								rebuiltVerb += " "
							}
							rebuiltVerb += x
						}

						if len(candidate.infixes) == 0 || identicalRunes(rebuiltVerb, strings.ReplaceAll(searchNaviWord, "-", " ")) {
							results = AppendAndAlphabetize(results, a)
						} else if participle {
							// In case we have a [word]-susi
							rebuiltHyphen := strings.ReplaceAll(searchNaviWord, "-", " ")
							if identicalRunes("a"+rebuiltVerb, rebuiltHyphen) {
								// a-v<us>erb and a-v<awn>erb
								results = AppendAndAlphabetize(results, a)
							} else if identicalRunes(rebuiltVerb+"a", rebuiltHyphen) {
								// v<us>erb-a and v<awn>erb-a
								results = AppendAndAlphabetize(results, a)
							} else if rebuiltVerb[0] == '\'' && identicalRunes("a"+rebuiltVerb[1:], rebuiltHyphen) {
								// a-'<us>em
								results = AppendAndAlphabetize(results, a)
							} else if rebuiltVerb[len(rebuiltVerb)-1] == '\'' && identicalRunes(rebuiltVerb[:len(rebuiltVerb)-1]+"a", rebuiltHyphen) {
								// fp<us>e'a
								results = AppendAndAlphabetize(results, a)
							} else if firstInfixes == "us" {
								if len(results) == 0 {
									results = AppendAndAlphabetize(results, infixError(searchNaviWord, rebuiltVerbForest, c.IPA))
								}
							}
						} else if gerund { // ti is needed to weed out non-productive tì-verbs
							if len(results) == 0 {
								results = AppendAndAlphabetize(results, infixError(searchNaviWord, rebuiltVerbForest, c.IPA))
							}
						} else {
							if len(results) == 0 {
								results = AppendAndAlphabetize(results, infixError(searchNaviWord, rebuiltVerbForest, c.IPA))
							}
						}
					}
				} else if candidate.insistPOS == "nì." {
					posNoun := pos
					if len(candidate.infixes) == 0 && (posNoun == "adj." || posNoun == "pn.") {
						a := c
						a.Affixes.Lenition = candidate.lenition
						a.Affixes.Prefix = candidate.prefixes
						a.Affixes.Suffix = candidate.suffixes
						results = AppendAndAlphabetize(results, a)
					}
				} else if len(candidate.infixes) == 0 {
					a := c
					a.Affixes.Lenition = candidate.lenition
					a.Affixes.Prefix = candidate.prefixes
					a.Affixes.Suffix = candidate.suffixes
					results = AppendAndAlphabetize(results, a)
				}
			}
		}
	}
	return
}
