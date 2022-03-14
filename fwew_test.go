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
	// assure dict, so tests wont fail
	AssureDict()

	// call flag.Parse() here if TestMain uses flags
	if testing.CoverMode() != "" {
		debugMode = true
	}
	os.Exit(m.Run())
}

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
			//w1.DE != w2.DE ||
			//w1.EN != w2.EN ||
			//w1.ET != w2.ET ||
			//w1.FR != w2.FR ||
			//w1.HU != w2.HU ||
			//w1.NL != w2.NL ||
			//w1.PL != w2.PL ||
			//w1.RU != w2.RU ||
			//w1.SV != w2.SV ||
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
		name: "tsatan",
		args: args{
			searchNaviText: "tsatan",
		},
		want: []Word{
			{
				ID: "160",
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
				ID: "7408",
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
				ID: "160",
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
				ID: "56",
				Navi: "'eveng",
				Affixes: affix{
					Prefix: []string{
						"pe",
						"pxe",
					},
				},
			},
		},
	},
	{
		name: "pepfil",
		args: args{
			searchNaviText: "pepfil",
		},
		want: []Word{
			{
				ID:   "6616",
				Navi: "pil",
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
				ID:   "12989",
				Navi: "fil",
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
		name: "verìn",
		args: args{
			searchNaviText: "verìn",
		},
		want: []Word{
			{
				ID:   "6884",
				Navi: "vrrìn",
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
				ID:   "4984",
				Navi: "tswa'",
				Affixes: affix{
					Prefix: []string{
						"ketsuk",
					},
				},
			},
			{
				ID:      "7352",
				Navi:    "ketsuktswa'",
				Affixes: affix{},
			},
			{
				ID:   "7356",
				Navi: "tsuktswa'",
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
			{
				ID:      "6156",
				Navi:    "tìtusaron",
				Affixes: affix{},
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
		name: "pawnengsì",
		args: args{
			searchNaviText: "pawnengsì",
		},
		want: []Word{
			{
				ID:   "1512",
				Navi: "peng",
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
						"eyä",
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
			{
				ID:   "7676",
				Navi: "fwampop",
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
				ID:   "1608",
				Navi: "pxi",
			},
			{
				ID:   "1616",
				Navi: "pxiut",
			},
			{
				ID:   "8604",
				Navi: "litx",
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
}

func TestTranslateFromNavi(t *testing.T) {
	for _, tt := range naviWords {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := TranslateFromNavi(tt.args.searchNaviText); err != nil || !wordSimpleEqual(got, tt.want, true) {
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
				TranslateFromNavi(bm.args.searchNaviText)
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
				TranslateFromNavi(line)
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
			if gotResults, _ := Random(tt.args.amount, nil); !(len(gotResults) == tt.args.amount) {
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
