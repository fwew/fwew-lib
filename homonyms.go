package fwew_lib

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

var homonymsArray = []string{"", "", ""}
var candidates2 []string
var homoMap = map[string]int{}
var lenitors = []string{"px", "p", "ts", "tx", "t", "kx", "k", "'"}
var lenitionMap = map[string]string{
	"ts": "s",
	"t":  "s",
	"tx": "t",
	"p":  "f",
	"px": "p",
	"k":  "h",
	"kx": "k",
	"'":  "",
}

func DuplicateDetector(query string) bool {
	result := false
	query = " " + query + " "

	for i := 0; i < len(homonymsArray); i++ {
		temp := " " + homonymsArray[i] + " "
		if strings.Contains(temp, query) {
			return true
		}
	}

	return result
}

// Check for ones that are the exact same, no affixes needed
func StageOne() error {
	tempHoms := []string{}

	err := runOnFile(func(word Word) error {
		standardizedWord := word.Navi
		badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\"„“”«»`

		// remove all the sketchy chars from arguments
		for _, c := range badChars {
			standardizedWord = strings.ReplaceAll(standardizedWord, string(c), "")
		}

		// normalize tìftang character
		standardizedWord = strings.ReplaceAll(standardizedWord, "’", "'")
		standardizedWord = strings.ReplaceAll(standardizedWord, "‘", "'")

		// find everything lowercase
		standardizedWord = strings.ToLower(standardizedWord)

		// If the word appears more than once, record it
		if _, ok := dictHash[standardizedWord]; ok {
			found := false
			for _, a := range tempHoms {
				if a == standardizedWord {
					found = true
					break
				}
			}
			if !found {
				tempHoms = append(tempHoms, standardizedWord)
			}
		}
		if strings.Contains(standardizedWord, "é") {
			noAcute := strings.ReplaceAll(standardizedWord, "é", "e")
			found := false
			for _, a := range tempHoms {
				if a == noAcute {
					found = true
					break
				}
			}
			if !found {
				tempHoms = append(tempHoms, noAcute)
				tempHoms = append(tempHoms, standardizedWord)
			}
		}

		return nil
	})

	// Reverse the order to make accidental and new homonyms easier to see
	// Also make it a string for easier searching
	i := len(tempHoms)
	for i > 0 {
		i--
		homonymsArray[0] += tempHoms[i] + " "
	}

	homonymsArray[0] = strings.TrimSuffix(homonymsArray[0], " ")

	if err != nil {
		log.Printf("Error in homonyms stage 1: %s", err)
		return err
	}

	return nil
}

// Helper to turn a string into a list of known words

// Check for ones that are the exact same, no affixes needed
func StageTwo() error {
	tempHoms := []string{}

	err := runOnFile(func(word Word) error {
		standardizedWord := word.Navi
		// If the word can conjugate into something else, record it
		results, err := TranslateFromNaviHash(standardizedWord, true)
		if err == nil && len(results[0]) > 2 {
			allNaviWords := ""
			for i, a := range results[0] {
				if i != 0 { //&& i < 3 {
					tempHoms = append(tempHoms, a.Navi)
					allNaviWords += a.Navi + " "
				}
			}

			homoMap[allNaviWords] = 1
			fmt.Println(strconv.Itoa(len(results[0])) + " " + allNaviWords + " " + standardizedWord)
		}

		//Lenited forms, too
		found := false
		for _, a := range lenitors {
			if strings.HasPrefix(word.Navi, a) {
				//fmt.Println(word.Navi)
				word.Navi = strings.TrimPrefix(word.Navi, a)
				word.Navi = lenitionMap[a] + word.Navi
				found = true
				break
			}
		}
		if found {
			//fmt.Println(word.Navi)
			// If the word can conjugate into something else, record it
			results, err := TranslateFromNaviHash(word.Navi, true)
			if err == nil && len(results[0]) > 2 {
				allNaviWords := ""
				for i, a := range results[0] {
					if i != 0 { //&& i < 3 {
						tempHoms = append(tempHoms, a.Navi)
						allNaviWords += a.Navi + " "
					}
				}

				homoMap[allNaviWords] = 1
				fmt.Println(strconv.Itoa(len(results[0])) + " " + allNaviWords + " " + word.Navi)
			}
		}

		return nil
	})

	// Reverse the order to make accidental and new homonyms easier to see
	// Also make it a string for easier searching
	i := len(tempHoms)
	for i > 0 {
		i--
		homonymsArray[1] += tempHoms[i] + " "
	}

	homonymsArray[1] = strings.TrimSuffix(homonymsArray[1], " ")

	if err != nil {
		log.Printf("Error in homonyms stage 2: %s", err)
		return err
	}

	//fmt.Println(homonymsArray[1])

	return nil
}

