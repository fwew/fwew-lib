package fwew_lib

import (
	"slices"
	"strings"
)

// TranslateFromNaviHash Translate some navi text.
// !! Multiple words are supported !!
// This will return a 2D array of Words that fit the input text
// The first word will only contain the query put into the "translate" command
// One Navi-Word can have multiple meanings and words (e.g., synonyms)
func TranslateFromNaviHash(searchNaviWords string, checkFixes bool, strict bool, allowReef bool) (results [][]Word) {
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
		j, newWords := translateFromNaviHashHelper(dict, i, allWords, checkFixes, strict, allowReef)

		for _, newWord := range newWords {
			// Set up a receptacle for words
			results = append(results, []Word{})
			results[len(results)-1] = append(results[len(results)-1], newWord...)
		}

		if len(results[len(results)-1]) > 1 && len(strings.Split(results[len(results)-1][1].Navi, " ")) > 1 {
			var newQuery strings.Builder
			kOffset := 0
			for k := range strings.Split(results[len(results)-1][1].Navi, " ") {
				if i+k+kOffset >= len(allWords) {
					break
				}
				if allWords[i+k+kOffset] == "ke" || strings.ReplaceAll(allWords[i+k+kOffset], "e", "ä") == "rä'ä" {
					kOffset += 1
				}
				if k != 0 {
					newQuery.WriteString(" ")
				}
				newQuery.WriteString(allWords[i+k+kOffset])
				if strings.HasSuffix(allWords[i+k+kOffset], "-susi") {
					break
				}
			}
			results[len(results)-1][0].Navi = newQuery.String()
		}
		i += j
		i++
	}

	return
}

