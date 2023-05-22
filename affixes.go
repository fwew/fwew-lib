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
	"fmt"
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

	//denying prefixes on tsari, tsar, tsat, tsal
	declensedPronouns := []string{"tsal", "tsat", "tsar", "tsari"}
	checkWord := []string{w.Navi}
	targetTmp := []string{target}
	if Contains(checkWord, declensedPronouns) || Contains(targetTmp, declensedPronouns) {
		w.Affixes.Prefix = []string{""}
		return ""
	}

	// pull this out of the switch because the pos data for verbs is so irregular,
	// the switch condition would be like 25 possibilities long
	if strings.HasPrefix(w.PartOfSpeech, "v") ||
		strings.HasPrefix(w.PartOfSpeech, svin) || w.PartOfSpeech == "" {
		inf := w.Affixes.Infix
		//trying to detect tì-us / sì-us
		rootDiscriminate := 0
		targetDiscriminate := 0
		compareWord := w.Navi
		compareTarget := target
		flagTius := 0

		//checking if root word (w.Navi) has tì or sì by itself
		if (len(inf) > 0 && (inf[0] == "us")) && (strings.Contains(target, "tì") || strings.Contains(target, "sì")) {
			for {
				if strings.Contains(compareWord, "tì") {
					rootDiscriminate = rootDiscriminate + 1
					compareWord = strings.Replace(compareWord, "tì", "", 1)
				} else if strings.Contains(compareWord, "sì") {
					rootDiscriminate = rootDiscriminate + 1
					compareWord = strings.Replace(compareWord, "sì", "", 1)
				} else {
					break
				}
			}

			//checking if target word (previousAttempt) has tì or sì by itself excluding -sì suffix and -usì-
			for {
				if strings.Contains(compareTarget, "usì") {
					compareTarget = strings.Replace(compareTarget, "usì", "", 1)
				} else if strings.HasSuffix(compareTarget, "sì") {
					compareTarget = strings.TrimSuffix(compareTarget, "sì")
				} else if strings.Contains(compareTarget, "tì") {
					targetDiscriminate = targetDiscriminate + 1
					compareTarget = strings.Replace(compareTarget, "tì", "", 1)
				} else if strings.Contains(compareTarget, "sì") {
					targetDiscriminate = targetDiscriminate + 1
					compareTarget = strings.Replace(compareTarget, "sì", "", 1)
				} else {
					break
				}
			}

			//checking if verb has tì or sì as valid infix syllable
			if strings.Contains(w.InfixLocations, "t<0><1><2>ì") ||
				strings.Contains(w.InfixLocations, "s<0><1><2>ì") ||
				strings.Contains(w.InfixLocations, "t<0><1>ì") ||
				strings.Contains(w.InfixLocations, "s<0><1>ì") {
				rootDiscriminate = rootDiscriminate - 1
			}

			if targetDiscriminate > rootDiscriminate {
				flagTius = 1
			} else {
				flagTius = 2
			}
		}

		//detecting participles
		if (len(inf) > 0 && (inf[0] == "us" || inf[0] == "awn")) &&
			(flagTius == 2 || flagTius == 0) {
			reString = "^(a)?"
			//detecting tì-us or sì-us noun
		} else if ((len(inf) > 0 && inf[0] == "us") && (strings.Contains(target, "tì") || strings.Contains(target, "sì"))) || flagTius == 1 {
			reString = "^(pep|pem|pe|fray|tsay|fay|pay|fra|fì|tsa)?(ay|me|pxe|pe)?(fne)?(munsna)?(tì|sì)?"
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
		case num:
			reString = "^(a)?"
		default:
			return previousAttempt // Not a type that has a prefix, return word without attempting.
		}
	}

	// me/pxe/pe + e/'e vowel merge
	if HasPrefixStrArr(target, []string{"me", "pxe", "pe"}) {
		if strings.HasPrefix(previousAttempt, "e") {
			reString = reString + "(e)?"
			previousAttempt = previousAttempt[1:]
		} else if strings.HasPrefix(previousAttempt, "'e") {
			reString = reString + "('e)?"
			previousAttempt = previousAttempt[2:]
		}
	} else if strings.HasPrefix(target, "fne") {
		// prefix boundary vowel merges (ee)
		if strings.HasPrefix(previousAttempt, "e") {
			reString = reString + "(e)?"
			previousAttempt = previousAttempt[1:]
		}
	} else if HasPrefixStrArr(target, []string{"fra", "tsa", "munsna"}) {
		// prefix boundary vowel merges (aa)
		if strings.HasPrefix(previousAttempt, "a") {
			reString = reString + "(a)?"
			previousAttempt = previousAttempt[1:]
		}
	} else if strings.HasPrefix(target, "fì") {
		// prefix boundary vowel (ìì)
		if strings.HasPrefix(previousAttempt, "ì") {
			reString = reString + "(ì)?"
			previousAttempt = previousAttempt[2:]
		}
	}

	// soaiä replacement
	if w.Navi == "soaia" && strings.HasSuffix(target, "soaiä") {
		previousAttempt = strings.Replace(previousAttempt, "soaia", "soai", -1)
	}

	// meuiä replacement
	if w.Navi == "meuia" && strings.HasSuffix(target, "meuiä") {
		previousAttempt = strings.Replace(previousAttempt, "meuia", "meui", -1)
	}

	// kemuiä replacement
	if w.Navi == "kemuia" && (strings.HasSuffix(target, "kemuiä") || strings.HasSuffix(target, "hemuiä")) {
		previousAttempt = strings.Replace(previousAttempt, "emuia", "emui", -1)
	}

	// aungiä replacement
	if w.Navi == "aungia" && strings.HasSuffix(target, "ungiä") {
		previousAttempt = strings.Replace(previousAttempt, "ungia", "ungi", -1)
	}

	// tìftiä replacement
	if w.Navi == "tìftia" && (strings.HasSuffix(target, "tìftiä") || strings.HasSuffix(target, "sìftiä")) {
		previousAttempt = strings.Replace(previousAttempt, "ìftia", "ìfti", -1)
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

	// sì-us lenition check
	if Contains(matchPrefixes, []string{"sì"}) {
		w.Affixes.Lenition = append(w.Affixes.Lenition, "t→s")
	}

	// only allow lenition after lenition-causing prefixes when prefixes and lenition present
	lenPre := []string{"pep", "pem", "pe", "fray", "tsay", "fay", "pay", "ay", "me", "pxe"}
	lenTry := ""
	lenResult := ""
	rootFirst := ""
	startWithD := []string{"kx", "px", "tx", "ts"}
	excludeStart := []string{"'rr", "'ll"}
	strTry := previousAttempt
	strRoot := w.Navi

	//trying to take first phoneme from attempt word
	if len(attempt) > 0 {
		lenTry = attempt[0:1]
	} else if len(previousAttempt) > 0 {
		if !HasPrefixStrArr(strRoot, excludeStart) {
			if HasPrefixStrArr(strTry, startWithD) {
				lenTry = previousAttempt[0:2]
			} else {
				lenTry = previousAttempt[0:1]
			}
		}
	}

	lenResult = w.plainLenite(lenTry)

	//trying to take first phoneme from root word
	if !HasPrefixStrArr(strRoot, excludeStart) {
		if HasPrefixStrArr(strRoot, startWithD) {
			rootFirst = w.Navi[0:2]
		} else {
			rootFirst = w.Navi[0:1]
		}
	}

	//hack for disappeared tìftang
	if rootFirst == "'" && lenResult == "" {
		lenResult = "-"
	}

	//treat sì- prefix like already lenited
	if Contains(matchPrefixes, []string{"sì"}) {
		lenResult = ""
	}

	if len(w.Affixes.Lenition) > 0 && len(matchPrefixes) > 0 {
		//this check if word starts from 'll or 'rr
		if HasPrefixStrArr(strRoot, excludeStart) {
			return previousAttempt
		}
		//don't lenite if there are those prefixes
		if Contains(matchPrefixes, []string{"fne", "munsna", "nì", "tì"}) {
			//handling -usia snowflakes
			if Contains(matchPrefixes, []string{"sì"}) && strings.HasSuffix(target, "usia") {
				previousAttempt = strings.Replace(previousAttempt, "usia", "", 1)
			}
			return previousAttempt
		}
		//don't lenite if there are those prefixes
		if Contains(matchPrefixes, []string{"fì", "tsa", "fra"}) && !(Contains(matchPrefixes, []string{"me", "pxe"})) {
			//handling -usia snowflakes
			if Contains(matchPrefixes, []string{"sì"}) && strings.HasSuffix(target, "usia") {
				previousAttempt = strings.Replace(previousAttempt, "usia", "", 1)
			}
			return previousAttempt
		}
		//unleniting where should not
		//lenpre+fne|munsna+(lenited root or tì- prefix)
	} else if Contains(matchPrefixes, lenPre) && !Contains(matchPrefixes, []string{"fne", "munsna"}) &&
		(((lenResult != rootFirst) && (lenResult != "") && (rootFirst != "")) ||
			(Contains(matchPrefixes, []string{"tì"}))) {
		//handling -usia snowflakes
		if strings.HasSuffix(previousAttempt, "usia") {
			previousAttempt = strings.Replace(previousAttempt, "usia", "", 1)
		}
		return previousAttempt
	}

	// build what prefixes to put on
	for _, p := range matchPrefixes {
		attempt = attempt + p
	}

	previousAttempt = attempt + previousAttempt

	matchPrefixes = DeleteElement(matchPrefixes, "e")
	//delete a- only if not number or adjective
	if (w.PartOfSpeech != "num.") && (w.PartOfSpeech != "adj.") {
		matchPrefixes = DeleteElement(matchPrefixes, "a")
	}

	//undelete a- for participles
	inf := w.Affixes.Infix
	if (Contains(inf, []string{"awn"}) || Contains(inf, []string{"us"})) &&
		!(Contains(matchPrefixes, []string{"tì"}) || Contains(matchPrefixes, []string{"sì"})) {
		matchPrefixes = append(matchPrefixes, "a")
	}

	matchPrefixes = DeleteElement(matchPrefixes, "ì")

	if ArrCount(matchPrefixes, "pe") == 2 {
		matchPrefixes = DeleteElement(matchPrefixes, "pe")
		matchPrefixes = append([]string{"pe", "pxe"}, matchPrefixes...)
	}

	if len(matchPrefixes) > 0 {
		// sì-us lenition append
		if Contains(matchPrefixes, []string{"sì"}) {
			matchPrefixes = DeleteElement(matchPrefixes, "sì")
			matchPrefixes = append(matchPrefixes, "tì")
		}
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
		// -to as suffix to nouns and made -sì|-to as separate option
		nSufRe string = "(nga'|tsyìp|tu|fkeyk)?(o)?(pe)?(mungwrr|kxamlä|tafkip|pxisre|pximaw|ftumfa|mìkam|nemfa|takip|lisre|talun|krrka|teri|fkip|pxaw|pxel|luke|rofa|fpi|ftu|kip|vay|lok|maw|sìn|sre|few|kam|kay|nuä|sko|yoa|äo|eo|fa|hu|ka|mì|na|ne|ta|io|uo|ro|wä|ìri|ìl|eyä|yä|ä|it|ri|ru|ti|ur|l|r|t)?(to|sì)?$"
		ngey   string = "ngey"
	)

	// hardcoded hack for tseyä
	if target == "tseyä" && w.Navi == "tsaw" {
		w.Affixes.Suffix = []string{"yä"}
		return "tseyä"
	}

	//denying suffixes on tsari, tsar, tsat, tsal, aylaru, sneyä
	declensedPronouns := []string{"tsal", "tsat", "tsar", "tsari", "aylaru", "sneyä", "sat"}
	checkWord := []string{w.Navi}
	targetTmp := []string{target}
	if Contains(checkWord, declensedPronouns) || Contains(targetTmp, declensedPronouns) {
		w.Affixes.Suffix = []string{""}
		return ""
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
		// meuiä support
	} else if w.Navi == "meuia" && strings.HasSuffix(target, "meuiä") {
		previousAttempt = strings.Replace(previousAttempt, "meuia", "meui", -1)
		reString = previousAttempt + reString
		// kemuiä support
	} else if w.Navi == "kemuia" && (strings.HasSuffix(target, "kemuiä") || strings.HasSuffix(target, "hemuiä")) {
		previousAttempt = strings.Replace(previousAttempt, "emuia", "emui", -1)
		reString = previousAttempt + reString
		// aungiä support
	} else if w.Navi == "aungia" && strings.HasSuffix(target, "ungiä") {
		previousAttempt = strings.Replace(previousAttempt, "ungia", "ungi", -1)
		reString = previousAttempt + reString
		// tìftiä support
	} else if w.Navi == "tìftia" && (strings.HasSuffix(target, "tìftiä") || strings.HasSuffix(target, "sìftiä")) {
		previousAttempt = strings.Replace(previousAttempt, "ìftia", "ìfti", -1)
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
		oldTarget       = target
		corTarget       = target
	)

	// hardcode for ner (n<0><er>rr)
	if w.Navi == "nrr" && (strings.HasSuffix(target, "er")) {
		w.InfixLocations = strings.Replace(w.InfixLocations, "<1><2>rr", "<1>rr", 1)
	}

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
		// one syllable <ol>ll
	} else if strings.Contains(reString, "<2>ll") {
		reString = strings.Replace(reString, "<1>", pos1InfixRe, 1)
		reString = strings.Replace(reString, "ll", "(ll)?", 1)
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

	// here goes er to rr change
	if strings.Contains(pos1InfixString, "er") && strings.Contains(reString, "?(rr)?") && w.Navi != "nrr" {
		corTarget = strings.Replace(corTarget, "er", "rr", 1)
		matchInfixes = remove(matchInfixes, "er")
		checkComment := fmt.Sprintf("correct form of %s is %s", oldTarget, corTarget)
		w.Affixes.Comment = []string{checkComment}
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
	attempt = strings.Replace(attempt, "olll", "ol", 1)
	attempt = strings.Replace(attempt, "errr", "er", 1)

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

// short table of all the possible lenitions
var shortLenitionTable = [4][2]string{
	{"kx, px, tx", "k, p, t"},
	{"k, p, t", "h, f, s"},
	{"ts", "s"},
	{"'", ""},
}

func GetShortLenitionTable() [][2]string {
	return shortLenitionTable[:]
}

// table of all the possible translations of "that"
var thatTable = [9][5]string{
	{"Case", "Noun", "   Clause Wrapper   ", "", ""},
	{" ", " ", "Prox.", "Dist.", "Answer "},
	{"====", "=====", "=====", "======", "======="},
	{"Sub.", "tsaw", "fwa", "tsawa", "teynga "},
	{"Agt.", "tsal", "fula", "tsala", "teyngla"},
	{"Pat.", "tsat", "futa", "tsata", "teyngta"},
	{"Gen.", "tseyä", "N/A", "N/A", ""},
	{"Dat.", "tsar", "fura", "tsara", ""},
	{"Top.", "tsari", "furia", "tsaria", ""},
}

func GetThatTable() [][5]string {
	return thatTable[:]
}

// table of all the possible translations of "that"
var otherThats = [8][3]string{
	{"tsa'u", "n.", "that (thing)"},
	{"tsakem", "n.", "that (action)"},
	{"fmawnta", "sbd.", "that news"},
	{"fayluta", "sbd.", "these words"},
	{"tsnì", "sbd.", "that (function word)"},
	{"tsonta", "conj.", "to (with kxìm)"},
	{"kuma/akum", "conj.", "that (as a result)"},
	{"a", "part.", "clause  level attributive marker"},
}

func GetOtherThats() [][3]string {
	return otherThats[:]
}

// plain lenite for backward checks
func (w *Word) plainLenite(tries string) string {
	for _, v := range lenitionTable {
		if strings.HasPrefix(strings.ToLower(tries), v[0]) {
			tries = strings.Replace(tries, v[0], v[1], 1)
			return tries
		}
	}
	return tries
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

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// Reconstruct is the main function of affixes.go, responsible for the affixing algorithm
// This will try to reconstruct a Word, so it matches with the target.
// Returns true if word got reconstructed into target!
func (w *Word) reconstruct(target string) bool {
	attempt := w.Navi

	//tì'usiä replacement
	var oldTarget = ""
	var origTarget = ""
	var corTarget = ""
	if w.Navi == "'ia" {
		if strings.HasSuffix(target, "tì'usiaä") || strings.HasSuffix(target, "sì'usiaä") {
			origTarget = target
			target = strings.Replace(target, "usiaä", "usiä", -1)
		} else if strings.HasSuffix(target, "tì'usiayä") || strings.HasSuffix(target, "sì'usiayä") {
			origTarget = target
			target = strings.Replace(target, "usiayä", "usiä", -1)
		}

		if strings.HasSuffix(target, "tì'usiä") || strings.HasSuffix(target, "sì'usiä") {
			oldTarget = target
			target = strings.Replace(target, "usiä", "usia", -1)
		}
	}
	//tìftusiä replacement
	if w.Navi == "ftia" {
		if strings.HasSuffix(target, "tìftusiaä") || strings.HasSuffix(target, "sìftusiaä") {
			origTarget = target
			target = strings.Replace(target, "usiaä", "usiä", -1)
		} else if strings.HasSuffix(target, "tìftusiayä") || strings.HasSuffix(target, "sìftusiayä") {
			origTarget = target
			target = strings.Replace(target, "usiayä", "usiä", -1)
		}

		if strings.HasSuffix(target, "tìftusiä") || strings.HasSuffix(target, "sìftusiä") {
			oldTarget = target
			target = strings.Replace(target, "usiä", "usia", -1)
		}
	}

	// detect errr in rr verb
	var rrVerbs = []string{"'rrko", "frrfen", "vrrìn", "nrr"}
	checkWord := []string{w.Navi}
	if Contains(checkWord, rrVerbs) && strings.Contains(target, "errr") {
		oldTarget = target
		target = strings.Replace(target, "errr", "er", 1)
	}
	// detect olll in ll verb
	var llVerbs = []string{"mll'an", "mllte", "plltxe", "pllhrr", "pllngay", "vll"}
	checkWord = []string{w.Navi}
	if Contains(checkWord, llVerbs) && strings.Contains(target, "olll") {
		oldTarget = target
		target = strings.Replace(target, "olll", "ol", 1)
	}

	// only try to infix verbs
	if strings.HasPrefix(w.PartOfSpeech, "v") || strings.HasPrefix(w.PartOfSpeech, svin) {
		attempt = w.infix(target)

		if debugMode {
			log.Printf("Navi: %s | Attempt: %s | Target: %s\n", w.Navi, attempt, target)
		}

		// altering errr output
		if Contains(checkWord, rrVerbs) && strings.Contains(oldTarget, "errr") {
			if w.Navi == "nrr" {
				corTarget = strings.Replace(oldTarget, "errr", "er", 1)
			} else {
				corTarget = strings.Replace(oldTarget, "errr", "rr", 1)
			}
			checkComment := fmt.Sprintf("correct form of %s is %s", oldTarget, corTarget)
			w.Affixes.Comment = []string{checkComment}
		}
		// altering olll output
		if Contains(checkWord, llVerbs) && strings.Contains(oldTarget, "olll") {
			corTarget = strings.Replace(oldTarget, "olll", "ol", 1)
			checkComment := fmt.Sprintf("correct form of %s is %s", oldTarget, corTarget)
			w.Affixes.Comment = []string{checkComment}
		}

		if attempt == target {
			return true
		}
	}

	// -iayä and -iaä replacements
	iaWords := []string{"aungia", "meuia", "soaia", "tìftia", "kemuia"}
	checkWord = []string{w.Navi}
	if (Contains(checkWord, iaWords) && strings.HasSuffix(target, "iayä")) && (!strings.HasSuffix(target, "usiayä")) {
		oldTarget = target
		target = strings.Replace(target, "iayä", "iä", -1)
	} else if (Contains(checkWord, iaWords) && strings.HasSuffix(target, "iaä")) && (!strings.HasSuffix(target, "usiaä")) {
		oldTarget = target
		target = strings.Replace(target, "iaä", "iä", -1)
	}

	attempt = w.prefix(target, attempt)

	if debugMode {
		log.Println("PREFIX w")
		log.Printf("Navi: %s | Attempt: %s | Target: %s\n", w.Navi, attempt, target)
	}

	//tìftusiä & tì'usiä append
	usiaWords := []string{"'ia", "ftia"}
	uscheckWord := []string{w.Navi}
	if Contains(uscheckWord, usiaWords) {
		if strings.HasSuffix(oldTarget, "usiä") {
			if strings.HasSuffix(origTarget, "usiayä") || strings.HasSuffix(origTarget, "usiaä") {
				checkComment := fmt.Sprintf("correct form of %s is %s", origTarget, oldTarget)
				w.Affixes.Comment = []string{checkComment}
			} else {
				w.Affixes.Suffix = append(w.Affixes.Suffix, "ä")
			}
			//now only if attempt is -usia too
			if strings.HasSuffix(attempt, "usia") {
				attempt = target
			}
		}
	}

	if attempt == target {
		return true
	}

	attempt = w.suffix(target, attempt)

	if debugMode {
		log.Println("SUFFIX w")
		log.Printf("Navi: %s | Attempt: %s | Target: %s\n", w.Navi, attempt, target)
	}

	// -iayä or -iaä results append
	checkWord = []string{w.Navi}
	iaSuffixesToSuppress := []string{"ä", "yä"}
	checkIaSuffixes := w.Affixes.Suffix
	if (Contains(checkWord, iaWords) && (strings.HasSuffix(oldTarget, "iayä") || strings.HasSuffix(oldTarget, "iaä"))) && (!strings.HasSuffix(target, "usiaä") || !strings.HasSuffix(target, "usiayä")) {
		checkComment := fmt.Sprintf("correct form of %s is %s", oldTarget, target)
		w.Affixes.Comment = []string{checkComment}
		if Contains(checkIaSuffixes, iaSuffixesToSuppress) {
			w.Affixes.Suffix = remove(w.Affixes.Suffix, "ä")
			w.Affixes.Suffix = remove(w.Affixes.Suffix, "yä")
		}
	}

	if attempt == target {
		return true
	}

	attempt = w.lenite(w.Navi)

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

	// -iayä or -iaä results append
	checkWord = []string{w.Navi}
	checkIaSuffixes = w.Affixes.Suffix
	if (Contains(checkWord, iaWords) && (strings.HasSuffix(oldTarget, "iayä") || strings.HasSuffix(oldTarget, "iaä"))) && (!strings.HasSuffix(target, "usiaä") || !strings.HasSuffix(target, "usiayä")) {
		checkComment := fmt.Sprintf("correct form of %s is %s", oldTarget, target)
		w.Affixes.Comment = []string{checkComment}
		if Contains(checkIaSuffixes, iaSuffixesToSuppress) {
			w.Affixes.Suffix = remove(w.Affixes.Suffix, "ä")
			w.Affixes.Suffix = remove(w.Affixes.Suffix, "yä")
		}
	}

	if attempt == target {
		return true
	}

	if debugMode {
		log.Println("GIVING UP")
	}

	return false
}
