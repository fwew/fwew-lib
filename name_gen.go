package fwew_lib

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type PhonemeTuple struct {
	value  int
	letter string
}

type Tuples []PhonemeTuple

func (s Tuples) Len() int {
	return len(s)
}

func (s Tuples) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Tuples) Less(i, j int) bool {
	// bigger values first here
	if s[i].value == s[j].value {
		return AlphabetizeHelper(s[i].letter, s[j].letter)
	}
	return s[i].value > s[j].value
}

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
	for i := 0; i < name_count; i++ {
		// Fill it with three names
		output += glottal_caps(string(single_name_gen(rand_if_zero(syllable_count[0]), dialect)))
		output += " te "
		output += glottal_caps(string(single_name_gen(rand_if_zero(syllable_count[1]), dialect)))
		output += " "
		output += glottal_caps(string(single_name_gen(rand_if_zero(syllable_count[2]), dialect)))

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
		if dialect == 2 && has("aäeìouù", get_last_rune(output, 1)) {
			ending2 = ending2[1:]
		}

		// Add the ending
		output += ending2 + "\n"
		if two_thousand_limit && len([]rune(output)) > 1914 {
			// (stopped at {count}. 2000 Character limit)
			output += strings.ReplaceAll(message_too_big["en"], "{count}", strconv.Itoa(i+1))
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
			noun_word := fast_random(allNouns)
			noun += strings.ReplaceAll(convertDialect(noun_word, dialect), "-", "")
		case 2:
			verb := fast_random(allVerbs)
			a := strings.Split(convertDialect(verb, dialect), " ")
			for k := 0; k < len(a); k++ {
				noun += a[k]
			}
			noun = strings.ReplaceAll(noun, "-", "")
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
				adj = convertDialect(adj_word, dialect)
				adj = strings.ReplaceAll(adj, "-", "")

				// If the adj starts with a in forest, we don't need another a
				if !two_word_noun && (strings.ToLower(string(adj[0])) != "a" || dialect != 1) {
					if (adj[:2] == "le" && adj != "ler" && adj != "leyr" && adj != "lewnga'") || adj == "lafyon" {
						adj = glottal_caps(adj) // le-adjectives
					} else {
						adj = "a" + glottal_caps(adj)
					}
				} else if two_word_noun && (adj[len(adj)-1] != 'a' || dialect != 1) {
					adj = glottal_caps(adj) + "a"
				} else {
					adj = glottal_caps(adj) // forest dialect a-adjectives like axpa or alaksi
				}
			case 3: //genitive noun
				adj_word := fast_random(allNouns)

				adj = strings.ToLower(adj_word.Navi)
				if adj == "tsko swizaw" {
					adj = "Tsko Swizawyä"
				} else if adj == "toruk makto" || adj == "torùk makto" {
					if dialect == 0 || dialect == 2 {
						adj = "Torùkä Maktoyuä"
					} else {
						adj = "Torukä Maktoyuä"
					}
				} else if adj == "mo a fngä'" {
					adj = "Moä a Fgnä'"
				} else {
					adj = convertDialect(adj_word, dialect)
					adjSplit := strings.Split(adj, " ")
					adj_rune := []rune(adjSplit[0])
					if has("aeìiä", string(adj_rune[len(adj_rune)-1])) {
						adjSplit[0] = adjSplit[0] + "yä"
					} else {
						adjSplit[0] = adjSplit[0] + "ä"
					}
					adj = ""
					for _, a := range adjSplit {
						adj += glottal_caps(a) + " "
					}
					adj = strings.TrimSuffix(adj, " ")
				}

				adj = strings.ReplaceAll(adj, "-", "")
			case 4: //origin noun
				adj_word := fast_random(allNouns)
				adj = strings.ToLower(adj_word.Navi)
				if adj == "tsko swizaw" {
					adj = "ta Tsko Swizaw"
				} else if adj == "toruk makto" || adj == "torùk makto" {
					if dialect == 0 || dialect == 2 {
						adj = "ta Torùkä Maktoyu"
					} else {
						adj = "ta Torukä Maktoyu"
					}
				} else if adj == "mo a fngä'" {
					adj = "ta Mo a Fgnä'"
				} else {
					adj = convertDialect(adj_word, dialect)
					if two_word_noun {
						adj = glottal_caps(adj) + "ta"
					} else {
						adj = "ta " + glottal_caps(adj)
					}
				}

				adj = strings.ReplaceAll(adj, "-", "")
			case 5: //participle verb
				infix := "us"
				find_verb := one_word_verb(allVerbs)
				// If it's transitive, 50% chance of <awn>
				if find_verb.PartOfSpeech[2] == 'r' && rand.Intn(2) == 0 {
					infix = "awn"
				}
				adj = find_verb.InfixDots
				switch dialect {
				case 2: // reef
					adj = quickReef(adj)
					fallthrough
				case 0: // interdialect
					adj = specialU(adj, find_verb.IPA)
				}

				adj = insert_infix(strings.Split(adj, " "), infix, dialect)
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
				adj = find_verb.InfixDots
				switch dialect {
				case 2: // reef
					adj = quickReef(adj)
					fallthrough
				case 0: // interdialect
					adj = specialU(adj, find_verb.IPA)
				}

				adj = insert_infix(strings.Split(adj, " "), "us", dialect)

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
				adj = find_verb.InfixDots
				switch dialect {
				case 2: // reef
					adj = quickReef(adj)
					fallthrough
				case 0: // interdialect
					adj = specialU(adj, find_verb.IPA)
				}

				adj = insert_infix(strings.Split(adj, " "), "awn", dialect)
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
			output += " "
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

func GetPhonemeDistrosMap(lang string) (allDistros [][][]string) {
	// Non-English ones were pulled out of Google translate
	header_row := map[string][]string{
		"en": {"Onset", "Nucleus", "Coda"},          // English
		"de": {"Beginn", "Kern", "Coda"},            // German (Deutsch)
		"es": {"Inicio", "Núcleo", "Coda"},          // Spanish (Español)
		"et": {"Algus", "tuum", "Coda"},             // Estonian (Eesti)
		"fr": {"Début", "Noyau", "Coda"},            // French (Français)
		"hu": {"Szótagkezdet", "Szótagmag", "Coda"}, // Hungarian (Magyar)
		"ko": {"초성(두음)", "중성(음절핵)", "종성(말음)"},       // Korean (한국어)
		"nl": {"Begin", "Kern", "Coda"},             // Dutch (Nederlands)
		"pl": {"Początek", "Jądro", "Kod"},          // Polish (Polski)
		"pt": {"Início", "Núcleo", "Coda"},          // Portuguese (Português)
		"ru": {"Начало", "Ядро", "Кода"},            // Russian (Русский)
		"sv": {"Debut", "Nucleus", "Coda"},          // Swedish (Svenska)
		"tr": {"Başlangıç", "çekirdek", "Kodası"},   // Turkish (Türkçe)
		"uk": {"Початок", "Ядро", "Кода"},           // Ukrainian (Українська)
	}

	cluster_name := map[string]string{
		"en": "Consonant Clusters",        // English
		"de": "Konsonantengruppen",        // German (Deutsch)
		"es": "Grupos de consonantes",     // Spanish (Español)
		"et": "Konsonantide klastrid",     // Estonian (Eesti)
		"fr": "Groupes de consonnes",      // French (Français)
		"hu": "Mássalhangzócsoportok",     // Hungarian (Magyar)
		"ko": "자음군",                       // Korean (한국어)
		"nl": "Medeklinkerclusters",       // Dutch (Nederlands)
		"pl": "Zbiory spółgłosek",         // Polish (Polski)
		"pt": "Aglomerados de consoantes", // Portuguese (Português)
		"ru": "Согласные кластеры",        // Russian (Русский)
		"sv": "Konsonantkluster",          // Swedish (Svenska)
		"tr": "Ünsüz harfler",             // Turkish (Türkçe)
		"uk": "Збори приголосних",         // Ukrainian (Українська)
	}

	// Default to English
	header_lang := []string{"Onset", "Nucleus", "Coda"}
	cluster_lang := "Consonant Clusters"

	if a, ok := header_row[lang]; ok {
		header_lang = a
	}
	if a, ok := cluster_name[lang]; ok {
		cluster_lang = a
	}

	allDistros = [][][]string{
		{header_lang},
		{{cluster_lang, "f", "s", "ts"}},
	}

	// Convert them to tuples for sorting
	onset_tuples := []PhonemeTuple{}
	for key, val := range onset_map {
		onset_tuples = append(onset_tuples, PhonemeTuple{val, key})
	}
	sort.Sort(Tuples(onset_tuples))

	nucleus_tuples := []PhonemeTuple{}
	for key, val := range nucleus_map {
		nucleus_tuples = append(nucleus_tuples, PhonemeTuple{val, key})
	}
	sort.Sort(Tuples(nucleus_tuples))

	coda_tuples := []PhonemeTuple{}
	for key, val := range coda_map {
		coda_tuples = append(coda_tuples, PhonemeTuple{val, key})
	}
	sort.Sort(Tuples(coda_tuples))

	// Probably not needed but just in case any other number exceeds it
	max_len := len(onset_tuples)
	if len(nucleus_tuples) > max_len {
		max_len = len(nucleus_tuples)
	}
	if len(coda_tuples) > max_len {
		max_len = len(coda_tuples)
	}

	// Put them into a 2d string array
	i := 0
	for i < max_len {
		allDistros[0] = append(allDistros[0], []string{})
		c := len(allDistros[0]) - 1

		if i < len(onset_tuples) {
			allDistros[0][c] = append(allDistros[0][c], onset_tuples[i].letter+" "+strconv.Itoa(onset_tuples[i].value))
		} else {
			allDistros[0][c] = append(allDistros[0][c], "")
		}

		if i < len(nucleus_tuples) {
			allDistros[0][c] = append(allDistros[0][c], nucleus_tuples[i].letter+" "+strconv.Itoa(nucleus_tuples[i].value))
		} else {
			allDistros[0][c] = append(allDistros[0][c], "")
		}

		if i < len(coda_tuples) {
			allDistros[0][c] = append(allDistros[0][c], coda_tuples[i].letter+" "+strconv.Itoa(coda_tuples[i].value))
		} else {
			allDistros[0][c] = append(allDistros[0][c], "")
		}
		i += 1
	}

	// Cluster time
	cluster_1 := []string{"f", "s", "ts"}
	cluster_2 := []string{"k", "kx", "l", "m", "n", "ng", "p",
		"px", "t", "tx", "r", "w", "y"}

	for _, a := range cluster_2 {
		allDistros[1] = append(allDistros[1], []string{a})
		c := len(allDistros[1]) - 1
		for _, b := range cluster_1 {
			allDistros[1][c] = append(allDistros[1][c], strconv.Itoa(cluster_map[b][a]))
		}
	}

	return
}
