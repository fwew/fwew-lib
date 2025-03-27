package fwew_lib

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

const dictFileName = "dictionary-v2.txt"

var dictionary []Word
var dictHashLoose map[string][]Word
var dictHashStrict map[string][]Word
var dictHashStrictReef map[string][]Word
var dictionaryCached bool
var dictHashCached bool
var dictHash2 MetaDict
var dictHash2Parenthesis MetaDict
var dictHash2Cached bool
var homonyms string
var oddballs string
var multiIPA string

type MetaDict struct {
	EN map[string][]string
	DE map[string][]string
	ES map[string][]string
	ET map[string][]string
	FR map[string][]string
	HU map[string][]string
	KO map[string][]string
	NL map[string][]string
	PL map[string][]string
	PT map[string][]string
	RU map[string][]string
	SV map[string][]string
	TR map[string][]string
	UK map[string][]string
}

var letterMap = map[rune]int{
	' ': -1, '\'': 0, 'a': 1, '2': 2, '3': 3,
	'ä': 4, 'e': 5, '4': 6, '5': 7,
	'f': 8, 'h': 9, 'i': 10, 'ì': 11,
	'j': 12, 'k': 13, 'q': 14, 'l': 15,
	'1': 16, 'm': 17, 'n': 18, 'g': 19,
	'o': 20, 'p': 21, 'b': 22, 'r': 23,
	'0': 24, 's': 25, 't': 26, 'c': 27,
	'd': 28, 'u': 29, 'v': 30, 'w': 31,
	'y': 32, 'z': 33, '-': 34,
}

var nkx = []string{}
var nkxSub = map[string]string{}

// A mutex to ensure concurrent requests to the
// dictionary and phoneme counts will not cause
// the program to crash
var universalLock sync.Mutex
var phonoLock sync.Mutex

// helper for nkx for shortest words first
func shortestFirst(array []string, input string) []string {
	newArray := []string{}
	found := false
	for _, a := range array {
		if !found && len(a) > len(input) {
			newArray = append(newArray, input)
		}
		newArray = append(newArray, a)
	}
	if !found {
		newArray = append(newArray, input)
	}
	return newArray
}

// check if a file exists
func fileExists(filepath string) bool {
	fileStat, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}

	return !fileStat.IsDir()
}

// the dictionary file can be places into:
// - <workingDir>/dictionary.txt
// - <workingDir>/.fwew/dictionary.txt
// - <homeDir>/.fwew/dictionary.txt
func FindDictionaryFile() string {
	wd, err := os.Getwd()
	if err == nil {
		path := filepath.Join(wd, ".fwew", dictFileName)
		if fileExists(path) {
			return path
		}

		path = filepath.Join(wd, dictFileName)
		if fileExists(path) {
			return path
		}
	}

	path := filepath.Join(texts["dataDir"], dictFileName)
	if fileExists(path) {
		return path
	}

	return ""
}

func AlphabetizeHelper(a string, b string) bool {
	aCompacted := []rune(strings.ReplaceAll(compress(strings.ToLower(a)), "-", ""))

	// Start in the middle
	bCompacted := []rune(strings.ReplaceAll(compress(strings.ToLower(b)), "-", ""))
	lowestLen := len(aCompacted)
	if lowestLen > len(bCompacted) {
		lowestLen = len(bCompacted)
	}
	// compare an individual word
	for j := 0; j < lowestLen; j++ {
		// If the new letter is bigger, wait until it gets
		if letterMap[aCompacted[j]] < letterMap[bCompacted[j]] {
			return true
		} else if letterMap[aCompacted[j]] > letterMap[bCompacted[j]] {
			return false
		}
		// If equal, continue
	}

	// longer words go after
	return len(aCompacted) < len(bCompacted)
}

