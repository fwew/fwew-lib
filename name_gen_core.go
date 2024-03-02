package fwew_lib

/*
 *	The core of the name generator
 *
 *	All functions here except PhonemeFrequency are helper functions for
 * 	name_gen.go.  PhonemeFrequency is called on startup for info for
 *	the name generator.
 */

import (
	"math/rand"
	"strings"
	"unicode"
	"unicode/utf8"
)

/* To help deduce phonemes */
var romanization = map[string]string{
	// Vowels
	"a": "a", "i": "i", "ɪ": "ì",
	"o": "o", "ɛ": "e", "u": "u",
	"æ": "ä",
	// Diphthongs
	"aw": "aw", "ɛj": "ey",
	"aj": "ay", "ɛw": "ew",
	// Psuedovowels
	"ṛ": "rr", "ḷ": "ll",
	// Consonents
	"t": "t", "p": "p", "ʔ": "'",
	"n": "n", "k": "k", "l": "l",
	"s": "s", "ɾ": "r", "j": "y",
	"t͡s": "ts", "t'": "tx", "m": "m",
	"v": "v", "w": "w", "h": "h",
	"ŋ": "ng", "z": "z", "k'": "kx",
	"p'": "px", "f": "f", "r": "r",
	// Reef dialect
	"b": "px", "d": "tx", "g": "kx",
	"ʃ": "sy", "tʃ": "tsy", "ʊ": "ù",
	// mistakes and rarities
	"ʒ": "tsy", "": "", " ": ""}

/* The likelihood of each letter appearing in a specific part of a Na'vi syllable.
 * They're ordered most common first to save time in linear search (common case fast).
 * Someday they'll be calculated from dictionary-v3.txt upon startup. */
var onset_likelihood = [21]int{640, 498, 402, 398, 389, 359, 313, 311, 306, 196,
	193, 190, 188, 185, 180, 158, 140, 104, 83, 81, 76}
var onset_letters = [21]string{"t", "", "n", "k", "l", "s", "'", "p", "r", "y",
	"ts", "m", "tx", "v", "w", "h", "ng", "z", "kx", "px", "f"}
var onset_map = map[string]int{"t": 0, "": 0, "n": 0, "k": 0, "l": 0, "s": 0,
	"'": 0, "p": 0, "r": 0, "y": 0, "ts": 0, "m": 0, "tx": 0, "v": 0, "w": 0,
	"h": 0, "ng": 0, "z": 0, "kx": 0, "px": 0, "f": 0}

/* The clusters aren't arranged in this order, but they're only 1/8 of the onsets anyway.
 * Still waiting for a word with tspx. */
var cluster_likelihood = [39]int{
	19, 6, 15, 14, 17, 9, 31, 5, 29, 27, 20, 28, 66,
	14, 8, 26, 10, 32, 18, 15, 21, 41, 4, 49, 32, 57,
	5, 8, 11, 7, 1, 7, 2, 0, 16, 5, 9, 19, 64}
var cluster_letters = [39]string{
	"fk", "fkx", "fl", "fm", "fn", "fng", "fp", "fpx", "ft", "ftx", "fr", "fw", "fy",
	"sk", "skx", "sl", "sm", "sn", "sng", "sp", "spx", "st", "stx", "sr", "sw", "sy",
	"tsk", "tskx", "tsl", "tsm", "tsn", "tsng", "tsp", "tspx", "tst", "tstx", "tsr", "tsw", "tsy"}
var cluster_map = map[string]map[string]int{
	"f": {
		"k": 0, "kx": 0, "l": 0, "m": 0, "n": 0, "ng": 0, "p": 0,
		"px": 0, "t": 0, "tx": 0, "r": 0, "w": 0, "y": 0,
	},
	"s": {
		"k": 0, "kx": 0, "l": 0, "m": 0, "n": 0, "ng": 0, "p": 0,
		"px": 0, "t": 0, "tx": 0, "r": 0, "w": 0, "y": 0,
	},
	"ts": {
		"k": 0, "kx": 0, "l": 0, "m": 0, "n": 0, "ng": 0, "p": 0,
		"px": 0, "t": 0, "tx": 0, "r": 0, "w": 0, "y": 0,
	},
}

