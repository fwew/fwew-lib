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

// Package main contains all the things
package fwew_lib

import (
	"log"
	"math/rand"
	"strings"
)

// Global
const (
	space string = " "
)

var debugMode bool

func intersection(a, b string) (c string) {
	m := make(map[rune]bool)
	for _, r := range a {
		m[r] = true
	}
	for _, r := range b {
		if _, ok := m[r]; ok {
			c += string(r)
		}
	}
	return
}

func (w *Word) similarity(other string) float64 {
	if w.Navi == other {
		return 1.0
	}
	if len(w.Navi) > len(other)+1 {
		return 0.0
	}
	if w.Navi == "nga" && other == "ngey" {
		return 1.0
	}
	if w.Navi == "'ia" && strings.HasSuffix(other, "ì'usiä") {
		return 1.0
	}
	vowels := "aäeiìoulr"
	w0v := intersection(w.Navi, vowels)
	w1v := intersection(other, vowels)
	wis := intersection(w.Navi, other)
	wav := intersection(w0v, other)
	if len(w0v) > len(w1v) {
		return 0.0
	}
	if len(wav) == 0 {
		return 0.0
	}
	scc := len(wis)
	iratio := float64(scc) / float64(len(w.Navi))
	lratio := float64(len(w.Navi)) / float64(len(other))
	return (iratio + lratio) / 2
}

func identicalRunes(first string, second string) bool {
	a := []rune(first)
	b := []rune(second)

	if len(a) != len(b) {
		return false
	}

	for i, c := range a {
		if b[i] != c {
			return false
		}
	}

	return true
}

