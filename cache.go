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

	err := runOnFile(func(word Word) error {
		dictionary[word.LangCode] = append(dictionary[word.LangCode], word)
		return nil
	})
	if err != nil {
		log.Printf("Error caching dictionary: %s", err)
		// uncache dict, to be save
		UncacheDict()
		return err
	}

	dictionaryCached = true

	return nil
}

func UncacheDict() {
	dictionary = make(map[string][]Word)
	dictionaryCached = false
}

// This will run the function `f` inside the cache or the file directly.
// Use this to get words out of the dictionary
// function `f` is called on every single line that matches the langCode!
func RunOnDict(lang string, f func(word Word) error) (err error) {
	if dictionaryCached {
		for _, word := range dictionary[lang] {
			err = f(word)
			if err != nil {
				return
			}
		}
	} else {
		err = runOnFile(func(word Word) error {
			if word.LangCode == lang {
				err = f(word)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	return
}

func runOnFile(f func(word Word) error) error {
	dictionaryFile := findDictionaryFile()
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

	for scanner.Scan() {
		// get a single line out of the dict
		line := scanner.Text()

		// Split line at \t so we get all information
		fields := strings.Split(line, "\t")

		// Put the stuff from fields into the Word struct
		err = f(newWord(fields))
		if err != nil {
			return err
		}
	}

	return nil
}
