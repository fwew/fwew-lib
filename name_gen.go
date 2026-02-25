package fwew_lib

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// SingleNames generates single given names
func SingleNames(nameCount int, dialect int, syllableCount int) (output string) {
	universalLock.Lock()
	defer universalLock.Unlock()
	// Make sure the numbers are good
	if nameCount > 50 || nameCount <= 0 || syllableCount > 4 || syllableCount < 0 {
		return "Max name count is 50, max syllable count is 4"
	}

	// Charts and variables used for formatting
	output = ""

	// Fill the chart with names
	for range nameCount {
		output += glottalCaps(singleNameGen(randIfZero(syllableCount), dialect)) + "\n"
	}

	return output
}

// FullNames generates full Na'vi names
func FullNames(ending string, nameCount int, dialect int, syllableCount [3]int, twoThousandLimit bool) (output string) {
	universalLock.Lock()
	defer universalLock.Unlock()
	// Make sure the numbers are good
	if nameCount > 50 || nameCount <= 0 {
		return "Max name count is 50, max syllable count is 4"
	}

	for i := range 3 {
		if syllableCount[i] > 4 || syllableCount[i] < 0 {
			return "Max name count is 50, max syllable count is 4"
		}
	}

	// Charts and variables used for formatting
	output = ""

	endings := map[string]string{
		"'itu":  "descendent",
		"'itan": "son",
		"'ite":  "daughter",
	}

	randomize := true

	if _, ok := endings[ending]; ok {
		randomize = false
	}

	// Fill the chart with names
	for i := range nameCount {
		// Fill it with three names
		output += glottalCaps(singleNameGen(randIfZero(syllableCount[0]), dialect))
		output += " te "
		output += glottalCaps(singleNameGen(randIfZero(syllableCount[1]), dialect))
		output += " "
		output += glottalCaps(singleNameGen(randIfZero(syllableCount[2]), dialect))

		ending2 := ending
		if randomize {
			pick := rand.Intn(3)
			switch pick {
			case 0:
				ending2 = "'itan"
			case 1:
				ending2 = "'ite"
			case 2:
				ending2 = "'itu"
			}
		}

		// we don't want Neytiri''itan
		if output[len(output)-1] == '\'' {
			output = output[:len(output)-1]
		}

		// In reef dialect, glottal stops between nonidentical vowels are dropped
		if dialect == 2 && hasAt("aäeìouù", output, 1) {
			ending2 = ending2[1:]
		}

		// Add the ending
		output += ending2 + "\n"
		if twoThousandLimit && len([]rune(output)) > 1914 {
			// (stopped at {count}. 2000-Character limit)
			output += strings.ReplaceAll(messageTooBig["en"], "{count}", strconv.Itoa(i+1))
			break
		}

		// We want to know what the message that exceeded 2000 characters looked like
		if len([]rune(output)) > 2000 {
			fmt.Println(output)
			fmt.Println("Made a name message with " + strconv.Itoa(i+1) + " names.")
		}
	}

	return output
}

// NameAlu generates <name> alu <noun> <adjective> names
func NameAlu(nameCount int, dialect int, syllableCount int, nounMode int, adjMode int) (output string) {
	// Make sure the numbers are good
	if nameCount > 50 || nameCount <= 0 || syllableCount > 4 || syllableCount < 0 {
		return "Max name count is 50, max syllable count is 4"
	}

	// A single function that allows all these to be acquired with only one dictionary search
	allNouns, allAdjectives, allVerbs, allTransitiveVerbs := sortedWords()

	output = ""

	// This isn't at the top because sortedWords calls List, which uses the same lock
	universalLock.Lock()
	defer universalLock.Unlock()

	for range nameCount {
		output += glottalCaps(singleNameGen(randIfZero(syllableCount), dialect))

		/* Noun */
		twoWordNoun := false
		noun := getNameAluNoun(getNMode(nounMode), dialect, allNouns, allVerbs)

		/* Alu */
		output += " alu"

		/* Noun */
		if len(strings.Split(noun, " ")) > 1 {
			twoWordNoun = true
		} else {
			output += " " + glottalCaps(noun)
		}

		/* Adjective */
		if adjMode != 1 {
			adj := getNameAluAdjective(getAdjMode(adjMode), allAdjectives, dialect, twoWordNoun, allNouns, allVerbs, allTransitiveVerbs)

			if len(adj) > 1 {
				output += " " + adj
			}
		}

		/* Two word noun */
		if twoWordNoun {
			output += " "
			nounWords := strings.SplitSeq(noun, " ")
			for a := range nounWords {
				output += glottalCaps(a) + " "
			}
			output = output[:len(output)-1]
		}

		output += "\n"
	}

	return output
}
