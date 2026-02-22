package fwew_lib

func GetLenitionTable() [][2]string {
	return lenitionTable[:]
}

func GetShortLenitionTable() [][2]string {
	return shortLenitionTable[:]
}

func GetThatTable() [][5]string {
	return thatTable[:]
}

func GetOtherThats() [][3]string {
	return otherThats[:]
}

// GetMultiwordWords returns all words with spaces
func GetMultiwordWords() map[string][][]string {
	universalLock.Lock()
	defer universalLock.Unlock()
	return multiwordWords
}

// GetHomonyms returns all words with multiple definitions
func GetHomonyms() (results [][]Word, err error) {
	return TranslateFromNaviHash(homonyms, false, false, false)
}

// GetOddballs returns all words with non-standard phonotactics
func GetOddballs() (results [][]Word, err error) {
	return TranslateFromNaviHash(oddballs, true, false, false)
}

// GetMultiIPA returns all words with multiple definitions
func GetMultiIPA() (results [][]Word, err error) {
	return TranslateFromNaviHash(multiIPA, false, false, false)
}
