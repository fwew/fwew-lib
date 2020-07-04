//	This file is part of Fwew.
//	Fwew is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU General Public License as published by
// 	the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
//	Fwew is distributed in the hope that it will be useful,
//	but WITHOUT ANY WARRANTY; without even implied warranty of
//	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	GNU General Public License for more details.
//
//	You should have received a copy of the GNU General Public License
//	along with Fwew.  If not, see http://gnu.org/licenses/

// Package main contains all the things. affixes.go handles affix parsing of input.
package fwew_lib

import (
	"log"
	"regexp"
	"strings"
)

const (
	adj   string = "adj."
	adv   string = "adv."
	dem   string = "dem."
	inter string = "inter."
	n     string = "n."
	num   string = "num."
	pn    string = "pn."
	propN string = "prop.n."
	vin   string = "vin."
	svin  string = "svin."
)

// try to add prefixes to the word. Return the attempt with placed prefixes
// Has to be provided with a previousAttempt, the word to go from and add prefixes to.
func (w *Word) prefix(target string, previousAttempt string) string {
	var (
		re            *regexp.Regexp
		reString      string
		attempt       string
		matchPrefixes []string
	)

	// pull this out of the switch because the pos data for verbs is so irregular,
	// the switch condition would be like 25 possibilities long
	if strings.HasPrefix(w.PartOfSpeech, "v") ||
		strings.HasPrefix(w.PartOfSpeech, svin) || w.PartOfSpeech == "" {
		inf := w.Affixes.Infix
		if len(inf) > 0 && (inf[0] == "us" || inf[0] == "awn") {
			reString = "(a|tì)?"
		} else if strings.Contains(target, "ketsuk") || strings.Contains(target, "tsuk") {
			reString = "(a)?(ketsuk|tsuk)?"
		} else if strings.Contains(target, "siyu") && w.PartOfSpeech == "vin." {
			reString = "^(pep|pem|pe|fray|tsay|fay|pay|fra|fì|tsa)?(ay|me|pxe|pe)?(fne)?(munsna)?"
		}
	} else {
		switch w.PartOfSpeech {
		case n, pn, propN:
			reString = "^(pep|pem|pe|fray|tsay|fay|pay|fra|fì|tsa)?(ay|me|pxe|pe)?(fne)?(munsna)?"
		case adj:
			reString = "^(nìk|nì|a)?(ke|a)?"
		default:
			return previousAttempt // Not a type that has a prefix, return word without attempting.
		}
	}

	if strings.HasPrefix(target, "me") || strings.HasPrefix(target, "pxe") || strings.HasPrefix(target, "pe") {
		if strings.HasPrefix(previousAttempt, "e") {
			reString = reString + "(e)?"
			previousAttempt = previousAttempt[1:]
		} else if strings.HasPrefix(previousAttempt, "'e") {
			reString = reString + "('e)?"
			previousAttempt = previousAttempt[2:]
		}
	}

	// soaiä replacement
	if w.Navi == "soaia" && strings.HasSuffix(target, "soaiä") {
		previousAttempt = strings.Replace(previousAttempt, "soaia", "soai", -1)
	}

	reString = reString + previousAttempt + ".*"

	if debugMode {
		log.Printf("Prefix reString: %s\n", reString)
	}

	re = regexp.MustCompile(reString)
	tmp := re.FindAllStringSubmatch(target, -1)
	if len(tmp) > 0 && len(tmp[0]) >= 1 {
		matchPrefixes = tmp[0][1:]
	}
	matchPrefixes = DeleteEmpty(matchPrefixes)

	if debugMode {
		log.Printf("matchPrefixes: %s\n", matchPrefixes)
	}

	// no productive prefixes found; why bother to continue?
	if len(matchPrefixes) == 0 {
		return previousAttempt
	}
	// only allow lenition after lenition-causing prefixes when prefixes and lenition present
	if len(w.Affixes.Lenition) > 0 && len(matchPrefixes) > 0 {
		if Contains(matchPrefixes, []string{"fne", "munsna"}) {
			return previousAttempt
		}
		lenPre := []string{"pep", "pem", "pe", "fray", "tsay", "fay", "pay", "ay", "me", "pxe"}
		if Contains(matchPrefixes, []string{"fì", "tsa", "fra"}) && !Contains(matchPrefixes, lenPre) {
			return previousAttempt
		}
	}

	// build what prefixes to put on
	for _, p := range matchPrefixes {
		attempt = attempt + p
	}

	previousAttempt = attempt + previousAttempt

	matchPrefixes = DeleteElement(matchPrefixes, "e")
	if len(matchPrefixes) > 0 {
		w.Affixes.Prefix = append(w.Affixes.Prefix, matchPrefixes...)
	}

	return previousAttempt
}

