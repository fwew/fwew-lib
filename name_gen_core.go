package fwew_lib

/*
 *	The core of the name generator
 *
 *	All functions here except PhonemeFrequency are helper functions for
 * 	name_gen.go.  PhonemeFrequency is called on startup for info for
 *	the name generator.
 */

import (
	"fmt"
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
var onsetLikelihood = [21]int{640, 498, 402, 398, 389, 359, 313, 311, 306, 196,
	193, 190, 188, 185, 180, 158, 140, 104, 83, 81, 76}
var onsetLetters = [21]string{"t", "", "n", "k", "l", "s", "'", "p", "r", "y",
	"ts", "m", "tx", "v", "w", "h", "ng", "z", "kx", "px", "f"}
var onsetMap = map[string]int{"t": 0, "": 0, "n": 0, "k": 0, "l": 0, "s": 0,
	"'": 0, "p": 0, "r": 0, "y": 0, "ts": 0, "m": 0, "tx": 0, "v": 0, "w": 0,
	"h": 0, "ng": 0, "z": 0, "kx": 0, "px": 0, "f": 0}

/* The clusters aren't arranged in this order, but they're only 1/8 of the onsets anyway.
 * Still waiting for a word with tspx. */
var clusterLikelihood = [39]int{
	19, 6, 15, 14, 17, 9, 31, 5, 29, 27, 20, 28, 66,
	14, 8, 26, 10, 32, 18, 15, 21, 41, 4, 49, 32, 57,
	5, 8, 11, 7, 1, 7, 2, 0, 16, 5, 9, 19, 64}
var clusterLetters = [39]string{
	"fk", "fkx", "fl", "fm", "fn", "fng", "fp", "fpx", "ft", "ftx", "fr", "fw", "fy",
	"sk", "skx", "sl", "sm", "sn", "sng", "sp", "spx", "st", "stx", "sr", "sw", "sy",
	"tsk", "tskx", "tsl", "tsm", "tsn", "tsng", "tsp", "tspx", "tst", "tstx", "tsr", "tsw", "tsy"}
var clusterMap = map[string]map[string]int{
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

var nucleusLikelihood = [14]int{1226, 1021, 760, 704, 615, 564, 277, 209, 187, 158, 153, 152, 70, 61}
var nucleusLetters = [14]string{"a", "e", "ì", "o", "u", "i", "ä", "aw", "ey", "ù", "rr", "ay", "ew", "ll"}
var nucleusMap = map[string]int{"a": 0, "e": 0, "ì": 0, "o": 0, "u": 0, "i": 0,
	"ä": 0, "aw": 0, "ey": 0, "ù": 0, "rr": 0, "ay": 0, "ew": 0, "ll": 0}

var codaLikelihood = [13]int{3938, 497, 405, 288, 258, 192, 179, 176, 133, 69, 48, 22, 18}
var codaLetters = [13]string{"", "n", "m", "ng", "l", "k", "p", "'", "r", "t", "kx", "px", "tx"}
var codaMap = map[string]int{"": 0, "n": 0, "m": 0, "ng": 0, "l": 0,
	"k": 0, "p": 0, "'": 0, "r": 0, "t": 0, "kx": 0, "px": 0, "tx": 0}

var validTripleConsonants = map[string]map[string]map[string]int{
	"'": {
		"f": {
			"y": 0,
		},
	},
}

var multiwordWords = map[string][][]string{}
var multiwordWordsLoose = map[string][][]string{}
var multiwordWordsReef = map[string][][]string{}

/* Calculated on startup to assist the random number generators and letter selector */
var maxOnset = 0
var maxNonCluster = 0
var maxNucleus = 0
var maxCoda = 0

/* Helper function to find the start of a string */
func firstRune(word string) (letter rune) {
	r := []rune(word)
	return r[0]
}

/* Get the nth to last letter of a string */
func getLastRune(word string, n int) (letter rune) {
	r := []rune(word)
	if n > len(r) {
		n = len(r)
	}
	return r[len(r)-n]
}

/* Take n letters off the end of a string */
func shaveRune(word string, n int) (letter string) {
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

	vowels := "aäeiìouù"

	for i, a := range runes {
		if i != 0 && i != len(runes)-1 && a == '\'' {
			if hasAt(vowels, output, i+1) && hasAt(vowels, output, i-1) {
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
func insertInfix(verb []string, infix string) (output string) {
	output = ""
	foundInfix := false
	for j := 0; j < len(verb); j++ {
		someVerb := []rune(verb[j])
		for k := 0; k < len(someVerb); k++ {
			if someVerb[k] == '.' {
				if !foundInfix {
					output += infix
					foundInfix = true
				}
			} else {
				output += string(someVerb[k])
			}
		}
		if j+1 < len(verb) {
			output += "-"
		}
	}
	return glottalCaps(output)
}

// Assistant function for name generating functions
func randIfZero(n int) (x int) {
	if n == 0 {
		return rand.Intn(3) + 2
	}
	return n
}

/* Is it a vowel? (for when the psuedovowel bool won't work) */
func isVowel(letter rune) (found bool) {
	// Also arranged from most to least common (not accounting for diphthongs)
	vowels := []rune{'a', 'e', 'u', 'ì', 'o', 'i', 'ä', 'ù'}
	// Linear search
	for _, a := range vowels {
		if letter == a {
			return true
		}
	}
	return false
}

/* Randomly select an onset for a Na'vi syllable */
func getOnset() (onset string, cluster bool) {
	selector := rand.Intn(maxOnset)
	// Clusters
	if selector > maxNonCluster { // If the number is too high for the non-cluster onsets,
		selector -= maxNonCluster // you get to skip all of them.  It saves time.
		// Linear search
		for i := 0; i < len(clusterLikelihood); i++ {
			if selector < clusterLikelihood[i] {
				return clusterLetters[i], true
			}
			selector -= clusterLikelihood[i]
		}
		return clusterLetters[len(clusterLetters)-1], true
	}
	// Non-clusters (single consonants)
	// Linear search
	for i := 0; i < len(onsetLikelihood); i++ {
		if selector < onsetLikelihood[i] {
			return onsetLetters[i], false
		}
		selector -= onsetLikelihood[i]
	}
	return onsetLetters[len(onsetLetters)-1], false

}

/* Get a random Na'vi nucleus */
func getNucleus() (onset string) {
	selector := rand.Intn(maxNucleus)
	// Linear search
	for i := 0; i < len(nucleusLikelihood); i++ {
		if selector < nucleusLikelihood[i] {
			return nucleusLetters[i]
		}
		selector -= nucleusLikelihood[i]
	}
	return nucleusLetters[len(nucleusLetters)-1]
}

/* Get a random Na'vi coda */
func getCoda() (onset string) {
	selector := rand.Intn(maxCoda)
	// Linear search
	for i := 0; i < len(codaLikelihood); i++ {
		if selector < codaLikelihood[i] {
			return codaLetters[i]
		}
		selector -= codaLikelihood[i]
	}
	return codaLetters[len(codaLetters)-1]
}

// Helper function for name-alu()
func oneWordVerb(verbList []Word) (words Word) {
	word := fastRandom(verbList)
	findVerb := strings.Split(word.InfixDots, " ")

	/* The second condition here is a clever and efficient little thing
	 * one word and not si: allowed (e.g. "takuk")
	 * two words and not si: disallowed (e.g. "tswìk kxenerit")
	 * one word and si: disallowed ("si" only)
	 * two words and si: allowed (e.g. "unil si")
	 * Any three-word verb: disallowed ("eltur tìtxen si" only)
	 * != is used as an exclusive "or"
	 */
	for (len(findVerb) == 2) != (findVerb[len(findVerb)-1] == "s..i") {
		word = fastRandom(verbList)
		findVerb = strings.Split(word.InfixDots, " ")
	}
	return word
}

/* Helper function: turn ejectives into voiced plosives for reef */
func reefPlosives(letter rune) (voiced rune) {
	if letter == 'p' {
		return 'b'
	} else if letter == 't' {
		return 'd'
	} else if letter == 'k' {
		return 'g'
	}
	return '' // How we know if it's an error
}

/* Helper function: Replace an ejective with a voiced plosive. */
func reefEjective(name string) (reefyName string) {
	onsetNew := ""
	lastThird := getLastRune(name, 3)

	if lastThird == 'x' { // Adjacent ejectives become adjacent voiced plosives, too
		onsetNew += string(reefPlosives(getLastRune(name, 4)))
	} else if lastThird == 'n' && getLastRune(name, 2) == 'k' {
		onsetNew += "-" // disambiguate on-gi vs o-ngi
	}

	onsetNew += string(reefPlosives(getLastRune(name, 2)))

	if lastThird == 'x' {
		return shaveRune(name, 4) + onsetNew
	}

	return shaveRune(name, 2) + onsetNew
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
	if strings.Contains(output, " or ") {
		output = strings.Split(output, " or ")[0]
	}
	output = strings.ReplaceAll(output, "_", "")
	return output
}

/* Randomly construct a phonotactically valid Na'vi word
 * Dialect codes: 0 = interdialect, 1 = forest, 2 = reef */
func singleNameGen(syllableCount int, dialect int) (name string) {
	phonoLock.Lock()
	defer phonoLock.Unlock()
	// Sometimes these things might be referenced across loops
	name = ""
	onset := ""
	nucleus := ""
	coda := ""
	psuedovowel := false // But not this
	cluster := false

	// Make a name with len syllables
	for i := 0; i < syllableCount; i++ {
		onset, cluster = getOnset()

		// Triple consonants are whitelisted
		if cluster && len(coda) > 0 { // don't want errors
			if !(coda == "t" && onset[0] == 's') { // t-s-kx is valid as ts-kx
				firstClusterNum := 1
				if onset[0] == 't' {
					firstClusterNum = 2
				}
				firstCluster := onset[:firstClusterNum]
				secondCluster := onset[firstClusterNum:]
				if _, ok := validTripleConsonants[coda][firstCluster][secondCluster]; ok {
					// Do nothing.  We found a valid triple
				} else {
					onset = secondCluster
				}
			}
		}

		nucleus = getNucleus()

		psuedovowel = false

		// These will be important later
		onsetlength := utf8.RuneCountInString(onset)
		namelength := utf8.RuneCountInString(name)

		// Get psuedovowel status
		if nucleus == "rr" || nucleus == "ll" {
			psuedovowel = true
			// Disallow onsets from imitating the psuedovowel
			if onsetlength > 0 {
				if getLastRune(onset, 1) == 'l' || getLastRune(onset, 1) == 'r' {
					onset = "'"
				}
				// If no onset, disallow the previous coda from imitating the psuedovowel
			} else if namelength > 0 {
				if getLastRune(name, 1) == firstRune(nucleus) || isVowel(getLastRune(name, 1)) {
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
			if onsetlength == 0 && namelength > 0 && getLastRune(name, 1) == firstRune(nucleus) {
				onset = "y"
			}
		} else if nucleusMap["ù"] == 0 { //no psuedovowel or forest dialect
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
			if firstRune(onset) == getLastRune(name, length) {
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
			coda = getCoda()
		} else {
			coda = ""
		}

		// reef dialect stuff
		if dialect == 2 && namelength > 1 { // In reef dialect,
			if getLastRune(name, 1) == 'x' { // if there's an ejective in the onset
				if namelength > 2 {
					// that's not in a cluster,
					lastRune := getLastRune(name, 3)
					if !(lastRune == 's' || lastRune == 'f') {
						// it becomes a voiced plosive
						name = reefEjective(name)
					}
				} else {
					name = reefEjective(name)
				}
			} else if !psuedovowel && getLastRune(name, 1) == '\'' && getLastRune(name, 2) != firstRune(nucleus) {
				// 'a'aw is optionally 'aaw (the generator leaves it in)
				if isVowel(getLastRune(name, 2)) { // Does `kaw'it` become `kawit` in reef?
					name = shaveRune(name, 1)
				}
			}
		}

		// Finish the syllable
		name += nucleus + coda
	}
	return name
}

func glottalCaps(input string) (output string) {
	a := []rune(input)
	n := 1
	output = ""
	if a[0] == '\'' {
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

func fastRandom(wordList []Word) (results Word) {
	dictLength := len(wordList)

	return wordList[rand.Intn(dictLength)]
}

// What is the nth rune of word?
func nthRune(word string, n int) string {
	i := 0
	for _, r := range word {
		if i == n {
			return string(r)
		}
		i += 1
	}

	return ""
}

// Does ipa contain any character from word as its nth letter?
func hasAt(word string, ipa string, n int) (output bool) {
	// negative index
	if n < 0 {
		n = len([]rune(ipa)) + n
	}

	i := 0
	for _, s := range ipa {
		if i == n {
			for _, r := range word {
				if r == s {
					return true
				}
			}
			break // save a few compute cycles
		}
		i += 1
	}

	return false
}

// SortedWords is a helper function for name-alu
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

// PhonemeDistros is called on startup to feed and compile dictionary information into the name generator
func PhonemeDistros() {
	phonoLock.Lock()
	defer phonoLock.Unlock()
	// get the dict
	words, err := List([]string{}, 0)

	clear(multiwordWords)
	clear(multiwordWordsLoose)
	clear(multiwordWordsReef)

	if err != nil || len(words) == 0 {
		return
	}

	//set the maps to zero

	//Onsets
	for i := 0; i < len(onsetLetters); i++ {
		onsetMap[onsetLetters[i]] = 0
	}

	//Clusters
	cluster1Full := []string{"f", "s", "ts"}
	cluster2Full := []string{"k", "kx", "l", "m", "n", "ng", "p",
		"px", "t", "tx", "r", "w", "y"}
	for i := 0; i < len(cluster1Full); i++ {
		for j := 0; j < len(cluster2Full); j++ {
			clusterMap[cluster1Full[i]][cluster2Full[j]] = 0
		}
	}

	//Nuclei
	for i := 0; i < len(nucleusLikelihood); i++ {
		nucleusMap[nucleusLetters[i]] = 0
	}

	//Codas
	for i := 0; i < len(codaLikelihood); i++ {
		codaMap[codaLetters[i]] = 0
	}

	//syllable_map := map[string]int{}

	// Look through all the words
	for i := 0; i < len(words); i++ {
		word := strings.Split(words[i].IPA, " ")

		// Piggybacking off of the frequency script to get all words with spaces
		allWords := strings.Split(strings.ToLower(words[i].Navi), " ")
		if len(allWords) > 1 {
			newWords := dialectCrunch(allWords, true, false)
			newWordsReef := dialectCrunch(allWords, true, true)
			if _, ok := multiwordWordsLoose[newWords[0]]; ok {
				// Ensure no duplicates
				appended := false

				// Append in a way that makes the longer words first
				var temp [][]string
				for _, j := range multiwordWordsLoose[newWords[0]] {
					if !appended && len([]rune(newWords[1])) > len([]rune(j[0])) {
						temp = append(temp, newWords[1:])
						appended = true
					}
					temp = append(temp, j)
				}
				if len(temp) <= len(multiwordWordsLoose[newWords[0]]) {
					temp = append(temp, newWords[1:])
				}

				multiwordWordsLoose[newWords[0]] = temp
			} else {
				multiwordWordsLoose[newWords[0]] = [][]string{newWords[1:]}
			}

			if _, ok := multiwordWordsReef[newWordsReef[0]]; ok {
				// Ensure no duplicates
				appended := false

				// Append in a way that makes the longer words first
				var temp [][]string
				for _, j := range multiwordWordsReef[newWordsReef[0]] {
					if !appended && len([]rune(newWordsReef[1])) > len([]rune(j[0])) {
						temp = append(temp, newWordsReef[1:])
						appended = true
					}
					temp = append(temp, j)
				}
				if len(temp) <= len(multiwordWordsReef[newWordsReef[0]]) {
					temp = append(temp, newWordsReef[1:])
				}

				multiwordWordsReef[newWordsReef[0]] = temp
			} else {
				multiwordWordsReef[newWordsReef[0]] = [][]string{newWordsReef[1:]}
			}

			if _, ok := multiwordWords[allWords[0]]; ok {
				// Ensure no duplicates
				appended := false

				// Append in a way that makes the longer words first
				var temp [][]string
				for _, j := range multiwordWords[allWords[0]] {
					if !appended && len([]rune(allWords[1])) > len([]rune(j[0])) {
						temp = append(temp, allWords[1:])
						appended = true
					}
					temp = append(temp, j)
				}
				if len(temp) <= len(multiwordWords[allWords[0]]) {
					temp = append(temp, allWords[1:])
				}

				multiwordWords[allWords[0]] = temp
			} else {
				multiwordWords[allWords[0]] = [][]string{allWords[1:]}
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

				onsetIfCluster := [2]string{"", ""}

				//roman_syllable := ""

				// ts
				if len(syllable) >= 4 && syllable[0:4] == "t͡s" {
					onsetIfCluster[0] = "ts"
					//tsp
					if hasAt("ptk", syllable, 3) {
						if nthRune(syllable, 4) == "'" {
							// ts + ejective onset
							clusterMap["ts"][romanization[syllable[4:6]]] = clusterMap["ts"][romanization[syllable[4:6]]] + 1
							onsetIfCluster[1] = romanization[syllable[4:6]]
							//roman_syllable += "ts" + romanization[syllable[4:6]]
							syllable = syllable[6:]
						} else {
							// ts + unvoiced plosive
							clusterMap["ts"][romanization[string(syllable[4])]] = clusterMap["ts"][romanization[string(syllable[4])]] + 1
							onsetIfCluster[1] = romanization[string(syllable[4])]
							//roman_syllable += "ts" + romanization[string(syllable[4])]
							syllable = syllable[5:]
						}
					} else if hasAt("lɾmnŋwj", syllable, 3) {
						// ts + other consonent
						clusterMap["ts"][romanization[nthRune(syllable, 3)]] = clusterMap["ts"][romanization[nthRune(syllable, 3)]] + 1
						onsetIfCluster[1] = romanization[nthRune(syllable, 3)]
						//roman_syllable += "ts" + romanization[nth_rune(syllable, 3)]
						syllable = syllable[4+len(nthRune(syllable, 3)):]
					} else {
						// ts without a cluster
						onsetMap["ts"] = onsetMap["ts"] + 1
						//roman_syllable += "ts"
						syllable = syllable[4:]
					}
				} else if hasAt("fs", syllable, 0) {
					//
					onsetIfCluster[0] = string(syllable[0])
					if hasAt("ptk", syllable, 1) {
						if nthRune(syllable, 2) == "'" {
							// f/s + ejective onset
							clusterMap[string(syllable[0])][romanization[syllable[1:3]]] = clusterMap[string(syllable[0])][romanization[syllable[1:3]]] + 1
							onsetIfCluster[1] = romanization[syllable[1:3]]
							//roman_syllable += string(syllable[0]) + romanization[syllable[1:3]]
							syllable = syllable[3:]
						} else {
							// f/s + unvoiced plosive
							clusterMap[string(syllable[0])][romanization[string(syllable[1])]] = clusterMap[string(syllable[0])][romanization[string(syllable[1])]] + 1
							onsetIfCluster[1] = romanization[string(syllable[1])]
							//roman_syllable += string(syllable[0]) + romanization[string(syllable[1])]
							syllable = syllable[2:]
						}
					} else if hasAt("lɾmnŋwj", syllable, 1) {
						// f/s + other consonent
						clusterMap[string(syllable[0])][romanization[nthRune(syllable, 1)]] = clusterMap[string(syllable[0])][romanization[nthRune(syllable, 1)]] + 1
						onsetIfCluster[1] = romanization[nthRune(syllable, 1)]
						//roman_syllable += string(syllable[0]) + romanization[nth_rune(syllable, 1)]
						syllable = syllable[1+len(nthRune(syllable, 1)):]
					} else {
						// f/s without a cluster
						onsetMap[string(syllable[0])] = onsetMap[string(syllable[0])] + 1
						//roman_syllable += string(syllable[0])
						syllable = syllable[1:]
					}
				} else if hasAt("ptk", syllable, 0) {
					if nthRune(syllable, 1) == "'" {
						// ejective
						onsetMap[romanization[syllable[0:2]]] = onsetMap[romanization[syllable[0:2]]] + 1
						//roman_syllable += romanization[syllable[0:2]]
						syllable = syllable[2:]
					} else {
						// unvoiced plosive
						onsetMap[romanization[string(syllable[0])]] = onsetMap[romanization[string(syllable[0])]] + 1
						//roman_syllable += romanization[string(syllable[0])]
						syllable = syllable[1:]
					}
				} else if hasAt("ʔlɾhmnŋvwjzbdg", syllable, 0) {
					// other normal onset
					onsetMap[romanization[nthRune(syllable, 0)]] = onsetMap[romanization[nthRune(syllable, 0)]] + 1
					//roman_syllable += romanization[nth_rune(syllable, 0)]
					syllable = syllable[len(nthRune(syllable, 0)):]
				} else if hasAt("ʃʒ", syllable, 0) {
					// one sound representd as a cluster
					if nthRune(syllable, 0) == "ʃ" {
						clusterMap["s"]["y"] = clusterMap["s"]["y"] + 1
						//roman_syllable += "sy"
					} else if nthRune(syllable, 0) == "ʒ" {
						clusterMap["ts"]["y"] = clusterMap["ts"]["y"] + 1
						//roman_syllable += "tsy"
					}
					syllable = syllable[len(nthRune(syllable, 0)):]
				} else {
					// no onset
					onsetMap[""] = onsetMap[""] + 1
				}

				/* Found a triple consonant? */
				if coda != "" && onsetIfCluster[1] != "" {
					if val, ok := validTripleConsonants[coda][onsetIfCluster[0]][onsetIfCluster[1]]; ok {
						validTripleConsonants[coda][onsetIfCluster[0]][onsetIfCluster[1]] = val + 1
					} else if _, ok := validTripleConsonants[coda][onsetIfCluster[0]]; ok {
						validTripleConsonants[coda][onsetIfCluster[0]][onsetIfCluster[1]] = 1
					} else if _, ok := validTripleConsonants[coda]; ok {
						validTripleConsonants[coda][onsetIfCluster[0]] = make(map[string]int)
						validTripleConsonants[coda][onsetIfCluster[0]][onsetIfCluster[1]] = 1
					} else {
						validTripleConsonants[coda] = make(map[string]map[string]int)
						validTripleConsonants[coda][onsetIfCluster[0]] = make(map[string]int)
						validTripleConsonants[coda][onsetIfCluster[0]][onsetIfCluster[1]] = 1
					}
				}
				//#    table_manager_supercluster(coda, start_cluster)
				//#    coda = ""syllable[l
				//#    start_cluster = ""

				/*
				 * Nucleus
				 */
				if len(syllable) > 1 && hasAt("jw", syllable, 1) {
					//diphthong
					nucleusMap[romanization[syllable[0:len(nthRune(syllable, 0))+1]]] = nucleusMap[romanization[syllable[0:len(nthRune(syllable, 0))+1]]] + 1
					//roman_syllable += romanization[syllable[0:len(nth_rune(syllable, 0))+1]]
					syllable = string([]rune(syllable)[2:])
				} else if len(syllable) > 1 && hasAt("lr", syllable, 0) {
					nucleusMap[romanization[syllable[0:3]]] = nucleusMap[romanization[syllable[0:3]]] + 1
					//roman_syllable += romanization[syllable[0:3]]
					continue
				} else {
					//vowel
					nucleusMap[romanization[nthRune(syllable, 0)]] = nucleusMap[romanization[nthRune(syllable, 0)]] + 1
					//roman_syllable += romanization[nth_rune(syllable, 0)]
					if len(syllable) == 0 {
						fmt.Println("Invalid word: " + words[i].ID + " - " + words[i].Navi + " - " + words[i].IPA)
					} else {
						syllable = string([]rune(syllable)[1:])
					}
				}

				/*
				 * Coda
				 */

				if len(syllable) == 0 || nthRune(syllable, 0) == "s" {
					codaMap[""] = codaMap[""] + 1 //oìsss only
					coda = ""
				} else {
					if syllable == "k̚" {
						codaMap["k"] = codaMap["k"] + 1
						coda = "k"
					} else if syllable == "p̚" {
						codaMap["p"] = codaMap["p"] + 1
						coda = "p"
					} else if syllable == "t̚" {
						codaMap["t"] = codaMap["t"] + 1
						coda = "t"
					} else if syllable == "ʔ̚" {
						codaMap["'"] = codaMap["'"] + 1
						coda = "'"
					} else {
						if syllable[0] == 'k' && len(syllable) > 1 {
							codaMap["kx"] = codaMap["kx"] + 1
							coda = "kx"
						} else {
							codaMap[romanization[syllable]] = codaMap[romanization[syllable]] + 1
							coda = romanization[syllable]
						}
					}
				}
				/*roman_syllable += coda

				// Finally see if there is a good syllable frequency here
				if _, ok := syllable_map[roman_syllable]; !ok {
					syllable_map[roman_syllable] = 1
				} else {
					syllable_map[roman_syllable] = syllable_map[roman_syllable] + 1
				}*/
			}
		}
	}

	// Show the phoneme map sorted
	/*syllable_tuples := []PhonemeTuple{}
	for key, val := range syllable_map {
		syllable_tuples = append(syllable_tuples, PhonemeTuple{val, key})
	}
	sort.Sort(Tuples(syllable_tuples))

	for _, a := range syllable_tuples {
		fmt.Println(a)
	}*/

	maxNonCluster = 0
	maxOnset = 0
	maxNucleus = 0
	maxCoda = 0

	// Copy everything from the maps to the arrays

	//Onsets
	for i := 0; i < len(onsetLikelihood); i++ {
		onsetLikelihood[i] = onsetMap[onsetLetters[i]]
		maxOnset += onsetMap[onsetLetters[i]]
	}

	//Clusters
	maxNonCluster = maxOnset

	superI := 0
	for i := 0; i < len(cluster1Full); i++ {
		for j := 0; j < len(cluster2Full); j++ {
			clusterLetters[superI] = cluster1Full[i] + cluster2Full[j]
			clusterLikelihood[superI] = clusterMap[cluster1Full[i]][cluster2Full[j]]
			maxOnset += clusterMap[cluster1Full[i]][cluster2Full[j]]
			superI++
		}
	}

	//Nuclei
	for i := 0; i < len(nucleusLikelihood); i++ {
		nucleusLikelihood[i] = nucleusMap[nucleusLetters[i]]
		maxNucleus += nucleusMap[nucleusLetters[i]]
	}

	//Codas
	for i := 0; i < len(codaLikelihood); i++ {
		codaLikelihood[i] = codaMap[codaLetters[i]]
		maxCoda += codaMap[codaLetters[i]]
	}
}
