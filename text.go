package fwew_lib

import (
	"os/user"
	"path/filepath"
)

var texts = map[string]string{}

const (
	dictFileName = "dictionary-v2.txt"
	mdBold       = "**"
	mdItalic     = "*"
	newline      = "\n"
	valNull      = "NULL"
	space        = " "
)

func init() {
	currentUser, _ := user.Current()

	// main program strings
	texts["name"] = "fwew"

	// List
	texts["listHelpEN"] = `Help for /list Command aka Fwew List Filter Expressions (LFEs)

Part of Speech
want a list of nouns? verbs? adjectives?

  pos <condition> <your_text>

  <condition>:
    starts not-starts
    ends   not-ends
    is     not-is 
    has    not-has
    like   not-like

Word Characteristics
want words that start/end with or have certain letters?

  word <condition> <your_text>

  <condition>: 
    starts starts-any starts-all starts-none
    ends   ends-any   ends-all   ends-none
    has    has-any    has-all    has-none
    like   like-any   like-all   like-none

    not-starts   matches
    not-ends
    not-has
    not-like

Chronology
want the last 20 words that came out? the first 10 words that came out?

  words first <your_number>
  words last <your_number>

Syllable Count
want 1-syllable words? words with more than 2?

  syllables <comparison> <your_number>

Stressed Syllable Location
want words that are stressed on the first syllable? last? 3rd one?

  stress <comparison> <your_number>

Word Length
want 5-letter words? note: "kx", "ng", etc. each count as 1 "Na'vi letter"

  length <comparison> <your_number>

<comparison>:
  <  (less than)
  <= (less than or equal to)
  =  (equal to)
  >= (greater than or equal to)
  >  (greater than)
  != (not equal to)`

	texts["/listDesc"] = `list all words that meet given criteria`
	texts["/listUsage"] = texts["listHelpEN"]
	texts["/listExample"] = "list syllables = 3 and pos has vtr."

	// Random
	texts["/randomDesc"] = "show given <number> of random entries. <what>, <cond>, and <spec> work the same way as with /list"
	texts["/randomUsage"] = "random <number> [where <what> <cond> <spec> [and <what> <cond> <spec> ...]]"
	texts["/randomExample"] = "random 5 where pos is n."

	// More List & Random
	// <what> strings
	texts["w_pos"] = "pos"
	texts["w_word"] = "word"
	texts["w_words"] = "words"
	texts["w_syllables"] = "syllables"
	texts["w_stress"] = "stress"
	texts["w_length"] = "length"
	// <cond> strings
	texts["c_is"] = "is"
	texts["c_has"] = "has"
	texts["c_has-any"] = "has-any"
	texts["c_has-all"] = "has-all"
	texts["c_has-none"] = "has-none"
	texts["c_like"] = "like"
	texts["c_like-any"] = "like-any"
	texts["c_like-all"] = "like-all"
	texts["c_like-none"] = "like-none"
	texts["c_starts"] = "starts"
	texts["c_starts-any"] = "starts-any"
	texts["c_starts-all"] = "starts-all"
	texts["c_starts-none"] = "starts-none"
	texts["c_ends"] = "ends"
	texts["c_ends-any"] = "ends-any"
	texts["c_ends-all"] = "ends-all"
	texts["c_ends-none"] = "ends-none"
	texts["c_not-is"] = "not-is"
	texts["c_not-has"] = "not-has"
	texts["c_not-like"] = "not-like"
	texts["c_not-starts"] = "not-starts"
	texts["c_not-ends"] = "not-ends"
	texts["c_first"] = "first"
	texts["c_last"] = "last"
	texts["c_matches"] = "matches"

	// file strings
	texts["homeDir"], _ = filepath.Abs(currentUser.HomeDir)
	texts["dataDir"] = filepath.Join(texts["homeDir"], ".fwew")
	texts["dictURL"] = "https://tirea.learnnavi.org/dictionarydata/" + dictFileName

	// general message strings
	texts["src"] = "source"
}

// Text function is the accessor for []string texts
func Text(s string) string {
	if _, ok := texts[s]; ok {
		return texts[s]
	}
	return TextNotFound.Error() + ": " + s
}

