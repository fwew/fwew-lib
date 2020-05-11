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

// Package main contains all the things. affixes_test.go tests fwew.go functions.
package fwew_lib

import (
	"bufio"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	if testing.CoverMode() != "" {
		debugMode = true
	}
	os.Exit(m.Run())
}

/*
import (
	"flag"
	"testing"
)

func TestSimilarity(t *testing.T) {
	f0 := similarity("fmetok", "fmetok")
	if f0 != 1.0 {
		t.Errorf("Wanted %f, Got %f", 1.0, f0)
	}

	f1 := similarity("meoauniaea", "eltu")
	if f1 != 0.0 {
		t.Errorf("Wanted %f, Got %f", 0.0, f1)
	}
}

// helper function for TestFwew, basically a means to consider two Word structs equal
func testEqualWord(w0, w1 Word) bool {
	if w0.ID == w1.ID && w0.Navi == w1.Navi {
		return true
	}
	return false
}

func TestFwew(t *testing.T) {
	// Set relevant option flags
	configuration = ReadConfig()
	reverse = flag.Bool("r", false, Text("usageR"))
	language = flag.String("l", configuration.Language, Text("usageL"))
	posFilter = flag.String("p", configuration.PosFilter, Text("usageP"))
	useAffixes = flag.Bool("a", configuration.UseAffixes, Text("usageA"))
	flag.Parse()

	var w Word

	w0 := fwew("fmetok")[0]
	w = Word{ID: "392", Navi: "fmetok"}
	if !testEqualWord(w, w0) {
		t.Errorf("Wanted %s, Got %s\n", w, w0)
	}

	w1 := fwew("")
	if w1 != nil {
		t.Errorf("empty string did not yield empty Word slice\n")
	}

	w2 := fwew("tseyä")[0]
	w = Word{ID: "5268", Navi: "tsaw"}
	// if w3.ID != "5268" && w3.Navi != "tsaw" {
	if !testEqualWord(w, w2) {
		t.Errorf("Wanted %s, Got %s\n", w, w2)
	}

	w5 := fwew("oey")[0]
	w = Word{ID: "1380", Navi: "oe"}
	if !testEqualWord(w, w5) {
		t.Errorf("Wanted %s, Got %s\n", w, w5)
	}

	w6 := fwew("ngey")[0]
	w = Word{ID: "1348", Navi: "nga"}
	if !testEqualWord(w, w6) {
		t.Errorf("Wanted %s, Got %s\n", w, w6)
	}

	*reverse = true
	w7 := fwew("test")[0]
	w = Word{ID: "392", Navi: "fmetok"}
	if !testEqualWord(w, w7) {
		t.Errorf("Wanted %s, Got %s\n", w, w7)
	}

	*useAffixes = false
	*reverse = false
	w8 := fwew("fmetok")
	if len(w8) != 1 {
		t.Errorf("Wanted 1 word, Got %d\n", len(w8))
	}
}

// helper function for TestSyllableCount, basically cuts down on repetitive code
func testEqualInt(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Errorf("Wanted %d, Got %d\n", expected, actual)
	}
}

func TestSyllableCount(t *testing.T) {
	var w Word

	w = Word{Navi: "nari si"}
	testEqualInt(t, 3, syllableCount(w))

	w = Word{Navi: "lu"}
	testEqualInt(t, 1, syllableCount(w))

	w = Word{Navi: "ätxäle si"}
	testEqualInt(t, 4, syllableCount(w))

	w = Word{Navi: "tireapängkxo"}
	testEqualInt(t, 5, syllableCount(w))

	w = Word{Navi: "tìng tseng"}
	testEqualInt(t, 2, syllableCount(w))
}
*/

