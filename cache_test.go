package fwew_lib

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func Test_cacheDict(t *testing.T) {
	// cache dict and test only one entry

	log.Println("huhu")
	fmt.Println("buh!")

	word := Word{
		ID:             "32",
		LangCode:       "de",
		Navi:           "'awve",
		IPA:            "ˈʔaw.vɛ",
		InfixLocations: "NULL",
		PartOfSpeech:   "adj.",
		Definition:     "erste/r/s",
		Source:         "ASG",
		Stressed:       "1",
		Syllables:      "'aw-ve",
		InfixDots:      "NULL",
		Affixes:        nil,
		Attempt:        "",
	}

	err := CacheDict()
	if err != nil {
		t.Errorf("Error caching Dictionary!!")
	}
	if !reflect.DeepEqual(dictionary["de"][8], word) {
		t.Errorf("Read wrong word from cache: expected \"%s\", but got \"%s\"", word, dictionary["de"][8])
	}
}
