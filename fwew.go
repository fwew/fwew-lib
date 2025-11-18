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
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
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
	"æ": "ä", "õ": "õ", //võvä' only
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
	badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\"„“”«»`

	// remove all the sketchy chars from arguments
	for _, c := range badChars {
		searchNaviWords = strings.ReplaceAll(searchNaviWords, string(c), " ")
	}

	// Recognize line breaks and turn them into spaces
	searchNaviWords = strings.ReplaceAll(searchNaviWords, "\n", " ")

	// No leading or trailing spaces
	searchNaviWords = strings.TrimSpace(searchNaviWords)

	// normalize tìftang character
	searchNaviWords = strings.ReplaceAll(searchNaviWords, "’", "'")
	searchNaviWords = strings.ReplaceAll(searchNaviWords, "‘", "'")

	// find everything lowercase
	searchNaviWords = strings.ToLower(searchNaviWords)

	// Get rid of all double spaces
	for strings.Contains(searchNaviWords, "  ") {
		searchNaviWords = strings.ReplaceAll(searchNaviWords, "  ", " ")
	}

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
func TranslateFromNaviHash(searchNaviWords string, checkFixes bool, strict bool, allowReef bool) (results [][]Word, err error) {
	universalLock.Lock()
	defer universalLock.Unlock()
	searchNaviWords = clean(searchNaviWords)

	// No Results if empty string after removing sketch chars
	if len(searchNaviWords) == 0 {
		return
	}

	allWords := strings.Split(clean(searchNaviWords), " ")

	i := 0

	results = [][]Word{}

	dict := &dictHashLoose

	if !allowReef {
		dict = &dictHashStrict
	} else if strict {
		dict = &dictHashStrictReef
	}

	for i < len(allWords) {
		// Skip empty words or ridiculously long words
		// 50 was chosen because a quick and dirty program found the max
		// Na'vi word length is 43 (before adding sì to the end)
		if len(allWords[i]) == 0 || len([]rune(allWords[i])) > 50 {
			i++
			continue
		}
		j, newWords, error2 := TranslateFromNaviHashHelper(dict, i, allWords, checkFixes, strict, allowReef)
		if error2 == nil {
			for _, newWord := range newWords {
				// Set up receptacle for words
				results = append(results, []Word{})
				results[len(results)-1] = append(results[len(results)-1], newWord...)
			}
		}

		if len(results[len(results)-1]) > 1 && len(strings.Split(results[len(results)-1][1].Navi, " ")) > 1 {
			newQuery := ""
			kOffset := 0
			for k := range strings.Split(results[len(results)-1][1].Navi, " ") {
				if i+k+kOffset >= len(allWords) {
					break
				}
				if allWords[i+k+kOffset] == "ke" || strings.ReplaceAll(allWords[i+k+kOffset], "e", "ä") == "rä'ä" {
					kOffset += 1
				}
				if k != 0 {
					newQuery += " "
				}
				newQuery += allWords[i+k+kOffset]
				if strings.HasSuffix(allWords[i+k+kOffset], "-susi") {
					break
				}
			}
			results[len(results)-1][0].Navi = newQuery
		}
		i += j
		i++
	}

	return
}

