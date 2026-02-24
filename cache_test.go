package fwew_lib

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func Test_cacheDict(t *testing.T) {
	word := Word{
		ID:   "4",
		Navi: "'ampi",
	}

	err := cacheDictHash()
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

func TestGetDictSizeSimple(t *testing.T) {
	size := GetDictSizeSimple()
	expected := len(dictionary)
	if size != expected {
		t.Errorf("Dictionary size mismatch: %d != %d", size, expected)
	}
}

func TestGetDictSize(t *testing.T) {
	size, err := GetDictSize("en")
	expected := len(dictionary)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(size, strconv.Itoa(expected)) {
		t.Errorf("Dictionary size mismatch: %s != %d", size, expected)
	}
}

func Test_UpdateDict(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		content, err := os.ReadFile(filepath.Join(Text("dataDir"), dictFileName))
		if err != nil {
			t.Fatalf("failed to read downloaded file: %v", err)
		}
		_, err2 := w.Write(content)
		if err2 != nil {
			t.Fatalf("failed to write downloaded file: %v", err2)
		}
	}))
	defer server.Close()

	// Test your function with the mock server URL
	prevURL := Text("dictURL")
	texts["dictURL"] = server.URL
	err := UpdateDict()
	if err != nil {
		t.Fatalf("downloadFile failed: %v", err)
	}
	texts["dictURL"] = prevURL

	// Verify the file was created with the correct content
	content, err := os.ReadFile(filepath.Join(Text("dataDir"), dictFileName))
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}

	if len(content) < 100_000 {
		t.Fatalf("dictionary file smaller than expected")
	}
}

func Test_StartEverything(t *testing.T) {
	status := StartEverything()
	status = strings.TrimSpace(status)
	if !strings.HasPrefix(status, "Caches loaded") {
		t.Errorf("[%s]", status)
	}
	if !strings.HasSuffix(status, "seconds") {
		t.Errorf("[%s]", status)
	}
}

func Test_StopEverything(t *testing.T) {
	status := StopEverything()
	status = strings.TrimSpace(status)
	if status != "Caches cleared" {
		t.Errorf("[%s]", status)
	}
}
