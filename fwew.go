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
func TranslateFromNaviHash(searchNaviWords string, checkFixes bool) (results [][]Word, err error) {
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

	// No Results if empty string after removing sketch chars
	if len(searchNaviWords) == 0 {
		return
	}

	allWords := strings.Split(searchNaviWords, " ")

	i := 0

	for i < len(allWords) {
		searchNaviWord := allWords[i]
		// Make sure we have a place to put the found words
		results = append(results, []Word{})

		bareNaviWord := false
		// Find the word
		if _, ok := dictHash[searchNaviWord]; ok {
			bareNaviWord = true
			for _, b := range dictHash[searchNaviWord] {
				results[len(results)-1] = append(results[len(results)-1], b)
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
					// There could be more than one pair (win säpi and win si for example)
					for j, pairWord := range pairWordSet {
						found = false
						// Don't cause an index out of range error
						if i+j+1 >= len(allWords) {
							break
						} else {
							// Find all words the second word can represent
							secondWords := []Word{}

							// First by itself
							if pairWord == allWords[i+j+1] {
								found = true
								continue
							}

							// And then by its possible conjugations
							for _, b := range TestDeconjugations(allWords[i+j+1]) {
								secondWords = append(secondWords, b)
							}

							// Do any of the conjugations work?
							for _, b := range secondWords {
								if b.Navi == pairWord {
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
						for _, definition := range dictHash[fullWord] {
							// Replace the word
							results[len(results)-1] = []Word{definition}
							results[len(results)-1][0].Affixes = keepAffixes

							//fmt.Println(results[len(results)-1][0].Affixes)
							i += len(pairWordSet) + 1
						}
					}
				}
			}
			// We found it!
			if found {
				continue
			}
		}

		if checkFixes {
			// Find all possible unconjugated versions of the word
			for _, a := range TestDeconjugations(searchNaviWord) {
				results[len(results)-1] = append(results[len(results)-1], a)
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
						// There could be more than one pair (win säpi and win si for example)
						for j, pairWord := range pairWordSet {
							found = false
							// Don't cause an index out of range error
							if i+j+1 >= len(allWords) {
								break
							} else {
								// Find all words the second word can represent
								secondWords := []Word{}

								// First by itself
								if pairWord == allWords[i+j+1] {
									found = true
									continue
								}

								// And then by its possible conjugations
								for _, b := range TestDeconjugations(allWords[i+j+1]) {
									secondWords = append(secondWords, b)
								}

								// Do any of the conjugations work?
								for _, b := range secondWords {
									if b.Navi == pairWord {
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
							fullWord := a.Navi
							for _, pairWord := range pairWordSet {
								fullWord += " " + pairWord
							}
							for _, definition := range dictHash[fullWord] {
								// Replace the word
								keepAffixes = addAffixes(keepAffixes, a.Affixes)

								results[len(results)-1] = []Word{definition}
								results[len(results)-1][0].Affixes = keepAffixes

								//fmt.Println(results[len(results)-1][0].Affixes)
								i += len(pairWordSet)
							}
						}
					}
				}
			}
		}
		i++
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

func TranslateToNaviHash(searchWord string, langCode string) (results [][]Word) {
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

	results = [][]Word{}

	for _, word := range strings.Split(searchWord, " ") {
		results = append(results, []Word{})
		switch langCode {
		case "de":
			for _, a := range SearchNatlangWord(dictHash2.DE, word) {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		case "en":
			for _, a := range SearchNatlangWord(dictHash2.EN, word) {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		case "et":
			for _, a := range SearchNatlangWord(dictHash2.ET, word) {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		case "fr":
			for _, a := range SearchNatlangWord(dictHash2.FR, word) {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		case "hu":
			for _, a := range SearchNatlangWord(dictHash2.HU, word) {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		case "nl":
			for _, a := range SearchNatlangWord(dictHash2.NL, word) {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		case "pl":
			for _, a := range SearchNatlangWord(dictHash2.PL, word) {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		case "ru":
			for _, a := range SearchNatlangWord(dictHash2.RU, word) {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		case "sv":
			for _, a := range SearchNatlangWord(dictHash2.SV, word) {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		case "tr":
			for _, a := range SearchNatlangWord(dictHash2.TR, word) {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		default:
			// If we get an odd language code, return English
			for _, a := range SearchNatlangWord(dictHash2.EN, searchWord) {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		}
	}

	return
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

// Get all words with spaces
func GetMultiwordWords() map[string][][]string {
	return multiword_words
}

func StartEverything() {
	AssureDict()
	CacheDictHash()
	CacheDictHash2()
	PhonemeDistros()
}
