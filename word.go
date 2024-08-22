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

// Package fwew_lib contains all the things. word.go is home to the Word struct.
package fwew_lib

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Word is a struct that contains all the data about a given word
type Word struct {
	ID             string
	Navi           string
	IPA            string
	InfixLocations string
	PartOfSpeech   string
	Source         string
	Stressed       string
	Syllables      string
	InfixDots      string
	DE             string
	EN             string
	ES             string
	ET             string
	FR             string
	HU             string
	KO             string
	NL             string
	PL             string
	PT             string
	RU             string
	SV             string
	TR             string
	UK             string
	Affixes        affix
}

// affixes has its own type, so it is automatically copied :)
type affix struct {
	Prefix   []string
	Infix    []string
	Suffix   []string
	Lenition []string
	Comment  []string
}

func addAffixes(a affix, z affix) (w affix) {
	z.Prefix = append(z.Prefix, a.Prefix...)
	z.Infix = append(z.Infix, a.Infix...)
	z.Suffix = append(z.Suffix, a.Suffix...)
	z.Lenition = append(z.Lenition, a.Lenition...)
	z.Comment = append(z.Comment, a.Comment...)
	return z
}

func (w Word) String() string {
	// this string only doesn't get translated or called from Text() because they're var names
	return fmt.Sprintf(""+
		"Id: %s\n"+
		"Navi: %s\n"+
		"IPA: %s\n"+
		"InfixLocations: %s\n"+
		"PartOfSpeech: %s\n"+
		"Source: %s\n"+
		"Stressed: %s\n"+
		"Syllables: %s\n"+
		"InfixDots: %s\n"+
		"DE: %s\n"+
		"EN: %s\n"+
		"ES: %s\n"+
		"ET: %s\n"+
		"FR: %s\n"+
		"HU: %s\n"+
		"KO: %s\n"+
		"NL: %s\n"+
		"PL: %s\n"+
		"PT: %s\n"+
		"RU: %s\n"+
		"SV: %s\n"+
		"TR: %s\n"+
		"UK: %s\n"+
		"Affixes: %v\n",
		w.ID,
		w.Navi,
		w.IPA,
		w.InfixLocations,
		w.PartOfSpeech,
		w.Source,
		w.Stressed,
		w.Syllables,
		w.InfixDots,
		w.DE,
		w.EN,
		w.ES,
		w.ET,
		w.FR,
		w.HU,
		w.KO,
		w.NL,
		w.PL,
		w.PT,
		w.RU,
		w.SV,
		w.TR,
		w.UK,
		w.Affixes,
	)
}

// Make a simple word to show what query led to this word
func simpleWord(name string) Word {
	var word Word
	word.Navi = name
	return word
}

// Initialize Word with one row of the dictionary.
func newWord(dataFields []string, order dictPos) Word {
	var word Word
	word.ID = dataFields[order.idField]
	word.Navi = dataFields[order.navField]
	word.IPA = dataFields[order.ipaField]
	word.InfixLocations = dataFields[order.infField]
	word.PartOfSpeech = dataFields[order.posField]
	word.Source = dataFields[order.srcField]
	word.Stressed = dataFields[order.stsField]
	word.Syllables = dataFields[order.sylField]
	word.InfixDots = dataFields[order.ifdField]
	word.DE = dataFields[order.deField]
	word.EN = dataFields[order.enField]
	word.ES = dataFields[order.esField]
	word.ET = dataFields[order.etField]
	word.FR = dataFields[order.frField]
	word.HU = dataFields[order.huField]
	word.KO = dataFields[order.koField]
	word.NL = dataFields[order.nlField]
	word.PL = dataFields[order.plField]
	word.PT = dataFields[order.ptField]
	word.RU = dataFields[order.ruField]
	word.SV = dataFields[order.svField]
	word.TR = dataFields[order.trField]
	word.UK = dataFields[order.ukField]
	return word
}

