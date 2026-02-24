package fwew_lib

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var dictionary []Word
var dictHashLoose map[string][]Word
var dictHashStrict map[string][]Word
var dictHashStrictReef map[string][]Word
var dictionaryCached bool
var dictHashCached bool
var dictHash2 metaDict
var dictHash2Parenthesis metaDict
var dictHash2Cached bool
var homonyms string
var oddballs string
var multiIPA string

type metaDict struct {
	EN map[string][]string
	DE map[string][]string
	ES map[string][]string
	ET map[string][]string
	FR map[string][]string
	HU map[string][]string
	IT map[string][]string
	KO map[string][]string
	NL map[string][]string
	PL map[string][]string
	PT map[string][]string
	RU map[string][]string
	SV map[string][]string
	TR map[string][]string
	UK map[string][]string
}

var dictHash2MapPtrs = []*map[string][]string{
	&dictHash2.EN, &dictHash2.DE, &dictHash2.ES, &dictHash2.ET, &dictHash2.FR, &dictHash2.HU, &dictHash2.IT,
	&dictHash2.KO, &dictHash2.NL, &dictHash2.PL, &dictHash2.PT, &dictHash2.RU, &dictHash2.SV, &dictHash2.TR,
	&dictHash2.UK}

var dictHash2ParenthesisMapPtrs = []*map[string][]string{
	&dictHash2Parenthesis.EN, &dictHash2Parenthesis.DE, &dictHash2Parenthesis.ES, &dictHash2Parenthesis.ET,
	&dictHash2Parenthesis.FR, &dictHash2Parenthesis.HU, &dictHash2Parenthesis.IT, &dictHash2Parenthesis.KO,
	&dictHash2Parenthesis.NL, &dictHash2Parenthesis.PL, &dictHash2Parenthesis.PT, &dictHash2Parenthesis.RU,
	&dictHash2Parenthesis.SV, &dictHash2Parenthesis.TR, &dictHash2Parenthesis.UK}

var letterMap = map[rune]int{
	' ': -1, '\'': 0, 'a': 1, '2': 2, '3': 3,
	'ä': 4, 'e': 5, '4': 6, '5': 7,
	'f': 8, 'h': 9, 'i': 10, 'ì': 11,
	'j': 12, 'k': 13, 'q': 14, 'l': 15,
	'1': 16, 'm': 17, 'n': 18, 'g': 19,
	'o': 20, 'p': 21, 'b': 22, 'r': 23,
	'0': 24, 's': 25, 't': 26, 'c': 27,
	'd': 28, 'u': 29, 'v': 30, 'w': 31,
	'y': 32, 'z': 33, '-': 34,
}

var nkx []string
var nkxSub = map[string]string{}

// A mutex to ensure concurrent requests to the
// dictionary and phoneme counts will not cause
// the program to crash
var universalLock sync.Mutex
var phonoLock sync.Mutex

// helper for nkx for shortest words first
// TODO: revisit implementation after removing always-false boolean var
func shortestFirst(array []string, input string) []string {
	var newArray []string
	for _, a := range array {
		if len(a) > len(input) {
			newArray = append(newArray, input)
		}
		newArray = append(newArray, a)
	}
	newArray = append(newArray, input)

	return newArray
}

// check if a file exists
func fileExists(filepath string) bool {
	fileStat, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}

	return !fileStat.IsDir()
}

// findDictionaryFile returns the location of the dictionary file
// the dictionary file can be places into:
// - <workingDir>/dictionary.txt
// - <workingDir>/.fwew/dictionary.txt
// - <homeDir>/.fwew/dictionary.txt
func findDictionaryFile() string {
	wd, err := os.Getwd()
	if err == nil {
		path := filepath.Join(wd, ".fwew", dictFileName)
		if fileExists(path) {
			return path
		}

		path = filepath.Join(wd, dictFileName)
		if fileExists(path) {
			return path
		}
	}

	path := filepath.Join(Text("dataDir"), dictFileName)
	if fileExists(path) {
		return path
	}

	return ""
}

