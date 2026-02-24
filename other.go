package fwew_lib

import (
	"slices"
	"strconv"
)

// GetLenitionTable returns the lenition table
func GetLenitionTable() [8][2]string {
	return lenitionTable
}

// GetShortLenitionTable returns the shortened lenition table
func GetShortLenitionTable() [4][2]string {
	return shortLenitionTable
}

// GetThatTable returns the "that" table
func GetThatTable() [9][5]string {
	return thatTable
}

// GetOtherThats returns the "other 'that's" table
func GetOtherThats() [9][3]string {
	return otherThats
}

// GetMultiwordWords returns all words with spaces
func GetMultiwordWords() map[string][][]string {
	universalLock.Lock()
	defer universalLock.Unlock()
	return multiwordWords
}

// GetHomonyms returns all words with multiple definitions
func GetHomonyms() (results [][]Word) {
	return TranslateFromNaviHash(homonyms, false, false, false)
}

// GetOddballs returns all words with non-standard phonotactics
func GetOddballs() (results [][]Word) {
	return TranslateFromNaviHash(oddballs, true, false, false)
}

// GetMultiIPA returns all words with multiple definitions
func GetMultiIPA() (results [][]Word) {
	return TranslateFromNaviHash(multiIPA, false, false, false)
}

// GetPhonemeDistrosMap returns phoneme distribution data for a Na'vi language
func GetPhonemeDistrosMap(lang string) (allDistros [][][]string) {
	phonoLock.Lock()
	defer phonoLock.Unlock()

	// Default to English
	headerLang := defaultHeaderLang
	clusterLang := defaultClusterLang

	if a, ok := headerRow[lang]; ok {
		headerLang = a
	}
	if a, ok := clusterName[lang]; ok {
		clusterLang = a
	}

	allDistros = [][][]string{
		{headerLang},
		{{clusterLang, "f", "s", "ts"}},
	}

	// Convert them to tuples for sorting
	var onsetTuples []phonemeTuple
	for key, val := range onsetMap {
		onsetTuples = append(onsetTuples, phonemeTuple{val, key})
	}
	slices.SortFunc(tuples(onsetTuples), func(a, b phonemeTuple) int {
		return b.value - a.value
	})

	var nucleusTuples []phonemeTuple
	for key, val := range nucleusMap {
		nucleusTuples = append(nucleusTuples, phonemeTuple{val, key})
	}
	slices.SortFunc(tuples(nucleusTuples), func(a, b phonemeTuple) int {
		return b.value - a.value
	})

	var codaTuples []phonemeTuple
	for key, val := range codaMap {
		codaTuples = append(codaTuples, phonemeTuple{val, key})
	}
	slices.SortFunc(tuples(codaTuples), func(a, b phonemeTuple) int {
		return b.value - a.value
	})

	// Probably not needed but just in case any other number exceeds it
	maxLen := max(len(nucleusTuples), len(onsetTuples))
	if len(codaTuples) > maxLen {
		maxLen = len(codaTuples)
	}

	// Put them into a 2d string array
	i := 0
	for i < maxLen {
		allDistros[0] = append(allDistros[0], []string{})
		c := len(allDistros[0]) - 1

		allDistros[0][c] = appendPhoneme(allDistros[0][c], onsetTuples, i)
		allDistros[0][c] = appendPhoneme(allDistros[0][c], nucleusTuples, i)
		allDistros[0][c] = appendPhoneme(allDistros[0][c], codaTuples, i)

		i += 1
	}

	// Cluster time
	for _, a := range cluster2Full {
		allDistros[1] = append(allDistros[1], []string{a})
		c := len(allDistros[1]) - 1
		for _, b := range cluster1Full {
			allDistros[1][c] = append(allDistros[1][c], strconv.Itoa(clusterMap[b][a]))
		}
	}

	return
}

func appendPhoneme(row []string, tuples []phonemeTuple, i int) []string {
	if i < len(tuples) {
		return append(row, tuples[i].letter+" "+strconv.Itoa(tuples[i].value))
	}
	return append(row, "")
}