func wordSimpleEqual(w1a, w2a []Word, checkAffixes bool) bool {
	w1l := len(w1a)
	w2l := len(w2a)

	if w1l != w2l {
		return false
	}

	for j := 0; j < w1l; j++ {
		w1 := w1a[j]
		w2 := w2a[j]

		if w1.ID != w2.ID ||
			w1.LangCode != w2.LangCode ||
			w1.Navi != w2.Navi ||
			(checkAffixes &&
				(!reflect.DeepEqual(w1.Affixes.Prefix, w2.Affixes.Prefix) ||
					!reflect.DeepEqual(w1.Affixes.Infix, w2.Affixes.Infix) ||
					!reflect.DeepEqual(w1.Affixes.Suffix, w2.Affixes.Suffix) ||
					!reflect.DeepEqual(w1.Affixes.Lenition, w2.Affixes.Lenition))) {

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
			languageCode:   "de",
		},
		want: []Word{
			{
				ID:       "1120",
				LangCode: "de",
				Navi:     "mllte",
				Affixes: affix{
					Infix: []string{
						"ol",
					},
				},
			},
		},
	},
	{
		name: "pepfil",
		args: args{
			searchNaviText: "pepfil",
			languageCode:   "fr",
		},
		want: []Word{
			{
				ID:       "6616",
				LangCode: "fr",
				Navi:     "pil",
				Affixes: affix{
					Prefix: []string{
						"pep",
					},
					Lenition: []string{
						"p→f",
					},
				},
			},
			{
				ID:       "12989",
				LangCode: "fr",
				Navi:     "fil",
				Affixes: affix{
					Prefix: []string{
						"pep",
					},
				},
			},
		},
	},
	{
		name: "empty",
		args: args{
			searchNaviText: "",
			languageCode:   "de",
		},
		want: []Word{},
	}, // empty
	{
		name: "säpeykiyevatsi",
		args: args{
			searchNaviText: "säpeykiyevatsi",
			languageCode:   "hu",
		},
		want: []Word{
			{
				ID:       "1788",
				LangCode: "hu",
				Navi:     "si",
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
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:       "2380",
				LangCode: "en",
				Navi:     "tseng",
				Affixes:  affix{},
			},
		},
	}, // simple (no *fixes)
	{
		name: "luyu",
		args: args{
			searchNaviText: "luyu",
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:       "1044",
				LangCode: "en",
				Navi:     "lu",
				Affixes: affix{
					Infix: []string{
						"uy",
					},
				},
			},
		},
	},
	{
		name: "seiyi",
		args: args{
			searchNaviText: "seiyi",
			languageCode:   "fr",
		},
		want: []Word{
			{
				ID:       "1788",
				LangCode: "fr",
				Navi:     "si",
				Affixes: affix{
					Infix: []string{
						"ei",
					},
				},
			},
		},
	}, // ayi override
	{
		name: "zenuyeke",
		args: args{
			searchNaviText: "zenuyeke",
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:       "3968",
				LangCode: "en",
				Navi:     "zenke",
				Affixes: affix{
					Infix: []string{
						"uy",
					},
				},
			},
		},
	}, // zenke `yu` override
	{
		name: "verìn",
		args: args{
			searchNaviText: "verìn",
			languageCode:   "et",
		},
		want: []Word{
			{
				ID:       "6884",
				LangCode: "et",
				Navi:     "vrrìn",
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
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:       "4984",
				LangCode: "en",
				Navi:     "tswa'",
				Affixes: affix{
					Prefix: []string{
						"ketsuk",
					},
				},
			},
			{
				ID:       "7352",
				LangCode: "en",
				Navi:     "ketsuktswa'",
				Affixes:  affix{},
			},
			{
				ID:       "7356",
				LangCode: "en",
				Navi:     "tsuktswa'",
				Affixes: affix{
					Prefix: []string{
						"ke",
					},
				},
			},
		},
	}, // ke/ketsuk prefix
	{
		name: "tìtusaron",
		args: args{
			searchNaviText: "tìtusaron",
			languageCode:   "et",
		},
		want: []Word{
			{
				ID:       "2004",
				LangCode: "et",
				Navi:     "taron",
				Affixes: affix{
					Prefix: []string{
						"tì",
					},
					Infix: []string{
						"us",
					},
				},
			},
			{
				ID:       "6156",
				LangCode: "et",
				Navi:     "tìtusaron",
				Affixes:  affix{},
			},
		},
	}, // tì prefix
	{
		name: "fayioang",
		args: args{
			searchNaviText: "fayioang",
			languageCode:   "fr",
		},
		want: []Word{
			{
				ID:       "612",
				LangCode: "fr",
				Navi:     "ioang",
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
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:       "4804",
				LangCode: "en",
				Navi:     "soaia",
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
			languageCode:   "de",
		},
		want: []Word{
			{
				ID:       "5268",
				LangCode: "de",
				Navi:     "tsaw",
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
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:       "1380",
				LangCode: "en",
				Navi:     "oe",
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
			languageCode:   "nl",
		},
		want: []Word{
			{
				ID:       "1348",
				LangCode: "nl",
				Navi:     "nga",
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
			languageCode:   "de",
		},
		want: []Word{
			{
				ID:       "3676",
				LangCode: "de",
				Navi:     "'em",
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
			languageCode:   "fr",
		},
		want: []Word{
			{
				ID:       "3948",
				LangCode: "fr",
				Navi:     "wem",
				Affixes: affix{
					Suffix: []string{
						"tswo",
					},
				},
			},
		},
	}, // tswo suffix
	{
		name: "pawnengsì",
		args: args{
			searchNaviText: "pawnengsì",
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:       "1512",
				LangCode: "en",
				Navi:     "peng",
				Affixes: affix{
					Infix: []string{
						"awn",
					},
					Suffix: []string{
						"sì",
					},
				},
			},
		},
	}, // awn infix and some suffix
	{
		name: "tsuknumesì",
		args: args{
			searchNaviText: "tsuknumesì",
			languageCode:   "nl",
		},
		want: []Word{
			{
				ID:       "1340",
				LangCode: "nl",
				Navi:     "nume",
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
			languageCode:   "et",
		},
		want: []Word{
			{
				ID:       "5268",
				LangCode: "et",
				Navi:     "tsaw",
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
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:       "2344",
				LangCode: "en",
				Navi:     "tsamsiyu",
			},
		},
	},
	{
		name: "'ueyä",
		args: args{
			searchNaviText: "'ueyä",
			languageCode:   "de",
		},
		want: []Word{
			{
				ID:       "108",
				LangCode: "de",
				Navi:     "'u",
				Affixes: affix{
					Suffix: []string{
						"eyä",
					},
				},
			},
			{
				ID:       "4180",
				LangCode: "de",
				Navi:     "'uo",
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
			languageCode:   "et",
		},
		want: []Word{
			{
				ID:       "192",
				LangCode: "et",
				Navi:     "awnga",
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
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:       "428",
				LangCode: "en",
				Navi:     "fpi+",
			},
		},
	}, // "+" support
	{
		name: "pe",
		args: args{
			searchNaviText: "pe",
			languageCode:   "en",
		},
		want: []Word{
			{
				ID:       "1488",
				LangCode: "en",
				Navi:     "--pe+",
			},
		},
	}, // "--" support
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
				ID:       "12989",
				LangCode: "de",
				Navi:     "fil",
			},
		},
	},
	{
		name: "Tier",
		args: args{
			searchNaviText: "Tier",
			languageCode:   "de",
		},
		want: []Word{
			{
				ID:       "612",
				LangCode: "de",
				Navi:     "ioang",
			},
			{
				ID:       "1440",
				LangCode: "de",
				Navi:     "pa'li",
			},
			{
				ID:       "2124",
				LangCode: "de",
				Navi:     "tireaioang",
			},
			{
				ID:       "2744",
				LangCode: "de",
				Navi:     "yerik",
			},
			{
				ID:       "7676",
				LangCode: "de",
				Navi:     "fwampop",
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
				ID:       "4524",
				LangCode: "de",
				Navi:     "venzek",
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
				ID:       "1608",
				LangCode: "en",
				Navi:     "pxi",
			},
			{
				ID:       "1616",
				LangCode: "en",
				Navi:     "pxiut",
			},
			{
				ID:       "8604",
				LangCode: "en",
				Navi:     "litx",
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
				ID:       "936",
				LangCode: "en",
				Navi:     "kxll",
			},
			{
				ID:       "5256",
				LangCode: "en",
				Navi:     "kxll si",
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
				ID:       "9868",
				LangCode: "en",
				Navi:     "srä",
			},
		},
	},
}

func TestTranslateFromNavi(t *testing.T) {
	for _, tt := range naviWords {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateFromNavi(tt.args.searchNaviText, tt.args.languageCode); !wordSimpleEqual(got, tt.want, true) {
				t.Errorf("TranslateFromNavi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslateFromNaviCached(t *testing.T) {
	CacheDict()
	TestTranslateFromNavi(t)
	UncacheDict()
}

func BenchmarkTranslateFromNavi(b *testing.B) {
	for _, bm := range naviWords {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				TranslateFromNavi(bm.args.searchNaviText, bm.args.languageCode)
			}
		})
	}
}

func BenchmarkTranslateFromNaviCached(b *testing.B) {
	CacheDict()
	BenchmarkTranslateFromNavi(b)
	UncacheDict()
}

func BenchmarkTranslateFromNaviBig(b *testing.B) {
	open, err := os.Open("misc/random_words.txt")
	if err != nil {
		return
	}
	defer open.Close()

	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		line := scanner.Text()

		b.Run(line, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				TranslateFromNavi(line, "en")
			}
		})
	}
}

