package fwew_lib

import (
	"testing"
)

func Test_cacheDict(t *testing.T) {
	// cache dict and test only one entry

	word := Word{
		ID:             "4",
		Navi:           "'ampi",
		IPA:            "ˈʔ·am.p·i",
		InfixLocations: "'<0><1>amp<2>i",
		PartOfSpeech:   "vtr.",
		Source:         "ASG; http://naviteri.org/2012/11/renu-ayinanfyaya-the-senses-paradigm/ (27 November 2012)",
		Stressed:       "1",
		Syllables:      "'am-pi",
		InfixDots:      "'.amp.i",
		DE:             "berühren",
		EN:             "touch",
		ET:             "katsuma, puutuma",
		FR:             "toucher",
		HU:             "(meg)érint",
		NL:             "aanraken",
		PL:             "dotykać",
		RU:             "трогать",
		SV:             "beröra",
		TR:             "dokunmak",
	}

	err := CacheDictHash()
	if err != nil {
		t.Fatalf("Error caching Dictionary!!")
	}
	entry := dictHash["'ampi"]
	if !word.Equals(entry) {
		t.Errorf("Read wrong word from cache:\n"+
			"Id: \"%s\" == \"%s\"\n"+
			"Navi: \"%s\" == \"%s\"\n"+
			"IPA: \"%s\" == \"%s\"\n"+
			"InfixLocations: \"%s\" == \"%s\"\n"+
			"PartOfSpeech: \"%s\" == \"%s\"\n"+
			"Source: \"%s\" == \"%s\"\n"+
			"Stressed: \"%s\" == \"%s\"\n"+
			"Syllables: \"%s\" == \"%s\"\n"+
			"DE: \"%s\" == \"%s\"\n"+
			"EN: \"%s\" == \"%s\"\n"+
			"ET: \"%s\" == \"%s\"\n"+
			"FR: \"%s\" == \"%s\"\n"+
			"HU: \"%s\" == \"%s\"\n"+
			"NL: \"%s\" == \"%s\"\n"+
			"PL: \"%s\" == \"%s\"\n"+
			"RU: \"%s\" == \"%s\"\n"+
			"SV: \"%s\" == \"%s\"\n"+
			"TR: \"%s\" == \"%s\"\n"+
			"InfixDots: \"%s\" == \"%s\"\n",
			word.ID, entry.ID,
			word.Navi, entry.Navi,
			word.IPA, entry.IPA,
			word.InfixLocations, entry.InfixLocations,
			word.PartOfSpeech, entry.PartOfSpeech,
			word.Source, entry.Source,
			word.Stressed, entry.Stressed,
			word.Syllables, entry.Syllables,
			word.DE, entry.DE,
			word.EN, entry.EN,
			word.ET, entry.ET,
			word.FR, entry.FR,
			word.HU, entry.HU,
			word.NL, entry.NL,
			word.PL, entry.PL,
			word.RU, entry.RU,
			word.SV, entry.SV,
			word.TR, entry.TR,
			word.InfixDots, entry.InfixDots)
	}
}
