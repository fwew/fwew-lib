//	This file is part of Fwew.
//	Fwew is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU General Public License as published by
// 	the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
//	Fwew is distributed in the hope that it will be useful,
//	but WITHOUT ANY WARRANTY; without even implied warranty of
//	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	GNU General Public License for more details.
//
//	You should have received a copy of the GNU General Public License
//	along with Fwew.  If not, see http://gnu.org/licenses/
// cSpell: disable

// Package main contains all the things. affixes_test.go tests fwew.go functions.
package fwew_lib

import (
	"bufio"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	// assure dict, so tests wont fail
	_ = AssureDict()

	// call flag.Parse() here if TestMain uses flags
	if testing.CoverMode() != "" {
		debugMode = true
	}
	os.Exit(m.Run())
}

func wordSimpleEqual(w1a, w2a []Word) bool {
	w1l := len(w1a)
	w2l := len(w2a)

	if w1l != w2l {
		return false
	}

	for j := 0; j < w1l; j++ {
		w1 := w1a[j]
		w2 := w2a[j]

		if w1.ID != w2.ID || w1.Navi != w2.Navi ||
			(!reflect.DeepEqual(w1.Affixes.Prefix, w2.Affixes.Prefix) ||
				!reflect.DeepEqual(w1.Affixes.Infix, w2.Affixes.Infix) ||
				!reflect.DeepEqual(w1.Affixes.Suffix, w2.Affixes.Suffix) ||
				!reflect.DeepEqual(w1.Affixes.Lenition, w2.Affixes.Lenition)) {

			return false
		}
	}

	return true
}

type args struct {
	searchNaviText string
	languageCode   string
}