var nucleus_likelihood = [14]int{1226, 1021, 760, 704, 615, 564, 277, 209, 187, 158, 153, 152, 70, 61}
var nucleus_letters = [14]string{"a", "e", "ì", "o", "u", "i", "ä", "aw", "ey", "ù", "rr", "ay", "ew", "ll"}
var nucleus_map = map[string]int{"a": 0, "e": 0, "ì": 0, "o": 0, "u": 0, "i": 0,
	"ä": 0, "aw": 0, "ey": 0, "ù": 0, "rr": 0, "ay": 0, "ew": 0, "ll": 0}

var coda_likelihood = [13]int{3938, 497, 405, 288, 258, 192, 179, 176, 133, 69, 48, 22, 18}
var coda_letters = [13]string{"", "n", "m", "ng", "l", "k", "p", "'", "r", "t", "kx", "px", "tx"}
var coda_map = map[string]int{"": 0, "n": 0, "m": 0, "ng": 0, "l": 0,
	"k": 0, "p": 0, "'": 0, "r": 0, "t": 0, "kx": 0, "px": 0, "tx": 0}

var valid_triple_consonants = map[string]map[string]map[string]int{
	"'": {
		"f": {
			"y": 0,
		},
	},
}

var multiword_words = map[string][][]string{}

/* Calculated on startup to assist the random number generators and letter selector */
var max_onset = 0
var max_non_cluster = 0
var max_nucleus = 0
var max_coda = 0

/* Helper function to find the start of a string */
func first_rune(word string) (letter string) {
	r := []rune(word)
	return string(r[:1])
}

/* Get the nth to last letter of a string */
func get_last_rune(word string, n int) (letter string) {
	r := []rune(word)
	if n > len(r) {
		n = len(r)
	}
	return string(r[len(r)-n : len(r)-n+1])
}

/* Take n letters off the end of a string */
func shave_rune(word string, n int) (letter string) {
	r := []rune(word)
	if n > len(r) {
		n = len(r) + 1
	}
	return string(r[:len(r)-n])
}

// helper for insert-infix
func quickReef(input string) string {
	output := strings.ReplaceAll(input, "tsy", "ch")
	output = strings.ReplaceAll(output, "sy", "sh")

	ejectiveMap := map[string]string{"px": "b", "tx": "d", "kx": "g"}

	for _, e := range []string{"px", "tx", "kx"} {
		// Don't convert clusters
		output = strings.ReplaceAll(output, "f"+e, "1")
		output = strings.ReplaceAll(output, "s"+e, "2")
		// convert syllable-initial ejectives
		for _, a := range []string{"a", "e", "i", "o", "u", "ì", "ù", "ä", "rr", "ll"} {
			output = strings.ReplaceAll(output, e+a, ejectiveMap[e]+a)
		}
		// restore clusters
		output = strings.ReplaceAll(output, "1", "f"+e)
		output = strings.ReplaceAll(output, "2", "s"+e)
	}

	temp := ""
	runes := []rune(output)

	for i, a := range runes {
		if i != 0 && i != len(runes)-1 && a == rune('\'') {
			if is_vowel(string(runes[i+1])) && is_vowel(string(runes[i-1])) {
				if runes[i+1] != runes[i-1] {
					continue
				}
			}
		}
		temp += string(a)
	}

	output = temp

	return output
}

// helper for insert-infix
func specialU(input string, ipa string) string {
	split := strings.Split(input, "u")

	runes := []rune(ipa)
	output := ""

	i := 0
	for _, a := range runes {
		if a == 'u' {
			output += split[i] + "u"
			i++
		} else if a == 'ʊ' {
			output += split[i] + "ù"
			i++
		}
	}
	output += split[i]

	return output
}

/* Helper function for name-alu */
func insert_infix(verb []string, infix string, dialect int) (output string) {
	output = ""
	found_infix := false
	for j := 0; j < len(verb); j++ {
		some_verb := []rune(verb[j])
		for k := 0; k < len(some_verb); k++ {
			if some_verb[k] == '.' {
				if !found_infix {
					output += infix
					found_infix = true
				}
			} else {
				output += string(some_verb[k])
			}
		}
		if j+1 < len(verb) {
			output += "-"
		}
	}
	return output
}

// Assistant function for name generating functions
func rand_if_zero(n int) (x int) {
	if n == 0 {
		return rand.Intn(3) + 2
	} else {
		return n
	}
}

