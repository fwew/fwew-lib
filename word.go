package fwew_lib

import (
	"encoding/json"
	"fmt"
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
	IT             string
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

// affixes struct has its own type, so it is automatically copied :)
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

func (w *Word) String() string {
	s, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(s)
}

// equal returns true if Word b is a duplicate of Word a
func equal(a, b Word) bool {
	return a.ID == b.ID &&
		len(a.Affixes.Prefix) == len(b.Affixes.Prefix) &&
		len(a.Affixes.Suffix) == len(b.Affixes.Suffix) &&
		len(a.Affixes.Lenition) == len(b.Affixes.Lenition) &&
		len(a.Affixes.Infix) == len(b.Affixes.Infix)
}

// Make a simple word to show what query led to this word
func simpleWord(name string) Word {
	var word Word
	word.Navi = name
	return word
}

// Initialize Word with one row of the dictionary.
func newWord(dataFields []string, order dictPos) Word {
	var w Word

	// Keep these two lists in the same order.
	fields := []*string{
		&w.ID, &w.Navi, &w.IPA, &w.InfixLocations, &w.PartOfSpeech, &w.Source, &w.Stressed, &w.Syllables, &w.InfixDots,
		&w.DE, &w.EN, &w.ES, &w.ET, &w.FR, &w.HU, &w.IT, &w.KO, &w.NL, &w.PL, &w.PT, &w.RU, &w.SV, &w.TR, &w.UK,
	}

	indexes := []int{
		order.idField, order.navField, order.ipaField, order.infField, order.posField, order.srcField, order.stsField,
		order.sylField, order.ifdField,
		order.deField, order.enField, order.esField, order.etField, order.frField, order.huField, order.itField,
		order.koField, order.nlField, order.plField, order.ptField, order.ruField, order.svField, order.trField,
		order.ukField,
	}

	// Assign safely (in case an index is missing/out of range).
	set := func(dst *string, idx int) {
		if idx >= 0 && idx < len(dataFields) {
			*dst = dataFields[idx]
		}
	}

	for i := range fields {
		set(fields[i], indexes[i])
	}

	return w
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

	syl, err := formatSyllables(w, w.Syllables, withMarkdown)
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

	var langDefMap = map[string]string{
		"de": w.DE, "en": w.EN, "es": w.ES, "et": w.ET, "fr": w.FR, "hu": w.HU, "it": w.IT, "ko": w.KO,
		"nl": w.NL, "pl": w.PL, "pt": w.PT, "ru": w.RU, "sv": w.SV, "tr": w.TR, "uk": w.UK,
	}

	if langDef, ok := langDefMap[langCode]; ok {
		output += langDef
	} else {
		output += w.EN
	}

	if reef {
		reefy := ReefMe(w.IPA, false)
		output += "\n(Reef Na'vi: "
		formattedReef, err := formatSyllables(w, reefy[0], withMarkdown)
		if err != nil {
			return "", err
		}
		output += formattedReef
		if showIPA {
			output += " [" + reefy[1] + "]"
		}
		output += ")"
	}

	var affixMap = map[string]string{
		Text("prefixes"): strings.Join(runOn(w.Affixes.Prefix, applyPrefixNotation), ", "),
		Text("infixes"):  strings.Join(runOn(w.Affixes.Infix, applyInfixNotation), ", "),
		Text("suffixes"): strings.Join(runOn(w.Affixes.Suffix, applySuffixNotation), ", "),
		Text("lenition"): strings.Join(w.Affixes.Lenition, ", "),
		Text("comment"):  strings.Join(w.Affixes.Comment, ", "),
	}

	for k, v := range affixMap {
		if len(v) > 0 {
			output += fmt.Sprintf("%s%s- %s: %v", newline, indent, k, v)
		}
	}

	if showSource && w.Source != "" {
		output += newline + src
	}

	output += newline
	return
}

func formatSyllables(w *Word, syllables string, withMarkdown bool) (string, error) {
	var formattedSyllables strings.Builder
	var currentSyllableSet string
	allSpaces := strings.Split(syllables, " ")

	for i, a := range allSpaces {
		if a == "or" {
			currentSyllableSet = strings.Trim(currentSyllableSet, " -")
			input, err := doUnderline(w, currentSyllableSet, withMarkdown)
			if err != nil {
				return "", err
			}
			formattedSyllables.WriteString(input + " or ")
			currentSyllableSet = ""
			continue
		}

		currentSyllableSet += a + " "

		if i+1 == len(allSpaces) {
			currentSyllableSet = strings.Trim(currentSyllableSet, " -")
			input, err := doUnderline(w, currentSyllableSet, withMarkdown)
			if err != nil {
				return "", err
			}
			formattedSyllables.WriteString(input)
		}
	}

	return formattedSyllables.String(), nil
}

func doUnderline(w *Word, input string, markdown bool) (string, error) {
	syllables := w.Syllables
	if len(input) > 0 {
		syllables = input
	}
	if !strings.Contains(syllables, "-") || w.Stressed == "0" {
		return syllables, nil
	}

	mdUnderline := "__"
	shUnderlineA := "\033[4m"
	shUnderlineB := "\033[0m"

	// get it from the IPA
	var stressed []bool
	for a := range strings.SplitSeq(w.IPA, " ") {
		if a == "or" {
			break
		}
		for b := range strings.SplitSeq(a, ".") {
			if strings.Contains(b, "ˈ") {
				stressed = append(stressed, true)
			} else {
				stressed = append(stressed, false)
			}
		}
	}

	// runOn it from the IPA
	i := 0
	underlined := ""
	for a := range strings.SplitSeq(syllables, " ") {
		for b := range strings.SplitSeq(a, "-") {
			if i >= len(stressed) {
				break
			}
			if stressed[i] {
				if markdown {
					underlined += "-" + mdUnderline + b + mdUnderline
				} else {
					underlined += "-" + shUnderlineA + b + shUnderlineB
				}
			} else {
				underlined += "-" + b
			}
			i += 1
		}
		underlined += " "
	}

	underlined = strings.TrimPrefix(underlined, "-")
	underlined = strings.TrimSuffix(underlined, " ")
	underlined = strings.ReplaceAll(underlined, " -", " ")

	return underlined, nil
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
	itField  int // Italian definition
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
		case "it":
			pos.itField = i
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