var messageNonNaviLetters = map[string]string{
	"en": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // English
	// TODO
	"de": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // German (Deutsch)
	// TODO
	"es": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // Spanish (Español)
	// TODO
	"et": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // Estonian (Eesti)
	// TODO
	"fr": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // French (Français)
	// TODO
	"hu": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`",       // Hungarian (Magyar)
	"it": "**{oldWord}** Presenta lettere non in Na'vi: `{nonNaviLetters}`",  // Italian (Italiano)
	"ko": "**{oldWord}**에는 나비어에 존재하지 않는 낱말이 포함되어 있습니다. - `{nonNaviLetters}`", // Korean (한국어)
	// TODO
	"nl": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // Dutch (Nederlands)
	// TODO
	"pl": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // Polish (Polski)
	// TODO
	"pt": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // Portuguese (Português)
	// TODO
	"ru": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // Russian (Русский)
	// TODO
	"sv": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // Swedish (Svenska)
	// TODO
	"tr": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // Turkish (Türkçe)
	// TODO
	"uk": "**{oldWord}** Has letters not in Na'vi: `{nonNaviLetters}`", // Ukrainian (Українська)
}

var messageNoNuclei = map[string]string{
	"en": "**{oldWord}** Error: could not find any syllable nuclei", // English
	// TODO
	"de": "**{oldWord}** Error: could not find any syllable nuclei", // German (Deutsch)
	// TODO
	"es": "**{oldWord}** Error: could not find any syllable nuclei", // Spanish (Español)
	// TODO
	"et": "**{oldWord}** Error: could not find any syllable nuclei", // Estonian (Eesti)
	// TODO
	"fr": "**{oldWord}** Error: could not find any syllable nuclei", // French (Français)
	// TODO
	"hu": "**{oldWord}** Error: could not find any syllable nuclei",        // Hungarian (Magyar)
	"it": "**{oldWord}** Errore: non si è trovato alcuno nucleo sillabico", // Italian (Italiano)
	"ko": "**{oldWord}**에서 음절핵(중성)에 해당하는 요소를 찾을 수 없습니다.",                   // Korean (한국어)
	// TODO
	"nl": "**{oldWord}** Error: could not find any syllable nuclei", // Dutch (Nederlands)
	// TODO
	"pl": "**{oldWord}** Error: could not find any syllable nuclei", // Polish (Polski)
	// TODO
	"pt": "**{oldWord}** Error: could not find any syllable nuclei", // Portuguese (Português)
	// TODO
	"ru": "**{oldWord}** Error: could not find any syllable nuclei", // Russian (Русский)
	// TODO
	"sv": "**{oldWord}** Error: could not find any syllable nuclei", // Swedish (Svenska)
	// TODO
	"tr": "**{oldWord}** Error: could not find any syllable nuclei", // Turkish (Türkçe)
	// TODO
	"uk": "**{oldWord}** Error: could not find any syllable nuclei", // Ukrainian (Українська)
}

var messageInvalidConsonants = map[string]string{
	"en": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // English
	// TODO
	"de": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // German (Deutsch)
	// TODO
	"es": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // Spanish (Español)
	// TODO
	"et": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // Estonian (Eesti)
	// TODO
	"fr": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // French (Français)
	// TODO
	"hu": "**{oldWord}** Invalid consonant combination: `{badConsonants}`",        // Hungarian (Magyar)
	"it": "**{oldWord}** Combinazione consonantica non valida: `{badConsonants}`", // Italian (Italiano)
	"ko": "**{oldWord}**에 유효하지 않은 조합이 발견되었습니다. - `{badConsonants}`",               // Korean (한국어)
	// TODO
	"nl": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // Dutch (Nederlands)
	// TODO
	"pl": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // Polish (Polski)
	// TODO
	"pt": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // Portuguese (Português)
	// TODO
	"ru": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // Russian (Русский)
	// TODO
	"sv": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // Swedish (Svenska)
	// TODO
	"tr": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // Turkish (Türkçe)
	// TODO
	"uk": "**{oldWord}** Invalid consonant combination: `{badConsonants}`", // Ukrainian (Українська)
}

var messageNeededVowel = map[string]string{
	"en": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // English
	// TODO
	"de": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // German (Deutsch)
	// TODO
	"es": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // Spanish (Español)
	// TODO
	"et": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // Estonian (Eesti)
	// TODO
	"fr": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // French (Français)
	// TODO
	"hu": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`",               // Hungarian (Magyar)
	"it": "**{oldWord}** necessita di una vocale, dittongo o semivocale qui: `{breakdown}`",         // Italian (Italiano)
	"ko": "**{oldWord}** 에 유효하지 않은 자음 조합이 발견되었습니다. 다음 위치에 모음 또는 준모음(음절자음)을 추가해주세요. - `{breakdown}`", // Korean (한국어)
	// TODO
	"nl": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // Dutch (Nederlands)
	// TODO
	"pl": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // Polish (Polski)
	// TODO
	"pt": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // Portuguese (Português)
	// TODO
	"ru": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // Russian (Русский)
	// TODO
	"sv": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // Swedish (Svenska)
	// TODO
	"tr": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // Turkish (Türkçe)
	// TODO
	"uk": "**{oldWord}** Needs a vowel, diphthong or psuedovowel here: `{breakdown}`", // Ukrainian (Українська)
}