var naviWords = []struct {
	name string
	args args
	want []Word
}{
	{
		name: "molte",
		args: args{
			searchNaviText: "molte",
		},
		want: []Word{
			{
				ID:   "1120",
				Navi: "mllte",
				Affixes: affix{
					Infix: []string{
						"ol",
					},
				},
			},
		},
	},
	{
		name: "a'aw",
		args: args{
			searchNaviText: "a'aw",
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:   "6224",
				Navi: "'a'aw",
				Affixes: affix{
					Lenition: []string{
						"'a→a",
					},
				},
			},
			{
				ID:   "12",
				Navi: "'aw",
				Affixes: affix{
					Prefix: []string{
						"a",
					},
				},
			},
		},
	},
	{
		name: "apukap",
		args: args{
			searchNaviText: "apukap",
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:   "4316",
				Navi: "pukap",
				Affixes: affix{
					Prefix: []string{
						"a",
					},
				},
			},
		},
	},
	{
		name: "tsatan",
		args: args{
			searchNaviText: "tsatan",
		},
		want: []Word{
			{
				ID:   "160",
				Navi: "atan",
				Affixes: affix{
					Prefix: []string{
						"tsa",
					},
				},
			},
		},
	},
	{
		name: "fìlva",
		args: args{
			searchNaviText: "fìlva",
		},
		want: []Word{
			{
				ID:   "7408",
				Navi: "ìlva",
				Affixes: affix{
					Prefix: []string{
						"fì",
					},
				},
			},
		},
	},
	{
		name: "fratan",
		args: args{
			searchNaviText: "fratan",
		},
		want: []Word{
			{
				ID:   "160",
				Navi: "atan",
				Affixes: affix{
					Prefix: []string{
						"fra",
					},
				},
			},
		},
	},
	{
		name: "pepeveng",
		args: args{
			searchNaviText: "pepeveng",
		},
		want: []Word{
			{
				ID:   "56",
				Navi: "'eveng",
				Affixes: affix{
					Prefix: []string{
						"pepe",
					},
				},
			},
		},
	},
	{
		name: "pepefil",
		args: args{
			searchNaviText: "pepefil",
		},
		want: []Word{
			{
				ID:   "12989",
				Navi: "fil",
				Affixes: affix{
					Prefix: []string{
						"pepe",
					},
				},
			},
			{
				ID:   "6616",
				Navi: "pil",
				Affixes: affix{
					Prefix: []string{
						"pepe",
					},
					Lenition: []string{
						"p→f",
					},
				},
			},
		},
	},
	{
		name: "empty",
		args: args{
			searchNaviText: "",
		},
		want: []Word{},
	}, // empty
	{
		name: "säpeykiyevatsi",
		args: args{
			searchNaviText: "säpeykiyevatsi",
		},
		want: []Word{
			{
				ID:   "1788",
				Navi: "si",
				Affixes: affix{
					Infix: []string{
						"äp",
						"eyk",
						"iyev",
						"ats",
					},
				},
			},
		},
	},
	{
		name: "tseng",
		args: args{
			searchNaviText: "tseng",
		},
		want: []Word{
			{
				ID:      "2380",
				Navi:    "tseng",
				Affixes: affix{},
			},
		},
	}, // simple (no *fixes)
	{
		name: "luyu",
		args: args{
			searchNaviText: "luyu",
		},
		want: []Word{
			{
				ID:   "1044",
				Navi: "lu",
				Affixes: affix{
					Infix: []string{
						"uy",
					},
				},
			},
			{
				ID:   "1044",
				Navi: "lu",
				Affixes: affix{
					Suffix: []string{
						"yu",
					},
				},
			},
		},
	},
	{
		name: "seiyi",
		args: args{
			searchNaviText: "seiyi",
		},
		want: []Word{
			{
				ID:   "1788",
				Navi: "si",
				Affixes: affix{
					Infix: []string{
						"eiy",
					},
				},
			},
		},
	}, // ayi override
	{
		name: "zenuyeke",
		args: args{
			searchNaviText: "zenuyeke",
		},
		want: []Word{
			{
				ID:   "3968",
				Navi: "zenke",
				Affixes: affix{
					Infix: []string{
						"uy",
					},
				},
			},
		},
	}, // zenke `yu` override
	{
		name: "ferfen",
		args: args{
			searchNaviText: "ferfen",
		},
		want: []Word{
			{
				ID:   "464",
				Navi: "frrfen",
				Affixes: affix{
					Infix: []string{
						"er",
					},
				},
			},
		},
	}, // handle <er>rr
	{
		name: "ketsuktswa'",
		args: args{
			searchNaviText: "ketsuktswa'",
		},
		want: []Word{
			{
				ID:      "7352",
				Navi:    "ketsuktswa'",
				Affixes: affix{},
			},
			{
				ID:   "4984",
				Navi: "tswa'",
				Affixes: affix{
					Prefix: []string{
						"ketsuk",
					},
				},
			},
		},
	}, // ke/ketsuk prefix
	{
		name: "tìtusaron",
		args: args{
			searchNaviText: "tìtusaron",
		},
		want: []Word{
			{
				ID:   "2004",
				Navi: "taron",
				Affixes: affix{
					Prefix: []string{
						"tì",
					},
					Infix: []string{
						"us",
					},
				},
			},
		},
	}, // tì prefix
	{
		name: "fayioang",
		args: args{
			searchNaviText: "fayioang",
		},
		want: []Word{
			{
				ID:   "612",
				Navi: "ioang",
				Affixes: affix{
					Prefix: []string{
						"fay",
					},
				},
			},
		},
	},
	{
		name: "tsasoaiä",
		args: args{
			searchNaviText: "tsasoaiä",
		},
		want: []Word{
			{
				ID:   "4804",
				Navi: "soaia",
				Affixes: affix{
					Prefix: []string{
						"tsa",
					},
					Suffix: []string{
						"ä",
					},
				},
			},
		},
	}, // soaiä replacement
	{
		name: "tseyä",
		args: args{
			searchNaviText: "tseyä",
		},
		want: []Word{
			{
				ID:   "5268",
				Navi: "tsaw",
				Affixes: affix{
					Suffix: []string{
						"yä",
					},
				},
			},
		},
	}, // tseyä override
	{
		name: "oey",
		args: args{
			searchNaviText: "oey",
		},
		want: []Word{
			{
				ID:   "1380",
				Navi: "oe",
				Affixes: affix{
					Suffix: []string{
						"y",
					},
				},
			},
		},
	}, // oey override
	{
		name: "ngey",
		args: args{
			searchNaviText: "ngey",
		},
		want: []Word{
			{
				ID:   "1348",
				Navi: "nga",
				Affixes: affix{
					Suffix: []string{
						"y",
					},
				},
			},
		},
	}, // ngey override
	{
		name: "tì'usemä",
		args: args{
			searchNaviText: "tì'usemä",
		},
		want: []Word{
			{
				ID:   "3676",
				Navi: "'em",
				Affixes: affix{
					Prefix: []string{
						"tì",
					},
					Infix: []string{
						"us",
					},
					Suffix: []string{
						"ä",
					},
				},
			},
		},
	},
	{
		name: "wemtswo",
		args: args{
			searchNaviText: "wemtswo",
		},
		want: []Word{
			{
				ID:   "3948",
				Navi: "wem",
				Affixes: affix{
					Suffix: []string{
						"tswo",
					},
				},
			},
		},
	}, // tswo suffix
	{
		name: "tsuknumesì",
		args: args{
			searchNaviText: "tsuknumesì",
		},
		want: []Word{
			{
				ID:   "1340",
				Navi: "nume",
				Affixes: affix{
					Prefix: []string{
						"tsuk",
					},
					Suffix: []string{
						"sì",
					},
				},
			},
		},
	}, // tsuk prefix and some suffix
	{
		name: "tsamungwrr",
		args: args{
			searchNaviText: "tsamungwrr",
		},
		want: []Word{
			{
				ID:   "5268",
				Navi: "tsaw",
				Affixes: affix{
					Suffix: []string{
						"mungwrr",
					},
				},
			},
		},
	},
	{
		name: "tsamsiyu",
		args: args{
			searchNaviText: "tsamsiyu",
		},
		want: []Word{
			{
				ID:   "2344",
				Navi: "tsamsiyu",
			},
		},
	},
	{
		name: "'ueyä",
		args: args{
			searchNaviText: "'ueyä",
		},
		want: []Word{
			{
				ID:   "108",
				Navi: "'u",
				Affixes: affix{
					Suffix: []string{
						"yä", "o",
					},
				},
			},
			{
				ID:   "4180",
				Navi: "'uo",
				Affixes: affix{
					Suffix: []string{
						"yä",
					},
				},
			},
		},
	}, // o -> e vowel shift for pronouns with -yä
	{
		name: "awngeyä",
		args: args{
			searchNaviText: "awngeyä",
		},
		want: []Word{
			{
				ID:   "192",
				Navi: "awnga",
				Affixes: affix{
					Suffix: []string{
						"yä",
					},
				},
			},
		},
	}, // a -> e vowel shift for pronouns with -yä
	{
		name: "fpi",
		args: args{
			searchNaviText: "fpi",
		},
		want: []Word{
			{
				ID:   "428",
				Navi: "fpi+",
			},
		},
	}, // "+" support
	{
		name: "a'awnem",
		args: args{
			searchNaviText: "a'awnem",
		},
		want: []Word{
			{
				ID:   "3676",
				Navi: "'em",
				Affixes: affix{
					Prefix: []string{"a"},
					Infix:  []string{"awn"},
				},
			},
		},
	}, // end-attributed verb with tìftang
	{
		name: "aawnem",
		args: args{
			searchNaviText: "aawnem",
		},
		want: []Word{
			{
				ID:   "3676",
				Navi: "'em",
				Affixes: affix{
					Prefix: []string{"a"},
					Infix:  []string{"awn"},
				},
			},
		},
	}, // end-attributed verb with removed tìftang in reef
	{
		name: "fpuse'a",
		args: args{
			searchNaviText: "fpuse'a",
		},
		want: []Word{
			{
				ID:   "420",
				Navi: "fpe'",
				Affixes: affix{
					Suffix: []string{"a"},
					Infix:  []string{"us"},
				},
			},
		},
	}, // start-attributed verb with tìftang
	{
		name: "fpusea",
		args: args{
			searchNaviText: "fpusea",
		},
		want: []Word{
			{
				ID:   "420",
				Navi: "fpe'",
				Affixes: affix{
					Suffix: []string{"a"},
					Infix:  []string{"us"},
				},
			},
		},
	}, // start-attributed verb with removed tìftang in reef
	{
		name: "tsukruna",
		args: args{
			searchNaviText: "tsukruna",
		},
		want: []Word{
			{
				ID:   "1724",
				Navi: "run",
				Affixes: affix{
					Suffix: []string{"a"},
					Prefix: []string{"tsuk"},
				},
			},
		},
	}, // findable
	{
		name: "atsukrun",
		args: args{
			searchNaviText: "atsukrun",
		},
		want: []Word{
			{
				ID:   "1724",
				Navi: "run",
				Affixes: affix{
					Prefix: []string{"a", "tsuk"},
				},
			},
		},
	}, // findable
	{
		name: "tìngusä'än",
		args: args{
			searchNaviText: "tìngusä'än",
		},
		want: []Word{
			{
				ID:   "9632",
				Navi: "ngä'än",
				Affixes: affix{
					Prefix: []string{"tì"},
					Infix:  []string{"us"},
				},
			},
		},
	}, // suffering
	{
		name: "LosÄntsyelesì",
		args: args{
			searchNaviText: "LosÄntsyelesì",
		},
		want: []Word{
			{
				ID:   "10152",
				Navi: "LosÄntsyelesì",
			},
		},
	}, // Los Angeles
	{
		name: "teyngteri",
		args: args{
			searchNaviText: "teyngteri",
		},
		want: []Word{
			{
				ID:   "2136",
				Navi: "tì'eyng",
				Affixes: affix{
					Suffix: []string{"teri"},
				},
			},
		},
	}, // About the answer
	{
		name: "yaìlä",
		args: args{
			searchNaviText: "yaìlä",
		},
		want: []Word{
			{
				ID:   "2724",
				Navi: "ya",
				Affixes: affix{
					Suffix: []string{"ìlä"},
				},
			},
		},
	}, // Through the air
	{
		name: "yawä",
		args: args{
			searchNaviText: "yawä",
		},
		want: []Word{
			{
				ID:   "2724",
				Navi: "ya",
				Affixes: affix{
					Suffix: []string{"wä"},
				},
			},
		},
	}, // Against the air
	{
		name: "yaftumfa",
		args: args{
			searchNaviText: "yaftumfa",
		},
		want: []Word{
			{
				ID:   "2724",
				Navi: "ya",
				Affixes: affix{
					Suffix: []string{"ftumfa"},
				},
			},
		},
	}, // Out of the air
	{
		name: "hivùm",
		args: args{
			searchNaviText: "hivùm",
		},
		want: []Word{
			{
				ID:   "588",
				Navi: "hum",
				Affixes: affix{
					Infix: []string{"iv"},
				},
			},
		},
	}, // Leave (subjunctive, reef)
	{
		name: "ila",
		args: args{
			searchNaviText: "ila",
		},
		want: []Word{
			{
				ID:   "648",
				Navi: "ìlä+",
			},
		},
	}, // See if it can search words without diacritics
	{
		name: "wayila",
		args: args{
			searchNaviText: "wayila",
		},
		want: []Word{
			{
				ID:   "2692",
				Navi: "way",
				Affixes: affix{
					Suffix: []string{"ìlä"},
				},
			},
		},
	}, // See if it can search words without diacritics
	{
		name: "sangi",
		args: args{
			searchNaviText: "sangi",
		},
		want: []Word{
			{
				ID:   "1788",
				Navi: "si",
				Affixes: affix{
					Infix: []string{"äng"},
				},
			},
		},
	}, // See if it can search words without diacritics
	{
		name: "za'utswo",
		args: args{
			searchNaviText: "za'utswo",
		},
		want: []Word{
			{
				ID:   "2792",
				Navi: "za'u",
				Affixes: affix{
					Suffix: []string{"tswo"},
				},
			},
		},
	}, // tìftangs and suffixes
	{
		name: "heykìmangheiam",
		args: args{
			searchNaviText: "heykìmangheiam",
		},
		want: []Word{
			{
				ID:   "3716",
				Navi: "hangham",
				Affixes: affix{
					Infix: []string{"eyk", "ìm", "ei"},
				},
			},
		},
	}, // eyk ìm ei
	{
		name: "heykìyangheiam",
		args: args{
			searchNaviText: "heykìyangheiam",
		},
		want: []Word{
			{
				ID:   "3716",
				Navi: "hangham",
				Affixes: affix{
					Infix: []string{"eyk", "ìy", "ei"},
				},
			},
		},
	}, // eyk ìy ei
	{
		name: "lìmu",
		args: args{
			searchNaviText: "lìmu",
		},
		want: []Word{
			{
				ID:   "1044",
				Navi: "lu",
				Affixes: affix{
					Infix: []string{"ìm"},
				},
			},
		},
	}, // ìm
	{
		name: "lìmu",
		args: args{
			searchNaviText: "lìmu",
		},
		want: []Word{
			{
				ID:   "1044",
				Navi: "lu",
				Affixes: affix{
					Infix: []string{"ìm"},
				},
			},
		},
	}, // ìy
}
var englishWords = []struct {
	name string
	args args
	want []Word
}{
	{
		name: "Spielzeug",
		args: args{
			searchNaviText: "Spielzeug",
			languageCode:   "de",
		},
		want: []Word{
			{
				ID:   "12989",
				Navi: "fil",
			},
		},
	},
	{
		name: "Zehe",
		args: args{
			searchNaviText: "Zehe",
			languageCode:   "de",
		},
		want: []Word{
			{
				ID:   "4524",
				Navi: "venzek",
			},
		},
	},
	{
		name: "sharp",
		args: args{
			searchNaviText: "sharp",
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:   "8604",
				Navi: "litx",
			},
			{
				ID:   "1608",
				Navi: "pxi",
			},
			{
				ID:   "1616",
				Navi: "pxiut",
			},
		},
	},
	{
		name: "charge",
		args: args{
			searchNaviText: "charge",
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:   "936",
				Navi: "kxll",
			},
			{
				ID:   "5256",
				Navi: "kxll si",
			},
		},
	},
	{
		name: "cloth",
		args: args{
			searchNaviText: "cloth",
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:   "9868",
				Navi: "srä",
			},
		},
	},
	{
		name: "merhaba",
		args: args{
			searchNaviText: "merhaba",
			languageCode:   "tr",
		},
		want: []Word{
			{
				ID:   "692",
				Navi: "kaltxì",
			},
			{
				ID:   "11380",
				Navi: "kxì",
			},
		},
	},
}

