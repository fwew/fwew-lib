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

import "testing"

// zero, negative, decimal, >77777, numbers where at least one or more digit is zero

func Test_NaviToNumber(t *testing.T) {
	type testCaseStruct struct {
		word     string
		expected int
	}
	testCases := []testCaseStruct{
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
		{"pxezamrrvomrr", 0o355},
		{"kizavomrr", 0o715},
		{"kizapuvomun", 0o762},
		{"mezamtsìvohin", 0o247},
		{"pxezamevol", 0o320},
		{"kizampxevohin", 0o737},
		{"tsìzamun", 0o402},
		{"kizatsìvomrr", 0o745},
		{"zamsìng", 0o104},
		{"pxezamvosìng", 0o314},
		{"mrrzamrrvomun", 0o552},
		{"pxezamevolaw", 0o321},
		{"puzamevolaw", 0o621},
		{"mezamrrvol", 0o250},
		{"zafu", 0o106},

		// 4-digit numbers
		{"kivozamtsìzamkivofu", 0o7476},
		{"puvozamzamkivofu", 0o6176},
		{"pxevozampey", 0o3003},
		{"kivozamtsìzamkivomun", 0o7472},
		{"kivozampuvolaw", 0o7061},
		{"puvozamzampxevofu", 0o6136},
		{"pxevozammezamfu", 0o3206},
		{"vozamkizatsìvomrr", 0o1745},
		{"pxevozamvomun", 0o3012},
		{"pxevozampuzamtsìvol", 0o3640},
		{"tsìvozammrrzampuvosìng", 0o4564},
		{"tsìvozampuzampxevofu", 0o4636},
		{"puvozamtsìzamrr", 0o6405},
		{"mrrvozamzamvofu", 0o5116},
		{"tsìvozamrrzammevomrr", 0o4525},
		{"vozampxezampxevohin", 0o1337},

		{"pxezazamkizamtsìvomrr", 0o30745},
		{"pxezazamtsìvozamtsìvofu", 0o34046},
		{"pxezazammevozampxezamhin", 0o32307},
		{"puzazammevozammezampxevol", 0o62230},
		{"pxezazammrrvofu", 0o30056},
		{"kizazamvozamhin", 0o71007},
		{"mrrzazammevozampuzam", 0o52600},
		{"kizazamsìng", 0o70004},
		{"zazammevozam", 0o12000},
		{"tsìzazamvozamzamhin", 0o41107},
		{"zazampxevozammezamkivofu", 0o13276},
		{"zazamkivozamzamkivo", 0o17170},
		{"pxezazampuvozammrrzammrrvomrr", 0o36555},
		{"mrrzazampxezamtsìvohin", 0o50347},
		{"zazammevozamkizammrrvofu", 0o12756},
		{"mrrzazamkivozammrrzammrrvosìng", 0o57554},
		{"mrrzazamvosìng", 0o50014},
		{"zazammevozamkivomun", 0o12072},
		{"mrrzazammevozamkizamkivo", 0o52770},
		{"kizazamkivozamkizamkivohin", 0o77777},
	}

	for _, testCase := range testCases {
		if number := NaviToNumber(testCase.word); number != testCase.expected {
			t.Errorf("Translated number of \"%s\" was incorrect: expected %#o, but got %#o", testCase.word, testCase.expected, number)
		}
	}
}
