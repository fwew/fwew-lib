package fwew_lib

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
