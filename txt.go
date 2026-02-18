package fwew_lib

import (
	"os/user"
	"path/filepath"
)

var usr, _ = user.Current()
var texts = map[string]string{}

func init() {
	// main program strings
	texts["name"] = "fwew"

	// slash-commands Help
	texts["slashCommandHelp"] = `commands:
/set [options]
  	set given options (separated by space) or show currently set options
  	type "/help" for valid options (use without the '-' prefix)
  	notes:
  	c is a function and is implemented as /config (see below)
  	v is a function and is implemented as /version (see below)
  	f is a function not supported in interactive mode
/unset [options]
  	alias for /set [options]
/<option>
  	shortcut alias for /set <option>
/list <what> <cond> <spec> [and <what> <cond> <spec> ...]
  	list all words that meet given criteria
  	<what> is any one of: pos, word, words, syllables, stress
  	<cond> depends on the <what> used:
  	  <what>    | valid <cond>
  	  ----------|------------------------------------
  	  pos       | any one of: is, has, like
  	  word      | any one of: starts, ends, has, like
  	  words     | any one of: first, last
  	  syllables | any one of: <, <=, =, >=, >
	  stress    | any one of: <, <=, =, >=, >
  	<spec> depends on the <cond> used:
  	  <cond>                       | valid <spec>
  	  -----------------------------|----------------------------
  	  is, has, starts, ends        | any string of letter(s)
  	  <, <=, =, >=, >, first, last | any whole number > 0
  	  like                         | any string of letter(s) and
  	                               |     wildcard percent-sign(s)
/random <number>
/random <number> where <what> <cond> <spec> [and <what> <cond> <spec> ...]
  	show given <number> of random entries
  	<what>, <cond>, and <spec> work the same way as with /list
/random random
/random random where <what> <cond> <spec> [and <what> <cond> <spec> ...]
  	show random number of random entries
  	<what>, <cond>, and <spec> work the same way as with /list
/lenition
  	display the lenition table
/len
  	shortcut alias for /lenition
/update
  	download and update the dictionary file
/config <option> <value>
/config [option=value ...]
  	update the default options in the config file
  	type "/config" to see valid options and their current default values
  	valid values are "true" or "false" for all except Language and PosFilter
  	Language: type "/help" for supported language codes
  	PosFilter: any part of speech abbreviation (including '.' at the end)
  	<option> and <value> are not case-sensitive
/commands
  	show this commands help text
/help
  	show main help text
/version
  	show version information
/exit
  	exit/quit the program (aliases /quit /q /wc)`

	// List
	texts["/listDesc"] = `list all words that meet given criteria`
	texts["/listUsage"] = `list <what> <cond> <spec> [and <what> <cond> <spec> ...]
<what> is any one of: pos, word, words, syllables, stress
<cond> depends on the <what> used:
  <what>    | valid <cond>
  ----------|------------------------------------
  pos       | any one of: is, has, like
  word      | any one of: starts, ends, has, like
  words     | any one of: first, last
  syllables | any one of: <, <=, =, >=, >
  stress    | any one of: <, <=, =, >=, >
<spec> depends on the <cond> used:
  <cond>                 | valid <spec>
  -----------------------|----------------------------
  is, has, starts, ends  | any string of letter(s)
  <, <=, =, >=, >        | any whole number > 0
  first, last            | any whole number > 0
  like                   | any string of letter(s) and
                         |     wildcard percent-sign(s)`
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
	texts["cset"] = "currently set"
	texts["set"] = "set"
	texts["unset"] = "unset"
	texts["pre"] = "Prefixes"
	texts["inf"] = "Infixes"
	texts["suf"] = "Suffixes"
	texts["src"] = "source"
	texts["configSaved"] = "config file successfully updated\n"

	// error message strings
	texts["none"] = "no results\n"
	texts["noTextError"] = "err 0: text not found:"
	texts["noDataError"] = "err 1: failed to open dictionary file (" + texts["dictionary"] + ")"
	texts["fileError"] = "err 2: failed to open configuration file (" + texts["config"] + ")"
	texts["noOptionError"] = "err 3: invalid option"
	texts["invalidIntError"] = "err 4: input must be a decimal integer in range 0 <= n <= 32767 or octal integer in range 0 <= n <= 77777"
	texts["invalidOctalError"] = "err 5: invalid octal integer"
	texts["invalidDecimalError"] = "err 6: invalid decimal integer"
	texts["invalidLanguageError"] = "err 7: invalid language option"
	texts["invalidPOSFilterError"] = "err 8: invalid part of speech filter"
	texts["dictCloseError"] = "err 9: failed to close dictionary file (" + texts["dictionary"] + ")"
	texts["noFileError"] = "err 10: failed to open file"
	texts["fileCloseError"] = "err 11: failed to close input file"
	texts["configSyntaxError"] = "err 12: invalid syntax for config"
	texts["configOptionError"] = "err 13: invalid config option"
	texts["configValueError"] = "err 14: invalid config value for"
	texts["invalidNumericError"] = "err 15: invalid numeric digits"
	texts["downloadError"] = "err 16: could not download dictionary update"
}

// Text function is the accessor for []string texts
func Text(s string) string {
	if _, ok := texts[s]; ok {
		return texts[s]
	}
	return texts["noTextError"] + " " + s
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