var messagePsuedovowelsCantCoda = map[string]string{
	"en": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // English
	// TODO
	"de": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // German (Deutsch)
	// TODO
	"es": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // Spanish (Español)
	// TODO
	"et": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // Estonian (Eesti)
	// TODO
	"fr": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // French (Français)
	// TODO
	"hu": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`",                         // Hungarian (Magyar)
	"it": "**{oldWord}** Le semivocali non hanno mai coda: `{breakdown}`",                        // Italian (Italiano)
	"ko": "**{oldWord}**에 유효하지 않은 자음 조합이 발견되었습니다. 준모음(음절자음)은 말음(종성)을 가질 수 없습니다. - `{breakdown}`", // Korean (한국어)
	// TODO
	"nl": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // Dutch (Nederlands)
	// TODO
	"pl": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // Polish (Polski)
	// TODO
	"pt": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // Portuguese (Português)
	// TODO
	"ru": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // Russian (Русский)
	// TODO
	"sv": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // Swedish (Svenska)
	// TODO
	"tr": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // Turkish (Türkçe)
	// TODO
	"uk": "**{oldWord}** Psuedovowels can't accept codas: `{breakdown}`", // Ukrainian (Українська)
}

var messagePsuedovowelsMustOnset = map[string]string{
	"en": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // English
	// TODO
	"de": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // German (Deutsch)
	// TODO
	"es": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // Spanish (Español)
	// TODO
	"et": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // Estonian (Eesti)
	// TODO
	"fr": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // French (Français)
	// TODO
	"hu": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`",                           // Hungarian (Magyar)
	"it": "**{oldWord}** le semivocali devono avere un inizio: `{breakdown}`",                    // Italian (Italiano)
	"ko": "**{oldWord}**에 유효하지 않은 자음 조합이 발견되었습니다. 준모음(음절자음)은 반드시 두음(초성)이 필요합니다. - `{breakdown}`", // Korean (한국어)
	// TODO
	"nl": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // Dutch (Nederlands)
	// TODO
	"pl": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // Polish (Polski)
	// TODO
	"pt": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // Portuguese (Português)
	// TODO
	"ru": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // Russian (Русский)
	// TODO
	"sv": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // Swedish (Svenska)
	// TODO
	"tr": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // Turkish (Türkçe)
	// TODO
	"uk": "**{oldWord}** Psuedovowels must have onsets: `{breakdown}`", // Ukrainian (Українська)
}

var messageReefDialect = map[string]string{
	"en": " (In reef dialect.  Forest dialect {breakdown})", // English
	// TODO
	"de": " (In reef dialect.  Forest dialect {breakdown})", // German (Deutsch)
	// TODO
	"es": " (In reef dialect.  Forest dialect {breakdown})", // Spanish (Español)
	// TODO
	"et": " (In reef dialect.  Forest dialect {breakdown})", // Estonian (Eesti)
	// TODO
	"fr": " (In reef dialect.  Forest dialect {breakdown})", // French (Français)
	// TODO
	"hu": " (In reef dialect.  Forest dialect {breakdown})",                  // Hungarian (Magyar)
	"it": " (Nel dialetto del reef. Nel dialetto della foresta {breakdown})", // Italian (Italiano)
	"ko": " (산호초 방언 한정 - 숲 방언: {breakdown})",                                 // Korean (한국어)
	// TODO
	"nl": " (In reef dialect.  Forest dialect {breakdown})", // Dutch (Nederlands)
	// TODO
	"pl": " (In reef dialect.  Forest dialect {breakdown})", // Polish (Polski)
	// TODO
	"pt": " (In reef dialect.  Forest dialect {breakdown})", // Portuguese (Português)
	// TODO
	"ru": " (In reef dialect.  Forest dialect {breakdown})", // Russian (Русский)
	// TODO
	"sv": " (In reef dialect.  Forest dialect {breakdown})", // Swedish (Svenska)
	// TODO
	"tr": " (In reef dialect.  Forest dialect {breakdown})", // Turkish (Türkçe)
	// TODO
	"uk": " (In reef dialect.  Forest dialect {breakdown})", // Ukrainian (Українська)
}

var messageValid = map[string]string{
	"en": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // English
	// TODO
	"de": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // German (Deutsch)
	// TODO
	"es": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // Spanish (Español)
	// TODO
	"et": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // Estonian (Eesti)
	// TODO
	"fr": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // French (Français)
	// TODO
	"hu": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // Hungarian (Magyar)
	"it": "**{oldWord}** Valida: `{breakdown}` con {syllable_count} sillabe {syllable_forest}",   // Italian (Italiano)
	"ko": "**{oldWord}**는 `{breakdown}`의 {syllable_count}음절로 구성된 유효한 단어입니다. {syllable_forest}",   // Korean (한국어)
	// TODO
	"nl": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // Dutch (Nederlands)
	// TODO
	"pl": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // Polish (Polski)
	// TODO
	"pt": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // Portuguese (Português)
	// TODO
	"ru": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // Russian (Русский)
	// TODO
	"sv": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // Swedish (Svenska)
	// TODO
	"tr": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // Turkish (Türkçe)
	// TODO
	"uk": "**{oldWord}** Valid: `{breakdown}` with {syllable_count} syllables {syllable_forest}", // Ukrainian (Українська)
}

