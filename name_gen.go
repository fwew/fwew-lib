package fwew_lib

import (
	"fmt"
	"math/rand"
	"slices"
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

// SingleNames generates a single given names
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
	for i := 0; i < nameCount; i++ {
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

	for i := 0; i < 3; i++ {
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
	for i := 0; i < nameCount; i++ {
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
	allNouns, allAdjectives, allVerbs, allTransitiveVerbs := SortedWords()

	output = ""

	// This isn't at the top because SortedWords calls List, which uses the same lock
	universalLock.Lock()
	defer universalLock.Unlock()

	for i := 0; i < nameCount; i++ {
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
			for k := 0; k < len(a); k++ {
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
			if adjMode == 0 {
				// "something" mode
				amode = rand.Intn(8) - 1
				if amode <= 2 {
					// 50% chance of normal adjective
					amode = 2
				} else if amode >= 5 {
					// Verb participles get two sides of the die
					amode = 5
				}
			} else if adjMode == -1 {
				// "any" mode
				amode = rand.Intn(5) + 1
			} else {
				amode = adjMode
			}

			adj := ""
			switch amode {
			// no case 1 (no adjective)
			case 2: // normal adjective
				adjWord := fastRandom(allAdjectives)
				adj = convertDialect(adjWord, dialect)
				adj = strings.ReplaceAll(adj, "-", "")

				// If the adj starts with a in forest, we don't need another a
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
				// If the adj starts with a in forest, we don't need another a
				if !twoWordNoun && (adj[0] != 'a' || dialect != 1) {
					adj = "a" + glottalCaps(adj)
				} else if twoWordNoun && (adj[len(adj)-1] != 'a' || dialect != 1) {
					adj = glottalCaps(adj) + "a"
				} else {
					adj = glottalCaps(adj)
				}
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

				// If the adj starts with a in forest, we don't need another a
				if !twoWordNoun && (adj[0] != 'a' || dialect != 1) {
					adj = "a" + glottalCaps(adj)
				} else if twoWordNoun && (adj[len(adj)-1] != 'a' || dialect != 1) {
					adj = glottalCaps(adj) + "a"
				} else {
					adj = glottalCaps(adj)
				}
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
				// If the adj starts with a in forest, we don't need another a
				if !twoWordNoun && (adj[0] != 'a' || dialect != 1) {
					adj = "a" + glottalCaps(adj)
				} else if twoWordNoun && (adj[len(adj)-1] != 'a' || dialect != 1) {
					adj = glottalCaps(adj) + "a"
				} else {
					adj = glottalCaps(adj)
				}
			}

			if len(adj) > 1 {
				output += " " + adj
			}
		}

		if twoWordNoun {
			output += " "
			nounWords := strings.Split(noun, " ")
			for _, a := range nounWords {
				output += glottalCaps(a) + " "
			}
			output = output[:len(output)-1]
		}

		output += "\n"
	}

	return output
}

func GetPhonemeDistrosMap(lang string) (allDistros [][][]string) {
	phonoLock.Lock()
	defer phonoLock.Unlock()
	// Non-English ones were pulled out of Google Translate unless it says VERIFIED
	headerRow := map[string][]string{
		"en": {"Onset", "Nucleus", "Coda"},          // English
		"de": {"Beginn", "Kern", "Coda"},            // German (Deutsch)
		"es": {"Inicio", "Núcleo", "Coda"},          // Spanish (Español)
		"et": {"Algus", "tuum", "Coda"},             // Estonian (Eesti)
		"fr": {"Début", "Noyau", "Coda"},            // French (Français)
		"hu": {"Szótagkezdet", "Szótagmag", "Coda"}, // Hungarian (Magyar)
		"it": {"Inizio", "Nucleo", "Coda"},          // Italian (Italiano)
		"ko": {"초성(두음)", "중성(음절핵)", "종성(말음)"},       // Korean (한국어)
		"nl": {"Begin", "Kern", "Coda"},             // Dutch (Nederlands)
		"pl": {"Początek", "Jądro", "Kod"},          // Polish (Polski)
		"pt": {"Início", "Núcleo", "Coda"},          // Portuguese (Português)
		"ru": {"Инициаль", "Ядро", "Финаль"},        // VERIFIED: Russian (Русский)
		"sv": {"Debut", "Nucleus", "Coda"},          // Swedish (Svenska)
		"tr": {"Başlangıç", "çekirdek", "Kodası"},   // Turkish (Türkçe)
		"uk": {"Початок", "Ядро", "Кода"},           // Ukrainian (Українська)
	}

	clusterName := map[string]string{
		"en": "Consonant Clusters",        // English
		"de": "Konsonantengruppen",        // German (Deutsch)
		"es": "Grupos de consonantes",     // Spanish (Español)
		"et": "Konsonantide klastrid",     // Estonian (Eesti)
		"fr": "Groupes de consonnes",      // French (Français)
		"hu": "Mássalhangzócsoportok",     // Hungarian (Magyar)
		"it": "Gruppi Consonantici",       // Italian (Italiano)
		"ko": "자음군",                       // Korean (한국어)
		"nl": "Medeklinkerclusters",       // Dutch (Nederlands)
		"pl": "Zbiory spółgłosek",         // Polish (Polski)
		"pt": "Aglomerados de consoantes", // Portuguese (Português)
		"ru": "Кластеры согласных",        // VERIFIED: Russian (Русский)
		"sv": "Konsonantkluster",          // Swedish (Svenska)
		"tr": "Ünsüz harfler",             // Turkish (Türkçe)
		"uk": "Збори приголосних",         // Ukrainian (Українська)
	}

	// Default to English
	headerLang := []string{"Onset", "Nucleus", "Coda"}
	clusterLang := "Consonant Clusters"

	if a, ok := headerRow[lang]; ok {
		headerLang = a
	}
	if a, ok := clusterName[lang]; ok {
		clusterLang = a
	}

	allDistros = [][][]string{
		{headerLang},
		{{clusterLang, "f", "s", "ts"}},
	}

	// Convert them to tuples for sorting
	var onsetTuples []PhonemeTuple
	for key, val := range onsetMap {
		onsetTuples = append(onsetTuples, PhonemeTuple{val, key})
	}
	slices.SortFunc(Tuples(onsetTuples), func(a, b PhonemeTuple) int {
		return b.value - a.value
	})

	var nucleusTuples []PhonemeTuple
	for key, val := range nucleusMap {
		nucleusTuples = append(nucleusTuples, PhonemeTuple{val, key})
	}
	slices.SortFunc(Tuples(nucleusTuples), func(a, b PhonemeTuple) int {
		return b.value - a.value
	})

	var codaTuples []PhonemeTuple
	for key, val := range codaMap {
		codaTuples = append(codaTuples, PhonemeTuple{val, key})
	}
	slices.SortFunc(Tuples(codaTuples), func(a, b PhonemeTuple) int {
		return b.value - a.value
	})

	// Probably not needed but just in case any other number exceeds it
	maxLen := len(onsetTuples)
	if len(nucleusTuples) > maxLen {
		maxLen = len(nucleusTuples)
	}
	if len(codaTuples) > maxLen {
		maxLen = len(codaTuples)
	}

	// Put them into a 2d string array
	i := 0
	for i < maxLen {
		allDistros[0] = append(allDistros[0], []string{})
		c := len(allDistros[0]) - 1

		if i < len(onsetTuples) {
			allDistros[0][c] = append(allDistros[0][c], onsetTuples[i].letter+" "+strconv.Itoa(onsetTuples[i].value))
		} else {
			allDistros[0][c] = append(allDistros[0][c], "")
		}

		if i < len(nucleusTuples) {
			allDistros[0][c] = append(allDistros[0][c], nucleusTuples[i].letter+" "+strconv.Itoa(nucleusTuples[i].value))
		} else {
			allDistros[0][c] = append(allDistros[0][c], "")
		}

		if i < len(codaTuples) {
			allDistros[0][c] = append(allDistros[0][c], codaTuples[i].letter+" "+strconv.Itoa(codaTuples[i].value))
		} else {
			allDistros[0][c] = append(allDistros[0][c], "")
		}
		i += 1
	}

	// Cluster time
	cluster1Full := []string{"f", "s", "ts"}
	cluster2Full := []string{"k", "kx", "l", "m", "n", "ng", "p",
		"px", "t", "tx", "r", "w", "y"}

	for _, a := range cluster2Full {
		allDistros[1] = append(allDistros[1], []string{a})
		c := len(allDistros[1]) - 1
		for _, b := range cluster1Full {
			allDistros[1][c] = append(allDistros[1][c], strconv.Itoa(clusterMap[b][a]))
		}
	}

	return
}
