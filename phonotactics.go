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

func MakeSyllableBreakdown(syllables []string) string {
	syllable_breakdown_temp := ""
	for i, a := range syllables {
		if i != len(syllables)-1 && i != 0 {
			syllable_breakdown_temp += "-"
		}
		syllable_breakdown_temp += a
	}
	return syllable_breakdown_temp
}

// Turns things like maw-ey into ma-wey, so the checker knows
// mawll is valid as ma-wll and not mistaken for maw-ll (invalid)
func NoDoubleDiphthongs(syllable_breakdown string) string {
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "2-0", "a-w0")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "2-1", "a-w1")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "3-0", "a-y0")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "3-1", "a-y1")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "4-0", "e-w0")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "4-1", "e-w1")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "5-0", "e-y0")
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "5-1", "e-y1")
	return syllable_breakdown
}

// See if a word is phonotactically valid in Na'vi
func IsValidNaviHelper(word string, lang string) string {
	// Protect against odd language values
	if _, ok := message_valid[lang]; !ok {
		lang = "en" // default to English
	}

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
	// Non-letters which are acceptable in certain contexts
	word = strings.ReplaceAll(word, "-g", ">G")
	word = strings.ReplaceAll(word, "●g", ">G")
	word = strings.ReplaceAll(word, "•g", ">G")
	word = strings.ReplaceAll(word, "·g", ">G")
	word = strings.ReplaceAll(word, "-", "")
	word = strings.ReplaceAll(word, "•", "")
	word = strings.ReplaceAll(word, "·", "")
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
		} else if int(a) < int(rune('a')) && int(a) > int(rune('G')) {
			nonNaviLetters += string(a)
		} else if int(a) < int(rune('G')) && int(a) > int(rune('>')) {
			nonNaviLetters += string(a)
		} else if int(a) < int(rune('>')) && int(a) > int(rune('\'')) {
			nonNaviLetters += string(a)
		} else if int(a) < int(rune('\'')) {
			nonNaviLetters += string(a)
		}
	}

	// ERROR 1a: letters not in Na'vi
	if len(nonNaviLetters) > 0 {
		message := strings.ReplaceAll(message_non_navi_letters[lang], "{oldWord}", oldWord)
		message = strings.ReplaceAll(message, "{nonNaviLetters}", nonNaviLetters)
		return "❌ " + message
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
	// But the user can say otherwise
	tempWord = strings.ReplaceAll(tempWord, ">G", "G")

	// ERROR 1b: letters not in Na'vi
	if badLetters != "" {
		message := strings.ReplaceAll(message_non_navi_letters[lang], "{oldWord}", oldWord)
		message = strings.ReplaceAll(message, "{nonNaviLetters}", badLetters)
		return "❌ " + message
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

	// ERROR 2: No syllable nuclei
	if len(word_nuclei) == 0 {
		message := strings.ReplaceAll(message_no_nuclei[lang], "{oldWord}", oldWord)
		return "❌ " + message
	}

	// Phase 2.1: Go through syllable boundaries
	syllable_breakdown := ""

	for i, a := range strings.Split(syllable_boundaries, ".") {
		a = strings.ReplaceAll(a, " ", "")
		if b, ok := letters_map[a]; ok {
			syllable_breakdown = syllable_breakdown + b
		} else { // ERROR 3: Invalid consonant combination
			message := strings.ReplaceAll(message_invalid_consonants[lang], "{oldWord}", oldWord)
			message = strings.ReplaceAll(message, "{badConsonants}", strings.ToLower(decompress(a)))
			return "❌ " + message
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

	// ERROR 4a: Incomplete syllables
	if !contains[0] {
		syllables[0] += "•"
		syllable_breakdown_2 := MakeSyllableBreakdown(syllables)
		syllable_breakdown_2 = NoDoubleDiphthongs(syllable_breakdown_2)
		message := strings.ReplaceAll(message_needed_vowel[lang], "{oldWord}", oldWord)
		message = strings.ReplaceAll(message, "{breakdown}", strings.ToLower(decompress(syllable_breakdown_2)))
		return "❌ " + message
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
			syllables[len(syllables)-1] += "•"
			syllable_breakdown_2 := MakeSyllableBreakdown(syllables)
			syllable_breakdown_2 = NoDoubleDiphthongs(syllable_breakdown_2)
			message := strings.ReplaceAll(message_needed_vowel[lang], "{oldWord}", oldWord)
			message = strings.ReplaceAll(message, "{breakdown}", strings.ToLower(decompress(syllable_breakdown_2)))
			return "❌ " + message
		}

		can_coda := false

		for _, a := range nuclei {
			if strings.HasSuffix(syllables[len(syllables)-2], string(a)) {
				can_coda = true
				break
			}
		}

		if !can_coda {
			syllables[len(syllables)-1] += "•"
			syllable_breakdown_2 := MakeSyllableBreakdown(syllables)
			syllable_breakdown_2 = NoDoubleDiphthongs(syllable_breakdown_2)
			message := strings.ReplaceAll(message_needed_vowel[lang], "{oldWord}", oldWord)
			message = strings.ReplaceAll(message, "{breakdown}", strings.ToLower(decompress(syllable_breakdown_2)))
			return "❌ " + message
		}

		syllable_breakdown = MakeSyllableBreakdown(syllables)
	}

	// Finally, psuedovowels cannot accept codas
	for _, a := range letters_end {
		if a != "" && (strings.Contains(syllable_breakdown, "0"+a) || strings.Contains(syllable_breakdown, "1"+a)) {
			message := strings.ReplaceAll(message_psuedovowels_cant_coda[lang], "{oldWord}", oldWord)
			message = strings.ReplaceAll(message, "{breakdown}", strings.ToLower(decompress(syllable_breakdown)))
			return "❌ " + message
		}
	}

	// Ensure no diphthong confuses the checker (as in "ewll" becoming "ew-ll" and not "e-wll")
	syllable_breakdown = NoDoubleDiphthongs(syllable_breakdown)

	if strings.Contains(syllable_breakdown, "-0-") || strings.Contains(syllable_breakdown, "-1-") ||
		strings.HasPrefix(syllable_breakdown, "0") || strings.HasPrefix(syllable_breakdown, "1") ||
		strings.HasSuffix(syllable_breakdown, "-0") || strings.HasSuffix(syllable_breakdown, "-1") {
		message := strings.ReplaceAll(message_psuedovowels_must_onset[lang], "{oldWord}", oldWord)
		message = strings.ReplaceAll(message, "{breakdown}", strings.ToLower(decompress(syllable_breakdown)))
		return "❌ " + message
	}

	if strings.Contains(syllable_breakdown, "0-r") || strings.Contains(syllable_breakdown, "1-l") {
		message := strings.ReplaceAll(message_triple_liquid[lang], "{oldWord}", oldWord)
		message = strings.ReplaceAll(message, "{breakdown}", strings.ToLower(decompress(syllable_breakdown)))
		return "❌ " + message
	}

	// If you reach here, the word is valid
	syllable_breakdown = strings.ToLower(decompress(syllable_breakdown))

	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "ng", "0")

	isReef := false
	if strings.ContainsAny(syllable_breakdown, "bdg") || strings.Contains(oldWord, "ch") || strings.Contains(oldWord, "sh") {
		syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "tsy", "ch")
		syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "sy", "sh")
		isReef = true
	}
	syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "0", "ng")

	// Identical adjacent vowels mean reef Na'vi
	if !isReef {
		found := false
		for _, a := range []string{"a", "ä", "e", "i", "ì", "o", "u", "ù"} {
			if strings.Contains(syllable_breakdown, a+"-"+a) {
				isReef = true
				found = true
				break
			}
		}
		if !found {
			if strings.Contains(syllable_breakdown, "i-ì") || strings.Contains(syllable_breakdown, "ì-i") {
				isReef = true
			}
		}
	}

	// So does ù
	if !isReef {
		if strings.Contains(syllable_breakdown, "ù") {
			isReef = true
		}
	}

	// Double diphthongs are usually not genuine in Na'vi
	// For example, mawey is ma-wey (not maw-ey) and kxeyey is kxe-yey (not kxey-ey)
	for _, a := range []rune{'a', 'ä', 'e', 'i', 'ì', 'o', 'u', 'ù'} {
		syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "w-"+string(a), "-w"+string(a))
		syllable_breakdown = strings.ReplaceAll(syllable_breakdown, "y-"+string(a), "-y"+string(a))
	}

	// If reef dialect is present, show what forest looks like
	syllable_forest := ""
	if isReef {
		syllable_forest = strings.ReplaceAll(syllable_breakdown, "sh", "sy")
		syllable_forest = strings.ReplaceAll(syllable_forest, "ng", "0")
		syllable_forest = strings.ReplaceAll(syllable_forest, "ch", "tsy")
		syllable_forest = strings.ReplaceAll(syllable_forest, "b", "px")
		syllable_forest = strings.ReplaceAll(syllable_forest, "d", "tx")
		syllable_forest = strings.ReplaceAll(syllable_forest, "g", "kx")
		for _, a := range []string{"a", "ä", "e", "i", "ì", "o", "u", "ù"} {
			// First pass: a-a-a-a-a-a is seen as [a-a]-[a-a]-[a-a] and becomes a-ya-a-ya-a-ya
			syllable_forest = strings.ReplaceAll(syllable_forest, a+"-"+a, a+"-y"+a)
			// Second pass: a-ya-a-ya-a-ya is seen as a-y[a-a]-y[a-a]-ya and becomes a-ya-ya-ya-ya-ya
			syllable_forest = strings.ReplaceAll(syllable_forest, a+"-"+a, a+"-y"+a)
		}
		syllable_forest = strings.ReplaceAll(syllable_forest, "i-ì", "i-yì")
		syllable_forest = strings.ReplaceAll(syllable_forest, "ì-i", "ì-yi")
		syllable_forest = strings.ReplaceAll(syllable_forest, "0", "ng")
		syllable_forest = strings.ReplaceAll(message_reef_dialect[lang], "{breakdown}", syllable_forest)
	}

	syllable_count := len(strings.Split(syllable_breakdown, "-"))
	message := valid_message(syllable_count, lang)
	message = strings.ReplaceAll(message, "{oldWord}", oldWord)
	message = strings.ReplaceAll(message, "{breakdown}", syllable_breakdown)
	message = strings.ReplaceAll(message, "{syllable_forest}", syllable_forest)

	return "✅ " + message
}

func IsValidNavi(word string, lang string, two_thousand_limit bool) string {
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
		for _, a := range []string{"B", "D", "G"} {
			for _, b := range []string{"B", "D", "G"} {
				if a != b {
					letters_map[a+b] = a + "-" + b
				}
			}
		}
	}
	results := ""
	for i, a := range strings.Split(word, " ") {
		newLine := IsValidNaviHelper(a, lang) + "\n"
		if two_thousand_limit && len([]rune(results))+len([]rune(newLine)) > 1914 {
			// (stopped at {count}. 2000 Character limit)
			results += strings.ReplaceAll(message_too_big[lang], "{count}", strconv.Itoa(i+1))
			break
		}
		results += newLine
	}
	return results
}
