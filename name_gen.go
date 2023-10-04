package fwew_lib

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

/*
 * Name generators
 */
func SingleNames(name_count int, dialect int, syllable_count int) (output string) {
	// Make sure the numbers are good
	if name_count > 50 || name_count <= 0 || syllable_count > 4 || syllable_count < 0 {
		return "Max name count is 50, max syllable count is 4"
	}

	// Charts and variables used for formatting
	output = ""

	// Fill the chart with names
	for i := 0; i < name_count; i++ {
		output += glottal_caps(string(single_name_gen(rand_if_zero(syllable_count), dialect))) + "\n"
	}

	return output
}

func FullNames(ending string, name_count int, dialect int, syllable_count [3]int, two_thousand_limit bool) (output string) {
	// Make sure the numbers are good
	if name_count > 50 || name_count <= 0 {
		return "Max name count is 50, max syllable count is 4"
	}

	for i := 0; i < 3; i++ {
		if syllable_count[i] > 4 || syllable_count[i] < 0 {
			return "Max name count is 50, max syllable count is 4"
		}
	}

	// Charts and variables used for formatting
	output = ""

	// Fill the chart with names
	for i := 0; i < name_count; i++ {
		// Fill it with three names
		output += glottal_caps(string(single_name_gen(rand_if_zero(syllable_count[0]), dialect)))
		output += " te "
		output += glottal_caps(string(single_name_gen(rand_if_zero(syllable_count[1]), dialect)))
		output += " "
		output += glottal_caps(string(single_name_gen(rand_if_zero(syllable_count[2]), dialect)))

		ending2 := ending

		// we don't want Neytiri''itan
		if output[len(output)-1] == '\'' {
			output = output[:len(output)-1]
		}

		// In reef dialect, glottal stops between nonidentical vowels are dropped
		if dialect == 2 && has("aäeìouù", get_last_rune(output, 1)) {
			ending2 = ending[1:]
		}

		// Add the ending
		output += ending2 + "\n"
		if two_thousand_limit && len([]rune(output)) > 1914 {
			output += "(stopped at " + strconv.Itoa(i+1) + ". 2000 Character limit)"
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

func NameAlu(name_count int, dialect int, syllable_count int, noun_mode int, adj_mode int) (output string) {
	// Make sure the numbers are good
	if name_count > 50 || name_count <= 0 || syllable_count > 4 || syllable_count < 0 {
		return "Max name count is 50, max syllable count is 4"
	}

	// A single function that allows all these to be acquired with only one dictionary search
	allNouns, allAdjectives, allVerbs, allTransitiveVerbs := SortedWords()

	output = ""

	for i := 0; i < name_count; i++ {
		output += glottal_caps(string(single_name_gen(rand_if_zero(syllable_count), dialect)))

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

		two_word_noun := false

		noun := ""
		switch nmode {
		case 1:
			noun_word := fast_random(allNouns).Navi
			noun += noun_word
		case 2:
			verb := fast_random(allVerbs).Navi
			a := strings.Split(verb, " ")
			for k := 0; k < len(a); k++ {
				noun += a[k]
			}
			noun += "yu"
		default:
			return "Error: unknown noun type"
		}

		output += " alu"

		if len(strings.Split(noun, " ")) > 1 {
			two_word_noun = true
		} else {
			output += " " + glottal_caps(noun)
		}

		if adj_mode != 1 {
			// Adjective
			amode := 0
			if adj_mode == 0 {
				// "something" mode
				amode = rand.Intn(8) - 1
				if amode <= 2 {
					// 50% chance of normal adjective
					amode = 2
				} else if amode >= 5 {
					// Verb participles get two sides of the die
					amode = 5
				}
			} else if adj_mode == -1 {
				// "any" mode
				amode = rand.Intn(5) + 1
			} else {
				amode = adj_mode
			}

			adj := ""
			switch amode {
			// no case 1 (no adjective)
			case 2: //nomal adjective
				adj_word := fast_random(allAdjectives)
				adj = adj_word.Navi

				// If the adj starts with a in forest, we don't need another a
				if !two_word_noun && (strings.ToLower(string(adj[0])) != "a" || dialect != 1) {
					if !(adj[:2] == "le" && adj != "ler" && adj != "leyr") {
						adj = "a" + glottal_caps(adj)
					}
				} else if two_word_noun && (adj[len(adj)-1] != 'a' || dialect != 1) {
					adj = glottal_caps(adj) + "a"
				} else {
					adj = glottal_caps(adj)
				}
			case 3: //genitive noun
				adj_word := fast_random(allNouns)

				adj = adj_word.Navi
				if adj == "tsko swizaw" {
					adj = "Tsko Swizawyä"
				} else if adj == "toruk makto" {
					adj = "Torukä Maktoyuä"
				} else if adj == "mo a fngä'" {
					adj = "Moä a Fgnä'"
				} else {
					adj_rune := []rune(adj)
					if has("aeìiä", string(adj_rune[len(adj_rune)-1])) {
						adj = glottal_caps(adj) + "yä"
					} else {
						adj = glottal_caps(adj) + "ä"
					}
				}
			case 4: //origin noun
				adj_word := fast_random(allNouns)
				adj = adj_word.Navi
				if adj == "tsko swizaw" {
					adj = "ta Tsko Swizaw"
				} else if adj == "toruk makto" {
					adj = "ta Torukä Makto"
				} else if adj == "mo a fngä'" {
					adj = "ta Mo a Fgnä'"
				} else {
					adj = "ta " + glottal_caps(adj)
				}
			case 5: //participle verb
				infix := "us"
				find_verb := one_word_verb(allVerbs)
				// If it's transitive, 50% chance of <awn>
				if find_verb.PartOfSpeech[2] == 'r' && rand.Intn(2) == 0 {
					infix = "awn"
				}
				adj = insert_infix(strings.Split(find_verb.InfixDots, " "), infix)
				// If the adj starts with a in forest, we don't need another a
				if !two_word_noun && (adj[0] != 'a' || dialect != 1) {
					adj = "a" + glottal_caps(adj)
				} else if two_word_noun && (adj[len(adj)-1] != 'a' || dialect != 1) {
					adj = glottal_caps(adj) + "a"
				} else {
					adj = glottal_caps(adj)
				}
			case 6: //active participle verb
				find_verb := one_word_verb(allVerbs)
				adj = insert_infix(strings.Split(find_verb.InfixDots, " "), "us")
				// If the adj starts with a in forest, we don't need another a
				if !two_word_noun && (adj[0] != 'a' || dialect != 1) {
					adj = "a" + glottal_caps(adj)
				} else if two_word_noun && (adj[len(adj)-1] != 'a' || dialect != 1) {
					adj = glottal_caps(adj) + "a"
				} else {
					adj = glottal_caps(adj)
				}
			case 7: //passive participle verb
				find_verb := one_word_verb(allTransitiveVerbs)
				adj = insert_infix(strings.Split(find_verb.InfixDots, " "), "awn")
				// If the adj starts with a in forest, we don't need another a
				if !two_word_noun && (adj[0] != 'a' || dialect != 1) {
					adj = "a" + glottal_caps(adj)
				} else if two_word_noun && (adj[len(adj)-1] != 'a' || dialect != 1) {
					adj = glottal_caps(adj) + "a"
				} else {
					adj = glottal_caps(adj)
				}
			}

			if len(adj) > 1 {
				output += " " + adj
			}
		}

		if two_word_noun {
			noun_words := strings.Split(noun, " ")
			for _, a := range noun_words {
				output += glottal_caps(a) + " "
			}
			output = output[:len(output)-1]
		}

		output += "\n"
	}

	return output
}

func GetPhonemeDistrosMap() (allDistros map[string]map[string]map[string]int) {
	allDistros = map[string]map[string]map[string]int{
		"Clusters": cluster_map,
		"Others": {
			"Onsets": onset_map,
			"Nuclei": nucleus_map,
			"Codas":  coda_map,
		},
	}
	return
}
