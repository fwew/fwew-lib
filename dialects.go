package fwew_lib

import (
	"strconv"
	"strings"
)

func dialectCrunch(query []string, guaranteedForest bool, allowReef bool) []string {
	var newQuery []string
	for _, a := range query {
		oldQuery := a

		// When caching, we are guaranteed forest words and don't need anything in this block
		if !guaranteedForest && allowReef {
			for i, b := range nkx {
				// make sure words like tìkankxan show up
				a = strings.ReplaceAll(a, strconv.Itoa(i), "")
				a = strings.ReplaceAll(a, b, strconv.Itoa(i))
			}
			// don't accidentally make every ng into nkx
			a = strings.ReplaceAll(a, "?", "")
			a = strings.ReplaceAll(a, "ng", "?")
			// un-soften ejectives
			a = strings.ReplaceAll(a, "b", "px")
			a = strings.ReplaceAll(a, "d", "tx")
			a = strings.ReplaceAll(a, "g", "kx")
			// these too
			a = strings.ReplaceAll(a, "ch", "tsy")
			a = strings.ReplaceAll(a, "sh", "sy")
			a = strings.ReplaceAll(a, "?", "ng")
			for i, b := range nkx {
				// make sure words like tìkankxan show up
				a = strings.ReplaceAll(a, strconv.Itoa(i), nkxSub[b])
			}
		}

		if allowReef {
			nucleusCount := 0
			// remove reef tìftangs
			for i, b := range []string{"a", "ä", "e", "i", "ì", "o", "u", "ù", "ll", "rr"} {
				if strings.Contains(a, b) {
					nucleusCount += strings.Count(a, b)
					for j, c := range []string{"a", "ä", "e", "i", "ì", "o", "u", "ù", "ll", "rr"} {
						if i < 8 && j < 8 {
							a = strings.ReplaceAll(a, b+"'"+c, b+c)
						}
					}
				}
			}
			if nucleusCount > 1 && strings.Contains(a, "ä") {
				// and to make sure every ä is possibly an e
				a = strings.ReplaceAll(a, "ä", "e")
			}

			// "eo" and "äo" are different, so the distinction must remain
			if strings.HasSuffix(oldQuery, "äo") || strings.HasSuffix(oldQuery, "ä'o") {
				a = strings.TrimSuffix(a, "eo") + "äo"
			}
		}

		newQuery = append(newQuery, a)
	}
	return newQuery
}

