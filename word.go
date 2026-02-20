package fwew_lib

import (
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

func (w *Word) String() string {
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
		"IT: %s\n"+
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
		w.IT,
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
	word.IT = dataFields[order.itField]
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

	sylBuilder := ""
	syl := ""
	allSpaces := strings.Split(w.Syllables, " ")
	for i, a := range allSpaces {
		if a == "or" {
			sylBuilder = strings.Trim(sylBuilder, " -")
			input, err := w.doUnderline(sylBuilder, withMarkdown)
			if err != nil {
				return "", err
			}
			syl += input + " or "
			sylBuilder = ""
			continue
		}
		sylBuilder += a + " "
		if i+1 == len(allSpaces) {
			sylBuilder = strings.Trim(sylBuilder, " -")
			input, err := w.doUnderline(sylBuilder, withMarkdown)
			if err != nil {
				return "", err
			}
			syl += input
		}
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
	case "it":
		output += w.IT
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
		output += "\n(Reef Na'vi: "
		reefBreakdownInput := ""
		allSpaces := strings.Split(reefy[0], " ")
		for i, a := range allSpaces {
			if a == "or" {
				reefBreakdownInput = strings.Trim(reefBreakdownInput, " -")
				input, err := w.doUnderline(reefBreakdownInput, withMarkdown)
				if err != nil {
					return "", err
				}
				output += input + " or "
				reefBreakdownInput = ""
				continue
			}
			reefBreakdownInput += a + " "
			if i+1 == len(allSpaces) {
				reefBreakdownInput = strings.Trim(reefBreakdownInput, " -")
				input, err := w.doUnderline(reefBreakdownInput, withMarkdown)
				if err != nil {
					return "", err
				}
				output += input
			}
		}
		if showIPA {
			output += " [" + reefy[1] + "]"
		}
		output += ")"
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

	mdUnderline := "__"
	shUnderlineA := "\033[4m"
	shUnderlineB := "\033[0m"

	// get it from the IPA
	var stressed []bool
	for _, a := range strings.Split(w.IPA, " ") {
		if a == "or" {
			break
		}
		for _, b := range strings.Split(a, ".") {
			if strings.Contains(b, "ˈ") {
				stressed = append(stressed, true)
			} else {
				stressed = append(stressed, false)
			}
		}
	}

	// apply it from the IPA
	i := 0
	underlined := ""
	for _, a := range strings.Split(syllables, " ") {
		for _, b := range strings.Split(a, "-") {
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
