package fwew_lib

import (
	"math"
	"slices"
	"strings"
	//"fmt"
)

type ConjugationCandidate struct {
	Word      string
	Lenition  []string
	Prefixes  []string
	Suffixes  []string
	Infixes   []string
	InsistPOS string
}

func candidateDupe(candidate ConjugationCandidate) (c ConjugationCandidate) {
	a := ConjugationCandidate{}
	a.Word = candidate.Word
	a.Lenition = candidate.Lenition
	a.Prefixes = candidate.Prefixes
	a.Infixes = candidate.Infixes
	a.Suffixes = candidate.Suffixes
	a.InsistPOS = candidate.InsistPOS
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

var prefixes1Nouns = []string{"fì", "tsa", "fi"}
var prefixes1NounsLenition = []string{"pay", "fay"}
var prefixes1lenition = []string{"pxe", "ay", "me"}
var stemPrefixes = []string{"fne", "sna", "munsna"}
var verbPrefixes = []string{"tsuk", "ketsuk"}
var caseEndings = map[string]bool{
	"ìl":  true,
	"l":   true,
	"it":  true,
	"ti":  true,
	"t":   true,
	"ur":  true,
	"ru":  true,
	"r":   true,
	"yä":  true,
	"ä":   true,
	"ìri": true,
	"ri":  true,
	"ye":  true,
	"e":   true,
	"il":  true,
	"iri": true,
}

var adposuffixes = []string{
	// adpositions that can be mistaken for case endings
	"pxel",                                                     //"agentive"
	"mungwrr",                                                  //"dative"
	"kxamlä", "ìlä", "wä", "nuä", "kxamle", "ìle", "we", "nue", //"genitive"
	"ilä", "ile", //"genitive"
	"teri", //"topical"
	// Case endings
	"ìl", "l", "it", "ti", "t", "ur", "ru", "r", "yä", "ä", "ìri", "ri", "ye", "e", "il", "iri",
	// Sorted alphabetically by their reverse forms
	"ftumfa", "nemfa", "rofa", "ka", "fa", "na", "ta", "ya", "yoa", "krrka", "ftuopa", //-a
	"lisre", "pxisre", "sre", "luke", "ne", //-e
	"fpi", "mi", //-i
	"mì",                    //-ì
	"lok",                   //-k
	"mìkam", "kam", "mikam", //-m
	"ken", "sìn", "talun", "sin", //-n
	"äo", "eo", "io", "uo", "ro", "to", "sko", //-o
	"tafkip", "takip", "fkip", "kip", //-p
	"ftu", "hu", "sru", //-u
	"pximaw", "maw", "pxaw", "few", "raw", //-w
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
var stemSuffixes = []string{"tsyìp", "tsyip", "fkeyk"}
var verbSuffixes = []string{"tswo", "yu", "tseng"}

var infixes = map[rune][]string{
	rune('a'): {"ay", "asy", "aly", "ary", "am", "alm", "arm", "ats", "awn", "ap", "ang"},
	rune('ä'): {"äng", "äpeyk", "äp"},
	rune('e'): {"epeyk", "ep", "eng", "er", "ei", "eiy", "eyk"},
	rune('i'): {"iy", "isy", "ily", "iry", "im", "ilm", "irm", "iyev", "iv", "ilv", "irv", "imv"},
	rune('ì'): {"ìy", "ìsy", "ìly", "ìry", "ìm", "ìlm", "ìrm", "ìyev"},
	rune('o'): {"ol"},
	rune('u'): {"us", "uy"},
}

var prefirst = []string{"äp", "äpeyk", "eyk", "epeyk", "ep"}
var first = []string{"ay", "asy", "aly", "ary", "ìy", "iy", "ìsy", "ìly", "ìry", "ol", "er", "ìm", "im",
	"ìlm", "ìrm", "am", "alm", "arm", "ìyev", "iyev", "iv", "ilv", "irv", "imv", "us", "awn", "isy", "ily", "iry", "ilm", "irm"}
var second = []string{"ei", "eiy", "äng", "eng", "ang", "uy", "ats"}

var prefirstMap = map[string]bool{"äp": true, "äpeyk": true, "eyk": true, "ep": true, "epeyk": true}
var firstMap = map[string]bool{"ay": true, "asy": true, "aly": true, "ary": true, "ìy": true, "iy": true, "ìsy": true,
	"ìly": true, "ìry": true, "ol": true, "er": true, "ìm": true, "im": true, "ìlm": true,
	"ìrm": true, "am": true, "alm": true, "arm": true, "ìyev": true, "iyev": true,
	"iv": true, "ilv": true, "irv": true, "imv": true, "us": true, "awn": true, "isy": true, "ily": true, "iry": true, "ilm": true, "irm": true}
var secondMap = map[string]bool{"ei": true, "eiy": true, "äng": true, "uy": true, "ats": true, "ap": true, "ang": true, "eng": true}

var unreefFixes = map[string]string{
	"eng":    "äng",
	"epeyk":  "äpeyk",
	"ep":     "äp",
	"ye":     "yä",
	"e":      "ä",
	"we":     "wä",
	"ìle":    "ìlä",
	"nue":    "nuä",
	"kxamle": "kxamlä",
}

var unstrictFixes = map[string]string{
	// Below this are for diacritic independence of i and ì
	"fi":    "fì",
	"ilä":   "ìlä",
	"mi":    "mì",
	"mikam": "mìkam",
	"sin":   "sìn",
	"il":    "ìl",
	"iri":   "ìri",
	"ap":    "äp",
	"ang":   "äng",
	"iy":    "ìy",
	"im":    "ìm",
	"ilm":   "ìlm",
	"isy":   "ìsy",
	"tsyip": "tsyìp",
}

var weirdNounSuffixes = map[string]string{
	// For "tsa" with case endings
	// Canonized in:
	// https://naviteri.org/2011/08/new-vocabulary-clothing/comment-page-1/#comment-912
	"tsa":   "tsaw",
	"teyng": "tì'eyng",
	// The a re-appears when case endings are added (it uses a instead of ì)
	"oenga":   "oeng",
	"ayoenga": "ayoeng",
	"pxoenga": "pxoeng",
	// Foreign nouns
	"'ìnglìs":      "'ìnglìsì",
	"'inglis":      "'ìnglìsì",
	"keln":         "kelnì",
	"kerìsmìs":     "kerìsmìsì",
	"kerismis":     "kerìsmìsì",
	"kìreys":       "kìreysì", // https://naviteri.org/2011/09/miscellaneous-vocabulary/
	"kireys":       "kìreysì",
	"tsìräf":       "tsìräfì",
	"tsiräf":       "tsìräfì",
	"nìyu york":    "nìyu yorkì",
	"niyu york":    "nìyu yorkì",
	"nu york":      "nu yorkì", // https://naviteri.org/2013/01/awvea-posti-zisita-amip-first-post-of-the-new-year/
	"päts":         "pätsì",
	"post":         "postì",
	"losäntsyeles": "losäntsyelesì",
	"york":         "yorkì", // For a program called Litxap
}

var forbiddenTsa = [][]string{{"fì", "tsa", "fay", "tsay", "ay", "pe"}, {}, {"pe"}}
var forbiddenEyk = [][]string{{}, {"äp", "eyk"}, {}}
var forbiddenAy = [][]string{{"fay", "tsay", "ay", "pe"}, {}, {}}
var forbiddenTsyìp = [][]string{{}, {}, {"tsyìp"}}
var forbiddenTsaw = [][]string{{"fì", "tsa", "fi", "pay", "fay", "pxe", "ay", "me"}, {}, slices.Concat(stemSuffixes,
	[]string{"ìl", "it", "ur", "ä", "ìri", "e", "il", "iri"})}
var forbiddenTsat = [][]string{{"fì", "tsa", "fi", "pay", "fay", "pxe", "ay", "me"}, {}, slices.Concat(adposuffixes, stemSuffixes,
	[]string{"ìl", "l", "it", "ti", "t", "ur", "ru", "r", "yä", "ä", "ìri", "ri", "ye", "e", "il", "iri"})}

// This cannot be autogenerated since there's sängop, kxal and other fake ones
// Pe is in the beginning and end so peupe won't appear
var productiveCompounds = map[string][][]string{
	"fìtseng":        forbiddenTsa,
	"tsatseng":       forbiddenTsa,
	"'upe":           forbiddenTsa,
	"ayfo":           forbiddenAy,
	"aynga":          forbiddenAy,
	"ayoe":           forbiddenAy,
	"ayoeng":         forbiddenAy,
	"fìpo":           forbiddenTsa,
	"fìkem":          forbiddenTsa,
	"fì'u'":          forbiddenTsa,
	"kempe":          forbiddenTsa,
	"krrpe":          forbiddenTsa,
	"pefya":          forbiddenTsa,
	"pehem":          forbiddenTsa,
	"pehrr":          forbiddenTsa,
	"pelun":          forbiddenTsa,
	"peseng":         forbiddenTsa,
	"peu":            forbiddenTsa,
	"tsa'u":          forbiddenTsa,
	"tsengpe":        forbiddenTsa,
	"tsakrr":         forbiddenTsa,
	"steyki":         forbiddenEyk,
	"fratseng":       forbiddenTsa,
	"fratrr":         forbiddenTsa,
	"tsengo":         forbiddenTsa,
	"holpxaype":      forbiddenTsa,
	"hìmtxampe":      forbiddenTsa,
	"lì'upe":         forbiddenTsa,
	"pelì'u":         forbiddenTsa,
	"fìtrr":          forbiddenTsa,
	"fìtxon":         forbiddenTsa,
	"ayu":            forbiddenAy,
	"'itetsyìp":      forbiddenTsyìp,
	"sa'nutsyìp":     forbiddenTsyìp,
	"säspxintsyìp":   forbiddenTsyìp,
	"utraltsyìp":     forbiddenTsyìp,
	"puktsyìp":       forbiddenTsyìp,
	"tswintsyìp":     forbiddenTsyìp,
	"taronyutsyìp":   forbiddenTsyìp,
	"oetsyìp":        forbiddenTsyìp,
	"ngatsyìp":       forbiddenTsyìp,
	"txeptsyìp":      forbiddenTsyìp,
	"tìpängkxotsyìp": forbiddenTsyìp,
	"ramtsyìp":       forbiddenTsyìp,
	"swawtsyìp":      forbiddenTsyìp,
	"skxirtsyìp":     forbiddenTsyìp,
	"tsongtsyìp":     forbiddenTsyìp,
	"reykol":         forbiddenEyk,
	"srungtsyìp":     forbiddenTsyìp,
	"vezeyko":        forbiddenEyk,
	"zeyko":          forbiddenEyk,
	"ingyentsyìp":    forbiddenTsyìp,
	"späpeng":        forbiddenEyk,
	"tsyeytsyìp":     forbiddenTsyìp,
	"nantangtsyìp":   forbiddenTsyìp,
	"vultsyìp":       forbiddenTsyìp,
	"'opinvultsyìp":  forbiddenTsyìp,
	"pela'a":         forbiddenTsa,
	"la'ape":         forbiddenTsa,
	"pelìmsim":       forbiddenTsa,
	"lìmsimpe":       forbiddenTsa,
	"aysupe":         forbiddenTsa,
	"pxesupe":        forbiddenTsa,
	"tutepe":         forbiddenTsa,
	"tstunkemtsyìp":  forbiddenTsyìp,
	"leykek":         forbiddenEyk,
	"penunyol":       forbiddenTsa,
	"nunyolpe":       forbiddenTsa,
	"pengimpup":      forbiddenTsa,
	"ngimpuppe":      forbiddenTsa,
	"säfleltsyìp":    forbiddenTsyìp,
	"mawuptsyìp":     forbiddenTsyìp,
	"trrpxìvitsyìp":  forbiddenTsyìp,
	"pesrrpxì":       forbiddenTsa,
	"trrpxìpe":       forbiddenTsa,
	"pehrrlik":       forbiddenTsa,
	"krrlikpe":       forbiddenTsa,
	"pamtsyìp":       forbiddenTsyìp,
	"tsaw":           forbiddenTsaw,
	"tsal":           forbiddenTsat,
	"tsat":           forbiddenTsat,
	"tsar":           forbiddenTsat,
	"tsari":          forbiddenTsat,
	"tseyä":          forbiddenTsaw,
}

func isDuplicate(input ConjugationCandidate) bool {
	if a, ok := candidateMap[input.Word]; ok {
		if input.InsistPOS == a.InsistPOS {
			if len(input.Prefixes) == len(a.Prefixes) && len(input.Suffixes) == len(a.Suffixes) {
				if len(input.Infixes) == len(a.Infixes) {
					return true
				}
			}
		}
	}
	return false
}

func isDuplicateFix(fixes []string, fix string, strict bool, allowReef bool) (newFixes []string, added bool) {
	if newfix, ok := unreefFixes[fix]; ok {
		if !allowReef {
			return fixes, false
		}
		fix = newfix
	}
	if newfix, ok := unstrictFixes[fix]; ok {
		if strict {
			return fixes, false
		}
		fix = newfix
	}
	if fix == "ile" {
		if !strict && allowReef {
			fix = "ìlä"
		} else {
			return fixes, false
		}
	}
	for _, a := range fixes {
		if fix == a {
			return fixes, false
		}
	}
	fixes = append(fixes, fix)
	return fixes, true
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

// Helper for infix detection
func verifyInfix(existing []string, new string) (bool, []string) {
	if _, ok := prefirstMap[new]; ok {
		if existing[0] == "" {
			return true, []string{new, existing[1], existing[2]}
		}
	} else if _, ok := firstMap[new]; ok {
		if existing[1] == "" {
			return true, []string{existing[0], new, existing[2]}
		}
	} else if _, ok := secondMap[new]; ok {
		if existing[2] == "" {
			return true, []string{existing[0], existing[1], new}
		}
	}

	return false, existing
}

func verifyCaseEnding(noun string, ending string) bool {
	// error prevention
	if len(noun) == 0 {
		return false
	}

	if get_last_rune(noun, 1) == 'i' && (ending == "ä" || ending == "e") {
		//soaiä, tìftiä
		return true
	}
	// Don't check adpositions
	if _, ok := caseEndings[ending]; !ok {
		return true
	}
	// Non-standard conjugations
	if noun == "omatikaya" && ending == "ä" {
		return true
	}
	diphthongs := map[string]bool{
		"ay": true,
		"aw": true,
		"ey": true,
		"ew": true,
	}
	vowels := map[string]bool{
		"a": true,
		"ä": true,
		"e": true,
		"é": true,
		"i": true,
		"ì": true,
		"o": true,
		"u": true,
		"ù": true,
	}
	nounEnding := ""
	if len(noun) >= 2 {
		nounEnding = noun[len(noun)-2:]
	}
	if _, ok := diphthongs[nounEnding]; ok {
		nounEnding := noun[len(noun)-2:]
		//ewur isn't valid
		if nounEnding == "ew" && ending == "ur" {
			return false
		}
		// Diphthong
		diphthongEndings := map[string]bool{
			"ìl": true,
			"ti": true,
			"it": true,
			"ru": true,
			"ur": true,
			"ä":  true,
			"e":  true,
			"ri": true,
		}
		if _, ok := diphthongEndings[ending]; ok {
			return true
		} else {
			lastRune := get_last_rune(noun, 1)
			switch lastRune {
			case 'y':
				// ayt, eyt
				if ending == "t" {
					return true
				}
			case 'w':
				// ewr, awr
				if ending == "r" {
					return true
				}
			}
		}
	} else if _, ok := vowels[string(get_last_rune(noun, 1))]; ok {
		lastVowel := get_last_rune(noun, 1)
		if lastVowel == 'u' || lastVowel == 'o' {
			// No oyä or ayä
			if ending == "yä" || ending == "ye" {
				return false
			}
		}
		vowelEndings := map[string]bool{
			"l":  true,
			"t":  true,
			"ti": true,
			"ru": true,
			"r":  true,
			"yä": true,
			"ä":  true,
			"ye": true,
			"e":  true,
			"ri": true,
		}
		if _, ok := vowelEndings[ending]; ok {
			return true
		}
	} else {
		// Consonant or psuedovowel
		otherEndings := map[string]bool{
			"ìl":  true,
			"ti":  true,
			"it":  true,
			"ur":  true,
			"ä":   true,
			"e":   true,
			"ìri": true,
		}
		if _, ok := otherEndings[ending]; ok {
			return true
		}
		//'ri
		if get_last_rune(noun, 1) == '\'' && ending == "ru" {
			return true
		}
	}
	return false
}

func deconjugateHelper(input ConjugationCandidate, prefixCheck int, suffixCheck int, unlenite int8,
	infix []string, lastPrefix string, lastSuffix string, strict bool, allowReef bool) []ConjugationCandidate {
	if isDuplicate(input) {
		return candidates
	}

	vowels := "aäeiìouù"

	added := false

	if allowReef {
		// fneu checking for fne-'u
		if len(lastPrefix) > 0 && len(input.Word) > 0 && hasAt(vowels, lastPrefix, -1) && hasAt(vowels, input.Word, 0) {
			if !implContainsAny(prefixes1lenition, []string{lastPrefix}) { // do not do this for leniting prefixes
				newCandidate := candidateDupe(input)
				newCandidate.Word = "'" + newCandidate.Word
				deconjugateHelper(newCandidate, prefixCheck, suffixCheck, unlenite, infix, "", "", strict, allowReef)
			}
		}

		// fea checkeing for fe'a
		if len(lastSuffix) > 0 && len(input.Word) > 0 && hasAt(vowels, lastSuffix, 0) && hasAt(vowels, input.Word, -1) {
			newCandidate := candidateDupe(input)
			newCandidate.Word += "'"
			deconjugateHelper(newCandidate, prefixCheck, suffixCheck, unlenite, infix, "", "", strict, allowReef)
		}
	}

	// Exceptions for how words conjugate
	if len(input.Suffixes) == 1 {
		if validWord, ok := weirdNounSuffixes[input.Word]; ok {
			input.Word = validWord
			if !isDuplicate(input) {
				candidates = append(candidates, input)
				candidateMap[input.Word] = input
			}
			return candidates
		}
	}

	if len(input.Infixes) > 0 && implContainsAny(input.Infixes, []string{"ats", "uy"}) {
		// for the cases of zen<ats>eke and zen<uy>eke
		// confirmed in here: https://forum.learnnavi.org/index.php?msg=493217
		if input.Word == "zeneke" {
			input.Word = "zenke"
			if !isDuplicate(input) {
				candidates = append(candidates, input)
				candidateMap[input.Word] = input
			}
			return candidates
		}
	}

	candidates = append(candidates, input)
	candidateMap[input.Word] = input

	// Add a way for e to become ä again if we're down to 1 syllable
	if !strict && allowReef && len([]rune(input.Word)) < 8 && (len(input.Prefixes) > 0 ||
		len(input.Infixes) > 0 || len(input.Suffixes) > 0) && strings.Contains(input.Word, "e") {
		// could be tskxäpx (7 letters 1 syllable)
		newCandidate := candidateDupe(input)
		newCandidate.Word = strings.ReplaceAll(newCandidate.Word, "e", "ä")
		deconjugateHelper(newCandidate, prefixCheck, suffixCheck, unlenite, infix, "", "", strict, allowReef)
	}

	newString := ""

	if input.InsistPOS == "n." || input.InsistPOS == "any" {
		// For [word] si becoming [word]tswo
		if strings.HasSuffix(input.Word, "tswo") {
			newCandidate := candidateDupe(input)
			newCandidate.Word = strings.TrimSuffix(input.Word, "tswo") + " si"
			newCandidate.InsistPOS = "v."
			newCandidate.Suffixes, added = isDuplicateFix(newCandidate.Suffixes, "tswo", strict, allowReef)
			if added && !isDuplicate(newCandidate) {
				candidates = append(candidates, newCandidate)
				candidateMap[input.Word] = input
			}
		}
	}

	if input.InsistPOS == "adj." || input.InsistPOS == "any" {
		// For lrrtok-susi and others
		if strings.HasSuffix(input.Word, "-susi") || strings.HasSuffix(input.Word, "-susia") {
			found := false
			trimmedWord := strings.TrimSuffix(input.Word, "-susi")
			aPosition := 0
			if strings.HasSuffix(input.Word, "-susia") {
				trimmedWord = strings.TrimSuffix(input.Word, "-susia")
				aPosition = 1
			}

			if strict {
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
			} else {
				for _, pairWordSet := range multiword_words_loose[trimmedWord] {
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
			}

			if !found && aPosition == 0 && strings.HasPrefix(trimmedWord, "a") {
				noA := strings.TrimPrefix(trimmedWord, "a")

				if strict {
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
				} else {
					for _, pairWordSet := range multiword_words_loose[noA] {
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
			}

			if !isDuplicate(input) {
				candidates = append(candidates, input)
				candidateMap[input.Word] = input
			} // to bump the real candidate into recognition

			if found {
				newCandidate := candidateDupe(input)
				newCandidate.Word = trimmedWord + " si"
				if aPosition == -1 {
					newCandidate.Word = strings.TrimPrefix(trimmedWord, "a") + " si"
					newCandidate.Prefixes = append(newCandidate.Prefixes, "a")
				}
				newCandidate.Infixes = []string{"us"}
				newCandidate.InsistPOS = "v."
				if aPosition == 1 {
					newCandidate.Suffixes = append(newCandidate.Suffixes, "a")
				}
				if !isDuplicate(newCandidate) {
					candidates = append(candidates, newCandidate)
					candidateMap[input.Word] = input
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
		if strings.HasPrefix(input.Word, "a") && input.InsistPOS != "n." && !strings.HasPrefix(input.InsistPOS, "ad") {
			// No nouns, adpositions or adverbs
			newCandidate := candidateDupe(input)
			newCandidate.Word = input.Word[1:]
			newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, "a", strict, allowReef)
			if added {
				newCandidate.InsistPOS = "adj."
				deconjugateHelper(newCandidate, 1, suffixCheck, -1, []string{}, "a", "", strict, allowReef)
				newCandidate.InsistPOS = "v."
				deconjugateHelper(newCandidate, 1, suffixCheck, -1, []string{"", "", ""}, "a", "", strict, allowReef)
			}
		} else if strings.HasPrefix(input.Word, "nì") {
			newCandidate := candidateDupe(input)
			newCandidate.Word = strings.TrimPrefix(input.Word, "nì")
			newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, "nì", strict, allowReef)
			if added {
				newCandidate.InsistPOS = "nì."
				// No other affixes allowed
				deconjugateHelper(newCandidate, 10, 10, -1, []string{}, "nì", "", strict, allowReef) // No other fixes
			}
		} else if !strict && strings.HasPrefix(input.Word, "ni") {
			newCandidate := candidateDupe(input)
			newCandidate.Word = strings.TrimPrefix(input.Word, "ni")
			newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, "nì", strict, allowReef)
			if added {
				newCandidate.InsistPOS = "nì."
				// No other affixes allowed
				deconjugateHelper(newCandidate, 10, 10, -1, []string{}, "nì", "", strict, allowReef) // No other fixes
			}
		}
		fallthrough
	case 1:
		for _, element := range verbPrefixes {
			// If it has a prefix
			if strings.HasPrefix(input.Word, element) {
				// remove it
				newCandidate := candidateDupe(input)
				newCandidate.Word = strings.TrimPrefix(input.Word, element)
				newCandidate.InsistPOS = "v."
				newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, element, strict, allowReef)
				if !added {
					continue
				}
				deconjugateHelper(newCandidate, 10, 10, -1, []string{}, element, "", strict, allowReef)

				// check "tsatan", "tan" and "atan"
				newCandidate.Word = string(get_last_rune(element, 1)) + newCandidate.Word
				deconjugateHelper(newCandidate, 10, 10, -1, []string{}, element, "", strict, allowReef)
			}
		}
		fallthrough
	case 2:
		// Non-lenition prefixes for nouns only
		if input.InsistPOS == "any" || input.InsistPOS == "n." {
			for _, element := range prefixes1Nouns {
				// If it has a prefix
				if strings.HasPrefix(input.Word, element) {
					// remove it
					newString = strings.TrimPrefix(input.Word, element)

					newCandidate := candidateDupe(input)
					newCandidate.Word = newString
					newCandidate.InsistPOS = "n."
					newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, element, strict, allowReef)
					if !added {
						continue
					}
					deconjugateHelper(newCandidate, 3, suffixCheck, -1, []string{}, element, "", strict, allowReef)

					// check "tsatan", "tan" and "atan"
					newCandidate.Word = string(get_last_rune(element, 1)) + newString
					deconjugateHelper(newCandidate, 3, suffixCheck, -1, []string{}, element, "", strict, allowReef)
				}
			}
		}

		// This one will demand this makes it use lenition
		for _, element := range prefixes1NounsLenition {
			// If it has a lenition-causing prefix
			if strings.HasPrefix(input.Word, element) {
				lenited := false
				newString = strings.TrimPrefix(input.Word, element)

				newCandidate := candidateDupe(input)
				newCandidate.Word = newString
				newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, element, strict, allowReef)
				if !added {
					continue
				}
				newCandidate.InsistPOS = "n."

				// Could it be pekoyu (pe + 'ekoyu, not pe + kxoyu)
				if hasAt(vowels, element, -1) {
					// check "pxeyktan", "yktan" and "eyktan"
					newCandidate.Word = string(get_last_rune(element, 1)) + newString
					deconjugateHelper(newCandidate, 5, suffixCheck, -1, []string{}, element, "", strict, allowReef)

					// check "pxeylan", "ylan" and "'eylan"
					newCandidate.Word = "'" + newCandidate.Word
					deconjugateHelper(newCandidate, 5, suffixCheck, -1, []string{}, element, "", strict, allowReef)
				}

				// find out the possible unlenited forms
				for _, oldPrefix := range unlenitionLetters {
					// If it has a letter that could have changed for lenition,
					if strings.HasPrefix(newString, oldPrefix) {
						// put all possibilities in the candidates
						lenited = true

						for _, newPrefix := range unlenition[oldPrefix] {
							newCandidate.Word = newPrefix + strings.TrimPrefix(newString, oldPrefix)
							if oldPrefix != newPrefix {
								newCandidate.Lenition = []string{newPrefix + "→" + oldPrefix}
							}
							deconjugateHelper(newCandidate, 5, suffixCheck, -1, []string{}, oldPrefix, "", strict, allowReef)
						}
						break // We don't want the "ts" to become "txs"
					}
				}
				if !lenited {
					newCandidate.Word = newString
					deconjugateHelper(newCandidate, 5, suffixCheck, -1, []string{}, element, "", strict, allowReef)
				}
			}
		}

		if input.InsistPOS == "any" || input.InsistPOS == "n." {
			// This one will demand this makes it use lenition
			if strings.HasPrefix(input.Word, "pe") {
				lenited := false
				newString = strings.TrimPrefix(input.Word, "pe")

				newCandidate := candidateDupe(input)
				newCandidate.Word = newString
				newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, "pe", strict, allowReef)
				if added {
					newCandidate.InsistPOS = "n."

					// Could it be pekoyu (pe + 'ekoyu, not pe + kxoyu)
					if hasAt(vowels, "pe", -1) {
						// check "pxeyktan", "yktan" and "eyktan"
						newCandidate.Word = string(get_last_rune("pe", 1)) + newString
						deconjugateHelper(newCandidate, 3, suffixCheck, -1, []string{}, "pe", "", strict, allowReef)

						// check "pxeylan", "ylan" and "'eylan"
						newCandidate.Word = "'" + newCandidate.Word
						deconjugateHelper(newCandidate, 3, suffixCheck, -1, []string{}, "pe", "", strict, allowReef)
					}

					// find out the possible unlenited forms
					for _, oldPrefix := range unlenitionLetters {
						// If it has a letter that could have changed for lenition,
						if strings.HasPrefix(newString, oldPrefix) {
							// put all possibilities in the candidates
							lenited = true

							for _, newPrefix := range unlenition[oldPrefix] {
								newCandidate.Word = newPrefix + strings.TrimPrefix(newString, oldPrefix)
								if oldPrefix != newPrefix {
									newCandidate.Lenition = []string{newPrefix + "→" + oldPrefix}
								}
								deconjugateHelper(newCandidate, 3, suffixCheck, -1, []string{}, oldPrefix, "", strict, allowReef)
							}
							break // We don't want the "ts" to become "txs"
						}
					}
					if !lenited {
						newCandidate.Word = newString
						deconjugateHelper(newCandidate, 3, suffixCheck, -1, []string{}, "pe", "", strict, allowReef)
					}
				}
			}
		}

		fallthrough
	case 3:
		if input.InsistPOS == "any" || input.InsistPOS == "n." {
			// If it has a prefix
			if strings.HasPrefix(input.Word, "fra") {
				// remove it
				newString = strings.TrimPrefix(input.Word, "fra")

				newCandidate := candidateDupe(input)
				newCandidate.Word = newString
				newCandidate.InsistPOS = "n."
				newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, "fra", strict, allowReef)
				if added {
					deconjugateHelper(newCandidate, 4, suffixCheck, -1, []string{}, "fra", "", strict, allowReef)

					// check "tsatan", "tan" and "atan"
					newCandidate.Word = string(get_last_rune("fra", 1)) + newString
					deconjugateHelper(newCandidate, 4, suffixCheck, -1, []string{}, "fra", "", strict, allowReef)
				}
			}
		}
		fallthrough
	case 4:
		if input.InsistPOS == "any" || input.InsistPOS == "n." {
			// This one will demand this makes it use lenition
			for _, element := range prefixes1lenition {
				// If it has a lenition-causing prefix
				if strings.HasPrefix(input.Word, element) {
					lenited := false
					newString = strings.TrimPrefix(input.Word, element)

					newCandidate := candidateDupe(input)
					newCandidate.Word = newString
					newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, element, strict, allowReef)
					if !added {
						continue
					}
					newCandidate.InsistPOS = "n."

					// Could it be pekoyu (pe + 'ekoyu, not pe + kxoyu)
					if hasAt(vowels, element, -1) {
						// check "pxeyktan", "yktan" and "eyktan"
						newCandidate.Word = string(get_last_rune(element, 1)) + newString
						deconjugateHelper(newCandidate, 5, suffixCheck, -1, []string{}, element, "", strict, allowReef)

						// check "pxeylan", "ylan" and "'eylan"
						newCandidate.Word = "'" + newCandidate.Word
						deconjugateHelper(newCandidate, 5, suffixCheck, -1, []string{}, element, "", strict, allowReef)
					}

					// find out the possible unlenited forms
					for _, oldPrefix := range unlenitionLetters {
						// If it has a letter that could have changed for lenition,
						if strings.HasPrefix(newString, oldPrefix) {
							// put all possibilities in the candidates
							lenited = true

							for _, newPrefix := range unlenition[oldPrefix] {
								newCandidate.Word = newPrefix + strings.TrimPrefix(newString, oldPrefix)
								if oldPrefix != newPrefix {
									newCandidate.Lenition = []string{newPrefix + "→" + oldPrefix}
								}
								deconjugateHelper(newCandidate, 5, suffixCheck, -1, []string{}, oldPrefix, "", strict, allowReef)
							}
							break // We don't want the "ts" to become "txs"
						}
					}
					if !lenited {
						newCandidate.Word = newString
						deconjugateHelper(newCandidate, 5, suffixCheck, -1, []string{}, element, "", strict, allowReef)
					}
				}
			}
		}
		fallthrough
	case 5:
		if input.InsistPOS == "any" || input.InsistPOS == "n." {
			for _, element := range stemPrefixes {
				// If it has a prefix
				if strings.HasPrefix(input.Word, element) {
					// remove it
					newCandidate := candidateDupe(input)
					newCandidate.Word = strings.TrimPrefix(input.Word, element)
					newCandidate.InsistPOS = "n."
					newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, element, strict, allowReef)
					if !added {
						continue
					}
					deconjugateHelper(newCandidate, 6, suffixCheck, -1, []string{}, element, "", strict, allowReef)

					// check "tsatan", "tan" and "atan"
					newCandidate.Word = string(get_last_rune(element, 1)) + newCandidate.Word
					deconjugateHelper(newCandidate, 6, suffixCheck, -1, []string{}, element, "", strict, allowReef)
				}
			}
		}
		fallthrough
	case 6:
		if strings.HasPrefix(input.Word, "tì") {
			if input.InsistPOS == "any" || input.InsistPOS == "n." {
				// remove it
				newCandidate := candidateDupe(input)
				newCandidate.Word = strings.TrimPrefix(input.Word, "tì")
				newCandidate.InsistPOS = "v."
				newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, "tì", strict, allowReef)
				if added {
					deconjugateHelper(newCandidate, 10, 10, -1, []string{"", "", ""}, "tì", "", strict, allowReef) // No other prefixes allowed

					newCandidate.Word = "ì" + newCandidate.Word
					deconjugateHelper(newCandidate, 10, 10, -1, []string{"", "", ""}, "tì", "", strict, allowReef) // Or any additional suffixes
				}
			}
		} else if !strict && strings.HasPrefix(input.Word, "ti") {
			if input.InsistPOS == "any" || input.InsistPOS == "n." {
				// remove it
				newCandidate := candidateDupe(input)
				newCandidate.Word = strings.TrimPrefix(input.Word, "ti")
				newCandidate.InsistPOS = "v."
				newCandidate.Prefixes, added = isDuplicateFix(newCandidate.Prefixes, "tì", strict, allowReef)
				if added {
					deconjugateHelper(newCandidate, 10, 10, -1, []string{"", "", ""}, "tì", "", strict, allowReef) // No other prefixes allowed

					newCandidate.Word = "ì" + newCandidate.Word
					deconjugateHelper(newCandidate, 10, 10, -1, []string{"", "", ""}, "tì", "", strict, allowReef) // Or any additional suffixes
				}
			}
		}
	}

	switch suffixCheck {
	case 0:
		// Made sì its own suffix and no suffixes can come after it
		if len(input.Suffixes) == 0 && strings.HasSuffix(input.Word, "sì") {
			newCandidate := candidateDupe(input)
			newCandidate.Word = strings.TrimSuffix(newCandidate.Word, "sì")
			newCandidate.Suffixes = append(newCandidate.Suffixes, "sì")
			deconjugateHelper(newCandidate, newPrefixCheck, 1, unlenite, infix, "", "sì", strict, allowReef)
		} else if !strict && len(input.Suffixes) == 0 && strings.HasSuffix(input.Word, "si") {
			newCandidate := candidateDupe(input)
			newCandidate.Word = strings.TrimSuffix(newCandidate.Word, "si")
			newCandidate.Suffixes = append(newCandidate.Suffixes, "sì")
			deconjugateHelper(newCandidate, newPrefixCheck, 1, unlenite, infix, "", "sì", strict, allowReef)
		}
		// special case: short genitives of pronouns like "oey" and "ngey"
		if input.InsistPOS == "any" || input.InsistPOS == "n." {
			if strings.HasSuffix(input.Word, "y") {
				// oey to oe
				newCandidate := candidateDupe(input)
				newCandidate.Word = strings.TrimSuffix(input.Word, "y")
				newCandidate.InsistPOS = "pn."
				newCandidate.Suffixes, added = isDuplicateFix(newCandidate.Suffixes, "y", strict, allowReef)
				if added {
					deconjugateHelper(newCandidate, newPrefixCheck, 10, unlenite, []string{}, "", "y", strict, allowReef)

					// ngey to nga
					if strings.HasSuffix(newCandidate.Word, "e") {
						newCandidate.Word = strings.TrimSuffix(newCandidate.Word, "e") + "a"
						newCandidate.InsistPOS = "pn."
						deconjugateHelper(newCandidate, newPrefixCheck, 10, unlenite, []string{}, "", "y", strict, allowReef)
					}
				}
			}
		}
		fallthrough
	case 1:
		if input.InsistPOS == "any" || input.InsistPOS == "n." {
			for _, oldSuffix := range adposuffixes {
				// If it has one of them,
				if strings.HasSuffix(input.Word, oldSuffix) {
					newString = strings.TrimSuffix(input.Word, oldSuffix)

					// Make sure you're using a valid case ending
					if !verifyCaseEnding(newString, oldSuffix) {
						continue
					}

					newCandidate := candidateDupe(input)
					newCandidate.Word = newString
					newCandidate.InsistPOS = "n."
					newCandidate.Suffixes, added = isDuplicateFix(newCandidate.Suffixes, oldSuffix, strict, allowReef)
					if !added {
						continue
					}
					// all set to 2 to avoid mengeyä -> mengo -> me + 'eng + o
					deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", oldSuffix, strict, allowReef)

					if oldSuffix == "ä" && !strings.HasSuffix(input.Word, "yä") && strings.HasSuffix(input.Word, "iä") { // Don't make peyä -> yä -> ya (air)
						// soaiä, tìftiä, etx.
						newString += "a"
						newCandidate.Word = newString
						deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", oldSuffix, strict, allowReef)
					} else if allowReef && oldSuffix == "e" && !strings.HasSuffix(input.Word, "ye") && strings.HasSuffix(input.Word, "ie") {
						// reef of above
						newString += "a"
						newCandidate.Word = newString
						deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", "ä", strict, allowReef)
					} else if (oldSuffix == "yä" || (allowReef && oldSuffix == "ye")) && strings.HasSuffix(newString, "e") {
						// A one-off
						if newString == "tse" {
							newCandidate.Word = "tsaw"
							deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", oldSuffix, strict, allowReef)
						}
						// ngeyä -> nga
						newCandidate.Word = strings.TrimSuffix(newString, "e") + "a"
						deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", oldSuffix, strict, allowReef)
						// oengeyä
						newCandidate.Word = strings.TrimSuffix(newString, "e")
						if newCandidate.Word == "oeng" { //no mengeyä -> meng -> me + 'eng
							deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", oldSuffix, strict, allowReef)
						}
						// sneyä -> sno
						newCandidate.Word = strings.TrimSuffix(newString, "e") + "o"
						deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", oldSuffix, strict, allowReef)
					} else if !strict && oldSuffix == "ye" && strings.HasSuffix(newString, "e") {
						// reef of above
						if newString == "tse" {
							newCandidate.Word = "tsaw"
							deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", "yä", strict, allowReef)
						}
						// ngeye -> nga
						newCandidate.Word = strings.TrimSuffix(newString, "e") + "a"
						deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", "yä", strict, allowReef)
						// oengeye
						newCandidate.Word = strings.TrimSuffix(newString, "e")
						if newCandidate.Word == "oeng" { //no mengeyä -> meng -> me + 'eng
							deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", "yä", strict, allowReef)
						}
						// sneye -> sno
						newCandidate.Word = strings.TrimSuffix(newString, "e") + "o"
						deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", "yä", strict, allowReef)
					} else if vowels, ok := vowelSuffixes["yä"]; ok {
						for _, vowel := range vowels {
							// Make sure zekwä-äo is recognized
							if strings.HasSuffix(newString, vowel+"-") {
								newString = strings.TrimSuffix(newString, "-")
								newCandidate.Word = newString
								deconjugateHelper(newCandidate, newPrefixCheck, 2, unlenite, []string{}, "", "yä", strict, allowReef)
							}
						}
					}
				}
			}
		}
		fallthrough
	case 2:
		if input.InsistPOS == "any" || input.InsistPOS == "n." {
			if strings.HasSuffix(input.Word, "pe") {
				newString = strings.TrimSuffix(input.Word, "pe")

				newCandidate := candidateDupe(input)
				newCandidate.Word = newString
				newCandidate.InsistPOS = "n."
				newCandidate.Suffixes, added = isDuplicateFix(newCandidate.Suffixes, "pe", strict, allowReef)
				if added {
					deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, []string{}, "", "pe", strict, allowReef)
				}
			}
		}
		fallthrough
	case 3:
		// If it has one of them,
		if strings.HasSuffix(input.Word, "a") && input.InsistPOS != "n." && !strings.HasPrefix(input.InsistPOS, "ad") {
			// No nouns, adpositions or adverbs
			newString = strings.TrimSuffix(input.Word, "a")

			newCandidate := candidateDupe(input)
			newCandidate.Word = newString
			newCandidate.InsistPOS = "adj."
			newCandidate.Suffixes, added = isDuplicateFix(newCandidate.Suffixes, "a", strict, allowReef)
			if added {
				deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, []string{"", "", ""}, "", "a", strict, allowReef)
				newCandidate.InsistPOS = "v."
				deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, []string{"", "", ""}, "", "a", strict, allowReef)
			}
		}

		fallthrough
	case 4: // -o suffix "some"
		if input.InsistPOS == "any" || input.InsistPOS == "n." {
			if strings.HasSuffix(input.Word, "o") {
				newString = strings.TrimSuffix(input.Word, "o")

				newCandidate := candidateDupe(input)
				newCandidate.Word = newString
				newCandidate.InsistPOS = "n."
				newCandidate.Suffixes, added = isDuplicateFix(newCandidate.Suffixes, "o", strict, allowReef)
				if added {
					deconjugateHelper(newCandidate, newPrefixCheck, 4, unlenite, []string{}, "", "o", strict, allowReef)

					// Make sure fya'o-o is recognized
					if vowels, ok := vowelSuffixes["o"]; ok {
						for _, vowel := range vowels {
							// Make sure fya'o-o is recognized
							if strings.HasSuffix(newString, vowel+"-") {
								newString = strings.TrimSuffix(newString, "-")
								newCandidate.Word = newString
								deconjugateHelper(newCandidate, newPrefixCheck, 5, unlenite, []string{}, "", "o", strict, allowReef)
							}
						}
					}
				}
			}
		}
		fallthrough
	case 5:
		if input.InsistPOS == "any" || input.InsistPOS == "n." {
			for _, oldSuffix := range stemSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input.Word, oldSuffix) {
					newString = strings.TrimSuffix(input.Word, oldSuffix)

					//candidates = append(candidates, newString)
					newCandidate := candidateDupe(input)
					newCandidate.Word = newString
					newCandidate.InsistPOS = "n."
					newCandidate.Suffixes, added = isDuplicateFix(newCandidate.Suffixes, oldSuffix, strict, allowReef)
					if !added {
						continue
					}
					deconjugateHelper(newCandidate, newPrefixCheck, 6, unlenite, []string{}, "", oldSuffix, strict, allowReef)
				}
			}
		}
		fallthrough
	case 6:
		// If it has one of them,
		if input.InsistPOS == "any" || input.InsistPOS == "n." {
			// verb suffixes change things from verbs to nouns, that's why we check for noun status
			for _, oldSuffix := range verbSuffixes {
				// If it has one of them,
				if strings.HasSuffix(input.Word, oldSuffix) {
					newString = strings.TrimSuffix(input.Word, oldSuffix)
					newCandidate := candidateDupe(input)
					newCandidate.Word = newString
					newCandidate.InsistPOS = "v."

					newCandidate.Suffixes, added = isDuplicateFix(newCandidate.Suffixes, oldSuffix, strict, allowReef)
					if !added {
						continue
					}
					deconjugateHelper(newCandidate, 10, 10, unlenite, []string{}, "", oldSuffix, strict, allowReef) // Don't allow any other prefixes
					// They may turn the InsistPOS back into a noun

					if oldSuffix == "yu" && strings.HasSuffix(newString, "si") {
						newCandidate.Word = strings.TrimSuffix(newString, "si") + " si"
						deconjugateHelper(newCandidate, 10, 10, unlenite, []string{}, "", oldSuffix, strict, allowReef) // don't allow any other prefixes or suffixes
					}
				}
			}
		}
	}

	if len(infix) == 3 && len(input.Infixes) < 3 {
		// Maybe someone else came in with stripped infixes
		if len(input.Word) > 2 && input.Word[len(input.Word)-3] != ' ' &&
			strings.HasSuffix(input.Word, "si") && !strings.HasSuffix(input.Word, "usi") &&
			!strings.HasSuffix(input.Word, "atsi") {
			newCandidate := candidateDupe(input)
			newCandidate.Word = strings.TrimSuffix(input.Word, "si") + " si"
			newCandidate.InsistPOS = "v."
			deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, infix, "", "", strict, allowReef)
		} else { // If there is a "si", we don't need to check for infixes
			// Check for infixes
			runes := []rune(input.Word)
			for i, c := range runes {
				// Infixes can only begin with vowels
				if is_vowel(c) {
					shortString := string(runes[i:])
					for _, newInfix := range infixes[c] {
						available, newInfixes := verifyInfix(infix, newInfix)
						if available && strings.HasPrefix(shortString, newInfix) {
							newCandidate := candidateDupe(input)
							newCandidate.Word = string(runes[:i]) + strings.TrimPrefix(shortString, newInfix)
							cont := false
							newCandidate.Infixes, cont = isDuplicateFix(newCandidate.Infixes, newInfix, strict, allowReef)
							if !cont {
								continue
							}
							newCandidate.InsistPOS = "v."
							deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, newInfixes, "", "", strict, allowReef)

							if newInfix == "ol" {
								newCandidate := candidateDupe(input)
								newCandidate.Word = string(runes[:i]) + "ll" + strings.TrimPrefix(shortString, newInfix)
								newCandidate.Infixes, _ = isDuplicateFix(newCandidate.Infixes, newInfix, strict, allowReef)
								newCandidate.InsistPOS = "v."
								deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, newInfixes, "", "", strict, allowReef)
							} else if newInfix == "er" {
								newCandidate := candidateDupe(input)
								newCandidate.Word = string(runes[:i]) + "rr" + strings.TrimPrefix(shortString, newInfix)
								newCandidate.Infixes, _ = isDuplicateFix(newCandidate.Infixes, newInfix, strict, allowReef)
								newCandidate.InsistPOS = "v."
								deconjugateHelper(newCandidate, newPrefixCheck, suffixCheck, unlenite, newInfixes, "", "", strict, allowReef)
							}
						}
					}
				}
			}
		}
	}

	// Short lenition check
	if unlenite != -1 {
		for _, oldPrefix := range unlenitionLetters {
			// If it has a letter that could have changed for lenition,
			if strings.HasPrefix(input.Word, oldPrefix) {
				// put all possibilities in the candidates
				for _, newPrefix := range unlenition[oldPrefix] {
					newCandidate := candidateDupe(input)
					newString = newPrefix + strings.TrimPrefix(input.Word, oldPrefix)
					newCandidate.Word = newString
					if oldPrefix != newPrefix {
						newCandidate.Lenition = []string{newPrefix + "→" + oldPrefix}
					}
					deconjugateHelper(newCandidate, prefixCheck, suffixCheck, -1, []string{}, "", "", strict, allowReef)
				}
				break // We don't want the "ts" to become "txs"
			}
		}
	}

	return candidates
}

