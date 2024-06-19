package printing

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"
	"github.com/mdwhatcott/testing/should"
)

func TestColorsSuite(t *testing.T) {
	should.Run(&ColorsSuite{T: should.New(t)}, should.Options.UnitTests())
}

type ColorsSuite struct {
	*should.T
}

func (this *ColorsSuite) Test() {
	var out bytes.Buffer
	inner := NewVerbatimPrinter(&out)
	outer := NewColorPrinter(&out, inner)
	input := `{"a": [1,2,3,null,true,false ],"b":"hi" }asdf`
	for token := range lexing.Lex(strings.NewReader(input)) {
		outer.Print(token)
	}
	this.So(out.String(), should.Equal,
		"\x1b[36m{\x1b[0m\x1b[34m\"a\"\x1b[0m\x1b[36m:\x1b[0m \x1b[36m[\x1b[0m\x1b[33m1\x1b[0m\x1b[36m,\x1b[0m\x1b[33m2\x1b[0m\x1b[36m,\x1b[0m\x1b[33m3\x1b[0m\x1b[36m,\x1b[0m\x1b[37mnull\x1b[0m\x1b[36m,\x1b[0m\x1b[32mtrue\x1b[0m\x1b[36m,\x1b[0m\x1b[35mfalse\x1b[0m \x1b[36m]\x1b[0m\x1b[36m,\x1b[0m\x1b[34m\"b\"\x1b[0m\x1b[36m:\x1b[0m\x1b[34m\"hi\"\x1b[0m \x1b[36m}\x1b[0m\x1b[31masdf\x1b[0m",
	)
}