var englishNoParen = []struct {
	name string
	args args
	want []Word
}{
	{
		name: "Tier",
		args: args{
			searchNaviText: "Tier",
			languageCode:   "de",
		},
		want: []Word{
			{
				ID:   "612",
				Navi: "ioang",
			},
			{
				ID:   "1440",
				Navi: "pa'li",
			},
			{
				ID:   "2124",
				Navi: "tireaioang",
			},
			{
				ID:   "2744",
				Navi: "yerik",
			},
		},
	},
}

var englishParen = []struct {
	name string
	args args
	want []Word
}{
	{
		name: "Tier",
		args: args{
			searchNaviText: "Tier",
			languageCode:   "de",
		},
		want: []Word{
			{
				ID:   "440",
				Navi: "fpxafaw",
			},
			{
				ID:   "7676",
				Navi: "fwampop",
			},
			{
				ID:   "612",
				Navi: "ioang",
			},
			{
				ID:   "1440",
				Navi: "pa'li",
			},
			{
				ID:   "10704",
				Navi: "seyto",
			},
			{
				ID:   "2124",
				Navi: "tireaioang",
			},
			{
				ID:   "2744",
				Navi: "yerik",
			},
		},
	},
}

func TestTranslateFromNaviCached(t *testing.T) {
	var (
		err1 error
		err2 error
	)

	err1 = CacheDictHash()
	err2 = CacheDictHash2()

	if err1 != nil {
		t.Errorf("TranslateFromNaviCached() Failed to CacheDictHash")
	}
	if err2 != nil {
		t.Errorf("TranslateFromNaviCached() Failed to CacheDictHash2")
	}

	for _, a := range adposuffixes {
		word := "tsun" + a

		if newfix, ok := unreefFixes[a]; ok {
			a = newfix
		}

		affixes := affix{Suffix: []string{a}}
		wordWord := Word{Navi: "tsun", ID: "13353", Affixes: affixes}
		want := []Word{wordWord}
		t.Run(word, func(t *testing.T) {
			got, err := TranslateFromNaviHash(word, true, false)
			if err == nil && word == "" && got != nil {
				t.Errorf("TranslateFromNaviCached() got = %v, want %v", got, want)
			} else if err == nil && len(want) == 0 && len(got) > 0 {
				t.Errorf("TranslateFromNaviCached() got = %v, want %v", got, want)
			} else if err != nil || len(want) > 0 && len(got) > 0 && !wordSimpleEqual(got[0][1:], want) {
				t.Errorf("TranslateFromNaviCached() = %v, want %v", got[0][1:], want)
			}
		})
	}

	for _, tt := range naviWords {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TranslateFromNaviHash(tt.args.searchNaviText, true, false)
			if err == nil && tt.args.searchNaviText == "" && got != nil {
				t.Errorf("TranslateFromNaviCached() got = %v, want %v", got, tt.want)
			} else if err == nil && len(tt.want) == 0 && len(got) > 0 {
				t.Errorf("TranslateFromNaviCached() got = %v, want %v", got, tt.want)
			} else if err != nil || len(tt.want) > 0 && len(got) > 0 && !wordSimpleEqual(got[0][1:], tt.want) {
				t.Errorf("TranslateFromNaviCached() = %v, want %v", got[0][1:], tt.want)
			}
		})
	}

	UncacheHashDict()
	UncacheHashDict2()
}