func AppendAndAlphabetize(words []Word, word Word) []Word {
	// Ensure it's not a duplicate
	for _, a := range words {
		if word.ID == a.ID {
			if len(word.Affixes.Prefix) == len(a.Affixes.Prefix) &&
				len(word.Affixes.Suffix) == len(a.Affixes.Suffix) &&
				len(word.Affixes.Lenition) == len(a.Affixes.Lenition) &&
				len(word.Affixes.Infix) == len(a.Affixes.Infix) {
				return words
			}
		}
	}
	// new array
	switch len(words) {
	case 0:
		return []Word{word}
	case 1:
		var newWords = []Word{}
		if AlphabetizeHelper(words[0].Syllables, word.Syllables) {
			newWords = []Word{words[0], word}
		} else {
			newWords = []Word{word, words[0]}
		}
		return newWords
	case 2:
		var newWords = []Word{}
		if AlphabetizeHelper(word.Syllables, words[0].Syllables) {
			newWords = []Word{word, words[0], words[1]}
		} else if AlphabetizeHelper(words[1].Syllables, word.Syllables) {
			newWords = []Word{words[0], words[1], word}
		} else {
			newWords = []Word{words[0], word, words[1]}
		}
		return newWords
	}

	// start in the middle
	halfway := len(words) / 2

	// Copy the first half
	newWords := make([]Word, len(words[:halfway]))
	copy(newWords, words[:halfway])

	// Copy the second half
	oldWords := make([]Word, len(words[halfway:]))
	copy(oldWords, words[halfway:])

	// compare an individual word
	if AlphabetizeHelper(word.Syllables, words[halfway].Syllables) {
		// Copy the first half
		newWords = AppendAndAlphabetize(newWords, word)

		// Join them
		newWords = append(newWords, oldWords...)
	} else {
		// Copy the second half
		oldWords = AppendAndAlphabetize(oldWords, word)

		// Join them
		newWords = append(newWords, oldWords...)
	}

	return newWords
}

// Helper to find empty definitions
func NullDef(definition string) bool {
	return strings.ToUpper(definition) == "NULL" || len(strings.Trim(definition, " ")) < 1
}

// If a definition is not available in a certain language, default to English
func EnglishIfNull(word Word) Word {
	// English
	if NullDef(word.EN) {
		word.EN = "(no definition)"
	}

	// German (Deutsch)
	if NullDef(word.DE) {
		word.DE = word.EN
	}

	// Spanish (Español)
	if NullDef(word.ES) {
		word.ES = word.EN
	}

	// Estonian (Eesti)
	if NullDef(word.ET) {
		word.ET = word.EN
	}

	// French (Français)
	if NullDef(word.FR) {
		word.FR = word.EN
	}

	// Hungarian (Magyar)
	if NullDef(word.HU) {
		word.HU = word.EN
	}

	// Korean (한국어)
	if NullDef(word.KO) {
		word.KO = word.EN
	}

	// Dutch (Nederlands)
	if NullDef(word.NL) {
		word.NL = word.EN
	}

	// Polish (Polski)
	if NullDef(word.PL) {
		word.PL = word.EN
	}

	// Portuguese (Português)
	if NullDef(word.PT) {
		word.PT = word.EN
	}

	// Russian (Русский)
	if NullDef(word.RU) {
		word.RU = word.EN
	}

	// Swedish (Svenska)
	if NullDef(word.SV) {
		word.SV = word.EN
	}

	// Turkish (Türkçe)
	if NullDef(word.TR) {
		word.TR = word.EN
	}

	// Ukrainian (Українська)
	if NullDef(word.UK) {
		word.UK = word.EN
	}

	return word
}