// try to add suffixes to the word. Return the attempt with placed suffixes
// Has to be provided with a previousAttempt, the word to go from and add suffixes to.
func (w *Word) suffix(target string, previousAttempt string) string {
	var (
		re            *regexp.Regexp
		tmp           [][]string
		reString      string
		attempt       string
		matchSuffixes []string
	)
	const (
		adjSufRe string = "(a|sì)?$"
		nSufRe   string = "(nga'|tsyìp|tu)?(o)?(pe)?(mungwrr|kxamlä|tafkip|pxisre|pximaw|ftumfa|mìkam|nemfa|takip|lisre|talun|krrka|teri|fkip|pxaw|pxel|luke|rofa|fpi|ftu|kip|vay|lok|maw|sìn|sre|few|kam|kay|nuä|sko|yoa|äo|eo|fa|hu|ka|mì|na|ne|ta|io|uo|ro|wä|sì|ìri|ìl|eyä|yä|ä|it|ri|ru|ti|ur|l|r|t)?$"
		ngey     string = "ngey"
	)

	// hardcoded hack for tseyä
	if target == "tseyä" && w.Navi == "tsaw" {
		w.Affixes.Suffix = []string{"yä"}
		return "tseyä"
	}

	// hardcoded hack for oey
	if target == "oey" && w.Navi == "oe" {
		w.Affixes.Suffix = []string{"y"}
		return "oey"
	}

	// hardcoded hack for ngey
	if target == ngey && w.Navi == "nga" {
		w.Affixes.Suffix = []string{"y"}
		return ngey
	}

	// verbs
	if !strings.Contains(w.PartOfSpeech, adv) &&
		strings.Contains(w.PartOfSpeech, "v") || w.PartOfSpeech == "" {
		inf := w.Affixes.Infix
		pre := w.Affixes.Prefix
		// word is verb with <us> or <awn>
		if len(inf) == 1 && (inf[0] == "us" || inf[0] == "awn") {
			// it's a tì-<us> gerund; treat it like a noun
			if len(pre) > 0 && ContainsStr(pre, "tì") && inf[0] == "us" {
				reString = nSufRe
				// Just a regular <us> or <awn> verb
			} else {
				reString = adjSufRe
			}
			// It's a tsuk/ketsuk adj from a verb
		} else if len(inf) == 0 && Contains(pre, []string{"tsuk", "ketsuk"}) {
			reString = adjSufRe
		} else if strings.Contains(target, "tswo") {
			reString = "(tswo)?" + nSufRe
		} else {
			reString = "(yu)?$"
		}
	} else {
		switch w.PartOfSpeech {
		// nouns and noun-likes
		case n, pn, propN, inter, dem, "dem., pn.":
			reString = nSufRe
			// adjectives
		case adj:
			reString = adjSufRe
		// numbers
		case num:
			reString = "(ve)?(a)?"
		default:
			return previousAttempt // Not a type that has a suffix, return word without attempting.
		}
	}

	// soaiä support
	if w.Navi == "soaia" && strings.HasSuffix(target, "soaiä") {
		previousAttempt = strings.Replace(previousAttempt, "soaia", "soai", -1)
		reString = previousAttempt + reString
		// o -> e vowel shift support
	} else if strings.HasSuffix(previousAttempt, "o") {
		reString = strings.Replace(previousAttempt, "o", "[oe]", -1) + reString
		// a -> e vowel shift support
	} else if strings.HasSuffix(previousAttempt, "a") {
		reString = strings.Replace(previousAttempt, "a", "[ae]", -1) + reString
	} else if w.Navi == "tsaw" {
		tsaSuf := []string{
			"mungwrr", "kxamlä", "tafkip", "pxisre", "pximaw", "ftumfa", "mìkam", "nemfa", "takip", "lisre", "talun",
			"krrka", "teri", "fkip", "pxaw", "pxel", "luke", "rofa", "fpi", "ftu", "kip", "vay", "lok", "maw", "sìn", "sre",
			"few", "kam", "kay", "nuä", "sko", "yoa", "äo", "eo", "fa", "hu", "ka", "mì", "na", "ne", "ta", "io", "uo",
			"ro", "wä", "ìri", "ri", "ru", "ti", "r"}
		for _, s := range tsaSuf {
			if strings.HasSuffix(target, "tsa"+s) || strings.HasSuffix(target, "sa"+s) {
				previousAttempt = strings.Replace(previousAttempt, "aw", "a", 1)
				reString = previousAttempt + reString
			}
		}
	} else {
		reString = previousAttempt + reString
	}

	if debugMode {
		log.Printf("Suffix reString: %s\n", reString)
	}

	re = regexp.MustCompile(reString)
	if strings.HasSuffix(target, "siyu") {
		tmp = re.FindAllStringSubmatch(strings.Replace(target, "siyu", " siyu", -1), -1)
	} else {
		tmp = re.FindAllStringSubmatch(target, -1)
	}
	if len(tmp) > 0 && len(tmp[0]) >= 1 {
		matchSuffixes = tmp[0][1:]
	}
	matchSuffixes = DeleteEmpty(matchSuffixes)

	if debugMode {
		log.Printf("matchSuffixes: %s\n", matchSuffixes)
	}

	// no productive prefixes found; why bother to continue?
	if len(matchSuffixes) == 0 {
		return previousAttempt
	}

	// build what prefixes to put on
	for _, p := range matchSuffixes {
		attempt = attempt + p
	}

	// o -> e vowel shift support for pronouns with -yä
	if w.PartOfSpeech == pn && ContainsStr(matchSuffixes, "yä") {
		if strings.HasSuffix(previousAttempt, "o") {
			previousAttempt = strings.TrimSuffix(previousAttempt, "o") + "e"
			// a -> e vowel shift support
		} else if strings.HasSuffix(previousAttempt, "a") {
			previousAttempt = strings.TrimSuffix(previousAttempt, "a") + "e"
		}
	}
	previousAttempt = previousAttempt + attempt
	if strings.Contains(previousAttempt, " ") && strings.HasSuffix(previousAttempt, "siyu") {
		previousAttempt = strings.Replace(previousAttempt, " siyu", "siyu", -1)
	}

	w.Affixes.Suffix = append(w.Affixes.Suffix, matchSuffixes...)

	return previousAttempt
}