func AlphabetizeHelper(a string, b string) bool {
	aCompacted := []rune(strings.ReplaceAll(compress(strings.ToLower(a)), "-", ""))

	// Start in the middle
	bCompacted := []rune(strings.ReplaceAll(compress(strings.ToLower(b)), "-", ""))
	lowestLen := min(len(aCompacted), len(bCompacted))
	// compare an individual word
	for j := 0; j < lowestLen; j++ {
		// If the new letter is bigger, wait until it gets
		if letterMap[aCompacted[j]] < letterMap[bCompacted[j]] {
			return true
		} else if letterMap[aCompacted[j]] > letterMap[bCompacted[j]] {
			return false
		}
		// If equal, continue
	}

	// longer words go after
	return len(aCompacted) < len(bCompacted)
}

func appendAndAlphabetize(words []Word, word Word) []Word {
	// Ensure it's not a duplicate
	for _, a := range words {
		if equal(word, a) {
			return words
		}
	}
	// new array
	switch len(words) {
	case 0:
		return []Word{word}
	case 1:
		var newWords []Word
		if AlphabetizeHelper(words[0].Syllables, word.Syllables) {
			newWords = []Word{words[0], word}
		} else {
			newWords = []Word{word, words[0]}
		}
		return newWords
	case 2:
		var newWords []Word
		if AlphabetizeHelper(word.Syllables, words[0].Syllables) {
			newWords = []Word{word, words[0], words[1]}
		} else if AlphabetizeHelper(words[1].Syllables, word.Syllables) {
			newWords = []Word{words[0], words[1], word}
		} else {
			newWords = []Word{words[0], word, words[1]}
		}
		return newWords
	}

	// start in the middle
	halfway := len(words) / 2

	// Copy the first half
	newWords := make([]Word, len(words[:halfway]))
	copy(newWords, words[:halfway])

	// Copy the second half
	oldWords := make([]Word, len(words[halfway:]))
	copy(oldWords, words[halfway:])

	// compare an individual word
	if AlphabetizeHelper(word.Syllables, words[halfway].Syllables) {
		// Copy the first half
		newWords = appendAndAlphabetize(newWords, word)

		// Join them
		newWords = append(newWords, oldWords...)
	} else {
		// Copy the second half
		oldWords = appendAndAlphabetize(oldWords, word)

		// Join them
		newWords = append(newWords, oldWords...)
	}

	return newWords
}

// hasNullDef is a helper to find empty definitions
func hasNullDef(definition string) bool {
	return strings.ToUpper(definition) == "NULL" || len(strings.Trim(definition, " ")) < 1
}

// englishIfNull defaults a word to English if a definition is not available in a certain language
func englishIfNull(word Word) Word {
	// English
	if hasNullDef(word.EN) {
		word.EN = "(no definition)"
	}

	definitions := []*string{
		&word.DE, &word.ES, &word.ET, &word.FR, &word.HU, &word.IT, &word.KO,
		&word.NL, &word.PL, &word.PT, &word.RU, &word.SV, &word.TR, &word.UK,
	}

	for _, definition := range definitions {
		if hasNullDef(*definition) {
			*definition = word.EN
		}
	}

	return word
}

// romanizeSecondIPA is a helper function to get phonetic transcriptions of secondary pronunciations
// Only multiple IPA words will call this function
func romanizeSecondIPA(IPA string) string {
	// now Romanize the IPA
	IPA = strings.ReplaceAll(IPA, "ʊ", "u")
	IPA = strings.ReplaceAll(IPA, "õ", "o") // vonvä' as võvä' only
	word := strings.Split(IPA, " ")

	breakdown := ""

	// get the last one only
	for j := 2; j < len(word); j++ {
		word[j] = strings.ReplaceAll(word[j], "[", "")
		// "or" means there's more than one IPA in this word, and we only want one
		if word[j] == "or" {
			breakdown = ""
			continue
		}

		syllables := strings.Split(word[j], ".")

		/* Onset */
		for k := range syllables {
			syllable := strings.ReplaceAll(syllables[k], "·", "")
			syllable = strings.ReplaceAll(syllable, "ˈ", "")
			syllable = strings.ReplaceAll(syllable, "ˌ", "")

			breakdown += syllableToRoman(syllable)
		}
		breakdown += " "
	}
	return strings.TrimSuffix(breakdown, " ")
}