// Helper function to get phonetic transcriptions of secondary pronunciations
// Only multiple IPA words will call this function
func RomanizeSecondIPA(IPA string) string {
	// now Romanize the IPA
	IPA = strings.ReplaceAll(IPA, "ʊ", "u")
	IPA = strings.ReplaceAll(IPA, "õ", "o") // vonvä' as võvä' only
	word := strings.Split(IPA, " ")

	breakdown := ""

	// get the last one only
	for j := 2; j < len(word); j++ {
		word[j] = strings.ReplaceAll(word[j], "[", "")
		// "or" means there's more than one IPA in this word, and we only want one
		if word[j] == "or" {
			breakdown = ""
			continue
		}

		syllables := strings.Split(word[j], ".")

		/* Onset */
		for k := 0; k < len(syllables); k++ {
			syllable := strings.ReplaceAll(syllables[k], "·", "")
			syllable = strings.ReplaceAll(syllable, "ˈ", "")
			syllable = strings.ReplaceAll(syllable, "ˌ", "")

			// tsy
			if strings.HasPrefix(syllable, "tʃ") {
				breakdown += "ch"
				syllable = strings.TrimPrefix(syllable, "tʃ")
			} else if len(syllable) >= 4 && syllable[0:4] == "t͡s" {
				// ts
				breakdown += "ts"
				//tsp
				if hasAt("ptk", syllable, 3) {
					if nth_rune(syllable, 4) == "'" {
						// ts + ejective onset
						breakdown += romanization2[syllable[4:6]]
						syllable = syllable[6:]
					} else {
						// ts + unvoiced plosive
						breakdown += romanization2[string(syllable[4])]
						syllable = syllable[5:]
					}
				} else if hasAt("lɾmnŋwj", syllable, 3) {
					// ts + other consonent
					breakdown += romanization2[nth_rune(syllable, 3)]
					syllable = syllable[4+len(nth_rune(syllable, 3)):]
				} else {
					// ts without a cluster
					syllable = syllable[4:]
				}
			} else if hasAt("fs", syllable, 0) {
				//
				breakdown += nth_rune(syllable, 0)
				if hasAt("ptk", syllable, 1) {
					if nth_rune(syllable, 2) == "'" {
						// f/s + ejective onset
						breakdown += romanization2[syllable[1:3]]
						syllable = syllable[3:]
					} else {
						// f/s + unvoiced plosive
						breakdown += romanization2[string(syllable[1])]
						syllable = syllable[2:]
					}
				} else if hasAt("lɾmnŋwj", syllable, 1) {
					// f/s + other consonent
					breakdown += romanization2[nth_rune(syllable, 1)]
					syllable = syllable[1+len(nth_rune(syllable, 1)):]
				} else {
					// f/s without a cluster
					syllable = syllable[1:]
				}
			} else if hasAt("ptk", syllable, 0) {
				if nth_rune(syllable, 1) == "'" {
					// ejective
					breakdown += romanization2[syllable[0:2]]
					syllable = syllable[2:]
				} else {
					// unvoiced plosive
					breakdown += romanization2[string(syllable[0])]
					syllable = syllable[1:]
				}
			} else if hasAt("ʔlɾhmnŋvwjzbdg", syllable, 0) {
				// other normal onset
				breakdown += romanization2[nth_rune(syllable, 0)]
				syllable = syllable[len(nth_rune(syllable, 0)):]
			} else if hasAt("ʃʒ", syllable, 0) {
				// one sound representd as a cluster
				if nth_rune(syllable, 0) == "ʃ" {
					breakdown += "sh"
				}
				syllable = syllable[len(nth_rune(syllable, 0)):]
			}

			/*
			 * Nucleus
			 */
			if len(syllable) > 1 && hasAt("jw", syllable, 1) {
				//diphthong
				breakdown += romanization2[syllable[0:len(nth_rune(syllable, 0))+1]]
				syllable = string([]rune(syllable)[2:])
			} else if len(syllable) > 1 && hasAt("lr", syllable, 0) {
				breakdown += romanization2[syllable[0:3]]
				continue
			} else {
				//vowel
				breakdown += romanization2[nth_rune(syllable, 0)]
				syllable = string([]rune(syllable)[1:])
			}

			/*
			 * Coda
			 */
			if len(syllable) > 0 {
				if nth_rune(syllable, 0) == "s" {
					breakdown += "sss" //oìsss only
				} else {
					if syllable == "k̚" {
						breakdown += "k"
					} else if syllable == "p̚" {
						breakdown += "p"
					} else if syllable == "t̚" {
						breakdown += "t"
					} else if syllable == "ʔ̚" {
						breakdown += "'"
					} else {
						if syllable[0] == 'k' && len(syllable) > 1 {
							breakdown += "kx"
						} else {
							breakdown += romanization2[syllable]
						}
					}
				}
			}
		}
		breakdown += " "
	}
	return strings.TrimSuffix(breakdown, " ")
}

func UncacheDict() {
	dictionaryCached = false
	dictionary = []Word{}
}

func CacheDict() error {
	var err error

	UncacheDict()
	err = runOnDB(func(word Word) error {
		dictionary = append(dictionary, word)
		return nil
	})

	if err == nil {
		fmt.Println("cache 0 loaded (SQL)")
	} else {
		UncacheDict()
		err = runOnFile(func(word Word) error {
			dictionary = append(dictionary, word)
			return nil
		})
		//fmt.Println("cache 0 loaded (File)")
	}

	if err != nil {
		UncacheDict()
		return err
	}

	dictionaryCached = true

	return nil
}