// Helper for TestDeconjugations
func allIConfigs(input string, discrimRune rune, replaceRune rune, strict bool, allowReef bool) []string {
	discrim := string(discrimRune)
	replace := string(replaceRune)
	cCount := strings.Count(input, discrim)
	results := []string{input}
	buffer := strings.Builder{}

	// Prevent a DoS attack with alalalalalalalalal
	// The most as and äs are ay-sna-{word}-kxamlä
	// {word} can be tsawlapxangrr, sälätxayn, tawsyuratan,
	// atanzaw, sä'anla, sängä'än, lanay'ka, lamaytxa,
	// txantsawtsray, särangal or säkahena
	if cCount == 0 || cCount > 7 {
		return results
	}

	aArray := []string{}
	for range cCount {
		aArray = append(aArray, discrim)
	}

	splitString := strings.Split(input, discrim)

	// Count with a and ä like a binary number
	for i := range int(math.Round(math.Pow(2, float64(cCount)))) {
		tempI := i
		for a := range len(aArray) {
			pow := int(math.Round(math.Pow(2, float64(a+1))))
			if tempI%pow == 0 {
				aArray[a] = discrim
			} else {
				aArray[a] = replace
				tempI -= int(math.Round(math.Pow(2, float64(a))))
			}
		}
		buffer.WriteString(splitString[0])
		for i, split := range aArray {
			buffer.WriteString(split)
			buffer.WriteString(splitString[i+1])
		}

		newAConfig := dialectCrunch([]string{buffer.String()}, false, strict, allowReef)[0]

		buffer.Reset()

		if newAConfig == input {
			continue
		}

		results = append(results, newAConfig)
	}

	return results
}

