package fwew_lib

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var dictionary []Word
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
		dictionary = append(dictionary, newWord(fields))
	}

	dictionaryCached = true

	return nil
}
