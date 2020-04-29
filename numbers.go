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

// Package fwew_lib contains all the things. numbers.go contains all the stuff for the number parsing
package fwew_lib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var naviVocab = [][]string{
	// 0 1 2 3 4 5 6 7 actual
	{"", "'aw", "mune", "pxey", "tsìng", "mrr", "pukap", "kinä"},
	// 0 1 2 3 4 5 6 7 last digit
	{"", "aw", "mun", "pey", "sìng", "mrr", "fu", "hin"},
	// 0 1 2 3 4 5 6 7 first or middle digit
	{"", "", "me", "pxe", "tsì", "mrr", "pu", "ki"},
	// 0 1 2 3 4 powers of 8
	{"", "vo", "zam", "vozam", "zazam"},
	// 0 1 2 3 4 powers of 8 last digit
	{"", "l", "", "", ""},
}

// "word number portion": octal value
// the upper array is the digit.
var numTable = []map[string]int{
	{
		"kizazam":  0o70000,
		"kizaza":   0o70000,
		"puzazam":  0o60000,
		"puzaza":   0o60000,
		"mrrzazam": 0o50000,
		"mrrzaza":  0o50000,
		"rrzazam":  0o50000,
		"rrzaza":   0o50000,
		"tsìzazam": 0o40000,
		"tsìzaza":  0o40000,
		"pxezazam": 0o30000,
		"pxezaza":  0o30000,
		"mezazam":  0o20000,
		"mezaza":   0o20000,
		"ezazam":   0o20000,
		"ezaza":    0o20000,
		"zazam":    0o10000,
		"zaza":     0o10000,
	},
	{
		"kivozam":  0o7000,
		"kivoza":   0o7000,
		"puvozam":  0o6000,
		"puvoza":   0o6000,
		"mrrvozam": 0o5000,
		"mrrvoza":  0o5000,
		"rrvozam":  0o5000,
		"rrvoza":   0o5000,
		"tsìvozam": 0o4000,
		"tsìvoza":  0o4000,
		"pxevozam": 0o3000,
		"pxevoza":  0o3000,
		"mevozam":  0o2000,
		"mevoza":   0o2000,
		"evozam":   0o2000,
		"evoza":    0o2000,
		"vozam":    0o1000,
		"voza":     0o1000,
	},
	{
		"kizam":  0o700,
		"kiza":   0o700,
		"puzam":  0o600,
		"puza":   0o600,
		"mrrzam": 0o500,
		"mrrza":  0o500,
		"rrzam":  0o500,
		"rrza":   0o500,
		"tsìzam": 0o400,
		"tsìza":  0o400,
		"pxezam": 0o300,
		"pxeza":  0o300,
		"mezam":  0o200,
		"meza":   0o200,
		"ezam":   0o200,
		"eza":    0o200,
		"zam":    0o100,
		"za":     0o100,
	},
	{
		"kivol":  0o70,
		"kivo":   0o70,
		"puvol":  0o60,
		"puvo":   0o60,
		"mrrvol": 0o50,
		"mrrvo":  0o50,
		"rrvol":  0o50,
		"rrvo":   0o50,
		"tsìvol": 0o40,
		"tsìvo":  0o40,
		"pxevol": 0o30,
		"pxevo":  0o30,
		"mevol":  0o20,
		"mevo":   0o20,
		"evol":   0o20,
		"evo":    0o20,
		"vol":    0o10,
		"vo":     0o10,
	},
	{
		"hin":  0o7,
		"fu":   0o6,
		"mrr":  0o5,
		"rr":   0o5,
		"sìng": 0o4,
		"pey":  0o3,
		"mun":  0o2,
		"un":   0o2,
		"aw":   0o1,
	},
}

// The regex values for the different values.
// The upper array is the digit.
var numTableRegexp = [][]string{
	{
		"kizazam?",
		"puzazam?",
		"m?rrzazam?",
		"tsìzazam?",
		"pxezazam?",
		"m?ezazam?",
		"zazam?",
	},
	{
		"kivozam?",
		"puvozam?",
		"m?rrvozam?",
		"tsìvozam?",
		"pxevozam?",
		"m?evozam?",
		"vozam?",
	},
	{
		"kizam?",
		"puzam?",
		"m?rrzam?",
		"tsìzam?",
		"pxezam?",
		"m?ezam?",
		"zam?",
	},
	{
		"kivol?",
		"puvol?",
		"m?rrvol?",
		"tsìvol?",
		"pxevol?",
		"m?evol?",
		"vol?",
	},
	{
		"hin",
		"fu",
		"mrr",
		"rr",
		"sìng",
		"pey",
		"mun",
		"un",
		"aw",
	},
}