func uncacheDict() {
	dictionaryCached = false
	dictionary = []Word{}
}

func cacheDict() error {
	var err error

	uncacheDict()
	err = runOnDB(func(word Word) error {
		dictionary = append(dictionary, word)
		return nil
	})

	if err == nil {
		fmt.Println("cache 0 loaded")
	} else {
		uncacheDict()
		err = runOnFile(func(word Word) error {
			dictionary = append(dictionary, word)
			return nil
		})
		//fmt.Println("cache 0 loaded (File)")
	}

	if err != nil {
		uncacheDict()
		return err
	}

	dictionaryCached = true

	return nil
}

func cacheDictHash() error {
	err := cacheDictHashOrig(true)
	if err == nil {
		fmt.Println("cache 1 loaded")
		return nil
	}
	err = cacheDictHashOrig(false)
	//fmt.Println("cache 1 loaded (File)")
	return err
}

// cacheDictHashOrig will cache the whole dictionary (Na'vi to natural language).
// Please call this if you want to translate multiple words or running infinitely (e.g., CLI-go-prompt, discord-bot)
func cacheDictHashOrig(mysql bool) error {
	// dont run if already is cached
	if len(dictHashLoose) != 0 {
		return nil
	}
	dictHashLoose = make(map[string][]Word)
	dictHashStrict = make(map[string][]Word)
	dictHashStrictReef = make(map[string][]Word)

	var tempHoms []string

	//Clear to avoid duplicates
	multiIPA = ""

	var f = func(word Word) error {
		standardizedWord := word.Navi
		badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\"„“”«»`

		// remove all the sketchy chars from arguments
		for _, c := range badChars {
			standardizedWord = strings.ReplaceAll(standardizedWord, string(c), "")
		}

		// normalize tìftang character
		standardizedWord = strings.ReplaceAll(standardizedWord, "’", "'")
		standardizedWord = strings.ReplaceAll(standardizedWord, "‘", "'")

		// find everything lowercase
		standardizedWord = strings.ToLower(standardizedWord)

		// Make sure we know of every word with nkx
		if strings.Contains(standardizedWord, "nkx") {
			fakeNG := strings.ReplaceAll(standardizedWord, "nkx", "ng")
			nkx = shortestFirst(nkx, fakeNG)
			nkxSub[fakeNG] = standardizedWord
		}

		standardizedWordArray := dialectCrunch(strings.Split(standardizedWord, " "), true, true)
		standardizedWordLoose := ""
		for i, a := range standardizedWordArray {
			if i != 0 {
				standardizedWordLoose += " "
			}
			standardizedWordLoose += a
		}

		strictReefArray := dialectCrunch(strings.Split(standardizedWord, " "), true, true)
		strictReef := ""
		for i, a := range strictReefArray {
			if i != 0 {
				strictReef += " "
			}
			strictReef += a
		}

		// If the word appears more than once, record it
		if _, ok := dictHashStrict[standardizedWord]; ok {
			found := slices.Contains(tempHoms, standardizedWord)
			if !found {
				tempHoms = append(tempHoms, standardizedWord)
			}
		}

		if strings.Contains(standardizedWord, "é") {
			noAcute := strings.ReplaceAll(standardizedWord, "é", "e")
			found := slices.Contains(tempHoms, noAcute)
			if !found {
				tempHoms = append(tempHoms, noAcute)
				tempHoms = append(tempHoms, standardizedWord)
			}
		}

		word = englishIfNull(word)
		dictHashLoose[standardizedWordLoose] = append(dictHashLoose[standardizedWordLoose], word)
		dictHashStrictReef[strictReef] = append(dictHashStrictReef[strictReef], word)
		dictHashStrict[standardizedWord] = append(dictHashStrict[standardizedWord], word)

		//find words with multiple IPAs
		if strings.Contains(word.IPA, " or ") {
			multiIPA += word.Navi + " "
			secondTerm := romanizeSecondIPA(word.IPA)
			if secondTerm != standardizedWord {
				dictHashLoose[dialectCrunch([]string{secondTerm}, true, true)[0]] = append(dictHashLoose[dialectCrunch([]string{secondTerm}, true, true)[0]], word)
				dictHashStrictReef[dialectCrunch([]string{secondTerm}, true, true)[0]] = append(dictHashStrictReef[dialectCrunch([]string{secondTerm}, true, true)[0]], word)
				dictHashStrict[secondTerm] = append(dictHashStrict[secondTerm], word)
			}
		}

		// See whether or not it violates normal phonotactic rules like Jakesully or Oìsss
		valid := true
		for a := range strings.SplitSeq(IsValidNavi(word.Navi, "en", false), "\n") {
			// Check every word.  If one of them isn't good, write down the word
			if len(a) > 0 && (!strings.Contains(a, "Valid:") || strings.Contains(a, "reef")) {
				valid = false
				break
			}
		}
		if !valid {
			oddballs += word.Navi + " "
		}

		return nil
	}

	var err error

	if mysql {
		if err = runOnDB(f); err != nil {
			uncacheHashDict()
			return err
		}
	}

	if err = runOnFile(f); err != nil {
		uncacheHashDict()
		return err
	}

	// Reverse the order to make accidental and new homonyms easier to see
	// Also make it a string for easier searching
	i := len(tempHoms)
	for i > 0 {
		i--
		homonyms += tempHoms[i] + " "
	}

	homonyms = strings.TrimSuffix(homonyms, " ")

	dictHashCached = true

	return nil
}