func Deconjugate(input string, strict bool, allowReef bool) []ConjugationCandidate {
	candidates = []ConjugationCandidate{} //empty array of strings
	candidateMap = map[string]ConjugationCandidate{}
	newCandidate := ConjugationCandidate{}
	newCandidate.Word = input
	newCandidate.InsistPOS = "any"
	deconjugateHelper(newCandidate, 0, 0, 0, []string{"", "", ""}, "", "", strict, allowReef)

	candidates = candidates[1:]
	return candidates
}

func TestDeconjugations(dict *map[string][]Word, searchNaviWord string, strict bool, allowReef bool, umlaut bool) (results []Word) {
	conjugations := Deconjugate(searchNaviWord, strict, allowReef)

	searchNaviWord = strings.ReplaceAll(searchNaviWord, "ù", "u")

	allAConfigs := []string{searchNaviWord}
	allIAConfigs := []string{}

	//For using a to search ä
	if !strict {
		allAConfigs = append(allAConfigs, allIConfigs(searchNaviWord, 'a', 'ä', strict, allowReef)...)

		for _, config := range allAConfigs {
			allIAConfigs = append(allIAConfigs, config)
			allIAConfigs = append(allIAConfigs, allIConfigs(config, 'i', 'ì', strict, allowReef)...)
		}

		for _, a := range allIAConfigs {
			newCandidate := ConjugationCandidate{Word: a, InsistPOS: "any"}
			conjugations = append(conjugations, newCandidate)
			conjugations = append(conjugations, Deconjugate(a, strict, allowReef)...)
		}

		// For using i to search ì
	}

	// Reintroduce the umlaut if needed
	if allowReef && umlaut {
		for _, a := range conjugations {
			nucleusCount := 0
			for _, b := range []string{"a", "ä", "e", "i", "ì", "o", "u", "ù", "ll", "rr"} {
				nucleusCount += strings.Count(a.Word, b)
			}
			if nucleusCount == 1 && strings.Contains(a.Word, "e") {
				a.Word = strings.ReplaceAll(a.Word, "e", "ä")
				conjugations = append(conjugations, a)
			}
		}
	}

	for _, candidate := range conjugations {
		// Avoid fìfìtseng, zeykeyko and peupe
		skip := false
		if affixes, ok := productiveCompounds[candidate.Word]; ok {
			if len(affixes[0]) > 0 {
				for _, a := range candidate.Prefixes {
					for _, b := range affixes[0] {
						if a == b {
							skip = true
							break
						}
					}
					if skip {
						break
					}
				}
			}
			if !skip {
				for _, a := range candidate.Infixes {
					for _, b := range affixes[1] {
						if a == b {
							skip = true
							break
						}
					}
					if skip {
						break
					}
				}
			}
			if !skip {
				for _, a := range candidate.Suffixes {
					for _, b := range affixes[2] {
						if a == b {
							skip = true
							break
						}
					}
					if skip {
						break
					}
				}
			}
			if skip {
				continue
			}
		}
		if candidate.Word == "kive" && umlaut {
			continue
		}
		a := strings.ReplaceAll(candidate.Word, "ù", "u")

		standardizedWordArray := strings.Split(a, " ")
		if !strict {
			standardizedWordArray = dialectCrunch(standardizedWordArray, false, strict, allowReef)
		}

		a = ""
		for i, b := range standardizedWordArray {
			if i != 0 {
				a += " "
			}
			a += b
		}

		if allowReef {
			a = dialectCrunch([]string{a}, false, true, true)[0]
		}

		for _, c := range (*dict)[a] {
			for _, pos := range strings.Split(c.PartOfSpeech, ",") {
				pos = strings.ReplaceAll(pos, " ", "")

				// An inter. can act like a noun or an adjective, so it gets special treatment
				if pos == "inter." && candidate.InsistPOS[0] != 'v' && len(candidate.Infixes) == 0 {
					dupe := false
					for _, b := range results {
						if b.Navi == c.Navi {
							dupe = true
							break
						}
					}
					if !dupe {
						a := c
						a.Affixes.Lenition = candidate.Lenition
						a.Affixes.Prefix = candidate.Prefixes
						a.Affixes.Infix = candidate.Infixes
						a.Affixes.Suffix = candidate.Suffixes
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
				if len(candidate.Infixes) == 1 && candidate.Infixes[0] == "us" {
					// Reverse search is more likely to find it immediately
					for i := len(candidate.Prefixes) - 1; i >= 0; i-- {
						if candidate.Prefixes[i] == "tì" {
							gerund = true
							break
						}
					}
					if !gerund {
						participle = true
					}
				} else if len(candidate.Infixes) > 0 {
					// Now reverse search is just gratuitous
					for i := len(candidate.Infixes) - 1; i >= 0; i-- {
						if candidate.Infixes[i] == "us" || candidate.Infixes[i] == "awn" {
							participle = true
							break
						}
					}
				}

				// If the InsistPOS and found word agree they are nouns
				if len(candidate.Suffixes) < 3 && len(candidate.Suffixes) > 0 && candidate.Suffixes[0] == "tswo" {
					if pos[0] == 'v' {
						siVerb := false
						if len(candidate.Infixes) == 0 {
							if strict {
								if _, ok := multiword_words[candidate.Word]; ok {
									for _, b := range multiword_words[candidate.Word] {
										if b[0] == "si" {
											siVerb = true
											a := c
											a.Affixes.Lenition = candidate.Lenition
											a.Affixes.Prefix = candidate.Prefixes
											a.Affixes.Infix = candidate.Infixes
											a.Affixes.Suffix = candidate.Suffixes
											results = AppendAndAlphabetize(results, a)
											break
										}
									}
								}
								if !siVerb {
									a := c
									a.Affixes.Lenition = candidate.Lenition
									a.Affixes.Prefix = candidate.Prefixes
									a.Affixes.Infix = candidate.Infixes
									a.Affixes.Suffix = candidate.Suffixes
									results = AppendAndAlphabetize(results, a)
								}
							} else {
								if _, ok := multiword_words_loose[candidate.Word]; ok {
									for _, b := range multiword_words_loose[candidate.Word] {
										if b[0] == "si" {
											siVerb = true
											a := c
											a.Affixes.Lenition = candidate.Lenition
											a.Affixes.Prefix = candidate.Prefixes
											a.Affixes.Infix = candidate.Infixes
											a.Affixes.Suffix = candidate.Suffixes
											results = AppendAndAlphabetize(results, a)
											break
										}
									}
								}
								if !siVerb {
									a := c
									a.Affixes.Lenition = candidate.Lenition
									a.Affixes.Prefix = candidate.Prefixes
									a.Affixes.Infix = candidate.Infixes
									a.Affixes.Suffix = candidate.Suffixes
									results = AppendAndAlphabetize(results, a)
								}
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
						if strings.Contains(searchNaviWord, noTìftang) || strings.Contains(searchNaviWord, dialectCrunch([]string{rebuiltVerb}, false, strict, allowReef)[0]) {
							a := c
							a.Affixes.Lenition = candidate.Lenition
							a.Affixes.Prefix = candidate.Prefixes
							a.Affixes.Infix = candidate.Infixes
							a.Affixes.Suffix = candidate.Suffixes
							results = AppendAndAlphabetize(results, a)
						} else if len(results) == 0 {
							results = AppendAndAlphabetize(results, infixError(searchNaviWord, "tì"+rebuiltVerb, c.IPA))
						}
					}
				} else if candidate.InsistPOS == "n." {
					// n., pn., Prop.n. and inter. (but not vin.)
					if len(candidate.Infixes) == 0 {
						if (pos[0] != 'v' && strings.HasSuffix(pos, "n.")) || pos == "inter." {
							a := c
							a.Affixes.Lenition = candidate.Lenition
							a.Affixes.Prefix = candidate.Prefixes
							a.Affixes.Suffix = candidate.Suffixes
							results = AppendAndAlphabetize(results, a)
						}
					}
				} else if candidate.InsistPOS == "pn." {
					// pn.
					if len(candidate.Infixes) == 0 && strings.HasSuffix(pos, "pn.") {
						a := c
						a.Affixes.Lenition = candidate.Lenition
						a.Affixes.Prefix = candidate.Prefixes
						a.Affixes.Suffix = candidate.Suffixes
						results = AppendAndAlphabetize(results, a)
					}
				} else if candidate.InsistPOS == "adj." {
					posNoun := pos
					if len(candidate.Infixes) == 0 && (posNoun == "adj." || posNoun == "num.") {
						a := c
						a.Affixes.Lenition = candidate.Lenition
						a.Affixes.Prefix = candidate.Prefixes
						a.Affixes.Suffix = candidate.Suffixes
						results = AppendAndAlphabetize(results, a)
					}
				} else if candidate.InsistPOS == "v." {
					posNoun := pos
					if strings.HasPrefix(posNoun, "v") {
						// Verbs with -tswo or -yu cannot have infixes
						if len(candidate.Suffixes) > 0 {
							for i := len(candidate.Suffixes) - 1; i >= 0; i-- {
								if candidate.Suffixes[i] == "a" {
									attributed = true
									break
								}
							}
							// Forward search fixs the "a" before "yu" and "tswo"
							for i := len(candidate.Suffixes) - 1; i >= 0; i-- {
								for _, j := range verbSuffixes {
									if candidate.Suffixes[i] == j {
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

						if len(candidate.Prefixes) > 0 {
							// Reverse search is more likely to find it immediately
							for i := len(candidate.Prefixes) - 1; i >= 0; i-- {
								if candidate.Prefixes[i] == "a" {
									attributed = true
								} else if candidate.Prefixes[i] == "tì" {
									// we found gerunds up top, so this isn't needed
									looseTì = true
									break
								} else {
									for _, j := range verbPrefixes {
										if candidate.Prefixes[i] == j {
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
						if attributed && (len(candidate.Infixes) == 0 || infixBan) && !tsuk {
							continue
						}

						// Take action on tsuk-verb-yus and a-verb-tswos
						if doubleBan || (attributed && !tsuk && infixBan) || looseTì {
							continue
						}

						a := c
						a.Affixes.Lenition = candidate.Lenition
						a.Affixes.Prefix = candidate.Prefixes
						a.Affixes.Suffix = candidate.Suffixes
						a.Affixes.Infix = candidate.Infixes

						if infixBan {
							if len(candidate.Infixes) > 0 {
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
						if c.InfixLocations == "z<0><1>en<2>ke" && implContainsAny(candidate.Infixes, []string{"ats", "uy"}) {
							rebuiltVerb = "z<0><1>en<2>eke"
						}
						firstInfixes := ""

						for _, newInfix := range candidate.Infixes {
							if _, ok := prefirstMap[newInfix]; ok {
								firstInfixes += newInfix
								rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<0>", firstInfixes)
								if (allowReef && newInfix == "epeyk") || newInfix == "äpeyk" {
									newCandidateInfixes := []string{}
									for _, newInfix2 := range candidate.Infixes {
										// äpeyk gets split
										if (allowReef && newInfix2 == "epeyk") || newInfix2 == "äpeyk" {
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
						for _, newInfix := range candidate.Infixes {
							if _, ok := firstMap[newInfix]; ok {
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
						for _, newInfix := range candidate.Infixes {
							if !strict && (newInfix == "eng" || newInfix == "ang") {
								rebuiltVerb = strings.ReplaceAll(rebuiltVerb, "<2>", "äng")
								break
							} else if _, ok := secondMap[newInfix]; ok {
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
						rebuiltVerbArray := strings.Split(rebuiltVerb, " ")
						if !strict || allowReef {
							rebuiltVerbArray = dialectCrunch(rebuiltVerbArray, false, strict, allowReef)
						}

						rebuiltVerb = ""
						for k, x := range rebuiltVerbArray {
							if k != 0 {
								rebuiltVerb += " "
							}
							rebuiltVerb += x
						}

						/*searchNaviWordSquish := strings.ReplaceAll(searchNaviWord, "ù", "u")
						searchNaviWordSquish = strings.ReplaceAll(searchNaviWordSquish, "-", " ")

						if !strict {
							searchNaviWordSquish = dialectCrunch([]string{searchNaviWordSquish}, false)[0]
						}*/

						if len(candidate.Infixes) == 0 || implContainsAny([]string{rebuiltVerb}, allAConfigs) {
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
				} else if candidate.InsistPOS == "nì." {
					posNoun := pos
					if len(candidate.Infixes) == 0 && (posNoun == "adj." || posNoun == "pn.") {
						a := c
						a.Affixes.Lenition = candidate.Lenition
						a.Affixes.Prefix = candidate.Prefixes
						a.Affixes.Suffix = candidate.Suffixes
						results = AppendAndAlphabetize(results, a)
					}
				} else if len(candidate.Infixes) == 0 {
					a := c
					a.Affixes.Lenition = candidate.Lenition
					a.Affixes.Prefix = candidate.Prefixes
					a.Affixes.Suffix = candidate.Suffixes
					results = AppendAndAlphabetize(results, a)
				}
			}
		}
	}
	return
}
