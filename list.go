package fwew_lib

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Filter the dictionary based on the args.
// args can be empty, if so, the whole Dict will be returned (This also happens if < 3 args are given)
// It will try to always get 3 args and an `and` in between. If less than 3 exist, than it will wil return the previous results.
func List(args []string, checkDigraphs uint8) (results []Word, err error) {
	results, err = GetFullDict()

	if err != nil {
		return
	}

	for len(args) >= 3 {
		// get 3 args and remove 4th
		simpleArgs := args[0:3]

		results, err = listWords(simpleArgs, results, checkDigraphs)
		if err != nil {
			return
		}

		// TODO: check if args[3] is something different than "and" (only if other things will be supported)

		// remove first 4 elements
		if len(args) > 4 {
			args = args[4:]
		} else {
			break
		}
	}

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

func listWords(args []string, words []Word, checkDigraphs uint8) (results []Word, err error) {
	var (
		what = strings.ToLower(args[0])
		cond = strings.ToLower(args[1])
		spec = strings.ToLower(args[2])
	)

	// /list what cond spec
	// /list pos starts v
	// /list pos ends m.
	// /list pos has svin.
	// /list pos is v.
	// /list pos like *
	// /list pos not-starts v
	// /list pos not-ends m.
	// /list pos not-has svin.
	// /list pos not-is v.
	// /list pos not-like *
	// /list word starts ft
	// /list word ends ang
	// /list word has ts
	// /list word like *
	// /list word not-starts ft
	// /list word not-ends ang
	// /list word not-has ts
	// /list word not-like *
	// /list words first 20
	// /list words last 30
	// /list syllables > 1
	// /list syllables = 2
	// /list syllables <= 3
	// /list stress = 1

	wordsLen := len(words)

	switch what {
	case Text("w_word"):
		spec = preventCompressBug(spec)
	}

	for i, word := range words {
		switch what {
		case Text("w_pos"):
			pos := strings.ReplaceAll(word.PartOfSpeech, ".", "")
			pos = strings.ToLower(pos)
			spec = strings.ReplaceAll(spec, ".", "")
			switch cond {
			case Text("c_starts"):
				if strings.HasPrefix(pos, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_ends"):
				if strings.HasSuffix(pos, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_is"):
				if pos == spec {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_has"):
				if strings.Contains(pos, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_like"):
				if Glob(spec, pos) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_not-starts"):
				if !strings.HasPrefix(pos, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_not-ends"):
				if !strings.HasSuffix(pos, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_not-is"):
				if pos != spec {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_not-has"):
				if !strings.Contains(pos, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_not-like"):
				if !Glob(spec, pos) {
					results = AppendAndAlphabetize(results, word)
				}
			}
		case Text("w_word"):
			syll := word.Syllables
			naviWord := word.Navi
			switch checkDigraphs {
			case 1:
				spec = compress(spec)
				fallthrough
			case 2:
				syll = compress(syll)
				naviWord = compress(naviWord)
			}

			syll = strings.ReplaceAll(syll, "-", "")
			plus := false
			if spec[len(spec)-1] == '+' {
				plus = true
			}
			switch cond {
			case Text("c_starts"):
				if strings.HasPrefix(syll, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_ends"):
				if plus && strings.HasSuffix(naviWord, spec) {
					results = AppendAndAlphabetize(results, word)
				} else if strings.HasSuffix(syll, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_has"):
				if plus && strings.HasSuffix(naviWord, spec) {
					results = AppendAndAlphabetize(results, word)
				} else if strings.Contains(syll, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_like"):
				if Glob(spec, syll) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_not-starts"):
				if !strings.HasPrefix(syll, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_not-ends"):
				if !strings.HasSuffix(syll, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_not-has"):
				if plus && !strings.HasSuffix(naviWord, spec) {
					results = AppendAndAlphabetize(results, word)
				} else if !strings.Contains(syll, spec) {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_not-like"):
				if !Glob(spec, syll) {
					results = AppendAndAlphabetize(results, word)
				}
			}
		case Text("w_words"):
			specNumber, err1 := strconv.Atoi(spec)
			if err1 != nil {
				log.Printf("%s (%s)\n", Text("invalidNumericError"), spec)
				err = InvalidNumber.wrap(err1)
				return
			}
			switch cond {
			case Text("c_first"):
				if i < specNumber {
					results = AppendAndAlphabetize(results, word)
				}
			case Text("c_last"):
				if i >= wordsLen-specNumber && i <= wordsLen {
					results = AppendAndAlphabetize(results, word)
				}
			}
		case Text("w_syllables"):
			ispec, err1 := strconv.Atoi(spec)
			if err1 != nil {
				fmt.Println(Text("invalidDecimalError"))
				err = InvalidNumber.wrap(err1)
				return
			}
			switch cond {
			case "<":
				if word.SyllableCount() < ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case "<=":
				if word.SyllableCount() <= ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case "=":
				if word.SyllableCount() == ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case ">=":
				if word.SyllableCount() >= ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case ">":
				if word.SyllableCount() > ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case "!=":
				if word.SyllableCount() != ispec {
					results = AppendAndAlphabetize(results, word)
				}
			}
		case Text("w_stress"):
			ispec, err1 := strconv.Atoi(spec)
			if err1 != nil {
				fmt.Println(Text("invalidDecimalError"))
				err = InvalidNumber.wrap(err1)
				return
			}
			istress, err1 := strconv.Atoi(word.Stressed)
			if err1 != nil {
				fmt.Println(Text("invalidDecimalError"))
				err = InvalidNumber.wrap(err1)
				return
			}
			switch cond {
			case "<":
				if istress < ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case "<=":
				if istress <= ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case "=":
				if ispec < 0 {
					if word.SyllableCount()+ispec+1 == istress {
						results = AppendAndAlphabetize(results, word)
					}
				} else if istress == ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case ">=":
				if istress >= ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case ">":
				if istress > ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case "!=":
				if ispec < 0 {
					if word.SyllableCount()+ispec+1 != istress {
						results = AppendAndAlphabetize(results, word)
					}
				} else if istress != ispec {
					results = AppendAndAlphabetize(results, word)
				}
			}
		case Text("w_length"):
			ispec, err1 := strconv.Atoi(spec)
			if err1 != nil {
				fmt.Println(Text("invalidDecimalError"))
				err = InvalidNumber.wrap(err1)
				return
			}
			ilength := utf8.RuneCountInString(compress(word.Syllables))
			switch cond {
			case "<":
				if ilength < ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case "<=":
				if ilength <= ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case "=":
				if ispec < 0 {
					if word.SyllableCount()+ispec+1 == ilength {
						results = AppendAndAlphabetize(results, word)
					}
				} else if ilength == ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case ">=":
				if ilength >= ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case ">":
				if ilength > ispec {
					results = AppendAndAlphabetize(results, word)
				}
			case "!=":
				if ispec < 0 {
					if word.SyllableCount()+ispec+1 != ilength {
						results = AppendAndAlphabetize(results, word)
					}
				} else if ilength != ispec {
					results = AppendAndAlphabetize(results, word)
				}
			}
		}
	}

	return
}