// Translate a Na'vi number word to the actual integer.
// Na'vi numbers are octal values, so the integer is defined as octal number, and can easily be displayed as decimal number.
func NaviToNumber(input string) int {
	input = strings.ToLower(input)
	// kew
	if input == "kew" {
		return 0
	}

	// 'aw mune pxey tsìng mrr pukap kinä
	// literal numbers 1-7
	for i, w := range naviVocab[0] {
		if input == w && w != "" {
			return i
		}
	}

	// build regexp for all other numbers
	// regex for big values
	var regexpString string
	for _, digit := range numTableRegexp {
		regexpString += "("
		first := true

		for _, number := range digit {
			if !first {
				regexpString += "|"
			}
			regexpString += number
			first = false
		}
		regexpString += ")?"
	}

	re := regexp.MustCompile(regexpString)
	tmp := re.FindStringSubmatch(input)
	var n int
	if len(tmp) > 0 && len(tmp[0]) > 0 {
		for i, v := range tmp[1:] {
			n += numTable[i][v]
		}
	}
	return n
}

//------------------------------------------------------------------------

func wordify(input string) string {
	rev := Reverse(input)
	output := ""
	if len(input) == 1 {
		if input == "0" {
			return "kew"
		}
		inty, _ := strconv.Atoi(input)
		return naviVocab[0][inty]
	}
	for i, d := range rev {
		switch i {
		case 0: // 7777[7]
			output = naviVocab[1][int(d-'0')] + output
			if int(d-'0') == 1 && rev[1] != '0' {
				output = naviVocab[4][1] + output
			}
		case 1: // 777[7]7
			if int(d-'0') > 0 {
				output = naviVocab[2][int(d-'0')] + naviVocab[3][1] + output
			}
		case 2: // 77[7]77
			if int(d-'0') > 0 {
				output = naviVocab[2][int(d-'0')] + naviVocab[3][2] + output
			}
		case 3: // 7[7]777
			if int(d-'0') > 0 {
				output = naviVocab[2][int(d-'0')] + naviVocab[3][3] + output
			}
		case 4: // [7]7777
			if int(d-'0') > 0 {
				output = naviVocab[2][int(d-'0')] + naviVocab[3][4] + output
			}
		}
	}
	for _, d := range []string{"01", "02", "03", "04", "05", "06", "07"} {
		if rev[0:2] == d {
			output = output + naviVocab[4][1]
		}
	}
	output = strings.Replace(output, "mm", "m", -1)
	return output
}

// Convert is the main number conversion function
func Convert(input string, reverse bool) string {
	output := ""
	if reverse {
		i, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return fmt.Sprintf("%s: %s\n", Text("invalidDecimalError"), input)
		}
		if !Valid(i, reverse) {
			return fmt.Sprintf("%s\n", Text("invalidIntError"))
		}
		o := strconv.FormatInt(int64(i), 8)
		output += fmt.Sprintf("Octal: %s\n", o)
		output += fmt.Sprintf("Na'vi: %s\n", wordify(o))
	} else {
		var io int64
		var err error
		if IsLetter(input) {
			//io, err = strconv.ParseInt(unwordify(input), 8, 64)
		} else {
			io, err = strconv.ParseInt(input, 8, 64)
		}
		if err != nil {
			return fmt.Sprintf("%s: %s\n", Text("invalidOctalError"), input)
		}
		if !Valid(io, reverse) {
			return fmt.Sprintf("%s\n", Text("invalidIntError"))
		}
		d := strconv.FormatInt(int64(io), 10)
		o := strconv.FormatInt(int64(io), 8)
		output += fmt.Sprintf("Decimal: %s\n", d)
		if IsLetter(input) {
			output += fmt.Sprintf("Octal: %s\n", o)
		} else {
			output += fmt.Sprintf("Na'vi: %s\n", wordify(input))
		}
	}
	return output
}
