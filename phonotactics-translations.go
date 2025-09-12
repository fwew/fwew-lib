package fwew_lib

import (
	"strconv"
	"strings"
)

var message_non_navi_letters = map[string]string{
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

var message_no_nuclei = map[string]string{
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

var message_invalid_consonants = map[string]string{
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

var message_needed_vowel = map[string]string{
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

var message_psuedovowels_cant_coda = map[string]string{
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

var message_psuedovowels_must_onset = map[string]string{
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

var message_triple_liquid = map[string]string{
	"en": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // English
	// TODO
	"de": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // German (Deutsch)
	// TODO
	"es": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // Spanish (Español)
	// TODO
	"et": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // Estonian (Eesti)
	// TODO
	"fr": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // French (Français)
	// TODO
	"hu": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // Hungarian (Magyar)
	"it": "**{oldWord}** triple R o L non sono ammesse: `{breakdown}`",  // Italian (Italiano)
	"ko": "**{oldWord}** 연속되는 세개의 R 또는 L은 사용 불가능합니다. - `{breakdown}`",   // Korean (한국어)
	// TODO
	"nl": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // Dutch (Nederlands)
	// TODO
	"pl": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // Polish (Polski)
	// TODO
	"pt": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // Portuguese (Português)
	// TODO
	"ru": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // Russian (Русский)
	// TODO
	"sv": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // Swedish (Svenska)
	// TODO
	"tr": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // Turkish (Türkçe)
	// TODO
	"uk": "**{oldWord}** Triple Rs or Ls aren't allowed: `{breakdown}`", // Ukrainian (Українська)
}

var message_reef_dialect = map[string]string{
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

var message_valid = map[string]string{
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

var message_identical_adjacent_letters = map[string]string{
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

var message_psuedovowel_and_consonant = map[string]string{
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

func valid_message(syllable_count int, lang string) string {
	if lang == "en" {
		if syllable_count == 1 {
			message := strings.ReplaceAll(message_valid[lang], "syllables", "syllable")
			message = strings.ReplaceAll(message, "{syllable_count}", strconv.Itoa(syllable_count))
			return message
		}
	}
	return strings.ReplaceAll(message_valid[lang], "{syllable_count}", strconv.Itoa(syllable_count))
}

var message_too_big = map[string]string{
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