// Helper for StageThree, based on reconstruct from affixes.go
func reconjugateNouns(input Word, inputNavi string, prefixCheck int, suffixCheck int, unlenite int8) error {
	switch prefixCheck {
	case 0:
		fallthrough
	case 1:
		fallthrough
	case 2:
		// Non-lenition prefixes for nouns only
		for _, element := range prefixes1Nouns {
			newWord := element + inputNavi
			candidates2 = append(candidates2, newWord)
			reconjugateNouns(input, newWord, 4, suffixCheck, -1)
		}
		fallthrough
	case 3:
		// This one will demand this makes it use lenition
		for _, element := range prefixes1lenition {
			// If it has a lenition-causing prefix
			newWord := element + inputNavi
			candidates2 = append(candidates2, newWord)
			reconjugateNouns(input, newWord, 4, suffixCheck, -1)
		}
		fallthrough
	case 4:
		for _, element := range stemPrefixes {
			// If it has a lenition-causing prefix
			newWord := element + inputNavi
			candidates2 = append(candidates2, newWord)
			reconjugateNouns(input, newWord, 5, suffixCheck, -1)
		}
		//fallthrough
	}

	switch suffixCheck {
	case 0: // adpositions, sì, o, case endings
		for _, element := range adposuffixes {
			newWord := inputNavi + element
			candidates2 = append(candidates2, newWord)
			reconjugateNouns(input, newWord, prefixCheck, 3, -1)
		}
		fallthrough
	case 1:
		fallthrough
	case 2:
		for _, element := range determinerSuffixes {
			newWord := inputNavi + element
			candidates2 = append(candidates2, newWord)
			reconjugateNouns(input, newWord, prefixCheck, 3, -1)
		}
		fallthrough
	case 3:
		fallthrough
	case 4:
		for _, element := range stemSuffixes {
			newWord := inputNavi + element
			candidates2 = append(candidates2, newWord)
			reconjugateNouns(input, newWord, prefixCheck, 5, -1)
		}
	}

	return nil
}

// Helper for ReconjugateVerbs
func removeBrackets(input string) string {
	input = strings.ReplaceAll(input, "<0>", "")
	input = strings.ReplaceAll(input, "<1>", "")
	input = strings.ReplaceAll(input, "<2>", "")
	return input
}

// Helper for StageThree, based on reconstruct from affixes.go
func reconjugateVerbs(inputNavi string, prefirstUsed bool, firstUsed bool, secondUsed bool) error {
	candidates2 = append(candidates2, removeBrackets(inputNavi))
	if !prefirstUsed {
		for _, a := range prefirst {
			reconjugateVerbs(strings.ReplaceAll(inputNavi, "<0>", a), true, firstUsed, secondUsed)
		}
		reconjugateVerbs(strings.ReplaceAll(inputNavi, "<0>", "äpeyk"), true, firstUsed, secondUsed)
	} else if !firstUsed {
		for _, a := range first {
			reconjugateVerbs(strings.ReplaceAll(inputNavi, "<1>", a), prefirstUsed, true, secondUsed)
		}
	} else if !secondUsed {
		for _, a := range second {
			reconjugateVerbs(strings.ReplaceAll(inputNavi, "<2>", a), prefirstUsed, firstUsed, true)
		}
	}

	return nil
}

