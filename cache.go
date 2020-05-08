package fwew_lib

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//var dictionary map[string][]Word
var dictionary = make(map[string][]Word)
var dictionaryCached bool

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
func findDictionaryFile() string {
	const dictFileName = "dictionary.txt"

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

// This will cache the whole dictionary.
// Please call this, if you want to translate multiple words or running infinitely (e.g. CLI-go-prompt, discord-bot)
func CacheDict() error {
	// dont run if already is cached
	if len(dictionary) != 0 {
		return nil
	}

	runOnFile(func(word Word) {
		dictionary[word.LangCode] = append(dictionary[word.LangCode], word)
	})

	dictionaryCached = true

	return nil
}

// This will run the function `f` inside the cache or the file directly.
// Use this to get words out of the dictionary
// function `f` is called on every single line that matches the langCode!
func RunOnDict(lang string, f func(word Word)) {
	if dictionaryCached {
		for _, word := range dictionary[lang] {
			f(word)
		}
	} else {
		runOnFile(func(word Word) {
			if word.LangCode == lang {
				f(word)
			}
		})
	}
}

func runOnFile(f func(word Word)) error {
	dictionaryFile := findDictionaryFile()
	if dictionaryFile == "" {
		return errors.New("no dictionary found")
	}

	file, err := os.Open(dictionaryFile)
	if err != nil {
		log.Printf("Error opening the dictionary file %s", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// get a single line out of the dict
		line := scanner.Text()

		// Split line at \t so we get all information
		fields := strings.Split(line, "\t")

		// Put the stuff from fields into the Word struct
		f(newWord(fields))
	}

	return nil
}
