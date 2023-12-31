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
	ET map[string][]string
	FR map[string][]string
	HU map[string][]string
	NL map[string][]string
	PL map[string][]string
	RU map[string][]string
	SV map[string][]string
	TR map[string][]string
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

func AppendAndAlphabetize(words []Word, word Word) []Word {
	newWords := []Word{}
	newWordCompacted := []rune(compress(strings.ToLower(word.Syllables)))
	oldWordCompacted := []rune("bdgj")
	i := 0
	found := false
	// loop through all the words
	for i = 0; i < len(words); i++ {
		if !found {
			oldWordCompacted = []rune(compress(strings.ToLower(words[i].Syllables)))
			lowestLen := len(newWordCompacted)
			if lowestLen > len(oldWordCompacted) {
				lowestLen = len(oldWordCompacted)
			}
			// compare an individual word
			for j := 0; j < lowestLen; j++ {
				// If the new letter is bigger, wait until it gets
				if letterMap[newWordCompacted[j]] > letterMap[oldWordCompacted[j]] {
					break
				} else if letterMap[newWordCompacted[j]] < letterMap[oldWordCompacted[j]] {
					found = true
					newWords = append(newWords, word)
					break
				}
				// If equal, continue
			}
		}
		newWords = append(newWords, words[i])
	}
	if !found {
		newWords = append(newWords, word)
	}
	return newWords
}

// This will cache the whole dictionary.
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
		badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\`

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
	/* English */
	standardizedWord := natlangWords
	badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\`

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

func CacheDictHash2() error {
	// dont run if already is cached
	if len(dictHash2.EN) != 0 {
		return nil
	} else {
		dictHash2.EN = make(map[string][]string)
		dictHash2.DE = make(map[string][]string)
		dictHash2.ET = make(map[string][]string)
		dictHash2.FR = make(map[string][]string)
		dictHash2.HU = make(map[string][]string)
		dictHash2.NL = make(map[string][]string)
		dictHash2.PL = make(map[string][]string)
		dictHash2.RU = make(map[string][]string)
		dictHash2.SV = make(map[string][]string)
		dictHash2.TR = make(map[string][]string)
	}

	// Set up the whole thing

	err := runOnFile(func(word Word) error {
		standardizedWord := word.Navi
		badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\`

		// remove all the sketchy chars from arguments
		for _, c := range badChars {
			standardizedWord = strings.ReplaceAll(standardizedWord, string(c), "")
		}
		dictHash2.EN = AssignWord(dictHash2.EN, word.EN, strings.ToLower(standardizedWord))
		dictHash2.DE = AssignWord(dictHash2.DE, word.DE, strings.ToLower(standardizedWord))
		dictHash2.ET = AssignWord(dictHash2.ET, word.ET, strings.ToLower(standardizedWord))
		dictHash2.FR = AssignWord(dictHash2.FR, word.FR, strings.ToLower(standardizedWord))
		dictHash2.HU = AssignWord(dictHash2.HU, word.HU, strings.ToLower(standardizedWord))
		dictHash2.NL = AssignWord(dictHash2.NL, word.NL, strings.ToLower(standardizedWord))
		dictHash2.PL = AssignWord(dictHash2.PL, word.PL, strings.ToLower(standardizedWord))
		dictHash2.RU = AssignWord(dictHash2.RU, word.RU, strings.ToLower(standardizedWord))
		dictHash2.SV = AssignWord(dictHash2.SV, word.SV, strings.ToLower(standardizedWord))
		dictHash2.TR = AssignWord(dictHash2.TR, word.TR, strings.ToLower(standardizedWord))
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
	dictHash2.ET = nil
	dictHash2.FR = nil
	dictHash2.HU = nil
	dictHash2.NL = nil
	dictHash2.PL = nil
	dictHash2.RU = nil
	dictHash2.SV = nil
	dictHash2.TR = nil
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
