package fwew_lib

import (
	"os/user"
	"path/filepath"
)

var usr, _ = user.Current()
var texts = map[string]string{}

const dictFileName = "dictionary-v2.txt"

func init() {
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
	texts["homeDir"], _ = filepath.Abs(usr.HomeDir)
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

// table of all the possible lenitions
var lenitionTable = [8][2]string{
	{"kx", "k"},
	{"px", "p"},
	{"tx", "t"},
	{"k", "h"},
	{"p", "f"},
	{"ts", "s"},
	{"t", "s"},
	{"'", ""},
}

// short table of all the possible lenitions
var shortLenitionTable = [4][2]string{
	{"kx, px, tx", "k, p, t"},
	{"k, p, t", "h, f, s"},
	{"ts", "s"},
	{"'", ""},
}

// table of all the possible translations of "that"
var thatTable = [9][5]string{
	{"Case", "Noun", "   Clause Wrapper   ", "", ""},
	{" ", " ", "Prox.", "Dist.", "Answer "},
	{"====", "=====", "=====", "======", "========"},
	{"Sub.", "tsaw", "fwa", "tsawa", "teynga  "},
	{"Agt.", "tsal", "fula", "tsala", "teyngla "},
	{"Pat.", "tsat", "futa", "tsata", "teyngta "},
	{"Gen.", "tseyä", "N/A", "N/A", "teyngä  "},
	{"Dat.", "tsar", "fura", "tsara", "teyngra "},
	{"Top.", "tsari", "furia", "tsaria", "teyngria"},
}

// table of all the possible translations of "that"
var otherThats = [9][3]string{
	{"tsa-", "pre.", "that"},
	{"tsa'u", "n.", "that (thing)"},
	{"tsakem", "n.", "that (action)"},
	{"fmawnta", "sbd.", "that news"},
	{"fayluta", "sbd.", "these words"},
	{"tsnì", "sbd.", "that (function word)"},
	{"tsonta", "conj.", "to (with kxìm)"},
	{"kuma/akum", "conj.", "that (as a result)"},
	{"a", "part.", "clause level attributive marker"},
}
