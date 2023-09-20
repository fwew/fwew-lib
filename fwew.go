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

// Translate some navi text.
// !! Only one word is allowed, if spaces are found, they will be treated like part of the word !!
// This will return an array of Words, that fit the input text
// One Navi-Word can have multiple meanings and words (e.g. synonyms)
func TranslateFromNavi(searchNaviWord string, checkFixes bool) (results []Word, err error) {
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

		if checkFixes && word.reconstruct(searchNaviWord) {
			//when it's a verb ending on -uyu, it adds one more to output
			if Contains(word.Affixes.Comment, []string{"flagUYU"}) {
				word.Affixes.Comment = []string{}
				results = append(results, word)
				word2 := word.CloneWordStruct()
				word2.Affixes.Infix = []string{}
				word2.Affixes.Suffix = []string{"yu"}
				results = append(results, word2)
			} else {
				word.Navi = naviWord
				results = append(results, word)
			}
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
		case "tr":
			wordString += word.TR
		}
		wordString = StripChars(wordString, ",;.:?!()")
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

	if err == nil {
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

/*
 * Name generators
 */
func SingleNames(name_count int, dialect int, syllable_count int) (output string) {
	// Make sure the numbers are good
	if name_count > 50 || name_count <= 0 || syllable_count > 4 || syllable_count < 0 {
		return "Max name count is 50, max syllable count is 4"
	}

	phoneme_distros()
	calculate_rand_params()

	// Charts and variables used for formatting
	output = ""

	// Fill the chart with names
	for i := 0; i < name_count; i++ {
		output += string(single_name_gen(rand_if_zero(syllable_count), dialect)) + "\n"
	}

	return output
}

func FullNames(ending string, name_count int, dialect int, syllable_count [3]int) (output string) {
	// Make sure the numbers are good
	if name_count > 50 || name_count <= 0 {
		return "Max name count is 50, max syllable count is 4"
	}

	for i := 0; i < 3; i++ {
		if syllable_count[i] > 4 || syllable_count[i] < 0 {
			return "Max name count is 50, max syllable count is 4"
		}
	}

	phoneme_distros()
	calculate_rand_params()

	// Charts and variables used for formatting
	output = ""

	// Fill the chart with names
	for i := 0; i < name_count; i++ {
		output += string(single_name_gen(rand_if_zero(syllable_count[0]), dialect))
		output += " te "
		output += string(single_name_gen(rand_if_zero(syllable_count[1]), dialect))
		output += " "
		output += string(single_name_gen(rand_if_zero(syllable_count[2]), dialect))
		output += ending + "\n"
	}

	return output
}

func NameAlu(name_count int, dialect int, syllable_count int, noun_mode int, adj_mode int) (output string) {
	// Make sure the numbers are good
	if name_count > 50 || name_count <= 0 || syllable_count > 4 || syllable_count < 0 {
		return "Max name count is 50, max syllable count is 4"
	}

	phoneme_distros()
	calculate_rand_params()

	output = ""

	for i := 0; i < name_count; i++ {
		output += string(single_name_gen(rand_if_zero(syllable_count), dialect))

		/* Noun */
		nmode := 0
		if noun_mode != 1 && noun_mode != 2 {
			nmode = rand.Intn(5) // 80% chance of normal noun
			if nmode == 4 {
				nmode = 2
			} else {
				nmode = 1
			}
		} else {
			nmode = noun_mode
		}

		//two_word_noun := false

		noun := ""
		switch nmode {
		case 1:
			noun_word, err := Random(1, []string{"pos is n."})
			if err != nil {
				fmt.Println(noun_word)
				fmt.Println(err)
				noun += noun_word[0].Navi + " "
			} else {
				return "Error in normal noun"
			}
		case 2:
			verb, err := Random(1, []string{"pos starts v"})
			if err != nil {
				noun += verb[0].Navi + "yu "
			} else {
				return "Error in verb-er"
			}
		default:
			return "Error: unknown noun type"
		}

		/*if len(strings.Split(noun, " ")) > 1 {
			print(noun)
			two_word_noun = true
		} else {
			output += noun + " "
		}

		// Adjective
		amode := 0
		if adj_mode == 0 {
			// "something" mode
			amode = rand.Intn(8) - 1
			if amode <= 2 {
				// 50% chance of normal adjective
				amode = 2
			} else if amode <= 5 {
				// Verb participles get two sides of the die
				amode = 5
			}
		} else if adj_mode < 1 || adj_mode > 7 {
			// "any" mode
			amode = rand.Intn(5)
		} else {
			amode = adj_mode
		}

		adj := ""
		switch amode {
		// no case 1 (no adjective)
		case 2: //nomal adjective
			adj_word, err := Random(1, []string{"pos is adj."})
			if err == nil {
				adj = adj_word[0].Navi
				// If the adj starts with a in forest, we don't need another a
				if !two_word_noun && (adj[0] != 'a' || dialect != 1) {
					if !(adj[:2] == "le" && adj != "ler" && adj != "leyr") {
						adj = "a" + adj
					}
				} else if two_word_noun && (adj[len(adj)-1] != 'a' || dialect != 1) {
					adj += "a "
				}
			} else {
				return "Error in adjective"
			}
		case 3: //genitive noun
			adj_word, err := Random(1, []string{"pos is n."})
			if err == nil {
				adj = adj_word[0].Navi
				if adj == "tsko swizaw" {
					adj = "Tsko Swizawyä"
				} else if adj == "toruk makto" {
					adj = "Torukä Maktoyuä"
				} else if adj == "mo a fngä'" {
					adj = "Moä a Fgnä'"
				} else {
					adj_rune := []rune(adj)
					if has("aeìiä", string(adj_rune[len(adj_rune)-1])) {
						adj += "yä"
					} else {
						adj += "ä"
					}
				}
			} else {
				return "Error in genitive noun"
			}
		case 4: //origin noun
			adj_word, err := Random(1, []string{"pos is n."})
			if err == nil {
				if adj == "tsko swizaw" {
					adj = "ta Tsko Swizaw"
				} else if adj == "toruk makto" {
					adj = "ta Torukä Makto"
				} else if adj == "mo a fngä'" {
					adj = "ta Mo a Fgnä'"
				} else {
					adj = "ta " + adj_word[0].Navi
				}
			} else {
				return "Error in origin noun"
			}
		case 5: //participle verb
			infix := "us"
			find_verb := []string{"Eltur", "tìtxen", "s..i"}
			for len(find_verb) == 2 != (len(find_verb) > 1 && find_verb[len(find_verb)-1] != "s..i") {
				adj_word, err := Random(1, []string{"pos starts v"})
				if err == nil {
					find_verb = strings.Split(adj_word[0].InfixDots, " ")
					if adj_word[0].PartOfSpeech == "vtr" && rand.Intn(2) == 1 {
						infix = "awn"
					}
				} else {
					return "Error in participle verb"
				}
			}
			adj = insert_infix(find_verb, infix)
		case 6: //active participle verb
			find_verb := []string{"Eltur", "tìtxen", "s..i"}
			for len(find_verb) == 2 != (len(find_verb) > 1 && find_verb[len(find_verb)-1] != "s..i") {
				adj_word, err := Random(1, []string{"pos starts v"})
				if err == nil {
					find_verb = strings.Split(adj_word[0].InfixDots, " ")
				} else {
					return "Error in active participle verb"
				}
			}
			adj = insert_infix(find_verb, "us")
		case 7: //passive participle verb
			find_verb := []string{"Eltur", "tìtxen", "s..i"}
			for len(find_verb) == 1 {
				adj_word, err := Random(1, []string{"pos starts vtr"})
				if err == nil {
					find_verb = strings.Split(adj_word[0].InfixDots, " ")
				} else {
					return "Error in passive participle verb"
				}
			}
			adj = insert_infix(find_verb, "us")
		}

		output += adj*/

		//if two_word_noun {
		output += " " + noun
		//}

		output += "\n"

		//fmt.Println(strconv.Itoa(nmode) + " " + strconv.Itoa(amode))
	}

	return output
}