// searchTerms turns a definition into its searchable terms
func searchTerms(input string, excludeParen bool) []string {
	badChars := `~@#$%^&*()[]{}<>_/.,;:!?|+\"„“”«»`

	input = strings.ReplaceAll(input, "(", " (")

	// remove anything in parentheses to avoid clogging search results
	var tempString strings.Builder
	parenthesis := false
	for _, c := range input {
		if excludeParen {
			if c == '(' {
				parenthesis = true
			} else if c == ')' {
				parenthesis = false
				continue
			}
		}

		if !parenthesis {
			tempString.WriteString(string(c))
		}
	}
	input = tempString.String()

	// remove all the sketchy chars from arguments
	for _, c := range badChars {
		input = strings.ReplaceAll(input, string(c), "")
	}

	// normalize tìftang character
	input = strings.ReplaceAll(input, "’", "'")
	input = strings.ReplaceAll(input, "‘", "'")

	// find everything lowercase
	input = strings.ToLower(input)

	return strings.Split(input, " ")
}

// assignWord is a helper function for cacheDictHash2
func assignWord(wordMap map[string][]string, natlangWords string, naviWord string, excludeParen bool) (result map[string][]string) {
	newWords := searchTerms(natlangWords, excludeParen)

	for i := range newWords {
		duplicate := false
		for j := 0; j < len(wordMap[newWords[i]]); j++ {
			if wordMap[newWords[i]][j] == naviWord {
				duplicate = true
				break
			}
		}
		if !duplicate {
			wordMap[newWords[i]] = append(wordMap[newWords[i]], naviWord)
		}
	}
	return wordMap
}

// cacheDictHash2 caches Natural languages to Na'vi
func cacheDictHash2() error {
	err := cacheDictHash2Orig(true)
	if err == nil {
		fmt.Println("cache 2 loaded")
	} else {
		err = cacheDictHash2Orig(false)
		//fmt.Println("cache 2 loaded (File)")
	}
	return err
}