// CloneWordStruct is basically a copy constructor for Word struct
// Basically not needed, cause go copies things by itself.
// Only string arrays in Affixes are pointers and therefore need manual copy.
func (w *Word) CloneWordStruct() Word {
	// Copy struct to new instance
	nw := *w

	// copy the arrays manually
	copy(nw.Affixes.Prefix, w.Affixes.Prefix)
	copy(nw.Affixes.Infix, w.Affixes.Infix)
	copy(nw.Affixes.Suffix, w.Affixes.Suffix)
	copy(nw.Affixes.Lenition, w.Affixes.Lenition)
	copy(nw.Affixes.Comment, w.Affixes.Comment)

	return nw
}

func (w *Word) Equals(other Word) bool {
	return w.ID == other.ID &&
		w.Navi == other.Navi &&
		w.IPA == other.IPA &&
		w.InfixLocations == other.InfixLocations &&
		w.PartOfSpeech == other.PartOfSpeech &&
		w.Source == other.Source &&
		w.Stressed == other.Stressed &&
		w.Syllables == other.Syllables &&
		w.InfixDots == other.InfixDots &&
		w.DE == other.DE &&
		w.EN == other.EN &&
		w.ES == other.ES &&
		w.ET == other.ET &&
		w.FR == other.FR &&
		w.HU == other.HU &&
		w.KO == other.KO &&
		w.NL == other.NL &&
		w.PL == other.PL &&
		w.PT == other.PT &&
		w.RU == other.RU &&
		w.SV == other.SV &&
		w.TR == other.TR &&
		w.UK == other.UK &&
		reflect.DeepEqual(w.Affixes, other.Affixes)
}

func (w *Word) SyllableCount() int {
	var numSyllables int
	var vowels = []string{"a", "ä", "e", "é", "i", "ì", "o", "u", "ll", "rr"}
	var word = strings.ToLower(w.Navi)
	for _, p := range vowels {
		numSyllables += strings.Count(word, p)
	}
	return numSyllables
}

const (
	mdBold   = "**"
	mdItalic = "*"
	newline  = "\n"
	valNull  = "NULL"
)

func (w *Word) ToOutputLine(
	i string,
	withMarkdown, showIPA, showInfixes, showDashed, showInfDots, showSource, reef bool,
	langCode string,
) (output string, err error) {
	num := fmt.Sprintf("[%s]", i)
	nav := w.Navi
	ipa := fmt.Sprintf("[%s]", w.IPA)
	pos := "" + w.PartOfSpeech
	inf := "" + w.InfixLocations
	src := fmt.Sprintf("%s: %s", Text("src"), w.Source)
	ifd := "" + w.InfixDots

	syl, err := w.doUnderline("", withMarkdown)
	if err != nil {
		return "", err
	}

	if withMarkdown {
		nav = mdBold + nav + mdBold
		pos = mdItalic + pos + mdItalic
	}

	output += num + space + nav + space

	if showIPA {
		output += strings.ReplaceAll(ipa, "ʊ", "u") + space
		if strings.Contains(ipa, "ʊ") {
			output += "or " + ipa + space
		}
	}

	if showInfixes && w.InfixLocations != valNull {
		output += inf + space
	}

	if showDashed {
		output += "(" + syl
		if showInfDots && w.InfixDots != valNull {
			output += "," + space
		} else {
			output += ")" + space
		}
	}

	if showInfDots && w.InfixDots != valNull {
		if !showDashed {
			output += "("
		}
		output += ifd + ")" + space
	}

	output += pos + space

	switch langCode {
	case "de":
		output += w.DE
	case "en":
		output += w.EN
	case "es":
		output += w.ES
	case "et":
		output += w.ET
	case "fr":
		output += w.FR
	case "hu":
		output += w.HU
	case "ko":
		output += w.KO
	case "nl":
		output += w.NL
	case "pl":
		output += w.PL
	case "pt":
		output += w.PT
	case "ru":
		output += w.RU
	case "sv":
		output += w.SV
	case "tr":
		output += w.TR
	case "uk":
		output += w.UK
	default:
		output += w.EN
	}

	if reef {
		reefy := ReefMe(w.IPA, false)
		reefBreakdown, err := w.doUnderline(reefy[0], withMarkdown)
		if err != nil {
			return "", err
		}
		output += "\n(Reef Na'vi: " + reefBreakdown + " [" + reefy[1] + "])"
	}

	if len(w.Affixes.Prefix) > 0 {
		output += newline + fmt.Sprintf("Prefixes: %s", w.Affixes.Prefix)
	}
	if len(w.Affixes.Infix) > 0 {
		output += newline + fmt.Sprintf("Infixes: %s", w.Affixes.Infix)
	}
	if len(w.Affixes.Suffix) > 0 {
		output += newline + fmt.Sprintf("Suffixes: %s", w.Affixes.Suffix)
	}
	if len(w.Affixes.Lenition) > 0 {
		output += newline + fmt.Sprintf("Lenition: %s", w.Affixes.Lenition)
	}
	if len(w.Affixes.Comment) > 0 {
		output += newline + fmt.Sprintf("Comment: %s", w.Affixes.Comment)
	}
	if showSource && w.Source != "" {
		output += newline + src
	}

	output += newline
	return
}