// try to add infixes to the word. Returns the attempt with placed infixes
func (w *Word) infix(target string) string {
	// Have we already attempted infixes?
	if w.Affixes.Infix != nil {
		return ""
	}

	// Does the word even have infix positions??
	if w.InfixLocations == "\\N" {
		return ""
	}

	var (
		re              *regexp.Regexp
		reString        string
		attempt         string
		pos0InfixRe     = "(äp)?(eyk)?"
		pos1InfixRe     = "(ìyev|iyev|ìlm|ìly|ìrm|ìry|ìsy|alm|aly|arm|ary|asy|ìm|imv|ilv|irv|ìy|am|ay|er|iv|ol|us|awn)?"
		pos2InfixRe     = "(eiy|ei|äng|eng|ats|uy)?"
		pos0InfixString string
		pos1InfixString string
		pos2InfixString string
		matchInfixes    []string
	)

	// Hardcode hack for z**enke
	if w.Navi == "zenke" && (strings.Contains(target, "uy") || strings.Contains(target, "ats")) {
		w.InfixLocations = strings.Replace(w.InfixLocations, "ke", "eke", 1)
	}

	reString = strings.Replace(w.InfixLocations, "<0>", pos0InfixRe, 1)

	// handle <ol>ll and <er>rr
	if strings.Contains(reString, "<1>ll") {
		reString = strings.Replace(reString, "<1>ll", pos1InfixRe+"(ll)?", 1)
	} else if strings.Contains(w.InfixLocations, "<1>rr") {
		reString = strings.Replace(reString, "<1>rr", pos1InfixRe+"(rr)?", 1)
	} else {
		reString = strings.Replace(reString, "<1>", pos1InfixRe, 1)
	}
	reString = strings.Replace(reString, "<2>", pos2InfixRe, 1)

	if debugMode {
		log.Printf("Infix reString: %s\n", reString)
	}

	re, _ = regexp.Compile(reString)
	tmp := re.FindAllStringSubmatch(target, -1)
	if len(tmp) > 0 && len(tmp[0]) >= 1 {
		matchInfixes = tmp[0][1:]
	}
	matchInfixes = DeleteEmpty(matchInfixes)
	matchInfixes = DeleteElement(matchInfixes, "ll")
	matchInfixes = DeleteElement(matchInfixes, "rr")

	for _, i := range matchInfixes {
		if i == "äp" || i == "eyk" {
			pos0InfixString = pos0InfixString + i
		} else if ContainsStr([]string{"eiy", "ei", "äng", "eng", "ats", "uy"}, i) {
			pos2InfixString = i
		} else {
			pos1InfixString = i
		}
	}

	attempt = strings.Replace(w.InfixLocations, "<0>", pos0InfixString, 1)
	attempt = strings.Replace(attempt, "<1>", pos1InfixString, 1)
	attempt = strings.Replace(attempt, "<2>", pos2InfixString, 1)

	// eiy override?
	if ContainsStr(matchInfixes, "eiy") {
		eiy := Index(matchInfixes, "eiy")
		matchInfixes[eiy] = "ei"
	}

	if debugMode {
		log.Printf("matchInfixes: %s\n", matchInfixes)
	}

	// handle <ol>ll and <er>rr
	if strings.Contains(attempt, "olll") {
		attempt = strings.Replace(attempt, "olll", "ol", 1)
	} else if strings.Contains(attempt, "errr") {
		attempt = strings.Replace(attempt, "errr", "er", 1)
	}

	if len(matchInfixes) != 0 {
		w.Affixes.Infix = append(w.Affixes.Infix, matchInfixes...)
	}

	return attempt
}

