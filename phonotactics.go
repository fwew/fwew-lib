package fwew_lib

import (
	"strconv"
	"strings"
)

var cluster_1 = []string{"f", "s", "c"}
var cluster_2 = []string{"k", "q", "l", "m", "n", "g", "p",
	"b", "t", "d", "r", "w", "y"}
var letters_start = []string{"", "p", "t", "k", "b", "d", "q", "'",
	"m", "n", "g", "r", "l", "w", "y",
	"f", "v", "s", "z", "c", "h", "B", "D", "G"}
var letters_end = []string{"", "p", "t", "k", "b", "d", "q", "'",
	"m", "n", "l", "r", "g"}

var letters_map = map[string]string{}

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
	word = strings.ReplaceAll(word, "ch", "tsy")
	word = strings.ReplaceAll(word, "sh", "sy")
	word = strings.ReplaceAll(word, "-", "")
	word = strings.TrimSuffix(word, "+")

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
		return "❌ **" + oldWord + "** Has letters not in Na'vi: " + nonNaviLetters
	}

	// Phase 1: don't confuse the digraph compression things
	firstCheckLetters := map[rune]bool{
		'q': false, // kx
		'b': true,  // px
		'd': true,  // tx
		'g': true,  // ng
		'c': false, // ts
		'0': false, // rr
		'1': false, // ll
		'2': false, // aw
		'3': false, // ay
		'4': false, // ew
		'5': false, // ey
	}

	badLetters := ""
	tempWord := ""
	for _, a := range []rune(word) {
		found := false
		if voiced_plosive, ok := firstCheckLetters[a]; ok {
			if voiced_plosive {
				found = true
				tempWord += strings.ToUpper(string(a))
			} else {
				badLetters = badLetters + string(a)
			}
		}

		if !found {
			tempWord += string(a)
		}
	}

	// G is allowed as part of "ng"
	tempWord = strings.ReplaceAll(tempWord, "nG", "ng")

	if badLetters != "" {
		return "❌ **" + oldWord + "** Invalid letters: `" + badLetters + "`"
	}

	// Phase 2: Compress digraphs and divide into syllable boundaries
	compressed := compress(tempWord)
	nuclei := []rune{
		'a', 'ä', 'e', 'i', 'ì', 'o', 'u', 'ù', // vowels
		'0', '1', '2', '3', '4', '5', // diphthongs and psuedovowels
	}

	syllable_boundaries := ""
	word_nuclei := []rune{}
	for _, a := range []rune(compressed) {
		found := false
		for _, b := range nuclei {
			if a == b {
				found = true
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
		return "❌ **" + oldWord + "** Error: could not find any syllable nuclei"
	}

	// Phase 2.1: Go through syllable boundaries
	syllable_breakdown := ""

	for i, a := range strings.Split(syllable_boundaries, ".") {
		a = strings.ReplaceAll(a, " ", "")
		if b, ok := letters_map[a]; ok {
			syllable_breakdown = syllable_breakdown + b
		} else {
			return "❌ **" + oldWord + "** Invalid consonant combination: `" + strings.ToLower(decompress(a)) + "`"
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
		return "❌ **" + oldWord + "** Incomplete syllables: `" + strings.ToLower(decompress(syllable_breakdown)) + "`"
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
			return "❌ **" + oldWord + "** Incomplete syllables: `" + strings.ToLower(decompress(syllable_breakdown)) + "`"
		}

		can_coda := false

		for _, a := range nuclei {
			if strings.HasSuffix(syllables[len(syllables)-2], string(a)) {
				can_coda = true
				break
			}
		}

		if !can_coda {
			return "❌ **" + oldWord + "** Incomplete syllables: `" + strings.ToLower(decompress(syllable_breakdown)) + "`"
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
			return "❌ **" + oldWord + "** Psuedovowels can't accept codas: `" + decompress(strings.ToLower(syllable_breakdown)) + "`"
		}
	}

	// Ensure no diphthong confuses the checker (as in "ewll" becoming "ew-ll" and not "e-wll")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "2-0", "a-w0")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "2-1", "a-w1")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "3-0", "a-y0")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "3-1", "a-y1")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "4-0", "e-w0")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "4-1", "e-w1")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "5-0", "e-y0")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "5-1", "e-y1")

	if strings.Contains(syllable_breakdown, "-0-") || strings.Contains(syllable_breakdown, "-1-") ||
		strings.HasPrefix(syllable_breakdown, "0") || strings.HasPrefix(syllable_breakdown, "1") ||
		strings.HasSuffix(syllable_breakdown, "-0") || strings.HasSuffix(syllable_breakdown, "-1") {
		return "❌ **" + oldWord + "** Psuedovowels must have onsets: `" + decompress(strings.ToLower(syllable_breakdown)) + "`"
	}

	if strings.Contains(syllable_breakdown, "0-r") || strings.Contains(syllable_breakdown, "1-l") {
		return "❌ **" + oldWord + "** Triple Rs or Ls aren't allowed: `" + decompress(strings.ToLower(syllable_breakdown)) + "`"
	}

	// If you reach here, the word is valid
	syllable_breakdown = strings.ToLower(decompress(syllable_breakdown))

	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "ng", "0")

	isReef := ""
	if strings.ContainsAny(syllable_breakdown, "bdg") || strings.Contains(oldWord, "ch") || strings.Contains(oldWord, "sh") {
		syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "tsy", "ch")
		syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "sy", "sh")
		isReef = " (in reef dialect)"
	}
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "0", "ng")

	// Double diphthongs are usually not genuine in Na'vi
	// For example, mawey is ma-wey (not maw-ey) and kxeyey is kxe-yey (not kxey-ey)
	for _, a := range []rune{'a', 'ä', 'e', 'i', 'ì', 'o', 'u', 'ù'} {
		syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "w-"+string(a), "-w"+string(a))
		syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "y-"+string(a), "-y"+string(a))
	}

	syllable_word := " syllables"
	syllable_count := len(strings.Split(syllable_breakdown, "-"))
	if syllable_count == 1 {
		syllable_word = " syllable"
	}
	return "✅ **" + oldWord + "** Valid: `" + syllable_breakdown + "` with " + strconv.Itoa(syllable_count) + syllable_word + isReef
}

func IsValidNavi(word string) string {
	// Let it know of valid syllable boundaries
	if len(letters_map) == 0 {
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

		// Reef dialect can make words like "adge" and "egdu"
		for _, a := range []string{"b", "d", "g"} {
			for _, b := range []string{"b", "d", "g"} {
				if a != b {
					letters_map[a+b] = a + "-" + b
				}
			}
		}
	}
	results := ""
	for i, a := range strings.Split(word, " ") {
		newLine := IsValidNaviHelper(a) + "\n"
		if len(results)+len(newLine) > 1914 {
			results += "⛔ (stopped at " + strconv.Itoa(i+1) + ". 2000 Character limit)"
			break
		}
		results += newLine
	}
	return results
}