// appendToFront is a helper for translateFromNaviHashHelper
func appendToFront(words []Word, input Word) []Word {
	// Ensure it's not a duplicate
	for i, a := range words {
		if i != 0 && equal(input, a) {
			return words
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

// isVerb is a helper for translateFromNaviHashHelper
func isVerb(dict *map[string][]Word, input string, comparator string, strict bool, allowReef bool) (result bool, affixes Word) {
	affixes = simpleWord(input)
	_, possibilities := translateFromNaviHashHelper(dict, 0, []string{input}, true, strict, allowReef)
	_, possibilities2 := translateFromNaviHashHelper(dict, 0, []string{comparator}, true, strict, allowReef)
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
				if slices.Contains(b.Affixes.Prefix, prefix) {
					return false, affixes
				}
			}

			for _, suffix := range verbSuffixes {
				if slices.Contains(b.Affixes.Suffix, suffix) {
					return false, affixes
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
	return isRealVerb && pairFound && !unknownInfix, affixes
}

func stripAndCheckGlottalStop(a Word, containsGlottalStop bool) (string, bool) {
	strippedA := a.Navi
	if len(a.Affixes.Prefix) == 0 {
		strippedA = strings.TrimPrefix(strippedA, "'")
	}
	if len(a.Affixes.Suffix) == 0 {
		strippedA = strings.TrimSuffix(strippedA, "'")
	}
	if containsGlottalStop && !strings.Contains(strippedA, "'") {
		return strippedA, true // make sure tsa'u doesn't return tsa-au
	}
	return strippedA, false
}

func handleNegativeVerb(dict *map[string][]Word, allWords []string, i, j int, pairWord string, strict, allowReef bool, results *[][]Word) (bool, affix, int) {
	if i+j+2 < len(allWords) && (allWords[i+j+1] == "ke" || allWords[i+j+1] == "rä'ä" || allWords[i+j+1] == "ree") {
		validVerb, itsAffixes := isVerb(dict, allWords[i+j+2], pairWord, strict, allowReef)
		if validVerb {
			if len(*results) == 1 {
				*results = append(*results, []Word{simpleWord(allWords[i+j+1])})
				for _, b := range (*dict)[allWords[i+j+1]] {
					(*results)[1] = appendToFront((*results)[1], b)
				}
			}
			return true, itsAffixes.Affixes, 1
		}
	}
	return false, affix{}, 0
}

func filterDeconjugatedVerbs(dict *map[string][]Word, word string, strict, allowReef, containsUmlaut bool) []Word {
	var secondWords []Word
	for _, b := range testDeconjugations(dict, word, strict, allowReef, containsUmlaut) {
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

		secondWords = appendAndAlphabetize(secondWords, b)
	}
	return secondWords
}

func populateMultiwordResults(dict *map[string][]Word, a string, results *[][]Word, keepAffixes affix) {
	for _, definition := range (*dict)[a] {
		// Replace the word
		if len(*results) > 0 && len((*results)[0]) > 1 && ((*results)[0][1].Navi == "ke" || (*results)[0][1].Navi == "rä'ä") {
			// Get the query it's looking for
			(*results)[0][len((*results)[0])-1].Navi = (*results)[0][1].Navi
			(*results)[1] = appendToFront((*results)[1], definition)
			(*results)[1][1].Affixes = keepAffixes
		} else {
			// Get the query it's looking for
			(*results)[0] = appendToFront((*results)[0], definition)
			(*results)[0][1].Affixes = keepAffixes
		}
	}
}

func translateFromNaviHashHelper(dict *map[string][]Word, start int, allWords []string, checkFixes bool, strict bool, allowReef bool) (steps int, results [][]Word) {
	i := start

	var containsUmlaut []bool
	var containsGlottalStop []bool

	var tempResults []Word
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
				containsGlottalStop = append(containsGlottalStop, true)
			} else {
				containsGlottalStop = append(containsGlottalStop, false)
			}
		}

		results = [][]Word{{simpleWord(allWords[i])}}

		allWords = dialectCrunch(allWords, false, allowReef)

		searchNaviWord = allWords[i]

		// Find the word
		a := strings.ReplaceAll(searchNaviWord, "ù", "u")

		if _, ok := (*dict)[a]; ok {
			//bareNaviWord = true
			for _, b := range (*dict)[a] {
				results[len(results)-1] = appendAndAlphabetize(results[len(results)-1], b)
			}
		}

		// If one searches kivä, make sure kive doesn't show up
		for _, a := range results[len(results)-1] {
			if containsUmlaut[i] && !strings.Contains(strings.ToLower(a.Navi), "ä") {
				continue // ä can unstress to e, but not the other way around
			}
			_, skip := stripAndCheckGlottalStop(a, containsGlottalStop[i])
			if skip {
				continue
			}
			tempResults = append(tempResults, a)
		}
	} else {
		for range len(allWords) {
			containsUmlaut = append(containsUmlaut, true)
			containsGlottalStop = append(containsGlottalStop, false)
		}

		searchNaviWord = allWords[i]
		results = [][]Word{{simpleWord(allWords[i])}}

		// Find the word
		a := strings.ReplaceAll(searchNaviWord, "ù", "u")

		if _, ok := (*dict)[a]; ok {
			for _, b := range (*dict)[a] {
				results[len(results)-1] = appendAndAlphabetize(results[len(results)-1], b)
			}
		}
		// TODO: this is unreachable code because allowReef will always be false by this point.
		//else if allowReef {
		//	noUmlaut:= strings.ReplaceAll(a, "ä", "e")
		//	if _, ok:= (*dict)[noUmlaut]; ok {
		//		for _, b := range (*dict)[noUmlaut] {
		//			results[len(results)-1] = appendAndAlphabetize(results[len(results)-1], b)
		//		}
		//	}
		//}

		tempResults = append(tempResults, results[len(results)-1]...)
	}

	results[len(results)-1] = tempResults

	foundAlready := false

	// Bunch of duplicate code for the edge case of eltur tìtxen si and others like it
	//if !bareNaviWord {
	found := false
	// See if it is in the list known to start multiword words
	multiwords := &multiwordWords
	if !strict {
		multiwords = &multiwordWordsLoose
	} else if allowReef {
		multiwords = &multiwordWordsReef
	}
	if _, ok := (*multiwords)[searchNaviWord]; ok {
		// If so, loop through it
		for _, pairWordSet := range (*multiwords)[searchNaviWord] {
			if foundAlready {
				break
			}

			keepAffixes := *new(affix)

			extraWord := 0

			var revert strings.Builder
			revert.WriteString(results[0][0].Navi)
			// There could be more than one pair (win säpi and win si, for example)
			for j, pairWord := range pairWordSet {
				found = false
				// Don't cause an index-out-of-range error
				if i+j+1 >= len(allWords) {
					break
				} else {
					// For "[word] ke si and [word] rä'ä si"
					foundNegative, negativeAffixes, negativeExtraWord := handleNegativeVerb(dict, allWords, i, j, pairWord, strict, allowReef, &results)
					if foundNegative {
						extraWord = negativeExtraWord
						found = true
						foundAlready = true
						revert.WriteString(" " + allWords[i+j+1+extraWord])
						keepAffixes = negativeAffixes
						j += 1
						continue
					}

					// Verbs don't just come after ke or rä'ä
					validVerb, itsAffixes := isVerb(dict, allWords[i+j+1], pairWord, strict, allowReef)
					if validVerb {
						found = true
						foundAlready = true
						revert.WriteString(" " + allWords[i+j+1])
						keepAffixes = itsAffixes.Affixes
						continue
					}

					// Find all words the second word can represent
					var secondWords []Word

					// First by itself
					if pairWord == allWords[i+j+1] {
						found = true
						revert.WriteString(" " + allWords[i+j+1])
						continue
					}

					// And then by its possible conjugations
					secondWords = filterDeconjugatedVerbs(dict, allWords[i+j+1], strict, allowReef, containsUmlaut[i])

					// Do any of the conjugations work?
					for _, b := range secondWords {

						if b.Navi == pairWord {
							revert.WriteString(" " + b.Navi)
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
				results[0][0].Navi = revert.String()
				var fullWord strings.Builder
				fullWord.WriteString(searchNaviWord)
				for _, pairWord := range pairWordSet {
					fullWord.WriteString(" " + pairWord)
				}

				results[0] = []Word{results[0][0]}
				a := strings.ReplaceAll(fullWord.String(), "ù", "u")

				populateMultiwordResults(dict, a, &results, keepAffixes)
				i += len(pairWordSet) + extraWord
			}
		}
	}
	//}

	if checkFixes {
		var newResults []Word

		if !foundAlready {
			if len(results) > 0 && len(results[0]) > 0 {
				if !(strings.ToLower(results[len(results)-1][0].Navi) != searchNaviWord && strings.HasPrefix(strings.ToLower(results[len(results)-1][0].Navi), searchNaviWord)) {
					// Find all possible unconjugated versions of the word
					newResults = testDeconjugations(dict, searchNaviWord, strict, allowReef, containsUmlaut[i])
				}
			} else {
				// Find all possible unconjugated versions of the word
				newResults = testDeconjugations(dict, searchNaviWord, strict, allowReef, containsUmlaut[i])
			}
		}

		var tempNewResults []Word

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
			_, skip := stripAndCheckGlottalStop(a, containsGlottalStop[i])
			if skip {
				continue
			}
			tempNewResults = append(tempNewResults, a)
		}

		// Do not duplicate
		alreadyHere := results[len(results)-1]
		for _, a := range tempNewResults {
			add := true
			for _, b := range alreadyHere {
				if equal(b, a) {
					add = false
					break
				}
			}
			if add {
				results[len(results)-1] = append(results[len(results)-1], a)
			}
		}

		// Check if the word could have more than one word
		found := false
		// Find the results words

		var revert strings.Builder
		revert.WriteString(results[0][0].Navi)

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
						// There could be more than one pair (win säpi and win si, for example)
						for j, pairWord := range pairWordSet {
							found = false
							// Don't cause an index-out-of-range error
							if i+j+1 >= len(allWords) {
								break
							} else {
								// For "[word] ke si and [word] rä'ä si"
								foundNegative, negativeAffixes, negativeExtraWord := handleNegativeVerb(dict, allWords, i, j, pairWord, strict, allowReef, &results)
								if foundNegative {
									extraWord = negativeExtraWord
									found = true
									foundAlready = true
									revert.WriteString(" " + allWords[i+j+1+extraWord])
									keepAffixes = negativeAffixes
									j += 1

									continue
								}

								// Find all words the second word can represent
								var secondWords []Word

								allWord := allWords[i+j+1]

								if !strict || allowReef {
									pairWord = dialectCrunch([]string{pairWord}, false, allowReef)[0]
									allWord = dialectCrunch([]string{allWord}, false, allowReef)[0]
								}

								// First by itself
								if pairWord == allWord {
									found = true
									revert.WriteString(" " + allWords[i+j+1])
									continue
								}

								// And then by its possible conjugations
								secondWords = filterDeconjugatedVerbs(dict, allWords[i+j+1], strict, allowReef, containsUmlaut[i])

								// Do any of the conjugations work?
								for _, b := range secondWords {
									if b.Navi == pairWord {
										revert.WriteString(" " + b.Navi)
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
							results[0][0].Navi = revert.String()
							var fullWord strings.Builder
							fullWord.WriteString(newSearch)
							for _, pairWord := range pairWordSet {
								fullWord.WriteString(" " + pairWord)
							}

							results[0] = []Word{results[0][0]}
							a := strings.ReplaceAll(fullWord.String(), "ù", "u")
							if !strict {
								a = dialectCrunch([]string{a}, false, allowReef)[0]
							}

							populateMultiwordResults(dict, a, &results, keepAffixes)
							i += len(pairWordSet) + extraWord
						}
					}
				}
			}
		}
	}

	// If we found nothing, at least return the query
	if len(results[0]) == 0 {
		return i - start, [][]Word{{simpleWord(searchNaviWord)}}
	}

	if len(results) == 2 {
		temp := results[0]
		results[0] = results[1]
		results[1] = temp
	}

	return i - start, results
}

func searchNatlangWord(wordMap map[string][]string, searchWord string) (results []Word) {

	// No Results if empty string after removing sketch chars
	if len(searchWord) == 0 {
		return
	}

	// Find the word
	if _, ok := wordMap[searchWord]; !ok {
		return results
	}

	firstResults := wordMap[searchWord]

	for i := range firstResults {
		for _, c := range dictHashStrict[firstResults[i]] {
			results = appendAndAlphabetize(results, c)
		}
	}

	return
}

func TranslateToNaviHash(searchWord string, langCode string) (results [][]Word) {
	universalLock.Lock()
	defer universalLock.Unlock()
	searchWord = clean(searchWord)

	results = [][]Word{}

	for word := range strings.SplitSeq(searchWord, " ") {
		// Skip empty words
		if len(word) == 0 {
			continue
		}
		results = append(results, []Word{})
		for _, a := range translateToNaviHashHelper(&dictHash2Parenthesis, word, langCode) {
			results[len(results)-1] = appendAndAlphabetize(results[len(results)-1], a)
		}
		// Append the query to the front of the list
		tempResults := []Word{simpleWord(word)}
		tempResults = append(tempResults, results[len(results)-1]...)
		results[len(results)-1] = tempResults
	}
	return
}

func translateToNaviHashHelper(dictionary *metaDict, searchWord string, langCode string) (results []Word) {
	results = []Word{}
	var langMetaMap = getLangMetaDictMap(dictionary)

	// Verify the search query is actually in the definition
	for _, a := range searchNatlangWord(langMetaMap[langCode], searchWord) {
		var langMap2 = getLangDefMap(&a)
		searchWords := searchTerms(langMap2[langCode], false)
		found := slices.Contains(searchWords, searchWord)
		if found {
			results = appendAndAlphabetize(results, a)
		}
	}

	return
}

// BidirectionalSearch Search in both directions.  The language context is with Eywa now :ipu:
// !! Multiple words are supported !!
// This will return a 2D array of Words that fit the input text
// One Word can have multiple meanings and words (e.g., synonyms)
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
		j, newWords := translateFromNaviHashHelper(ourDict, i, allWords, checkFixes, false, allowReef)

		var NaviIDs []string

		for _, newWord := range newWords {
			// Set up a receptacle for words
			results = append(results, []Word{})
			results[len(results)-1] = append(results[len(results)-1], newWord...)
			if len(newWord) > 1 {
				NaviIDs = append(NaviIDs, newWord[1].ID)
			}
		}

		// Search for natural language words
		var natlangWords []Word
		for _, a := range translateToNaviHashHelper(&dictHash2, allWords[i], langCode) {
			// Do not duplicate if the Na'vi word is in the definition
			if contains(NaviIDs, []string{a.ID}) {
				continue
			}
			// We want them alphabetized with their fellow natlang words...
			natlangWords = appendAndAlphabetize(natlangWords, a)
		}

		// ...but not with the Na'vi words
		results[len(results)-1] = append(results[len(results)-1], natlangWords...)

		i += j

		i++
	}
	return
}