func BenchmarkTranslateFromNavi(b *testing.B) {
	for _, bm := range naviWords {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = TranslateFromNaviHash(bm.args.searchNaviText, true, false)
			}
		})
	}
}

func BenchmarkTranslateFromNaviCached(b *testing.B) {
	_ = CacheDictHash()
	BenchmarkTranslateFromNavi(b)
	UncacheHashDict()
}

func BenchmarkTranslateFromNaviBig(b *testing.B) {
	open, err := os.Open("misc/random_words.txt")
	if err != nil {
		return
	}
	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			panic(err)
		}
	}(open)

	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		line := scanner.Text()

		b.Run(line, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = TranslateFromNaviHash(line, true, false)
			}
		})
	}
}

func BenchmarkTranslateFromNaviBigCached(b *testing.B) {
	_ = CacheDictHash()
	BenchmarkTranslateFromNaviBig(b)
	UncacheHashDict()
}

func TestTranslateToNaviCached(t *testing.T) {
	var (
		err1 error
		err2 error
	)

	err1 = CacheDictHash()
	err2 = CacheDictHash2()

	if err1 != nil {
		t.Errorf("TranslateToNaviCached() Failed to CacheDictHash")
	}
	if err2 != nil {
		t.Errorf("TranslateToNaviCached() Failed to CacheDictHash2")
	}

	for _, tt := range englishWords {
		t.Run(tt.name, func(t *testing.T) {
			gotResults := TranslateToNaviHash(tt.args.searchNaviText, tt.args.languageCode)
			if !wordSimpleEqual(gotResults[0][1:], tt.want) {
				t.Errorf("TranslateToNavi() = %v, want %v", gotResults[0][1:], tt.want)
			}
		})
	}

	for _, tt := range englishParen {
		t.Run(tt.name, func(t *testing.T) {
			gotResults := TranslateToNaviHash(tt.args.searchNaviText, tt.args.languageCode)
			if !wordSimpleEqual(gotResults[0][1:], tt.want) {
				t.Errorf("TranslateToNavi() = %v, want %v", gotResults[0][1:], tt.want)
			}
		})
	}

	UncacheHashDict()
	UncacheHashDict2()
}