func cacheDictHash2Orig(mysql bool) error {
	// dont run if already is cached
	if len(dictHash2.EN) != 0 {
		return nil
	}

	for _, dictHash2Map := range dictHash2MapPtrs {
		*dictHash2Map = make(map[string][]string)
	}

	for _, dictHash2ParenthesisMap := range dictHash2ParenthesisMapPtrs {
		*dictHash2ParenthesisMap = make(map[string][]string)
	}

	// Set up the whole thing

	var setUpTheWholeThing = func(word Word) error {
		standardizedWord := strings.ToLower(word.Navi)
		standardizedWord = strings.ReplaceAll(standardizedWord, "+", "")
		var definitions = []string{word.EN, word.DE, word.ES, word.ET, word.FR, word.HU, word.IT, word.KO, word.NL,
			word.PL, word.PT, word.RU, word.SV, word.TR, word.UK}
		for i, definition := range definitions {
			if !hasNullDef(definition) {
				*dictHash2MapPtrs[i] = assignWord(*dictHash2MapPtrs[i], definition, standardizedWord, true)
				*dictHash2ParenthesisMapPtrs[i] = assignWord(*dictHash2ParenthesisMapPtrs[i], definition, standardizedWord, false)
			}
		}
		return nil
	}

	var err error
	if mysql {
		err = runOnDB(setUpTheWholeThing)
		if err != nil {
			uncacheHashDict2()
			return err
		}
	} else {
		err = runOnFile(setUpTheWholeThing)
		if err != nil {
			log.Printf("Error caching dictionary: %s", err)
			uncacheHashDict2()
			return err
		}
	}

	dictHash2Cached = true

	return nil
}

func uncacheHashDict() {
	dictHashCached = false
	dictHashLoose = nil
	dictHashStrict = nil
	homonyms = ""
	oddballs = ""
}

func uncacheHashDict2() {
	for _, dictHash2Map := range dictHash2MapPtrs {
		*dictHash2Map = nil
	}
	for _, dictHash2ParenthesisMap := range dictHash2ParenthesisMapPtrs {
		*dictHash2ParenthesisMap = nil
	}
	dictHash2Cached = false
}

func runOnDB(f func(word Word) error) error {
	user := os.Getenv("FW_USER")
	pass := os.Getenv("FW_PASS")
	host := os.Getenv("FW_HOST")
	name := os.Getenv("FW_DB")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, name)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return FailedToOpenDatabase.wrap(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(FailedToCloseDatabase.wrap(err).Error())
		}
	}(db)

	rows, err1 := db.Query("SELECT " +
		"m.id, m.navi, m.ipa, m.infixes, m.partOfSpeech, s.source, b.stressed, b.syllables, b.infixDots, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'de') AS de, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'en') AS en, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'es') AS es, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'et') AS et, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'fr') AS fr, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'hu') AS hu, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'it') AS it, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'ko') AS ko, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'nl') AS nl, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'pl') AS pl, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'pt') AS pt, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'ru') AS ru, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'sv') AS sv, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'tr') AS tr, " +
		"(SELECT localized FROM fwedit_localizedWords AS l WHERE l.id = m.id AND languageCode = 'uk') AS uk " +
		"FROM fwedit_metaWords AS m " +
		"INNER JOIN fwedit_sources AS s ON (m.id = s.id) " +
		"INNER JOIN fwedit_breakdown AS b ON (s.id = b.id)")

	if err1 != nil {
		return err1
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(FailedToCloseDatabase.wrap(err).Error())
		}
	}(rows)

	var w Word
	var de, en, es, et, fr, hu, it, ko, nl, pl, pt, ru, sv, tr, uk []byte

	for rows.Next() {
		err = rows.Scan(&w.ID, &w.Navi, &w.IPA, &w.InfixLocations, &w.PartOfSpeech, &w.Source, &w.Stressed,
			&w.Syllables, &w.InfixDots, &de, &en, &es, &et, &fr, &hu, &it, &ko, &nl, &pl, &pt, &ru, &sv, &tr, &uk)

		if err != nil {
			return err
		}

		w.DE = string(de)
		w.EN = string(en)
		w.ES = string(es)
		w.ET = string(et)
		w.FR = string(fr)
		w.HU = string(hu)
		w.IT = string(it)
		w.KO = string(ko)
		w.NL = string(nl)
		w.PL = string(pl)
		w.PT = string(pt)
		w.RU = string(ru)
		w.SV = string(sv)
		w.TR = string(tr)
		w.UK = string(uk)

		err = f(w)

		if err != nil {
			return err
		}
	}

	return nil
}

