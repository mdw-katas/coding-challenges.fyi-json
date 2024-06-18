package printing

import (
	"bytes"
	"io"

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

func (this *pretty) nested() bool {
	if len(this.state) == 0 {
		return false
	}
	last := this.state[len(this.state)-1]
	return last == lexing.TokenArrayStart
}

func (this *pretty) Print(token lexing.Token) {
	switch token.Type {
	case lexing.TokenArrayStart, lexing.TokenObjectStart:
		if this.nested() {
			this.indent()
		}
		this.write(token.Value)
		this.state = append(this.state, token.Type)
		this.items = append(this.items, 0)
		this.awaitingArrayValue = true
	case lexing.TokenArrayStop, lexing.TokenObjectStop:
		this.state = this.state[:len(this.state)-1]
		if this.items[len(this.items)-1] > 0 {
			this.indent()
		}
		this.items = this.items[:len(this.items)-1]
		this.write(token.Value)
	case lexing.TokenNull, lexing.TokenTrue, lexing.TokenFalse, lexing.TokenString, lexing.TokenNumber:
		if len(this.state) > 0 {
			this.items[len(this.items)-1]++
			if !this.awaitingObjectValue || this.awaitingArrayValue {
				this.indent()
			}
			this.write(token.Value)
		}
		this.awaitingObjectValue = false
		this.awaitingArrayValue = false
	case lexing.TokenComma, lexing.TokenIllegal:
		this.write(token.Value)
	case lexing.TokenColon:
		this.awaitingObjectValue = true
		this.write(token.Value)
		this.write(space)
	}
}

func (this *pretty) write(data []byte) {
	_, _ = this.out.Write(data)
}
func (this *pretty) indent() {
	this.write(newline)
	this.write(bytes.Repeat(indent, len(this.state)))
}

var (
	space   = []byte(" ")
	newline = []byte("\n")
	indent  = []byte("  ")
)
