package fwew_lib

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const dictFileName = "dictionary-v2.txt"

var dictionary []Word
var dictHash map[string][]Word
var dictionaryCached bool
var dictHashCached bool
var dictHash2 MetaDict
var dictHash2Cached bool
var homonyms string

type MetaDict struct {
	EN map[string][]string
	DE map[string][]string
	ES map[string][]string
	ET map[string][]string
	FR map[string][]string
	HU map[string][]string
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
	aCompacted := []rune(strings.ReplaceAll((strings.ToLower(a)), "-", ""))

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
	// new array
	switch len(words) {
	case 0:
		return []Word{word}
	case 1:
		var newWords = []Word{}
		if word.ID == words[0].ID {
			return words
		}
		if AlphabetizeHelper(words[0].Syllables, word.Syllables) {
			newWords = []Word{words[0], word}
		} else {
			newWords = []Word{word, words[0]}
		}
		return newWords
	case 2:
		var newWords = []Word{}
		if word.ID == words[0].ID {
			return words
		}
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

// If a definition is not available in a certain language, default to English
func EnglishIfNull(word Word) Word {
	// English
	if word.EN == "NULL" {
		word.EN = "(no definition)"
	}

	// German (Deutsch)
	if word.DE == "NULL" {
		word.DE = word.EN
	}

	// Spanish (Español)
	if word.ES == "NULL" {
		word.ES = word.EN
	}

	// Estonian (Eesti)
	if word.ET == "NULL" {
		word.ET = word.EN
	}

	// French (Français)
	if word.FR == "NULL" {
		word.FR = word.EN
	}

	// Hungarian (Magyar)
	if word.HU == "NULL" {
		word.HU = word.EN
	}

	// Dutch (Nederlands)
	if word.NL == "NULL" {
		word.NL = word.EN
	}

	// Polish (Polski)
	if word.PL == "NULL" {
		word.PL = word.EN
	}

	// Portuguese (Português)
	if word.PT == "NULL" {
		word.PT = word.EN
	}

	// Russian (Русский)
	if word.RU == "NULL" {
		word.RU = word.EN
	}

	// Swedish (Svenska)
	if word.SV == "NULL" {
		word.SV = word.EN
	}

	// Turkish (Türkçe)
	if word.TR == "NULL" {
		word.TR = word.EN
	}

	// Ukrainian (Magyar)
	if word.UK == "NULL" {
		word.UK = word.EN
	}

	return word
}

// This will cache the whole dictionary (Na'vi to natural language).
// Please call this, if you want to translate multiple words or running infinitely (e.g. CLI-go-prompt, discord-bot)
func CacheDictHash() error {
	// dont run if already is cached
	if len(dictHash) != 0 {
		return nil
	} else {
		dictHash = make(map[string][]Word)
	}

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

		word = EnglishIfNull(word)
		dictHash[standardizedWord] = append(dictHash[standardizedWord], word)
		return nil
	})

	// Reverse the order to make accidental and new homonyms easier to see
	// Also make it a string for easier searching
	i := len(tempHoms)
	for i > 0 {
		i--
		homonyms += tempHoms[i] + " "
	}

	homonyms = strings.TrimSuffix(homonyms, " ")

	if err != nil {
		log.Printf("Error caching dictionary: %s", err)
		// uncache dict, to be save
		UncacheHashDict()
		return err
	}

	dictHashCached = true

	return nil
}

