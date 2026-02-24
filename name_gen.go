package fwew_lib

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type phonemeTuple struct {
	value  int
	letter string
}

type tuples []phonemeTuple

func (t tuples) Len() int {
	return len(t)
}

func (t tuples) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t tuples) Less(i, j int) bool {
	// bigger values first here
	if t[i].value == t[j].value {
		return AlphabetizeHelper(t[i].letter, t[j].letter)
	}
	return t[i].value > t[j].value
}

func formatAdj(adj string, twoWordNoun bool, dialect int) string {
	if !twoWordNoun && (adj[0] != 'a' || dialect != 1) {
		return "a" + glottalCaps(adj)
	} else if twoWordNoun && (adj[len(adj)-1] != 'a' || dialect != 1) {
		return glottalCaps(adj) + "a"
	}
	return glottalCaps(adj)
}

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
		nmode := 0
		if nounMode != 1 && nounMode != 2 {
			nmode = rand.Intn(5) // 80% chance of normal noun
			if nmode == 4 {
				nmode = 2
			} else {
				nmode = 1
			}
		} else {
			nmode = nounMode
		}

		twoWordNoun := false

		noun := ""
		switch nmode {
		case 1:
			nounWord := fastRandom(allNouns)
			noun += strings.ReplaceAll(convertDialect(nounWord, dialect), "-", "")
		default: // case 2:
			verb := fastRandom(allVerbs)
			a := strings.Split(convertDialect(verb, dialect), " ")
			for k := range a {
				noun += a[k]
			}
			noun = strings.ReplaceAll(noun, "-", "")
			noun += "yu"
		}

		output += " alu"

		if len(strings.Split(noun, " ")) > 1 {
			twoWordNoun = true
		} else {
			output += " " + glottalCaps(noun)
		}

		if adjMode != 1 {
			// Adjective
			amode := 0
			switch adjMode {
			case 0:
				// "something" mode
				amode = rand.Intn(8) - 1
				if amode <= 2 {
					// 50% chance of normal adjective
					amode = 2
				} else if amode >= 5 {
					// Verb participles get two sides of the die
					amode = 5
				}
			case -1:
				// "any" mode
				amode = rand.Intn(5) + 1
			default:
				amode = adjMode
			}

			adj := ""
			switch amode {
			// no case 1 (no adjective)
			case 2: // normal adjective
				adjWord := fastRandom(allAdjectives)
				adj = convertDialect(adjWord, dialect)
				adj = strings.ReplaceAll(adj, "-", "")

				// If the adj starts with "a" in forest, we don't need another a
				if !twoWordNoun && (strings.ToLower(string(adj[0])) != "a" || dialect != 1) {
					if (adj[:2] == "le" && adj != "ler" && adj != "leyr" && adj != "lewnga'") || adj == "lafyon" {
						adj = glottalCaps(adj) // le-adjectives
					} else {
						adj = "a" + glottalCaps(adj)
					}
				} else if twoWordNoun && (adj[len(adj)-1] != 'a' || dialect != 1) {
					adj = glottalCaps(adj) + "a"
				} else {
					adj = glottalCaps(adj) // forest dialect a-adjectives like axpa or alaksi
				}
			case 3: //genitive noun
				adjWord := fastRandom(allNouns)

				adj = strings.ToLower(adjWord.Navi)
				switch adj {
				case "tsko swizaw":
					adj = "Tsko Swizawyä"
				case "toruk makto", "torùk makto":
					if dialect == 0 || dialect == 2 {
						adj = "Torùkä Maktoyuä"
					} else {
						adj = "Torukä Maktoyuä"
					}
				case "mo a fngä'":
					adj = "Moä a Fgnä'"
				default:
					adj = convertDialect(adjWord, dialect)
					adjSplit := strings.Split(adj, " ")
					if hasAt("aeìiä", adjSplit[0], -1) {
						adjSplit[0] = adjSplit[0] + "yä"
					} else {
						adjSplit[0] = adjSplit[0] + "ä"
					}
					adj = ""
					for _, a := range adjSplit {
						adj += glottalCaps(a) + " "
					}
					adj = strings.TrimSuffix(adj, " ")
				}

				adj = strings.ReplaceAll(adj, "-", "")
			case 4: //origin noun
				adjWord := fastRandom(allNouns)
				adj = strings.ToLower(adjWord.Navi)
				switch adj {
				case "tsko swizaw":
					adj = "ta Tsko Swizaw"
				case "toruk makto", "torùk makto":
					if dialect == 0 || dialect == 2 {
						adj = "ta Torùkä Maktoyu"
					} else {
						adj = "ta Torukä Maktoyu"
					}
				case "mo a fngä'":
					adj = "ta Mo a Fgnä'"
				default:
					adj = convertDialect(adjWord, dialect)
					if twoWordNoun {
						adj = glottalCaps(adj) + "ta"
					} else {
						adj = "ta " + glottalCaps(adj)
					}
				}

				adj = strings.ReplaceAll(adj, "-", "")
			case 5: //participle verb
				infix := "us"
				findVerb := oneWordVerb(allVerbs)
				// If it's transitive, 50% chance of <awn>
				if findVerb.PartOfSpeech[2] == 'r' && rand.Intn(2) == 0 {
					infix = "awn"
				}
				adj = findVerb.InfixDots
				switch dialect {
				case 2: // reef
					adj = quickReef(adj)
					fallthrough
				case 0: // interdialect
					adj = specialU(adj, findVerb.IPA)
				}

				adj = insertInfix(strings.Split(adj, " "), infix)
				// If the adj starts with "a" in forest, we don't need another a
				adj = formatAdj(adj, twoWordNoun, dialect)
			case 6: //active participle verb
				findVerb := oneWordVerb(allVerbs)
				adj = findVerb.InfixDots
				switch dialect {
				case 2: // reef
					adj = quickReef(adj)
					fallthrough
				case 0: // interdialect
					adj = specialU(adj, findVerb.IPA)
				}

				adj = insertInfix(strings.Split(adj, " "), "us")

				// If the adj starts with "a" in forest, we don't need another a
				adj = formatAdj(adj, twoWordNoun, dialect)
			case 7: //passive participle verb
				findVerb := oneWordVerb(allTransitiveVerbs)
				adj = findVerb.InfixDots
				switch dialect {
				case 2: // reef
					adj = quickReef(adj)
					fallthrough
				case 0: // interdialect
					adj = specialU(adj, findVerb.IPA)
				}

				adj = insertInfix(strings.Split(adj, " "), "awn")
				// If the adj starts with "a" in forest, we don't need another a
				adj = formatAdj(adj, twoWordNoun, dialect)
			}

			if len(adj) > 1 {
				output += " " + adj
			}
		}

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