func BenchmarkTranslateFromNaviBigCached(b *testing.B) {
	CacheDict()
	BenchmarkTranslateFromNaviBig(b)
	UncacheDict()
}

func TestTranslateToNavi(t *testing.T) {
	for _, tt := range englishWords {
		t.Run(tt.name, func(t *testing.T) {
			if gotResults := TranslateToNavi(tt.args.searchNaviText, tt.args.languageCode); !wordSimpleEqual(gotResults, tt.want, false) {
				t.Errorf("TranslateToNavi() = %v, want %v", gotResults, tt.want)
			}
		})
	}
}

func TestTranslateToNaviCached(t *testing.T) {
	CacheDict()
	TestTranslateToNavi(t)
	UncacheDict()
}

func BenchmarkTranslateToNaviBig(b *testing.B) {
	open, err := os.Open("misc/random_words_english.txt")
	if err != nil {
		return
	}
	defer open.Close()

	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		line := scanner.Text()

		b.Run(line, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				TranslateToNavi(line, "en")
			}
		})
	}
}

func BenchmarkTranslateToNaviBigCached(b *testing.B) {
	CacheDict()
	BenchmarkTranslateToNaviBig(b)
	UncacheDict()
}

func TestRandom(t *testing.T) {
	type args struct {
		amount   int
		langCode string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				amount:   1,
				langCode: "en",
			},
		},
		{
			name: "2",
			args: args{
				amount:   2,
				langCode: "de",
			},
		},
		{
			name: "3",
			args: args{
				amount:   3,
				langCode: "fr",
			},
		},
		{
			name: "4",
			args: args{
				amount:   4,
				langCode: "nl",
			},
		},
		{
			name: "5",
			args: args{
				amount:   5,
				langCode: "pl",
			},
		},
		{
			name: "500",
			args: args{
				amount:   500,
				langCode: "en",
			},
		},
		{
			name: "2000",
			args: args{
				amount:   2000,
				langCode: "en",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResults := Random(tt.args.amount, tt.args.langCode); !(len(gotResults) == tt.args.amount) {
				t.Errorf("Random() = %v, want %v", len(gotResults), tt.args.amount)
			}
		})
	}
}

func TestRandomCached(t *testing.T) {
	CacheDict()
	TestRandom(t)
	UncacheDict()
}
