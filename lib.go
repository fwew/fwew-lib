//	This file is part of Fwew.
//	Fwew is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU General Public License as published by
// 	the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
//	Fwew is distributed in the hope that it will be useful,
//	but WITHOUT ANY WARRANTY; without even implied warranty of
//	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	GNU General Public License for more details.
//
//	You should have received a copy of the GNU General Public License
//	along with Fwew.  If not, see http://gnu.org/licenses/

// Package main contains all the things. lib.go handles common functions.
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
	"unicode"
)

// Contains returns true if anything in q is also in s
func Contains(s []string, q []string) bool {
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

// ContainsStr returns true if q is in s
func ContainsStr(s []string, q string) bool {
	if len(q) == 0 || len(s) == 0 {
		return false
	}
	for _, x := range s {
		if q == x {
			return true
		}
	}
	return false
}

// ContainsStrArr returns true if anything in q is in s
func HasPrefixStrArr(s string, q []string) bool {
	if len(q) == 0 || len(s) == 0 {
		return false
	}
	for _, x := range q {
		if strings.HasPrefix(s, x) {
			return true
		}
	}
	return false
}

// ArrCount returns the number of occurrences of q in s
func ArrCount(s []string, q string) int {
	dict := make(map[string]int)
	for _, x := range s {
		dict[x]++
	}
	return dict[q]
}

// DeleteElement "deletes" all occurrences of q in s
// actually returns a new string slice containing the original minus all q
func DeleteElement(s []string, q string) []string {
	if len(s) == 0 {
		return s
	}
	var r []string
	for _, str := range s {
		if str != q {
			r = append(r, str)
		}
	}
	return r
}

// DeleteEmpty "deletes" all empty string entries in s
// actually returns a new string slice containing all non-empty strings in s
func DeleteEmpty(s []string) []string {
	return DeleteElement(s, "")
}

// Index return the index of q in s
func Index(s []string, q string) int {
	for i, v := range s {
		if v == q {
			return i
		}
	}
	return -1
}

// IndexStr return the index of q in s
func IndexStr(s string, q rune) int {
	for i, v := range s {
		if v == q {
			return i
		}
	}
	return -1
}

// IsLetter returns true if s is an alphabet character or apostrophe
func IsLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || r == '\'' || r == '‘' {
			return true
		}
	}
	return false
}

// Reverse returns the reversed version of s
func Reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, r := range s {
		n--
		runes[n] = r
	}
	return string(runes[n:])
}

// StripChars strips all the characters in chr out of str
func StripChars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if !strings.ContainsRune(chr, r) {
			return r
		}
		return -1
	}, str)
}

// DownloadDict downloads the latest released version of the dictionary file and saves it to the given filepath.
// You can give an empty string as filepath param, to update the found dictionary file.
func DownloadDict(filepath string) error {
	var url = Text("dictURL")

	// only try to find dictionary-file if no path is given
	if filepath == "" {
		filepath = FindDictionaryFile()
	}

	// if still no filepath is given, error out
	if filepath == "" {
		return DictFileNotFoundError
	}

	// download the new dict
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// create new file, this will remove the old file, if it exists
	err = os.Mkdir(path.Dir(filepath), 0755)
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

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
	Version.DictBuild = SHA1Hash(filepath)

	return nil
}

// GLOB https://github.com/ryanuber/go-glob
// The character which is treated like a glob
const GLOB = "%"

// Glob will test a string pattern, potentially containing globs, against a
// subject string. The result is a simple true/false, determining whether or
// not the glob pattern matched the subject text.
func Glob(pattern, subj string) bool {
	// Empty pattern can only match empty subject
	if pattern == "" {
		return subj == pattern
	}

	// If the pattern _is_ a glob, it matches everything
	if pattern == GLOB {
		return true
	}

	parts := strings.Split(pattern, GLOB)

	if len(parts) == 1 {
		// No globs in pattern, so test for equality
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

// SHA1Hash gets hash of dictionary file
func SHA1Hash(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
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
