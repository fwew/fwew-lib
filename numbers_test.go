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

// Package fwew_lib contains all the things. numbers_test.go tests numbers.go functions.
package fwew_lib

import (
	"fmt"
	"testing"
)

// zero, negative, decimal, >77777, numbers where at least one or more digit is zero

var testCases = []struct {
	word   string
	number int
}{
	// 0...
	{"kew", 0},

	// simple numbers
	{"'aw", 1},
	{"mune", 2},
	{"pxey", 3},
	{"tsìng", 4},
	{"mrr", 5},
	{"pukap", 6},
	{"kinä", 7},

	// 2-digit numbers
	{"vol", 0o10},
	{"volaw", 0o11},
	{"vomun", 0o12},
	{"vopey", 0o13},
	{"vosìng", 0o14},
	{"vomrr", 0o15},
	{"vofu", 0o16},
	{"vohin", 0o17},

	{"mevol", 0o20},
	{"mevolaw", 0o21},
	{"mevomun", 0o22},
	{"mevopey", 0o23},
	{"mevosìng", 0o24},
	{"mevomrr", 0o25},
	{"mevofu", 0o26},
	{"mevohin", 0o27},

	{"pxevol", 0o30},
	{"pxevolaw", 0o31},
	{"pxevomun", 0o32},
	{"pxevopey", 0o33},
	{"pxevosìng", 0o34},
	{"pxevomrr", 0o35},
	{"pxevofu", 0o36},
	{"pxevohin", 0o37},

	{"tsìvol", 0o40},
	{"tsìvolaw", 0o41},
	{"tsìvomun", 0o42},
	{"tsìvopey", 0o43},
	{"tsìvosìng", 0o44},
	{"tsìvomrr", 0o45},
	{"tsìvofu", 0o46},
	{"tsìvohin", 0o47},

	{"mrrvol", 0o50},
	{"mrrvolaw", 0o51},
	{"mrrvomun", 0o52},
	{"mrrvopey", 0o53},
	{"mrrvosìng", 0o54},
	{"mrrvomrr", 0o55},
	{"mrrvofu", 0o56},
	{"mrrvohin", 0o57},

	{"puvol", 0o60},
	{"puvolaw", 0o61},
	{"puvomun", 0o62},
	{"puvopey", 0o63},
	{"puvosìng", 0o64},
	{"puvomrr", 0o65},
	{"puvofu", 0o66},
	{"puvohin", 0o67},

	{"kivol", 0o70},
	{"kivolaw", 0o71},
	{"kivomun", 0o72},
	{"kivopey", 0o73},
	{"kivosìng", 0o74},
	{"kivomrr", 0o75},
	{"kivofu", 0o76},
	{"kivohin", 0o77},

	// 3-digit full numbers
	{"zam", 0o100},
	{"mezam", 0o200},
	{"pxezam", 0o300},
	{"tsìzam", 0o400},
	{"mrrzam", 0o500},
	{"puzam", 0o600},
	{"kizam", 0o700},

	// 4-digit full numbers
	{"vozam", 0o1000},
	{"mevozam", 0o2000},
	{"pxevozam", 0o3000},
	{"tsìvozam", 0o4000},
	{"mrrvozam", 0o5000},
	{"puvozam", 0o6000},
	{"kivozam", 0o7000},

	// 5-digit full numbers
	{"zazam", 0o10000},
	{"mezazam", 0o20000},
	{"pxezazam", 0o30000},
	{"tsìzazam", 0o40000},
	{"mrrzazam", 0o50000},
	{"puzazam", 0o60000},
	{"kizazam", 0o70000},

	// some random numbers in all ranges
	// 3-digit numbers
	{"zam", 0o100},    // zam
	{"zamaw", 0o101},  // zam
	{"zamun", 0o102},  // za
	{"zapey", 0o103},  // za
	{"zasìng", 0o104}, // za
	{"zamrr", 0o105},  // za
	{"zafu", 0o106},   // za
	{"zahin", 0o107},  // za

	{"zavol", 0o110},    // za
	{"zavolaw", 0o111},  // za
	{"zavomun", 0o112},  // za
	{"zavopey", 0o113},  // za
	{"zavosìng", 0o114}, // za
	{"zavomrr", 0o115},  // za
	{"zavofu", 0o116},   // za
	{"zavohin", 0o117},  // za

	{"zamevol", 0o120},   // zam
	{"zampxevol", 0o130}, // zam
	{"zamtsìvol", 0o140}, // zam
	{"zamrrvol", 0o150},  // zam
	{"zampuvol", 0o160},  // zam
	{"zamkivol", 0o170},  // zam

	{"mezamtsìvohin", 0o247},
	{"mezamrrvol", 0o250}, // xx0
	{"pxezavosìng", 0o314},
	{"pxezamevol", 0o320},
	{"pxezamevolaw", 0o321},
	{"pxezamrrvomrr", 0o355},
	{"tsìzamun", 0o402}, // x0x
	{"mrrzamrrvomun", 0o552},
	{"puzamevolaw", 0o621},
	{"kizafu", 0o706},
	{"kizavomrr", 0o715},
	{"kizampxevohin", 0o737},
	{"kizamtsìvomrr", 0o745},
	{"kizampuvomun", 0o762},

	// 4-digit numbers
	{"vozampxezampxevohin", 0o1337},
	{"vozamkizamtsìvomrr", 0o1745},
	{"mevozamrrzam", 0o2500},        // xx00
	{"pxevozampey", 0o3003},         // x00x
	{"pxevozamvomun", 0o3012},       // x0xx
	{"pxevozamezafu", 0o3206},       // xx0x
	{"pxevozampuzamtsìvol", 0o3640}, // xxx0
	{"tsìvozamrrzamevomrr", 0o4525},
	{"tsìvozamrrzampuvosìng", 0o4564},
	{"tsìvozampuzampxevofu", 0o4636},
	{"mrrvozamzavofu", 0o5116},
	{"puvozamzampxevofu", 0o6136},
	{"puvozamzamkivofu", 0o6176},
	{"puvozamtsìzamrr", 0o6405},
	{"kivozampuvolaw", 0o7061},
	{"kivozamtsìzamkivomun", 0o7472},
	{"kivozamtsìzamkivofu", 0o7476},

	{"zazamevozam", 0o12000}, // xx000
	{"zazamevozamkivomun", 0o12072},
	{"zazamevozamkizamrrvofu", 0o12756},
	{"zazampxevozamezamkivofu", 0o13276},
	{"zazamkivozamzamkivol", 0o17170},
	{"pxezazamrrvofu", 0o30056},          // x00xx
	{"pxezazamkizamtsìvomrr", 0o30745},   // x0xxx
	{"pxezazamevozampxezahin", 0o32307},  // xxx0x
	{"pxezazamtsìvozamtsìvofu", 0o34046}, // xx0xx
	{"pxezazampuvozamrrzamrrvomrr", 0o36555},
	{"tsìzazamvozamzahin", 0o41107},
	{"mrrzazamvosìng", 0o50014},
	{"mrrzazampxezamtsìvohin", 0o50347},
	{"mrrzazamevozampuzam", 0o52600}, // xxx00
	{"mrrzazamevozamkizamkivol", 0o52770},
	{"mrrzazamkivozamrrzamrrvosìng", 0o57554},
	{"puzazamevozamezampxevol", 0o62230}, // xxxx0
	{"kizazamsìng", 0o70004},             // x000x
	{"kizazamvozamhin", 0o71007},         // xx00x
	{"kizazamkivozamkizamkivohin", 0o77777},
}

func Test_NaviToNumber(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Translate number %#o", testCase.number), func(t *testing.T) {
			if number, err := NaviToNumber(testCase.word); err != nil || number != testCase.number {
				t.Errorf("Translated number of \"%s\" was incorrect: expected \"%#o\", but got \"%#o\"", testCase.word, testCase.number, number)
			}
		})
	}
}

func Test_NumberToNavi(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Translate %#o", testCase.number), func(t *testing.T) {
			if word, err := NumberToNavi(testCase.number); err != nil || word != testCase.word {
				t.Errorf("Translated word of number \"%#o\" was incorrect: expected \"%s\", but got \"%s\"", testCase.number, testCase.word, word)
			}
		})
	}
}