func runOnFile(f func(word Word) error) error {
	dictionaryFile := findDictionaryFile()
	if dictionaryFile == "" {
		return NoDictionary
	}

	file, err := os.Open(dictionaryFile)
	if err != nil {
		log.Println(FailedToOpenDictFile.wrap(err).Error())
		return FailedToOpenDictFile.wrap(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(FailedToCloseDictFile.wrap(err).Error())
		}
	}(file)

	scanner := bufio.NewScanner(file)

	var first = true
	var pos dictPos
	for scanner.Scan() {
		// get a single line out of the dict
		line := scanner.Text()

		// Split line at \t so we get all information
		fields := strings.Split(line, "\t")

		// When first, then this is the header
		if first {
			pos = readDictPos(fields)
			first = false
		} else {
			// Put the stuff from fields into the Word struct
			err = f(newWord(fields, pos))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetFullDict() (allWords []Word, err error) {
	// No need for the lock because only List() calls it
	if dictionaryCached {
		firstWordID, _ := strconv.Atoi(dictionary[0].ID)
		if firstWordID > 100 {
			slices.SortFunc(dictionary, func(a, b Word) int {
				a1, _ := strconv.Atoi(a.ID)
				b1, _ := strconv.Atoi(b.ID)
				return a1 - b1
			})
		}
		allWords = dictionary
	} else {
		err = runOnFile(func(word Word) error {
			allWords = append(allWords, word)
			return nil
		})
		return
	}
	return
}

// GetDictSizeSimple returns dictionary size as a number
func GetDictSizeSimple() (count int) {
	universalLock.Lock()
	defer universalLock.Unlock()
	return len(dictionary)
}

// GetDictSize returns dictionary size as a complete sentence
func GetDictSize(lang string) (count string, err error) {
	universalLock.Lock()
	defer universalLock.Unlock()
	// Count words
	amount := 0
	if dictionaryCached {
		amount = len(dictionary)
	} else {
		err = runOnFile(func(word Word) error {
			amount++
			return nil
		})
	}

	// Put the word count into a complete sentence
	count = strconv.Itoa(amount)
	count = getDictSizeMessage(lang, count)

	return
}

// UpdateDict updates the dictionary file.
// universalLock will hopefully prevent anything from accessing
// the dict while updating
func UpdateDict() error {
	universalLock.Lock()
	defer universalLock.Unlock()
	err := downloadDict("")
	if err != nil {
		return FailedToDownload.wrap(err)
	}

	err = cacheDict()
	if err != nil {
		log.Println(FailedToCache.wrap(err))
		return FailedToCache.wrap(err)
	}

	if dictHashCached {
		uncacheHashDict()
	}

	err = cacheDictHash()
	if err != nil {
		log.Println(FailedToCache.wrap(err))
		return FailedToCache.wrap(err)
	}

	if dictHash2Cached {
		uncacheHashDict2()
	}

	err = cacheDictHash2()
	if err != nil {
		log.Println(FailedToCache.wrap(err))
		return FailedToCache.wrap(err)
	}

	return nil
}

// AssureDict will assure that the dictionary exists.
// If no dictionary is found, it will be downloaded next to the executable.
func AssureDict() error {
	// check if dict already exists
	file := findDictionaryFile()
	if file != "" {
		return nil
	}

	// if it doesn't, put it in ~/.fwew/
	path := filepath.Join(Text("dataDir"), dictFileName)

	err := downloadDict(path)
	if err != nil {
		return err
	}

	return nil
}

// StartEverything connects to the dictionary, loads the cache maps, and generates phoneme statistics. Returns a status
// including how long everything took.
func StartEverything() string {
	universalLock.Lock()
	start := time.Now()
	var errors = []error{
		AssureDict(),
		cacheDict(),
		cacheDictHash(),
		cacheDictHash2(),
	}
	for _, err := range errors {
		if err != nil {
			log.Println(err)
		}
	}
	universalLock.Unlock()
	PhonemeDistros()
	elapsed := strconv.FormatFloat(time.Since(start).Seconds(), 'f', -1, 64)
	return fmt.Sprintln("Caches loaded,  took " + elapsed + " seconds")
}

// StopEverything clears the caches and returns a status string.
func StopEverything() string {
	universalLock.Lock()
	uncacheDict()
	uncacheHashDict()
	uncacheHashDict2()
	universalLock.Unlock()
	return fmt.Sprintln("Caches cleared")
}
