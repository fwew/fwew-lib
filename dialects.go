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

	// Unstressed ä becomes e
	ipaSyllables := strings.Split(ipa, ".")
	if len(ipaSyllables) > 1 {
		newIpa := ""
		for _, a := range ipaSyllables {
			newIpa += "."
			if !strings.Contains(a, "ˈ") {
				newIpa += strings.ReplaceAll(a, "æ", "ɛ")
			} else {
				newIpa += a
			}
		}

		ipa = newIpa
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
			if strings.HasPrefix(ipaReef, a) {
				ipaReef = soften[a] + strings.TrimPrefix(ipaReef, a)
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

		temp := ""
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
			temp += string(a)
		}

		ipaReef = temp
	}

	ipaReef = strings.TrimPrefix(ipaReef, ".")

	ipaReef = strings.ReplaceAll(ipaReef, "*.", " ")

	// now Romanize the reef IPA
	word := strings.Split(ipaReef, " ")

	breakdown = ""

	for j := 0; j < len(word); j++ {
		word[j] = strings.ReplaceAll(word[j], "]", "")
		word[j] = strings.ReplaceAll(word[j], "[", "")
		// "or" means there's more than one IPA in this word, and we only want one
		if word[j] == "or" {
			breakdown += "or "
			continue
		}

		syllables := strings.Split(word[j], ".")

		/* Onset */
		for k := 0; k < len(syllables); k++ {
			breakdown += "-"

			stressed := false
			syllable := strings.ReplaceAll(syllables[k], "·", "")
			if strings.Contains(syllable, "ˈ") {
				stressed = true
				breakdown += "__"
			}
			syllable = strings.ReplaceAll(syllable, "ˈ", "")
			syllable = strings.ReplaceAll(syllable, "ˌ", "")

			// tsy
			if strings.HasPrefix(syllable, "tʃ") {
				breakdown += "ch"
				syllable = strings.TrimPrefix(syllable, "tʃ")
			} else if len(syllable) >= 4 && syllable[0:4] == "t͡s" {
				// ts
				breakdown += "ts"
				//tsp
				if hasAt("ptk", syllable, 3) {
					if nthRune(syllable, 4) == "'" {
						// ts + ejective onset
						breakdown += romanization2[syllable[4:6]]
						syllable = syllable[6:]
					} else {
						// ts + unvoiced plosive
						breakdown += romanization2[string(syllable[4])]
						syllable = syllable[5:]
					}
				} else if hasAt("lɾmnŋwj", syllable, 3) {
					// ts + other consonent
					breakdown += romanization2[nthRune(syllable, 3)]
					syllable = syllable[4+len(nthRune(syllable, 3)):]
				} else {
					// ts without a cluster
					syllable = syllable[4:]
				}
			} else if hasAt("fs", syllable, 0) {
				//
				breakdown += nthRune(syllable, 0)
				if hasAt("ptk", syllable, 1) {
					if nthRune(syllable, 2) == "'" {
						// f/s + ejective onset
						breakdown += romanization2[syllable[1:3]]
						syllable = syllable[3:]
					} else {
						// f/s + unvoiced plosive
						breakdown += romanization2[string(syllable[1])]
						syllable = syllable[2:]
					}
				} else if hasAt("lɾmnŋwj", syllable, 1) {
					// f/s + other consonent
					breakdown += romanization2[nthRune(syllable, 1)]
					syllable = syllable[1+len(nthRune(syllable, 1)):]
				} else {
					// f/s without a cluster
					syllable = syllable[1:]
				}
			} else if hasAt("ptk", syllable, 0) {
				if nthRune(syllable, 1) == "'" {
					// ejective
					breakdown += romanization2[syllable[0:2]]
					syllable = syllable[2:]
				} else {
					// unvoiced plosive
					breakdown += romanization2[string(syllable[0])]
					syllable = syllable[1:]
				}
			} else if hasAt("ʔlɾhmnŋvwjzbdg", syllable, 0) {
				// other normal onset
				breakdown += romanization2[nthRune(syllable, 0)]
				syllable = syllable[len(nthRune(syllable, 0)):]
			} else if hasAt("ʃʒ", syllable, 0) {
				// one sound representd as a cluster
				if nthRune(syllable, 0) == "ʃ" {
					breakdown += "sh"
				}
				syllable = syllable[len(nthRune(syllable, 0)):]
			}

			/*
			 * Nucleus
			 */
			if len(syllable) > 1 && hasAt("jw", syllable, 1) {
				//diphthong
				breakdown += romanization2[syllable[0:len(nthRune(syllable, 0))+1]]
				syllable = string([]rune(syllable)[2:])
			} else if len(syllable) > 1 && hasAt("lr", syllable, 0) {
				//psuedovowel
				breakdown += romanization2[syllable[0:3]]
				continue // psuedovowels can't coda
			} else {
				//vowel
				breakdown += romanization2[nthRune(syllable, 0)]
				syllable = string([]rune(syllable)[1:])
			}

			/*
			 * Coda
			 */
			if len(syllable) > 0 {
				if nthRune(syllable, 0) == "s" {
					breakdown += "sss" //oìsss only
				} else {
					if syllable == "k̚" {
						breakdown += "k"
					} else if syllable == "p̚" {
						breakdown += "p"
					} else if syllable == "t̚" {
						breakdown += "t"
					} else if syllable == "ʔ̚" {
						breakdown += "'"
					} else {
						if syllable[0] == 'k' && len(syllable) > 1 {
							breakdown += "kx"
						} else {
							breakdown += romanization2[syllable]
						}
					}
				}
			}

			if stressed {
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
