package printing

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"
	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/util/git"
	"github.com/mdwhatcott/testing/should"
)

func TestPrettyPrinter(t *testing.T) {
	input, err := os.ReadFile(filepath.Join(git.RootDirectory(), "lib", "lexing", "testdata", "pass1.json"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join(git.RootDirectory(), "lib", "printing", "pretty_test_expected.txt"))
	if err != nil {
		t.Fatal(err)
	}
	out := &bytes.Buffer{}
	printer := NewPrettyPrinter(out)
	for token := range lexing.Lex(bytes.NewReader(input)) {
		printer.Print(token)
	}
	fmt.Println(out.String())
	should.So(t, out.String(), should.Equal, string(expected))
}
