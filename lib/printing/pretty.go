package printing

import (
	"fmt"
	"io"
	"strings"

	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"
)

type pretty struct {
	out   io.Writer
	state []lexing.TokenType
	items []int

	awaitingArrayValue  bool
	awaitingObjectValue bool
}

func NewPrettyPrinter(out io.Writer) Printer {
	return &pretty{out: out}
}

func (this *pretty) Print(token lexing.Token) {
	switch token.Type {
	case lexing.TokenArrayStart, lexing.TokenObjectStart:
		_, _ = this.out.Write(token.Value)
		this.state = append(this.state, token.Type)
		this.items = append(this.items, 0)
		this.awaitingArrayValue = true
	case lexing.TokenArrayStop, lexing.TokenObjectStop:
		this.state = this.state[:len(this.state)-1]
		if this.items[len(this.items)-1] > 0 {
			_, _ = fmt.Fprintln(this.out)
			_, _ = io.WriteString(this.out, strings.Repeat("  ", len(this.state)))
		}
		this.items = this.items[:len(this.items)-1]
		_, _ = this.out.Write(token.Value)
	case lexing.TokenNull, lexing.TokenTrue, lexing.TokenFalse, lexing.TokenString, lexing.TokenNumber:
		if len(this.state) > 0 {
			this.items[len(this.items)-1]++
			if !this.awaitingObjectValue || this.awaitingArrayValue {
				_, _ = fmt.Fprintln(this.out)
				_, _ = io.WriteString(this.out, strings.Repeat("  ", len(this.state)))
			}
			_, _ = this.out.Write(token.Value)
		}
		this.awaitingObjectValue = false
		this.awaitingArrayValue = false
	case lexing.TokenComma, lexing.TokenIllegal:
		_, _ = this.out.Write(token.Value)
	case lexing.TokenColon:
		this.awaitingObjectValue = true
		_, _ = this.out.Write(token.Value)
		_, _ = io.WriteString(this.out, " ")
	}
}
