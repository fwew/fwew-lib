package fwew_lib

import (
	"testing"
)

func Test_cacheDict(t *testing.T) {
	// cache dict and test only one entry

	word := Word{
		ID:   "4",
		Navi: "'ampi",
	}

	err := CacheDictHash()
	if err != nil {
		t.Fatalf("Error caching Dictionary!!")
	}

	entry := dictHashLoose["'ampi"]

	if word.ID != entry[0].ID {
		t.Errorf("Read wrong word from cache:\n"+
			"Id: \"%s\" == \"%s\"\n"+
			"Navi: \"%s\" == \"%s\"\n",
			word.ID, entry[0].ID,
			word.Navi, entry[0].Navi)
	}
}
