package fwew_lib

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// List filters the dictionary based on the args.
// args can be empty, if so, the whole Dict will be returned (This also happens if < 3 args are given)
// It will try to always get 3 args and an `and` in between. If less than 3 exist, than it will wil return the previous
// results.
func List(args []string, checkDigraphs uint8) (results []Word, err error) {
	universalLock.Lock()
	defer universalLock.Unlock()
	results, err = GetFullDict()

	if err != nil {
		return
	}

	i := 0
	for i < len(args) {
		args[i] = strings.ReplaceAll(args[i], ",", ", ")
		i++
	}

	for len(args) >= 3 {
		// get 3 args and remove 4th
		simpleArgs := args[0:3]

		results, err = listWords(simpleArgs, results, checkDigraphs)
		if err != nil {
			return
		}

		// remove first 4 elements
		if len(args) > 4 {
			args = args[4:]
		} else {
			break
		}
	}

	return
}

func listWords(args []string, words []Word, checkDigraphs uint8) (results []Word, err error) {
	what := strings.ToLower(args[0])
	wordsLen := len(words)

	for i, word := range words {
		switch what {
		case Text("w_pos"):
			results = filterPos(results, word, args)
		case Text("w_word"):
			results = filterWord(results, word, args, checkDigraphs)
		case Text("w_words"):
			results, err = filterWords(results, word, args, wordsLen, i)
		case Text("w_syllables"):
			results, err = filterNumeric(results, word, args)
		case Text("w_stress"):
			results, err = filterNumeric(results, word, args)
		case Text("w_length"):
			results, err = filterNumeric(results, word, args)
		}
	}

	return
}

func ListHelp(lang string) (helpText string) {
	lang = strings.ToUpper(lang)

	// TODO remove this after other languages are translated
	if lang != "EN" {
		lang = "EN"
	}

	helpText = Text("listHelp" + lang)

	return
}

func filterPos(results []Word, word Word, args []string) []Word {
	var (
		cond = strings.ToLower(args[1])
		spec = strings.ReplaceAll(strings.ToLower(args[2]), ".", "")
		pos  = strings.ReplaceAll(strings.ToLower(word.PartOfSpeech), ".", "")
	)

	condMap := map[string]bool{
		Text("c_starts"):     strings.HasPrefix(pos, spec),
		Text("c_ends"):       strings.HasSuffix(pos, spec),
		Text("c_is"):         pos == spec,
		Text("c_has"):        strings.Contains(pos, spec),
		Text("c_like"):       glob(spec, pos),
		Text("c_not-starts"): !strings.HasPrefix(pos, spec),
		Text("c_not-ends"):   !strings.HasSuffix(pos, spec),
		Text("c_not-is"):     pos != spec,
		Text("c_not-has"):    !strings.Contains(pos, spec),
		Text("c_not-like"):   !glob(spec, pos),
	}

	if condMap[cond] {
		return append(results, word)
	}

	return results
}

func filterWord(results []Word, word Word, args []string, checkDigraphs uint8) []Word {
	var (
		cond = strings.ToLower(args[1])
		spec = preventCompressBug(strings.ToLower(args[2]))
	)

	syllables := word.Syllables
	navi := word.Navi

	switch checkDigraphs {
	case 1: // 1: compress spec and syllables (consider all digraphs)
		spec = compress(strings.ToLower(spec))
		fallthrough
	case 2: // 2: compress syllables, but not spec (find fake digraphs)
		syllables = compress(strings.ToLower(syllables))
		navi = compress(strings.ToLower(navi))
	}

	syllables = strings.ReplaceAll(syllables, "-", "")
	plus := spec[len(spec)-1] == '+'

	condMap := map[string]bool{
		Text("c_starts"):      strings.HasPrefix(syllables, spec),
		Text("c_starts-any"):  satisfiesAny(word, cond, spec),
		Text("c_starts-all"):  satisfiesAll(word, cond, spec),
		Text("c_starts-none"): satisfiesNone(word, cond, spec),
		Text("c_ends"):        plus && strings.HasSuffix(navi, spec) || strings.HasSuffix(syllables, spec),
		Text("c_ends-any"):    satisfiesAny(word, cond, spec),
		Text("c_ends-all"):    satisfiesAll(word, cond, spec),
		Text("c_ends-none"):   satisfiesNone(word, cond, spec),
		Text("c_has"):         plus && strings.Contains(navi, spec) || strings.Contains(syllables, spec),
		Text("c_has-any"):     satisfiesAny(word, cond, spec),
		Text("c_has-all"):     satisfiesAll(word, cond, spec),
		Text("c_has-none"):    satisfiesNone(word, cond, spec),
		Text("c_like"):        glob(spec, syllables),
		Text("c_like-any"):    satisfiesAny(word, cond, spec),
		Text("c_like-all"):    satisfiesAll(word, cond, spec),
		Text("c_like-none"):   satisfiesNone(word, cond, spec),
		Text("c_not-starts"):  !strings.HasPrefix(syllables, spec),
		Text("c_not-ends"):    !strings.HasSuffix(syllables, spec),
		Text("c_not-has"):     plus && !strings.Contains(navi, spec) || !strings.Contains(syllables, spec),
		Text("c_not-like"):    !glob(spec, syllables),
		Text("c_matches"):     spec != "+" && regexp.MustCompile(spec).MatchString(navi),
	}

	if condMap[cond] {
		return append(results, word)
	}

	return results
}

