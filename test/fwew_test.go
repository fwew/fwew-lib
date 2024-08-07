package test

import (
	"bufio"
	"github.com/fwew/fwew-lib/v5"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	// assure dict, so tests wont fail
	err := fwew_lib.AssureDict()

	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func wordSimpleEqual(w1a, w2a []fwew_lib.Word) bool {
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
	want []fwew_lib.Word
}{
	{
		name: "molte",
		args: args{
			searchNaviText: "molte",
		},
		want: []fwew_lib.Word{
			{
				ID:   "1120",
				Navi: "mllte",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "6224",
				Navi: "'a'aw",
				Affixes: fwew_lib.Affix{
					Lenition: []string{
						"'a→a",
					},
				},
			},
			{
				ID:   "12",
				Navi: "'aw",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "4316",
				Navi: "pukap",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "160",
				Navi: "atan",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "7408",
				Navi: "ìlva",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "160",
				Navi: "atan",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "56",
				Navi: "'eveng",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "12989",
				Navi: "fil",
				Affixes: fwew_lib.Affix{
					Prefix: []string{
						"pepe",
					},
				},
			},
			{
				ID:   "6616",
				Navi: "pil",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{},
	}, // empty
	{
		name: "säpeykiyevatsi",
		args: args{
			searchNaviText: "säpeykiyevatsi",
		},
		want: []fwew_lib.Word{
			{
				ID:   "1788",
				Navi: "si",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:      "2380",
				Navi:    "tseng",
				Affixes: fwew_lib.Affix{},
			},
		},
	}, // simple (no *fixes)
	{
		name: "luyu",
		args: args{
			searchNaviText: "luyu",
		},
		want: []fwew_lib.Word{
			{
				ID:   "1044",
				Navi: "lu",
				Affixes: fwew_lib.Affix{
					Infix: []string{
						"uy",
					},
				},
			},
			{
				ID:   "1044",
				Navi: "lu",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "1788",
				Navi: "si",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "3968",
				Navi: "zenke",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "6884",
				Navi: "vrrìn",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:      "7352",
				Navi:    "ketsuktswa'",
				Affixes: fwew_lib.Affix{},
			},
			{
				ID:   "4984",
				Navi: "tswa'",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:      "6156",
				Navi:    "tìtusaron",
				Affixes: fwew_lib.Affix{},
			},
			{
				ID:   "2004",
				Navi: "taron",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "612",
				Navi: "ioang",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "4804",
				Navi: "soaia",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "5268",
				Navi: "tsaw",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "1380",
				Navi: "oe",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "1348",
				Navi: "nga",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "3676",
				Navi: "'em",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "3948",
				Navi: "wem",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "1340",
				Navi: "nume",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "5268",
				Navi: "tsaw",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
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
		want: []fwew_lib.Word{
			{
				ID:   "108",
				Navi: "'u",
				Affixes: fwew_lib.Affix{
					Suffix: []string{
						"yä", "o",
					},
				},
			},
			{
				ID:   "4180",
				Navi: "'uo",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
			{
				ID:   "192",
				Navi: "awnga",
				Affixes: fwew_lib.Affix{
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
		want: []fwew_lib.Word{
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
	want []fwew_lib.Word
}{
	{
		name: "Spielzeug",
		args: args{
			searchNaviText: "Spielzeug",
			languageCode:   "de",
		},
		want: []fwew_lib.Word{
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
		want: []fwew_lib.Word{
			{
				ID:   "440",
				Navi: "fpxafaw",
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
			{
				ID:   "10704",
				Navi: "seyto",
			},
		},
	},
	{
		name: "Zehe",
		args: args{
			searchNaviText: "Zehe",
			languageCode:   "de",
		},
		want: []fwew_lib.Word{
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
		want: []fwew_lib.Word{
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
		want: []fwew_lib.Word{
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
		want: []fwew_lib.Word{
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
		want: []fwew_lib.Word{
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

func TestTranslateFromNaviCached(t *testing.T) {
	var err error

	err = fwew_lib.CacheDictHash()
	err = fwew_lib.CacheDictHash2()

	if err != nil {
		t.Fatalf("Could not Cache Dict")
	}

	for _, tt := range naviWords {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fwew_lib.TranslateFromNaviHash(tt.args.searchNaviText, true)
			if err != nil {
				t.Fatalf("TranslateFromNaviHash() returned error: %v", err)
			}

			if len(got) > 0 && len(got[0]) > 0 && !wordSimpleEqual(got[0][1:], tt.want) {
				t.Fatalf("TranslateFromNavi() = %v, want %v", got[0][1:], tt.want)
			}
		})
	}

	fwew_lib.UncacheHashDict()
	fwew_lib.UncacheHashDict2()
}

func BenchmarkTranslateFromNavi(b *testing.B) {
	for _, bm := range naviWords {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = fwew_lib.TranslateFromNaviHash(bm.args.searchNaviText, true)
			}
		})
	}
}

func BenchmarkTranslateFromNaviCached(b *testing.B) {
	err := fwew_lib.CacheDictHash()
	if err != nil {
		b.Fatalf("Could not Cache Dict")
	}
	BenchmarkTranslateFromNavi(b)
	fwew_lib.UncacheHashDict()
}

func BenchmarkTranslateFromNaviBig(b *testing.B) {
	open, err := os.Open("misc/random_words.txt")
	if err != nil {
		return
	}
	defer func(open *os.File) {
		err := open.Close()
		if err != nil {

		}
	}(open)

	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		line := scanner.Text()

		b.Run(line, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = fwew_lib.TranslateFromNaviHash(line, true)
			}
		})
	}
}

func BenchmarkTranslateFromNaviBigCached(b *testing.B) {
	err := fwew_lib.CacheDictHash()

	if err != nil {
		b.Fatalf("Could not Cache Dict")
	}

	BenchmarkTranslateFromNaviBig(b)

	fwew_lib.UncacheHashDict()
}

func TestTranslateToNaviCached(t *testing.T) {
	err := fwew_lib.CacheDictHash2()

	if err != nil {
		t.Fatalf("Could not Cache Dict")
	}

	for _, tt := range englishWords {
		t.Run(tt.name, func(t *testing.T) {
			gotResults := fwew_lib.TranslateToNaviHash(tt.args.searchNaviText, tt.args.languageCode)
			if !wordSimpleEqual(gotResults[0][1:], tt.want) {
				t.Errorf("TranslateToNavi() = %v, want %v", gotResults, tt.want)
			}
		})
	}

	fwew_lib.UncacheHashDict2()
}

func BenchmarkTranslateToNaviBig(b *testing.B) {
	open, err := os.Open("misc/random_words_english.txt")

	if err != nil {
		return
	}

	defer func(open *os.File) {
		err := open.Close()
		if err != nil {

		}
	}(open)

	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		line := scanner.Text()

		b.Run(line, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fwew_lib.TranslateToNaviHash(line, "en")
			}
		})
	}
}

func BenchmarkTranslateToNaviBigCached(b *testing.B) {
	err := fwew_lib.CacheDictHash2()

	if err != nil {
		b.Fatalf("Could not Cache Dict")
	}

	BenchmarkTranslateToNaviBig(b)

	fwew_lib.UncacheHashDict2()
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
			if gotResults, _ := fwew_lib.Random(tt.args.amount, nil, 1); !(len(gotResults) == tt.args.amount) {
				t.Errorf("Random() = %v, want %v", len(gotResults), tt.args.amount)
			}
		})
	}
}

func TestRandomCached(t *testing.T) {
	var err error
	err = fwew_lib.CacheDictHash()
	err = fwew_lib.CacheDictHash2()

	if err != nil {
		t.Fatalf("Unable to Cache Dict")
	}

	TestRandom(t)

	fwew_lib.UncacheHashDict()
	fwew_lib.UncacheHashDict2()
}