// Helper function for CacheDictHash2
func AssignWord(wordmap map[string][]string, natlangWords string, naviWord string) (result map[string][]string) {
	standardizedWord := natlangWords
	badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\"„“”«»`

	// remove anything in parenthesis to avoid clogging search results
	tempString := ""
	parenthesis := false
	for _, c := range natlangWords {
		if c == '(' {
			parenthesis = true
		} else if c == ')' {
			parenthesis = false
			continue
		}

		if !parenthesis {
			tempString += string(c)
		}
	}
	standardizedWord = tempString

	// remove all the sketchy chars from arguments
	for _, c := range badChars {
		standardizedWord = strings.ReplaceAll(standardizedWord, string(c), "")
	}

	// normalize tìftang character
	standardizedWord = strings.ReplaceAll(standardizedWord, "’", "'")
	standardizedWord = strings.ReplaceAll(standardizedWord, "‘", "'")

	// find everything lowercase
	standardizedWord = strings.ToLower(standardizedWord)
	newWords := strings.Split(standardizedWord, " ")

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
		dictHash2.NL = make(map[string][]string)
		dictHash2.PL = make(map[string][]string)
		dictHash2.PT = make(map[string][]string)
		dictHash2.RU = make(map[string][]string)
		dictHash2.SV = make(map[string][]string)
		dictHash2.TR = make(map[string][]string)
		dictHash2.UK = make(map[string][]string)
	}

	// Set up the whole thing

	err := runOnFile(func(word Word) error {
		standardizedWord := strings.ToLower(word.Navi)
		standardizedWord = strings.ReplaceAll(standardizedWord, "+", "")

		// English
		if word.EN != "NULL" {
			dictHash2.EN = AssignWord(dictHash2.EN, word.EN, standardizedWord)
		}

		// German (Deutsch)
		if word.DE != "NULL" {
			dictHash2.DE = AssignWord(dictHash2.DE, word.DE, standardizedWord)
		}

		// Spanish (Español)
		if word.ES != "NULL" {
			dictHash2.ES = AssignWord(dictHash2.ES, word.ES, standardizedWord)
		}

		// Estonian (Eesti)
		if word.ET != "NULL" {
			dictHash2.ET = AssignWord(dictHash2.ET, word.ET, standardizedWord)
		}

		// French (Français)
		if word.FR != "NULL" {
			dictHash2.FR = AssignWord(dictHash2.FR, word.FR, standardizedWord)
		}

		// Hungarian (Magyar)
		if word.HU != "NULL" {
			dictHash2.HU = AssignWord(dictHash2.HU, word.HU, standardizedWord)
		}

		// Dutch (Nederlands)
		if word.NL != "NULL" {
			dictHash2.NL = AssignWord(dictHash2.NL, word.NL, standardizedWord)
		}

		// Polish (Polski)
		if word.PL != "NULL" {
			dictHash2.PL = AssignWord(dictHash2.PL, word.PL, standardizedWord)
		}

		// Portuguese (Português)
		if word.PT != "NULL" {
			dictHash2.PT = AssignWord(dictHash2.PT, word.PT, standardizedWord)
		}

		// Russian (Русский)
		if word.RU != "NULL" {
			dictHash2.RU = AssignWord(dictHash2.RU, word.RU, standardizedWord)
		}

		// Swedish (Svenska)
		if word.SV != "NULL" {
			dictHash2.SV = AssignWord(dictHash2.SV, word.SV, standardizedWord)
		}

		// Turkish (Türkçe)
		if word.TR != "NULL" {
			dictHash2.TR = AssignWord(dictHash2.TR, word.TR, standardizedWord)
		}

		// Ukrainian (Magyar)
		if word.UK != "NULL" {
			dictHash2.UK = AssignWord(dictHash2.UK, word.UK, standardizedWord)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error caching dictionary: %s", err)
		// uncache dict, to be save
		UncacheHashDict2()
		return err
	}

	dictHash2Cached = true

	return nil
}

func UncacheHashDict() {
	dictHashCached = false
	dictHash = nil
}

func UncacheHashDict2() {
	dictHash2Cached = false
	dictHash2.EN = nil
	dictHash2.DE = nil
	dictHash2.ES = nil
	dictHash2.ET = nil
	dictHash2.FR = nil
	dictHash2.HU = nil
	dictHash2.NL = nil
	dictHash2.PL = nil
	dictHash2.PT = nil
	dictHash2.RU = nil
	dictHash2.SV = nil
	dictHash2.TR = nil
	dictHash2.UK = nil
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

func GetDictSize() (amount int, err error) {
	if dictionaryCached {
		amount = len(dictionary)
	} else {
		err = runOnFile(func(word Word) error {
			amount++
			return nil
		})
	}
	return
}

// Update the dictionary.txt.
// Be careful to not do anything with the dict-file, while update is in progress
func UpdateDict() error {
	err := DownloadDict("")
	if err != nil {
		log.Println(Text("downloadError"))
		return err
	}

	if dictionaryCached {
		UncacheHashDict()
		err = CacheDictHash()
		if err != nil {
			log.Printf("Error caching dict after updating ... Cache disabled")
			return err
		}
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

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(wd, dictFileName)

	err = DownloadDict(path)
	if err != nil {
		return err
	}

	return nil
}
