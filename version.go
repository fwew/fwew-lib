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
	5, 27, 2,
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
