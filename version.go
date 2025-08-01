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

// Package fwew_lib contains all the things. version.go handles program version.
package fwew_lib

import "fmt"

type version struct {
	Major, Minor, Patch int
	Label               string
	Name                string
	DictBuild           string
}

// Version is a printable version struct containing program version information
var Version = version{
	5, 24, 3,
	"",
	"Kanua Kenten",
	"",
}

func init() {
	file := FindDictionaryFile()
	if file != "" {
		Version.DictBuild = SHA1Hash(file)
	}
}

func (v version) String() string {
	if v.Label != "" {
		return fmt.Sprintf("%s: %d.%d.%d-%s \"%s\"\ndictionary %s",
			Text("name"), v.Major, v.Minor, v.Patch, v.Label, v.Name, v.DictBuild)
	}

	return fmt.Sprintf("%s %d.%d.%d \"%s\"\ndictionary %s",
		Text("name"), v.Major, v.Minor, v.Patch, v.Name, v.DictBuild)
}
