package fwew_lib

import (
	"strings"
	"testing"
)

func Test_SingleNames(t *testing.T) {
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
