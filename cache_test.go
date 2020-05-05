package fwew_lib

import (
	"testing"
)

func Test_cacheDict(t *testing.T) {
	// cache dict and test only one entry

	//4	en	'ampi	ˈʔ·am.p·i	'<0><1>amp<2>i	vtr.	touch	ASG; http://naviteri.org/2012/11/renu-ayinanfyaya-the-senses-paradigm/ (27 November 2012)	1	'am-pi	'.amp.i
	// 4	de	'ampi	ˈʔ·am.p·i	'<0><1>amp<2>i	vtr.	berühren	ASG; http://naviteri.org/2012/11/renu-ayinanfyaya-the-senses-paradigm/ (27 November 2012)	1	'am-pi	'.amp.i
	word := Word{
		ID:             "4",
		LangCode:       "de",
		Navi:           "'ampi",
		IPA:            "ˈʔ·am.p·i",
		InfixLocations: "'<0><1>amp<2>i",
		PartOfSpeech:   "vtr.",
		Definition:     "berühren",
		Source:         "ASG; http://naviteri.org/2012/11/renu-ayinanfyaya-the-senses-paradigm/ (27 November 2012)",
		Stressed:       "1",
		Syllables:      "'am-pi",
		InfixDots:      "'.amp.i",
		Affixes:        nil,
		Attempt:        "",
	}

	err := CacheDict()
	if err != nil {
		t.Errorf("Error caching Dictionary!!")
	}
	entry := dictionary["de"][0]
	if !word.Equals(entry) {
		t.Errorf("Read wrong word from cache:\n"+
			"Id: \"%s\" == \"%s\"\n"+
			"LangCode: \"%s\" == \"%s\"\n"+
			"Navi: \"%s\" == \"%s\"\n"+
			"Attempt: \"%s\" == \"%s\"\n"+
			"IPA: \"%s\" == \"%s\"\n"+
			"InfixLocations: \"%s\" == \"%s\"\n"+
			"PartOfSpeech: \"%s\" == \"%s\"\n"+
			"Definition: \"%s\" == \"%s\"\n"+
			"Source: \"%s\" == \"%s\"\n"+
			"Stressed: \"%s\" == \"%s\"\n"+
			"Syllables: \"%s\" == \"%s\"\n"+
			"InfixDots: \"%s\" == \"%s\"\n"+
			"Affixes: \"%s\" == \"%s\"",
			word.ID, entry.ID,
			word.LangCode, entry.LangCode,
			word.Navi, entry.Navi,
			word.Attempt, entry.Attempt,
			word.IPA, entry.IPA,
			word.InfixLocations, entry.InfixLocations,
			word.PartOfSpeech, entry.PartOfSpeech,
			word.Definition, entry.Definition,
			word.Source, entry.Source,
			word.Stressed, entry.Stressed,
			word.Syllables, entry.Syllables,
			word.InfixDots, entry.InfixDots,
			word.Affixes, entry.Affixes)
	}
}