func filterWords(results []Word, word Word, args []string, wordsLen, index int) (filtered []Word, err error) {
	var (
		cond = args[1]
		spec = args[2]
	)

	specNumber, err1 := strconv.Atoi(spec)
	if err1 != nil {
		err = InvalidNumber.wrap(err1)
		return
	}

	condMap := map[string]bool{
		Text("c_first"): index < specNumber,
		Text("c_last"):  index >= wordsLen-specNumber && index <= wordsLen,
	}

	if condMap[cond] {
		filtered = append(results, word)
		return
	}

	filtered = results
	return
}

func filterNumeric(results []Word, word Word, args []string) (filtered []Word, err error) {
	var (
		what = args[0]
		cond = args[1]
		spec = args[2]
	)

	ispec, err1 := strconv.Atoi(spec)
	if err1 != nil {
		err = InvalidNumber.wrap(err1)
		return
	}
	if ispec < 0 {
		syllDash := strings.ReplaceAll(word.Syllables, " ", "-")
		syllArr := strings.Split(syllDash, "-")
		ispec += len(syllArr) + 1
	}

	istress, err2 := strconv.Atoi(word.Stressed)
	if err2 != nil {
		err = InvalidNumber.wrap(err2)
		return
	}

	whatMap := map[string]int{
		Text("w_syllables"): word.SyllableCount(),
		Text("w_stress"):    istress,
		Text("w_length"):    utf8.RuneCountInString(compress(strings.ToLower(word.Syllables))),
	}

	condMap := map[string]bool{
		"<":  whatMap[what] < ispec,
		"<=": whatMap[what] <= ispec,
		"=":  whatMap[what] == ispec,
		">=": whatMap[what] >= ispec,
		">":  whatMap[what] > ispec,
		"!=": whatMap[what] != ispec,
	}

	if condMap[cond] {
		filtered = append(results, word)
		return
	}

	filtered = results
	return
}

func preventCompressBug(input string) string {
	// Be sure nothing can contaminate the data to compress
	removeChars := []string{"q", "b", "d", "c", "0", "1", "2", "3", "4", "5"}
	for _, char := range removeChars {
		input = strings.ReplaceAll(input, char, ";")
	}

	// We don't want to interrupt any g that's part of ng
	input = strings.ReplaceAll(input, "ng", "[-]")
	input = strings.ReplaceAll(input, "g", ";")
	input = strings.ReplaceAll(input, "[-]", "ng")

	return input
}

func satisfiesAny(word Word, cond, spec string) bool {
	var (
		syllables = word.Syllables
		navi      = word.Navi
		specs     = strings.Split(spec, ",")
		plus      = spec[len(spec)-1] == '+'
	)

	switch cond {
	case Text("c_starts-any"):
		for _, s := range specs {
			if strings.HasPrefix(syllables, s) {
				return true
			}
		}
	case Text("c_ends-any"):
		for _, s := range specs {
			if plus && strings.HasSuffix(navi, s) || strings.HasSuffix(syllables, s) {
				return true
			}
		}
	case Text("c_has-any"):
		for _, s := range specs {
			if plus && strings.Contains(navi, s) || strings.Contains(syllables, s) {
				return true
			}
		}
	case Text("c_like-any"):
		for _, s := range specs {
			if glob(s, syllables) {
				return true
			}
		}
	}

	return false
}

func satisfiesAll(word Word, cond, spec string) bool {
	var (
		syllables = word.Syllables
		navi      = word.Navi
		specs     = strings.Split(spec, ",")
		plus      = spec[len(spec)-1] == '+'
	)

	switch cond {
	case Text("c_starts-all"):
		for _, s := range specs {
			if !strings.HasPrefix(syllables, s) {
				return false
			}
		}
	case Text("c_ends-all"):
		for _, s := range specs {
			if !(plus && strings.HasSuffix(navi, s) || strings.HasSuffix(syllables, s)) {
				return false
			}
		}
	case Text("c_has-all"):
		for _, s := range specs {
			if !(plus && strings.Contains(navi, s) || strings.Contains(syllables, s)) {
				return false
			}
		}
	case Text("c_like-all"):
		for _, s := range specs {
			if !glob(s, syllables) {
				return false
			}
		}
	}

	return true
}

func satisfiesNone(word Word, cond, spec string) bool {
	var (
		syllables = word.Syllables
		navi      = word.Navi
		specs     = strings.Split(spec, ",")
		plus      = spec[len(spec)-1] == '+'
	)

	switch cond {
	case Text("c_starts-none"):
		for _, s := range specs {
			if strings.HasPrefix(syllables, s) {
				return false
			}
		}
	case Text("c_ends-none"):
		for _, s := range specs {
			if plus && strings.HasSuffix(navi, s) || strings.HasSuffix(syllables, s) {
				return false
			}
		}
	case Text("c_has-none"):
		for _, s := range specs {
			if plus && strings.Contains(navi, s) || strings.Contains(syllables, s) {
				return false
			}
		}
	case Text("c_like-none"):
		for _, s := range specs {
			if glob(s, syllables) {
				return false
			}
		}
	}

	return true
}
