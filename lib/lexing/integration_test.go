package lexing

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mdwhatcott/testing/should"
)

func TestIntegration(t *testing.T) {
	_, here, _, _ := runtime.Caller(0)
	hereDir := filepath.Dir(here)
	testdata := filepath.Join(hereDir, "testdata")
	listing, err := os.ReadDir(testdata)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range listing {
		t.Run(entry.Name(), func(t *testing.T) {
			file, err := os.Open(filepath.Join(testdata, entry.Name()))
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = file.Close() }()

			should.So(t, isValid(file), should.Equal, strings.HasPrefix(entry.Name(), "pass"))
		})
	}
}

func isValid(input io.Reader) bool {
	var last Token
	for token := range Lex(input) {
		last = token
	}
	return last.Type != "" && last.Type != TokenIllegal
}
