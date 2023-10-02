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
	if !word.Equals(entry[0]) {
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
			word.ID, entry[0].ID,
			word.Navi, entry[0].Navi,
			word.IPA, entry[0].IPA,
			word.InfixLocations, entry[0].InfixLocations,
			word.PartOfSpeech, entry[0].PartOfSpeech,
			word.Source, entry[0].Source,
			word.Stressed, entry[0].Stressed,
			word.Syllables, entry[0].Syllables,
			word.DE, entry[0].DE,
			word.EN, entry[0].EN,
			word.ET, entry[0].ET,
			word.FR, entry[0].FR,
			word.HU, entry[0].HU,
			word.NL, entry[0].NL,
			word.PL, entry[0].PL,
			word.RU, entry[0].RU,
			word.SV, entry[0].SV,
			word.TR, entry[0].TR,
			word.InfixDots, entry[0].InfixDots)
	}
}