// Helper for TranslateFromNaviHashHelper
func AppendToFront(words []Word, input Word) []Word {
	// Ensure it's not a duplicate
	for i, a := range words {
		if i != 0 && input.ID == a.ID {
			if len(input.Affixes.Prefix) == len(a.Affixes.Prefix) && len(input.Affixes.Suffix) == len(a.Affixes.Suffix) &&
				len(input.Affixes.Lenition) == len(a.Affixes.Lenition) && len(input.Affixes.Infix) == len(a.Affixes.Infix) {
				return words
			}
		}
	}
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
func IsVerb(dict *map[string][]Word, input string, comparator string, strict bool, allowReef bool) (result bool, affixes Word) {
	affixes = simpleWord(input)
	_, possibilities, err := TranslateFromNaviHashHelper(dict, 0, []string{input}, true, strict, allowReef)
	_, possibilities2, err2 := TranslateFromNaviHashHelper(dict, 0, []string{comparator}, true, strict, allowReef)
	if err != nil || err2 != nil {
		return false, affixes
	}
	isRealVerb := false
	pairFound := false
	unknownInfix := false
	for _, a := range possibilities {
		for i, b := range a {
			// Don't check the empty first row
			if i == 0 {
				continue
			}

			for _, prefix := range verbPrefixes {
				for _, ourPrefixes := range b.Affixes.Prefix {
					if prefix == ourPrefixes {
						return false, affixes
					}
				}
			}

			for _, suffix := range verbSuffixes {
				for _, ourSuffixes := range b.Affixes.Suffix {
					if suffix == ourSuffixes {
						return false, affixes
					}
				}
			}

			// Make sure it's a verb
			if len(b.PartOfSpeech) > 0 && b.PartOfSpeech[0] == 'v' {
				for _, c := range b.Affixes.Infix {
					// <us> and <awn> are participles, so they become adjectives
					if c == "us" || c == "awn" {
						return false, affixes
					}
				}
				isRealVerb = true
			}
			// Make sure it's also found in the multiword word set
			for _, c := range possibilities2 {
				for j, d := range c {
					// Don't check the empty first row
					if j == 0 {
						continue
					}
					if d.ID == b.ID {
						affixes = b
						// Infix check is to make sure "win säpi" doesn't become "win si"
						// Make sure d doesn't have an infix that b has
						pairFound = true
						miniMap := map[string]bool{}
						for _, e := range b.Affixes.Infix {
							miniMap[e] = true
						}
						for _, f := range d.Affixes.Infix {
							if _, ok := miniMap[f]; !ok {
								unknownInfix = true
								break
							}
						}
					}
				}
			}
		}
	}
	return (isRealVerb && pairFound && !unknownInfix), affixes
}

func TranslateFromNaviHashHelper(dict *map[string][]Word, start int, allWords []string, checkFixes bool, strict bool, allowReef bool) (steps int, results [][]Word, err error) {
	i := start

	containsUmlaut := []bool{}
	containsTìftang := []bool{}

	tempResults := []Word{}
	searchNaviWord := ""

	// don't crunch more than once
	if !strict || allowReef {
		for _, a := range allWords {
			if strings.Contains(a, "ä") {
				containsUmlaut = append(containsUmlaut, true)
			} else {
				containsUmlaut = append(containsUmlaut, false)
			}

			strippedA := strings.TrimPrefix(strings.TrimSuffix(a, "'"), "'")
			if strings.Contains(strippedA, "'") {
				containsTìftang = append(containsTìftang, true)
			} else {
				containsTìftang = append(containsTìftang, false)
			}
		}

		results = [][]Word{{simpleWord(allWords[i])}}

		allWords = dialectCrunch(allWords, false, strict, allowReef)

		searchNaviWord = allWords[i]

		// Find the word
		a := strings.ReplaceAll(searchNaviWord, "ù", "u")

		if _, ok := (*dict)[a]; ok {
			//bareNaviWord = true
			for _, b := range (*dict)[a] {
				results[len(results)-1] = AppendAndAlphabetize(results[len(results)-1], b)
			}
		}

		// If one searches kivä, make sure kive doesn't show up
		for _, a := range results[len(results)-1] {
			if containsUmlaut[i] && !strings.Contains(strings.ToLower(a.Navi), "ä") {
				continue // ä can unstress to e, but not the other way around
			}
			strippedA := a.Navi
			if len(a.Affixes.Prefix) == 0 {
				strippedA = strings.TrimPrefix(strippedA, "'")
			}
			if len(a.Affixes.Suffix) == 0 {
				strippedA = strings.TrimSuffix(strippedA, "'")
			}
			if containsTìftang[i] && !strings.Contains(strippedA, "'") {
				continue // make sure tsa'u doesn't return tsa-au
			}
			tempResults = append(tempResults, a)
		}
	} else {
		for range len(allWords) {
			containsUmlaut = append(containsUmlaut, true)
			containsTìftang = append(containsTìftang, false)
		}

		searchNaviWord = allWords[i]
		results = [][]Word{{simpleWord(allWords[i])}}

		// Find the word
		a := strings.ReplaceAll(searchNaviWord, "ù", "u")

		if _, ok := (*dict)[a]; ok {
			for _, b := range (*dict)[a] {
				results[len(results)-1] = AppendAndAlphabetize(results[len(results)-1], b)
			}
		} else if allowReef {
			noUmlaut := strings.ReplaceAll(a, "ä", "e")
			if _, ok := (*dict)[noUmlaut]; ok {
				for _, b := range (*dict)[noUmlaut] {
					results[len(results)-1] = AppendAndAlphabetize(results[len(results)-1], b)
				}
			}
		}

		tempResults = append(tempResults, results[len(results)-1]...)
	}

	results[len(results)-1] = tempResults

	foundAlready := false

	// Bunch of duplicate code for the edge case of eltur tìtxen si and others like it
	//if !bareNaviWord {
	found := false
	// See if it is in the list known to start multiword words
	multiwords := &multiword_words
	if !strict {
		multiwords = &multiword_words_loose
	} else if allowReef {
		multiwords = &multiword_words_reef
	}
	if _, ok := (*multiwords)[searchNaviWord]; ok {
		// If so, loop through it
		for _, pairWordSet := range (*multiwords)[searchNaviWord] {
			if foundAlready {
				break
			}

			keepAffixes := *new(affix)

			extraWord := 0

			revert := results[0][0].Navi
			// There could be more than one pair (win säpi and win si for example)
			for j, pairWord := range pairWordSet {
				found = false
				// Don't cause an index out of range error
				if i+j+1 >= len(allWords) {
					break
				} else {
					// For "[word] ke si and [word] rä'ä si"
					if i+j+2 < len(allWords) && (allWords[i+j+1] == "ke" || allWords[i+j+1] == "rä'ä") {
						validVerb, itsAffixes := IsVerb(dict, allWords[i+j+2], pairWord, strict, allowReef)
						if validVerb {
							extraWord = 1
							if len(results) == 1 {
								results = append(results, []Word{simpleWord(allWords[i+j+1])})
								for _, b := range (*dict)[allWords[i+j+1]] {
									results[1] = AppendToFront(results[1], b)
								}
							}
							found = true
							foundAlready = true
							revert += " " + allWords[i+j+2]
							keepAffixes = itsAffixes.Affixes
							j += 1
							continue
						}
					}

					// Verbs don't just come after ke or rä'ä
					validVerb, itsAffixes := IsVerb(dict, allWords[i+j+1], pairWord, strict, allowReef)
					if validVerb {
						found = true
						foundAlready = true
						revert += " " + allWords[i+j+1]
						keepAffixes = itsAffixes.Affixes
						continue
					}

					// Find all words the second word can represent
					secondWords := []Word{}

					// First by itself
					if pairWord == allWords[i+j+1] {
						found = true
						revert += " " + allWords[i+j+1]
						continue
					}

					// And then by its possible conjugations
					for _, b := range TestDeconjugations(dict, allWords[i+j+1], strict, allowReef, containsUmlaut[i]) {
						breakAdding := false
						for _, prefix := range verbPrefixes {
							for _, ourPrefixes := range b.Affixes.Prefix {
								if prefix == ourPrefixes {
									breakAdding = true
								}
							}
							if breakAdding {
								break
							}
						}

						if !breakAdding {
							for _, suffix := range verbSuffixes {
								for _, ourSuffixes := range b.Affixes.Suffix {
									if suffix == ourSuffixes {
										breakAdding = true
									}
								}
								if breakAdding {
									break
								}
							}
						}

						if breakAdding {
							continue
						}

						secondWords = AppendAndAlphabetize(secondWords, b)
					}

					// Do any of the conjugations work?
					for _, b := range secondWords {

						if b.Navi == pairWord {
							revert += " " + b.Navi
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
				results[0][0].Navi = revert
				fullWord := searchNaviWord
				for _, pairWord := range pairWordSet {
					fullWord += " " + pairWord
				}

				results[0] = []Word{results[0][0]}
				a := strings.ReplaceAll(fullWord, "ù", "u")

				for _, definition := range (*dict)[a] {
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
	//}

	if checkFixes {
		newResults := []Word{}

		if !foundAlready {
			if len(results) > 0 && len(results[0]) > 0 {
				if !(strings.ToLower(results[len(results)-1][0].Navi) != searchNaviWord && strings.HasPrefix(strings.ToLower(results[len(results)-1][0].Navi), searchNaviWord)) {
					// Find all possible unconjugated versions of the word
					newResults = TestDeconjugations(dict, searchNaviWord, strict, allowReef, containsUmlaut[i])
				}
			} else {
				// Find all possible unconjugated versions of the word
				newResults = TestDeconjugations(dict, searchNaviWord, strict, allowReef, containsUmlaut[i])
			}
		}

		tempNewResults := []Word{}

		// If one searches for ke, don't search for kä
		for _, a := range newResults {
			nucleusCount := 0
			for _, b := range []string{"a", "ä", "e", "i", "ì", "o", "u", "ù", "ll", "rr"} {
				nucleusCount += strings.Count(a.Navi, b)
			}
			if nucleusCount == 1 {
				if !containsUmlaut[i] && !strings.Contains(searchNaviWord, "a") && strings.Contains(a.Navi, "ä") {
					continue
				}
			}
			strippedA := a.Navi
			if len(a.Affixes.Prefix) == 0 {
				strippedA = strings.TrimPrefix(strippedA, "'")
			}
			if len(a.Affixes.Suffix) == 0 {
				strippedA = strings.TrimSuffix(strippedA, "'")
			}
			if containsTìftang[i] && !strings.Contains(strippedA, "'") {
				continue // make sure tsa'u doesn't return tsa-au
			}
			tempNewResults = append(tempNewResults, a)
		}

		// Do not duplicate
		alreadyHere := results[len(results)-1]
		for _, a := range tempNewResults {
			add := true
			for _, b := range alreadyHere {
				if b.ID == a.ID {
					if len(b.Affixes.Prefix) == len(a.Affixes.Prefix) &&
						len(b.Affixes.Suffix) == len(a.Affixes.Suffix) &&
						len(b.Affixes.Lenition) == len(a.Affixes.Lenition) &&
						len(b.Affixes.Infix) == len(a.Affixes.Infix) {
						add = false
						break
					}
				}
			}
			if add {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		}

		// Check if the word could have more than one word
		found := false
		// Find the results words

		revert := results[0][0].Navi

		for _, a := range results[len(results)-1] {
			breakAdding2 := false

			for _, prefix := range verbPrefixes {
				for _, ourPrefixes := range a.Affixes.Prefix {
					if prefix == ourPrefixes {
						breakAdding2 = true
					}
				}
				if breakAdding2 {
					break
				}
			}

			if !breakAdding2 {
				for _, suffix := range verbSuffixes {
					for _, ourSuffixes := range a.Affixes.Suffix {
						if suffix == ourSuffixes {
							breakAdding2 = true
						}
					}
					if breakAdding2 {
						break
					}
				}
			}

			// No tìngyu mikyun
			if !breakAdding2 {
				// See if it is in the list known to start multiword words
				if _, ok := (*multiwords)[a.Navi]; ok {
					// If so, loop through it
					for _, pairWordSet := range (*multiwords)[a.Navi] {
						if foundAlready {
							break
						}

						newSearch := a.Navi

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
								// For "[word] ke si and [word] rä'ä si"
								if i+j+2 < len(allWords) && (allWords[i+j+1] == "ke" || allWords[i+j+1] == "ree") {
									validVerb, itsAffixes := IsVerb(dict, allWords[i+j+2], pairWord, strict, allowReef)
									if validVerb {
										extraWord = 1
										if len(results) == 1 {
											results = append(results, []Word{simpleWord(allWords[i+j+1])})
											for _, b := range (*dict)[allWords[i+j+1]] {
												results[1] = AppendToFront(results[1], b)
											}
										}
										found = true
										foundAlready = true
										revert += " " + allWords[i+j+2]
										keepAffixes = itsAffixes.Affixes
										j += 1

										continue
									}
								}

								// Find all words the second word can represent
								secondWords := []Word{}

								allWord := allWords[i+j+1]

								if !strict || allowReef {
									pairWord = dialectCrunch([]string{pairWord}, false, strict, allowReef)[0]
									allWord = dialectCrunch([]string{allWord}, false, strict, allowReef)[0]
								}

								// First by itself
								if pairWord == allWord {
									found = true
									revert += " " + allWords[i+j+1]
									continue
								}

								// And then by its possible conjugations
								for _, b := range TestDeconjugations(dict, allWords[i+j+1], strict, allowReef, containsUmlaut[i]) {
									breakAdding := false
									for _, prefix := range verbPrefixes {
										for _, ourPrefixes := range b.Affixes.Prefix {
											if prefix == ourPrefixes {
												breakAdding = true
											}
										}
										if breakAdding {
											break
										}
									}

									if !breakAdding {
										for _, suffix := range verbSuffixes {
											for _, ourSuffixes := range b.Affixes.Suffix {
												if suffix == ourSuffixes {
													breakAdding = true
												}
											}
											if breakAdding {
												break
											}
										}
									}

									if breakAdding {
										continue
									}
									secondWords = AppendAndAlphabetize(secondWords, b)
								}

								// Do any of the conjugations work?
								for _, b := range secondWords {
									if b.Navi == pairWord {
										revert += " " + b.Navi
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
							results[0][0].Navi = revert
							fullWord := newSearch
							for _, pairWord := range pairWordSet {
								fullWord += " " + pairWord
							}

							results[0] = []Word{results[0][0]}
							a := strings.ReplaceAll(fullWord, "ù", "u")
							if !strict {
								a = dialectCrunch([]string{a}, false, strict, allowReef)[0]
							}

							for _, definition := range (*dict)[a] {
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
		for _, c := range dictHashStrict[firstResults[i]] {
			results = AppendAndAlphabetize(results, c)
		}
	}

	return
}

func TranslateToNaviHash(searchWord string, langCode string) (results [][]Word) {
	universalLock.Lock()
	defer universalLock.Unlock()
	searchWord = clean(searchWord)

	results = [][]Word{}

	for _, word := range strings.Split(searchWord, " ") {
		// Skip empty words
		if len(word) == 0 {
			continue
		}
		results = append(results, []Word{})
		for _, a := range TranslateToNaviHashHelper(&dictHash2Parenthesis, word, langCode) {
			results[len(results)-1] = AppendAndAlphabetize(results[len(results)-1], a)
		}
		// Append the query to the front of the list
		tempResults := []Word{simpleWord(word)}
		tempResults = append(tempResults, results[len(results)-1]...)
		results[len(results)-1] = tempResults
	}
	return
}

func TranslateToNaviHashHelper(dictionary *MetaDict, searchWord string, langCode string) (results []Word) {
	results = []Word{}
	switch langCode {
	case "de": // German
		for _, a := range SearchNatlangWord((*dictionary).DE, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.DE, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "en": // English
		for _, a := range SearchNatlangWord((*dictionary).EN, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.EN, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "es": // Spanish
		for _, a := range SearchNatlangWord((*dictionary).ES, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.ES, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "et": // Estonian
		for _, a := range SearchNatlangWord((*dictionary).ET, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.ET, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "fr": // French
		for _, a := range SearchNatlangWord((*dictionary).FR, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.FR, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "hu": // Hungarian
		for _, a := range SearchNatlangWord((*dictionary).HU, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.HU, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "ko": // Korean
		for _, a := range SearchNatlangWord((*dictionary).KO, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.KO, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "nl": // Dutch
		for _, a := range SearchNatlangWord((*dictionary).NL, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.NL, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "pl": // Polish
		for _, a := range SearchNatlangWord((*dictionary).PL, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.PL, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "pt": // Portuguese
		for _, a := range SearchNatlangWord((*dictionary).PT, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.PT, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "ru": // Russian
		for _, a := range SearchNatlangWord((*dictionary).RU, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.RU, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "sv": // Swedish
		for _, a := range SearchNatlangWord((*dictionary).SV, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.SV, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "tr": // Turkish
		for _, a := range SearchNatlangWord((*dictionary).TR, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.TR, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	case "uk": // Ukrainian
		for _, a := range SearchNatlangWord((*dictionary).UK, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.UK, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	default:
		// If we get an odd language code, return English
		for _, a := range SearchNatlangWord((*dictionary).EN, searchWord) {
			// Verify the search query is actually in the definition
			searchWords := SearchTerms(a.EN, false)
			found := false
			for _, d := range searchWords {
				if d == searchWord {
					found = true
					break
				}
			}
			if found {
				results = AppendAndAlphabetize(results, a)
			}
		}
	}

	return
}

// Translate some text.  The language context is with Eywa now :ipu:
// !! Multiple words are supported !!
// This will return a 2D array of Words, that fit the input text
// One Word can have multiple meanings and words (e.g. synonyms)
func BidirectionalSearch(searchNaviWords string, checkFixes bool, langCode string, allowReef bool) (results [][]Word, err error) {
	universalLock.Lock()
	defer universalLock.Unlock()
	searchNaviWords = clean(searchNaviWords)

	// No Results if empty string after removing sketch chars
	if len(searchNaviWords) == 0 {
		return
	}

	allWords := strings.Split(searchNaviWords, " ")

	i := 0

	ourDict := &dictHashLoose
	if !allowReef {
		ourDict = &dictHashStrict
	}

	results = [][]Word{}
	for i < len(allWords) {
		// Search for Na'vi words
		j, newWords, error2 := TranslateFromNaviHashHelper(ourDict, i, allWords, checkFixes, false, allowReef)

		NaviIDs := []string{}
		if error2 == nil {
			for _, newWord := range newWords {
				// Set up receptacle for words
				results = append(results, []Word{})
				results[len(results)-1] = append(results[len(results)-1], newWord...)
				if len(newWord) > 1 {
					NaviIDs = append(NaviIDs, newWord[1].ID)
				}
			}
		}

		// Search for natural language words
		natlangWords := []Word{}
		for _, a := range TranslateToNaviHashHelper(&dictHash2, allWords[i], langCode) {
			// Do not duplicate if the Na'vi word is in the definition
			if implContainsAny(NaviIDs, []string{a.ID}) {
				continue
			}
			// We want them alphabetized with their fellow natlang words...
			natlangWords = AppendAndAlphabetize(natlangWords, a)
		}

		// ...but not with the Na'vi words
		results[len(results)-1] = append(results[len(results)-1], natlangWords...)

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
	universalLock.Lock()
	defer universalLock.Unlock()
	return multiword_words
}

// Get all words with multiple definitions
func GetHomonyms() (results [][]Word, err error) {
	return TranslateFromNaviHash(homonyms, false, false, false)
}

// Get all words with non-standard phonotactics
func GetOddballs() (results [][]Word, err error) {
	return TranslateFromNaviHash(oddballs, true, false, false)
}

// Get all words with multiple definitions
func GetMultiIPA() (results [][]Word, err error) {
	return TranslateFromNaviHash(multiIPA, false, false, false)
}

/* Is it a vowel? (for when the psuedovowel bool won't work) */
func is_vowel_ipa(letter string) (found bool) {
	// Also arranged from most to least common (not accounting for diphthongs)
	vowels := []string{"a", "ɛ", "ɪ", "o", "u", "i", "æ", "ʊ"}
	// Linear search
	for _, a := range vowels {
		if letter == a {
			return true
		}
	}
	return false
}

func dialectCrunch(query []string, guaranteedForest bool, strict bool, allowReef bool) []string {
	newQuery := []string{}
	for _, a := range query {
		oldQuery := a

		// When caching, we are guaranteed forest words and don't need anything in this block
		if !guaranteedForest && allowReef {
			for i, b := range nkx {
				// make sure words like tìkankxan show up
				a = strings.ReplaceAll(a, strconv.Itoa(i), "")
				a = strings.ReplaceAll(a, b, strconv.Itoa(i))
			}
			// don't accidentally make every ng into nkx
			a = strings.ReplaceAll(a, "?", "")
			a = strings.ReplaceAll(a, "ng", "?")
			// unsoften ejectives
			a = strings.ReplaceAll(a, "b", "px")
			a = strings.ReplaceAll(a, "d", "tx")
			a = strings.ReplaceAll(a, "g", "kx")
			// these too
			a = strings.ReplaceAll(a, "ch", "tsy")
			a = strings.ReplaceAll(a, "sh", "sy")
			a = strings.ReplaceAll(a, "?", "ng")
			for i, b := range nkx {
				// make sure words like tìkankxan show up
				a = strings.ReplaceAll(a, strconv.Itoa(i), nkxSub[b])
			}
		}

		if allowReef {
			nucleusCount := 0
			// remove reef tìftangs
			for i, b := range []string{"a", "ä", "e", "i", "ì", "o", "u", "ù", "ll", "rr"} {
				if strings.Contains(a, b) {
					nucleusCount += strings.Count(a, b)
					for j, c := range []string{"a", "ä", "e", "i", "ì", "o", "u", "ù", "ll", "rr"} {
						if i < 8 && j < 8 {
							a = strings.ReplaceAll(a, b+"'"+c, b+c)
						}
					}
				}
			}
			if nucleusCount > 1 && strings.Contains(a, "ä") {
				// and to make sure every ä is possibly an e
				a = strings.ReplaceAll(a, "ä", "e")
			}

			// "eo" and "äo" are different, so the distinction must remain
			if strings.HasSuffix(oldQuery, "äo") || strings.HasSuffix(oldQuery, "ä'o") {
				a = strings.TrimSuffix(a, "eo") + "äo"
			}
		}

		newQuery = append(newQuery, a)
	}
	return newQuery
}

func ReefMe(ipa string, inter bool) []string {
	if ipa == "ʒɛjk'.ˈsu:.li" { // Obsolete path
		return []string{"jake-__sùl__-ly", "ʒɛjk'.ˈsʊ:.li"}
	} else if strings.ReplaceAll(ipa, "·", "") == "ˈzɛŋ.kɛ" { // only IPA not to match the Romanization
		return []string{"__zen__-ke", "ˈz·ɛŋ·.kɛ"}
	} else if ipa == "ɾæ.ˈʔæ" || ipa == "ˈɾæ.ʔæ" { // we hear this in Avatar 2
		return []string{"rä-__'ä__ or rä-__ä__", "ɾæ.ˈʔæ] or [ɾæ.ˈæ"}
	}

	// Replace the spaces so as not to confuse strings.Split()
	ipa = strings.ReplaceAll(ipa, " ", "*.")

	// Unstressed ä becomes e
	ipa_syllables := strings.Split(ipa, ".")
	if len(ipa_syllables) > 1 {
		new_ipa := ""
		for _, a := range ipa_syllables {
			new_ipa += "."
			if !strings.Contains(a, "ˈ") {
				new_ipa += strings.ReplaceAll(a, "æ", "ɛ")
			} else {
				new_ipa += a
			}
		}

		ipa = new_ipa
	}

	breakdown := ""
	ejectives := []string{"p'", "t'", "k'"}
	soften := map[string]string{
		"p'": "b",
		"t'": "d",
		"k'": "g",
	}

	// Reefify the IPA first
	ipaReef := strings.ReplaceAll(ipa, "·", "")
	if !inter {
		// atxkxe and ekxtxu
		for _, a := range ejectives {
			for _, b := range ejectives {
				ipaReef = strings.ReplaceAll(ipaReef, a+".ˈ"+b, soften[a]+".ˈ"+soften[b])
				ipaReef = strings.ReplaceAll(ipaReef, a+"."+b, soften[a]+"."+soften[b])
			}
		}

		// Ejectives before vowels and diphthongs become voiced plosives regardless of syllable boundaries
		for _, a := range ejectives {
			if strings.HasPrefix(ipaReef, a) {
				ipaReef = soften[a] + strings.TrimPrefix(ipaReef, a)
			}
			ipaReef = strings.ReplaceAll(ipaReef, ".ˈ"+a, ".ˈ"+soften[a])
			ipaReef = strings.ReplaceAll(ipaReef, "."+a, "."+soften[a])

			for _, b := range []string{"a", "ɛ", "ɪ", "o", "u", "i", "æ", "ʊ"} {
				ipaReef = strings.ReplaceAll(ipaReef, a+".ˈ"+b, soften[a]+".ˈ"+b)
				ipaReef = strings.ReplaceAll(ipaReef, a+"."+b, soften[a]+"."+b)
			}
		}

		ipaReef = strings.ReplaceAll(ipaReef, "t͡sj", "tʃ")
		ipaReef = strings.ReplaceAll(ipaReef, "sj", "ʃ")

		temp := ""
		runes := []rune(ipaReef)

		// Glottal stops between two vowels are removed
		for i, a := range runes {
			if i != 0 && i != len(runes)-1 && a == 'ʔ' {
				if runes[i-1] == '.' && i > 1 {
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
				} else if runes[i-1] == 'ˈ' && i > 2 {
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

	ipaReef = strings.TrimPrefix(ipaReef, ".")

	ipaReef = strings.ReplaceAll(ipaReef, "*.", " ")

	// now Romanize the reef IPA
	word := strings.Split(ipaReef, " ")

	breakdown = ""

	for j := 0; j < len(word); j++ {
		word[j] = strings.ReplaceAll(word[j], "]", "")
		word[j] = strings.ReplaceAll(word[j], "[", "")
		// "or" means there's more than one IPA in this word, and we only want one
		if word[j] == "or" {
			breakdown += "or "
			continue
		}

		syllables := strings.Split(word[j], ".")

		/* Onset */
		for k := 0; k < len(syllables); k++ {
			breakdown += "-"

			stressed := false
			syllable := strings.ReplaceAll(syllables[k], "·", "")
			if strings.Contains(syllable, "ˈ") {
				stressed = true
				breakdown += "__"
			}
			syllable = strings.ReplaceAll(syllable, "ˈ", "")
			syllable = strings.ReplaceAll(syllable, "ˌ", "")

			// tsy
			if strings.HasPrefix(syllable, "tʃ") {
				breakdown += "ch"
				syllable = strings.TrimPrefix(syllable, "tʃ")
			} else if len(syllable) >= 4 && syllable[0:4] == "t͡s" {
				// ts
				breakdown += "ts"
				//tsp
				if hasAt("ptk", syllable, 3) {
					if nth_rune(syllable, 4) == "'" {
						// ts + ejective onset
						breakdown += romanization2[syllable[4:6]]
						syllable = syllable[6:]
					} else {
						// ts + unvoiced plosive
						breakdown += romanization2[string(syllable[4])]
						syllable = syllable[5:]
					}
				} else if hasAt("lɾmnŋwj", syllable, 3) {
					// ts + other consonent
					breakdown += romanization2[nth_rune(syllable, 3)]
					syllable = syllable[4+len(nth_rune(syllable, 3)):]
				} else {
					// ts without a cluster
					syllable = syllable[4:]
				}
			} else if hasAt("fs", syllable, 0) {
				//
				breakdown += nth_rune(syllable, 0)
				if hasAt("ptk", syllable, 1) {
					if nth_rune(syllable, 2) == "'" {
						// f/s + ejective onset
						breakdown += romanization2[syllable[1:3]]
						syllable = syllable[3:]
					} else {
						// f/s + unvoiced plosive
						breakdown += romanization2[string(syllable[1])]
						syllable = syllable[2:]
					}
				} else if hasAt("lɾmnŋwj", syllable, 1) {
					// f/s + other consonent
					breakdown += romanization2[nth_rune(syllable, 1)]
					syllable = syllable[1+len(nth_rune(syllable, 1)):]
				} else {
					// f/s without a cluster
					syllable = syllable[1:]
				}
			} else if hasAt("ptk", syllable, 0) {
				if nth_rune(syllable, 1) == "'" {
					// ejective
					breakdown += romanization2[syllable[0:2]]
					syllable = syllable[2:]
				} else {
					// unvoiced plosive
					breakdown += romanization2[string(syllable[0])]
					syllable = syllable[1:]
				}
			} else if hasAt("ʔlɾhmnŋvwjzbdg", syllable, 0) {
				// other normal onset
				breakdown += romanization2[nth_rune(syllable, 0)]
				syllable = syllable[len(nth_rune(syllable, 0)):]
			} else if hasAt("ʃʒ", syllable, 0) {
				// one sound representd as a cluster
				if nth_rune(syllable, 0) == "ʃ" {
					breakdown += "sh"
				}
				syllable = syllable[len(nth_rune(syllable, 0)):]
			}

			/*
			 * Nucleus
			 */
			if len(syllable) > 1 && hasAt("jw", syllable, 1) {
				//diphthong
				breakdown += romanization2[syllable[0:len(nth_rune(syllable, 0))+1]]
				syllable = string([]rune(syllable)[2:])
			} else if len(syllable) > 1 && hasAt("lr", syllable, 0) {
				//psuedovowel
				breakdown += romanization2[syllable[0:3]]
				continue // psuedovowels can't coda
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

			if stressed {
				breakdown += "__"
			}
		}
		breakdown += " "
	}

	breakdown = strings.TrimPrefix(breakdown, "-")
	breakdown = strings.ReplaceAll(breakdown, " -", " ")
	breakdown = strings.TrimSuffix(breakdown, " ")

	// If there's a tìftang between two identical vowels, the tìftang is optional
	shortString := strings.ReplaceAll(strings.ReplaceAll(ipaReef, "ˈ", ""), ".", "")
	for _, a := range []string{"a", "ɛ", "ɪ", "o", "u", "i", "æ", "ʊ"} {
		if strings.Contains(shortString, a+"ʔ"+a) {
			// fix IPA
			noTìftangIPA := strings.ReplaceAll(ipaReef, a+".ˈʔ"+a, a+".ˈ"+a)
			noTìftangIPA = strings.ReplaceAll(noTìftangIPA, a+".ʔ"+a, a+"."+a)
			noTìftangIPA = strings.ReplaceAll(noTìftangIPA, a+"ʔ."+a, a+"."+a)
			noTìftangIPA = strings.ReplaceAll(noTìftangIPA, a+"ʔ.ˈ"+a, a+".ˈ"+a)

			ipaReef += "] or [" + noTìftangIPA
		}
	}

	// fix breakdown
	shortString = strings.ReplaceAll(breakdown, "-", "")
	for _, a := range []string{"a", "e", "ì", "o", "u", "i", "ä", "ù"} {
		if strings.Contains(shortString, a+"'"+a) {
			noTìftangBreakdown := strings.ReplaceAll(breakdown, a+"-'"+a, a+"-"+a)
			noTìftangBreakdown = strings.ReplaceAll(noTìftangBreakdown, a+"'-"+a, a+"-"+a)

			breakdown += " or " + noTìftangBreakdown
		}

	}

	return []string{breakdown, ipaReef}
}

func StartEverything() string {
	universalLock.Lock()
	start := time.Now()
	var errors = []error{
		AssureDict(),
		CacheDict(),
		CacheDictHash(),
		CacheDictHash2(),
	}
	for _, err := range errors {
		if err != nil {
			log.Println(err)
		}
	}
	universalLock.Unlock()
	PhonemeDistros()
	elapsed := strconv.FormatFloat(time.Since(start).Seconds(), 'f', -1, 64)
	return fmt.Sprintln("Everything is cached.  Took " + elapsed + " seconds")
}