/* Is it a vowel? (for when the psuedovowel bool won't work) */
func is_vowel(letter string) (found bool) {
	// Also arranged from most to least common (not accounting for diphthongs)
	vowels := []string{"a", "e", "u", "ì", "o", "i", "ä", "ù"}
	// Linear search
	for i := 0; i < 8; i++ {
		if letter == vowels[i] {
			return true
		}
	}
	return false
}

/* Randomly select an onset for a Na'vi syllable */
func get_onset() (onset string, cluster bool) {
	selector := rand.Intn(max_onset)
	// Clusters
	if selector > max_non_cluster { // If the number is too high for the non-cluster onsets,
		selector -= max_non_cluster // you get to skip all of them.  It saves time.
		// Linear search
		for i := 0; i < len(cluster_likelihood); i++ {
			if selector < cluster_likelihood[i] {
				return cluster_letters[i], true
			}
			selector -= cluster_likelihood[i]
		}
		return cluster_letters[len(cluster_letters)-1], true
	} else { // Non-clusters (single consonants)
		// Linear search
		for i := 0; i < len(onset_likelihood); i++ {
			if selector < onset_likelihood[i] {
				return onset_letters[i], false
			}
			selector -= onset_likelihood[i]
		}
		return onset_letters[len(onset_letters)-1], false
	}
}

/* Get a random Na'vi nucleus */
func get_nucleus() (onset string) {
	selector := rand.Intn(max_nucleus)
	// Linear search
	for i := 0; i < len(nucleus_likelihood); i++ {
		if selector < nucleus_likelihood[i] {
			return nucleus_letters[i]
		}
		selector -= nucleus_likelihood[i]
	}
	return nucleus_letters[len(nucleus_letters)-1]
}

/* Get a random Na'vi coda */
func get_coda() (onset string) {
	selector := rand.Intn(max_coda)
	// Linear search
	for i := 0; i < len(coda_likelihood); i++ {
		if selector < coda_likelihood[i] {
			return coda_letters[i]
		}
		selector -= coda_likelihood[i]
	}
	return coda_letters[len(coda_letters)-1]
}

// Helper function for name-alu()
func one_word_verb(verbList []Word) (words Word) {
	word := fast_random(verbList)
	find_verb := strings.Split(word.InfixDots, " ")

	/* The second condition here is a clever and efficient little thing
	 * one word and not si: allowed (e.g. "takuk")
	 * two words and not si: disallowed (e.g. "tswìk kxenerit")
	 * one word and si: disallowed ("si" only)
	 * two words and si: allowed (e.g. "unil si")
	 * Any three-word verb: disallowed ("eltur tìtxen si" only)
	 * != is used as an exclusive "or"
	 */
	for (len(find_verb) == 2) != (find_verb[len(find_verb)-1] == "s..i") {
		word = fast_random(verbList)
		find_verb = strings.Split(word.InfixDots, " ")
	}
	return word
}

/* Helper function: turn ejectives into voiced plosives for reef */
func reef_plosives(letter string) (voiced string) {
	if letter == "p" {
		return "b"
	} else if letter == "t" {
		return "d"
	} else if letter == "k" {
		return "g"
	}
	return "" // How we know if it's an error
}

/* Helper function: Replace an ejective with a voiced plosive. */
func reef_ejective(name string) (reefy_name string) {
	onset_new := ""
	last_third := get_last_rune(name, 3)

	if last_third == "x" { // Adjacent ejectives become adjacent voiced plosives, too
		onset_new += reef_plosives(get_last_rune(name, 4))
	} else if last_third == "n" && get_last_rune(name, 2) == "k" {
		onset_new += "-" // disambiguate on-gi vs o-ngi
	}

	onset_new += reef_plosives(get_last_rune(name, 2))

	if last_third == "x" {
		return shave_rune(name, 4) + onset_new
	}

	return shave_rune(name, 2) + onset_new
}

// Helper function for name-alu
func convertDialect(word Word, dialect int) string {
	output := ""
	switch dialect {
	case 0: // interdialect
		output += ReefMe(word.IPA, true)[0]
	case 2: // reef
		output += ReefMe(word.IPA, false)[0]
	default: // forest
		output += word.Navi
	}
	return output
}

/* Randomly construct a phonotactically valid Na'vi word
 * Dialect codes: 0 is interdialect, 1 is forest, 2 is reef */
