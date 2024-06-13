package printing

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"
	"github.com/mdwhatcott/testing/should"
)

func TestVerbatimPrinter(t *testing.T) {
	out := &bytes.Buffer{}
	input := `{"a": [1,2,3 ],"b":"hi" }`
	printer := NewVerbatimPrinter(out)
	for token := range lexing.Lex(strings.NewReader(input)) {
		printer.Print(token)
	}
	should.So(t, out.String(), should.Equal, input)
}
