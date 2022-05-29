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
func TranslateFromNavi(searchNaviWord string) (results []Word, err error) {
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

	err = RunOnDict(func(word Word) error {
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
	RunOnDict(func(word Word) error {
		var wordString string
		switch langCode {
		case "de":
			wordString += word.DE
		case "en":
			wordString += word.EN
		case "et":
			wordString += word.ET
		case "fr":
			wordString += word.FR
		case "hu":
			wordString += word.HU
		case "nl":
			wordString += word.NL
		case "pl":
			wordString += word.PL
		case "ru":
			wordString += word.RU
		case "sv":
			wordString += word.SV
		}
		wordString = StripChars(wordString, ",;.:?!")
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

	rand.Seed(time.Now().UnixNano())

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
