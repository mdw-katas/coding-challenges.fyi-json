package printing

import (
	"io"

	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"
)

type compact struct {
	out io.Writer
}

func NewCompactPrinter(out io.Writer) Printer {
	return &compact{out: out}
}

func (this *compact) Print(token lexing.Token) {
	if token.Type == lexing.TokenWhitespace {
		return
	}
	_, _ = this.out.Write(token.Value)
}