func CacheDictHash() error {
	err := CacheDictHashOrig(true)
	if err == nil {
		fmt.Println("cache 1 loaded (SQL)")
	} else {
		err = CacheDictHashOrig(false)
		//fmt.Println("cache 1 loaded (File)")
	}
	return err
}

// This will cache the whole dictionary (Na'vi to natural language).
// Please call this, if you want to translate multiple words or running infinitely (e.g. CLI-go-prompt, discord-bot)
func CacheDictHashOrig(mysql bool) error {
	// dont run if already is cached
	if len(dictHashLoose) != 0 {
		return nil
	} else {
		dictHashLoose = make(map[string][]Word)
		dictHashStrict = make(map[string][]Word)
		dictHashStrictReef = make(map[string][]Word)
	}

	tempHoms := []string{}

	//Clear to avoid duplicates
	multiIPA = ""

	var f = func(word Word) error {
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

		// Make sure we know of every word with nkx
		if strings.Contains(standardizedWord, "nkx") {
			fakeNG := strings.ReplaceAll(standardizedWord, "nkx", "ng")
			nkx = shortestFirst(nkx, fakeNG)
			nkxSub[fakeNG] = standardizedWord
		}

		standardizedWordArray := dialectCrunch(strings.Split(standardizedWord, " "), true, false, true)
		standardizedWordLoose := ""
		for i, a := range standardizedWordArray {
			if i != 0 {
				standardizedWordLoose += " "
			}
			standardizedWordLoose += a
		}

		strictReefArray := dialectCrunch(strings.Split(standardizedWord, " "), true, true, true)
		strictReef := ""
		for i, a := range strictReefArray {
			if i != 0 {
				strictReef += " "
			}
			strictReef += a
		}

		// If the word appears more than once, record it
		if _, ok := dictHashStrict[standardizedWord]; ok {
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

		word = EnglishIfNull(word)
		dictHashLoose[standardizedWordLoose] = append(dictHashLoose[standardizedWordLoose], word)
		dictHashStrictReef[strictReef] = append(dictHashStrictReef[strictReef], word)
		dictHashStrict[standardizedWord] = append(dictHashStrict[standardizedWord], word)

		//find words with multiple IPAs
		if strings.Contains(word.IPA, " or ") {
			multiIPA += word.Navi + " "
			secondTerm := RomanizeSecondIPA(word.IPA)
			if secondTerm != standardizedWord {
				dictHashLoose[dialectCrunch([]string{secondTerm}, true, false, true)[0]] = append(dictHashLoose[dialectCrunch([]string{secondTerm}, true, false, true)[0]], word)
				dictHashStrictReef[dialectCrunch([]string{secondTerm}, true, true, true)[0]] = append(dictHashStrictReef[dialectCrunch([]string{secondTerm}, true, true, true)[0]], word)
				dictHashStrict[secondTerm] = append(dictHashStrict[secondTerm], word)
			}
		}

		// See whether or not it violates normal phonotactic rules like Jakesully or Oìsss
		valid := true
		for _, a := range strings.Split(IsValidNavi(word.Navi, "en", false), "\n") {
			// Check every word.  If one of them isn't good, write down the word
			if len(a) > 0 && (!strings.Contains(a, "Valid:") || strings.Contains(a, "reef")) {
				valid = false
				break
			}
		}
		if !valid {
			oddballs += word.Navi + " "
		}

		return nil
	}

	var err error
	if mysql {
		err = runOnDB(f)
		if err != nil {
			UncacheHashDict()
			return err
		}
	} else {
		err = runOnFile(f)
		if err != nil {
			log.Printf("Error caching dictionary: %s", err)
			UncacheHashDict()
			return err
		}
	}

	// Reverse the order to make accidental and new homonyms easier to see
	// Also make it a string for easier searching
	i := len(tempHoms)
	for i > 0 {
		i--
		homonyms += tempHoms[i] + " "
	}

	homonyms = strings.TrimSuffix(homonyms, " ")

	dictHashCached = true

	return nil
}

// Turn a definition into its searchable terms
func SearchTerms(input string, excludeParen bool) []string {
	badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\"„“”«»`

	input = strings.ReplaceAll(input, "(", " (")

	// remove anything in parenthesis to avoid clogging search results
	tempString := ""
	parenthesis := false
	for _, c := range input {
		if excludeParen {
			if c == '(' {
				parenthesis = true
			} else if c == ')' {
				parenthesis = false
				continue
			}
		}

		if !parenthesis {
			tempString += string(c)
		}
	}
	input = tempString

	// remove all the sketchy chars from arguments
	for _, c := range badChars {
		input = strings.ReplaceAll(input, string(c), "")
	}

	// normalize tìftang character
	input = strings.ReplaceAll(input, "’", "'")
	input = strings.ReplaceAll(input, "‘", "'")

	// find everything lowercase
	input = strings.ToLower(input)

	return strings.Split(input, " ")
}

// Helper function for CacheDictHash2
func AssignWord(wordmap map[string][]string, natlangWords string, naviWord string, excludeParen bool) (result map[string][]string) {
	newWords := SearchTerms(natlangWords, excludeParen)

	for i := 0; i < len(newWords); i++ {
		duplicate := false
		for j := 0; j < len(wordmap[newWords[i]]); j++ {
			if wordmap[newWords[i]][j] == naviWord {
				duplicate = true
				break
			}
		}
		if !duplicate {
			wordmap[newWords[i]] = append(wordmap[newWords[i]], naviWord)
		}
	}
	return wordmap
}

// Natural languages to Na'vi
func CacheDictHash2() error {
	err := CacheDictHash2Orig(true)
	if err == nil {
		fmt.Println("cache 2 loaded (SQL)")
	} else {
		err = CacheDictHash2Orig(false)
		//fmt.Println("cache 2 loaded (File)")
	}
	return err
}

func CacheDictHash2Orig(mysql bool) error {
	// dont run if already is cached
	if len(dictHash2.EN) != 0 {
		return nil
	} else {
		dictHash2.EN = make(map[string][]string)
		dictHash2.DE = make(map[string][]string)
		dictHash2.ES = make(map[string][]string)
		dictHash2.ET = make(map[string][]string)
		dictHash2.FR = make(map[string][]string)
		dictHash2.HU = make(map[string][]string)
		dictHash2.KO = make(map[string][]string)
		dictHash2.NL = make(map[string][]string)
		dictHash2.PL = make(map[string][]string)
		dictHash2.PT = make(map[string][]string)
		dictHash2.RU = make(map[string][]string)
		dictHash2.SV = make(map[string][]string)
		dictHash2.TR = make(map[string][]string)
		dictHash2.UK = make(map[string][]string)

		dictHash2Parenthesis.EN = make(map[string][]string)
		dictHash2Parenthesis.DE = make(map[string][]string)
		dictHash2Parenthesis.ES = make(map[string][]string)
		dictHash2Parenthesis.ET = make(map[string][]string)
		dictHash2Parenthesis.FR = make(map[string][]string)
		dictHash2Parenthesis.HU = make(map[string][]string)
		dictHash2Parenthesis.KO = make(map[string][]string)
		dictHash2Parenthesis.NL = make(map[string][]string)
		dictHash2Parenthesis.PL = make(map[string][]string)
		dictHash2Parenthesis.PT = make(map[string][]string)
		dictHash2Parenthesis.RU = make(map[string][]string)
		dictHash2Parenthesis.SV = make(map[string][]string)
		dictHash2Parenthesis.TR = make(map[string][]string)
		dictHash2Parenthesis.UK = make(map[string][]string)
	}

	// Set up the whole thing

	var setUpTheWholeThing = func(word Word) error {
		standardizedWord := strings.ToLower(word.Navi)
		standardizedWord = strings.ReplaceAll(standardizedWord, "+", "")

		// English
		if !NullDef(word.EN) {
			dictHash2.EN = AssignWord(dictHash2.EN, word.EN, standardizedWord, true)
			dictHash2Parenthesis.EN = AssignWord(dictHash2Parenthesis.EN, word.EN, standardizedWord, false)
		}

		// German (Deutsch)
		if !NullDef(word.DE) {
			dictHash2.DE = AssignWord(dictHash2.DE, word.DE, standardizedWord, true)
			dictHash2Parenthesis.DE = AssignWord(dictHash2Parenthesis.DE, word.DE, standardizedWord, false)
		}

		// Spanish (Español)
		if !NullDef(word.ES) {
			dictHash2.ES = AssignWord(dictHash2.ES, word.ES, standardizedWord, true)
			dictHash2Parenthesis.ES = AssignWord(dictHash2Parenthesis.ES, word.ES, standardizedWord, false)
		}

		// Estonian (Eesti)
		if !NullDef(word.ET) {
			dictHash2.ET = AssignWord(dictHash2.ET, word.ET, standardizedWord, true)
			dictHash2Parenthesis.ET = AssignWord(dictHash2Parenthesis.ET, word.ET, standardizedWord, false)
		}

		// French (Français)
		if !NullDef(word.FR) {
			dictHash2.FR = AssignWord(dictHash2.FR, word.FR, standardizedWord, true)
			dictHash2Parenthesis.FR = AssignWord(dictHash2Parenthesis.FR, word.FR, standardizedWord, false)
		}

		// Hungarian (Magyar)
		if !NullDef(word.HU) {
			dictHash2.HU = AssignWord(dictHash2.HU, word.HU, standardizedWord, true)
			dictHash2Parenthesis.HU = AssignWord(dictHash2Parenthesis.HU, word.HU, standardizedWord, false)
		}

		// Korean (한국어)
		if !NullDef(word.KO) {
			dictHash2.KO = AssignWord(dictHash2.KO, word.KO, standardizedWord, true)
			dictHash2Parenthesis.KO = AssignWord(dictHash2Parenthesis.KO, word.KO, standardizedWord, false)
		}

		// Dutch (Nederlands)
		if !NullDef(word.NL) {
			dictHash2.NL = AssignWord(dictHash2.NL, word.NL, standardizedWord, true)
			dictHash2Parenthesis.NL = AssignWord(dictHash2Parenthesis.NL, word.NL, standardizedWord, false)
		}

		// Polish (Polski)
		if !NullDef(word.PL) {
			dictHash2.PL = AssignWord(dictHash2.PL, word.PL, standardizedWord, true)
			dictHash2Parenthesis.PL = AssignWord(dictHash2Parenthesis.PL, word.PL, standardizedWord, false)
		}

		// Portuguese (Português)
		if !NullDef(word.PT) {
			dictHash2.PT = AssignWord(dictHash2.PT, word.PT, standardizedWord, true)
			dictHash2Parenthesis.PT = AssignWord(dictHash2Parenthesis.PT, word.PT, standardizedWord, false)
		}

		// Russian (Русский)
		if !NullDef(word.RU) {
			dictHash2.RU = AssignWord(dictHash2.RU, word.RU, standardizedWord, true)
			dictHash2Parenthesis.RU = AssignWord(dictHash2Parenthesis.RU, word.RU, standardizedWord, false)
		}

		// Swedish (Svenska)
		if !NullDef(word.SV) {
			dictHash2.SV = AssignWord(dictHash2.SV, word.SV, standardizedWord, true)
			dictHash2Parenthesis.SV = AssignWord(dictHash2Parenthesis.SV, word.SV, standardizedWord, false)
		}

		// Turkish (Türkçe)
		if !NullDef(word.TR) {
			dictHash2.TR = AssignWord(dictHash2.TR, word.TR, standardizedWord, true)
			dictHash2Parenthesis.TR = AssignWord(dictHash2Parenthesis.TR, word.TR, standardizedWord, false)
		}

		// Ukrainian (Українська)
		if !NullDef(word.UK) {
			dictHash2.UK = AssignWord(dictHash2.UK, word.UK, standardizedWord, true)
			dictHash2Parenthesis.UK = AssignWord(dictHash2Parenthesis.UK, word.UK, standardizedWord, false)
		}
		return nil
	}

	var err error
	if mysql {
		err = runOnDB(setUpTheWholeThing)
		if err != nil {
			UncacheHashDict2()
			return err
		}
	} else {
		err = runOnFile(setUpTheWholeThing)
		if err != nil {
			log.Printf("Error caching dictionary: %s", err)
			UncacheHashDict2()
			return err
		}
	}

	dictHash2Cached = true

	return nil
}

func UncacheHashDict() {
	dictHashCached = false
	dictHashLoose = nil
	dictHashStrict = nil
	homonyms = ""
	oddballs = ""
}

func UncacheHashDict2() {
	dictHash2Cached = false
	dictHash2.EN = nil
	dictHash2.DE = nil
	dictHash2.ES = nil
	dictHash2.ET = nil
	dictHash2.FR = nil
	dictHash2.HU = nil
	dictHash2.KO = nil
	dictHash2.NL = nil
	dictHash2.PL = nil
	dictHash2.PT = nil
	dictHash2.RU = nil
	dictHash2.SV = nil
	dictHash2.TR = nil
	dictHash2.UK = nil

	dictHash2Parenthesis.EN = nil
	dictHash2Parenthesis.DE = nil
	dictHash2Parenthesis.ES = nil
	dictHash2Parenthesis.ET = nil
	dictHash2Parenthesis.FR = nil
	dictHash2Parenthesis.HU = nil
	dictHash2Parenthesis.KO = nil
	dictHash2Parenthesis.NL = nil
	dictHash2Parenthesis.PL = nil
	dictHash2Parenthesis.PT = nil
	dictHash2Parenthesis.RU = nil
	dictHash2Parenthesis.SV = nil
	dictHash2Parenthesis.TR = nil
	dictHash2Parenthesis.UK = nil
}

// This will run the function `f` inside the cache or the file directly.
// Use this to get words out of the dictionary
// function `f` is called on every single line in the dictionary!
func RunOnDict(f func(word Word) error) (err error) {
	if dictionaryCached {
		for _, word := range dictionary {
			err = f(word)
			if err != nil {
				return
			}
		}
	} else {
		err = runOnFile(func(word Word) error {
			err = f(word)
			if err != nil {
				return err
			}
			return nil
		})
	}

	return
}

func runOnDB(f func(word Word) error) error {
	user := os.Getenv("FW_USER")
	pass := os.Getenv("FW_PASS")
	host := os.Getenv("FW_HOST")
	name := os.Getenv("FW_DB")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, name)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err1 := db.Query("SELECT " +
		"m.id, m.navi, m.ipa, m.infixes, m.partOfSpeech, s.source, b.stressed, b.syllables, b.infixDots, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'de') AS de, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'en') AS en, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'es') AS es, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'et') AS et, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'fr') AS fr, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'hu') AS hu, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'ko') AS ko, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'nl') AS nl, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'pl') AS pl, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'pt') AS pt, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'ru') AS ru, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'sv') AS sv, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'tr') AS tr, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'uk') AS uk " +
		"FROM fwedit_metaWords AS m " +
		"INNER JOIN fwedit_sources AS s ON (m.id = s.id) " +
		"INNER JOIN fwedit_breakdown AS b ON (s.id = b.id)")

	if err1 != nil {
		return err1
	}

	var w Word
	var de, en, es, et, fr, hu, ko, nl, pl, pt, ru, sv, tr, uk []byte

	for rows.Next() {
		err = rows.Scan(&w.ID, &w.Navi, &w.IPA, &w.InfixLocations, &w.PartOfSpeech, &w.Source, &w.Stressed,
			&w.Syllables, &w.InfixDots, &de, &en, &es, &et, &fr, &hu, &ko, &nl, &pl, &pt, &ru, &sv, &tr, &uk)

		if err != nil {
			return err
		}

		w.DE = string(de)
		w.EN = string(en)
		w.ES = string(es)
		w.ET = string(et)
		w.FR = string(fr)
		w.HU = string(hu)
		w.KO = string(ko)
		w.NL = string(nl)
		w.PL = string(pl)
		w.PT = string(pt)
		w.RU = string(ru)
		w.SV = string(sv)
		w.TR = string(tr)
		w.UK = string(uk)

		err = f(w)

		if err != nil {
			return err
		}
	}

	return nil
}

