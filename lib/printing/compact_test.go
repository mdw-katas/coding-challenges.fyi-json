package printing

import (
	"bytes"
	"testing"

	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"
	"github.com/mdwhatcott/testing/should"
)

func TestCompactPrinter(t *testing.T) {
	out := &bytes.Buffer{}
	input := `{"a": [1,2,3 ],"b":"hi" }`
	expected := `{"a":[1,2,3],"b":"hi"}`
	printer := NewCompactPrinter(out)
	for token := range lexing.Lex([]byte(input)) {
		printer.Print(token)
	}
	should.So(t, out.String(), should.Equal, expected)
}
