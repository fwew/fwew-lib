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

/* To help deduce phonemes */
var romanization2 = map[string]string{
	// Vowels
	"a": "a", "i": "i", "ɪ": "ì",
	"o": "o", "ɛ": "e", "u": "u",
	"æ": "ä",
	// Diphthongs
	"aw": "aw", "ɛj": "ey",
	"aj": "ay", "ɛw": "ew",
	// Psuedovowels
	"ṛ": "rr", "ḷ": "ll",
	// Consonents
	"t": "t", "p": "p", "ʔ": "'",
	"n": "n", "k": "k", "l": "l",
	"s": "s", "ɾ": "r", "j": "y",
	"t͡s": "ts", "t'": "tx", "m": "m",
	"v": "v", "w": "w", "h": "h",
	"ŋ": "ng", "z": "z", "k'": "kx",
	"p'": "px", "f": "f", "r": "r",
	// Reef dialect
	"b": "b", "d": "d", "g": "g",
	"ʃ": "sh", "tʃ": "ch", "ʊ": "ù",
	// mistakes and rarities
	"ʒ": "ch", "": "", " ": ""}

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

func clean(searchNaviWords string) (words string) {
	badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\`

	// remove all the sketchy chars from arguments
	for _, c := range badChars {
		searchNaviWords = strings.ReplaceAll(searchNaviWords, string(c), "")
	}

	// normalize tìftang character
	searchNaviWords = strings.ReplaceAll(searchNaviWords, "’", "'")
	searchNaviWords = strings.ReplaceAll(searchNaviWords, "‘", "'")

	// find everything lowercase
	searchNaviWords = strings.ToLower(searchNaviWords)

	// normalize ù
	searchNaviWords = strings.ReplaceAll(searchNaviWords, "ù", "u")

	// No Results if empty string after removing sketch chars
	if len(searchNaviWords) == 0 {
		return
	}

	return searchNaviWords
}

// Translate some navi text.
// !! Multiple words are supported !!
// This will return a 2D array of Words that fit the input text
// The first word will only contain the query put into the translate command
// One Navi-Word can have multiple meanings and words (e.g. synonyms)
func TranslateFromNaviHash(searchNaviWords string, checkFixes bool) (results [][]Word, err error) {
	searchNaviWords = clean(searchNaviWords)

	// No Results if empty string after removing sketch chars
	if len(searchNaviWords) == 0 {
		return
	}

	allWords := strings.Split(clean(searchNaviWords), " ")

	i := 0

	results = [][]Word{}

	for i < len(allWords) {

		j, newWords, error2 := TranslateFromNaviHashHelper(i, allWords, checkFixes)
		if error2 == nil {
			for _, newWord := range newWords {
				// Set up receptacle for words
				results = append(results, []Word{})
				for _, newResult := range newWord {
					results[len(results)-1] = append(results[len(results)-1], newResult)
				}
			}
		}
		i += j
		i++
	}

	return
}

// Helper for TranslateFromNaviHashHelper
func AppendToFront(words []Word, input Word) []Word {
	// Get the query it's looking for
	dummyWord := []Word{words[0]}
	// Append it to the front of the list
	i := 1
	dummyWord = append(dummyWord, input)
	for i < len(words) {
		dummyWord = append(dummyWord, words[i])
		i++
	}
	// Make it the list
	return dummyWord
}

// Helper for TranslateFromNaviHashHelper
func IsVerb(input string) (result bool) {
	_, possibilities, err := TranslateFromNaviHashHelper(0, []string{input}, true)
	if err != nil {
		return false
	}
	for _, a := range possibilities {
		for _, b := range a {
			if len(b.PartOfSpeech) > 0 && b.PartOfSpeech[0] == 'v' {
				for _, c := range b.Affixes.Infix {
					// <us> and <awn> are participles, so they become adjectives
					if c == "us" || c == "awn" {
						return false
					}
				}
				return true
			}
		}
	}
	return false
}

func TranslateFromNaviHashHelper(start int, allWords []string, checkFixes bool) (steps int, results [][]Word, err error) {
	i := start

	searchNaviWord := allWords[i]
	results = [][]Word{{simpleWord(searchNaviWord)}}

	bareNaviWord := false
	// Find the word
	if _, ok := dictHash[searchNaviWord]; ok {
		bareNaviWord = true

		for _, b := range dictHash[searchNaviWord] {
			results[len(results)-1] = AppendAndAlphabetize(results[len(results)-1], b)
		}
	}

	// Bunch of duplicate code for the edge case of eltur tìtxen si and others like it
	if !bareNaviWord {
		found := false
		// See if it is in the list known to start multiword words
		if _, ok := multiword_words[searchNaviWord]; ok {
			// If so, loop through it
			for _, pairWordSet := range multiword_words[searchNaviWord] {
				keepAffixes := *new(affix)

				extraWord := 0
				// There could be more than one pair (win säpi and win si for example)
				for j, pairWord := range pairWordSet {
					found = false
					// Don't cause an index out of range error
					if i+j+1 >= len(allWords) {
						break
					} else {
						// For "[word] ke si and [word] rä'ä si"
						if (allWords[i+j+1] == "ke" || allWords[i+j+1] == "rä'ä") && IsVerb(allWords[i+j+2]) {
							extraWord = 1
							if len(results) == 1 {
								results = append(results, []Word{simpleWord(allWords[i+j+1])})
								for _, b := range dictHash[allWords[i+j+1]] {
									results[1] = AppendToFront(results[1], b)
								}
							}
							found = true
							j += 1
						}
						// Find all words the second word can represent
						secondWords := []Word{}

						// First by itself
						if pairWord == allWords[i+j+1] {
							found = true
							results[0][0].Navi += " " + allWords[i+j+1]
							continue
						}

						// And then by its possible conjugations
						for _, b := range TestDeconjugations(allWords[i+j+1]) {
							secondWords = AppendAndAlphabetize(secondWords, b)
						}

						// Do any of the conjugations work?
						for _, b := range secondWords {
							if b.Navi == pairWord {
								results[0][0].Navi += " " + allWords[i+j+1]
								found = true
								keepAffixes = addAffixes(keepAffixes, b.Affixes)
							}
						}

						// Chain is broken.  Exit.
						if !found {
							break
						}
					}
				}
				if found {
					fullWord := searchNaviWord
					for _, pairWord := range pairWordSet {
						fullWord += " " + pairWord
					}

					results[0] = []Word{results[0][0]}

					for _, definition := range dictHash[fullWord] {
						// Replace the word
						if len(results) > 0 && len(results[0]) > 1 && (results[0][1].Navi == "ke" || results[0][1].Navi == "rä'ä") {
							// Get the query it's looking for
							results[0][len(results[0])-1].Navi = results[0][1].Navi
							results[1] = AppendToFront(results[1], definition)
							results[1][1].Affixes = keepAffixes
						} else {
							// Get the query it's looking for
							results[0] = AppendToFront(results[0], definition)
							results[0][1].Affixes = keepAffixes
						}
					}
					i += len(pairWordSet) + extraWord
				}
			}
		}
	}

	if checkFixes {
		if len(results) > 0 && len(results[0]) > 0 {
			if !(strings.ToLower(results[len(results)-1][0].Navi) != searchNaviWord && strings.HasPrefix(strings.ToLower(results[len(results)-1][0].Navi), searchNaviWord)) {
				// Find all possible unconjugated versions of the word
				for _, a := range TestDeconjugations(searchNaviWord) {
					results[len(results)-1] = AppendAndAlphabetize(results[len(results)-1], a)
				}
			}
		} else {
			// Find all possible unconjugated versions of the word
			for _, a := range TestDeconjugations(searchNaviWord) {
				results[len(results)-1] = AppendAndAlphabetize(results[len(results)-1], a)
			}
		}

		// Check if the word could have more than one word
		found := false
		// Find the results words
		for _, a := range results[len(results)-1] {
			// See if it is in the list known to start multiword words
			if _, ok := multiword_words[a.Navi]; ok {
				// If so, loop through it
				for _, pairWordSet := range multiword_words[a.Navi] {
					keepAffixes := *new(affix)
					keepAffixes = addAffixes(keepAffixes, a.Affixes)

					extraWord := 0
					// There could be more than one pair (win säpi and win si for example)
					for j, pairWord := range pairWordSet {
						found = false
						// Don't cause an index out of range error
						if i+j+1 >= len(allWords) {
							break
						} else {
							if (allWords[i+j+1] == "ke" || allWords[i+j+1] == "rä'ä") && IsVerb(allWords[i+j+2]) {
								extraWord = 1
								if len(results) == 1 {
									results = append(results, []Word{simpleWord(allWords[i+j+1])})
									for _, b := range dictHash[allWords[i+j+1]] {
										results[1] = AppendToFront(results[1], b)
									}
								}

								j += 1
							}
							// Find all words the second word can represent
							secondWords := []Word{}

							// First by itself
							if pairWord == allWords[i+j+1] {
								found = true
								results[0][0].Navi += " " + allWords[i+j+1]
								continue
							}

							// And then by its possible conjugations
							for _, b := range TestDeconjugations(allWords[i+j+1]) {
								secondWords = AppendAndAlphabetize(secondWords, b)
							}

							// Do any of the conjugations work?
							for _, b := range secondWords {
								if b.Navi == pairWord {
									found = true
									results[0][0].Navi += " " + allWords[i+j+1]
									keepAffixes = addAffixes(keepAffixes, b.Affixes)
								}
							}

							// Chain is broken.  Exit.
							if !found {
								break
							}
						}
					}
					if found {
						fullWord := a.Navi
						for _, pairWord := range pairWordSet {
							fullWord += " " + pairWord
						}

						results[0] = []Word{results[0][0]}

						for _, definition := range dictHash[fullWord] {
							// Replace the word
							if len(results) > 0 && len(results[0]) > 1 && (results[0][1].Navi == "ke" || results[0][1].Navi == "rä'ä") {
								// Get the query it's looking for
								results[0][len(results[0])-1].Navi = results[0][1].Navi
								results[1] = AppendToFront(results[1], definition)
								results[1][1].Affixes = keepAffixes
							} else {
								// Get the query it's looking for
								results[0] = AppendToFront(results[0], definition)
								results[0][1].Affixes = keepAffixes
							}
						}
						i += len(pairWordSet) + extraWord
					}
				}
			}
		}
	}

	// If we found nothing, at least return the query
	if len(results[0]) == 0 {
		return i - start, [][]Word{{simpleWord(searchNaviWord)}}, nil
	}

	if len(results) == 2 {
		temp := results[0]
		results[0] = results[1]
		results[1] = temp
	}

	return i - start, results, nil
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
			results = AppendAndAlphabetize(results, c)
		}
	}

	return
}

func TranslateToNaviHash(searchWord string, langCode string) (results [][]Word) {
	searchWord = clean(searchWord)

	// No Results if empty string after removing sketch chars
	if len(searchWord) == 0 {
		return
	}

	results = [][]Word{}

	for _, word := range strings.Split(searchWord, " ") {
		results = append(results, []Word{})
		for _, a := range TranslateToNaviHashHelper(word, langCode) {
			results[len(results)-1] = AppendAndAlphabetize(results[len(results)-1], a)
		}
		// Append the query to the front of the list
		tempResults := []Word{simpleWord(word)}
		for _, b := range results[len(results)-1] {
			tempResults = append(tempResults, b)
		}
		results[len(results)-1] = tempResults
	}

	return
}

func TranslateToNaviHashHelper(searchWord string, langCode string) (results []Word) {
	results = []Word{}
	switch langCode {
	case "de":
		for _, a := range SearchNatlangWord(dictHash2.DE, searchWord) {
			results = AppendAndAlphabetize(results, a)
		}
	case "en":
		for _, a := range SearchNatlangWord(dictHash2.EN, searchWord) {
			results = AppendAndAlphabetize(results, a)
		}
	case "et":
		for _, a := range SearchNatlangWord(dictHash2.ET, searchWord) {
			results = AppendAndAlphabetize(results, a)
		}
	case "fr":
		for _, a := range SearchNatlangWord(dictHash2.FR, searchWord) {
			results = AppendAndAlphabetize(results, a)
		}
	case "hu":
		for _, a := range SearchNatlangWord(dictHash2.HU, searchWord) {
			results = AppendAndAlphabetize(results, a)
		}
	case "nl":
		for _, a := range SearchNatlangWord(dictHash2.NL, searchWord) {
			results = AppendAndAlphabetize(results, a)
		}
	case "pl":
		for _, a := range SearchNatlangWord(dictHash2.PL, searchWord) {
			results = AppendAndAlphabetize(results, a)
		}
	case "ru":
		for _, a := range SearchNatlangWord(dictHash2.RU, searchWord) {
			results = AppendAndAlphabetize(results, a)
		}
	case "sv":
		for _, a := range SearchNatlangWord(dictHash2.SV, searchWord) {
			results = AppendAndAlphabetize(results, a)
		}
	case "tr":
		for _, a := range SearchNatlangWord(dictHash2.TR, searchWord) {
			results = AppendAndAlphabetize(results, a)
		}
	default:
		// If we get an odd language code, return English
		for _, a := range SearchNatlangWord(dictHash2.EN, searchWord) {
			results = AppendAndAlphabetize(results, a)
		}
	}

	return
}

// Translate some text.  The language context is with Eywa now :ipu:
// !! Multiple words are supported !!
// This will return a 2D array of Words, that fit the input text
// One Word can have multiple meanings and words (e.g. synonyms)
func BidirectionalSearch(searchNaviWords string, checkFixes bool, langCode string) (results [][]Word, err error) {
	searchNaviWords = clean(searchNaviWords)

	// No Results if empty string after removing sketch chars
	if len(searchNaviWords) == 0 {
		return
	}

	allWords := strings.Split(searchNaviWords, " ")

	i := 0

	results = [][]Word{}
	for i < len(allWords) {

		// Search for Na'vi words
		j, newWords, error2 := TranslateFromNaviHashHelper(i, allWords, checkFixes)
		if error2 == nil {
			for _, newWord := range newWords {
				// Set up receptacle for words
				results = append(results, []Word{})
				for _, newResult := range newWord {
					results[len(results)-1] = append(results[len(results)-1], newResult)
				}
			}
		}

		// Search for natural language words
		natlangWords := []Word{}
		for _, a := range TranslateToNaviHashHelper(allWords[i], langCode) {
			// We want them alphabetized with their fellow natlang words...
			natlangWords = AppendAndAlphabetize(natlangWords, a)
		}

		for _, a := range natlangWords {
			// ...but not with the Na'vi words
			results[len(results)-1] = append(results[len(results)-1], a)
		}

		i += j

		i++
	}

	return
}

// Get random words out of the dictionary.
// If args are applied, the dict will be filtered for args before random words are chosen.
// args will be put into the `List()` algorithm.
func Random(amount int, args []string, checkDigraphs uint8) (results []Word, err error) {
	allWords, err := List(args, checkDigraphs)

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

// Get all words with spaces
func GetMultiwordWords() map[string][][]string {
	return multiword_words
}

// Get all words with multiple definitions
func GetHomonyms() (results [][]Word, err error) {
	return TranslateFromNaviHash(homonyms, false)
}

func EjectiveSoftener(ipa string, oldLetter string, newLetter string) (newIpa string) {
	ipa = "." + ipa

	for i, k := range []string{"t͡s", "s", "f"} {
		ipa = strings.ReplaceAll(ipa, k+oldLetter, string(i))
	}

	ipa = strings.ReplaceAll(ipa, "."+oldLetter, "."+newLetter)
	ipa = strings.ReplaceAll(ipa, ".ˈ"+oldLetter, ".ˈ"+newLetter)

	for i, k := range []string{"t͡s", "s", "f"} {
		ipa = strings.ReplaceAll(ipa, string(i), k+oldLetter)
	}

	ipa = strings.TrimPrefix(ipa, ".")

	return ipa
}

/* Is it a vowel? (for when the psuedovowel bool won't work) */
func is_vowel_ipa(letter string) (found bool) {
	// Also arranged from most to least common (not accounting for diphthongs)
	vowels := []string{"a", "ɛ", "u", "ɪ", "o", "i", "æ", "ʊ"}
	// Linear search
	for i := 0; i < 8; i++ {
		if letter == vowels[i] {
			return true
		}
	}
	return false
}

func ReefMe(ipa string, inter bool) []string {
	if ipa == "ʒɛjk'.ˈsu:.li" {
		return []string{"jake-sùl-ly", "ʒɛjk'.ˈsʊ:.li"}
	} else if ipa == "ˈz·ɛŋ.kɛ" {
		return []string{"zen-ke", "ˈz·ɛŋ.kɛ"}
	}

	breakdown := ""

	// Reefify the IPA first
	ipaReef := strings.ReplaceAll(ipa, "·", "")
	if !inter {
		ipaReef = EjectiveSoftener(ipaReef, "p'", "b")
		ipaReef = EjectiveSoftener(ipaReef, "t'", "d")
		ipaReef = EjectiveSoftener(ipaReef, "k'", "g")

		ipaReef = strings.ReplaceAll(ipaReef, "t͡sj", "tʃ")
		ipaReef = strings.ReplaceAll(ipaReef, "sj", "ʃ")

		temp := ""
		runes := []rune(ipaReef)

		for i, a := range runes {
			if i != 0 && i != len(runes)-1 && a == 'ʔ' {
				if runes[i-1] == '.' {
					if is_vowel_ipa(string(runes[i+1])) && is_vowel_ipa(string(runes[i-2])) {
						if runes[i+1] != runes[i-2] {
							continue
						}
					}
				} else if runes[i+1] == '.' {
					if is_vowel_ipa(string(runes[i+2])) && is_vowel_ipa(string(runes[i-1])) {
						if runes[i+2] != runes[i-1] {
							continue
						}
					}
				} else if runes[i-1] == 'ˈ' && i > 1 {
					if is_vowel_ipa(string(runes[i+1])) && is_vowel_ipa(string(runes[i-3])) {
						if runes[i+1] != runes[i-3] {
							continue
						}
					}
				}
			}
			temp += string(a)
		}

		ipaReef = temp
	}

	// now Romanize the reef IPA
	word := strings.Split(ipaReef, " ")

	breakdown = ""

	for j := 0; j < len(word); j++ {
		word[j] = strings.ReplaceAll(word[j], "]", "")
		// "or" means there's more than one IPA in this word, and we only want one
		if word[j] == "or" {
			break
		}

		syllables := strings.Split(word[j], ".")

		/* Onset */
		for k := 0; k < len(syllables); k++ {
			syllable := strings.ReplaceAll(syllables[k], "·", "")
			syllable = strings.ReplaceAll(syllable, "ˈ", "")
			syllable = strings.ReplaceAll(syllable, "ˌ", "")

			breakdown += "-"

			// tsy
			if strings.HasPrefix(syllable, "tʃ") {
				breakdown += "ch"
				syllable = strings.TrimPrefix(syllable, "tʃ")
			} else if len(syllable) >= 4 && syllable[0:4] == "t͡s" {
				// ts
				breakdown += "ts"
				//tsp
				if has("ptk", nth_rune(syllable, 3)) {
					if nth_rune(syllable, 4) == "'" {
						// ts + ejective onset
						breakdown += romanization2[syllable[4:6]]
						syllable = syllable[6:]
					} else {
						// ts + unvoiced plosive
						breakdown += romanization2[string(syllable[4])]
						syllable = syllable[5:]
					}
				} else if has("lɾmnŋwj", nth_rune(syllable, 3)) {
					// ts + other consonent
					breakdown += romanization2[nth_rune(syllable, 3)]
					syllable = syllable[4+len(nth_rune(syllable, 3)):]
				} else {
					// ts without a cluster
					syllable = syllable[4:]
				}
			} else if has("fs", nth_rune(syllable, 0)) {
				//
				breakdown += nth_rune(syllable, 0)
				if has("ptk", nth_rune(syllable, 1)) {
					if nth_rune(syllable, 2) == "'" {
						// f/s + ejective onset
						breakdown += romanization2[syllable[1:3]]
						syllable = syllable[3:]
					} else {
						// f/s + unvoiced plosive
						breakdown += romanization2[string(syllable[1])]
						syllable = syllable[2:]
					}
				} else if has("lɾmnŋwj", nth_rune(syllable, 1)) {
					// f/s + other consonent
					breakdown += romanization2[nth_rune(syllable, 1)]
					syllable = syllable[1+len(nth_rune(syllable, 1)):]
				} else {
					// f/s without a cluster
					syllable = syllable[1:]
				}
			} else if has("ptk", nth_rune(syllable, 0)) {
				if nth_rune(syllable, 1) == "'" {
					// ejective
					breakdown += romanization2[syllable[0:2]]
					syllable = syllable[2:]
				} else {
					// unvoiced plosive
					breakdown += romanization2[string(syllable[0])]
					syllable = syllable[1:]
				}
			} else if has("ʔlɾhmnŋvwjzbdg", nth_rune(syllable, 0)) {
				// other normal onset
				breakdown += romanization2[nth_rune(syllable, 0)]
				syllable = syllable[len(nth_rune(syllable, 0)):]
			} else if has("ʃʒ", nth_rune(syllable, 0)) {
				// one sound representd as a cluster
				if nth_rune(syllable, 0) == "ʃ" {
					breakdown += "sh"
				}
				syllable = syllable[len(nth_rune(syllable, 0)):]
			}

			/*
			 * Nucleus
			 */
			if len(syllable) > 1 && has("jw", nth_rune(syllable, 1)) {
				//diphthong
				breakdown += romanization2[syllable[0:len(nth_rune(syllable, 0))+1]]
				syllable = string([]rune(syllable)[2:])
			} else if len(syllable) > 1 && has("lr", nth_rune(syllable, 0)) {
				breakdown += romanization2[syllable[0:3]]
				continue
			} else {
				//vowel
				breakdown += romanization2[nth_rune(syllable, 0)]
				syllable = string([]rune(syllable)[1:])
			}

			/*
			 * Coda
			 */
			if len(syllable) > 0 {
				if nth_rune(syllable, 0) == "s" {
					breakdown += "sss" //oìsss only
				} else {
					if syllable == "k̚" {
						breakdown += "k"
					} else if syllable == "p̚" {
						breakdown += "p"
					} else if syllable == "t̚" {
						breakdown += "t"
					} else if syllable == "ʔ̚" {
						breakdown += "'"
					} else {
						if syllable[0] == 'k' && len(syllable) > 1 {
							breakdown += "kx"
						} else {
							breakdown += romanization2[syllable]
						}
					}
				}
			}
		}
		breakdown += " "
	}

	breakdown = strings.TrimPrefix(breakdown, "-")
	breakdown = strings.ReplaceAll(breakdown, " -", " ")
	breakdown = strings.TrimSuffix(breakdown, " ")

	return []string{breakdown, ipaReef}
}

func StartEverything() {
	AssureDict()
	CacheDictHash()
	CacheDictHash2()
	PhonemeDistros()
}
