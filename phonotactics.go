package fwew_lib

import (
	"strconv"
	"strings"
)

var cluster1 = []string{"f", "s", "c"}
var cluster2 = []string{"k", "q", "l", "m", "n", "g", "p",
	"b", "t", "d", "r", "w", "y"}
var lettersStart = []string{"", "p", "t", "k", "b", "d", "q", "'",
	"m", "n", "g", "r", "l", "w", "y",
	"f", "v", "s", "z", "c", "h", "B", "D", "G"}
var lettersEnd = []string{"", "p", "t", "k", "b", "d", "q", "'",
	"m", "n", "l", "r", "g"}

var lettersMap = map[string]string{}

func MakeSyllableBreakdown(syllables []string) string {
	syllableBreakdownTemp := ""
	for i, a := range syllables {
		if i != len(syllables)-1 && i != 0 {
			syllableBreakdownTemp += "-"
		}
		syllableBreakdownTemp += a
	}
	return syllableBreakdownTemp
}

// NoDoubleDiphthongs turns things like maw-ey into ma-wey, so the checker knows
// mawll is valid as ma-wll and not mistaken for maw-ll (invalid)
func NoDoubleDiphthongs(syllableBreakdown string) string {
	syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "2-0", "a-w0")
	syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "2-1", "a-w1")
	syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "3-0", "a-y0")
	syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "3-1", "a-y1")
	syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "4-0", "e-w0")
	syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "4-1", "e-w1")
	syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "5-0", "e-y0")
	syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "5-1", "e-y1")
	return syllableBreakdown
}

func ResolveFakePsuedovowels(syllableBreakdown string) string {
	vowels := []string{"a", "ä", "e", "i", "ì", "o", "u", "ù"}
	replacement := map[string]string{"0": "r-r", "1": "l-l"}
	for _, pseudovowel := range []string{"0", "1"} {
		for _, a := range vowels {
			for _, b := range vowels {
				syllableBreakdown = strings.ReplaceAll(syllableBreakdown, a+"-"+pseudovowel+"-"+b, a+replacement[pseudovowel]+b)
			}
		}
	}

	return syllableBreakdown
}

