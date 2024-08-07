package test

import (
	"github.com/fwew/fwew-lib/v5"
	"testing"
)

func Test_cacheDict(t *testing.T) {
	var err error

	err = fwew_lib.CacheDictHash()
	err = fwew_lib.CacheDictHash2()

	if err != nil {
		t.Fatalf("Error caching Dictionary!")
	}

	dict, err := fwew_lib.GetFullDict()

	if err != nil || len(dict) == 0 {
		t.Fatalf("Error reading Dictionary!")
	}

	fwew_lib.UncacheHashDict()
	fwew_lib.UncacheHashDict2()
}