func single_name_gen(syllable_count int, dialect int) (name string) {
	// Sometimes these things might be referenced across loops
	name = ""
	onset := ""
	nucleus := ""
	coda := ""
	psuedovowel := false // But not this
	cluster := false

	// Make a name with len syllables
	for i := 0; i < syllable_count; i++ {
		onset, cluster = get_onset()

		// Triple consonants are whitelisted
		if cluster && len(coda) > 0 { // don't want errors
			if !(coda == "t" && onset[0] == 's') { // t-s-kx is valid as ts-kx
				first_cluster_num := 1
				if onset[0] == 't' {
					first_cluster_num = 2
				}
				first_cluster := onset[:first_cluster_num]
				second_cluster := onset[first_cluster_num:]
				if _, ok := valid_triple_consonants[coda][first_cluster][second_cluster]; ok {
					// Do nothing.  We found a valid triple
				} else {
					onset = second_cluster
				}
			}
		}

		nucleus = get_nucleus()

		psuedovowel = false

		// These will be important later
		onsetlength := utf8.RuneCountInString(onset)
		namelength := utf8.RuneCountInString(name)

		// Get psuedovowel status
		if nucleus == "rr" || nucleus == "ll" {
			psuedovowel = true
			// Disallow onsets from imitating the psuedovowel
			if onsetlength > 0 {
				if get_last_rune(onset, 1) == "l" || get_last_rune(onset, 1) == "r" {
					onset = "'"
				}
				// If no onset, disallow the previous coda from imitating the psuedovowel
			} else if namelength > 0 {
				if get_last_rune(name, 1) == first_rune(nucleus) || is_vowel(get_last_rune(name, 1)) {
					onset = "'"
				}
				// No onset or loader thing?  Needs a thing to start
			} else {
				onset = "'"
			}
			// No identical vowels togther in forest
		} else if dialect == 1 {
			if nucleus == "ù" { // As of September 2023, the ratio of u to ù
				nucleus = "u" // was almost exactly 4 to 1 (615 to 158)
			}
			if onsetlength == 0 && namelength > 0 && get_last_rune(name, 1) == first_rune(nucleus) {
				onset = "y"
			}
		} else if nucleus_map["ù"] == 0 { //no psuedovowel or forest dialect
			// If only we didn't have to hardcode the likelihood of ù compared to u :ìì:
			if nucleus == "u" && rand.Intn(5) == 0 { // As of September 2023, the ratio of u to ù
				nucleus = "ù" // was almost exactly 4 to 1 (615 to 158)
			}
		}

		// Now that the onsets have settled, make sure they don't repeat a letter from the coda
		// No "ng-n" "t-tx", "o'-lll" becoming "o'-'ll" or anything like that
		if onsetlength > 0 && namelength > 0 {
			length := utf8.RuneCountInString(coda)
			if length == 0 {
				length = 1
			}
			if first_rune(onset) == get_last_rune(name, length) {
				onset = ""
			}
		}

		// You shawm futa sy and tsy become sh and ch XD
		if dialect == 2 {
			if onset == "sy" {
				onset = "sh"
			} else if onset == "tsy" {
				onset = "ch"
			}
		}

		// Apply the onset
		name += onset

		// namelength has changed, we need a new one
		namelength += utf8.RuneCountInString(onset)

		/*
		 * coda
		 */
		if !psuedovowel {
			coda = get_coda()
		} else {
			coda = ""
		}

		// reef dialect stuff
		if dialect == 2 && namelength > 1 { // In reef dialect,
			if get_last_rune(name, 1) == "x" { // if there's an ejective in the onset
				if namelength > 2 {
					// that's not in a cluster,
					last_rune := get_last_rune(name, 3)
					if !(last_rune == "s" || last_rune == "f") {
						// it becomes a voiced plosive
						name = reef_ejective(name)
					}
				} else {
					name = reef_ejective(name)
				}
			} else if !psuedovowel && get_last_rune(name, 1) == "'" && get_last_rune(name, 2) != first_rune(nucleus) {
				// 'a'aw is optionally 'aaw (the generator leaves it in)
				if is_vowel(get_last_rune(name, 2)) { // Does kaw'it become kawit in reef?
					name = shave_rune(name, 1)
				}
			}
		}

		// Finish the syllable
		name += nucleus + coda
	}
	return name
}

