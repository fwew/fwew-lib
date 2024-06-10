package fwew_lib

import (
	"strings"
)

// See if a word is phonotactically valid in Na'vi
func IsValidNaviHelper(word string) string {
	oldWord := word
	// Phase 0: Clean up the word
	word = strings.ToLower(word)
	word = strings.Trim(word, " ")
	// normalize tìftang character
	word = strings.ReplaceAll(word, "’", "'")
	word = strings.ReplaceAll(word, "‘", "'")
	// Normalize acute accent
	word = strings.ReplaceAll(word, "á", "a")
	word = strings.ReplaceAll(word, "é", "e")
	word = strings.ReplaceAll(word, "í", "i")
	word = strings.ReplaceAll(word, "ó", "o")
	word = strings.ReplaceAll(word, "ú", "u")

	// Make sure it doesn't have any invalid letters
	// It used unicode values to ensure it has nothing invalid
	// We don't have to worry about uppercase letters because it handled them already
	nonNaviLetters := ""
	for _, a := range []rune(word) {
		if int(a) > int(rune('ù')) {
			nonNaviLetters += string(a)
		} else if int(a) < int(rune('ù')) && int(a) > int(rune('ì')) {
			nonNaviLetters += string(a)
		} else if int(a) < int(rune('ì')) && int(a) > int(rune('ä')) {
			nonNaviLetters += string(a)
		} else if int(a) < int(rune('ä')) && int(a) > int(rune('z')) {
			nonNaviLetters += string(a)
		} else if int(a) < int(rune('a')) && int(a) > int(rune('\'')) {
			nonNaviLetters += string(a)
		} else if int(a) < int(rune('\'')) {
			nonNaviLetters += string(a)
		}
	}

	if len(nonNaviLetters) > 0 {
		return oldWord + " Has letters not in Na'vi: " + nonNaviLetters
	}

	// Phase 1: don't confuse the digraph compression things
	firstCheckLetters := map[rune]bool{
		'q': true, // kx
		'b': true, // px
		'd': true, // tx
		'g': true, // ng
		'c': true, // ts
		'0': true, // rr
		'1': true, // ll
		'2': true, // aw
		'3': true, // ay
		'4': true, // ew
		'5': true, // ey
	}

	badLetters := ""
	for i, a := range []rune(word) {
		// G is allowed as part of "ng"
		if a == 'g' {
			if !(i > 0 && []rune(word)[i-1] == 'n') {
				badLetters = badLetters + string(a)
			}
			continue
		}
		if _, ok := firstCheckLetters[a]; ok {
			badLetters = badLetters + string(a)
		}
	}

	if badLetters != "" {
		return oldWord + " Invalid letters: " + badLetters
	}

	// Phase 2: Compress digraphs and divide into syllable boundaries
	compressed := compress(word)
	nuclei := []rune{
		'a', 'ä', 'e', 'i', 'ì', 'o', 'u', 'ù', // vowels
		'0', '1', '2', '3', '4', '5', // diphthongs and psuedovowels
	}
	psuedovowels := []bool{}

	syllable_boundaries := ""
	word_nuclei := []rune{}
	for _, a := range []rune(compressed) {
		found := false
		for _, b := range nuclei {
			if a == b {
				found = true
				if a == '0' || a == '1' {
					psuedovowels = append(psuedovowels, true)
				} else {
					psuedovowels = append(psuedovowels, false)
				}
				word_nuclei = append(word_nuclei, a)
				break
			}
		}
		if !found {
			syllable_boundaries = syllable_boundaries + string(a)
		} else {
			syllable_boundaries = syllable_boundaries + " ."
		}
	}

	if len(word_nuclei) == 0 {
		return "Error: could not find any syllable nuclei in " + oldWord
	}

	// Phase 2.1: Go through syllable boundaries
	cluster_1 := []string{"f", "s", "c"}
	cluster_2 := []string{"k", "q", "l", "m", "n", "g", "p",
		"b", "t", "d", "r", "w", "y"}
	letters_start := []string{"", "p", "t", "k", "b", "d", "q", "'",
		"m", "n", "g", "r", "l", "w", "y",
		"f", "v", "s", "z", "c", "h"}
	letters_end := []string{"", "p", "t", "k", "b", "d", "q", "'",
		"m", "n", "l", "r", "g"}

	letters_map := map[string]string{}
	for _, a := range letters_end {
		for _, b := range letters_start {
			// Do not assume a thing comes at the end of a word if it doesn't have to
			if !(a != "" && b == "") {
				letters_map[a+b] = a + "-" + b
			}
		}
		for _, b := range cluster_1 {
			for _, c := range cluster_2 {
				// Do not assume a thing comes at the end of a word if it doesn't have to
				if !(a != "" && b == "") {
					letters_map[a+b+c] = a + "-" + b + c
				}
			}
		}
	}

	syllable_breakdown := ""

	for i, a := range strings.Split(syllable_boundaries, ".") {
		a = strings.ReplaceAll(a, " ", "")
		if b, ok := letters_map[a]; ok {
			syllable_breakdown = syllable_breakdown + b
		} else {
			return oldWord + " Invalid consonants: \"" + decompress(a) + "\""
		}
		if i < len(word_nuclei) {
			syllable_breakdown = syllable_breakdown + string(word_nuclei[i])
		}
	}

	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, " ", "")
	syllable_breakdown = strings.TrimPrefix(syllable_breakdown, "-")
	syllable_breakdown = strings.TrimSuffix(syllable_breakdown, "-")

	// Phase 3: Clean up the word and do final checks
	syllables := strings.Split(syllable_breakdown, "-")
	contains := []bool{false, false}
	for _, a := range nuclei {
		if strings.Contains(syllables[0], string(a)) {
			contains[0] = true
		}
		if strings.Contains(syllables[len(syllables)-1], string(a)) {
			contains[1] = true
		}
	}

	if !contains[0] {
		return oldWord + " Incomplete syllables: \"" + decompress(syllable_breakdown) + "\""
	}

	if !contains[1] {
		can_end_a_word := false
		for _, a := range letters_end {
			if syllables[len(syllables)-1] == string(a) {
				can_end_a_word = true
				break
			}
		}

		if !can_end_a_word {
			return oldWord + " Incomplete syllables: \"" + decompress(syllable_breakdown) + "\""
		}

		can_coda := false

		for _, a := range nuclei {
			if strings.HasSuffix(syllables[len(syllables)-2], string(a)) {
				can_coda = true
				break
			}
		}

		if !can_coda {
			return oldWord + " Incomplete syllables: \"" + decompress(syllable_breakdown) + "\""
		}

		syllable_breakdown_temp := ""

		for i, a := range syllables {
			if i != len(syllables)-1 {
				syllable_breakdown_temp += "-"
			}
			syllable_breakdown_temp += a
		}

		syllable_breakdown = strings.TrimPrefix(syllable_breakdown_temp, "-")
	}

	// Finally, psuedovowels cannot accept codas
	for _, a := range letters_end {
		if a != "" && (strings.Contains(syllable_breakdown, "0"+a) || strings.Contains(syllable_breakdown, "1"+a)) {
			return oldWord + " Psuedovowels can't accept codas: " + decompress(syllable_breakdown)
		}
	}

	if strings.Contains(syllable_breakdown, "-0-") || strings.Contains(syllable_breakdown, "-1-") ||
		strings.HasPrefix(syllable_breakdown, "0") || strings.HasPrefix(syllable_breakdown, "1") ||
		strings.HasSuffix(syllable_breakdown, "-0") || strings.HasSuffix(syllable_breakdown, "-1") {
		return oldWord + " Psuedovowels must have onsets: " + decompress(syllable_breakdown)
	}

	return oldWord + " Valid: " + decompress(syllable_breakdown)
}

func IsValidNavi(word string) string {
	results := ""
	for _, a := range strings.Split(word, " ") {
		results += IsValidNaviHelper(a) + "\n"
	}
	return results
}