func (w *Word) doUnderline(input string, markdown bool) (string, error) {
	syllables := w.Syllables
	if len(input) > 0 {
		syllables = input
	}
	if !strings.Contains(syllables, "-") || w.Stressed == "0" {
		return syllables, nil
	}

	var err error
	mdUnderline := "__"
	shUnderlineA := "\033[4m"
	shUnderlineB := "\033[0m"
	dashed := syllables
	dSlice := strings.FieldsFunc(dashed, func(r rune) bool {
		return r == '-' || r == ' '
	})

	stressedIndex, err := strconv.Atoi(w.Stressed)
	if err != nil {
		return "", InvalidNumber.wrap(err)
	}
	stressedSyllable := dSlice[stressedIndex-1]

	if strings.Contains(stressedSyllable, " ") {
		tmp := strings.Split(stressedSyllable, " ")
		if markdown {
			tmp[0] = mdUnderline + tmp[0] + mdUnderline
		} else {
			tmp[0] = shUnderlineA + tmp[0] + shUnderlineB
		}
		stressedSyllable = strings.Join(tmp, " ")
		dSlice[stressedIndex-1] = stressedSyllable
		return strings.Join(dSlice, "-"), nil
	} else {
		if markdown {
			dSlice[stressedIndex-1] = mdUnderline + stressedSyllable + mdUnderline
		} else {
			dSlice[stressedIndex-1] = shUnderlineA + stressedSyllable + shUnderlineB
		}
		return strings.Join(dSlice, "-"), nil
	}
}

// This holds the positions, how the fields are sorted (with dictV2 the fields can have any order)
type dictPos struct {
	idField  int // Database ID
	navField int // Na'vi word
	ipaField int // IPA data
	infField int // Infix location data
	posField int // Part of Speech data
	srcField int // Source data
	stsField int // Stressed syllable #
	sylField int // syllable breakdown
	ifdField int // dot-style infix data
	deField  int // German definition
	enField  int // English definition
	esField  int // Spanish definition
	etField  int // Estonian definition
	frField  int // French definition
	huField  int // Hungarian definition
	koField  int // Korean definition
	nlField  int // Dutch definition
	plField  int // Polish definition
	ptField  int // Portuguese definition
	ruField  int // Russian definition
	svField  int // Swedish definition
	trField  int // Turkish definition
	ukField  int // Ukrainian definition
}

func readDictPos(headerFields []string) dictPos {
	var pos dictPos

	for i, field := range headerFields {
		switch field {
		case "id":
			pos.idField = i
		case "navi":
			pos.navField = i
		case "ipa":
			pos.ipaField = i
		case "infixes":
			pos.infField = i
		case "partOfSpeech":
			pos.posField = i
		case "source":
			pos.srcField = i
		case "stressed":
			pos.stsField = i
		case "syllables":
			pos.sylField = i
		case "infixDots":
			pos.ifdField = i
		case "de":
			pos.deField = i
		case "en":
			pos.enField = i
		case "es":
			pos.esField = i
		case "et":
			pos.etField = i
		case "fr":
			pos.frField = i
		case "hu":
			pos.huField = i
		case "ko":
			pos.koField = i
		case "nl":
			pos.nlField = i
		case "pl":
			pos.plField = i
		case "pt":
			pos.ptField = i
		case "ru":
			pos.ruField = i
		case "sv":
			pos.svField = i
		case "tr":
			pos.trField = i
		case "uk":
			pos.ukField = i
		}
	}

	return pos
}