func runOnFile(f func(word Word) error) error {
	dictionaryFile := FindDictionaryFile()
	if dictionaryFile == "" {
		return DictionaryNotFound
	}

	file, err := os.Open(dictionaryFile)
	if err != nil {
		log.Printf("Error opening the dictionary file %s", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var first = true
	var pos dictPos
	for scanner.Scan() {
		// get a single line out of the dict
		line := scanner.Text()

		// Split line at \t so we get all information
		fields := strings.Split(line, "\t")

		// When first then this is the header
		if first {
			pos = readDictPos(fields)
			first = false
		} else {
			// Put the stuff from fields into the Word struct
			err = f(newWord(fields, pos))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetFullDict() (allWords []Word, err error) {
	// No need for the lock because only List() calls it
	if dictionaryCached {
		allWords = dictionary
	} else {
		err = runOnFile(func(word Word) error {
			allWords = append(allWords, word)
			return nil
		})
		return
	}
	return
}

// Just a number
func GetDictSizeSimple() (count int) {
	universalLock.Lock()
	defer universalLock.Unlock()
	return len(dictionary)
}

// Return a complete sentence
func GetDictSize(lang string) (count string, err error) {
	universalLock.Lock()
	defer universalLock.Unlock()
	// Count words
	amount := 0
	if dictionaryCached {
		amount = len(dictionary)
	} else {
		err = runOnFile(func(word Word) error {
			amount++
			return nil
		})
	}

	// Put the word count into a complete sentence
	count = strconv.Itoa(amount)

	if lang == "en" { // English
		count = "There are " + count + " entries in the dictionary."
	} else if lang == "de" { // German (Deutsch)
		count = "Es sind " + count + " Einträge im Wörterbuch."
	} else if lang == "es" { // Spanish (Español)
		count = "There are " + count + " entries in the dictionary." // TODO
	} else if lang == "et" { // Estonian (Eesti)
		count = "There are " + count + " entries in the dictionary." // TODO
	} else if lang == "fr" { // French (Français)
		count = "Il y a " + count + " définitions dans le dictionnaire."
	} else if lang == "hu" { // Hungarian (Magyar)
		count = "There are " + count + " entries in the dictionary." // TODO
	} else if lang == "ko" { // Korean (한국어)
		count = "Fwew에는 " + count + "개의 단어가 등록되어 있습니다."
	} else if lang == "nl" { // Dutch (Nederlands)
		count = "There are " + count + " entries in the dictionary." // TODO
	} else if lang == "pl" { // Polish (Polski)
		count = "There are " + count + " entries in the dictionary." // TODO
	} else if lang == "pt" { // Portuguese (Português)
		count = "There are " + count + " entries in the dictionary." // TODO
	} else if lang == "ru" { // Russian (Русский)
		count = "There are " + count + " entries in the dictionary." // TODO
	} else if lang == "sv" { // Swedish (Svenska)
		count = "There are " + count + " entries in the dictionary." // TODO
	} else if lang == "tr" { // Turkish (Türkçe)
		count = "There are " + count + " entries in the dictionary." // TODO
	} else if lang == "uk" { // Ukrainian (Українська)
		count = "There are " + count + " entries in the dictionary." // TODO
	}

	return
}

// Update the dictionary.txt.
// universalLock will hopefully prevent anything from accessing
// the dict while updating
func UpdateDict() error {
	universalLock.Lock()
	defer universalLock.Unlock()
	err := DownloadDict("")
	if err != nil {
		log.Println(Text("downloadError"))
		return err
	}

	err = CacheDict()
	if err != nil {
		log.Printf("Error caching dict after updatig ... Cache disabled")
		return err
	}

	if dictHashCached {
		UncacheHashDict()
	}

	err = CacheDictHash()
	if err != nil {
		log.Printf("Error caching dict after updating ... Cache disabled")
		return err
	}

	if dictHash2Cached {
		UncacheHashDict2()
	}

	err = CacheDictHash2()
	if err != nil {
		log.Printf("Error caching dict after updating ... Cache disabled")
		return err
	}

	return nil
}

// AssureDict will assure, that the dictionary exists.
// If no dictionary is found, it will be downloaded next of the executable.
func AssureDict() error {
	// check if dict already exists
	file := FindDictionaryFile()
	if file != "" {
		return nil
	}

	// if it doesn't, put it in ~/.fwew/
	path := filepath.Join(texts["dataDir"], dictFileName)

	err := DownloadDict(path)
	if err != nil {
		return err
	}

	return nil
}