// IsValidNaviHelper see if a word is phonotactically valid in Na'vi
func IsValidNaviHelper(word string, lang string) string {
	// Protect against odd language values
	if _, ok := messageValid[lang]; !ok {
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
	// It used Unicode values to ensure it has nothing invalid
	// We don't have to worry about uppercase letters because it handled them already
	nonNaviLetters := ""
	for _, a := range []rune(word) {
		if int(a) > int('ù') {
			nonNaviLetters += string(a)
		} else if int(a) < int('ù') && int(a) > int('ì') {
			nonNaviLetters += string(a)
		} else if int(a) < int('ì') && int(a) > int('ä') {
			nonNaviLetters += string(a)
		} else if int(a) < int('ä') && int(a) > int('z') {
			nonNaviLetters += string(a)
		} else if int(a) < int('a') && int(a) > int('G') {
			nonNaviLetters += string(a)
		} else if int(a) < int('G') && int(a) > int('>') {
			nonNaviLetters += string(a)
		} else if int(a) < int('>') && int(a) > int('\'') {
			nonNaviLetters += string(a)
		} else if int(a) < int('\'') {
			nonNaviLetters += string(a)
		}
	}

	// ERROR 1a: letters not in Na'vi
	if len(nonNaviLetters) > 0 {
		message := strings.ReplaceAll(messageNonNaviLetters[lang], "{oldWord}", oldWord)
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
		if voicedPlosive, ok := firstCheckLetters[a]; ok {
			if voicedPlosive {
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
		message := strings.ReplaceAll(messageNonNaviLetters[lang], "{oldWord}", oldWord)
		message = strings.ReplaceAll(message, "{nonNaviLetters}", badLetters)
		return "❌ " + message
	}

	// Phase 2: Compress digraphs and divide into syllable boundaries

	compressed := compress(tempWord)
	nuclei := []rune{
		'a', 'ä', 'e', 'i', 'ì', 'o', 'u', 'ù', // vowels
		'0', '1', '2', '3', '4', '5', // diphthongs and psuedovowels
	}

	syllableBoundaries := ""
	var wordNuclei []rune
	for _, a := range []rune(compressed) {
		found := false
		for _, b := range nuclei {
			if a == b {
				found = true
				wordNuclei = append(wordNuclei, a)
				break
			}
		}
		if !found {
			syllableBoundaries = syllableBoundaries + string(a)
		} else {
			syllableBoundaries = syllableBoundaries + " ."
		}
	}

	// ERROR 2: No syllable nuclei
	if len(wordNuclei) == 0 {
		message := strings.ReplaceAll(messageNoNuclei[lang], "{oldWord}", oldWord)
		return "❌ " + message
	}

	// Phase 2.1: Go through syllable boundaries
	syllableBreakdown := ""

	for i, a := range strings.Split(syllableBoundaries, ".") {
		a = strings.ReplaceAll(a, " ", "")
		if b, ok := lettersMap[a]; ok {
			syllableBreakdown = syllableBreakdown + b
		} else { // ERROR 3: Invalid consonant combination
			message := strings.ReplaceAll(messageInvalidConsonants[lang], "{oldWord}", oldWord)
			message = strings.ReplaceAll(message, "{badConsonants}", strings.ToLower(decompress(a)))
			return "❌ " + message
		}
		if i < len(wordNuclei) {
			syllableBreakdown = syllableBreakdown + string(wordNuclei[i])
		}
	}

	syllableBreakdown = strings.ReplaceAll(syllableBreakdown, " ", "")
	syllableBreakdown = strings.TrimPrefix(syllableBreakdown, "-")
	syllableBreakdown = strings.TrimSuffix(syllableBreakdown, "-")

	// Phase 3: Clean up the word and do final checks
	syllables := strings.Split(syllableBreakdown, "-")
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
		syllableBreakdown2 := MakeSyllableBreakdown(syllables)
		syllableBreakdown2 = NoDoubleDiphthongs(syllableBreakdown2)
		message := strings.ReplaceAll(messageNeededVowel[lang], "{oldWord}", oldWord)
		message = strings.ReplaceAll(message, "{breakdown}", strings.ToLower(decompress(syllableBreakdown2)))
		return "❌ " + message
	}

	if !contains[1] {
		canEndAWord := false
		for _, a := range lettersEnd {
			if syllables[len(syllables)-1] == a {
				canEndAWord = true
				break
			}
		}

		if !canEndAWord {
			syllables[len(syllables)-1] += "•"
			syllableBreakdown2 := MakeSyllableBreakdown(syllables)
			syllableBreakdown2 = NoDoubleDiphthongs(syllableBreakdown2)
			message := strings.ReplaceAll(messageNeededVowel[lang], "{oldWord}", oldWord)
			message = strings.ReplaceAll(message, "{breakdown}", strings.ToLower(decompress(syllableBreakdown2)))
			return "❌ " + message
		}

		canCoda := false

		for _, a := range nuclei {
			if strings.HasSuffix(syllables[len(syllables)-2], string(a)) {
				canCoda = true
				break
			}
		}

		if !canCoda {
			syllables[len(syllables)-1] += "•"
			syllableBreakdown2 := MakeSyllableBreakdown(syllables)
			syllableBreakdown2 = NoDoubleDiphthongs(syllableBreakdown2)
			message := strings.ReplaceAll(messageNeededVowel[lang], "{oldWord}", oldWord)
			message = strings.ReplaceAll(message, "{breakdown}", strings.ToLower(decompress(syllableBreakdown2)))
			return "❌ " + message
		}

		syllableBreakdown = MakeSyllableBreakdown(syllables)
	}

	// Ensure no diphthong confuses the checker (as in "ewll" becoming "ew-ll" and not "e-wll")
	syllableBreakdown = NoDoubleDiphthongs(syllableBreakdown)

	syllableBreakdown = ResolveFakePsuedovowels(syllableBreakdown)

	if strings.Contains(syllableBreakdown, "-0-") || strings.Contains(syllableBreakdown, "-1-") ||
		strings.HasPrefix(syllableBreakdown, "0") || strings.HasPrefix(syllableBreakdown, "1") ||
		strings.HasSuffix(syllableBreakdown, "-0") || strings.HasSuffix(syllableBreakdown, "-1") {
		message := strings.ReplaceAll(messagePsuedovowelsMustOnset[lang], "{oldWord}", oldWord)
		message = strings.ReplaceAll(message, "{breakdown}", strings.ToLower(decompress(syllableBreakdown)))
		return "❌ " + message
	}

	// Finally, psuedovowels cannot accept codas
	for _, a := range lettersEnd {
		if a != "" && (strings.Contains(syllableBreakdown, "0"+a) || strings.Contains(syllableBreakdown, "1"+a)) {
			message := strings.ReplaceAll(messagePsuedovowelsCantCoda[lang], "{oldWord}", oldWord)
			message = strings.ReplaceAll(message, "{breakdown}", strings.ToLower(decompress(syllableBreakdown)))
			return "❌ " + message
		}
	}

	forestWarning := ""

	if strings.Contains(syllableBreakdown, "0-r") || strings.Contains(syllableBreakdown, "1-l") {
		forestWarning = messagePsuedovowelAndConsonant[lang]
	}

	// If you reach here, the word is valid
	syllableBreakdown = strings.ToLower(decompress(syllableBreakdown))

	syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "ng", "0")

	isReef := false
	if strings.ContainsAny(syllableBreakdown, "bdg") || strings.Contains(oldWord, "ch") || strings.Contains(oldWord, "sh") {
		syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "tsy", "ch")
		syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "sy", "sh")
		isReef = true
	}
	syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "0", "ng")

	// So does ù
	if !isReef {
		if strings.Contains(syllableBreakdown, "ù") {
			isReef = true
		}
	}

	// Double diphthongs are usually not genuine in Na'vi
	// For example, mawey is ma-wey (not maw-ey) and kxeyey is kxe-yey (not kxey-ey)
	for _, a := range []rune{'a', 'ä', 'e', 'i', 'ì', 'o', 'u', 'ù'} {
		syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "w-"+string(a), "-w"+string(a))
		syllableBreakdown = strings.ReplaceAll(syllableBreakdown, "y-"+string(a), "-y"+string(a))
	}

	// If reef dialect is present, show what forest looks like
	syllableForest := ""
	if isReef {
		syllableForest = strings.ReplaceAll(syllableBreakdown, "sh", "sy")
		syllableForest = strings.ReplaceAll(syllableForest, "ng", "0")
		syllableForest = strings.ReplaceAll(syllableForest, "ch", "tsy")
		syllableForest = strings.ReplaceAll(syllableForest, "b", "px")
		syllableForest = strings.ReplaceAll(syllableForest, "d", "tx")
		syllableForest = strings.ReplaceAll(syllableForest, "g", "kx")
		/*for _, a := range []string{"a", "ä", "e", "i", "ì", "o", "u", "ù"} {
			// First pass: a-a-a-a-a-a is seen as [a-a]-[a-a]-[a-a] and becomes a-ya-a-ya-a-ya
			syllable_forest = strings.ReplaceAll(syllable_forest, a+"-"+a, a+"-y"+a)
			// Second pass: a-ya-a-ya-a-ya is seen as a-y[a-a]-y[a-a]-ya and becomes a-ya-ya-ya-ya-ya
			syllable_forest = strings.ReplaceAll(syllable_forest, a+"-"+a, a+"-y"+a)
		}*/
		syllableForest = strings.ReplaceAll(syllableForest, "i-ì", "i-yì")
		syllableForest = strings.ReplaceAll(syllableForest, "ì-i", "ì-yi")
		syllableForest = strings.ReplaceAll(syllableForest, "0", "ng")
		syllableForest = strings.ReplaceAll(syllableForest, "ù", "u")
		syllableForest = strings.ReplaceAll(messageReefDialect[lang], "{breakdown}", syllableForest)
	}

	letterCheck := syllableBreakdown
	if isReef {
		letterCheck = syllableForest
	}

	// Identical adjacent vowels mean reef Na'vi
	if forestWarning == "" {
		found := false
		for _, a := range []string{"a", "ä", "e", "i", "ì", "o", "u", "ù", "k", "kx", "l", "m", "n", "ng", "p", "px", "r", "t", "tx", "w", "y"} {
			if strings.Contains(letterCheck, a+"-"+a) {
				found = true
				forestWarning = messageIdenticalAdjacentLetters[lang]
				break
			}
		}
		// px-p, tx-t, kx-k
		if !found {
			for _, a := range [][]string{{"kx", "k"}, {"px", "p"}, {"tx", "t"}} {
				if strings.Contains(letterCheck, a[0]+"-"+a[1]) {
					found = true
					forestWarning = messageIdenticalAdjacentLetters[lang]
					break
				}
			}
		}
		if !found {
			if strings.Contains(letterCheck, "i-ì") || strings.Contains(letterCheck, "ì-i") {
				forestWarning = messageIdenticalAdjacentLetters[lang]
			}
		}
	}

	syllableCount := len(strings.Split(syllableBreakdown, "-"))
	message := validMessage(syllableCount, lang)
	message = strings.ReplaceAll(message, "{oldWord}", oldWord)
	message = strings.ReplaceAll(message, "{breakdown}", syllableBreakdown)
	message = strings.ReplaceAll(message, "{syllable_forest}", syllableForest)
	message = message + forestWarning

	return "✅ " + message
}

