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
	"time"
)

// Global
const (
	idField  int    = 0  // dictionary.txt line Field 0 is Database ID
	lcField  int    = 1  // dictionary.txt line field 1 is Language Code
	navField int    = 2  // dictionary.txt line field 2 is Na'vi word
	ipaField int    = 3  // dictionary.txt line field 3 is IPA data
	infField int    = 4  // dictionary.txt line field 4 is Infix location data
	posField int    = 5  // dictionary.txt line field 5 is Part of Speech data
	defField int    = 6  // dictionary.txt line field 6 is Local definition
	srcField int    = 7  // dictionary.txt line field 7 is Source data
	stsField int    = 8  // dictionary.txt line field 8 is Stressed syllable #
	sylField int    = 9  // dictionary.txt line field 9 is syllable breakdown
	ifdField int    = 10 // dictionary.txt line field 10 is dot-style infix data
	space    string = " "
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

// Translate some navi text.
// !! Only one word is allowed, if spaces are found, they will be treated like part of the word !!
// This will return an array of Words, that fit the input text
// One Navi-Word can have multiple meanings and words (e.g. synonyms)
func TranslateFromNavi(searchNaviWord string, languageCode string) (results []Word) {
	badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\`

	// remove all the sketchy chars from arguments
	for _, c := range badChars {
		searchNaviWord = strings.ReplaceAll(searchNaviWord, string(c), "")
	}

	// normalize tìftang character
	searchNaviWord = strings.ReplaceAll(searchNaviWord, "’", "'")
	searchNaviWord = strings.ReplaceAll(searchNaviWord, "‘", "'")

	// No Results if empty string after removing sketch chars
	if len(searchNaviWord) == 0 {
		return
	}

	RunOnDict(languageCode, func(word Word) error {
		// save original Navi word, we want to add "+" or "--" later again
		naviWord := word.Navi

		// remove "+" and "--", we want to be able to search with and without those!
		word.Navi = strings.ReplaceAll(word.Navi, "+", "")
		word.Navi = strings.ReplaceAll(word.Navi, "--", "")
		word.Navi = strings.ToLower(word.Navi)

		if word.Navi == searchNaviWord {
			word.Navi = naviWord
			results = append(results, word)
			return nil
		}

		// skip words that obviously won't work
		s := word.similarity(searchNaviWord)

		if debugMode {
			log.Printf("Target: %s | Line: %s | [%f]\n", searchNaviWord, word.Navi, s)
		}

		if s < 0.50 && !strings.HasSuffix(searchNaviWord, "eyä") {
			return nil
		}

		if word.reconstruct(searchNaviWord) {
			word.Navi = naviWord
			results = append(results, word)
		}

		return nil
	})

	return
}

func TranslateToNavi(searchWord string, langCode string) (results []Word) {
	RunOnDict(langCode, func(word Word) error {
		wordString := StripChars(word.Definition, ",;")
		wordString = strings.ToLower(wordString)
		searchWord = strings.ToLower(searchWord)

		// whole-word matching
		for _, w := range strings.Split(wordString, space) {
			if w == searchWord {
				results = append(results, word)
				break
			}
		}

		return nil
	})
	return
}

// Get random words out of the dictionary.
// If args are applied, the dict will be filtered for args before random words are chosen.
// args will be put into the `List()` algorithm.
func Random(langCode string, amount int, args []string) (results []Word, err error) {
	allWords, err := List(args, langCode)

	if err != nil {
		log.Printf("Error getting fullDing: %s", err)
		return
	}

	dictLength := len(allWords)

	rand.Seed(time.Now().UnixNano())

	// create random number
	if amount <= 0 {
		amount = rand.Intn(dictLength)
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
