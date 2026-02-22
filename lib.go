package fwew_lib

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

// contains returns true if anything in q is also in s.
func contains(s []string, q []string) bool {
	if len(q) == 0 || len(s) == 0 {
		return false
	}
	for _, x := range q {
		for _, y := range s {
			if y == x {
				return true
			}
		}
	}
	return false
}

// identicalRunes returns true if the two strings have the same runes in the same order and false if they don't.
func identicalRunes(first string, second string) bool {
	a := []rune(first)
	b := []rune(second)

	if len(a) != len(b) {
		return false
	}

	for i, c := range a {
		if b[i] != c {
			return false
		}
	}

	return true
}

// downloadDict downloads the latest released version of the dictionary file and saves it to the given filepath.
// You can give an empty string as a filepath param to update the found dictionary file.
func downloadDict(filepath string) error {
	var url = Text("dictURL")

	// only try to find the dictionary file if no path is given
	if filepath == "" {
		filepath = findDictionaryFile()
	}

	// if still no filepath is given, return error
	if filepath == "" {
		return NoDictionary
	}

	// download the new dict
	resp, err := http.Get(url)
	if err != nil {
		return FailedToDownload.wrap(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// create a new file, this will remove the old file, if it exists
	err = os.Mkdir(path.Dir(filepath), 0755)
	out, err := os.Create(filepath)
	if err != nil {
		return FailedToDownload.wrap(err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	// save downloaded dict to the opened file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	err = out.Close()
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}

	// Update the hash in the version
	Version.DictBuild = sha1Hash(filepath)

	return nil
}

// The glob function will test a string pattern, potentially containing globs, against a
// subject string. The result is a simple true/false, determining whether
// the glob pattern matched the subject text.
// https://github.com/ryanuber/go-glob
func glob(pattern, subj string) bool {
	// The character which is treated like a glob
	const GLOB = "%"

	// Empty pattern can only match an empty subject
	if pattern == "" {
		return subj == pattern
	}

	// If the pattern _is_ a glob, it matches everything
	if pattern == GLOB {
		return true
	}

	parts := strings.Split(pattern, GLOB)

	if len(parts) == 1 {
		// No globs in the pattern, so test for equality
		return subj == pattern
	}

	leadingGlob := strings.HasPrefix(pattern, GLOB)
	trailingGlob := strings.HasSuffix(pattern, GLOB)
	end := len(parts) - 1

	// Go over the leading parts and ensure they match.
	for i := 0; i < end; i++ {
		idx := strings.Index(subj, parts[i])

		switch i {
		case 0:
			// Check the first section. Requires special handling.
			if !leadingGlob && idx != 0 {
				return false
			}
		default:
			// Check that the middle parts match.
			if idx < 0 {
				return false
			}
		}

		// Trim evaluated text from subj as we loop over the pattern.
		subj = subj[idx+len(parts[i]):]
	}

	// Reached the last section. Requires special handling.
	return trailingGlob || strings.HasSuffix(subj, parts[end])
}

// sha1Hash gets hash of the dictionary file
func sha1Hash(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(FailedToOpenDictFile.wrap(err))
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(FailedToCloseDictFile.wrap(err))
		}
	}(f)
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))[0:8]
}

// compress compresses or normalizes each digraph of the given string to a unique single character
// inverse of `func decompress(compressed string) string`
func compress(syllables string) string {
	syll := syllables

	ct := make(map[string]string)
	ct["kx"] = "q"
	ct["px"] = "b"
	ct["tx"] = "d"
	ct["ng"] = "g"
	ct["ts"] = "c"
	ct["rr"] = "0"
	ct["ll"] = "1"
	ct["aw"] = "2"
	ct["ay"] = "3"
	ct["ew"] = "4"
	ct["ey"] = "5"
	for key := range ct {
		syll = strings.Replace(syll, key, ct[key], -1)
	}

	return strings.Replace(syll, "-", "", -1)
}

func decompress(syllables string) string {
	syll := syllables

	ct := make(map[string]string)
	ct["q"] = "kx"
	ct["b"] = "px"
	ct["d"] = "tx"
	ct["g"] = "ng"
	ct["c"] = "ts"
	ct["0"] = "rr"
	ct["1"] = "ll"
	ct["2"] = "aw"
	ct["3"] = "ay"
	ct["4"] = "ew"
	ct["5"] = "ey"
	for key := range ct {
		syll = strings.Replace(syll, key, ct[key], -1)
	}

	return syll
}

/* Is it a vowel? (for when the pseudovowel bool won't work) */
func isVowelIpa(letter string) (found bool) {
	// Also arranged from most to least common (not accounting for diphthongs)
	vowels := []string{"a", "ɛ", "ɪ", "o", "u", "i", "æ", "ʊ"}
	// Linear search
	for _, a := range vowels {
		if letter == a {
			return true
		}
	}
	return false
}

func clean(searchNaviWords string) (words string) {
	badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\"„“”«»`

	// remove all the sketchy chars from arguments
	for _, c := range badChars {
		searchNaviWords = strings.ReplaceAll(searchNaviWords, string(c), " ")
	}

	// Recognize line breaks and turn them into spaces
	searchNaviWords = strings.ReplaceAll(searchNaviWords, "\n", " ")

	// No leading or trailing spaces
	searchNaviWords = strings.TrimSpace(searchNaviWords)

	// normalize tìftang character
	searchNaviWords = strings.ReplaceAll(searchNaviWords, "’", "'")
	searchNaviWords = strings.ReplaceAll(searchNaviWords, "‘", "'")

	// find everything lowercase
	searchNaviWords = strings.ToLower(searchNaviWords)

	// Get rid of all double spaces
	for strings.Contains(searchNaviWords, "  ") {
		searchNaviWords = strings.ReplaceAll(searchNaviWords, "  ", " ")
	}

	// No Results if empty string after removing sketch chars
	if len(searchNaviWords) == 0 {
		return
	}

	return searchNaviWords
}
