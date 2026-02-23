package fwew_lib

import (
	"log"
	"strings"
	"testing"
)

func setupTest() func(t *testing.T) {
	log.Println("setup test")
	_ = StartEverything()

	// Return a function to tear down the test
	return func(t *testing.T) {
		log.Println("teardown test")
		uncacheDict()
		uncacheHashDict()
		uncacheHashDict2()
	}
}

func Test_SingleNames(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var (
		nameCount     = 3
		dialect       = 0
		syllableCount = 2
	)
	names := SingleNames(nameCount, dialect, syllableCount)
	names = strings.TrimSpace(names)
	nameList := strings.Split(names, "\n")
	if len(nameList) != 3 {
		t.Errorf("Expected %d names, got %d", nameCount, len(names))
	}
}

func Test_FullNames(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var (
		ending        = "'ite"
		nameCount     = 8
		dialect       = 1
		syllableCount = [3]int{2, 3, 2}
	)
	names := FullNames(ending, nameCount, dialect, syllableCount, false)
	names = strings.TrimSpace(names)
	nameList := strings.Split(names, "\n")
	if len(nameList) != 8 {
		t.Errorf("Expected %d names, got %d", nameCount, len(names))
	}
}

func Test_NameAlu(t *testing.T) {
	teardownTest := setupTest()
	defer teardownTest(t)

	var (
		nameCount     = 2
		dialect       = 1
		syllableCount = 2
		nounMode      = 1
		adjMode       = 1
	)
	names := NameAlu(nameCount, dialect, syllableCount, nounMode, adjMode)
	names = strings.TrimSpace(names)
	nameList := strings.Split(names, "\n")
	if len(nameList) != 2 {
		t.Errorf("Expected %d names, got %d", nameCount, len(names))
	}
}