func reconjugate(word Word, allowPrefixes bool) {
	// remove "+" and "--", we want to be able to search with and without those!
	word.Navi = strings.ReplaceAll(word.Navi, "+", "")
	word.Navi = strings.ReplaceAll(word.Navi, "--", "")
	word.Navi = strings.ToLower(word.Navi)

	if word.PartOfSpeech == "pn." {
		candidates2 = append(candidates2, "nì"+word.Navi)
	}

	if word.PartOfSpeech == "n." || word.PartOfSpeech == "pn." || word.PartOfSpeech == "Prop.n." {
		reconjugateNouns(word, word.Navi, 0, 0, 0)
		//Lenited forms, too
		found := false

		for _, a := range lenitors {
			if strings.HasPrefix(word.Navi, a) {
				//fmt.Println(word.Navi)
				word.Navi = strings.TrimPrefix(word.Navi, a)
				word.Navi = lenitionMap[a] + word.Navi
				found = true
				break
			}
		}
		if found {
			//fmt.Println(word.Navi)
			reconjugateNouns(word, word.Navi, 0, 0, 0)
		}
	} else if word.PartOfSpeech[0] == 'v' {
		reconjugateVerbs(word.InfixLocations, false, false, false)
		//None of these can productively combine with infixes
		if allowPrefixes {
			// Gerunds
			candidates2 = append(candidates2, removeBrackets("tì"+strings.ReplaceAll(word.InfixLocations, "<1>", "us")))
			//candidates2 = append(candidates2, removeBrackets("nì"+strings.ReplaceAll(word.InfixLocations, "<1>", "awn")))
			// [verb]-able
			candidates2 = append(candidates2, "tsuk"+word.Navi)
			candidates2 = append(candidates2, "atsuk"+word.Navi)
			candidates2 = append(candidates2, "tsuk"+word.Navi+"a")
			candidates2 = append(candidates2, "ketsuk"+word.Navi)
			candidates2 = append(candidates2, "aketsuk"+word.Navi)
			candidates2 = append(candidates2, "ketsuk"+word.Navi+"a")
		}
		// Ability to [verb]
		candidates2 = append(candidates2, word.Navi+"tswo")
		reconjugateNouns(word, word.Navi+"tswo", 0, 0, 0)

	} else if word.PartOfSpeech == "adj." {
		candidates2 = append(candidates2, word.Navi+"a")
		if allowPrefixes {
			candidates2 = append(candidates2, "a"+word.Navi)
			candidates2 = append(candidates2, "nì"+word.Navi)
		}
	}
}

func StageThree() (err error) {
	start := time.Now()

	tempHoms := []string{}

	wordCount := 0

	err = RunOnDict(func(word Word) error {
		wordCount += 1

		// Progress counter
		if wordCount%100 == 0 {
			fmt.Println("On word " + strconv.Itoa(wordCount))
		}
		// save original Navi word, we want to add "+" or "--" later again
		//naviWord := word.Navi

		// No multiword words
		if !strings.Contains(word.Navi, " ") {
			candidates2 = []string{word.Navi} //empty array of strings

			// Get conjugations
			reconjugate(word, true)

			//Lenited forms, too
			found := false
			for _, a := range lenitors {
				if strings.HasPrefix(word.Navi, a) {
					//fmt.Println(word.Navi)
					word.Navi = strings.TrimPrefix(word.Navi, a)
					word.Navi = lenitionMap[a] + word.Navi
					found = true
					break
				}
			}
			if found {
				//fmt.Println(word.Navi)
				candidates2 = append(candidates2, word.Navi)
				reconjugate(word, false)
			}

			for _, a := range candidates2 {

				results, err := TranslateFromNaviHash(a, true)
				if err == nil && len(results) > 0 && len(results[0]) > 2 {

					allNaviWords := ""
					noDupes := []string{}
					for i, b := range results[0] {
						if i == 0 {
							continue
						}
						dupe := false
						for _, c := range noDupes {
							if c == b.Navi {
								dupe = true
								break
							}
						}
						if !dupe { //&& i < 3 {
							noDupes = append(noDupes, b.Navi)
							allNaviWords += b.Navi + " "
						}
					}

					// No duplicates
					if _, ok := homoMap[allNaviWords]; !ok {
						homoMap[allNaviWords] = 1

						if len(noDupes) > 1 {
							fmt.Println(word.PartOfSpeech + ": -" + a + " " + word.Navi + "- -" + allNaviWords)
							tempHoms = append(tempHoms, a)
						}
					}
				}
			}

		}

		return nil
	})

	fmt.Println(homoMap)
	fmt.Println(tempHoms)

	total_seconds := time.Since(start)

	log.Printf("Stage three took " + strconv.Itoa(int(math.Floor(total_seconds.Hours()))) + " hours, " +
		strconv.Itoa(int(math.Floor(total_seconds.Minutes()))%60) + " minutes and " +
		strconv.Itoa(int(total_seconds.Seconds())%60) + " seconds")

	return
}

// Do everything
func homonymSearch() {
	StageOne()
	StageTwo()
	StageThree()
}
