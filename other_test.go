package fwew_lib

import "testing"

func TestGetLenitionTable(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var lt = GetLenitionTable()
	if lt != lenitionTable {
		t.Errorf("Lenition table not loaded or doesn't match itself")
	}
}

func TestGetShortLenitionTable(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var st = GetShortLenitionTable()
	if st != shortLenitionTable {
		t.Errorf("Short lenition table not loaded or doesn't match itself")
	}
}

func TestGetThatTable(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var tt = GetThatTable()
	if tt != thatTable {
		t.Errorf("That table not loaded or doesn't match itself")
	}
}

func TestGetOtherThats(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var ott = GetOtherThats()
	if ott != otherThats {
		t.Errorf("Other thats table not loaded or doesn't match itself")
	}
}

func TestGetMultiwordWords(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var mw = GetMultiwordWords()
	if mw == nil {
		t.Errorf("Multiword words not loaded")
	}
	if len(mw) == 0 {
		t.Errorf("Multiword words is empty")
	}
}

func TestGetHomonyms(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var homs, err = GetHomonyms()
	if err != nil {
		t.Error(err)
	}
	if homs == nil {
		t.Errorf("Homonyms not loaded")
	}
	if len(homs) == 0 {
		t.Errorf("Homonyms is empty")
	}
}

func TestGetOddballs(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var o, err = GetOddballs()
	if err != nil {
		t.Error(err)
	}
	if o == nil {
		t.Errorf("Oddballs not loaded")
	}
	if len(o) == 0 {
		t.Errorf("Oddballs is empty")
	}
}

func TestGetMultiIPA(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var m, err = GetMultiIPA()
	if err != nil {
		t.Error(err)
	}
	if m == nil {
		t.Errorf("MultiIPA not loaded")
	}
	if len(m) == 0 {
		t.Errorf("MultiIPA is empty")
	}
}

func TestGetPhonemeDistrosMap(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var lang = "en"
	var p = GetPhonemeDistrosMap(lang)
	if p == nil {
		t.Errorf("Phoneme distros map not loaded")
	}
	if len(p) == 0 {
		t.Errorf("Phoneme distros map is empty")
	}
}