func ReefMe(ipa string, inter bool) []string {
	if ipa == "ʒɛjk'.ˈsu:.li" { // Obsolete path
		return []string{"jake-__sùl__-ly", "ʒɛjk'.ˈsʊ:.li"}
	} else if strings.ReplaceAll(ipa, "·", "") == "ˈzɛŋ.kɛ" { // only IPA not to match the Romanization
		return []string{"__zen__-ke", "ˈz·ɛŋ·.kɛ"}
	} else if ipa == "ɾæ.ˈʔæ" || ipa == "ˈɾæ.ʔæ" { // we hear this in Avatar 2
		return []string{"rä-__'ä__ or rä-__ä__", "ɾæ.ˈʔæ] or [ɾæ.ˈæ"}
	}

	// Replace the spaces so as not to confuse strings.Split()
	ipa = strings.ReplaceAll(ipa, " ", "*.")

	// Unstressed "ä" becomes "e"
	ipaSyllables := strings.Split(ipa, ".")
	if len(ipaSyllables) > 1 {
		var newIpa strings.Builder
		for _, a := range ipaSyllables {
			newIpa.WriteString(".")
			if !strings.Contains(a, "ˈ") {
				newIpa.WriteString(strings.ReplaceAll(a, "æ", "ɛ"))
			} else {
				newIpa.WriteString(a)
			}
		}

		ipa = newIpa.String()
	}

	breakdown := ""
	ejectives := []string{"p'", "t'", "k'"}
	soften := map[string]string{
		"p'": "b",
		"t'": "d",
		"k'": "g",
	}

	// Reefify the IPA first
	ipaReef := strings.ReplaceAll(ipa, "·", "")
	if !inter {
		// atxkxe and ekxtxu
		for _, a := range ejectives {
			for _, b := range ejectives {
				ipaReef = strings.ReplaceAll(ipaReef, a+".ˈ"+b, soften[a]+".ˈ"+soften[b])
				ipaReef = strings.ReplaceAll(ipaReef, a+"."+b, soften[a]+"."+soften[b])
			}
		}

		// Ejectives before vowels and diphthongs become voiced plosives regardless of syllable boundaries
		for _, a := range ejectives {
			if after, ok := strings.CutPrefix(ipaReef, a); ok {
				ipaReef = soften[a] + after
			}
			ipaReef = strings.ReplaceAll(ipaReef, ".ˈ"+a, ".ˈ"+soften[a])
			ipaReef = strings.ReplaceAll(ipaReef, "."+a, "."+soften[a])

			for _, b := range []string{"a", "ɛ", "ɪ", "o", "u", "i", "æ", "ʊ"} {
				ipaReef = strings.ReplaceAll(ipaReef, a+".ˈ"+b, soften[a]+".ˈ"+b)
				ipaReef = strings.ReplaceAll(ipaReef, a+"."+b, soften[a]+"."+b)
			}
		}

		ipaReef = strings.ReplaceAll(ipaReef, "t͡sj", "tʃ")
		ipaReef = strings.ReplaceAll(ipaReef, "sj", "ʃ")

		var temp strings.Builder
		runes := []rune(ipaReef)

		// Glottal stops between two vowels are removed
		for i, a := range runes {
			if i != 0 && i != len(runes)-1 && a == 'ʔ' {
				if runes[i-1] == '.' && i > 1 {
					if isVowelIpa(string(runes[i+1])) && isVowelIpa(string(runes[i-2])) {
						if runes[i+1] != runes[i-2] {
							continue
						}
					}
				} else if runes[i+1] == '.' {
					if isVowelIpa(string(runes[i+2])) && isVowelIpa(string(runes[i-1])) {
						if runes[i+2] != runes[i-1] {
							continue
						}
					}
				} else if runes[i-1] == 'ˈ' && i > 2 {
					if isVowelIpa(string(runes[i+1])) && isVowelIpa(string(runes[i-3])) {
						if runes[i+1] != runes[i-3] {
							continue
						}
					}
				}
			}
			temp.WriteString(string(a))
		}

		ipaReef = temp.String()
	}

	ipaReef = strings.TrimPrefix(ipaReef, ".")

	ipaReef = strings.ReplaceAll(ipaReef, "*.", " ")

	// now Romanize the reef IPA
	word := strings.Split(ipaReef, " ")

	breakdown = ""

	for j := range word {
		word[j] = strings.ReplaceAll(word[j], "]", "")
		word[j] = strings.ReplaceAll(word[j], "[", "")
		// "or" means there's more than one IPA in this word, and we only want one
		if word[j] == "or" {
			breakdown += "or "
			continue
		}

		syllables := strings.Split(word[j], ".")

		/* Onset */
		for k := range syllables {
			breakdown += "-"

			syllable := strings.ReplaceAll(syllables[k], "·", "")
			if strings.Contains(syllable, "ˈ") {
				breakdown += "__"
			}
			syllable = strings.ReplaceAll(syllable, "ˈ", "")
			syllable = strings.ReplaceAll(syllable, "ˌ", "")

			breakdown += syllableToRoman(syllable)

			if strings.Contains(syllables[k], "ˈ") {
				breakdown += "__"
			}
		}
		breakdown += " "
	}

	breakdown = strings.TrimPrefix(breakdown, "-")
	breakdown = strings.ReplaceAll(breakdown, " -", " ")
	breakdown = strings.TrimSuffix(breakdown, " ")

	// If there's a tìftang between two identical vowels, the tìftang is optional
	shortString := strings.ReplaceAll(strings.ReplaceAll(ipaReef, "ˈ", ""), ".", "")
	for _, a := range []string{"a", "ɛ", "ɪ", "o", "u", "i", "æ", "ʊ"} {
		if strings.Contains(shortString, a+"ʔ"+a) {
			// fix IPA
			noGlottalStopIPA := strings.ReplaceAll(ipaReef, a+".ˈʔ"+a, a+".ˈ"+a)
			noGlottalStopIPA = strings.ReplaceAll(noGlottalStopIPA, a+".ʔ"+a, a+"."+a)
			noGlottalStopIPA = strings.ReplaceAll(noGlottalStopIPA, a+"ʔ."+a, a+"."+a)
			noGlottalStopIPA = strings.ReplaceAll(noGlottalStopIPA, a+"ʔ.ˈ"+a, a+".ˈ"+a)

			ipaReef += "] or [" + noGlottalStopIPA
		}
	}

	// fix breakdown
	shortString = strings.ReplaceAll(breakdown, "-", "")
	for _, a := range []string{"a", "e", "ì", "o", "u", "i", "ä", "ù"} {
		if strings.Contains(shortString, a+"'"+a) {
			noGlottalStopBreakdown := strings.ReplaceAll(breakdown, a+"-'"+a, a+"-"+a)
			noGlottalStopBreakdown = strings.ReplaceAll(noGlottalStopBreakdown, a+"'-"+a, a+"-"+a)

			breakdown += " or " + noGlottalStopBreakdown
		}

	}

	return []string{breakdown, ipaReef}
}