// table of all the possible lenitions
var lenitionTable = [8][2]string{
	{"kx", "k"},
	{"px", "p"},
	{"tx", "t"},
	{"k", "h"},
	{"p", "f"},
	{"ts", "s"},
	{"t", "s"},
	{"'", ""},
}

func GetLenitionTable() [][2]string {
	return lenitionTable[:]
}

// Lenite the word, based on the attempt. The target is not relevant here, so not given.
// Returns the lenite attempt.
func (w *Word) lenite(attempt string) string {
	// Have we already attempted lenition?
	if w.Affixes.Lenition != nil {
		return attempt
	}

	for _, v := range lenitionTable {
		if strings.HasPrefix(strings.ToLower(w.Navi), v[0]) {
			attempt = strings.Replace(attempt, v[0], v[1], 1)
			w.Affixes.Lenition = append(w.Affixes.Lenition, v[0]+"→"+v[1])
			return attempt
		}
	}
	return attempt
}

// Reconstruct is the main function of affixes.go, responsible for the affixing algorithm
// This will try to reconstruct a Word, so it matches with the target.
// Returns true if word got reconstructed into target!
func (w *Word) reconstruct(target string) bool {
	attempt := w.Navi

	// only try to infix verbs
	if strings.HasPrefix(w.PartOfSpeech, "v") || strings.HasPrefix(w.PartOfSpeech, svin) {
		attempt = w.infix(target)

		if debugMode {
			log.Printf("Navi: %s | Attempt: %s | Target: %s\n", w.Navi, attempt, target)
		}

		if attempt == target {
			return true
		}
	}

	attempt = w.prefix(target, attempt)

	if debugMode {
		log.Println("PREFIX w")
		log.Printf("Navi: %s | Attempt: %s | Target: %s\n", w.Navi, attempt, target)
	}

	if attempt == target {
		return true
	}

	attempt = w.suffix(target, attempt)

	if debugMode {
		log.Println("SUFFIX w")
		log.Printf("Navi: %s | Attempt: %s | Target: %s\n", w.Navi, attempt, target)
	}

	if attempt == target {
		return true
	}

	attempt = w.lenite(attempt)

	if debugMode {
		log.Println("LENITE w")
		log.Printf("Navi: %s | Attempt: %s | Target: %s\n", w.Navi, attempt, target)
	}

	if attempt == target {
		return true
	}

	// try it another time, with different guess order!

	// clean up word
	w.Affixes = affix{}

	attempt = w.lenite(w.Navi)

	if debugMode {
		log.Println("LENITE wl")
		log.Printf("Navi: %s | Attempt: %s | Target: %s\n", w.Navi, attempt, target)
	}

	if attempt == target {
		return true
	}

	attempt = w.prefix(target, attempt)

	if debugMode {
		log.Println("PREFIX wl")
		log.Printf("Navi: %s | Attempt: %s | Target: %s\n", w.Navi, attempt, target)
	}

	if attempt == target {
		return true
	}

	attempt = w.suffix(target, attempt)

	if debugMode {
		log.Println("SUFFIX wl")
		log.Printf("Navi: %s | Attempt: %s | Target: %s\n", w.Navi, attempt, target)
	}

	if attempt == target {
		return true
	}

	if debugMode {
		log.Println("GIVING UP")
	}

	return false
}