func IsValidNavi(word string, lang string, twoThousandLimit bool) string {
	// Let it know of valid syllable boundaries
	if len(lettersMap) == 0 {
		for _, a := range lettersEnd {
			for _, b := range lettersStart {
				// Do not assume a thing comes at the end of a word if it doesn't have to
				if !(a != "" && b == "") {
					lettersMap[a+b] = a + "-" + b
				}
			}
			for _, b := range cluster1 {
				for _, c := range cluster2 {
					// Do not assume a thing comes at the end of a word if it doesn't have to
					if !(a != "" && b == "") {
						lettersMap[a+b+c] = a + "-" + b + c
					}
				}
			}
		}

		// Reef dialect can make words like "adge" and "egdu"
		for _, a := range []string{"B", "D", "G"} {
			for _, b := range []string{"B", "D", "G"} {
				if a != b {
					lettersMap[a+b] = a + "-" + b
				}
			}
		}
	}
	results := ""
	for i, a := range strings.Split(word, " ") {
		newLine := IsValidNaviHelper(a, lang) + "\n"
		if twoThousandLimit && len([]rune(results))+len([]rune(newLine)) > 1914 {
			// (stopped at {count}. 2000-Character limit)
			results += strings.ReplaceAll(messageTooBig[lang], "{count}", strconv.Itoa(i+1))
			break
		}
		results += newLine
	}
	return results
}
