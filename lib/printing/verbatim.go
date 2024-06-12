package printing

import (
	"io"

	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"
)

type verbatim struct {
	out io.Writer
}

func NewVerbatimPrinter(out io.Writer) Printer {
	return &verbatim{out: out}
}

func (this *verbatim) Print(token lexing.Token) {
	_, _ = this.out.Write(token.Value)
}
