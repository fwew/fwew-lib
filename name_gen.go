package fwew_lib

import (
	"math/rand"
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

	// Charts and variables used for formatting
	output = ""

	// Fill the chart with names
	for i := 0; i < name_count; i++ {
		output += glottal_caps(string(single_name_gen(rand_if_zero(syllable_count[0]), dialect)))
		output += " te "
		output += glottal_caps(string(single_name_gen(rand_if_zero(syllable_count[1]), dialect)))
		output += " "
		output += glottal_caps(string(single_name_gen(rand_if_zero(syllable_count[2]), dialect)))
		output += ending + "\n"
	}

	return output
}

func NameAlu(name_count int, dialect int, syllable_count int, noun_mode int, adj_mode int) (output string) {
	// Make sure the numbers are good
	if name_count > 50 || name_count <= 0 || syllable_count > 4 || syllable_count < 0 {
		return "Max name count is 50, max syllable count is 4"
	}

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
			noun_word, err := Random(1, []string{"pos", "is", "n."})
			if err != nil {
				return "Error in noun mode"
			}
			noun += noun_word[0].Navi + ""
		case 2:
			verb, err := Random(1, []string{"pos", "starts", "v"})
			if err != nil {
				return "Error in verb-er"
			}
			a := strings.Split(verb[0].Navi, " ")
			for k := 0; k < len(a); k++ {
				noun += a[k]
			}
			noun += "yu"
		default:
			return "Error: unknown noun type"
		}

		output += " alu "

		if len(strings.Split(noun, " ")) > 1 {
			two_word_noun = true
		} else {
			output += glottal_caps(noun) + " "
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
			adj_word, err := Random(1, []string{"pos", "is", "adj."})
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
			adj_word, err := Random(1, []string{"pos", "is", "n."})
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
			adj_word, err := Random(1, []string{"pos", "is", "n."})
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
			find_verb := one_word_verb(true)
			// If it's transitive, 50% chance of <awn>
			if find_verb.PartOfSpeech[2] == 'r' && rand.Intn(2) == 0 {
				infix = "awn"
			}
			adj = insert_infix(strings.Split(find_verb.InfixDots, " "), infix)
			// If the adj starts with a in forest, we don't need another a
			if !two_word_noun && (adj[0] != 'a' || dialect != 1) {
				adj = "a" + adj
			} else if two_word_noun && (adj[len(adj)-1] != 'a' || dialect != 1) {
				adj += "a "
			}
		case 6: //active participle verb
			find_verb := one_word_verb(true)
			adj = insert_infix(strings.Split(find_verb.InfixDots, " "), "us")
			// If the adj starts with a in forest, we don't need another a
			if !two_word_noun && (adj[0] != 'a' || dialect != 1) {
				adj = "a" + adj
			} else if two_word_noun && (adj[len(adj)-1] != 'a' || dialect != 1) {
				adj += "a "
			}
		case 7: //passive participle verb
			find_verb := one_word_verb(false)
			adj = insert_infix(strings.Split(find_verb.InfixDots, " "), "awn")
			// If the adj starts with a in forest, we don't need another a
			if !two_word_noun && (adj[0] != 'a' || dialect != 1) {
				adj = "a" + adj
			} else if two_word_noun && (adj[len(adj)-1] != 'a' || dialect != 1) {
				adj += "a "
			}
		}

		output += glottal_caps(adj)

		if two_word_noun {
			output += " " + glottal_caps(noun)
		}

		output += "\n"
	}

	return output
}