// Translate some navi text.
// !! Only one word is allowed, if spaces are found, they will be treated like part of the word !!
// This will return an array of Words, that fit the input text
// One Navi-Word can have multiple meanings and words (e.g. synonyms)
func TranslateFromNaviHash(searchNaviWord string, checkFixes bool) (results []Word, err error) {
	badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\`

	// remove all the sketchy chars from arguments
	for _, c := range badChars {
		searchNaviWord = strings.ReplaceAll(searchNaviWord, string(c), "")
	}

	// normalize tìftang character
	searchNaviWord = strings.ReplaceAll(searchNaviWord, "’", "'")
	searchNaviWord = strings.ReplaceAll(searchNaviWord, "‘", "'")

	// find everything lowercase
	searchNaviWord = strings.ToLower(searchNaviWord)

	// No Results if empty string after removing sketch chars
	if len(searchNaviWord) == 0 {
		return
	}

	// Find the word
	for _, b := range dictHash[searchNaviWord] {
		if _, ok := dictHash[searchNaviWord]; ok {
			results = append(results, b)
		}
	}

	if checkFixes {
		conjugations := deconjugate(searchNaviWord)
		for _, candidate := range conjugations {
			for _, c := range dictHash[candidate.word] {
				if _, ok := dictHash[candidate.word]; ok {
					if candidate.insistPOS == "n." {
						posNoun := c.PartOfSpeech
						//posNoun == "n." || posNoun == "prop.n." || posNoun == "pn."
						if strings.HasSuffix(posNoun, "n.") && !strings.HasPrefix(posNoun, "v") {
							a := c
							a.Affixes.Lenition = candidate.lenition
							a.Affixes.Prefix = candidate.prefixes
							a.Affixes.Suffix = candidate.suffixes
							results = append(results, a)
						}

					} else if candidate.insistPOS == "adj." {
						posNoun := c.PartOfSpeech
						if posNoun == "adj." || posNoun == "num." {
							a := c
							a.Affixes.Lenition = candidate.lenition
							a.Affixes.Prefix = candidate.prefixes
							a.Affixes.Suffix = candidate.suffixes
							results = append(results, a)
						}
					} else if candidate.insistPOS == "v." {
						posNoun := c.PartOfSpeech
						if strings.HasPrefix(posNoun, "v") {
							a := c
							a.Affixes.Lenition = candidate.lenition
							a.Affixes.Prefix = candidate.prefixes
							a.Affixes.Suffix = candidate.suffixes
							a.Affixes.Infix = candidate.infixes

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

							// normal infixes
							if identicalRunes(rebuiltVerb, searchNaviWord) {
								results = append(results, a)
							} else if identicalRunes("tì"+rebuiltVerb, searchNaviWord) && len(candidate.infixes) == 1 && candidate.infixes[0] == "us" {
								// tì + v<us>erb constructions
								results = append(results, a)
							} else if firstInfixes == "us" || firstInfixes == "awn" {
								if identicalRunes("a"+rebuiltVerb, searchNaviWord) {
									// a-v<us>erb and a-v<awn>erb
									results = append(results, a)
								} else if identicalRunes(rebuiltVerb+"a", searchNaviWord) {
									// v<us>erb-a and v<awn>erb-a
									results = append(results, a)
								}
							}
						}
					} else if candidate.insistPOS == "nì." {
						posNoun := c.PartOfSpeech
						if posNoun == "adj." || posNoun == "pn." {
							a := c
							a.Affixes.Lenition = candidate.lenition
							a.Affixes.Prefix = candidate.prefixes
							a.Affixes.Suffix = candidate.suffixes
							results = append(results, a)
						}
					} else {
						a := c
						a.Affixes.Lenition = candidate.lenition
						a.Affixes.Prefix = candidate.prefixes
						a.Affixes.Suffix = candidate.suffixes
						results = append(results, a)
					}
				}
			}
		}
	}

	return
}

func SearchNatlangWord(wordmap map[string][]string, searchWord string) (results []Word) {

	// No Results if empty string after removing sketch chars
	if len(searchWord) == 0 {
		return
	}

	// Find the word
	if _, ok := wordmap[searchWord]; !ok {
		return results
	}

	firstResults := wordmap[searchWord]

	for i := 0; i < len(firstResults); i++ {
		for _, c := range dictHash[firstResults[i]] {
			results = append(results, c)
		}
	}

	return
}

func TranslateToNaviHash(searchWord string, langCode string) (results []Word) {
	badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\`

	// remove all the sketchy chars from arguments
	for _, c := range badChars {
		searchWord = strings.ReplaceAll(searchWord, string(c), "")
	}

	// normalize tìftang character
	searchWord = strings.ReplaceAll(searchWord, "’", "'")
	searchWord = strings.ReplaceAll(searchWord, "‘", "'")

	// find everything lowercase
	searchWord = strings.ToLower(searchWord)

	switch langCode {
	case "de":
		return SearchNatlangWord(dictHash2.DE, searchWord)
	case "en":
		return SearchNatlangWord(dictHash2.EN, searchWord)
	case "et":
		return SearchNatlangWord(dictHash2.ET, searchWord)
	case "fr":
		return SearchNatlangWord(dictHash2.FR, searchWord)
	case "hu":
		return SearchNatlangWord(dictHash2.HU, searchWord)
	case "nl":
		return SearchNatlangWord(dictHash2.NL, searchWord)
	case "pl":
		return SearchNatlangWord(dictHash2.PL, searchWord)
	case "ru":
		return SearchNatlangWord(dictHash2.RU, searchWord)
	case "sv":
		return SearchNatlangWord(dictHash2.SV, searchWord)
	case "tr":
		return SearchNatlangWord(dictHash2.TR, searchWord)
	}

	// If we get an odd language code, return English
	return SearchNatlangWord(dictHash2.EN, searchWord)
}

// Get random words out of the dictionary.
// If args are applied, the dict will be filtered for args before random words are chosen.
// args will be put into the `List()` algorithm.
func Random(amount int, args []string) (results []Word, err error) {
	allWords, err := List(args)

	if err != nil {
		log.Printf("Error getting fullDing: %s", err)
		return
	}

	dictLength := len(allWords)

	if dictLength == 0 {
		return nil, NoResults
	}

	// create random number
	if amount <= 0 {
		amount = rand.Intn(dictLength) + 1
	}

	if amount > dictLength {
		return allWords, nil
	}

	// get random numbers for allWords array
	perm := rand.Perm(dictLength)

	for _, i := range perm[:amount] {
		results = append(results, allWords[i])
	}

	return
}

func StartEverything() {
	AssureDict()
	CacheDictHash()
	CacheDictHash2()
	PhonemeDistros()
}