func glottal_caps(input string) (output string) {
	a := []rune(input)
	n := 1
	output = ""
	if a[0] == rune('\'') {
		output += "'" + string(unicode.ToUpper(a[1]))
		n = 2
	} else {
		output += string(unicode.ToUpper(a[0]))
	}
	for ; n < len(a); n++ {
		output += string(a[n])
	}
	return output
}

func fast_random(wordList []Word) (results Word) {
	dictLength := len(wordList)

	return wordList[rand.Intn(dictLength)]
}

func nth_rune(word string, n int) (output string) {
	r := []rune(word)
	if n >= len(r) {
		return ""
	}
	return string(r[n])
}

func has(word string, character string) (output bool) {
	r := []rune(word)
	if len(character) == 0 {
		return false
	}
	c := []rune(character)[0]
	for i := 0; i < len(r); i++ {
		if c == r[i] {
			return true
		}
	}
	return false
}

// Helper function for name-alu
func SortedWords() (nouns []Word, adjectives []Word, verbs []Word, transitiveVerbs []Word) {
	words, err := List([]string{}, 0)

	if err != nil || len(words) == 0 {
		return
	}

	for i := 0; i < len(words); i++ {
		if words[i].PartOfSpeech == "n." {
			nouns = append(nouns, words[i])
		} else if words[i].PartOfSpeech == "adj." {
			adjectives = append(adjectives, words[i])
		} else if words[i].PartOfSpeech[0] == 'v' {
			verbs = append(verbs, words[i])
			if words[i].PartOfSpeech[2] == 'r' {
				transitiveVerbs = append(transitiveVerbs, words[i])
			}
		}
	}
	return
}

func EqualSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Called on startup to feed and compile dictionary information into the name generator
func PhonemeDistros() {
	// get the dict
	words, err := List([]string{}, 0)

	if err != nil || len(words) == 0 {
		return
	}

	//set the maps to zero

	//Onsets
	for i := 0; i < len(onset_letters); i++ {
		onset_map[onset_letters[i]] = 0
	}

	//Clusters
	cluster_1 := []string{"f", "s", "ts"}
	cluster_2 := []string{"k", "kx", "l", "m", "n", "ng", "p",
		"px", "t", "tx", "r", "w", "y"}
	for i := 0; i < len(cluster_1); i++ {
		for j := 0; j < len(cluster_2); j++ {
			cluster_map[cluster_1[i]][cluster_2[j]] = 0
		}
	}

	//Nuclei
	for i := 0; i < len(nucleus_likelihood); i++ {
		nucleus_map[nucleus_letters[i]] = 0
	}

	//Codas
	for i := 0; i < len(coda_likelihood); i++ {
		coda_map[coda_letters[i]] = 0
	}

	// Look through all the words
	for i := 0; i < len(words); i++ { // reverse so tìpengayt wrrzeykärìp is searchable
		word := strings.Split(words[i].IPA, " ")

		// Piggybacking off of the frequency script to get all words with spaces
		all_words := strings.Split(strings.ToLower(words[i].Navi), " ")
		if len(all_words) > 1 {
			if _, ok := multiword_words[all_words[0]]; ok {
				// Ensure no duplicates
				appended := false

				// Append in a way that makes the longer words first
				temp := [][]string{}
				for _, j := range multiword_words[all_words[0]] {
					if !appended && len([]rune(all_words[1])) > len([]rune(j[0])) {
						temp = append(temp, all_words[1:])
						appended = true
					}
					temp = append(temp, j)
				}
				if len(temp) <= len(multiword_words[all_words[0]]) {
					temp = append(temp, all_words[1:])
				}

				multiword_words[all_words[0]] = temp
			} else {
				multiword_words[all_words[0]] = [][]string{all_words[1:]}
			}
		}

		for j := 0; j < len(word); j++ {
			word[j] = strings.Replace(word[j], "]", "", 1500)
			// "or" means there's more than one IPA in this word, and we only want one
			if word[j] == "or" {
				break
			}

			syllables := strings.Split(word[j], ".")
			coda := ""

			/* Onset */
			for k := 0; k < len(syllables); k++ {
				syllable := strings.Replace(syllables[k], "·", "", 1500)
				syllable = strings.Replace(syllable, "ˈ", "", 1500)
				syllable = strings.Replace(syllable, "ˌ", "", 1500)

				onset_if_cluster := [2]string{"", ""}

				// ts
				if len(syllable) >= 4 && syllable[0:4] == "t͡s" {
					onset_if_cluster[0] = "ts"
					//tsp
					if has("ptk", nth_rune(syllable, 3)) {
						if nth_rune(syllable, 4) == "'" {
							// ts + ejective onset
							cluster_map["ts"][romanization[syllable[4:6]]] = cluster_map["ts"][romanization[syllable[4:6]]] + 1
							onset_if_cluster[1] = romanization[syllable[4:6]]
							syllable = syllable[6:]
						} else {
							// ts + unvoiced plosive
							cluster_map["ts"][romanization[string(syllable[4])]] = cluster_map["ts"][romanization[string(syllable[4])]] + 1
							onset_if_cluster[1] = romanization[string(syllable[4])]
							syllable = syllable[5:]
						}
					} else if has("lɾmnŋwj", nth_rune(syllable, 3)) {
						// ts + other consonent
						cluster_map["ts"][romanization[nth_rune(syllable, 3)]] = cluster_map["ts"][romanization[nth_rune(syllable, 3)]] + 1
						onset_if_cluster[1] = romanization[nth_rune(syllable, 3)]
						syllable = syllable[4+len(nth_rune(syllable, 3)):]
					} else {
						// ts without a cluster
						onset_map["ts"] = onset_map["ts"] + 1
						syllable = syllable[4:]
					}
				} else if has("fs", nth_rune(syllable, 0)) {
					//
					onset_if_cluster[0] = string(syllable[0])
					if has("ptk", nth_rune(syllable, 1)) {
						if nth_rune(syllable, 2) == "'" {
							// f/s + ejective onset
							cluster_map[string(syllable[0])][romanization[syllable[1:3]]] = cluster_map[string(syllable[0])][romanization[syllable[1:3]]] + 1
							onset_if_cluster[1] = romanization[syllable[1:3]]
							syllable = syllable[3:]
						} else {
							// f/s + unvoiced plosive
							cluster_map[string(syllable[0])][romanization[string(syllable[1])]] = cluster_map[string(syllable[0])][romanization[string(syllable[1])]] + 1
							onset_if_cluster[1] = romanization[string(syllable[1])]
							syllable = syllable[2:]
						}
					} else if has("lɾmnŋwj", nth_rune(syllable, 1)) {
						// f/s + other consonent
						cluster_map[string(syllable[0])][romanization[nth_rune(syllable, 1)]] = cluster_map[string(syllable[0])][romanization[nth_rune(syllable, 1)]] + 1
						onset_if_cluster[1] = romanization[nth_rune(syllable, 1)]
						syllable = syllable[1+len(nth_rune(syllable, 1)):]
					} else {
						// f/s without a cluster
						onset_map[string(syllable[0])] = onset_map[string(syllable[0])] + 1
						syllable = syllable[1:]
					}
				} else if has("ptk", nth_rune(syllable, 0)) {
					if nth_rune(syllable, 1) == "'" {
						// ejective
						onset_map[romanization[syllable[0:2]]] = onset_map[romanization[syllable[0:2]]] + 1
						syllable = syllable[2:]
					} else {
						// unvoiced plosive
						onset_map[romanization[string(syllable[0])]] = onset_map[romanization[string(syllable[0])]] + 1
						syllable = syllable[1:]
					}
				} else if has("ʔlɾhmnŋvwjzbdg", nth_rune(syllable, 0)) {
					// other normal onset
					onset_map[romanization[nth_rune(syllable, 0)]] = onset_map[romanization[nth_rune(syllable, 0)]] + 1
					syllable = syllable[len(nth_rune(syllable, 0)):]
				} else if has("ʃʒ", nth_rune(syllable, 0)) {
					// one sound representd as a cluster
					if nth_rune(syllable, 0) == "ʃ" {
						cluster_map["s"]["y"] = cluster_map["s"]["y"] + 1
					} else if nth_rune(syllable, 0) == "ʒ" {
						cluster_map["ts"]["y"] = cluster_map["ts"]["y"] + 1
					}
					syllable = syllable[len(nth_rune(syllable, 0)):]
				} else {
					// no onset
					onset_map[""] = onset_map[""] + 1
				}

				/* Found a triple consonant? */
				if coda != "" && onset_if_cluster[1] != "" {
					if val, ok := valid_triple_consonants[coda][onset_if_cluster[0]][onset_if_cluster[1]]; ok {
						valid_triple_consonants[coda][onset_if_cluster[0]][onset_if_cluster[1]] = val + 1
					} else if _, ok := valid_triple_consonants[coda][onset_if_cluster[0]]; ok {
						valid_triple_consonants[coda][onset_if_cluster[0]][onset_if_cluster[1]] = 1
					} else if _, ok := valid_triple_consonants[coda]; ok {
						valid_triple_consonants[coda][onset_if_cluster[0]] = make(map[string]int)
						valid_triple_consonants[coda][onset_if_cluster[0]][onset_if_cluster[1]] = 1
					} else {
						valid_triple_consonants[coda] = make(map[string]map[string]int)
						valid_triple_consonants[coda][onset_if_cluster[0]] = make(map[string]int)
						valid_triple_consonants[coda][onset_if_cluster[0]][onset_if_cluster[1]] = 1
					}
				}
				//#    table_manager_supercluster(coda, start_cluster)
				//#    coda = ""syllable[l
				//#    start_cluster = ""

				/*
				 * Nucleus
				 */
				if len(syllable) > 1 && has("jw", nth_rune(syllable, 1)) {
					//diphthong
					nucleus_map[romanization[syllable[0:len(nth_rune(syllable, 0))+1]]] = nucleus_map[romanization[syllable[0:len(nth_rune(syllable, 0))+1]]] + 1
					syllable = string([]rune(syllable)[2:])
				} else if len(syllable) > 1 && has("lr", nth_rune(syllable, 0)) {
					nucleus_map[romanization[syllable[0:3]]] = nucleus_map[romanization[syllable[0:3]]] + 1
					continue
				} else {
					//vowel
					nucleus_map[romanization[nth_rune(syllable, 0)]] = nucleus_map[romanization[nth_rune(syllable, 0)]] + 1
					syllable = string([]rune(syllable)[1:])
				}

				/*
				 * Coda
				 */

				if len(syllable) == 0 || nth_rune(syllable, 0) == "s" {
					coda_map[""] = coda_map[""] + 1 //oìsss only
					coda = ""
				} else {
					if syllable == "k̚" {
						coda_map["k"] = coda_map["k"] + 1
						coda = "kx"
					} else if syllable == "p̚" {
						coda_map["p"] = coda_map["p"] + 1
						coda = "px"
					} else if syllable == "t̚" {
						coda_map["t"] = coda_map["t"] + 1
						coda = "tx"
					} else if syllable == "ʔ̚" {
						coda_map["'"] = coda_map["'"] + 1
						coda = "'"
					} else {
						if syllable[0] == 'k' && len(syllable) > 1 {
							coda_map["kx"] = coda_map["kx"] + 1
							coda = "kx"
						} else {
							coda_map[romanization[syllable]] = coda_map[romanization[syllable]] + 1
							coda = romanization[syllable]
						}
					}
				}
			}
		}
	}

	max_non_cluster = 0
	max_onset = 0
	max_nucleus = 0
	max_coda = 0

	// Copy everything from the maps to the arrays

	//Onsets
	for i := 0; i < len(onset_likelihood); i++ {
		onset_likelihood[i] = onset_map[onset_letters[i]]
		max_onset += onset_map[onset_letters[i]]
	}

	//Clusters
	max_non_cluster = max_onset

	super_i := 0
	for i := 0; i < len(cluster_1); i++ {
		for j := 0; j < len(cluster_2); j++ {
			cluster_letters[super_i] = cluster_1[i] + cluster_2[j]
			cluster_likelihood[super_i] = cluster_map[cluster_1[i]][cluster_2[j]]
			max_onset += cluster_map[cluster_1[i]][cluster_2[j]]
			super_i++
		}
	}

	//Nuclei
	for i := 0; i < len(nucleus_likelihood); i++ {
		nucleus_likelihood[i] = nucleus_map[nucleus_letters[i]]
		max_nucleus += nucleus_map[nucleus_letters[i]]
	}

	//Codas
	for i := 0; i < len(coda_likelihood); i++ {
		coda_likelihood[i] = coda_map[coda_letters[i]]
		max_coda += coda_map[coda_letters[i]]
	}
}