func TestBidirectionalCached(t *testing.T) {
	var (
		err1 error
		err2 error
	)

	err1 = CacheDictHash()
	err2 = CacheDictHash2()

	if err1 != nil {
		t.Errorf("TranslateToNaviCached() Failed to CacheDictHash")
	}
	if err2 != nil {
		t.Errorf("TranslateToNaviCached() Failed to CacheDictHash2")
	}

	for _, tt := range englishWords {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, _ := BidirectionalSearch(tt.args.searchNaviText, true, tt.args.languageCode)
			if !wordSimpleEqual(gotResults[0][1:], tt.want) {
				t.Errorf("TranslateToNavi() = %v, want %v", gotResults[0][1:], tt.want)
			}
		})
	}

	for _, tt := range englishNoParen {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, _ := BidirectionalSearch(tt.args.searchNaviText, true, tt.args.languageCode)
			if !wordSimpleEqual(gotResults[0][1:], tt.want) {
				t.Errorf("TranslateToNavi() = %v, want %v", gotResults[0][1:], tt.want)
			}
		})
	}

	for _, tt := range naviWords {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BidirectionalSearch(tt.args.searchNaviText, true, tt.args.languageCode)
			if err == nil && tt.args.searchNaviText == "" && got != nil {
				t.Errorf("TranslateFromNaviCached() got = %v, want %v", got, tt.want)
			} else if err == nil && len(tt.want) == 0 && len(got) > 0 {
				t.Errorf("TranslateFromNaviCached() got = %v, want %v", got, tt.want)
			} else if err != nil || len(tt.want) > 0 && len(got) > 0 && !wordSimpleEqual(got[0][1:], tt.want) {
				t.Errorf("TranslateFromNaviCached() = %v, want %v", got[0][1:], tt.want)
			}
		})
	}

	UncacheHashDict()
	UncacheHashDict2()
}