var messageIdenticalAdjacentLetters = map[string]string{
	"en": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // English
	// TODO
	"de": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // German (Deutsch)
	// TODO
	"es": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Spanish (Español)
	// TODO
	"et": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Estonian (Eesti)
	// TODO
	"fr": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // French (Français)
	// TODO
	"hu": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Hungarian (Magyar)
	"it": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Italian (Italiano)
	"ko": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Korean (한국어)
	// TODO
	"nl": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Dutch (Nederlands)
	// TODO
	"pl": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Polish (Polski)
	// TODO
	"pt": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Portuguese (Português)
	// TODO
	"ru": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Russian (Русский)
	// TODO
	"sv": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Swedish (Svenska)
	// TODO
	"tr": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Turkish (Türkçe)
	// TODO
	"uk": " (Warning: identical adjacent similar letters are awkward in forest Na'vi)", // Ukrainian (Українська)
}

var messagePsuedovowelAndConsonant = map[string]string{
	"en": " (Warning: a consonant like the previous psuedovowel is awkward)", // English
	// TODO
	"de": " (Warning: a consonant like the previous psuedovowel is awkward)", // German (Deutsch)
	// TODO
	"es": " (Warning: a consonant like the previous psuedovowel is awkward)", // Spanish (Español)
	// TODO
	"et": " (Warning: a consonant like the previous psuedovowel is awkward)", // Estonian (Eesti)
	// TODO
	"fr": " (Warning: a consonant like the previous psuedovowel is awkward)", // French (Français)
	// TODO
	"hu": " (Warning: a consonant like the previous psuedovowel is awkward)", // Hungarian (Magyar)
	"it": " (Warning: a consonant like the previous psuedovowel is awkward)", // Italian (Italiano)
	"ko": " (Warning: a consonant like the previous psuedovowel is awkward)", // Korean (한국어)
	// TODO
	"nl": " (Warning: a consonant like the previous psuedovowel is awkward)", // Dutch (Nederlands)
	// TODO
	"pl": " (Warning: a consonant like the previous psuedovowel is awkward)", // Polish (Polski)
	// TODO
	"pt": " (Warning: a consonant like the previous psuedovowel is awkward)", // Portuguese (Português)
	// TODO
	"ru": " (Warning: a consonant like the previous psuedovowel is awkward)", // Russian (Русский)
	// TODO
	"sv": " (Warning: a consonant like the previous psuedovowel is awkward)", // Swedish (Svenska)
	// TODO
	"tr": " (Warning: a consonant like the previous psuedovowel is awkward)", // Turkish (Türkçe)
	// TODO
	"uk": " (Warning: a consonant like the previous psuedovowel is awkward)", // Ukrainian (Українська)
}

var messageTooBig = map[string]string{
	"en": "⛔ (stopped at {count}. 2000 Character limit)", // English
	// TODO
	"de": "⛔ (stopped at {count}. 2000 Character limit)", // German (Deutsch)
	// TODO
	"es": "⛔ (stopped at {count}. 2000 Character limit)", // Spanish (Español)
	// TODO
	"et": "⛔ (stopped at {count}. 2000 Character limit)", // Estonian (Eesti)
	// TODO
	"fr": "⛔ (stopped at {count}. 2000 Character limit)", // French (Français)
	// TODO
	"hu": "⛔ (stopped at {count}. 2000 Character limit)",    // Hungarian (Magyar)
	"it": "⛔ (fermato a {count}. Limite di 2000 caratteri)", // Italian (Italiano)
	// TODO
	"ko": "⛔ 	(출력값 초과: {count} - 최대 2000개의 결과까지 출력 가능합니다.)", // Korean (한국어)
	// TODO
	"nl": "⛔ (stopped at {count}. 2000 Character limit)", // Dutch (Nederlands)
	// TODO
	"pl": "⛔ (stopped at {count}. 2000 Character limit)", // Polish (Polski)
	// TODO
	"pt": "⛔ (stopped at {count}. 2000 Character limit)", // Portuguese (Português)
	// TODO
	"ru": "⛔ (stopped at {count}. 2000 Character limit)", // Russian (Русский)
	// TODO
	"sv": "⛔ (stopped at {count}. 2000 Character limit)", // Swedish (Svenska)
	// TODO
	"tr": "⛔ (stopped at {count}. 2000 Character limit)", // Turkish (Türkçe)
	// TODO
	"uk": "⛔ (stopped at {count}. 2000 Character limit)", // Ukrainian (Українська)
}