func BenchmarkTranslateToNaviBig(b *testing.B) {
	open, err := os.Open("misc/random_words_english.txt")
	if err != nil {
		return
	}
	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			panic(err)
		}
	}(open)

	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		line := scanner.Text()

		b.Run(line, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				TranslateToNaviHash(line, "en")
			}
		})
	}
}

func BenchmarkTranslateToNaviBigCached(b *testing.B) {
	_ = CacheDictHash2()
	BenchmarkTranslateToNaviBig(b)
	UncacheHashDict2()
}

func TestRandom(t *testing.T) {
	type args struct {
		amount int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				amount: 1,
			},
		},
		{
			name: "2",
			args: args{
				amount: 2,
			},
		},
		{
			name: "3",
			args: args{
				amount: 3,
			},
		},
		{
			name: "4",
			args: args{
				amount: 4,
			},
		},
		{
			name: "5",
			args: args{
				amount: 5,
			},
		},
		{
			name: "500",
			args: args{
				amount: 500,
			},
		},
		{
			name: "2000",
			args: args{
				amount: 2000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResults, _ := Random(tt.args.amount, nil, 1); !(len(gotResults) == tt.args.amount) {
				t.Errorf("Random() = %v, want %v", len(gotResults), tt.args.amount)
			}
		})
	}
}

func TestRandomCached(t *testing.T) {
	_ = CacheDictHash()
	TestRandom(t)
	UncacheHashDict()
}
